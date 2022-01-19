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

	req := &service.SearchUserReq{}

	req.Query = query
	result, err := s.s.SearchUser(header.MutateContext(c), req)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, transform(result.Data, "user"))
}

func (s *search) DepartmentMember(c *gin.Context) {
	query := c.Query("query")

	req := &service.DepartmentMemberReq{}

	req.Query = query
	result, err := s.s.DepartmentMember(header.MutateContext(c), req)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, transform(result.Data, "department"))
}

func (s *search) Subordinate(c *gin.Context) {
	query := c.Query("query")

	req := &service.SubordinateReq{}
	req.UserID = c.GetHeader("User-Id")

	req.Query = query
	result, err := s.s.Subordinate(header.MutateContext(c), req)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, transform(result.Data, "subordinate"))
}

func (s *search) Leader(c *gin.Context) {
	query := c.Query("query")

	req := &service.LeaderReq{}
	req.UserID = c.GetHeader("User-Id")

	req.Query = query
	result, err := s.s.Leader(header.MutateContext(c), req)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, transform(result.Data, "leader"))
}

func (s *search) RoleMember(c *gin.Context) {
	query := c.Query("query")

	req := &service.RoleMemberReq{}

	req.Query = query
	result, err := s.s.RoleMember(header.MutateContext(c), req)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, transform(result.Data, "role"))
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
