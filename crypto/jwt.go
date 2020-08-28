package crypto

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
	"xg/entity"
)
const secretKey = "LWUc9agXEJp0bdqj"

func GenerateToken(id, orgId, roleId int)(string, error){
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(12)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["uid"] = id
	claims["oid"] = orgId
	claims["rid"] = roleId
	token.Claims = claims

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func secret()jwt.Keyfunc{
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey),nil
	}
}

func ParseToken(token string)(*entity.JWTUser, error){
	tokenObj,err := jwt.Parse(token,secret())
	if err != nil{
		return nil, err
	}
	claim,ok := tokenObj.Claims.(jwt.MapClaims)
	if !ok{
		err = errors.New("cannot convert claim to mapclaim")
		return nil, err
	}
	//验证token，如果token被修改过则为false
	if  !tokenObj.Valid{
		err = errors.New("token is invalid")
		return nil, nil
	}

	uid := claim["uid"].(float64)
	rid := claim["rid"].(float64)
	oid := claim["oid"].(float64)
	return &entity.JWTUser{
		UserId: int(uid),
		OrgId:  int(oid),
		RoleId: int(rid),
	}, nil
}