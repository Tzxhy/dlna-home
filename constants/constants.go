package constants

import (
	"os"
)

var ROOT_PATH, _ = os.Getwd()

const TOKEN_COOKIE_NAME = "token"

const USER_ID_SPLITTER = ";"

// 需要和许多tag保持一致
const DIR_ROOT_ID = "ROOT"
