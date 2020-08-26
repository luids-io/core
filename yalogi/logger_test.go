// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package yalogi_test

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/luids-io/core/yalogi"
)

func ExampleLogger() {
	withYalogi := func(logger yalogi.Logger) {
		logger.Debugf("esto es debug")
		logger.Infof("esto es informativo")
		logger.Warnf("esto es una advertencia")
	}

	//instantiate the logrus logger
	logger := logrus.New()
	formatter := &logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	}
	logger.SetFormatter(formatter)
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(os.Stdout)

	//we can use logrus in func because satisfaces the interface yalogi
	withYalogi(logger)

	//we can convert a yalogi to a standar log
	standarlog := yalogi.NewStandard(logger, yalogi.Info)
	standarlog.Printf("mensaje a un log estandar %v informativo\n", 1234)

	// Output:
	//level=debug msg="esto es debug"
	//level=info msg="esto es informativo"
	//level=warning msg="esto es una advertencia"
	//level=info msg="mensaje a un log estandar 1234 informativo\n"
}
