package jwt

import (
	"app/pkg/config"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dgrijalva/jwt-go"
)

type TokenModel struct {
	Token string
	User  primitive.ObjectID
}

var tokenKey = ""

func init() {
	tokenKey = config.GetConfig().Token_KEY
}

func MakeToken(model *TokenModel) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = model.User
	claims["expDate"] = time.Now().Add(time.Hour * 500).Format("2006-01-02 15:04:05")

	t, err := token.SignedString([]byte(tokenKey))
	if err != nil {
		return "", err
	}

	return t, nil

}

func ParseToken(tokenStr string) (*TokenModel, error) {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenKey), nil
	})
	if err != nil {
		return nil, errors.New("token解析失败")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		userId := claims["userId"]
		expDate := claims["expDate"]
		timeNow := time.Now().Format("2006-01-02 15:04:05")
		if timeNow > expDate.(string) {
			return nil, errors.New("token已过期")
		}
		if userId == "" {
			return nil, errors.New("token解密，id为空")
		}

		model := &TokenModel{}
		model.User, _ = primitive.ObjectIDFromHex(userId.(string))

		return model, nil

	} else {
		return nil, errors.New("token claims 解析失败")
	}

}
