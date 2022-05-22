package utils

import (
	"strings"

	"gitee.com/tzxhy/dlna-home/constants"
)

func GenerateUid() string {
	return RandStringBytesMaskImprSrc(5)
}
func GenerateDid() string {
	return RandStringBytesMaskImprSrc(8)
}
func GenerateFid() string {
	return RandStringBytesMaskImprSrc(10)
}

func GenerateGid() string {
	return RandStringBytesMaskImprSrc(5)
}

func GenerateRid() string {
	return RandStringBytesMaskImprSrc(8)
}

func GeneratePassword() string {
	return RandStringBytesMaskImprSrc(16)
}

func GetUserIds(userIds string) *[]string {
	arr := strings.Split(userIds, constants.USER_ID_SPLITTER)
	return &arr
}
