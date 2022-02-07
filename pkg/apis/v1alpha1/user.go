package v1alpha1

type User struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Email     string `json:"email,omitempty"`
	CreatedAt int64  `json:"createdAt,omitempty"`
	JobNumber string `json:"job_number,omitempty" `
	Avatar    string `json:"avatar,omitempty" `
	//UseStatus 1:normal,-2:disable，-quit，-1：del,2:active
	UseStatus int `json:"use_status,omitempty" `
	// TenantID tenant id
	TenantID string `json:"tenant_id,omitempty" `
	// Gender 1：man,2woman
	Gender int `json:"gender,omitempty" `
	// Source where the info come from
	Source    string `json:"source,omitempty" `
	SelfEmail string `json:"self_email,omitempty" `
	// Departments arranged from the current user's department
	// to the top-level department.
	Departments [][]Department `json:"departments,omitempty"`

	// Leaders recent leader to top leader.
	Leaders [][]Leader `json:"leaders,omitempty"`

	Roles []Role `json:"roles,omitempty"`
}

type Department struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	PID  string `json:"pid,omitempty"`
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
