package main

import (
  "flag"
  "fmt"
  "io/ioutil"
  "os"
  "runtime"
  "gopkg.in/yaml.v2"
  "github.com/elastic/libbeat/logp"
  "github.com/elastic/packetbeat/config"
)

// You can overwrite these, e.g.: go build -ldflags "-X main.Version 1.0.0-beta3"
var Version = "0.0.1"
var Name = "turnbeat"

func main() {
  // Use our own FlagSet, because some libraries pollute the global one
  var cmdLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
  configfile := cmdLine.String("c", "./" + Name + ".yml", "Configuration file")
  printVersion := cmdLine.Bool("version", false, "Print version and exit")

  // Adds logging specific flags
  logp.CmdLineFlags(cmdLine)

  cmdLine.Parse(os.Args[1:])

  if *printVersion {
    fmt.Printf("Packetbeat version %s (%s)\n", Version, runtime.GOARCH)
    return
  }

  // configuration file
  filecontent, err := ioutil.ReadFile(*configfile)
  if err != nil {
    fmt.Printf("Fail to read %s: %s. Exiting.\n", *configfile, err)
    os.Exit(1)
  }
  if err = yaml.Unmarshal(filecontent, &config.ConfigSingleton); err != nil {
    fmt.Printf("YAML config parsing failed on %s: %s. Exiting.\n", *configfile, err)
    os.Exit(1)
  }

  logp.Init(Name, &config.ConfigSingleton.Logging)

  logp.Info("TurnBeat Started")
}
