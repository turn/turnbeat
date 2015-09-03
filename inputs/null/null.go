package null

/* an empty input, useful as a starting point for a new input */

import (
  "time"
  "github.com/johann8384/libbeat/common"
  "github.com/johann8384/libbeat/logp"
  "github.com/turn/turnbeat/inputs"
)

type NullInput struct {
  Type	string
}

func (l *NullInput) InputType() string {
  return "NullInput"
}

func (l *NullInput) InputVersion() string {
  return "0.0.1"
}

func (l *NullInput) Init(config inputs.MothershipConfig) error {
  l.Type = "null"
  logp.Info("[NullInput] Initialized")
  return nil
}

func (l *NullInput) Run(output chan common.MapStr) error {
  logp.Debug("[NullInput]", "Running Null Input")

  // dispatch thread here
  go func(output chan common.MapStr) {
    l.doStuff(output)
  }(output)

  return nil
}

func (l *NullInput) doStuff(output chan common.MapStr) {
  now := func() time.Time {
    t := time.Now()
    return t
  }

  // construct event and write it to channel
  event := common.MapStr{}

  text := "null event"
  event["message"] = &text
  event["type"] = l.Type

  event.EnsureTimestampField(now)
  event.EnsureCountField()

  output <- event


}
