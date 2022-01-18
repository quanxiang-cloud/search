package service

import (
	"fmt"

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

	querySchema            graphql.Schema
	departmentMemberSchema graphql.Schema
	subordinateSchema      graphql.Schema
	leaderSchema           graphql.Schema
	rolememberSchema       graphql.Schema
	postmemberSchema       graphql.Schema

	userRepo models.UserRepo
}

func (u *user) newSchema() error {
	err := u.query()
	if err != nil {
		return err
	}
	err = u.departmentMember()
	if err != nil {
		return err
	}
	err = u.query()
	if err != nil {
		return err
	}
	err = u.subordinate()
	if err != nil {
		return err
	}
	err = u.leader()
	if err != nil {
		return err
	}
	err = u.roleMember()
	if err != nil {
		return err
	}
	err = u.postmember()
	if err != nil {
		return err
	}
	return nil
}

func (u *user) resolve(p graphql.ResolveParams) (interface{}, error) {
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
}

func (u *user) query() error {
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

func (u *user) departmentMember() error {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "_departmentMember",
			Fields: graphql.Fields{
				"department": &graphql.Field{
					Type: users,
					Args: newPageFeild(graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					),
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						if p.Args["id"] == "" {
							return nil, fmt.Errorf("id is must")
						}

						// rename
						p.Args["departmentID"] = p.Args["id"]
						delete(p.Args, "id")

						return u.resolve(p)
					},
				},
			},
		}),
	})

	if err != nil {
		return err
	}

	u.departmentMemberSchema = schema
	return nil
}

func (u *user) subordinate() error {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "_subordinate",
			Fields: graphql.Fields{
				"subordinate": &graphql.Field{
					Type: users,
					Args: newPageFeild(graphql.FieldConfigArgument{}),
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						p.Args["leaderID"] = p.Source.(map[string]interface{})["userID"]
						return u.resolve(p)
					},
				},
			},
		}),
	})

	if err != nil {
		return err
	}

	u.subordinateSchema = schema
	return nil
}

func (u *user) leader() error {
	return nil
}

func (u *user) roleMember() error {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "_roleMember",
			Fields: graphql.Fields{
				"role": &graphql.Field{
					Type: users,
					Args: newPageFeild(graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"name": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					),
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						// rename
						p.Args["roleID"] = p.Args["id"]
						delete(p.Args, "id")

						p.Args["roleName"] = p.Args["name"]
						delete(p.Args, "name")

						return u.resolve(p)
					},
				},
			},
		}),
	})

	if err != nil {
		return err
	}

	u.rolememberSchema = schema
	return nil
}

// TODO
func (u *user) postmember() error {
	return nil
}
