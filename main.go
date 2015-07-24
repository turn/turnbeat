package main

import (
  "flag"
  "fmt"
  "io/ioutil"
  "os"
  "runtime"
  "gopkg.in/yaml.v2"
  "github.com/johann8384/libbeat/common"
  "github.com/johann8384/libbeat/common/droppriv"
  "github.com/johann8384/libbeat/publisher"
  "github.com/johann8384/libbeat/logp"
  "github.com/johann8384/libbeat/filters"
  "github.com/johann8384/libbeat/filters/nop"
  "github.com/johann8384/libbeat/filters/opentsdb"
  "github.com/turn/turnbeat/config"
  "github.com/turn/turnbeat/reader"
)

// You can overwrite these, e.g.: go build -ldflags "-X main.Version 1.0.0-beta3"
var Version = "0.0.1"
var Name = "turnbeat"

var EnabledFilterPlugins map[filters.Filter]filters.FilterPlugin = map[filters.Filter]filters.FilterPlugin{
  filters.NopFilter: new(nop.Nop),
  filters.OpenTSDBFilter: new (opentsdb.OpenTSDB),
}

func main() {
  // Use our own FlagSet, because some libraries pollute the global one
  var cmdLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
  configfile := cmdLine.String("c", "./" + Name + ".yml", "Configuration file")
  publishDisabled := cmdLine.Bool("N", false, "Disable actual publishing for testing")
  printVersion := cmdLine.Bool("version", false, "Print version and exit")

  // Adds logging specific flags
  logp.CmdLineFlags(cmdLine)

  cmdLine.Parse(os.Args[1:])

  if *printVersion {
    fmt.Printf("Turnbeat version %s (%s)\n", Version, runtime.GOARCH)
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

  logp.Info("Initializing output plugins")
  if err = publisher.Publisher.Init(*publishDisabled, config.ConfigSingleton.Output,
    config.ConfigSingleton.Shipper); err != nil {

    logp.Critical(err.Error())
    os.Exit(1)
  }

  logp.Info("Initializing filter plugins")
  for filter, plugin := range EnabledFilterPlugins {
    logp.Debug("main", "Registering Plugin: %s", filter)
    filters.Filters.Register(filter, plugin)
  }
  logp.Debug("main", "Filter Config: %s", config.ConfigSingleton.Filter)
  filters_plugins, err :=
    LoadConfiguredFilters(config.ConfigSingleton.Filter)
  if err != nil {
    logp.Critical("Error loading filter plugins: %v", err)
    os.Exit(1)
  }
  logp.Debug("main", "Filter plugins order: %v", filters_plugins)
  var afterInputsQueue chan common.MapStr
  if len(filters_plugins) > 0 {
    runner := NewFilterRunner(publisher.Publisher.Queue, filters_plugins)
    go func() {
      err := runner.Run()
      if err != nil {
        logp.Critical("Filters runner failed: %v", err)
      }
    }()
    afterInputsQueue = runner.FiltersQueue
  } else {
    // short-circuit the runner
    afterInputsQueue = publisher.Publisher.Queue
  }

  logp.Info("Initializing input plugins")
  if err = reader.Reader.Init(config.ConfigSingleton.Input); err != nil {
    logp.Critical(err.Error())
    os.Exit(1)
  }

  if err = droppriv.DropPrivileges(config.ConfigSingleton.RunOptions); err != nil {
    logp.Critical(err.Error())
    os.Exit(1)
  }

  logp.Info("Starting input plugins")
  if err = reader.Reader.Run(publisher.Publisher.Queue); err != nil {
    logp.Critical(err.Error())
    os.Exit(1)
  }

  logp.Info("Turnbeat Started")
  for {
    event := <-afterInputsQueue
    logp.Info("Event: %v", event)
  }
}
