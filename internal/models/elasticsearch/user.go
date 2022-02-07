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

func (u *user) Get(ctx context.Context, userID string) (*v1alpha1.User, error) {
	result, err := u.client.Search().
		Index(u.index()).
		Query(
			elastic.NewTermQuery("id", userID),
		).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	if len(result.Hits.Hits) == 0 {
		return nil, nil
	}

	user := new(v1alpha1.User)
	err = json.Unmarshal(result.Hits.Hits[0].Source, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *user) List(ctx context.Context, userIDs []interface{}) ([]*v1alpha1.User, error) {
	result, err := u.client.Search().
		Index(u.index()).
		Query(
			elastic.NewTermsQuery("id", userIDs...),
		).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]*v1alpha1.User, 0, len(userIDs))
	for _, hit := range result.Hits.Hits {
		user := new(v1alpha1.User)
		err := json.Unmarshal(hit.Source, user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (u *user) Search(ctx context.Context, query *v1alpha1.SearchUser, page, size int) ([]*v1alpha1.User, int64, error) {
	ql := u.client.Search().Index(u.index())

	switch {
	case query.DepartmentID != "":
		ql = ql.Query(elastic.NewTermQuery("departments.id", query.DepartmentID))
	case query.RoleID != "":
		ql = ql.Query(elastic.NewTermQuery("roles.id", query.RoleID))
	case query.LeaderID != "":
		ql = ql.Query(elastic.NewTermQuery("leaders.id", query.LeaderID))
	default:
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
		ql = ql.Query(boolQuery)
	}

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
		u.log.Error(err, "user search")
		return nil, 0, err
	}

	users := make([]*v1alpha1.User, 0, size)
	for _, hit := range result.Hits.Hits {
		user := new(v1alpha1.User)
		err := json.Unmarshal(hit.Source, user)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, result.Hits.TotalHits.Value, nil
}
