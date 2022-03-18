package v1alpha1

type User struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Email     string `json:"email,omitempty"`
	CreatedAt int64  `json:"createdAt,omitempty"`
	JobNumber string `json:"jobNumber,omitempty"`
	Avatar    string `json:"avatar,omitempty"`

	//UseStatus 1:normal,-2:disable，-quit，-1：del,2:active
	UseStatus int `json:"useStatus,omitempty"`

	// TenantID tenant id
	TenantID string `json:"tenantID,omitempty"`

	// Gender 1:man,2: woman
	Gender int `json:"gender,omitempty"`

	// Source where the info come from
	Source    string `json:"source,omitempty"`
	SelfEmail string `json:"selfEmail,omitempty"`
	Position  string `json:"position,omitempty"`

	// Departments arranged from the current user's department
	// to the top-level department.
	Departments [][]Department `json:"departments,omitempty"`

	// Leaders recent leader to top leader.
	Leaders [][]Leader `json:"leaders,omitempty"`

	Roles []Role `json:"roles,omitempty"`
}

type Leader struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Attr string `json:"attr,omitempty"`
}

type Role struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type SearchUser struct {
	TenantID string `json:"tenantID,omitempty"`

	Name      string `json:"name,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Email     string `json:"email,omitempty"`
	JobNumber string `json:"tobNumber,omitempty"`
	Gender    string `json:"gender,omitempty"`
	UseStatus int    `json:"useStatus,omitempty"`

	DepartmentName string `json:"departmentName,omitempty"`
	DepartmentID   string `json:"departmentID,omitempty"`

	RoleID   string `json:"roleID,omitempty"`
	RoleName string `json:"roleName,omitempty"`

	LeaderID string `json:"leaderID,omitempty"`

	OrderBy  []string `json:"orderBy,omitempty"`
	Position string   `json:"position,omitempty"`
}

const UserIndex = "user"
