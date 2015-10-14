package config

import (
	"github.com/johann8384/libbeat/common/droppriv"
	"github.com/johann8384/libbeat/logp"
	"github.com/johann8384/libbeat/outputs"
	"github.com/johann8384/libbeat/publisher"
	"github.com/turn/turnbeat/inputs"
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
