package service

import (
	"context"
	"sync"
	"xg/da"
	"xg/entity"
	"xg/log"
)

type AuthService struct {

}

func (a *AuthService) ListAuths(ctx context.Context) ([]*entity.Auth, error){
	auths, err := da.GetAuthModel().ListAuth(ctx)
	if err != nil{
		log.Warning.Printf("List auth failed, err: %v\n", err)
		return nil, err
	}
	res := make([]*entity.Auth, len(auths))
	for i := range auths {
		res[i] = (*entity.Auth)(auths[i])
	}
	return res, nil
}

var(
	_authService *AuthService
	__authServiceOnce sync.Once
)

func GetAuthService() *AuthService{
	__authServiceOnce.Do(func() {
		if _authService == nil{
			_authService = new(AuthService)
		}
	})
	return _authService
}