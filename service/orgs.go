package service

import (
	"context"
	"strings"
	"sync"
	"xg/da"
	"xg/db"
	"xg/entity"
)

type OrgService struct {
}


func (s *OrgService) CreateOrg(ctx context.Context, req *entity.CreateOrgRequest, operator *entity.JWTUser) (int, error) {
	status := entity.OrgStatusCreated
	if req.ParentID != 0 {
		parentOrg, err := da.GetOrgModel().GetOrgById(ctx, req.ParentID)
		if err != nil{
			return 0, err
		}
		status = parentOrg.Status
	}

	return da.GetOrgModel().CreateOrg(ctx, da.Org{
		Name:     req.Name,
		Subjects: strings.Join(req.Subjects, ","),
		Status:   status,
		Address: req.Address,
		ParentID: req.ParentID,
	})
}

func (s *OrgService) UpdateOrgById(ctx context.Context, req *entity.UpdateOrgRequest, operator *entity.JWTUser) error {

	return da.GetOrgModel().UpdateOrg(ctx, db.Get(), req.ID, da.Org{
		Subjects: strings.Join(req.Subjects, ","),
		Address: req.Address,
		//Status:   req.Status,
	})
}

func (s *OrgService) CheckOrgById(ctx context.Context, id, status int, operator *entity.JWTUser) error {

	org, err := da.GetOrgModel().GetOrgById(ctx, id)
	if err != nil {
		return err
	}
	if org.Status != entity.OrgStatusCreated {
		return ErrInvalidOrgStatus
	}
	if org.ParentID != 0 {
		return ErrNotSuperOrg
	}

	_, subOrgs, err := da.GetOrgModel().SearchOrgs(ctx, da.SearchOrgsCondition{
		ParentIDs: []int{org.ID},
	})
	if err != nil {
		return err
	}
	tx := db.Get().Begin()
	err = da.GetOrgModel().UpdateOrg(ctx, tx, id, da.Org{
		Status: status,
	})
	if err != nil {
		tx.Rollback()
		return err
	}
	for i := range subOrgs {
		err = da.GetOrgModel().UpdateOrg(ctx, tx, subOrgs[i].ID, da.Org{
			Status: status,
		})
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (s *OrgService) GetOrgById(ctx context.Context, orgId int) (*entity.Org, error) {
	org, err := da.GetOrgModel().GetOrgById(ctx, orgId)
	if err != nil {
		return nil, err
	}
	var subjects []string
	if len(org.Subjects) > 0 {
		subjects = strings.Split(org.Subjects, ",")
	}
	_, subOrgs, err := da.GetOrgModel().SearchOrgs(ctx, da.SearchOrgsCondition{
		ParentIDs: []int{orgId},
	})

	return &entity.Org{
		ID:       org.ID,
		Name:     org.Name,
		Subjects: subjects,
		Status:   org.Status,
		Address: org.Address,
		ParentID: org.ParentID,
		SubOrgs: ToOrgEntities(subOrgs),
	}, nil
}

func (s *OrgService) ListOrgs(ctx context.Context) (int, []*entity.Org, error) {
	count, orgs, err := da.GetOrgModel().SearchOrgs(ctx, da.SearchOrgsCondition{
		Status: []int{
			entity.OrgStatusCertified,
		},
	})
	if err != nil {
		return 0, nil, err
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
			Address: orgs[i].Address,
			ParentID: orgs[i].ParentID,
		}
	}
	return count, res, nil
}

func (s *OrgService) ListOrgsByStatus(ctx context.Context, status []int) (int, []*entity.Org, error) {
	count, orgs, err := da.GetOrgModel().SearchOrgs(ctx, da.SearchOrgsCondition{
		Status: status,
		ParentIDs: []int{0},
	})
	if err != nil {
		return 0, nil, err
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
			Address: orgs[i].Address,
			ParentID: orgs[i].ParentID,
		}
	}
	return count, res, nil
}



func (s *OrgService) SearchSubOrgs(ctx context.Context, condition da.SearchOrgsCondition) (int, []*entity.Org, error) {
	condition.Status = []int{entity.OrgStatusCertified}
	condition.IsSubOrg = true
	count, orgs, err := da.GetOrgModel().SearchOrgs(ctx, condition)
	if err != nil {
		return 0, nil, err
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
			Address: orgs[i].Address,
			ParentID: orgs[i].ParentID,
		}
	}
	return count, res, nil
}
func (o *OrgService) GetSubOrgs(ctx context.Context, orgId int)([]*da.Org, error){
	_, subOrgs, err := da.GetOrgModel().SearchOrgs(ctx, da.SearchOrgsCondition{
		ParentIDs: []int{orgId},
	})
	if err != nil{
		return nil, err
	}
	return subOrgs, nil
}
func ToOrgEntity(org *da.Org) *entity.Org{
	var subjects []string
	if len(org.Subjects) > 0 {
		subjects = strings.Split(org.Subjects, ",")
	}
	return &entity.Org{
		ID:       org.ID,
		Name:     org.Name,
		Subjects: subjects,
		Status:   org.Status,
		Address: org.Address,
		ParentID: org.ParentID,
	}
}

func ToOrgEntities(orgs []*da.Org) []*entity.Org{
	ret := make([]*entity.Org, len(orgs))
	for i := range orgs{
		ret[i] = ToOrgEntity(orgs[i])
	}
	return ret
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
