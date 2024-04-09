package token

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
)

type Config struct {
	key string
	identityKey string
}

// 该错误表示auth请求头为空
var ErrMissingHeader = errors.New(" the length of the Authorization header is zero ")

var (
	config = Config{"Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5", "identityKey"}
	once sync.Once
)

// 包级别配置Config
func Init(key string, identityKey string) {
	once.Do(func(){
		if key != "" {
			config.key = key
		}
		if identityKey != "" {
			config.identityKey = identityKey
		}
	})
}

// 使用key 解析 token
func Parse(tokenString string, key string) (string, error) {
	token, err := jwt.Parse(tokenString, func (token *jwt.Token) (interface{}, error)  {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(key), nil
	})

	if err != nil {
		return "", err
	}

	var identityKey string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		identityKey = claims[config.identityKey].(string)
	}

	return identityKey, nil
}

// 从请求头中获取 token 信息, 并解析
func ParseRequest(c *gin.Context) (string, error) {
	header := c.Request.Header.Get("Authorization")

	if len(header) == 0 {
		return "", ErrMissingHeader
	}

	var t string
	// Sscanf 按指定格式解析字符串并赋值
	fmt.Sscanf(header, "Bearer %s", &t)

	return Parse(t, config.key)
}

// 签发函数
func Sign(identityKey string) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		config.identityKey: identityKey,
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(100000 * time.Hour).Unix(),
	})
	// 签发
	tokenString, err = token.SignedString([]byte(config.key))

	// 因为用了同名的变量，最后return可以不用写
	return 
}



