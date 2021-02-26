package service

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"
	"xg/da"
	"xg/db"
	"xg/entity"
	"xg/log"

	"github.com/jinzhu/gorm"
)

var (
	ErrLevelToHigh = errors.New("level can't be more than 3")
)

type ISubjectService interface {
	ListSubjects(ctx context.Context, parentID int) ([]*entity.Subject, error)
	CreateSubject(ctx context.Context, req entity.CreateSubjectRequest) (int, error)
	BatchCreateSubject(ctx context.Context, reqs []*entity.CreateSubjectRequest) error
	ListSubjectsTree(ctx context.Context) ([]*entity.SubjectTreeNode, error)
}

type SubjectService struct {
	sync.RWMutex
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

func (s *SubjectService) ListSubjectsTree(ctx context.Context) ([]*entity.SubjectTreeNode, error) {
	subjects, err := da.GetSubjectModel().SearchSubject(ctx, da.SearchSubjectCondition{})
	if err != nil {
		log.Warning.Printf("Search subjects failed, err: %v\n", err)
		return nil, err
	}
	subjectsMap := make(map[int]*entity.SubjectTreeNode)
	for i := range subjects {
		subjectsMap[subjects[i].ID] = &entity.SubjectTreeNode{
			ID:    subjects[i].ID,
			Level: subjects[i].Level,
			Name:  subjects[i].Name,
			Label: subjects[i].Name,
			Title: subjects[i].Name,
			Value: subjects[i].Name,
			Key:   subjects[i].ID,
		}
	}

	rootSubjectsIDs := make([]int, 0)
	for i := range subjects {
		if subjects[i].ParentId != 0 {
			subjectsMap[subjects[i].ParentId].Children = append(subjectsMap[subjects[i].ParentId].Children, subjectsMap[subjects[i].ID])
		} else {
			rootSubjectsIDs = append(rootSubjectsIDs, subjects[i].ID)
		}
	}

	res := make([]*entity.SubjectTreeNode, 0)
	for i := range rootSubjectsIDs {
		fillSubjectsValue(subjectsMap[rootSubjectsIDs[i]], "")
		res = append(res, subjectsMap[rootSubjectsIDs[i]])
	}

	return res, nil
}

func (s *SubjectService) ListSubjectsAll(ctx context.Context) ([]*entity.SubjectTreeNode, error) {
	subjects, err := da.GetSubjectModel().SearchSubject(ctx, da.SearchSubjectCondition{})
	if err != nil {
		log.Warning.Printf("Search subjects failed, err: %v\n", err)
		return nil, err
	}
	subjectsMap := make(map[int]*entity.SubjectTreeNode)
	for i := range subjects {
		subjectsMap[subjects[i].ID] = &entity.SubjectTreeNode{
			ID:    subjects[i].ID,
			Level: subjects[i].Level,
			Name:  subjects[i].Name,
			Title: subjects[i].Name,
			Value: subjects[i].Name,
			Key:   subjects[i].ID,
		}
	}

	for i := range subjects {
		if subjects[i].ParentId != 0 {
			subjectsMap[subjects[i].ParentId].Children = append(subjectsMap[subjects[i].ParentId].Children, subjectsMap[subjects[i].ID])
		}
	}

	res := make([]*entity.SubjectTreeNode, 0)
	for k := range subjectsMap {
		fillSubjectsValue(subjectsMap[k], "")
		res = append(res, subjectsMap[k])
	}

	return res, nil
}

func (s *SubjectService) CreateSubject(ctx context.Context, req entity.CreateSubjectRequest) (int, error) {
	level := 1
	log.Info.Printf("CreateSubject, req: %#v\n", req)

	s.Lock()
	defer s.Unlock()

	if strings.TrimSpace(req.Name) == "" {
		log.Warning.Printf("CreateSubject invalid name, req: %#v\n", req)
		return -1, ErrInvalidSubjectName
	}

	err := s.checkSubjectsUnique(ctx, req.ParentId, []string{req.Name})
	if err != nil {
		log.Warning.Printf("Check duplicate name failed, req: %#v\n", req)
		return -1, err
	}

	if req.ParentId > 0 {
		parentSubject, err := da.GetSubjectModel().GetSubjectById(ctx, req.ParentId)
		if err != nil {
			log.Warning.Printf("Search subjects from sub orgs failed, req: %#v, err: %v\n", req, err)
			return 0, err
		}
		level = parentSubject.Level + 1
	}
	if level > 3 {
		log.Warning.Printf("Level is more than 3, req: %#v, level: %v\n", req, level)
		return 0, ErrLevelToHigh
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
	if err != nil {
		log.Warning.Printf("Search subjects from sub orgs failed, req: %#v, data: %#v, err: %v\n", req, data, err)
		return -1, err
	}
	return id, nil
}

func (s *SubjectService) BatchCreateSubject(ctx context.Context, reqs []*entity.CreateSubjectRequest) error {
	s.Lock()
	defer s.Unlock()

	//check unique name
	nameList := make([]string, len(reqs))
	parentID := 0
	for i := range reqs {
		//check name
		if strings.TrimSpace(reqs[i].Name) == "" {
			log.Warning.Printf("CreateSubject invalid name, req: %#v\n", reqs[i])
			return ErrInvalidSubjectName
		}
		nameList[i] = reqs[i].Name
		parentID = reqs[i].ParentId
	}
	if len(nameList) > 0 {
		err := s.checkSubjectsUnique(ctx, parentID, nameList)
		if err != nil {
			log.Warning.Printf("Check duplicate name failed, req: %#v\n", reqs)
			return err
		}
	}
	err := db.GetTrans(ctx, func(ctx context.Context, tx *gorm.DB) error {
		for i := range reqs {
			req := reqs[i]
			level := 1
			log.Info.Printf("CreateSubject, req: %#v\n", req)

			if req.ParentId > 0 {
				parentSubject, err := da.GetSubjectModel().GetSubjectById(ctx, req.ParentId)
				if err != nil {
					log.Warning.Printf("Search subjects from sub orgs failed, req: %#v, err: %v\n", req, err)
					return err
				}
				level = parentSubject.Level + 1
			}
			if level > 3 {
				log.Warning.Printf("Level is more than 3, req: %#v, level: %v\n", req, level)
				return ErrLevelToHigh
			}

			now := time.Now()
			data := da.Subject{
				Level:    level,
				ParentId: req.ParentId,
				Name:     req.Name,

				UpdatedAt: &now,
				CreatedAt: &now,
			}
			_, err := da.GetSubjectModel().CreateSubjectTx(ctx, tx, data)
			if err != nil {
				log.Warning.Printf("Search subjects from sub orgs failed, req: %#v, data: %#v, err: %v\n", req, data, err)
				return err
			}
		}
		return nil
	})

	return err
}

func fillSubjectsValue(subject *entity.SubjectTreeNode, prefix string) {
	if prefix != "" {
		subject.Value = prefix + "-" + subject.Name
	}
	for i := range subject.Children {
		fillSubjectsValue(subject.Children[i], subject.Value)
	}
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

func (s *SubjectService) checkSubjectsUnique(ctx context.Context, parentID int, name []string) error {
	condition := da.SearchSubjectCondition{
		Names: name,
	}
	if parentID == 0 {
		condition.RootSubject = true
	} else {
		condition.ParentId = parentID
	}
	subjects, err := da.GetSubjectModel().SearchSubject(ctx, condition)
	if err != nil {
		log.Error.Printf("Can't get subjects, condition: %#v, err: %#v", condition, err)
		return err
	}
	if len(subjects) > 0 {
		log.Warning.Printf("Duplicate subject name, subjects: %#v, err: %#v", subjects, err)
		return ErrDuplicateSubjectName
	}
	return nil
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
