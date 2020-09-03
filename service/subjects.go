package service

import (
	"context"
	"sync"
	"time"
	"xg/da"
	"xg/entity"
)

type SubjectService struct {
}

func (s *SubjectService) ListSubjects(ctx context.Context, parentID int) ([]*entity.Subject, error) {
	subjects, err := da.GetSubjectModel().SearchSubject(ctx, da.SearchSubjectCondition{
		ParentId: parentID,
	})
	if err != nil {
		return nil, err
	}
	res := make([]*entity.Subject, len(subjects))
	for i := range res {
		res[i] = &entity.Subject{
			ID:       subjects[i].ID,
			Level:    subjects[i].Level,
			ParentId: subjects[i].ParentId,
			Name:     subjects[i].Name,
		}
	}
	return res, nil
}

func (s *SubjectService) CreateSubject(ctx context.Context, req entity.CreateSubjectRequest) (int, error) {
	level := 1
	if req.ParentId > 0 {
		parentSubject ,err := da.GetSubjectModel().GetSubjectById(ctx, req.ParentId)
		if err != nil{
			return 0, err
		}
		level = parentSubject.Level + 1
	}
	now := time.Now()
	id, err := da.GetSubjectModel().CreateSubject(ctx, da.Subject{
		Level:    level,
		ParentId: req.ParentId,
		Name:     req.Name,

		UpdatedAt: &now,
		CreatedAt: &now,
	})
	return id, err
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
