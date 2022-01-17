package api

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/search/internal/service"
)

type search struct {
	s *service.Search
}

func (s *search) SearchUser(c *gin.Context) {
	query := c.Query("query")
	result, err := s.s.SearchUser(header.MutateContext(c), &service.SearchUserReq{
		Query: query,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, transform(result.Data, "user"))
}

func transform(data interface{}, name string) map[string]interface{} {
	result := map[string]interface{}{
		"code": 0,
	}
	if reflect.TypeOf(data).Kind() == reflect.Map {
		value := reflect.ValueOf(data).MapIndex(reflect.ValueOf(name))
		if value.CanInterface() {
			result["data"] = value.Interface()
		}
	}

	return result
}
