package ErrInfo

// key = ErrCode, string = ErrMsg
type ErrInfo struct {
	ErrCode int
	ErrMsg  string
}

var Json = ErrInfo{9000, "Incorrect JSON format"}
var Secret = ErrInfo{9001, "Incorrect secret"}
var Token = ErrInfo{9002, "Incorrect token"}
