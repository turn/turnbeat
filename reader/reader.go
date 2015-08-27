package reader

import (
  "encoding/json"
  "errors"

  "github.com/johann8384/libbeat/common"
  "github.com/johann8384/libbeat/logp"
  "github.com/turn/turnbeat/inputs"
  "github.com/turn/turnbeat/inputs/tcp"
  "github.com/turn/turnbeat/inputs/syslog"
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

var EnabledInputPlugins map[inputs.InputPlugin]inputs.InputInterface = map[inputs.InputPlugin]inputs.InputInterface{
  inputs.TcpInput:    new(tcp.TcpInput),
  inputs.SyslogInput:    new(syslog.SyslogInput),
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

  for inputId, plugin := range EnabledInputPlugins {
    inputName := inputId.String()
    input, exists := inputs[inputName]
    if exists && input.Enabled {
      err := plugin.Init(input)
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
  }

  return nil
}

func (reader *ReaderType) Run(output chan common.MapStr) error {
  for _, plugin := range reader.Input {
    err := plugin.Run(output)
    if err != nil {
      logp.Err("Fail to start input plugin %s : %s", plugin.InputType, err)
        return err
    } else {
      logp.Info("Started input plugin %s", plugin.InputType)
    }
  }
  return nil
}
