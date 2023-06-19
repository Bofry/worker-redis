package middleware

import "github.com/Bofry/trace"

type EventEvidence struct {
	traceID trace.TraceID
	spanID  trace.SpanID
	stream  string
}

func (e EventEvidence) ProcessingTraceID() trace.TraceID {
	return e.traceID
}

func (e EventEvidence) ProcessingSpanID() trace.SpanID {
	return e.spanID
}

func (e EventEvidence) Stream() string {
	return e.stream
}
