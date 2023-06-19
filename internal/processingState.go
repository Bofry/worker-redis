package internal

import "github.com/Bofry/trace"

type ProcessingState struct {
	Tracer *trace.SeverityTracer
	Span   *trace.SeveritySpan
	Stream string
}
