package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"xg/da"
	"xg/entity"
	"xg/route"
)

type APIClient struct {
}

type SearchPayRecordCondition struct {
	PayRecordIDList []int `json:"pay_record_id_list"`
	OrderIDList     []int `json:"order_id_list"`
	AuthorIDList    []int `json:"author_id_list"`
	Mode            int   `json:"mode"`

	OrderBy string `json:"order_by"`

	PageSize int `json:"page_size"`
	Page     int `json:"page"`
}

func (a *APIClient) ListAuths(ctx context.Context, token string) (*route.AuthsListResponse, error) {
	req := JSONRequest{
		URL:      "/api/auths",
		Method:   "GET",
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.AuthsListResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) CreateOrder(ctx context.Context, crq *entity.CreateOrderRequest, token string) (*route.IdResponse, error) {
	req := JSONRequest{
		URL:      "/api/order",
		Method:   "POST",
		JSONBody: crq,
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.IdResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) SignUpOrder(ctx context.Context, orq *entity.OrderPayRequest, token string) error {
	req := JSONRequest{
		URL:      fmt.Sprintf("/api/order/%v/signup", orq.OrderID),
		Method:   "PUT",
		JSONBody: orq,
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return err
	}
	responseObj := new(route.Response)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return err
	}
	return responseObj.Error()
}

func (a *APIClient) DepositOrder(ctx context.Context, orq *entity.OrderPayRequest, token string) error {
	req := JSONRequest{
		URL:      fmt.Sprintf("/api/order/%v/deposit", orq.OrderID),
		Method:   "PUT",
		JSONBody: orq,
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return err
	}
	responseObj := new(route.Response)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return err
	}
	return responseObj.Error()
}

func (a *APIClient) RevokeOrder(ctx context.Context, orderId int, token string) error {
	req := JSONRequest{
		URL:      fmt.Sprintf("/api/order/%v/revoke", orderId),
		Method:   "PUT",
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return err
	}
	responseObj := new(route.Response)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return err
	}
	return responseObj.Error()
}


func (a *APIClient) InvalidOrder(ctx context.Context, orderId int, token string) error {
	req := JSONRequest{
		URL:      fmt.Sprintf("/api/order/%v/invalid", orderId),
		Method:   "PUT",
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return err
	}
	responseObj := new(route.Response)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return err
	}
	return responseObj.Error()
}

func (a *APIClient) PayOrder(ctx context.Context, orq *entity.OrderPayRequest, token string) error {
	req := JSONRequest{
		URL:      fmt.Sprintf("/api/payment/%v/pay", orq.OrderID),
		Method:   "POST",
		JSONBody: orq,
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return err
	}
	responseObj := new(route.Response)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return err
	}
	return responseObj.Error()
}

func (a *APIClient) PaybackOrder(ctx context.Context, orq *entity.OrderPayRequest, token string) error {
	req := JSONRequest{
		URL:      fmt.Sprintf("/api/payment/%v/payback", orq.OrderID),
		Method:   "POST",
		JSONBody: orq,
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return err
	}
	responseObj := new(route.Response)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return err
	}
	return responseObj.Error()
}

func (a *APIClient) AcceptPayment(ctx context.Context, orderPayId int, token string) error {
	req := JSONRequest{
		URL:      fmt.Sprintf("/api/payment/%v/review/accept", orderPayId),
		Method:   "PUT",
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return err
	}
	responseObj := new(route.Response)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return err
	}
	return responseObj.Error()
}


func (a *APIClient) RejectPayment(ctx context.Context, orderPayId int, token string) error {
	req := JSONRequest{
		URL:      fmt.Sprintf("/api/payment/%v/review/reject", orderPayId),
		Method:   "PUT",
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return err
	}
	responseObj := new(route.Response)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return err
	}
	return responseObj.Error()
}

func (a *APIClient) AddOrderRemark(ctx context.Context, orderId int, content string, token string) error {
	req := JSONRequest{
		URL:      fmt.Sprintf("/api/order/%v/mark", orderId),
		Method:   "POST",
		JSONBody: &entity.OrderMarkRequest{
			Content: content,
		},
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return err
	}
	responseObj := new(route.Response)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return err
	}
	return responseObj.Error()
}

func (a *APIClient) SearchOrderPayRecords(ctx context.Context, condition *SearchPayRecordCondition, token string) (*route.OrderPaymentRecordListResponse, error) {
	query := make(map[string]string)
	query["pay_record_ids"] = buildInts(condition.PayRecordIDList)
	query["order_ids"] = buildInts(condition.OrderIDList)
	query["author_ids"] = buildInts(condition.AuthorIDList)
	query["mode"] = buildInt(condition.Mode)
	query["order_by"] = condition.OrderBy
	query["page"] = buildInt(condition.Page)
	query["page_size"] = buildInt(condition.PageSize)

	req := JSONRequest{
		URL:      "/api/payments/pending",
		Method:   "GET",
		Query: 	query,
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.OrderPaymentRecordListResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) SearchOrders(ctx context.Context, condition *entity.SearchOrderCondition, token string) (*route.OrderInfoListResponse, error) {
	query := make(map[string]string)
	query["student_ids"] = buildInts(condition.StudentIDList)
	query["to_org_ids"] = buildInts(condition.ToOrgIDList)
	query["intent_subjects"] = condition.IntentSubjects
	query["publisher_id"] = buildInts(condition.PublisherID)
	query["order_sources"] = buildInts(condition.OrderSourceList)
	query["status"] = buildInts(condition.Status)
	query["order_by"] = condition.OrderBy
	query["page"] = buildInt(condition.Page)
	query["page_size"] = buildInt(condition.PageSize)
	if condition.CreateStartAt != nil && condition.CreateEndAt != nil {
		query["create_start_at"] = buildInt(int(condition.CreateStartAt.Unix()))
		query["create_end_at"] = buildInt(int(condition.CreateEndAt.Unix()))
	}

	req := JSONRequest{
		URL:      "/api/orders",
		Method:   "GET",
		Query: 	query,
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.OrderInfoListResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) SearchOrderWithAuthor(ctx context.Context, condition *entity.SearchOrderCondition, token string) (*route.OrderInfoListResponse, error) {
	query := make(map[string]string)
	query["student_ids"] = buildInts(condition.StudentIDList)
	query["to_org_ids"] = buildInts(condition.ToOrgIDList)
	query["intent_subjects"] = condition.IntentSubjects
	query["order_sources"] = buildInts(condition.OrderSourceList)
	query["status"] = buildInts(condition.Status)
	query["order_by"] = condition.OrderBy
	query["page"] = buildInt(condition.Page)
	query["page_size"] = buildInt(condition.PageSize)
	if condition.CreateStartAt != nil && condition.CreateEndAt != nil {
		query["create_start_at"] = buildInt(int(condition.CreateStartAt.Unix()))
		query["create_end_at"] = buildInt(int(condition.CreateEndAt.Unix()))
	}

	req := JSONRequest{
		URL:      "/api/orders/author",
		Method:   "GET",
		Query: 	query,
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.OrderInfoListResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) SearchOrderWithOrgId(ctx context.Context, condition *entity.SearchOrderCondition, token string) (*route.OrderInfoListResponse, error) {
	query := make(map[string]string)
	query["student_ids"] = buildInts(condition.StudentIDList)
	query["to_org_ids"] = buildInts(condition.ToOrgIDList)
	query["intent_subjects"] = condition.IntentSubjects
	query["publisher_id"] = buildInts(condition.PublisherID)
	query["order_sources"] = buildInts(condition.OrderSourceList)
	query["status"] = buildInts(condition.Status)
	query["order_by"] = condition.OrderBy
	query["page"] = buildInt(condition.Page)
	query["page_size"] = buildInt(condition.PageSize)
	if condition.CreateStartAt != nil && condition.CreateEndAt != nil {
		query["create_start_at"] = buildInt(int(condition.CreateStartAt.Unix()))
		query["create_end_at"] = buildInt(int(condition.CreateEndAt.Unix()))
	}

	req := JSONRequest{
		URL:      "/api/orders/org",
		Method:   "GET",
		Query: 	query,
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.OrderInfoListResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) GetOrderById(ctx context.Context, orderId int, token string) (*route.OrderRecordResponse, error) {
	req := JSONRequest{
		URL:      fmt.Sprintf("/api/order/%v", orderId),
		Method:   "GET",
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.OrderRecordResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) CreateOrderSources(ctx context.Context, name string, token string) (*route.IdResponse, error) {
	req := JSONRequest{
		URL:      "/api/order_source",
		Method:   "POST",
		JSONBody: entity.CreateOrderSourceRequest{Name: name},
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.IdResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) ListOrderSources(ctx context.Context, token string) (*route.OrderSourcesListResponse, error) {
	req := JSONRequest{
		URL:      "/api/order_sources",
		Method:   "GET",
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.OrderSourcesListResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) CreateOrg(ctx context.Context, orq *entity.CreateOrgWithSubOrgsRequest, token string) (*route.IdResponse, error) {
	req := JSONRequest{
		URL:      "/api/org",
		Method:   "POST",
		JSONBody: orq,
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.IdResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) RevokeOrgById(ctx context.Context, id int, token string) error {
	req := JSONRequest{
		URL:      fmt.Sprintf("/api/org/%v/revoke", id),
		Method:   "PUT",
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return err
	}
	responseObj := new(route.Response)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return err
	}
	return responseObj.Error()
}

func (a *APIClient) UpdateOrgById(ctx context.Context, id int, params entity.UpdateOrgWithSubOrgsRequest, token string) error {
	req := JSONRequest{
		URL:      fmt.Sprintf("/api/org/%v", id),
		Method:   "PUT",
		JSONBody: params,
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return err
	}
	responseObj := new(route.Response)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return err
	}
	return responseObj.Error()
}

func (a *APIClient) UpdateSelfOrgById(ctx context.Context, params entity.UpdateOrgWithSubOrgsRequest, token string) error {
	req := JSONRequest{
		URL:      fmt.Sprintf("/api/org"),
		Method:   "PUT",
		JSONBody: params,
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return err
	}
	responseObj := new(route.Response)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return err
	}
	return responseObj.Error()
}

func (a *APIClient) RejectPendingOrgById(ctx context.Context, id int, token string) error {
	req := JSONRequest{
		URL:      fmt.Sprintf("/api/org/%v/review/reject", id),
		Method:   "PUT",
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return err
	}
	responseObj := new(route.Response)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return err
	}
	return responseObj.Error()
}

func (a *APIClient) ApprovePendingOrgById(ctx context.Context, id int, token string) error {
	req := JSONRequest{
		URL:      fmt.Sprintf("/api/org/%v/review/approve", id),
		Method:   "PUT",
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return err
	}
	responseObj := new(route.Response)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return err
	}
	return responseObj.Error()
}

func (a *APIClient) GetOrgById(ctx context.Context, orgId int, token string) (*route.OrgInfoResponse, error) {
	req := JSONRequest{
		URL:      fmt.Sprintf("/api/org/%v", orgId),
		Method:   "GET",
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.OrgInfoResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) GetOrgSubjectsById(ctx context.Context, orgId int, token string) (*route.OrgSubjectsResponse, error) {
	req := JSONRequest{
		URL:      fmt.Sprintf("/api/org/%v/subjects", orgId),
		Method:   "GET",
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.OrgSubjectsResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) ListOrgs(ctx context.Context, token string) (*route.OrgsListResponse, error) {
	req := JSONRequest{
		URL:      "/api/orgs",
		Method:   "GET",
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.OrgsListResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) ListPendingOrgs(ctx context.Context, token string) (*route.OrgsListResponse, error) {
	req := JSONRequest{
		URL:      "/api/orgs/pending",
		Method:   "GET",
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.OrgsListResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) SearchSubOrgs(ctx context.Context, condition da.SearchOrgsCondition, token string) (*route.SubOrgsListResponse, error) {
	query := make(map[string]string)
	query["subjects"] = buildStrings(condition.Subjects)
	query["address"] = condition.Address
	req := JSONRequest{
		URL:      "/api/orgs/campus",
		Method:   "GET",
		Query: 	query,
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.SubOrgsListResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) CreateRole(ctx context.Context, role *entity.CreateRoleRequest, token string) (*route.IdResponse, error) {
	req := JSONRequest{
		URL:      "/api/role",
		Method:   "POST",
		JSONBody: role,
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.IdResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) ListRole(ctx context.Context, token string) (*route.RolesResponse, error) {
	req := JSONRequest{
		URL:      "/api/roles",
		Method:   "GET",
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.RolesResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) Summary(ctx context.Context, token string) (*route.SummaryResponse, error) {
	req := JSONRequest{
		URL:      "/api/statistics/summary",
		Method:   "GET",
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.SummaryResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) Graph(ctx context.Context, token string) (*route.GraphResponse, error) {
	req := JSONRequest{
		URL:      "/api/statistics/graph",
		Method:   "GET",
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.GraphResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) OrgGraph(ctx context.Context, token string) (*route.PerformanceGraphResponse, error) {
	req := JSONRequest{
		URL:     "/api/statistics/graph/org",
		Method:   "GET",
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.PerformanceGraphResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) DispatchGraph(ctx context.Context, token string) (*route.PerformanceGraphResponse, error) {
	req := JSONRequest{
		URL:      "/api/statistics/graph/dispatch",
		Method:   "GET",
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.PerformanceGraphResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) EnterGraph(ctx context.Context, token string) (*route.AuthorPerformanceGraphResponse, error) {
	req := JSONRequest{
		URL:      "/api/statistics/graph/enter",
		Method:   "GET",
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.AuthorPerformanceGraphResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) CreateStudent(ctx context.Context, c *entity.CreateStudentRequest, token string) (*route.IDStatusResponse, error) {
	req := JSONRequest{
		URL:      "/api/student",
		Method:   "POST",
		JSONBody: c,
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.IDStatusResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) GetStudentById(ctx context.Context, id int, token string) (*route.StudentWithDetailsListResponse, error) {
	req := JSONRequest{
		URL:      fmt.Sprintf("/api/student/%v", id),
		Method:   "GET",
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.StudentWithDetailsListResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) SearchPrivateStudents(ctx context.Context, ss *entity.SearchStudentRequest, token string) (*route.StudentListResponse, error) {
	query := make(map[string]string)
	query["name"] = ss.Name
	query["telephone"] = ss.Telephone
	query["address"] = ss.Address
	query["author_id"] = buildInts(ss.AuthorIDList)
	query["intent_subjects"] = ss.IntentSubject
	query["order_by"] = ss.OrderBy
	query["page"] = buildInt(ss.Page)
	query["page_size"] = buildInt(ss.PageSize)

	req := JSONRequest{
		URL:      "/api/students/private",
		Method:   "GET",
		Query: 	query,
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.StudentListResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) SearchStudents(ctx context.Context, ss *entity.SearchStudentRequest, token string) (*route.StudentListResponse, error) {
	query := make(map[string]string)
	query["name"] = ss.Name
	query["telephone"] = ss.Telephone
	query["address"] = ss.Address
	query["author_id"] = buildInts(ss.AuthorIDList)
	query["intent_subjects"] = ss.IntentSubject
	query["order_by"] = ss.OrderBy
	query["page"] = buildInt(ss.Page)
	query["page_size"] = buildInt(ss.PageSize)

	req := JSONRequest{
		URL:      "/api/students",
		Method:   "GET",
		Query: 	query,
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.StudentListResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) ListSubjects(ctx context.Context, parentID int, token string) (*route.SubjectsObjResponse, error) {
	req := JSONRequest{
		URL:      fmt.Sprintf("/api/subjects/%v", parentID),
		Method:   "GET",
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.SubjectsObjResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) CreateSubject(ctx context.Context, csr entity.CreateSubjectRequest, token string) (*route.IdResponse, error) {
	req := JSONRequest{
		URL:      "/api/subject",
		Method:   "POST",
		JSONBody: csr,
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.IdResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) Login(ctx context.Context, name, password string) (*route.UserLoginResponse, error) {
	req := JSONRequest{
		URL:      "/api/user/login",
		Method:   "POST",
		JSONBody: &entity.UserLoginRequest{
			Name:     name,
			Password: password,
		},
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.UserLoginResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) UpdatePassword(ctx context.Context, newPassword string, token string) error {
	req := JSONRequest{
		URL:      "/api/user/password",
		Method:   "PUT",
		JSONBody: &entity.UserUpdatePasswordRequest{
			NewPassword: newPassword,
		},
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return err
	}
	responseObj := new(route.Response)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return err
	}
	return responseObj.Error()
}

func (a *APIClient) ResetPassword(ctx context.Context, userId int, token string) error {
	req := JSONRequest{
		URL:      fmt.Sprintf("/api/user/reset/%v", userId),
		Method:   "PUT",
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return err
	}
	responseObj := new(route.Response)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return err
	}
	return responseObj.Error()
}

func (a *APIClient) ListUserAuthority(ctx context.Context, token string) (*route.AuthorizationListResponse, error) {
	req := JSONRequest{
		URL:      "/api/user/authority",
		Method:   "GET",
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.AuthorizationListResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) ListUsers(ctx context.Context, token string) (*route.UserListResponse, error) {
	req := JSONRequest{
		URL:      "/api/users",
		Method:   "GET",
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.UserListResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func (a *APIClient) CreateUser(ctx context.Context, cur *entity.CreateUserRequest, token string) (*route.IdResponse, error) {
	req := JSONRequest{
		URL:      "/api/user",
		Method:   "POST",
		JSONBody: cur,
		Token:    token,
	}
	resp, err := req.DoRequest(ctx)
	if err != nil {
		return nil, err
	}
	responseObj := new(route.IdResponse)
	err = json.Unmarshal(resp, responseObj)
	if err != nil {
		return nil, err
	}
	return responseObj, responseObj.Error()
}

func SetBaseURL(url string) {
	BaseURL = url
}

func buildInt(d int) string{
	return strconv.Itoa(d)
}
func buildInts(d []int) string {
	intsStr := make([]string, 0)
	for i := range d{
		intsStr = append(intsStr, strconv.Itoa(d[i]))
	}
	return strings.Join(intsStr, ",")
}
func buildStrings(s []string) string{
	return strings.Join(s, ",")
}