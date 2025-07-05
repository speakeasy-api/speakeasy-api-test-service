package eventstreams

import (
	"fmt"
	"net/http"
	"time"
)

func pushChunks(rw http.ResponseWriter, chunks []string) {
	for _, chunk := range chunks {
		fmt.Fprintln(rw, chunk)

		if f, ok := rw.(http.Flusher); ok {
			f.Flush()
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func pushEvents(rw http.ResponseWriter, events [][]string) {
	for _, event := range events {
		for _, line := range event {
			fmt.Fprintln(rw, line)
		}
		fmt.Fprintln(rw, "")

		if f, ok := rw.(http.Flusher); ok {
			f.Flush()
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func HandleEventStreamJSON(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Add("Content-Type", "text/event-stream")

	pushEvents(rw, [][]string{
		{
			`data: {"content": "Hello"}`,
		},

		{
			`data: {"content": " "}`,
		},

		{
			`data: {"content": "world"}`,
		},

		{
			`data: {"content": "!"}`,
		},
	})
}

func HandleEventStreamText(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Add("Content-Type", "text/event-stream")

	pushEvents(rw, [][]string{
		{
			`data: Hello`,
		},

		{
			`data:  `,
		},

		{
			`data: world`,
		},

		{
			`data: !`,
		},
	})
}

func HandleEventStreamMultiLine(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Add("Content-Type", "text/event-stream")

	pushEvents(rw, [][]string{
		{
			`data: YHOO`,
			`data: +2`,
			`data: 10`,
		},
	})
}

func HandleEventStreamRich(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Add("Content-Type", "text/event-stream")

	pushEvents(rw, [][]string{
		{
			`id: job-1`,
			`event: completion`,
			`data: {"completion": "Hello", "stop_reason": null, "model": "jeeves-1"}`,
		},

		{
			`event: heartbeat`,
			`data: ping`,
			`retry: 3000`,
		},

		{
			`id: job-1`,
			`event: completion`,
			`data: {"completion": "world!", "stop_reason": "stop_sequence", "model": "jeeves-1"}`,
		},
	})
}

func HandleEventStreamChat(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Add("Content-Type", "text/event-stream")

	pushEvents(rw, [][]string{
		{
			`data: {"content": "Hello"}`,
		},

		{
			`data: {"content": " "}`,
		},

		{
			`data: {"content": "world"}`,
		},

		{
			`data: {"content": "!"}`,
		},

		{
			`data: [DONE]`,
		},
	})
}

func HandleEventStreamChatFlatten(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Add("Content-Type", "text/event-stream")

	pushEvents(rw, [][]string{
		{
			`data: {"content": "Hello"}`,
		},

		{
			`data: {"content": " "}`,
		},

		{
			`data: {"content": "world"}`,
		},

		{
			`data: {"content": "!"}`,
		},

		{
			`data: [DONE]`,
		},
	})
}

func HandleEventStreamChatChunked(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Add("Content-Type", "text/event-stream")

	pushChunks(rw, []string{
		"data: {\"content\": ",
		"\"Hello\"}\n\ndata: {\"content\": \" \"}",
		"data: {\"content\": \"world\"}",
		"data: {\"content\": \"!\"}\n\ndata: [DONE]\n",
		"\ndata: {\"content\": \"Post sentinel data\"}\n\n",
	})
}

func HandleEventStreamDifferentDataSchemas(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Add("Content-Type", "text/event-stream")

	pushEvents(rw, [][]string{
		{
			`id: event-1`,
			`event: message`,
			`data: {"content": "Here is your url"}`,
		},

		{
			`id: event-2`,
			`event: url`,
			`data: {"url": "https://example.com"}`,
		},

		{
			`id: event-3`,
			`event: message`,
			`data: {"content": "Have a great day!"}`,
		},
	})
}

func HandleEventStreamDifferentDataSchemasFlatten(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Add("Content-Type", "text/event-stream")

	pushEvents(rw, [][]string{
		{
			`id: event-1`,
			`event: message`,
			`data: {"content": "Here is your url"}`,
		},

		{
			`id: event-2`,
			`event: url`,
			`data: {"url": "https://example.com"}`,
		},

		{
			`id: event-3`,
			`event: message`,
			`data: {"content": "Have a great day!"}`,
		},
	})
}

func HandleEventStreamStayOpen(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "text/event-stream")
	rw.Header().Add("Cache-Control", "no-cache")
	rw.Header().Add("Connection", "keep-alive")

	// Send events 1, 2, 3 immediately
	fmt.Fprintln(rw, "data: event 1")
	fmt.Fprintln(rw, "")
	fmt.Fprintln(rw, "data: event 2")
	fmt.Fprintln(rw, "")
	fmt.Fprintln(rw, "data: event 3")
	fmt.Fprintln(rw, "")
	
	if f, ok := rw.(http.Flusher); ok {
		f.Flush()
	}

	// Wait 100ms then send event 4
	time.Sleep(100 * time.Millisecond)
	fmt.Fprintln(rw, "data: event 4")
	fmt.Fprintln(rw, "")
	
	if f, ok := rw.(http.Flusher); ok {
		f.Flush()
	}

	// Wait another 100ms then send sentinel event
	time.Sleep(100 * time.Millisecond)
	fmt.Fprintln(rw, "data: [SENTINEL]")
	fmt.Fprintln(rw, "")
	
	if f, ok := rw.(http.Flusher); ok {
		f.Flush()
	}

	// Keep the connection open until client closes
	// Monitor the request context to detect when client disconnects
	<-r.Context().Done()
}
