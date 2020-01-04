package errcode

var (
	Success        = 0
	NoPatternMatch = 400
	ServerBusy     = 500
)
var errorMap = map[int]string{
	Success:        "success",
	ServerBusy:     "server busy",
	NoPatternMatch: "uri path invalid",
}

func GetMsg(ret int) string {
	return errorMap[ret]
}
