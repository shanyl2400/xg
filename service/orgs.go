package service

import (
	"context"
	"strings"
	"sync"
	"xg/da"
	"xg/entity"
)

type OrgService struct {
}

func (s *OrgService) checkCreateAuth(ctx context.Context, operator *entity.JWTUser) (bool, error) {
	//TODO:Check auth
	return false, nil
}

func (s *OrgService) CreateOrg(ctx context.Context, req *entity.CreateOrgRequest, operator *entity.JWTUser) (int, error) {
	canStatus, err := s.checkCreateAuth(ctx, operator)
	if err != nil {
		return -1, err
	}
	if !canStatus {
		req.Status = entity.OrgStatusCreated
	}

	return da.GetOrgModel().CreateOrg(ctx, da.Org{
		Name:     req.Name,
		Subjects: strings.Join(req.Subjects, ","),
		Status:   req.Status,
	})
}

func (s *OrgService) UpdateOrgById(ctx context.Context, req *entity.UpdateOrgRequest, operator *entity.JWTUser) error {
	canStatus, err := s.checkCreateAuth(ctx, operator)
	if err != nil {
		return err
	}
	if !canStatus {
		return ErrNoAuthorizeToOperate
	}

	return da.GetOrgModel().UpdateOrg(ctx, req.ID, da.Org{
		Subjects: strings.Join(req.Subjects, ","),
		Status:   req.Status,
	})
}

func (s *OrgService) ListOrgs(ctx context.Context) ([]*entity.Org, error) {
	orgs, err := da.GetOrgModel().SearchOrgs(ctx, da.SearchOrgsCondition{
		Status: []int{
			entity.OrgStatusCertified,
		},
	})
	if err != nil {
		return nil, err
	}
	res := make([]*entity.Org, len(orgs))

	for i := range orgs {
		var subjects []string
		if len(orgs[i].Subjects) > 0 {
			subjects = strings.Split(orgs[i].Subjects, ",")
		}
		res[i] = &entity.Org{
			ID:       orgs[i].ID,
			Name:     orgs[i].Name,
			Subjects: subjects,
			Status:   orgs[i].Status,
		}
	}
	return res, nil
}

var (
	_orgService     *OrgService
	_orgServiceOnce sync.Once
)

func GetOrgService() *OrgService {
	_orgServiceOnce.Do(func() {
		if _orgService == nil {
			_orgService = new(OrgService)
		}
	})
	return _orgService
}
