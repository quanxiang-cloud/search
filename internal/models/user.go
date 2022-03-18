package models

import (
	"context"

	"github.com/quanxiang-cloud/search/pkg/apis/v1alpha1"
)

// UserRepo user interface
type UserRepo interface {
	Get(ctx context.Context, userID string) (*v1alpha1.User, error)
	List(ctx context.Context, userIDs []interface{}) ([]*v1alpha1.User, error)
	Search(ctx context.Context, query *v1alpha1.SearchUser, page, size int) ([]*v1alpha1.User, int64, error)
}
