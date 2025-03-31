package test

import (
	"context"

	"github.com/Bofry/worker-redis/tracing"
)

type TestStreamMessageCounter struct {
	MessageCount        int
	SuccessMessageCount int
	InvalidMessageCount int
	PanicCount          int
}

func (c *TestStreamMessageCounter) IncreaseMessageCount(ctx context.Context) int {
	tr := tracing.GetTracer(c)
	sp := tr.Start(ctx, "IncreaseMessageCount()")
	defer sp.End()

	c.MessageCount++
	return c.MessageCount
}

func (c *TestStreamMessageCounter) IncreaseSuccessMessageCount(ctx context.Context) int {
	tr := tracing.GetTracer(c)
	sp := tr.Start(ctx, "IncreaseSuccessMessageCount()")
	defer sp.End()

	c.SuccessMessageCount++
	return c.SuccessMessageCount
}

func (c *TestStreamMessageCounter) IncreaseInvalidMessageCount(ctx context.Context) int {
	tr := tracing.GetTracer(c)
	sp := tr.Start(ctx, "IncreaseInvalidMessageCount()")
	defer sp.End()

	c.InvalidMessageCount++
	return c.InvalidMessageCount
}

func (c *TestStreamMessageCounter) IncreasePanicCount(ctx context.Context) int {
	tr := tracing.GetTracer(c)
	sp := tr.Start(ctx, "IncreasePanicCount()")
	defer sp.End()

	c.PanicCount++
	return c.PanicCount
}
