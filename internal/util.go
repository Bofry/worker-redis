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

func IsMessageObserver(rv reflect.Value) bool {
	if rv.IsValid() {
		return rv.Type().AssignableTo(typeOfMessageObserver)
	}
	return false
}

func AsMessageObserver(rv reflect.Value) MessageObserver {
	if rv.IsValid() {
		if v, ok := rv.Convert(typeOfMessageObserver).Interface().(MessageObserver); ok {
			return v
		}
	}
	return nil
}

func IsMessageFilterAffinity(rv reflect.Value) bool {
	if rv.IsValid() {
		return rv.Type().AssignableTo(typeOfMessageFilterAffinity)
	}
	return false
}

func AsMessageFilterAffinity(rv reflect.Value) MessageFilterAffinity {
	if rv.IsValid() {
		if v, ok := rv.Convert(typeOfMessageFilterAffinity).Interface().(MessageFilterAffinity); ok {
			return v
		}
	}
	return nil
}

func isRedisBusyGroupError(err error) bool {
	return strings.HasPrefix(err.Error(), REDIS_BUSYGROUP_PREFIX)
}

func isMessageObserverAffinity(rv reflect.Value) bool {
	if rv.IsValid() {
		return rv.Type().AssignableTo(typeOfMessageObserverAffinity)
	}
	return false
}

func asMessageObserverAffinity(rv reflect.Value) MessageObserverAffinity {
	if rv.IsValid() {
		if v, ok := rv.Convert(typeOfMessageObserverAffinity).Interface().(MessageObserverAffinity); ok {
			return v
		}
	}
	return nil
}
