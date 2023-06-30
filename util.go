package redis

import "github.com/Bofry/worker-redis/internal"

func ContextBuilder() *internal.ContextBuilder {
	return internal.NewContextBuilder()
}
