package service

import (
	"context"
	"sync"
	"time"
	"xg/da"
	"xg/entity"
	"xg/log"
)

type ISubjectService interface{
	ListSubjects(ctx context.Context, parentID int) ([]*entity.Subject, error)
	CreateSubject(ctx context.Context, req entity.CreateSubjectRequest) (int, error)
}

type SubjectService struct {
}

func (s *SubjectService) ListSubjects(ctx context.Context, parentID int) ([]*entity.Subject, error) {
	log.Info.Printf("ListSubjects, parentID: %#v\n", parentID)
	subjects, err := da.GetSubjectModel().SearchSubject(ctx, da.SearchSubjectCondition{
		ParentId: parentID,
	})
	if err != nil {
		log.Warning.Printf("Search subjects failed, parentID: %#v, err: %v\n", parentID, err)
		return nil, err
	}

	parentIds := make([]int, 0)
	for i := range subjects {
		parentIds = append(parentIds, subjects[i].ParentId)
	}

	parentSubjects, err := da.GetSubjectModel().SearchSubject(ctx, da.SearchSubjectCondition{
		IDList: parentIds,
	})
	if err != nil {
		log.Warning.Printf("Search subjects from sub orgs failed, parentIds: %#v, err: %v\n", parentIds, err)
		return nil, err
	}
	parentMap := make(map[int]*da.Subject)
	for i := range parentSubjects {
		parentMap[parentSubjects[i].ID] = parentSubjects[i]
	}

	res := make([]*entity.Subject, len(subjects))
	for i := range res {
		res[i] = &entity.Subject{
			ID:       subjects[i].ID,
			Level:    subjects[i].Level,
			ParentId: subjects[i].ParentId,
			Name:     subjects[i].Name,
			Parent:   convertSubject(parentMap[subjects[i].ParentId]),
		}
	}
	return res, nil
}

func (s *SubjectService) CreateSubject(ctx context.Context, req entity.CreateSubjectRequest) (int, error) {
	level := 1
	log.Info.Printf("CreateSubject, req: %#v\n", req)
	if req.ParentId > 0 {
		parentSubject, err := da.GetSubjectModel().GetSubjectById(ctx, req.ParentId)
		if err != nil {
			log.Warning.Printf("Search subjects from sub orgs failed, req: %#v, err: %v\n", req, err)
			return 0, err
		}
		level = parentSubject.Level + 1
	}
	now := time.Now()
	data := da.Subject{
		Level:    level,
		ParentId: req.ParentId,
		Name:     req.Name,

		UpdatedAt: &now,
		CreatedAt: &now,
	}
	id, err := da.GetSubjectModel().CreateSubject(ctx, data)
	if err != nil{
		log.Warning.Printf("Search subjects from sub orgs failed, req: %#v, data: %#v, err: %v\n", req, data, err)
		return -1, err
	}
	return id, nil
}

func convertSubject(subject *da.Subject) *entity.Subject {
	if subject == nil {
		return nil
	}
	res := &entity.Subject{
		ID:       subject.ID,
		Level:    subject.Level,
		ParentId: subject.ParentId,
		Name:     subject.Name,
	}
	return res
}

var (
	_subjectService     *SubjectService
	_subjectServiceOnce sync.Once
)

func GetSubjectService() *SubjectService {
	_subjectServiceOnce.Do(func() {
		if _subjectService == nil {
			_subjectService = new(SubjectService)
		}
	})
	return _subjectService
}
