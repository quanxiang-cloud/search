package service

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-logr/logr"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/search/pkg/util"
)

type Search struct {
	log logr.Logger

	user
	department
}

func NewSearch(ctx context.Context, opts ...Option) (*Search, error) {
	search := &Search{
		log: util.LoggerFromContext(ctx).WithName("search"),
	}

	err := search.newSchema()
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(search)
	}

	return search, nil
}

func (s *Search) newSchema() error {
	s.user.log = s.log.WithName("user")
	s.user.newSchema()
	s.department.log = s.log.WithName("department")
	s.department.newSchema()
	return nil
}

type base struct {
	UserID       string `json:"userID,omitempty"`
	DepartmentID string `json:"departmentID,omitempty"`
	TenantID     string `json:"tenantID,omitempty"`
	Query        string `json:"query,omitempty"`
}

type SearchUserReq struct {
	base
}
type SearchUserResp struct {
	Data interface{}
}

type SearchDepartmentReq struct {
	base
}
type SearchDepartmentResp struct {
	Data interface{}
}

func (s *Search) SearchUser(ctx context.Context, req *SearchUserReq) (*SearchUserResp, error) {
	data, err := s.search(ctx, s.user.querySchema, req.base)
	if err != nil {
		return &SearchUserResp{}, err
	}

	return &SearchUserResp{
		Data: data,
	}, nil
}

func (s *Search) SearchDepartment(ctx context.Context, req *SearchDepartmentReq) (*SearchDepartmentResp, error) {
	data, err := s.search(ctx, s.department.querySchema, req.base)
	if err != nil {
		return &SearchDepartmentResp{}, err
	}

	return &SearchDepartmentResp{
		Data: data,
	}, nil
}

type DepartmentsByIDsReq struct {
	base
}

type DepartmentsByIDsResp struct {
	Data interface{}
}

func (s *Search) DepartmentByIDs(ctx context.Context, req *DepartmentsByIDsReq) (*DepartmentsByIDsResp, interface{}) {
	data, err := s.search(ctx, s.department.queryByIDsSchema, req.base)
	if err != nil {
		return &DepartmentsByIDsResp{}, err
	}

	return &DepartmentsByIDsResp{
		Data: data,
	}, nil
}

type DepartmentMemberReq struct {
	base
}

type DepartmentMemberResp struct {
	Data interface{}
}

func (s *Search) DepartmentMember(ctx context.Context, req *DepartmentMemberReq) (*DepartmentMemberResp, error) {
	data, err := s.search(ctx, s.departmentMemberSchema, req.base)
	if err != nil {
		return &DepartmentMemberResp{}, err
	}

	return &DepartmentMemberResp{
		Data: data,
	}, nil
}

type SubordinateReq struct {
	base
}

type SubordinateResp struct {
	Data interface{}
}

func (s *Search) Subordinate(ctx context.Context, req *SubordinateReq) (*SubordinateResp, error) {
	data, err := s.search(ctx, s.subordinateSchema, req.base)
	if err != nil {
		return &SubordinateResp{}, err
	}

	return &SubordinateResp{
		Data: data,
	}, nil
}

type LeaderReq struct {
	base
}

type LeaderResp struct {
	Data interface{}
}

func (s *Search) Leader(ctx context.Context, req *LeaderReq) (*LeaderResp, error) {
	data, err := s.search(ctx, s.leaderSchema, req.base)
	if err != nil {
		return &LeaderResp{}, err
	}

	return &LeaderResp{
		Data: data,
	}, nil
}

type RoleMemberReq struct {
	base
}

type RoleMemberResp struct {
	Data interface{}
}

func (s *Search) RoleMember(ctx context.Context, req *RoleMemberReq) (*RoleMemberResp, error) {
	data, err := s.search(ctx, s.rolememberSchema, req.base)
	if err != nil {
		return &RoleMemberResp{}, err
	}

	return &RoleMemberResp{
		Data: data,
	}, nil
}

type UserByIDsReq struct {
	base
}

type UserByIDsResp struct {
	Data interface{}
}

func (s *Search) UserByIDs(ctx context.Context, req *UserByIDsReq) (*UserByIDsResp, interface{}) {
	data, err := s.search(ctx, s.user.userByIDsSchema, req.base)
	if err != nil {
		return &UserByIDsResp{}, err
	}

	return &UserByIDsResp{
		Data: data,
	}, nil
}

func (s *Search) search(ctx context.Context, schema graphql.Schema, base base) (interface{}, error) {
	params := graphql.Params{
		Context:       ctx,
		Schema:        schema,
		RequestString: base.Query,
		RootObject: map[string]interface{}{
			"userID":       base.UserID,
			"departmentID": base.DepartmentID,
			"tenantID":     base.TenantID,
		},
	}

	result := graphql.Do(params)
	if len(result.Errors) > 0 {
		logErrors(ctx, s.log, result.Errors...)
		return nil, result.Errors[0]
	}

	return result.Data, nil
}

func logErrors(ctx context.Context, log logr.Logger, errors ...gqlerrors.FormattedError) {
	for _, err := range errors {
		log.Info(err.Message, header.GetRequestIDKV(ctx).Fuzzy()...)
	}
}

func mapToStruct(dst interface{}, src map[string]interface{}) error {
	if reflect.TypeOf(dst).Kind() != reflect.Ptr {
		return fmt.Errorf("dst must ptr")
	}

	body, err := json.Marshal(src)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, dst)
}

const (
	maxSize int = 999
)

func bindPageSize(src map[string]interface{}) (int, int) {
	if src == nil {
		return 1, maxSize
	}

	page, _ := src["page"].(int)
	size, _ := src["size"].(int)

	if size == 0 {
		size = maxSize
	}
	if page < 1 {
		page = 1
	}
	return page, size
}

func newPageFeild(src graphql.FieldConfigArgument) graphql.FieldConfigArgument {
	src["orderBy"] = &graphql.ArgumentConfig{
		Type: graphql.NewScalar(graphql.ScalarConfig{
			Name: "orderBy",
			Serialize: func(value interface{}) interface{} {
				return value
			},
			ParseValue: func(value interface{}) interface{} {
				return value
			},
			ParseLiteral: func(valueAST ast.Value) interface{} {
				switch valueAST := valueAST.(type) {
				case *ast.ListValue:
					ordeyBy := make([]string, 0, len(valueAST.Values))
					for _, value := range valueAST.Values {
						if vs, ok := value.GetValue().([]*ast.ObjectField); ok &&
							len(vs) == 1 {
							name := vs[0].Name.Value
							if vt, ok := vs[0].Value.GetValue().(string); ok &&
								strings.ToUpper(vt) == "ASC" {
								ordeyBy = append(ordeyBy, name)
								continue
							}
							ordeyBy = append(ordeyBy, "-"+name)
						}
					}
					return ordeyBy
				}
				return nil
			},
		}),
	}
	src["page"] = &graphql.ArgumentConfig{
		Type:         graphql.Int,
		DefaultValue: 0,
	}
	src["size"] = &graphql.ArgumentConfig{
		Type:         graphql.Int,
		DefaultValue: 10,
	}

	return src
}
