/*
 * Copyright (c) 2018 LI Zhennan
 *
 * Use of this work is governed by an MIT License.
 * You may find a license copy in project root.
 *
 */

// rly is an API for O'RLY cover generation
package main

import (
	"flag"
	"fmt"

	"github.com/nanmu42/orly/cmd/common"

	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"

	"github.com/pkg/errors"
)

var (
	configFile = flag.String("config", "config.toml", "config.toml file location for rly")
	w          WorkerPool
	logger     *zap.Logger
	// Version build params
	Version string
	// BuildDate build params
	BuildDate string
)

func init() {
	w := common.NewBufferedLumberjack(&lumberjack.Logger{
		Filename:   "rly.log",
		MaxSize:    300, // megabytes
		MaxBackups: 5,
		MaxAge:     28, // days
	}, 32*1024)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.Lock(w),
		zap.InfoLevel,
	)
	logger = zap.New(core)
}

func main() {
	var err error
	defer logger.Sync() // nolint: errcheck
	defer func() {
		if err != nil {
			logger.Fatal("fatal error",
				zap.Error(err),
			)
		}
	}()

	flag.Parse()

	fmt.Printf(`O'rly Generator API(%s)
built on %s

`, Version, BuildDate)

	err = C.LoadFrom(*configFile)
	if err != nil {
		err = errors.Wrap(err, "C.LoadFrom")
		return
	}

	err = initializeFactory()
	if err != nil {
		err = errors.Wrap(err, "initializeFactory")
		return
	}

	w = InitWorkerPool(C.WorkerNum, C.QueueLen, makeCover)

	router := setupRouter()
	startAPI(router, C.Port)
}
