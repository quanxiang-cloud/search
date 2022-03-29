package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/go-logr/zapr"
	"github.com/quanxiang-cloud/search/api"
	"github.com/quanxiang-cloud/search/internal/config"
	"github.com/quanxiang-cloud/search/pkg/util"
	"go.uber.org/zap"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "/configs/config.yml", "config path")
	flag.Parse()

	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
	}
	logger := zapr.NewLogger(zapLog)

	ctx := context.Background()
	ctx = util.SetCtx(ctx, util.ContextKey{}, logger)

	conf, err := config.New(ctx, configFile)
	if err != nil {
		panic(err)
	}

	router, err := api.NewRouter(ctx, conf)
	if err != nil {
		panic(err)
	}

	logger.Info("running...")
	router.Run(conf.Port)
}
