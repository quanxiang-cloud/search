package service

import (
	"github.com/go-logr/logr"
	"github.com/graphql-go/graphql"
	"github.com/quanxiang-cloud/search/internal/models"
)

var department = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "department",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"isLeader": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

var role = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "role",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var UserInfo = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "user",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"phone": &graphql.Field{
				Type: graphql.String,
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
			"createdAt": &graphql.Field{
				Type: graphql.Int,
			},
			"departments": &graphql.Field{
				Type: graphql.NewList(department),
			},
			"roles": &graphql.Field{
				Type: graphql.NewList(role),
			},
		},
	},
)

var users = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "users",
		Fields: graphql.Fields{
			"total": &graphql.Field{
				Type: graphql.Int,
			},
			"users": &graphql.Field{
				Type: graphql.NewList(UserInfo),
			},
		},
	},
)

type user struct {
	log logr.Logger

	userSchema graphql.Schema
	userRepo   models.UserRepo
}

func (u *user) newSchema() error {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "_queryUsers",
			Fields: graphql.Fields{
				"user": &graphql.Field{
					Type: users,
					Args: newPageFeild(graphql.FieldConfigArgument{
						"name": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"phone": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"email": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"departmentName": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"roleName": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					),

					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						query := &models.SearchUser{}
						err := mapToStruct(query, p.Args)
						if err != nil {
							u.log.Error(err, "bind args")
							return nil, err
						}
						page, size := bindPageSize(p.Args)
						users, total, err := u.userRepo.Search(p.Context,
							query,
							page, size,
						)
						if err != nil {
							u.log.Error(err, "search user")
							return nil, err
						}

						return struct {
							Total int64          `json:"total,omitempty"`
							Users []*models.User `json:"users,omitempty"`
						}{
							Total: total,
							Users: users,
						}, nil
					},
				},
			},
		}),
	})

	if err != nil {
		return err
	}

	u.userSchema = schema
	return nil
}
