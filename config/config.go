package config

import (
	"github.com/johann8384/libbeat/common/droppriv"
  "github.com/johann8384/libbeat/logp"
	"github.com/johann8384/libbeat/outputs"
	"github.com/johann8384/libbeat/publisher"
)

type Config struct {
	Output     map[string]outputs.MothershipConfig
	Shipper    publisher.ShipperConfig
	RunOptions droppriv.RunOptions
  Logging    logp.Logging
	Filter     map[string]interface{}
}

// Config Singleton
var ConfigSingleton Config
