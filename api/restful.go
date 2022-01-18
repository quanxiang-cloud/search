package api

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/db/elastic"
	ginlogger "github.com/quanxiang-cloud/cabin/tailormade/gin"
	"github.com/quanxiang-cloud/search/internal/config"
	"github.com/quanxiang-cloud/search/internal/service"
	"github.com/quanxiang-cloud/search/pkg/probe"
	"github.com/quanxiang-cloud/search/pkg/util"
)

type Router struct {
	router *gin.Engine
}

func NewRouter(ctx context.Context, conf *config.Config) (*Router, error) {
	e := gin.New()
	e.Use(ginlogger.GinLogger(), ginlogger.GinRecovery())

	log := util.LoggerFromContext(ctx).WithName("router")

	// FIXME logger with logr
	esClient, err := elastic.NewClient(&conf.Elasticsearch, logger.New(nil))
	if err != nil {
		log.Error(err, "new elastic client")
		return nil, err
	}

	v1 := e.Group("/api/v1/search")
	{
		searchService, err := service.NewSearch(ctx,
			service.WithES(ctx, esClient),
		)
		if err != nil {
			log.Error(err, "new user service")
			return nil, err
		}
		s := &search{
			s: searchService,
		}
		v1.GET("/user", s.SearchUser)
		v1.GET("/department/member", s.DepartmentMember)
		v1.GET("/subordinate", s.Subordinate)
		v1.GET("/role/member", s.RoleMember)
	}

	{
		probe := probe.New(util.LoggerFromContext(ctx))
		e.GET("liveness", func(c *gin.Context) {
			probe.LivenessProbe(c.Writer, c.Request)
		})

		e.Any("readiness", func(c *gin.Context) {
			probe.ReadinessProbe(c.Writer, c.Request)
		})

	}

	return &Router{
		router: e,
	}, nil
}

func (r *Router) Run(port string) {
	r.router.Run(port)
}
