package tcp

import (
  "errors"
  "github.com/johann8384/libbeat/common"
  "github.com/johann8384/libbeat/logp"
  "github.com/turn/turnbeat/inputs"
  "gopkg.in/mcuadros/go-syslog.v2"
)

type SyslogInput struct {
  Port       int /* the port to listen on */
  Type       string /* the type to add to events */
}

func (l *SyslogInput) InputType() string {
  return "SyslogInput"
}

func (l *SyslogInput) InputVersion() string {
  return "0.0.1"
}

func (l *SyslogInput) Init(config inputs.MothershipConfig) error {

  if config.Port == 0 {
    return errors.New("No Input Port specified")
  }
  l.Port = config.Port

  if config.Type == "" {
    return errors.New("No Event Type specified")
  }
  l.Type = config.Type

  logp.Info("[SyslogInput] Using Port %d", l.Port)
  logp.Info("[SyslogInput] Adding Event Type %s", l.Type)

  return nil
}

func (l *SyslogInput) Run(output chan common.MapStr) error {
  logp.Debug("SyslogInput", "Running Syslog Input")
  channel := make(syslog.LogPartsChannel)
  handler := syslog.NewChannelHandler(channel)

  server := syslog.NewServer()
  server.SetFormat(syslog.RFC5424)
  server.SetHandler(handler)
  server.ListenUDP("0.0.0.0:514")
  server.ListenTCP("0.0.0.0:514")

  server.Boot()

  go func(channel syslog.LogPartsChannel) {
    for logParts := range channel {
      logp.Debug("sysloginput", "%v", logParts)
    }
  }(channel)

  server.Wait()
  return nil
}