package server

import (
	"errors"
	"net/http"
	"time"

	"git.2dfire.net/zerodb/common/statics/codes"
	"git.2dfire.net/zerodb/keeper/pkg/glog"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//JWTAuth gin json web token middleware
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		msg := NewResult()
		token := c.DefaultQuery("token", "")
		if token == "" {
			msg.Code = codes.Failed
			msg.ErrorCode = codes.ParameterIsNil
			msg.ErrMsg = "token is nil"
			unauthorized(c, http.StatusOK, &msg)
			return
		}
		j := NewJWT()

		claims, err := j.ParseToken(token)
		if err != nil {
			msg.Code = codes.Failed
			msg.ErrorCode = codes.CommonError
			msg.ErrMsg = err.Error()
			glog.GLog.Errorln(err)
			unauthorized(c, http.StatusOK, &msg)
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}

func unauthorized(c *gin.Context, code int, msg interface{}) {
	c.Abort()
	c.JSON(code, msg)
}

// CustomClaims 载荷
type CustomClaims struct {
	UserID   string `json:"userId"`
	Password string `json:"password"`
	jwt.StandardClaims
}

var (
	TokenExpired     = errors.New("Token is expired.")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token")
	SignKey          = "ZeroDatabase!@)*#)!@U#@*!@!)"
	Issuer           = "ZeroDatabase"
)

type JWT struct {
	SigningKey []byte
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

func GetSignKey() string {
	return SignKey
}

func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}

//GenerateToken 创建
func (j *JWT) GenerateToken(userID, password string) (string, error) {
	expireTime := time.Now().Add(24 * time.Hour)

	claims := CustomClaims{
		UserID:   userID,
		Password: password,
	}
	claims.StandardClaims.ExpiresAt = expireTime.Unix()
	claims.StandardClaims.Issuer = Issuer

	return j.CreatToken(claims)
}

//CreatToken 创建
func (j *JWT) CreatToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

//ParseToken 解析
func (j *JWT) ParseToken(strToken string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(strToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			}
			return nil, TokenInvalid
		}
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

//RefreshToken 刷新
func (j *JWT) RefreshToken(strToken string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(strToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(24 * time.Hour).Unix()
		claims.StandardClaims.Issuer = Issuer
		return j.CreatToken(*claims)
	}
	return "", TokenInvalid
}
