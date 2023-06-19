package internal

import "strings"

const (
	REDIS_BUSYGROUP_PREFIX = "BUSYGROUP"
)

func isRedisBusyGroupError(err error) bool {
	return strings.HasPrefix(err.Error(), REDIS_BUSYGROUP_PREFIX)
}
