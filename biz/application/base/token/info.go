package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/xh-polaris/synapse/biz/conf"
	"github.com/xh-polaris/synapse/biz/types/cst"
)

type Info struct {
	RawToken    string
	BasicUserId string `json:"basic_user_id"`
	UserRole    string `json:"user_role"`
}

func SignJWT(tokenConf *conf.Token, info *Info) (string, error) {
	now := time.Now().UTC()

	claims := jwt.MapClaims{
		"iat":                now.Unix(),                                                    // sign time
		"nbf":                now.Unix(),                                                    // validate time
		"exp":                now.Add(time.Duration(tokenConf.Expire) * time.Second).Unix(), // expire time
		cst.TokenBasicUserID: info.BasicUserId,
		cst.TokenUserRole:    info.UserRole,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	sk, err := conf.GetSecretKey(tokenConf.SecretKey)
	if err != nil {
		return "", err
	}
	str, err := token.SignedString(sk)
	return str, err
}

func ParseJWT(tokenConf *conf.Token, str string) (*Info, error) {
	token, err := jwt.Parse(str, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return conf.GetPublicKey(tokenConf.PublicKey)
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		info := &Info{
			BasicUserId: claims["basic_user_id"].(string),
			UserRole:    claims["user_role"].(string),
			RawToken:    str,
		}
		return info, nil
	}
	return nil, errors.New("invalid token")
}
