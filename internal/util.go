package internal

import (
	"reflect"
	"strings"
)

const (
	REDIS_BUSYGROUP_PREFIX = "BUSYGROUP"
)

func IsMessageHandlerType(rt reflect.Type) bool {
	return rt.AssignableTo(typeOfMessageHandler)
}

func IsMessageHandler(rv reflect.Value) bool {
	if rv.IsValid() {
		return rv.Type().AssignableTo(typeOfMessageHandler)
	}
	return false
}

func AsMessageHandler(rv reflect.Value) MessageHandler {
	if rv.IsValid() {
		if v, ok := rv.Convert(typeOfMessageHandler).Interface().(MessageHandler); ok {
			return v
		}
	}
	return nil
}

func isRedisBusyGroupError(err error) bool {
	return strings.HasPrefix(err.Error(), REDIS_BUSYGROUP_PREFIX)
}

func isMessageObserverAffair(rv reflect.Value) bool {
	if rv.IsValid() {
		return rv.Type().AssignableTo(typeOfMessageObserverAffair)
	}
	return false
}

func asMessageObserverAffair(rv reflect.Value) MessageObserverAffair {
	if rv.IsValid() {
		if v, ok := rv.Convert(typeOfMessageObserverAffair).Interface().(MessageObserverAffair); ok {
			return v
		}
	}
	return nil
}
