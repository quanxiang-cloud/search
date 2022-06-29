package v1alpha1

type Department struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	PID      string `json:"pid,omitempty"`
	Attr     string `json:"attr,omitempty"`
	TenantID string `json:"tenantID"`
}

type SearchDepartment struct {
	TenantID string        `json:"tenantID,omitempty"`
	Name     string        `json:"name,omitempty"`
	Attr     []int         `json:"attr,omitempty"`
	IDS      []interface{} `json:"ids,omitempty"`
	OrderBy  []string      `json:"orderBy,omitempty"`
}

const DepartmentIndex = "department"
