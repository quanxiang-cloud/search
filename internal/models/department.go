package models

import (
	"context"
	"github.com/quanxiang-cloud/search/pkg/apis/v1alpha1"
)

// DepartmentRepo department interface
type DepartmentRepo interface {
	Get(ctx context.Context, depID string) (*v1alpha1.Department, error)
	List(ctx context.Context, depIDs []interface{}) ([]*v1alpha1.Department, error)
	Search(ctx context.Context, query *v1alpha1.SearchDepartment, page, size int) ([]*v1alpha1.Department, int64, error)
}
