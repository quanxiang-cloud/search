package models

import "context"

type User struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Email     string `json:"email,omitempty"`
	CreatedAt int64  `json:"createdAt,omitempty"`

	// Departments arranged from the current user's department
	// to the top-level department.
	Departments []Department `json:"departments,omitempty"`

	// Leaders recent leader to top leader.
	Leaders []Leader

	Roles []Role `json:"roles,omitempty"`
}

type Department struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	IsLeader int    `json:"isLeader,omitempty"`
}

type Leader struct {
	ID string
}

type Role struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type SearchUser struct {
	Name  string `json:"name,omitempty"`
	Phone string `json:"phone,omitempty"`
	Email string `json:"email,omitempty"`

	DepartmentName string `json:"departmentName,omitempty"`
	DepartmentID   string `json:"departmentID,omitempty"`

	RoleID   string `json:"roleID,omitempty"`
	RoleName string `json:"roleName,omitempty"`

	LeaderID string `json:"leaderID,omitempty"`

	OrderBy []string `json:"orderBy,omitempty"`
}

type UserRepo interface {
	Search(ctx context.Context, query *SearchUser, page, size int) ([]*User, int64, error)
}
