package inputs

import (
  "time"
  "github.com/johann8384/libbeat/common"
)

type MothershipConfig struct {
  Enabled            bool
  Port               int
  Flush_interval     *int
  Max_retries        *int
  Type               string
  Sleep_interval     int
  Tick_interval      int
  Minor_interval     int
  Major_interval     int
}

// Functions to be exported by an input plugin
type InputInterface interface {
  // Initialize the input plugin
  Init(config MothershipConfig) error
  // Be able to retrieve the config
  GetConfig() MothershipConfig
  InputType() string
  InputVersion() string
  // Run
  Run(chan common.MapStr) error
}

func (l *MothershipConfig) Normalize (global MothershipConfig)  {
  // default to global if none assigned
  if l.Tick_interval == 0 {
    l.Tick_interval = global.Tick_interval
  }
  if l.Minor_interval == 0 {
    l.Minor_interval = global.Minor_interval
  }
  if l.Tick_interval == 0 {
    l.Major_interval = global.Major_interval
  }

  // adjust intervals to minimums if too low
  if l.Tick_interval < 15 {
    l.Tick_interval = 15
  }
  if l.Minor_interval < 60 {
    l.Minor_interval = 60
  }
  if l.Major_interval < 900 {
    l.Major_interval = 900
  }
}

// PeriodTaskRunner
type taskRunner func(chan common.MapStr)

func PeriodicTaskRunner (l InputInterface, output chan common.MapStr, ti taskRunner, mi taskRunner, ma taskRunner) {
  mi_last := time.Now()
  ma_last := time.Now()
  config := l.GetConfig()

  for {
    ti(output)
    time.Sleep(time.Duration(config.Tick_interval) * time.Second)
    if (time.Since(mi_last) > time.Duration(config.Minor_interval) * time.Second) {
      mi(output)
      mi_last = time.Now()
    }
    if (time.Since(ma_last) > time.Duration(config.Major_interval) * time.Second) {
      ma(output)
      ma_last = time.Now()
    }
  }
}

func EmptyFunc (chan common.MapStr) {}
