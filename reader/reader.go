package reader

import (
  "encoding/json"
  "errors"
  "strings"
  "fmt"
  "github.com/johann8384/libbeat/common"
  "github.com/johann8384/libbeat/logp"
  "github.com/turn/turnbeat/inputs"
  "github.com/turn/turnbeat/inputs/tcp"
  "github.com/turn/turnbeat/inputs/syslog"
  "github.com/turn/turnbeat/inputs/procfs"
  "github.com/turn/turnbeat/inputs/null"
)

type ReaderType struct {
  name           string
  tags           []string
  disabled       bool
  Index          string
  Input          []inputs.InputInterface
  Queue          chan common.MapStr
}

var Reader ReaderType

/***********************************/
/* register new input plugins here */
/***********************************/
func newInputInstance(name string) inputs.InputInterface {
  logp.Info("creating new instance of type %s", name)
  switch name {
  case "tcp":
    return new(tcp.TcpInput)
  case "syslog":
    return new(syslog.SyslogInput)
  case "procfs":
    return new(procfs.ProcfsInput)
  case "null":
    return new(null.NullInput)
  }
  return nil
}

func (reader *ReaderType) PrintReaderEvent(event common.MapStr) {
  json, err := json.MarshalIndent(event, "", "  ")
  if err != nil {
    logp.Err("json.Marshal: %s", err)
  } else {
    logp.Debug("reader", "Reader: %s", string(json))
  }
}

func (reader *ReaderType) Init(inputs map[string]inputs.MothershipConfig) error {
  logp.Info("reader input config", inputs)

  for inputId, config := range inputs {
    // default instance 0
    inputName, instance := inputId, "0"
    if (strings.Contains(inputId, "_")) {
      // otherwise grok tcp_2 as inputName = tcp, instance = 2
      sv := strings.Split(inputId,"_")
      inputName, instance = sv[0], sv[1]
    }
    logp.Info(fmt.Sprintf("input type: %s instance: %s\n", inputName, instance))
    logp.Debug("reader", "instance config: %s", config)

    plugin := newInputInstance(inputName)
    if plugin != nil && config.Enabled {
      err := plugin.Init(config)
      if err != nil {
        logp.Err("Fail to initialize %s plugin as input: %s", inputName, err)
        return err
      } else {
        logp.Info("Initialized %s plugin as input", inputName)
      }
      reader.Input = append(reader.Input, plugin)
    }
  }

  if len(reader.Input) == 0 {
    logp.Info("No inputs are defined. Please define one under the input section.")
    return errors.New("No input are defined. Please define one under the input section.")
  } else {
    logp.Info("%d inputs defined", len(reader.Input))
  }

  return nil
}

func (reader *ReaderType) Run(output chan common.MapStr) error {
  logp.Info("Attempting to start %d inputs", len(reader.Input))

  for _, plugin := range reader.Input {
    err := plugin.Run(output)
    if err != nil {
      logp.Err("Fail to start input plugin %s : %s", plugin.InputType(), err)
        return err
    } else {
      logp.Info("Started input plugin %s", plugin.InputType())
    }
  }
  return nil
}
