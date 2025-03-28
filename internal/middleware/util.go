package middleware

import (
	"reflect"
	"unsafe"

	"github.com/Bofry/worker-redis/internal"
)

func isMessageHandler(rv reflect.Value) bool {
	return internal.IsMessageHandler(rv)
}

func asMessageHandler(rv reflect.Value) internal.MessageHandler {
	return internal.AsMessageHandler(rv)
}

func isMessageObserver(rv reflect.Value) bool {
	return internal.IsMessageObserver(rv)
}

func asMessageObserver(rv reflect.Value) internal.MessageObserver {
	return internal.AsMessageObserver(rv)
}

func isMessageFilterAffinity(rv reflect.Value) bool {
	return internal.IsMessageFilterAffinity(rv)
}

func asMessageFilterAffinity(rv reflect.Value) internal.MessageFilterAffinity {
	return internal.AsMessageFilterAffinity(rv)
}

func asRedisWorker(rv reflect.Value) *internal.RedisWorker {
	return reflect.NewAt(typeOfHost, unsafe.Pointer(rv.Pointer())).
		Interface().(*internal.RedisWorker)
}
