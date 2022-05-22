package utils

import (
	"fmt"
	"os"
	"time"
)

var file, _ = os.OpenFile("log.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

func WriteLog(data interface{}) {
	str := fmt.Sprint(data)

	file.WriteString(time.Now().Format(time.Kitchen) + ": " + str + "\n")
}
