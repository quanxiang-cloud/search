package service

import (
	"context"

	"github.com/olivere/elastic/v7"
	"github.com/quanxiang-cloud/search/internal/models/elasticsearch"
)

// Option option
type Option func(*Search)

func WithES(ctx context.Context, client *elastic.Client) Option {
	return func(s *Search) {
		s.userRepo = elasticsearch.NewUser(ctx, client)
		s.depRepo = elasticsearch.NewDepartment(ctx, client)
	}
}
