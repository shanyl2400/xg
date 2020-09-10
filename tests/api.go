package tests

import (
	"context"
	"xg/da"
	"xg/entity"
)

type APIClient struct {

}

func (a *APIClient) ListAuths(ctx context.Context) ([]*entity.Auth, error) {
	panic("implement me")
}

func (a *APIClient) CreateOrder(ctx context.Context, req *entity.CreateOrderRequest, operator *entity.JWTUser) (int, error) {
	panic("implement me")
}

func (a *APIClient) SignUpOrder(ctx context.Context, req *entity.OrderPayRequest, operator *entity.JWTUser) error {
	panic("implement me")
}

func (a *APIClient) RevokeOrder(ctx context.Context, orderId int, operator *entity.JWTUser) error {
	panic("implement me")
}

func (a *APIClient) PayOrder(ctx context.Context, req *entity.OrderPayRequest, operator *entity.JWTUser) error {
	panic("implement me")
}

func (a *APIClient) PaybackOrder(ctx context.Context, req *entity.OrderPayRequest, operator *entity.JWTUser) error {
	panic("implement me")
}

func (a *APIClient) ConfirmOrderPay(ctx context.Context, orderPayId int, status int, operator *entity.JWTUser) error {
	panic("implement me")
}

func (a *APIClient) AddOrderRemark(ctx context.Context, orderId int, content string, operator *entity.JWTUser) error {
	panic("implement me")
}

func (a *APIClient) SearchOrderPayRecords(ctx context.Context, condition *entity.SearchPayRecordCondition, operator *entity.JWTUser) (*entity.PayRecordInfoList, error) {
	panic("implement me")
}

func (a *APIClient) SearchOrders(ctx context.Context, condition *entity.SearchOrderCondition, operator *entity.JWTUser) (*entity.OrderInfoList, error) {
	panic("implement me")
}

func (a *APIClient) SearchOrderWithAuthor(ctx context.Context, condition *entity.SearchOrderCondition, operator *entity.JWTUser) (*entity.OrderInfoList, error) {
	panic("implement me")
}

func (a *APIClient) SearchOrderWithOrgId(ctx context.Context, condition *entity.SearchOrderCondition, operator *entity.JWTUser) (*entity.OrderInfoList, error) {
	panic("implement me")
}

func (a *APIClient) GetOrderById(ctx context.Context, orderId int, operator *entity.JWTUser) (*entity.OrderInfoWithRecords, error) {
	panic("implement me")
}

func (a *APIClient) CreateOrderSources(ctx context.Context, name string) (int, error) {
	panic("implement me")
}

func (a *APIClient) ListOrderSources(ctx context.Context) ([]*da.OrderSource, error) {
	panic("implement me")
}

func (a *APIClient) CreateOrg(ctx context.Context, req *entity.CreateOrgRequest, operator *entity.JWTUser) (int, error) {
	panic("implement me")
}

func (a *APIClient) CreateOrgWithSubOrgs(ctx context.Context, req *entity.CreateOrgWithSubOrgsRequest, operator *entity.JWTUser) (int, error) {
	panic("implement me")
}

func (a *APIClient) UpdateOrgById(ctx context.Context, req *entity.UpdateOrgRequest, operator *entity.JWTUser) error {
	panic("implement me")
}

func (a *APIClient) RevokeOrgById(ctx context.Context, id int, operator *entity.JWTUser) error {
	panic("implement me")
}

func (a *APIClient) CheckOrgById(ctx context.Context, id, status int, operator *entity.JWTUser) error {
	panic("implement me")
}

func (a *APIClient) GetOrgById(ctx context.Context, orgId int) (*entity.Org, error) {
	panic("implement me")
}

func (a *APIClient) GetOrgSubjectsById(ctx context.Context, orgId int) ([]string, error) {
	panic("implement me")
}

func (a *APIClient) ListOrgs(ctx context.Context) (int, []*entity.Org, error) {
	panic("implement me")
}

func (a *APIClient) ListOrgsByStatus(ctx context.Context, status []int) (int, []*entity.Org, error) {
	panic("implement me")
}

func (a *APIClient) SearchSubOrgs(ctx context.Context, condition da.SearchOrgsCondition) (int, []*entity.Org, error) {
	panic("implement me")
}

func (a *APIClient) CreateRole(ctx context.Context, name string, authList []int) (int, error) {
	panic("implement me")
}

func (a *APIClient) ListRole(ctx context.Context) ([]*entity.Role, error) {
	panic("implement me")
}

func (a *APIClient) SetRoleAuth(ctx context.Context, id int, ids []int) error {
	panic("implement me")
}

func (a *APIClient) GetRoleAuth(ctx context.Context, id int) ([]*entity.Auth, error) {
	panic("implement me")
}

func (a *APIClient) Summary(ctx context.Context) (*entity.SummaryInfo, error) {
	panic("implement me")
}

func (a *APIClient) SearchYearRecords(ctx context.Context, key string) ([]*entity.StatisticRecord, error) {
	panic("implement me")
}

func (a *APIClient) CreateStudent(ctx context.Context, c *entity.CreateStudentRequest, operator *entity.JWTUser) (int, int, error) {
	panic("implement me")
}

func (a *APIClient) UpdateStudent(ctx context.Context, id int, req *entity.UpdateStudentRequest) error {
	panic("implement me")
}

func (a *APIClient) GetStudentById(ctx context.Context, id int, operator *entity.JWTUser) (*entity.StudentInfosWithOrders, error) {
	panic("implement me")
}

func (a *APIClient) SearchPrivateStudents(ctx context.Context, ss *entity.SearchStudentRequest, operator *entity.JWTUser) (int, []*entity.StudentInfo, error) {
	panic("implement me")
}

func (a *APIClient) SearchStudents(ctx context.Context, ss *entity.SearchStudentRequest, operator *entity.JWTUser) (int, []*entity.StudentInfo, error) {
	panic("implement me")
}

func (a *APIClient) AddStudentNote(ctx context.Context, c entity.AddStudentNoteRequest) error {
	panic("implement me")
}

func (a *APIClient) ListSubjects(ctx context.Context, parentID int) ([]*entity.Subject, error) {
	panic("implement me")
}

func (a *APIClient) CreateSubject(ctx context.Context, req entity.CreateSubjectRequest) (int, error) {
	panic("implement me")
}

func (a *APIClient) Login(ctx context.Context, name, password string) (*entity.UserLoginResponse, error) {
	panic("implement me")
}

func (a *APIClient) UpdatePassword(ctx context.Context, newPassword string, operator *entity.JWTUser) error {
	panic("implement me")
}

func (a *APIClient) ResetPassword(ctx context.Context, userId int, operator *entity.JWTUser) error {
	panic("implement me")
}

func (a *APIClient) ListUserAuthority(ctx context.Context, operator *entity.JWTUser) ([]*entity.Auth, error) {
	panic("implement me")
}

func (a *APIClient) ListUsers(ctx context.Context) ([]*entity.UserInfo, error) {
	panic("implement me")
}

func (a *APIClient) CreateUser(ctx context.Context, req *entity.CreateUserRequest) (int, error) {
	panic("implement me")
}
