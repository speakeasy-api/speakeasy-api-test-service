package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/lingrino/go-fault"
)

type FaultSession struct {
	// RequestCount is the number of requests that have been made on this session.
	RequestCount int
	// Exhausted is set to true when all faults on this session have been exercised.
	Exhausted bool
	Settings  FaultSettings
}

// Describes the fault injection settings for a session. The fault chain is
// a series of fault injectors that are applied to the request in order. The
// order of faults is:
// - Delay
// - ConnectionClose
// - ConnectionReset
// - Reject
// - Error
type FaultSettings struct {
	// Number of times to close the connection.
	ConnectionCloseCount int `json:"connection_close_count"`

	// ConnectionResetCount is the number of times to reset the connection.
	ConnectionResetCount int `json:"connection_reset_count"`

	// DelayMS is the number of milliseconds to delay the request.
	DelayMS int64 `json:"delay_ms"`

	// DelayCount is the number of times to delay the request.
	DelayCount int `json:"delay_count"`

	// RejectCount is the number of times to reject the request without a response.
	// A value greater than 0 enables this fault injector.
	RejectCount int `json:"reject_count"`

	// ErrorCount is the number of times to return an error status code.
	// A value greater than 0 enables this fault injector.
	//
	// NOTE: Error injection only takes effect after all rejections have
	// resolved if both of these injectors are enabled.
	ErrorCount int `json:"error_count"`

	// ErrorCode is the status code to return when the error injector is enabled.
	ErrorCode int `json:"error_code"`
}

func Fault(h http.Handler) http.Handler {
	var sessions sync.Map
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqid := r.Header.Get("request-id")
		if reqid == "" {
			h.ServeHTTP(w, r)
			return
		}

		settinghdr := r.Header.Get("fault-settings")
		if settinghdr == "" {
			h.ServeHTTP(w, r)
			return
		}

		session := &FaultSession{}
		asession, found := sessions.Load(reqid)
		if found {
			session = asession.(*FaultSession)
			if session.Exhausted {
				h.ServeHTTP(w, r)
				return
			}
		}

		var settings FaultSettings
		err := json.Unmarshal([]byte(settinghdr), &settings)
		if err != nil {
			http.Error(w, "Invalid fault settings", http.StatusBadRequest)
			return
		}

		reqCount := session.RequestCount
		session.Settings = settings
		session.RequestCount++
		defer func() {
			sessions.Store(reqid, session)
		}()

		var faults []fault.Injector

		// Since multiple injectors can be enabled, need to count the number of
		// requests based on prior injector counts
		countOffset := 0

		if settings.DelayMS > 0 && reqCount < settings.DelayCount {
			inj, err := fault.NewSlowInjector(time.Millisecond * time.Duration(settings.DelayMS))
			if err != nil {
				http.Error(w, "Failed to build slow injector", http.StatusInternalServerError)
				return
			}

			faults = append(faults, inj)
		}

		// Delay injector does not increase the count offset.

		if settings.ConnectionCloseCount > 0 && reqCount < settings.ConnectionCloseCount+countOffset {
			faults = append(faults, &ConnectionErrorInjector{})
		}

		countOffset += settings.ConnectionCloseCount

		if settings.ConnectionResetCount > 0 && reqCount < settings.ConnectionResetCount+countOffset {
			faults = append(faults, &ConnectionErrorInjector{Reset: true})
		}

		countOffset += settings.ConnectionResetCount

		if settings.RejectCount > 0 && reqCount < settings.RejectCount+countOffset {
			inj, err := fault.NewRejectInjector()
			if err != nil {
				http.Error(w, "Failed to build reject injector", http.StatusInternalServerError)
				return
			}

			faults = append(faults, inj)
		}

		countOffset += settings.RejectCount

		if settings.ErrorCode > 0 && reqCount < (settings.ErrorCount+countOffset) {
			inj, err := fault.NewErrorInjector(settings.ErrorCode, fault.WithStatusText("Injected error"))
			if err != nil {
				http.Error(w, "Failed to build error injector", http.StatusInternalServerError)
				return
			}

			faults = append(faults, inj)
		}

		if len(faults) == 0 {
			session.Exhausted = true
			h.ServeHTTP(w, r)
			return
		}

		faultchain, err := fault.NewChainInjector(faults)
		if err != nil {
			http.Error(w, "Failed to build fault chain injector", http.StatusInternalServerError)
		}

		w.Header().Set("Faults-Enabled", "true")
		faultchain.Handler(h).ServeHTTP(w, r)
	})
}
