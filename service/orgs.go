package service

import (
	"context"
	"strings"
	"sync"
	"xg/da"
	"xg/db"
	"xg/entity"
	"xg/log"
	"xg/utils"

	"github.com/jinzhu/gorm"
)

type IOrgService interface {
	CreateOrg(ctx context.Context, req *entity.CreateOrgRequest, operator *entity.JWTUser) (int, error)
	CreateOrgWithSubOrgs(ctx context.Context, req *entity.CreateOrgWithSubOrgsRequest, operator *entity.JWTUser) (int, error)
	UpdateOrgById(ctx context.Context, req *entity.UpdateOrgRequest, operator *entity.JWTUser) error
	RevokeOrgById(ctx context.Context, id int, operator *entity.JWTUser) error
	CheckOrgById(ctx context.Context, id, status int, operator *entity.JWTUser) error
	GetOrgById(ctx context.Context, orgId int) (*entity.Org, error)
	GetOrgSubjectsById(ctx context.Context, orgId int) ([]string, error)
	ListOrgs(ctx context.Context, condition da.SearchOrgsCondition) (int, []*entity.Org, error)
	ListOrgsByStatus(ctx context.Context, status []int) (int, []*entity.Org, error)
	SearchSubOrgs(ctx context.Context, condition da.SearchOrgsCondition) (int, []*entity.SubOrgWithDistance, error)
	SearchPendingOrgs(ctx context.Context, condition da.SearchOrgsCondition) (int, []*entity.SubOrgWithDistance, error)

	UpdateOrgWithSubOrgs(ctx context.Context, orgId int, req *entity.UpdateOrgWithSubOrgsRequest, operator *entity.JWTUser) error
}

type OrgService struct {
}

func (s *OrgService) CreateOrg(ctx context.Context, req *entity.CreateOrgRequest, operator *entity.JWTUser) (int, error) {
	log.Info.Printf("create org, req: %#v\n", req)

	req.SupportRoleID = []int{entity.RoleOutOrg, entity.RoleSeniorOutOrg}
	id, err := s.createOrg(ctx, db.Get(), req, operator)
	if err != nil {
		log.Warning.Printf("Create org failed, req: %#v, err: %v\n", req, err)
		return -1, err
	}
	return id, nil
}

func (s *OrgService) CreateOrgWithSubOrgs(ctx context.Context, req *entity.CreateOrgWithSubOrgsRequest, operator *entity.JWTUser) (int, error) {
	log.Info.Printf("create org with sub orgs, req: %#v\n", req)
	cid, err := db.GetTransResult(ctx, func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
		req.OrgData.SupportRoleID = []int{entity.RoleOutOrg, entity.RoleSeniorOutOrg}
		cid, err := s.createOrg(ctx, tx, &req.OrgData, operator)
		if err != nil {
			log.Warning.Printf("Create org failed, req: %#v, err: %v\n", req, err)
			return -1, nil
		}
		for i := range req.SubOrgs {
			status := entity.OrgStatusCreated
			//req.SubOrgs[i].SupportRoleID = []int{entity.RoleOutOrg, entity.RoleSeniorOutOrg}

			//获取经纬度信息
			if req.SubOrgs[i].Longitude == 0 && req.SubOrgs[i].Latitude == 0 {
				cor, err := utils.GetAddressLocation(req.SubOrgs[i].Address + req.SubOrgs[i].AddressExt)
				if err != nil {
					log.Warning.Printf("Get address failed, req: %#v, err: %v\n", req, err)
				} else {
					req.SubOrgs[i].Latitude = cor.Latitude
					req.SubOrgs[i].Longitude = cor.Longitude
				}
			}
			_, err = da.GetOrgModel().CreateOrg(ctx, tx, da.Org{
				Name:          req.OrgData.Name + "-" + req.SubOrgs[i].Name,
				Subjects:      strings.Join(req.SubOrgs[i].Subjects, ","),
				Status:        status,
				Address:       req.SubOrgs[i].Address,
				AddressExt:    req.SubOrgs[i].AddressExt,
				ParentID:      cid,
				Telephone:     req.SubOrgs[i].Telephone,
				SupportRoleID: entity.IntArrayToString(req.SubOrgs[i].SupportRoleID),
				Latitude:      req.SubOrgs[i].Latitude,
				Longitude:     req.SubOrgs[i].Longitude,
			})
			if err != nil {
				log.Warning.Printf("Create sub org failed, req: %#v, err: %v\n", req, err)
				return -1, nil
			}
		}
		return cid, nil
	})
	if err != nil {
		return -1, err
	}
	return cid.(int), nil
}

func (s *OrgService) UpdateOrgById(ctx context.Context, req *entity.UpdateOrgRequest, operator *entity.JWTUser) error {
	return s.updateOrgById(ctx, db.Get(), req, operator)
}

func (s *OrgService) RevokeOrgById(ctx context.Context, id int, operator *entity.JWTUser) error {
	log.Info.Printf("revoke org, id: %#v\n", id)
	if id == 1 {
		log.Warning.Printf("Can't revoke root org, id: %v\n", id)
		return ErrOperateOnRootOrg
	}
	org, err := da.GetOrgModel().GetOrgById(ctx, db.Get(), id)
	if err != nil {
		log.Warning.Printf("Get revoke org failed, id: %v, err: %v\n", id, err)
		return err
	}
	if org.Status == entity.OrgStatusRejected {
		log.Info.Printf("Revoke org workded, id: %v\n", id)
		return nil
	}
	if org.ParentID != 0 {
		log.Info.Printf("Revoke sub org, id: %v, org: %#v\n", id, org)
		return ErrNotSuperOrg
	}

	subOrgs, err := da.GetOrgModel().GetOrgsByParentIdWithStatus(ctx, org.ID, []int{entity.OrgStatusCreated, entity.OrgStatusCertified, entity.OrgStatusRejected})
	if err != nil {
		log.Warning.Printf("Get orgs by parent failed, org: %#v, err: %v\n", org, err)
		return err
	}

	err = db.GetTrans(ctx, func(ctx context.Context, tx *gorm.DB) error {
		err = da.GetOrgModel().UpdateOrg(ctx, tx, id, da.Org{
			Status: entity.OrgStatusRevoked,
		})
		if err != nil {
			log.Warning.Printf("Update org status failed, org: %#v, err: %v\n", org, err)
			return err
		}
		for i := range subOrgs {
			err = da.GetOrgModel().UpdateOrg(ctx, tx, subOrgs[i].ID, da.Org{
				Status: entity.OrgStatusRevoked,
			})
			if err != nil {
				log.Warning.Printf("Update sub org status failed, sub org: %#v, err: %v\n", subOrgs[i], err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *OrgService) CheckOrgById(ctx context.Context, id, status int, operator *entity.JWTUser) error {
	log.Info.Printf("check org, id: %#v, status: %#v\n", id, status)
	org, err := da.GetOrgModel().GetOrgById(ctx, db.Get(), id)
	if err != nil {
		log.Warning.Printf("Get org failed, orgid: %#v, err: %v\n", id, err)
		return err
	}
	if org.Status != entity.OrgStatusCreated {
		log.Warning.Printf("Org status invalid, org: %#v, err: %v\n", org, err)
		return ErrInvalidOrgStatus
	}
	if org.ParentID != 0 {
		log.Warning.Printf("Org is sub org, org: %#v, err: %v\n", org, err)
		return ErrNotSuperOrg
	}

	subOrgs, err := da.GetOrgModel().GetOrgsByParentIdWithStatus(ctx, org.ID, []int{entity.OrgStatusCreated, entity.OrgStatusCertified, entity.OrgStatusRejected})
	if err != nil {
		log.Warning.Printf("Get orgs by parent id failed, org: %#v, err: %v\n", org, err)
		return err
	}

	err = db.GetTrans(ctx, func(ctx context.Context, tx *gorm.DB) error {
		err = da.GetOrgModel().UpdateOrg(ctx, tx, id, da.Org{
			Status: status,
		})
		if err != nil {
			log.Warning.Printf("Updaste org status failed, org: %#v, err: %v\n", org, err)
			return err
		}
		for i := range subOrgs {
			err = da.GetOrgModel().UpdateOrg(ctx, tx, subOrgs[i].ID, da.Org{
				Status: status,
			})
			if err != nil {
				log.Warning.Printf("Updaste sub org status failed, org: %#v, err: %v\n", subOrgs[i], err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *OrgService) GetOrgById(ctx context.Context, orgId int) (*entity.Org, error) {
	log.Info.Printf("get org by id, id: %#v\n", orgId)
	org, err := da.GetOrgModel().GetOrgById(ctx, db.Get(), orgId)
	if err != nil {
		log.Warning.Printf("Get org failed, orgId: %#v, err: %v\n", orgId, err)
		return nil, err
	}
	var subjects []string
	if len(org.Subjects) > 0 {
		subjects = strings.Split(org.Subjects, ",")
	}
	subOrgs, err := da.GetOrgModel().GetOrgsByParentIdWithStatus(ctx, orgId, []int{entity.OrgStatusCreated, entity.OrgStatusCertified})
	if err != nil {
		log.Warning.Printf("Get org by parent id failed, orgId: %#v, err: %v\n", orgId, err)
		return nil, err
	}
	return &entity.Org{
		ID:            org.ID,
		Name:          org.Name,
		Subjects:      subjects,
		Status:        org.Status,
		Address:       org.Address,
		AddressExt:    org.AddressExt,
		ParentID:      org.ParentID,
		Telephone:     org.Telephone,
		SupportRoleID: entity.StringToIntArray(org.SupportRoleID),
		SubOrgs:       ToOrgEntities(subOrgs),
	}, nil
}

func (s *OrgService) GetOrgSubjectsById(ctx context.Context, orgId int) ([]string, error) {
	log.Info.Printf("GetOrgSubjectsById, id: %#v\n", orgId)
	org, err := da.GetOrgModel().GetOrgById(ctx, db.Get(), orgId)
	if err != nil {
		log.Warning.Printf("Get org by id failed, orgId: %#v, err: %v\n", orgId, err)
		return nil, err
	}
	var subjects []string
	if len(org.Subjects) > 0 {
		subjects = strings.Split(org.Subjects, ",")
	}
	subOrgs, err := da.GetOrgModel().GetOrgsByParentId(ctx, orgId)
	for i := range subOrgs {
		if len(subOrgs[i].Subjects) > 0 {
			subjects = append(subjects, strings.Split(subOrgs[i].Subjects, ",")...)
		}
	}
	return subjects, nil
}

func (s *OrgService) ListOrgs(ctx context.Context, condition da.SearchOrgsCondition) (int, []*entity.Org, error) {
	log.Info.Printf("ListOrgs, condition: %#v\n", condition)
	condition.Status = []int{entity.OrgStatusCertified}
	condition.ParentIDs = []int{0}

	count, orgs, err := da.GetOrgModel().SearchOrgs(ctx, condition)
	if err != nil {
		log.Warning.Printf("Search org failed, condition: %#v, err: %v\n", condition, err)
		return 0, nil, err
	}
	res := make([]*entity.Org, len(orgs))

	for i := range orgs {
		var subjects []string
		if len(orgs[i].Subjects) > 0 {
			subjects = strings.Split(orgs[i].Subjects, ",")
		}
		res[i] = &entity.Org{
			ID:            orgs[i].ID,
			Name:          orgs[i].Name,
			Subjects:      subjects,
			Status:        orgs[i].Status,
			Address:       orgs[i].Address,
			AddressExt:    orgs[i].AddressExt,
			ParentID:      orgs[i].ParentID,
			Telephone:     orgs[i].Telephone,
			SupportRoleID: entity.StringToIntArray(orgs[i].SupportRoleID),
		}
	}
	return count, res, nil
}

func (s *OrgService) ListOrgsByStatus(ctx context.Context, status []int) (int, []*entity.Org, error) {
	log.Info.Printf("ListOrgsByStatus, status: %#v\n", status)
	condition := da.SearchOrgsCondition{
		Status:    status,
		ParentIDs: []int{0},
	}
	count, orgs, err := da.GetOrgModel().SearchOrgs(ctx, condition)
	if err != nil {
		log.Warning.Printf("Search org failed, condition: %#v, err: %v\n", condition, err)
		return 0, nil, err
	}
	res := make([]*entity.Org, len(orgs))

	for i := range orgs {
		var subjects []string
		if len(orgs[i].Subjects) > 0 {
			subjects = strings.Split(orgs[i].Subjects, ",")
		}
		res[i] = &entity.Org{
			ID:            orgs[i].ID,
			Name:          orgs[i].Name,
			Subjects:      subjects,
			Status:        orgs[i].Status,
			Address:       orgs[i].Address,
			AddressExt:    orgs[i].AddressExt,
			ParentID:      orgs[i].ParentID,
			Telephone:     orgs[i].Telephone,
			SupportRoleID: entity.StringToIntArray(orgs[i].SupportRoleID),
		}
	}
	return count, res, nil
}

func (s *OrgService) SearchSubOrgs(ctx context.Context, condition da.SearchOrgsCondition) (int, []*entity.SubOrgWithDistance, error) {
	condition.Status = []int{entity.OrgStatusCertified}
	condition.IsSubOrg = true

	return s.searchOrgs(ctx, condition)
}

func (s *OrgService) SearchPendingOrgs(ctx context.Context, condition da.SearchOrgsCondition) (int, []*entity.Org, error) {
	condition.Status = []int{entity.OrgStatusCreated}
	condition.ParentIDs = []int{0}
	count, orgs, err := da.GetOrgModel().SearchOrgs(ctx, condition)
	if err != nil {
		log.Warning.Printf("Search org failed, condition: %#v, err: %v\n", condition, err)
		return 0, nil, err
	}
	res := make([]*entity.Org, len(orgs))

	for i := range orgs {
		var subjects []string
		if len(orgs[i].Subjects) > 0 {
			subjects = strings.Split(orgs[i].Subjects, ",")
		}
		res[i] = &entity.Org{
			ID:            orgs[i].ID,
			Name:          orgs[i].Name,
			Subjects:      subjects,
			Status:        orgs[i].Status,
			Address:       orgs[i].Address,
			AddressExt:    orgs[i].AddressExt,
			ParentID:      orgs[i].ParentID,
			Telephone:     orgs[i].Telephone,
			SupportRoleID: entity.StringToIntArray(orgs[i].SupportRoleID),
		}
	}
	return count, res, nil
}

func (s *OrgService) searchOrgs(ctx context.Context, condition da.SearchOrgsCondition) (int, []*entity.SubOrgWithDistance, error) {
	log.Info.Printf("SearchSubOrgs, condition: %#v\n", condition)
	if condition.StudentID < 0 {
		log.Warning.Printf("Student id invalid, condition: %#v\n", condition)
		return 0, nil, ErrStudentIdNeeded
	}
	student, err := da.GetStudentModel().GetStudentById(ctx, condition.StudentID)
	if err != nil {
		log.Warning.Printf("Search student failed, condition: %#v, err: %v\n", condition, err)
		return 0, nil, err
	}

	count, orgs, err := da.GetOrgModel().SearchOrgsWithDistance(ctx, condition, &entity.Coordinate{
		Longitude: student.Longitude,
		Latitude:  student.Latitude,
	})
	if err != nil {
		log.Warning.Printf("Search org failed, condition: %#v, err: %v\n", condition, err)
		return 0, nil, err
	}
	res := make([]*entity.SubOrgWithDistance, 0)
	for i := range orgs {
		var subjects []string
		if len(orgs[i].Subjects) > 0 {
			allSubjects := s.filterSubjects(orgs[i].Subjects)
			if len(condition.Subjects) > 0 {
				subjects = s.filterTargetSubjects(allSubjects, condition.Subjects)
			} else {
				subjects = append(subjects, allSubjects...)
			}
		}
		if len(orgs[i].Subjects) > 0 && len(subjects) == 0 {
			continue
		}
		res = append(res, &entity.SubOrgWithDistance{
			ID:            orgs[i].ID,
			Name:          orgs[i].Name,
			Subjects:      subjects,
			Status:        orgs[i].Status,
			Address:       orgs[i].Address,
			AddressExt:    orgs[i].AddressExt,
			ParentID:      orgs[i].ParentID,
			Telephone:     orgs[i].Telephone,
			SupportRoleID: entity.StringToIntArray(orgs[i].SupportRoleID),
			Distance:      orgs[i].Distance,
		})
	}
	return count, res, nil
}
func (s *OrgService) UpdateOrgWithSubOrgs(ctx context.Context, orgId int, req *entity.UpdateOrgWithSubOrgsRequest, operator *entity.JWTUser) error {
	updateEntity, err := s.prepareUpdateSubOrgs(ctx, orgId, req, operator)
	if err != nil {
		return err
	}
	log.Trace.Printf("OrgId:%v, req:%#v\n", orgId, req)
	err = db.GetTrans(ctx, func(ctx context.Context, tx *gorm.DB) error {
		//更新主机构
		addr := req.OrgData.Address + req.OrgData.AddressExt
		if addr != "" {
			cor, err := utils.GetAddressLocation(addr)
			if err != nil {
				log.Warning.Printf("Get address failed, req: %#v, err: %v\n", req, err)
			} else {
				req.OrgData.Latitude = cor.Latitude
				req.OrgData.Longitude = cor.Longitude
			}
		}

		err = s.updateOrgById(ctx, tx, &entity.UpdateOrgRequest{
			ID:         orgId,
			Name: 		req.OrgData.Name,
			Address:    req.OrgData.Address,
			AddressExt: req.OrgData.AddressExt,
			Telephone:  req.OrgData.Telephone,
			Latitude:   req.OrgData.Latitude,
			Longitude:  req.OrgData.Longitude,
		}, operator)
		if err != nil {
			log.Warning.Printf("Update org failed, orgId: %#v, req: %#v, err: %v\n", orgId, req, err)
			return err
		}
		//增加
		for i := range updateEntity.InsertOrgList {

			//获取经纬度信息
			if updateEntity.InsertOrgList[i].Longitude == 0 && updateEntity.InsertOrgList[i].Latitude == 0 {
				cor, err := utils.GetAddressLocation(updateEntity.InsertOrgList[i].Address + updateEntity.InsertOrgList[i].AddressExt)
				if err != nil {
					log.Warning.Printf("Get address failed, req: %#v, err: %v\n", req, err)
				} else {
					updateEntity.InsertOrgList[i].Latitude = cor.Latitude
					updateEntity.InsertOrgList[i].Longitude = cor.Longitude
				}
			}
			_, err = da.GetOrgModel().CreateOrg(ctx, tx, da.Org{
				Name:       updateEntity.OrgInfo.Name + "-" + updateEntity.InsertOrgList[i].Name,
				Subjects:   strings.Join(updateEntity.InsertOrgList[i].Subjects, ","),
				Status:     updateEntity.OrgInfo.Status,
				Address:    updateEntity.InsertOrgList[i].Address,
				AddressExt: updateEntity.InsertOrgList[i].AddressExt,
				ParentID:   orgId,
				Telephone:  updateEntity.InsertOrgList[i].Telephone,
				//SupportRoleID: entity.IntArrayToString([]int{entity.RoleOutOrg, entity.RoleSeniorOutOrg}),
				Longitude: updateEntity.InsertOrgList[i].Longitude,
				Latitude:  updateEntity.InsertOrgList[i].Latitude,
			})
			if err != nil {
				log.Warning.Printf("Create sub org failed, orgId: %#v, req: %#v, err: %v\n", orgId, req, err)
				return err
			}
		}

		//修改
		for i := range updateEntity.UpdateOrgsList {

			//获取经纬度信息
			if updateEntity.UpdateOrgsList[i].Longitude == 0 && updateEntity.UpdateOrgsList[i].Latitude == 0 {
				cor, err := utils.GetAddressLocation(updateEntity.UpdateOrgsList[i].Address + updateEntity.UpdateOrgsList[i].AddressExt)
				if err != nil {
					log.Warning.Printf("Get address failed, req: %#v, err: %v\n", req, err)
				} else {
					updateEntity.UpdateOrgsList[i].Latitude = cor.Latitude
					updateEntity.UpdateOrgsList[i].Longitude = cor.Longitude
				}
			}
			namePairs := strings.Split(updateEntity.UpdateOrgsList[i].Name, "-")
			realName := namePairs[len(namePairs) - 1]

			err = da.GetOrgModel().UpdateOrg(ctx, tx, updateEntity.UpdateOrgsList[i].ID, da.Org{
				Name:       updateEntity.OrgInfo.Name + "-" +  realName,
				Subjects:   strings.Join(updateEntity.UpdateOrgsList[i].Subjects, ","),
				Status:     updateEntity.OrgInfo.Status,
				Address:    updateEntity.UpdateOrgsList[i].Address,
				AddressExt: updateEntity.UpdateOrgsList[i].AddressExt,
				Telephone:  updateEntity.UpdateOrgsList[i].Telephone,
				Longitude:  updateEntity.UpdateOrgsList[i].Longitude,
				Latitude:   updateEntity.UpdateOrgsList[i].Latitude,
			})
			if err != nil {
				log.Warning.Printf("Update sub org failed, orgId: %#v, req: %#v, err: %v\n", orgId, req, err)
				return err
			}
		}

		//删除
		for i := range updateEntity.DeletedIds {
			err = da.GetOrgModel().UpdateOrg(ctx, tx, updateEntity.DeletedIds[i], da.Org{
				Status: entity.OrgStatusRevoked,
			})
			// err = da.GetOrgModel().DeleteOrgById(ctx, tx, updateEntity.DeletedIds)
			if err != nil {
				log.Warning.Printf("Delete sub org failed, orgId: %#v, req: %#v, err: %v\n", orgId, req, err)
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (o *OrgService) filterSubjects(subject string) []string {
	var subjects []string
	if len(subject) > 0 {
		subjects = strings.Split(subject, ",")
	}
	return subjects
}

func (o *OrgService) filterTargetSubjects(subjects []string, targetSubjects []string) []string {
	ret := make([]string, 0)
	for i := range subjects {
		for j := range targetSubjects {
			if o.compareSubject(subjects[i], targetSubjects[j]) {
				ret = append(ret, subjects[i])
			}
		}
	}

	return ret
}

func (o *OrgService) prepareUpdateSubOrgs(ctx context.Context, orgId int, req *entity.UpdateOrgWithSubOrgsRequest, operator *entity.JWTUser) (*entity.UpdateSubOrgsEntity, error) {
	log.Info.Printf("update org with sub orgs, id: %#v, req: %#v\n", orgId, req)
	orgObj, err := da.GetOrgModel().GetOrgById(ctx, db.Get(), orgId)
	if err != nil {
		log.Warning.Printf("Get org failed, orgId: %#v, err: %v\n", orgId, err)
		return nil, err
	}
	if orgObj.ParentID != 0 {
		return nil, ErrNotSuperOrg
	}
	subOrgs, err := da.GetOrgModel().GetOrgsByParentIdWithStatus(ctx, orgId, []int{entity.OrgStatusCreated, entity.OrgStatusCertified, entity.OrgStatusRejected})
	if err != nil {
		log.Warning.Printf("Get sub orgs failed, orgId: %#v, err: %v\n", orgId, err)
		return nil, err
	}

	//标记新增，修改
	insertOrgsReq := make([]*entity.CreateOrUpdateOrgRequest, 0)
	updateOrgsReq := make([]*entity.CreateOrUpdateOrgRequest, 0)
	for i := range req.SubOrgs {
		if req.SubOrgs[i].ID == 0 {
			insertOrgsReq = append(insertOrgsReq, req.SubOrgs[i])
		} else {
			flag := false
			for j := range subOrgs {
				//若该id不在组织中，则忽略
				if req.SubOrgs[i].ID == subOrgs[j].ID {
					flag = true
				}
			}
			if flag {
				updateOrgsReq = append(updateOrgsReq, req.SubOrgs[i])
			}
		}
	}

	//标记删除
	deletedIds := make([]int, 0)
	for i := range subOrgs {
		flag := false
		for j := range updateOrgsReq {
			if updateOrgsReq[j].ID == subOrgs[i].ID {
				flag = true
				break
			}
		}
		if !flag {
			deletedIds = append(deletedIds, subOrgs[i].ID)
		}
	}

	//更新org信息
	if req.OrgData.Name != "" {
		orgObj.Name = req.OrgData.Name
	}
	if req.OrgData.Address != "" {
		orgObj.Address = req.OrgData.Address
	}
	if req.OrgData.AddressExt != "" {
		orgObj.Address = req.OrgData.AddressExt
	}
	if req.OrgData.Telephone != "" {
		orgObj.Address = req.OrgData.Telephone
	}
	orgInfo := &entity.Org{
		ID:         orgObj.ID,
		Name:       orgObj.Name,
		Address:    orgObj.Address,
		AddressExt: orgObj.AddressExt,
		ParentID:   orgObj.ParentID,
		Telephone:  orgObj.Telephone,
		Status:     orgObj.Status,
	}

	return &entity.UpdateSubOrgsEntity{
		OrgInfo:        orgInfo,
		UpdateOrgReq:   req.OrgData,
		InsertOrgList:  insertOrgsReq,
		UpdateOrgsList: updateOrgsReq,
		DeletedIds:     deletedIds,
	}, nil
}

func (o *OrgService) compareSubject(subjectA, subjectB string) bool {
	subjectAPairs := strings.Split(subjectA, "-")
	subjectBPairs := strings.Split(subjectB, "-")
	subjectALen := len(subjectAPairs)
	subjectBLen := len(subjectBPairs)

	//超过2段，比较前2段
	if subjectALen >= 2 && subjectBLen >= 2 {
		if subjectAPairs[0] == subjectBPairs[0] &&
			subjectAPairs[1] == subjectBPairs[1] {
			return true
		}
		return false
	}
	minLen := subjectALen
	if subjectBLen < subjectALen {
		minLen = subjectBLen
	}
	//不足2段，比较短的
	for i := 0; i < minLen; i++ {
		if subjectAPairs[i] != subjectBPairs[i] {
			return false
		}
	}
	return true
}

func (o *OrgService) GetSubOrgs(ctx context.Context, orgId int) ([]*da.Org, error) {
	subOrgs, err := da.GetOrgModel().GetOrgsByParentId(ctx, orgId)
	if err != nil {
		log.Warning.Printf("Get orgs by parent failed, orgId: %#v, err: %v\n", orgId, err)
		return nil, err
	}
	return subOrgs, nil
}

func (s *OrgService) createOrg(ctx context.Context, tx *gorm.DB, req *entity.CreateOrgRequest, operator *entity.JWTUser) (int, error) {
	status := entity.OrgStatusCreated
	if req.ParentID != 0 {
		parentOrg, err := da.GetOrgModel().GetOrgById(ctx, tx, req.ParentID)
		if err != nil {
			log.Warning.Printf("Create org failed, req: %#v, err: %v\n", req, err)
			return 0, err
		}
		status = parentOrg.Status
	}
	//获取经纬度信息
	if req.Longitude == 0 && req.Latitude == 0 {
		cor, err := utils.GetAddressLocation(req.Address + req.AddressExt)
		if err != nil {
			log.Warning.Printf("Get address failed, req: %#v, err: %v\n", req, err)
		} else {
			req.Latitude = cor.Latitude
			req.Longitude = cor.Longitude
		}
	}

	data := da.Org{
		Name:          req.Name,
		Subjects:      strings.Join(req.Subjects, ","),
		Status:        status,
		Address:       req.Address,
		AddressExt:    req.AddressExt,
		ParentID:      req.ParentID,
		Telephone:     req.Telephone,
		SupportRoleID: entity.IntArrayToString(req.SupportRoleID),
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
	}
	id, err := da.GetOrgModel().CreateOrg(ctx, tx, data)
	if err != nil {
		log.Warning.Printf("Create org failed, data: %#v, err: %v\n", data, err)
		return id, err
	}
	return id, nil
}

func (s *OrgService) updateOrgById(ctx context.Context, tx *gorm.DB, req *entity.UpdateOrgRequest, operator *entity.JWTUser) error {
	log.Info.Printf("update org, req: %#v\n", req)
	err := da.GetOrgModel().UpdateOrg(ctx, tx, req.ID, da.Org{
		Subjects:   strings.Join(req.Subjects, ","),
		Name: req.Name,
		Address:    req.Address,
		AddressExt: req.AddressExt,
		Telephone:  req.Telephone,
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
		//Status:   req.Status,
	})
	if err != nil {
		log.Warning.Printf("Update org failed, req: %#v, err: %v\n", req, err)
		return err
	}
	return nil
}

func ToOrgEntity(org *da.Org) *entity.Org {
	var subjects []string
	if len(org.Subjects) > 0 {
		subjects = strings.Split(org.Subjects, ",")
	}
	return &entity.Org{
		ID:            org.ID,
		Name:          org.Name,
		Subjects:      subjects,
		Status:        org.Status,
		Address:       org.Address,
		AddressExt:    org.AddressExt,
		ParentID:      org.ParentID,
		Telephone:     org.Telephone,
		SupportRoleID: entity.StringToIntArray(org.SupportRoleID),
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
