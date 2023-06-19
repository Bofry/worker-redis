package internal

var (
	replyCodeNames = []string{
		UNSET: "unset",
		PASS:  "pass",
		FAIL:  "fail",
		ABORT: "abort",
	}
)

type ReplyCode int

func (code ReplyCode) String() string {
	if (code > __reply_code_maximum__) || (code < __reply_code_minimum__) {
		return __reply_code_invalid_text__
	}
	return replyCodeNames[code]
}
