package elasticsearch

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/go-logr/logr"
	"github.com/olivere/elastic/v7"
	"github.com/quanxiang-cloud/search/internal/models"
	"github.com/quanxiang-cloud/search/pkg/apis/v1alpha1"
	"github.com/quanxiang-cloud/search/pkg/util"
)

type department struct {
	log    logr.Logger
	client *elastic.Client
}

func NewDepartment(ctx context.Context, client *elastic.Client) models.DepartmentRepo {
	return &department{
		log:    util.LoggerFromContext(ctx).WithName("department"),
		client: client,
	}
}

func (u *department) index() string {
	return "department"
}

func (u *department) Search(ctx context.Context, query *v1alpha1.SearchDepartment, page, size int) ([]*v1alpha1.Department, int64, error) {
	ql := u.client.Search().Index(u.index())

	mustQuery := make([]elastic.Query, 0)

	mustQuery = append(mustQuery, elastic.NewTermQuery("tenantID", query.TenantID))

	if query.ID != "" {
		mustQuery = append(mustQuery, elastic.NewTermQuery("id", query.ID))
	}
	if query.Name != "" {
		mustQuery = append(mustQuery, elastic.NewMatchPhrasePrefixQuery("name", query.Name))
	}

	ql = ql.Query(elastic.NewBoolQuery().Must(mustQuery...))

	for _, orderBy := range query.OrderBy {
		if strings.HasPrefix(orderBy, "-") {
			ql = ql.Sort(orderBy[1:], true)
			continue
		}
		ql = ql.Sort(orderBy, false)
	}
	ql = ql.Sort("name.keyword", true)

	result, err := ql.From((page - 1) * size).Size(size).
		Do(ctx)

	if err != nil {
		u.log.Error(err, "department search")
		return nil, 0, err
	}

	deps := make([]*v1alpha1.Department, 0, size)
	for _, hit := range result.Hits.Hits {
		dep := new(v1alpha1.Department)
		err := json.Unmarshal(hit.Source, dep)
		if err != nil {
			return nil, 0, err
		}
		deps = append(deps, dep)
	}

	return deps, result.Hits.TotalHits.Value, nil
}
