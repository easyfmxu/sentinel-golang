package api

import (
	"fmt"

	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/log/metric"
	"github.com/alibaba/sentinel-golang/core/system"
	"github.com/alibaba/sentinel-golang/util"
)

// Initialization func initialize the Sentinel's runtime environment, including:
// 		1. override global config, from manually config or yaml file or env variable
//		2. override global logger
// 		3. initiate core component async task, including: metric log, system statistic...
// InitDefault initializes Sentinel using the configuration from system
// environment and the default value.
func InitDefault() error {
	return initSentinel("")
}

// InitWithConfig initializes Sentinel using given config.
func InitWithConfig(confEntity *config.Entity) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	err = config.CheckValid(confEntity)
	if err != nil {
		return err
	}
	config.SetDefaultConfig(confEntity)
	if err = config.OverrideConfigFromEnvAndInitLog(); err != nil {
		return err
	}
	return initCoreComponents()
}

// Init loads Sentinel general configuration from the given YAML file
// and initializes Sentinel.
func InitWithConfigFile(configPath string) error {
	return initSentinel(configPath)
}

// initCoreComponents init core components with default config
// it's better SetDefaultConfig before initCoreComponents
func initCoreComponents() error {
	if config.MetricLogFlushIntervalSec() > 0 {
		if err := metric.InitTask(); err != nil {
			return err
		}
	}

	if config.SystemStatCollectIntervalMs() > 0 {
		system.InitCollector(config.SystemStatCollectIntervalMs())
	}

	if config.UseCacheTime() {
		util.StartTimeTicker()
	}

	return nil
}

func initSentinel(configPath string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	// Initialize general config and logging module.
	if err = config.InitConfig(configPath); err != nil {
		return err
	}
	return initCoreComponents()
}
