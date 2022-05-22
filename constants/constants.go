package constants

import (
	"os"
	"path/filepath"
)

var ROOT_PATH, _ = os.Getwd()

var UPLOAD_PATH = filepath.Join(ROOT_PATH, "upload")

const TOKEN_COOKIE_NAME = "token"

const USER_ID_SPLITTER = ";"

// 需要和许多tag保持一致
const DIR_ROOT_ID = "ROOT"
