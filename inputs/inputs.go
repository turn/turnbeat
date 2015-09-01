package inputs

import (
  "github.com/johann8384/libbeat/common"
)

type MothershipConfig struct {
  Enabled            bool
  Port               int
  Flush_interval     *int
  Max_retries        *int
  Type               string
}

// Functions to be exported by an input plugin
type InputInterface interface {
  // Initialize the input plugin
  Init(config MothershipConfig) error
  InputType() string
  InputVersion() string
  // Run
  Run(chan common.MapStr) error
}

// Input identifier
type InputPlugin uint16

// Input constants
const (
  Unknowninput InputPlugin = iota
  FileInput
  StdInput
  TcpInput
  SyslogInput
  ProcfsInput
)

// Input names
var InputNames = []string{
  "unknown",
  "file",
  "stdin",
  "tcp",
  "syslog",
  "procfs",
}

func (i InputPlugin) String() string {
  if int(i) >= len(InputNames) {
    return "impossible"
  }
  return InputNames[i]
}
