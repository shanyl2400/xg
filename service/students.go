package service

import (
	"context"
	"strings"
	"sync"
	"time"
	"xg/da"
	"xg/db"
	"xg/entity"
)

const (
	ChallengeDate = 60 * time.Hour * 24
)

type StudentService struct {
}

func (s *StudentService) CreateStudent(ctx context.Context, c *entity.CreateStudentRequest, operator *entity.JWTUser) (int, int, error) {
	status := entity.StudentCreated

	//检查冲单，查询相同手机号是否有学生
	total, students, err := da.GetStudentModel().SearchStudents(ctx, da.SearchStudentCondition{
		Telephone: c.Telephone,
		OrderBy:   "created_at",
		PageSize:  1,
		Page:      0,
	})
	if err != nil {
		return -1, -1, err
	}
	//找到对应的学生
	//TODO:自己挑战自己的情况？
	if total > 0 {
		now := time.Now()
		latestAdd := students[0].CreatedAt
		timeDiff := now.Sub(*latestAdd)
		if timeDiff > ChallengeDate {
			//挑战成功
			status = entity.StudentConflictSuccess
		} else {
			//挑战失败
			status = entity.StudentConflictFailed
		}
	}

	tx := db.Get().Begin()
	//添加学生记录
	id, err := da.GetStudentModel().CreateStudent(ctx, tx, da.Student{
		Name:          c.Name,
		Gender:        c.Gender,
		Telephone:     c.Telephone,
		Address:       c.Address,
		Email:         c.Email,
		IntentSubject: strings.Join(c.IntentSubject, ","),
		Status:        status,
		AuthorID:      operator.UserId,
		OrderSourceID: c.OrderSourceID,
		Note:          c.Note,
	})
	if err != nil {
		tx.Rollback()
		return -1, -1, err
	}
	err = GetStatisticsService().AddStudent(ctx, tx, 1)
	if err != nil {
		tx.Rollback()
		return -1, -1, err
	}
	tx.Commit()
	return id, status, nil
}

func (s *StudentService) UpdateStudent(ctx context.Context, id int, req *entity.UpdateStudentRequest) error {
	return da.GetStudentModel().UpdateStudent(ctx, id, da.Student{
		Name:          req.Name,
		Gender:        req.Gender,
		Telephone:     req.Telephone,
		Address:       req.Address,
		Email:         req.Email,
		IntentSubject: strings.Join(req.IntentSubject, ","),
		OrderSourceID: req.OrderSourceID,
	})
}

func (s *StudentService) GetStudentById(ctx context.Context, id int, operator *entity.JWTUser) (*entity.StudentInfosWithOrders, error) {
	//查询学生信息
	student, err := da.GetStudentModel().GetStudentById(ctx, id)
	if err != nil {
		return nil, err
	}

	//查询相关订单
	_, orders, err := da.GetOrderModel().SearchOrder(ctx, da.SearchOrderCondition{
		StudentIDList: []int{student.ID},
	})
	if err != nil {
		return nil, err
	}

	//查询用户和机构信息
	orgIds := make([]int, len(orders))
	publisherIds := make([]int, len(orders))
	for i := range orders {
		orgIds[i] = orders[i].ToOrgID
		publisherIds[i] = orders[i].PublisherID
	}
	publishers, err := da.GetUsersModel().SearchUsers(ctx, da.SearchUserCondition{
		IDList: publisherIds,
	})
	if err != nil {
		return nil, err
	}
	orgs, err := da.GetOrgModel().ListOrgs(ctx)
	if err != nil {
		return nil, err
	}

	orderSources, err := da.GetOrderSourceModel().ListOrderSources(ctx)
	if err != nil {
		return nil, err
	}

	//构建map
	publisherMap := make(map[int]*da.User)
	for i := range publishers {
		publisherMap[publishers[i].ID] = publishers[i]
	}
	orgMap := make(map[int]*da.Org)
	for i := range orgs {
		orgMap[orgs[i].ID] = orgs[i]
	}

	orderSourceMap := make(map[int]string)
	for i := range orderSources {
		orderSourceMap[orderSources[i].ID] = orderSources[i].Name
	}

	user, err := da.GetUsersModel().GetUserById(ctx, student.AuthorID)
	if err != nil {
		return nil, err
	}

	//构建返回数据
	res := &entity.StudentInfosWithOrders{
		StudentInfo: entity.StudentInfo{
			ID:              student.ID,
			Name:            student.Name,
			Gender:          student.Gender,
			Telephone:       student.Telephone,
			Address:         student.Address,
			AuthorID:        student.AuthorID,
			Email:           student.Email,
			AuthorName:      user.Name,
			OrderSourceID:   student.OrderSourceID,
			OrderSourceName: orderSourceMap[student.OrderSourceID],
			IntentSubject:   strings.Split(student.IntentSubject, ","),
			Status:          student.Status,
			Note:            student.Note,
		},
	}

	//构建订单数据
	res.Orders = make([]*entity.OrderInfoDetails, len(orders))
	for i := range orders {
		res.Orders[i] = &entity.OrderInfoDetails{
			OrderInfo: entity.OrderInfo{
				ID:            orders[i].ID,
				StudentID:     orders[i].StudentID,
				ToOrgID:       orders[i].ToOrgID,
				IntentSubject: strings.Split(orders[i].IntentSubjects, ","),
				PublisherID:   orders[i].PublisherID,
				Status:        orders[i].Status,
			},
			StudentName:      student.Name,
			StudentTelephone: student.Telephone,
			OrgName:          orgMap[orders[i].ToOrgID].Name,
			PublisherName:    publisherMap[orders[i].PublisherID].Name,
		}
	}

	return res, nil
}

func (s *StudentService) SearchPrivateStudents(ctx context.Context, ss *entity.SearchStudentRequest, operator *entity.JWTUser) (int, []*entity.StudentInfo, error) {
	ss.AuthorIDList = []int{operator.UserId}
	return s.SearchStudents(ctx, ss, operator)
}

func (s *StudentService) SearchStudents(ctx context.Context, ss *entity.SearchStudentRequest, operator *entity.JWTUser) (int, []*entity.StudentInfo, error) {
	total, students, err := da.GetStudentModel().SearchStudents(ctx, da.SearchStudentCondition{
		Name:         ss.Name,
		Telephone:    ss.Telephone,
		Address:      ss.Address,
		AuthorIDList: ss.AuthorIDList,
		OrderBy:      ss.OrderBy,
		PageSize:     ss.PageSize,
		Page:         ss.Page,
	})
	if err != nil {
		return 0, nil, err
	}
	authorIds := make([]int, len(students))
	for i := range students {
		authorIds[i] = students[i].AuthorID
	}

	users, err := da.GetUsersModel().SearchUsers(ctx, da.SearchUserCondition{
		IDList: authorIds,
	})
	if err != nil {
		return 0, nil, err
	}

	authorNameMaps := make(map[int]string)
	for i := range users {
		authorNameMaps[users[i].ID] = users[i].Name
	}

	res := make([]*entity.StudentInfo, len(students))
	for i := range students {
		res[i] = &entity.StudentInfo{
			ID:            students[i].ID,
			Name:          students[i].Name,
			Gender:        students[i].Gender,
			Telephone:     students[i].Telephone,
			Address:       students[i].Address,
			AuthorID:      students[i].AuthorID,
			Email:         students[i].Email,
			AuthorName:    authorNameMaps[students[i].AuthorID],
			IntentSubject: strings.Split(students[i].IntentSubject, ","),
			Status:        students[i].Status,
			Note:          students[i].Note,
		}
	}
	return total, res, nil
}

func (s *StudentService) AddStudentNote(ctx context.Context, c entity.AddStudentNoteRequest) error {
	panic("not implemented")
}

var (
	_studentService     *StudentService
	_studentServiceOnce sync.Once
)

func GetStudentService() *StudentService {
	_studentServiceOnce.Do(func() {
		if _studentService == nil {
			_studentService = new(StudentService)
		}
	})
	return _studentService
}
