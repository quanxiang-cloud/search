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
	req.TenantID = c.GetHeader("Tenant-Id")

	req.Query = query
	result, err := s.s.SearchUser(header.MutateContext(c), req)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, transform(result.Data, "query"))
}

func (s *search) DepartmentMember(c *gin.Context) {
	query := c.Query("query")

	req := &service.DepartmentMemberReq{}
	req.TenantID = c.GetHeader("Tenant-Id")

	req.Query = query
	result, err := s.s.DepartmentMember(header.MutateContext(c), req)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, transform(result.Data, "query"))
}

func (s *search) Subordinate(c *gin.Context) {
	query := c.Query("query")

	req := &service.SubordinateReq{}
	req.UserID = c.GetHeader("User-Id")
	req.TenantID = c.GetHeader("Tenant-Id")

	req.Query = query
	result, err := s.s.Subordinate(header.MutateContext(c), req)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, transform(result.Data, "query"))
}

func (s *search) Leader(c *gin.Context) {
	query := c.Query("query")

	req := &service.LeaderReq{}
	req.UserID = c.GetHeader("User-Id")
	req.TenantID = c.GetHeader("Tenant-Id")

	req.Query = query
	result, err := s.s.Leader(header.MutateContext(c), req)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, transform(result.Data, "query"))
}

func (s *search) RoleMember(c *gin.Context) {
	query := c.Query("query")

	req := &service.RoleMemberReq{}
	req.TenantID = c.GetHeader("Tenant-Id")

	req.Query = query
	result, err := s.s.RoleMember(header.MutateContext(c), req)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, transform(result.Data, "query"))
}

func (s *search) UserByIDs(c *gin.Context) {
	query := c.Query("query")

	req := &service.UserByIDsReq{}
	req.TenantID = c.GetHeader("Tenant-Id")

	req.Query = query
	result, err := s.s.UserByIDs(header.MutateContext(c), req)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, transform(result.Data, "query"))

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

func (s *search) SearchDepartment(c *gin.Context) {
	query := c.Query("query")

	req := &service.SearchDepartmentReq{}
	req.TenantID = c.GetHeader("Tenant-Id")

	req.Query = query
	result, err := s.s.SearchDepartment(header.MutateContext(c), req)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, transform(result.Data, "query"))
}

func (s *search) DepartmentsByIDs(c *gin.Context) {
	query := c.Query("query")

	req := &service.DepartmentsByIDsReq{}
	req.TenantID = c.GetHeader("Tenant-Id")

	req.Query = query
	result, err := s.s.DepartmentByIDs(header.MutateContext(c), req)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, transform(result.Data, "query"))

}
