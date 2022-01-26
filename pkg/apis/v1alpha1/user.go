package v1alpha1

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
	Leaders []Leader `json:"leaders,omitempty"`

	Roles []Role `json:"roles,omitempty"`
}

type Department struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Leader struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
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
