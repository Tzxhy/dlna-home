package utils

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

type MyClaims struct {
	Username string `json:"username"`
	UserId   string `json:"userid"`
	jwt.StandardClaims
}

var MySecret = []byte("qva5im03q96fnjaga1rnafp3qrsi8r")

const TokenExpireDuration = time.Hour * 24

// GenToken 生成JWT
func GenToken(username string, userid string) (string, error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		username, // 自定义字段
		userid,
		jwt.StandardClaims{
			ExpiresAt: jwt.NewTime(float64(time.Now().Add(TokenExpireDuration).Unix())), // 过期时间
			Issuer:    "root",                                                           // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
