package elasticsearch

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/go-logr/logr"
	"github.com/olivere/elastic/v7"
	"github.com/quanxiang-cloud/search/internal/models"
	"github.com/quanxiang-cloud/search/pkg/util"
)

type user struct {
	log    logr.Logger
	client *elastic.Client
}

func NewUser(ctx context.Context, client *elastic.Client) models.UserRepo {
	return &user{
		log:    util.LoggerFromContext(ctx).WithName("user"),
		client: client,
	}
}

func (u *user) index() string {
	return "user"
}

func (u *user) Search(ctx context.Context, query *models.SearchUser, page, size int) ([]*models.User, int64, error) {
	boolQuery := elastic.NewBoolQuery()

	mustQuery := make([]elastic.Query, 0)
	if query.Name != "" {
		mustQuery = append(mustQuery, elastic.NewMatchQuery("name", query.Name))
	}
	if query.Phone != "" {
		mustQuery = append(mustQuery, elastic.NewMatchPhrasePrefixQuery("phone", query.Phone))
	}
	if query.Email != "" {
		mustQuery = append(mustQuery, elastic.NewMatchPhrasePrefixQuery("email", query.Email))
	}
	if query.DepartmentName != "" {
		mustQuery = append(mustQuery, elastic.NewMatchQuery("departments.name", query.DepartmentName))
	}
	if query.RoleName != "" {
		mustQuery = append(mustQuery, elastic.NewMatchQuery("roles.name", query.RoleName))
	}

	boolQuery = boolQuery.Must(mustQuery...)

	ql := u.client.Search().
		Index(u.index()).
		Query(boolQuery)

	for _, orderBy := range query.OrderBy {
		if strings.HasPrefix(orderBy, "-") {
			ql = ql.Sort(orderBy[1:], true)
			continue
		}
		ql = ql.Sort(orderBy, false)
	}

	result, err := ql.From((page - 1) * size).Size(size).
		Do(ctx)
	if err != nil {
		u.log.Error(err, "user search")
		return nil, 0, err
	}

	var total int64
	users := make([]*models.User, 0, size)
	if result.Hits != nil {
		for _, hit := range result.Hits.Hits {
			user := new(models.User)
			err := json.Unmarshal(hit.Source, user)
			if err != nil {
				return nil, 0, err
			}
			users = append(users, user)
		}
		total = result.Hits.TotalHits.Value
	}

	return users, total, nil
}
