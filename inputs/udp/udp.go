package udp

import (
	"errors"
	"github.com/blacklightops/libbeat/common"
	"github.com/blacklightops/libbeat/logp"
	"github.com/blacklightops/turnbeat/inputs"
	"net"
	"time"
)

type UdpInput struct {
	Config inputs.MothershipConfig
	Port   int    /* the port to listen on */
	Type   string /* the type to add to events */
}

func (l *UdpInput) InputType() string {
	return "UdpInput"
}

func (l *UdpInput) InputVersion() string {
	return "0.0.1"
}

func (l *UdpInput) Init(config inputs.MothershipConfig) error {

	l.Config = config
	if config.Port == 0 {
		return errors.New("No Input Port specified")
	}
	l.Port = config.Port

	if config.Type == "" {
		return errors.New("No Event Type specified")
	}
	l.Type = config.Type

	logp.Debug("udpinput", "Using Port %d", l.Port)
	logp.Debug("udpinput", "Adding Event Type %s", l.Type)

	return nil
}

func (l *UdpInput) GetConfig() inputs.MothershipConfig {
	return l.Config
}

func (l *UdpInput) Run(output chan common.MapStr) error {
	logp.Info("[UdpInput] Running UDP Input")
  addr := net.UDPAddr{
    Port: l.Port,
    IP:   net.ParseIP("0.0.0.0"),
  }
  server, err := net.ListenUDP("udp", &addr)
  server.SetReadBuffer(1048576)

	if err != nil {
		logp.Err("couldn't start listening: " + err.Error())
		return nil
	}

	logp.Info("[UdpInput] Listening on port %d", l.Port)

  i := 0
  for {
    i++
    buf := make([]byte, 4096)
    rlen, addr, err := server.ReadFromUDP(buf)
    if err != nil {
  		logp.Err("couldn't read from UDP: " + err.Error())
    }
    go l.handlePacket(buf, rlen, i, addr, output)
  }
	return nil
}

func(l *UdpInput) handlePacket(buffer []byte, rlen int, count int, source *net.UDPAddr, output chan common.MapStr) {
	now := func() time.Time {
		t := time.Now()
		return t
	}

  text := string(buffer[0:rlen])

  logp.Debug("udpinputlines", "New Line: %s", &text)

	event := common.MapStr{}
	event["source"] = &source
	event["offset"] = rlen
	event["line"] = count
	event["message"] = text
	event["type"] = l.Type

	event.EnsureTimestampField(now)
	event.EnsureCountField()

	logp.Debug("udpinput", "InputEvent: %v", event)
	output <- event // ship the new event downstream
}
