package constants

type CodeWithTip struct {
	Code int
	Tip  string
}

const CODE_OK = 0

// 通用
const (
	CODE_PARAMS_NOT_VALID = 900_0000 + iota
)

var CODE_PARAMS_NOT_VALID_TIPS = &CodeWithTip{
	CODE_PARAMS_NOT_VALID,
	"参数无效",
}
