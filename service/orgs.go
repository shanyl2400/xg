package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"xg/da"
	"xg/db"
	"xg/entity"

	"github.com/jinzhu/gorm"
)

type OrgService struct {
}

func (s *OrgService) CreateOrg(ctx context.Context, req *entity.CreateOrgRequest, operator *entity.JWTUser) (int, error) {
	return s.createOrg(ctx, db.Get(), req, operator)
}

func (s *OrgService) CreateOrgWithSubOrgs(ctx context.Context, req *entity.CreateOrgWithSubOrgsRequest, operator *entity.JWTUser) (int, error) {
	tx := db.Get().Begin()
	cid, err := s.createOrg(ctx, tx, &req.OrgData, operator)
	if err != nil {
		tx.Rollback()
		return -1, nil
	}
	for i := range req.SubOrgs {
		status := entity.OrgStatusCreated
		_, err = da.GetOrgModel().CreateOrg(ctx, db.Get(), da.Org{
			Name:      req.OrgData.Name + "-" + req.SubOrgs[i].Name,
			Subjects:  strings.Join(req.SubOrgs[i].Subjects, ","),
			Status:    status,
			Address:   req.SubOrgs[i].Address,
			ParentID:  cid,
			Telephone: req.SubOrgs[i].Telephone,
		})
		if err != nil {
			tx.Rollback()
			return -1, nil
		}
	}
	tx.Commit()
	return cid, nil
}

func (s *OrgService) UpdateOrgById(ctx context.Context, req *entity.UpdateOrgRequest, operator *entity.JWTUser) error {

	return da.GetOrgModel().UpdateOrg(ctx, db.Get(), req.ID, da.Org{
		Subjects: strings.Join(req.Subjects, ","),
		Address:  req.Address,
		//Status:   req.Status,
	})
}

func (s *OrgService) RevokeOrgById(ctx context.Context, id int, operator *entity.JWTUser) error {
	if id == 1 {
		return ErrOperateOnRootOrg
	}
	org, err := da.GetOrgModel().GetOrgById(ctx, db.Get(), id)
	if err != nil {
		return err
	}
	if org.Status == entity.OrgStatusRejected {
		return nil
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
		Status: entity.OrgStatusRevoked,
	})
	if err != nil {
		tx.Rollback()
		return err
	}
	for i := range subOrgs {
		err = da.GetOrgModel().UpdateOrg(ctx, tx, subOrgs[i].ID, da.Org{
			Status: entity.OrgStatusRevoked,
		})
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (s *OrgService) CheckOrgById(ctx context.Context, id, status int, operator *entity.JWTUser) error {
	org, err := da.GetOrgModel().GetOrgById(ctx, db.Get(), id)
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
	org, err := da.GetOrgModel().GetOrgById(ctx, db.Get(), orgId)
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
		ID:        org.ID,
		Name:      org.Name,
		Subjects:  subjects,
		Status:    org.Status,
		Address:   org.Address,
		ParentID:  org.ParentID,
		Telephone: org.Telephone,
		SubOrgs:   ToOrgEntities(subOrgs),
	}, nil
}

func (s *OrgService) GetOrgSubjectsById(ctx context.Context, orgId int) ([]string, error) {
	org, err := da.GetOrgModel().GetOrgById(ctx, db.Get(), orgId)
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
	fmt.Println(subjects)
	for i := range subOrgs {
		fmt.Println(subOrgs[i].Subjects)
		if len(subOrgs[i].Subjects) > 0 {
			subjects = append(subjects, strings.Split(subOrgs[i].Subjects, ",")...)
		}
	}
	fmt.Println(subjects)
	return subjects, nil
}

func (s *OrgService) ListOrgs(ctx context.Context) (int, []*entity.Org, error) {
	count, orgs, err := da.GetOrgModel().SearchOrgs(ctx, da.SearchOrgsCondition{
		Status: []int{
			entity.OrgStatusCertified,
		},
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
			ID:        orgs[i].ID,
			Name:      orgs[i].Name,
			Subjects:  subjects,
			Status:    orgs[i].Status,
			Address:   orgs[i].Address,
			ParentID:  orgs[i].ParentID,
			Telephone: orgs[i].Telephone,
		}
	}
	return count, res, nil
}

func (s *OrgService) ListOrgsByStatus(ctx context.Context, status []int) (int, []*entity.Org, error) {
	count, orgs, err := da.GetOrgModel().SearchOrgs(ctx, da.SearchOrgsCondition{
		Status:    status,
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
			ID:        orgs[i].ID,
			Name:      orgs[i].Name,
			Subjects:  subjects,
			Status:    orgs[i].Status,
			Address:   orgs[i].Address,
			ParentID:  orgs[i].ParentID,
			Telephone: orgs[i].Telephone,
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
			ID:        orgs[i].ID,
			Name:      orgs[i].Name,
			Subjects:  subjects,
			Status:    orgs[i].Status,
			Address:   orgs[i].Address,
			ParentID:  orgs[i].ParentID,
			Telephone: orgs[i].Telephone,
		}
	}
	return count, res, nil
}
func (o *OrgService) GetSubOrgs(ctx context.Context, orgId int) ([]*da.Org, error) {
	_, subOrgs, err := da.GetOrgModel().SearchOrgs(ctx, da.SearchOrgsCondition{
		ParentIDs: []int{orgId},
	})
	if err != nil {
		return nil, err
	}
	return subOrgs, nil
}

func (s *OrgService) createOrg(ctx context.Context, tx *gorm.DB, req *entity.CreateOrgRequest, operator *entity.JWTUser) (int, error) {
	status := entity.OrgStatusCreated
	if req.ParentID != 0 {
		parentOrg, err := da.GetOrgModel().GetOrgById(ctx, tx, req.ParentID)
		if err != nil {
			return 0, err
		}
		status = parentOrg.Status
	}

	return da.GetOrgModel().CreateOrg(ctx, tx, da.Org{
		Name:      req.Name,
		Subjects:  strings.Join(req.Subjects, ","),
		Status:    status,
		Address:   req.Address,
		ParentID:  req.ParentID,
		Telephone: req.Telephone,
	})
}

func ToOrgEntity(org *da.Org) *entity.Org {
	var subjects []string
	if len(org.Subjects) > 0 {
		subjects = strings.Split(org.Subjects, ",")
	}
	return &entity.Org{
		ID:        org.ID,
		Name:      org.Name,
		Subjects:  subjects,
		Status:    org.Status,
		Address:   org.Address,
		ParentID:  org.ParentID,
		Telephone: org.Telephone,
	}
}

func ToOrgEntities(orgs []*da.Org) []*entity.Org {
	ret := make([]*entity.Org, len(orgs))
	for i := range orgs {
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
