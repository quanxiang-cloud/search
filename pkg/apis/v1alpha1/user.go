package v1alpha1

type User struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Email     string `json:"email,omitempty"`
	CreatedAt int64  `json:"createdAt,omitempty"`
	JobNumber string `json:"job_number,omitempty" ` //工号
	Avatar    string `json:"avatar,omitempty" `     //头像
	UseStatus int    `json:"use_status,omitempty" ` //状态：1正常，-2禁用，-3离职，-1删除，2激活==1 （与账号库相同）
	TenantID  string `json:"tenant_id,omitempty" `  //租户id
	Gender    int    `json:"gender,omitempty" `     //用户密码状态：0无，1男，2女
	Source    string `json:"source,omitempty" `     //信息来源
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
