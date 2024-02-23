package pagination

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
)

type LimitOffsetRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Page   int `json:"page"`
}

type CursorRequest struct {
	Cursor int `json:"cursor"`
}

type NonNumericCursorRequest struct {
	Cursor string `json:"cursor"`
}

type PaginationResponse struct {
	NumPages    int           `json:"numPages"`
	ResultArray []interface{} `json:"resultArray"`
}

// Insecure reversable hashing for string cursors
func hash(s string) (int, error) {
	return strconv.Atoi(s)
}
func unhash(h int) string {
	return strconv.Itoa(h)
}

const total = 20

func HandleLimitOffsetPage(w http.ResponseWriter, r *http.Request) {
	queryLimit := r.FormValue("limit")
	queryPage := r.FormValue("page")

	var pagination LimitOffsetRequest
	hasBody := true
	if err := json.NewDecoder(r.Body).Decode(&pagination); err != nil {
		hasBody = false
	}
	limit := getValue(queryLimit, hasBody, pagination.Limit)
	if limit == 0 {
		limit = 20
	}
	page := getValue(queryPage, hasBody, pagination.Page)

	start := (page - 1) * limit

	res := PaginationResponse{
		NumPages:    int(math.Ceil(float64(total) / float64(limit))),
		ResultArray: make([]interface{}, 0),
	}

	for i := start; i < total && len(res.ResultArray) < limit; i++ {
		res.ResultArray = append(res.ResultArray, i)
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(500)
	}
}

func HandleLimitOffsetOffset(w http.ResponseWriter, r *http.Request) {
	queryLimit := r.FormValue("limit")
	queryOffset := r.FormValue("offset")

	var pagination LimitOffsetRequest
	hasBody := true
	if err := json.NewDecoder(r.Body).Decode(&pagination); err != nil {
		hasBody = false
	}

	limit := getValue(queryLimit, hasBody, pagination.Limit)
	if limit == 0 {
		limit = 20
	}
	offset := getValue(queryOffset, hasBody, pagination.Offset)

	res := PaginationResponse{
		NumPages:    int(math.Ceil(float64(total) / float64(limit))),
		ResultArray: make([]interface{}, 0),
	}

	for i := offset; i < total && len(res.ResultArray) < limit; i++ {
		res.ResultArray = append(res.ResultArray, i)
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(500)
	}
}

func HandleCursor(w http.ResponseWriter, r *http.Request) {
	queryCursor := r.FormValue("cursor")

	var pagination CursorRequest
	hasBody := true
	if err := json.NewDecoder(r.Body).Decode(&pagination); err != nil {
		hasBody = false
	}

	cursor := getValue(queryCursor, hasBody, pagination.Cursor)

	res := PaginationResponse{
		NumPages:    0,
		ResultArray: make([]interface{}, 0),
	}

	for i := cursor + 1; i < total && len(res.ResultArray) < 15; i++ {
		res.ResultArray = append(res.ResultArray, i)
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(500)
	}
}

func HandleNonNumericCursor(w http.ResponseWriter, r *http.Request) {
	queryCursor := r.FormValue("cursor")
	var pagination NonNumericCursorRequest
	hasBody := true
	if err := json.NewDecoder(r.Body).Decode(&pagination); err != nil {
		hasBody = false
	}

	cursor := getNonNumericValue(queryCursor, hasBody, pagination.Cursor)

	res := PaginationResponse{
		NumPages:    0,
		ResultArray: make([]interface{}, 0),
	}
	var cursorI, _ = hash(cursor)
	for i := cursorI + 1; i < total && len(res.ResultArray) < 15; i++ {
		res.ResultArray = append(res.ResultArray, unhash(i))
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(500)
	}
}

func getValue(queryValue string, hasBody bool, paginationValue int) int {
	if hasBody {
		return paginationValue
	} else {
		value, err := strconv.Atoi(queryValue)
		if err != nil {
			return 0
		}
		return value
	}
}

func getNonNumericValue(queryValue string, hasBody bool, paginationValue string) string {
	if hasBody {
		return paginationValue
	} else {
		return queryValue
	}
}
