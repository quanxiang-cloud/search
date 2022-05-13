package elasticsearch

import (
	"context"
	"encoding/json"
	"github.com/go-logr/logr"
	"github.com/olivere/elastic/v7"
	"github.com/quanxiang-cloud/search/internal/models"
	"github.com/quanxiang-cloud/search/pkg/apis/v1alpha1"
	"github.com/quanxiang-cloud/search/pkg/util"
	"strings"
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

	if query.Name != "" {
		mustQuery = append(mustQuery, elastic.NewMatchPhrasePrefixQuery("name", query.Name))
	}
	if query.TenantID != "" {
		mustQuery = append(mustQuery, elastic.NewTermQuery("tenantID", query.TenantID))
	} else {
		mustQuery = append(mustQuery, elastic.NewExistsQuery("tenantID"))
	}
	if len(query.Attr) > 0 {
		for k := range query.Attr {
			mustQuery = append(mustQuery, elastic.NewTermQuery("attr", query.Attr[k]))
		}
	}
	ql = ql.Query(elastic.NewBoolQuery().Must(mustQuery...))

	for _, orderBy := range query.OrderBy {
		if strings.HasPrefix(orderBy, "-") {
			ql = ql.Sort(orderBy[1:], true)
			continue
		}
		ql = ql.Sort(orderBy, false)
	}

	ql = ql.Sort("id.keyword", true)

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

func (u *department) List(ctx context.Context, depIDs []interface{}) ([]*v1alpha1.Department, error) {
	var size = 0
	if len(depIDs) > 100 {
		size = 99
	} else {
		size = len(depIDs)
	}
	result, err := u.client.Search().
		Index(u.index()).
		Query(
			elastic.NewTermsQuery("id.keyword", depIDs...),
		).From(0).Size(size).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	deps := make([]*v1alpha1.Department, 0, len(depIDs))
	for _, hit := range result.Hits.Hits {
		dep := new(v1alpha1.Department)
		err := json.Unmarshal(hit.Source, dep)
		if err != nil {
			return nil, err
		}
		deps = append(deps, dep)
	}

	return deps, nil
}
