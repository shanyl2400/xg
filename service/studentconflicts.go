package service

import (
	"context"
	"errors"
	"sync"
	"xg/da"
	"xg/db"
	"xg/entity"
	"xg/log"
	"xg/utils"

	"github.com/jinzhu/gorm"
)

var (
	ErrMarkProcessedRecord  = errors.New("mark processed record")
	ErrNoValidStudentID     = errors.New("no valid student id")
	ErrInvalidRecordStatus  = errors.New("invlaid record status")
	ErrInvalidStudentStatus = errors.New("invlaid student status")
	ErrUnconflictStudent    = errors.New("unconflict student")
)

type IStudentConflictService interface {
	CreateOrUpdateStudentConflict(ctx context.Context, tx *gorm.DB, record entity.CreateStudentConflictRequest) error
	HandleStudentConflict(ctx context.Context, r entity.HandleStudentConflictRequest) error
	SearchStudentConflicts(ctx context.Context, ss da.SearchStudentConflictCondition) (int, []*entity.StudentConflictRecordDetails, error)

	UpdateConflictStudentStatus(ctx context.Context, r entity.HandleUpdateConflictStudentStatusRequest) error
	HandleStudentConflictRecordStatus(ctx context.Context, r entity.HandleStudentConflictStatusRequest) error
}

type StudentConflictService struct {
}

func (s *StudentConflictService) CreateOrUpdateStudentConflict(ctx context.Context, tx *gorm.DB, record entity.CreateStudentConflictRequest) error {
	//查询未处理的
	_, records, err := da.GetStudentConflictModel().SearchStudentConflicts(ctx, da.SearchStudentConflictCondition{
		Telephone: record.Telephone,
		Status:    []int{entity.StudentConflictStatusUnprocessed},
	})
	if err != nil {
		log.Error.Printf("Seach conflict records failed, record: %#v, err: %v\n", record, err)
		return err
	}
	//为创建过,则创建
	if len(records) == 0 {
		_, err := da.GetStudentConflictModel().CreateStudentConflict(ctx, tx, da.StudentConflict{
			Telephone: record.Telephone,
			Status:    entity.StudentConflictStatusUnprocessed,
			Total:     record.Total,
			AuthorID:  record.AuthorID,
		})
		if err != nil {
			log.Error.Printf("Create conflict records failed, record: %#v, err: %v\n", record, err)
			return err
		}
		return nil
	}

	//已创建过,则更新
	r := records[0]
	r.AuthorID = record.AuthorID
	r.Total = record.Total
	err = da.GetStudentConflictModel().UpdateStudentConflict(ctx, tx, r.ID, *r)
	if err != nil {
		log.Error.Printf("Update conflict records failed, record: %#v, err: %v\n", r, err)
		return err
	}
	return nil
}

func (s *StudentConflictService) HandleStudentConflictRecordStatus(ctx context.Context, r entity.HandleStudentConflictStatusRequest) error {
	err := s.checkConflictRecordStatus(ctx, r.Status)
	if err != nil {
		log.Error.Printf("Bad request req: %#v, err: %v\n", r, err)
		return err
	}

	record, err := da.GetStudentConflictModel().GetStudentConflictByID(ctx, r.RecordID)
	if err != nil {
		log.Error.Printf("Get conflict records failed, id: %#v, err: %v\n", r.RecordID, err)
		return err
	}

	record.Status = r.Status
	err = da.GetStudentConflictModel().UpdateStudentConflict(ctx, db.Get(), record.ID, *record)
	if err != nil {
		log.Error.Printf("Update conflict records failed, record: %#v, err: %v\n", record, err)
		return err
	}
	return nil
}
func (s *StudentConflictService) UpdateConflictStudentStatus(ctx context.Context, r entity.HandleUpdateConflictStudentStatusRequest) error {
	err := s.checkStudentStatus(ctx, r.Status)
	if err != nil {
		log.Error.Printf("Bad request req: %#v, err: %v\n", r, err)
		return ErrMarkProcessedRecord
	}

	student, err := da.GetStudentModel().GetStudentById(ctx, r.StudentID)
	if err != nil {
		log.Error.Printf("Search students failed, req: %#v, err: %v\n", r, err)
		return err
	}
	//check student in conflict record
	_, records, err := da.GetStudentConflictModel().SearchStudentConflicts(ctx, da.SearchStudentConflictCondition{
		Telephone: student.Telephone,
	})
	if err != nil {
		log.Error.Printf("Search students conflict records failed, student: %#v, err: %v\n", student, err)
		return err
	}
	if len(records) == 0 {
		log.Error.Printf("Conflict records not found, records: %#v, err: %v\n", records, err)
		return ErrUnconflictStudent
	}

	flag := false
	for i := range records {
		if records[i].Status == entity.StudentConflictStatusUnprocessed {
			flag = true
			break
		}
	}
	if !flag {
		log.Error.Printf("No processing conflict record, records: %#v, err: %v\n", records, err)
		return ErrUnconflictStudent
	}

	student.Status = r.Status
	err = da.GetStudentModel().UpdateStudent(ctx, db.Get(), r.StudentID, *student)
	if err != nil {
		log.Error.Printf("Update student status failed, student: %#v, err: %v\n", student, err)
		return err
	}

	return nil
}

func (s *StudentConflictService) HandleStudentConflict(ctx context.Context, r entity.HandleStudentConflictRequest) error {
	record, err := da.GetStudentConflictModel().GetStudentConflictByID(ctx, r.RecordID)
	if err != nil {
		log.Error.Printf("Get conflict records failed, id: %#v, err: %v\n", r.RecordID, err)
		return err
	}
	//排除已处理的conflict
	if record.Status == entity.StudentConflictStatusProcessed {
		log.Error.Printf("Mark processed conflict, record: %#v, err: %v\n", record, err)
		return ErrMarkProcessedRecord
	}

	_, students, err := da.GetStudentModel().SearchStudents(ctx, da.SearchStudentCondition{
		Telephone: record.Telephone,
	})
	if err != nil {
		log.Error.Printf("Search students failed, Telephone: %#v, err: %v\n", record.Telephone, err)
		return err
	}

	err = db.GetTrans(ctx, func(ctx context.Context, tx *gorm.DB) error {
		//更新冲突记录状态
		record.Status = entity.StudentConflictStatusProcessed
		err = da.GetStudentConflictModel().UpdateStudentConflict(ctx, tx, record.ID, *record)
		if err != nil {
			log.Error.Printf("Update conflict records failed, record: %#v, err: %v\n", record, err)
			return err
		}

		//update students
		//更新名单状态
		flag := false
		for i := range students {
			if students[i].ID != r.SelectStudentID &&
				(students[i].Status == entity.StudentConflictSuccess ||
					students[i].Status == entity.StudentCreated) {
				students[i].Status = entity.StudentConflictFailed
			} else if students[i].ID == r.SelectStudentID && students[i].Status == entity.StudentConflictFailed {
				flag = true
				students[i].Status = entity.StudentConflictSuccess
			} else {
				continue
			}
			err = da.GetStudentModel().UpdateStudent(ctx, tx, students[i].ID, *students[i])
			if err != nil {
				log.Error.Printf("Update student status failed, students: %#v, err: %v\n", students[i], err)
				return err
			}
		}
		if !flag {
			log.Error.Printf("no valid student id, students: %#v, studentID: %v, err: %v\n", students, r.SelectStudentID, err)
			return ErrNoValidStudentID
		}

		return nil
	})
	if err != nil {
		log.Error.Printf("handle conflict records failed, id: %#v, studentID: %v, err: %v\n", r.RecordID, r.SelectStudentID, err)
		return err
	}
	return nil
}
func (s *StudentConflictService) SearchStudentConflicts(ctx context.Context, ss da.SearchStudentConflictCondition) (int, []*entity.StudentConflictRecordDetails, error) {
	log.Info.Printf("SearchStudentConflicts, condition: %#v\n", ss)
	total, records, err := da.GetStudentConflictModel().SearchStudentConflicts(ctx, ss)
	if err != nil {
		log.Warning.Printf("SearchStudents failed, condition: %#v, err: %v\n", ss, err)
		return 0, nil, err
	}
	authorIds := make([]int, len(records))
	for i := range records {
		authorIds[i] = records[i].AuthorID
	}

	_, users, err := da.GetUsersModel().SearchUsers(ctx, da.SearchUserCondition{
		IDList: utils.UniqueInts(authorIds),
	})
	if err != nil {
		log.Warning.Printf("Get User failed, ids: %#v, req: %#v, err: %v\n", authorIds, ss, err)
		return 0, nil, err
	}

	authorNameMaps := make(map[int]string)
	for i := range users {
		authorNameMaps[users[i].ID] = users[i].Name
	}
	res := make([]*entity.StudentConflictRecordDetails, len(records))
	for i := range records {
		res[i] = &entity.StudentConflictRecordDetails{
			ID:         records[i].ID,
			Telephone:  records[i].Telephone,
			AuthorID:   records[i].AuthorID,
			AuthorName: authorNameMaps[records[i].AuthorID],
			Total:      records[i].Total,
			Status:     records[i].Status,
			UpdatedAt:  records[i].UpdatedAt,
		}
	}
	return total, res, nil
}

func (s *StudentConflictService) checkConflictRecordStatus(ctx context.Context, status int) error {
	if status != entity.StudentConflictStatusUnprocessed && status != entity.StudentConflictStatusProcessed {
		return ErrInvalidRecordStatus
	}
	return nil
}
func (s *StudentConflictService) checkStudentStatus(ctx context.Context, status int) error {
	if status != entity.StudentCreated && status != entity.StudentConflictFailed && status != entity.StudentConflictSuccess {
		return ErrInvalidStudentStatus
	}
	return nil
}

var (
	_studentConflictService     *StudentConflictService
	_studentConflictServiceOnce sync.Once
)

func GetStudentConflictService() *StudentConflictService {
	_studentConflictServiceOnce.Do(func() {
		if _studentConflictService == nil {
			_studentConflictService = new(StudentConflictService)
		}
	})
	return _studentConflictService
}
