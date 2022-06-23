package redis

import redis "github.com/Bofry/lib-redis-stream"

func NewAdminClient(opt *UniversalOptions) (*AdminClient, error) {
	return redis.NewAdminClient(opt)
}

func NewForwarder(opt *UniversalOptions) (*Forwarder, error) {
	return redis.NewForwarder(opt)
}
