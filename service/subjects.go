package service

import (
	"context"
	"sync"
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
