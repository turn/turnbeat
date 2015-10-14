package config

import (
	"github.com/blacklightops/libbeat/common/droppriv"
	"github.com/blacklightops/libbeat/logp"
	"github.com/blacklightops/libbeat/outputs"
	"github.com/blacklightops/libbeat/publisher"
	"github.com/blacklightops/turnbeat/inputs"
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
