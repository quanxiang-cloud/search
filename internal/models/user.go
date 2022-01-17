package models

import "context"

type User struct {
	ID          string       `json:"_id,omitempty"`
	Name        string       `json:"name,omitempty"`
	Phone       string       `json:"phone,omitempty"`
	Email       string       `json:"email,omitempty"`
	CreatedAt   int64        `json:"createdAt,omitempty"`
	Departments []Department `json:"departments,omitempty"`
	Roles       []Role       `json:"roles,omitempty"`
}

type Department struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	IsLeader int    `json:"isLeader,omitempty"`
}

type Role struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type SearchUser struct {
	Name           string   `json:"name,omitempty"`
	Phone          string   `json:"phone,omitempty"`
	Email          string   `json:"email,omitempty"`
	DepartmentName string   `json:"departmentName,omitempty"`
	RoleName       string   `json:"roleName,omitempty"`
	OrderBy        []string `json:"orderBy,omitempty"`
}

type UserRepo interface {
	Search(ctx context.Context, query *SearchUser, page, size int) ([]*User, int64, error)
}
