package service

import (
	"errors"
	"github.com/go-logr/logr"
	"github.com/graphql-go/graphql"
	"github.com/quanxiang-cloud/search/internal/models"
	"github.com/quanxiang-cloud/search/pkg/apis/v1alpha1"
)

var DepartmentInfo = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "department",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"attr": &graphql.Field{
				Type: graphql.Int,
			},
			"pid": &graphql.Field{
				Type: graphql.String,
			},
			"useStatus": &graphql.Field{
				Type: graphql.Int,
			},
			"tenantID": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var departments = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "departments",
		Fields: graphql.Fields{
			"total": &graphql.Field{
				Type: graphql.Int,
			},
			"departments": &graphql.Field{
				Type: graphql.NewList(DepartmentInfo),
			},
		},
	},
)

type department struct {
	log              logr.Logger
	querySchema      graphql.Schema
	queryByIDsSchema graphql.Schema
	depRepo          models.DepartmentRepo
}

func (u *department) newSchema() error {
	err := u.query()
	if err != nil {
		return err
	}
	err = u.getByIDs()
	if err != nil {
		return err
	}
	return nil
}

func (u *department) resolve(p graphql.ResolveParams) (interface{}, error) {
	query := &v1alpha1.SearchDepartment{
		TenantID: p.Source.(map[string]interface{})["tenantID"].(string),
	}
	err := mapToStruct(query, p.Args)
	if err != nil {
		u.log.Error(err, "bind args")
		return nil, err
	}
	page, size := bindPageSize(p.Args)
	deps, total, err := u.depRepo.Search(p.Context,
		query,
		page, size,
	)
	if err != nil {
		u.log.Error(err, "search department")
		return nil, err
	}

	return struct {
		Total      int64                  `json:"total,omitempty"`
		Department []*v1alpha1.Department `json:"departments,omitempty"`
	}{
		Total:      total,
		Department: deps,
	}, nil
}

func (u *department) query() error {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "_queryDepartments",
			Fields: graphql.Fields{
				"query": &graphql.Field{
					Type: departments,
					Args: newPageFeild(graphql.FieldConfigArgument{
						"attr": &graphql.ArgumentConfig{
							Type: graphql.NewList(graphql.Int),
						},
						"name": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					),
					Resolve: u.resolve,
				},
			},
		}),
	})

	if err != nil {
		return err
	}

	u.querySchema = schema

	return nil
}

func (u *department) getByIDs() error {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "_queryDepartmentsByIDs",
			Fields: graphql.Fields{
				"query": &graphql.Field{
					Type: departments,
					Args: graphql.FieldConfigArgument{
						"ids": &graphql.ArgumentConfig{
							Type: graphql.NewList(graphql.String),
						},
					},
					Resolve: u.getByIDsResolve,
				},
			},
		}),
	})
	if err != nil {
		return err
	}
	u.queryByIDsSchema = schema

	return nil
}

func (u *department) getByIDsResolve(p graphql.ResolveParams) (interface{}, error) {
	ids, ok := p.Args["ids"].([]interface{})
	if !ok {
		return nil, errors.New("invalid id type")
	}
	list, err := u.depRepo.List(p.Context, ids)
	if err != nil {
		u.log.Error(err, "search department")
		return nil, err
	}
	return struct {
		Departments []*v1alpha1.Department `json:"departments,omitempty"`
		Total       int                    `json:"total,omitempty"`
	}{
		Departments: list,
		Total:       len(list),
	}, nil
}
