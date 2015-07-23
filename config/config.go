package config

import (
	"github.com/johann8384/libbeat/common/droppriv"
  "github.com/johann8384/libbeat/logp"
	"github.com/johann8384/libbeat/outputs"
	"github.com/turn/turnbeat/inputs"
	"github.com/johann8384/libbeat/publisher"
)

type Config struct {
	Input      map[string]inputs.MothershipConfig
	Output     map[string]outputs.MothershipConfig
	Shipper    publisher.ShipperConfig
	RunOptions droppriv.RunOptions
  Logging    logp.Logging
	Filter     map[string]interface{}
}

// Config Singleton
var ConfigSingleton Config
