package service

import (
	"bytes"
	"context"
	"strings"
	"sync"
	"time"
	"xg/da"
	"xg/db"
	"xg/entity"
	"xg/log"
	"xg/utils"

	"github.com/jinzhu/gorm"
)

const (
	ChallengeDate = 60 * time.Hour * 24 * 6
)

type IStudentService interface {
	CreateStudent(ctx context.Context, c *entity.CreateStudentRequest, operator *entity.JWTUser) (int, int, error)
	UpdateStudent(ctx context.Context, id int, req *entity.UpdateStudentRequest) error
	GetStudentById(ctx context.Context, id int, operator *entity.JWTUser) (*entity.StudentInfosWithOrders, error)
	SearchPrivateStudents(ctx context.Context, ss *entity.SearchStudentRequest, operator *entity.JWTUser) (int, []*entity.StudentInfo, error)
	SearchStudents(ctx context.Context, ss *entity.SearchStudentRequest, operator *entity.JWTUser) (int, []*entity.StudentInfo, error)
	AddStudentNote(ctx context.Context, c entity.AddStudentNoteRequest) error

	UpdateStudentOrderCount(ctx context.Context, tx *gorm.DB, id int, count int) error
}

type StudentService struct {
	sync.Mutex
}

func (s *StudentService) CreateStudent(ctx context.Context, c *entity.CreateStudentRequest, operator *entity.JWTUser) (int, int, error) {
	status := entity.StudentCreated
	log.Info.Printf("CreateStudent, req: %#v\n", c)
	//检查冲单，查询相同手机号是否有学生
	condition := da.SearchStudentCondition{
		Telephone: c.Telephone,
		OrderBy:   "created_at",
		Status:    []int{entity.StudentCreated},
		PageSize:  1,
		Page:      0,
	}
	total, _, err := da.GetStudentModel().SearchStudents(ctx, condition)
	if err != nil {
		log.Warning.Printf("SearchStudents failed, condition: %#v, err: %v\n", condition, err)
		return -1, -1, err
	}
	//找到对应的学生
	if total > 0 {
		status = entity.StudentConflictFailed
	}
	//获取经纬度信息
	if c.Longitude == 0 && c.Latitude == 0 {
		cor, err := utils.GetAddressLocation(c.Address + c.AddressExt)
		if err != nil {
			log.Warning.Printf("Get address failed, req: %#v, err: %v\n", c, err)
		} else {
			c.Latitude = cor.Latitude
			c.Longitude = cor.Longitude
		}
	}
	student := da.Student{
		Name:           c.Name,
		Gender:         c.Gender,
		Telephone:      c.Telephone,
		Address:        c.Address,
		AddressExt:     c.AddressExt,
		Email:          c.Email,
		IntentSubject:  strings.Join(c.IntentSubject, ","),
		Status:         status,
		AuthorID:       operator.UserId,
		OrderSourceID:  c.OrderSourceID,
		OrderSourceExt: c.OrderSourceExt,
		Latitude:       c.Latitude,
		Longitude:      c.Longitude,
		Note:           c.Note,
	}
	log.Info.Printf("create student, student: %#v, err: %v\n", student, err)

	s.Lock()
	defer s.Unlock()
	id, err := db.GetTransResult(ctx, func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
		//若存在冲突，创建冲突记录
		if status == entity.StudentConflictFailed {
			req := entity.CreateStudentConflictRequest{
				Telephone: c.Telephone,
				AuthorID:  operator.UserId,
				Total:     total + 1,
			}
			err = GetStudentConflictService().CreateOrUpdateStudentConflict(ctx, tx, req)
			if err != nil {
				log.Warning.Printf("Create student conflict failed, student: %#v, err: %v\n", req, err)
				return -1, err
			}
		}

		//添加学生记录
		id, err := da.GetStudentModel().CreateStudent(ctx, tx, student)
		if err != nil {
			log.Warning.Printf("Create student failed, student: %#v, err: %v\n", student, err)
			return -1, err
		}
		//err = GetStatisticsService().AddStudent(ctx, tx, operator.UserId, 1)
		//if err != nil {
		//	log.Warning.Printf("Add student statistics failed, student: %#v, err: %v\n", student, err)
		//	return -1, err
		//}

		err = GetOrderStatisticsService().AddStudent(ctx, tx, operator.UserId, student.OrderSourceID)
		if err != nil {
			log.Warning.Printf("Add student new statistics failed, student: %#v, err: %v\n", student, err)
			return -1, err
		}
		return id, nil
	})
	if err != nil {
		return -1, -1, err
	}
	return id.(int), status, nil
}

func (s *StudentService) UpdateStudentOrderCount(ctx context.Context, tx *gorm.DB, id int, count int) error {
	log.Info.Printf("UpdateStudentOrderCount, id: %#v, count: %#v\n", id, count)
	//获取经纬度信息
	student, err := da.GetStudentModel().GetStudentById(ctx, id)
	if err != nil {
		log.Warning.Printf("Get student failed, id: %#v, err: %v\n", id, err)
		return err
	}
	data := da.Student{
		OrderCount: student.OrderCount + count,
	}
	err = da.GetStudentModel().UpdateStudent(ctx, tx, id, data)
	if err != nil {
		log.Warning.Printf("Update student failed, id: %#v, student: %#v, err: %v\n", id, student, err)
		return err
	}
	return nil
}

func (s *StudentService) UpdateStudent(ctx context.Context, id int, req *entity.UpdateStudentRequest) error {
	log.Info.Printf("UpdateStudent, id: %#v, req: %#v\n", id, req)
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
	data := da.Student{
		Name:           req.Name,
		Gender:         req.Gender,
		Telephone:      req.Telephone,
		Address:        req.Address,
		AddressExt:     req.AddressExt,
		Email:          req.Email,
		IntentSubject:  strings.Join(req.IntentSubject, ","),
		OrderSourceID:  req.OrderSourceID,
		OrderSourceExt: req.OrderSourceExt,
		Longitude:      req.Longitude,
		Latitude:       req.Latitude,
	}
	err := da.GetStudentModel().UpdateStudent(ctx, db.Get(), id, data)
	if err != nil {
		log.Warning.Printf("Update student failed, req: %#v, data: %#v, err: %v\n", req, data, err)
		return err
	}
	return nil
}

func (s *StudentService) GetStudentById(ctx context.Context, id int, operator *entity.JWTUser) (*entity.StudentInfosWithOrders, error) {
	//查询学生信息
	log.Info.Printf("GetStudentById, id: %#v\n", id)

	student, err := da.GetStudentModel().GetStudentById(ctx, id)
	if err != nil {
		log.Warning.Printf("Get student, id: %#v, err: %v\n", id, err)
		return nil, err
	}

	//查询相关订单
	_, orders, err := da.GetOrderModel().SearchOrder(ctx, da.SearchOrderCondition{
		StudentIDList: []int{student.ID},
	})
	if err != nil {
		log.Warning.Printf("SearchOrder, student: %#v, err: %v\n", student, err)
		return nil, err
	}

	//查询用户和机构信息
	orgIds := make([]int, len(orders))
	publisherIds := make([]int, len(orders))
	for i := range orders {
		orgIds[i] = orders[i].ToOrgID
		publisherIds[i] = orders[i].PublisherID
	}
	_, publishers, err := da.GetUsersModel().SearchUsers(ctx, da.SearchUserCondition{
		IDList: utils.UniqueInts(publisherIds),
	})
	if err != nil {
		log.Warning.Printf("SearchUsers failed, publisherIds: %#v, err: %v\n", publisherIds, err)
		return nil, err
	}
	orgs, err := da.GetOrgModel().ListOrgs(ctx)
	if err != nil {
		log.Warning.Printf("ListOrgs failed, err: %v\n", err)
		return nil, err
	}

	orderSources, err := da.GetOrderSourceModel().ListOrderSources(ctx)
	if err != nil {
		log.Warning.Printf("Get OrderSources failed, err: %v\n", err)
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
		log.Warning.Printf("Get User failed, student: %#v, err: %v\n", student, err)
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
			AddressExt:      student.AddressExt,
			AuthorID:        student.AuthorID,
			Email:           student.Email,
			AuthorName:      user.Name,
			OrderSourceID:   student.OrderSourceID,
			OrderSourceName: orderSourceMap[student.OrderSourceID],
			OrderSourceExt:  student.OrderSourceExt,
			IntentSubject:   strings.Split(student.IntentSubject, ","),
			Status:          student.Status,
			Note:            student.Note,
			OrderCount:      student.OrderCount,
			CreatedAt:       student.CreatedAt,
			UpdatedAt:       student.UpdatedAt,
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
				CreatedAt:     orders[i].CreatedAt,
				UpdatedAt:     orders[i].UpdatedAt,
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
	log.Info.Printf("SearchPrivateStudents, condition: %#v\n", ss)
	return s.SearchStudents(ctx, ss, operator)
}

func (s *StudentService) ExportStudents(ctx context.Context, ss *entity.SearchStudentRequest, operator *entity.JWTUser) ([]byte, error) {
	_, data, err := s.SearchStudents(ctx, ss, operator)
	if err != nil {
		log.Warning.Printf("SearchStudents failed, condition: %#v, err: %v\n", ss, err)
		return nil, err
	}
	file, err := StudentsToXlsx(ctx, data)
	if err != nil {
		log.Warning.Printf("Build xlsx file failed, condition: %#v, orderInfos: %#v, err: %v\n", ss, data, err)
		return nil, err
	}
	buf := &bytes.Buffer{}
	file.Write(buf)
	return buf.Bytes(), nil
}

func (s *StudentService) SearchStudents(ctx context.Context, ss *entity.SearchStudentRequest, operator *entity.JWTUser) (int, []*entity.StudentInfo, error) {
	log.Info.Printf("SearchStudents, condition: %#v\n", ss)
	condition := da.SearchStudentCondition{
		Name:            ss.Name,
		Telephone:       ss.Telephone,
		Keywords:        ss.Keywords,
		Address:         ss.Address,
		AuthorIDList:    ss.AuthorIDList,
		IntentString:    ss.IntentSubject,
		Status:          ss.Status,
		NoDispatchOrder: ss.NoDispatchOrder,
		OrderSourceIDs:  ss.OrderSourceIDs,
		CreatedStartAt:  ss.CreatedStartAt,
		CreatedEndAt:    ss.CreatedEndAt,
		OrderBy:         ss.OrderBy,
		PageSize:        ss.PageSize,
		Page:            ss.Page,
	}
	total, students, err := da.GetStudentModel().SearchStudents(ctx, condition)
	if err != nil {
		log.Warning.Printf("SearchStudents failed, condition: %#v, req: %#v, err: %v\n", condition, ss, err)
		return 0, nil, err
	}
	authorIds := make([]int, len(students))
	for i := range students {
		authorIds[i] = students[i].AuthorID
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

	res := make([]*entity.StudentInfo, len(students))
	for i := range students {
		res[i] = &entity.StudentInfo{
			ID:            students[i].ID,
			Name:          students[i].Name,
			Gender:        students[i].Gender,
			Telephone:     students[i].Telephone,
			Address:       students[i].Address,
			AddressExt:    students[i].AddressExt,
			AuthorID:      students[i].AuthorID,
			Email:         students[i].Email,
			AuthorName:    authorNameMaps[students[i].AuthorID],
			IntentSubject: strings.Split(students[i].IntentSubject, ","),
			Status:        students[i].Status,
			Note:          students[i].Note,
			OrderCount:    students[i].OrderCount,
			CreatedAt:     students[i].CreatedAt,
			UpdatedAt:     students[i].UpdatedAt,
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
