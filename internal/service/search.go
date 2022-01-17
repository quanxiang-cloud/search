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
	"github.com/quanxiang-cloud/search/pkg/util"
)

type Search struct {
	log logr.Logger

	user
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
	return nil
}

type SearchUserReq struct {
	UserID       string
	DepartmentID string
	Query        string
}

type SearchUserResp struct {
	Data interface{}
}

func (s *Search) SearchUser(ctx context.Context, req *SearchUserReq) (*SearchUserResp, error) {
	params := graphql.Params{
		Context:       ctx,
		Schema:        s.userSchema,
		RequestString: req.Query,
		RootObject: map[string]interface{}{
			"userID":       req.UserID,
			"departmentID": req.DepartmentID,
		},
	}

	result := graphql.Do(params)
	if len(result.Errors) > 0 {
		logErrors(s.log, result.Errors...)
		return &SearchUserResp{}, result.Errors[0]
	}

	return &SearchUserResp{
		Data: result.Data,
	}, nil
}

func logErrors(log logr.Logger, errors ...gqlerrors.FormattedError) {
	for _, err := range errors {
		log.Info(err.Message)
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
