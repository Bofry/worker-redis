package internal

import (
	"reflect"
	"strings"
)

const (
	REDIS_BUSYGROUP_PREFIX = "BUSYGROUP"
)

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
