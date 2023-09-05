package internal

import redis "github.com/Bofry/lib-redis-stream"

var (
	defaultStreamSetting = &StreamSetting{}
)

type StreamSetting struct {
	MessageStateKeyPrefix *string
}

func (s *StreamSetting) DecodeMessageContentOption() []redis.DecodeMessageContentOption {
	var options []redis.DecodeMessageContentOption
	if s.MessageStateKeyPrefix != nil {
		options = append(options, redis.WithMessageStateKeyPrefix(*s.MessageStateKeyPrefix))
	}
	return options
}
