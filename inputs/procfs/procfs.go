package procfs

import (
  "strconv"
  "os"
  "io/ioutil"
  "path"
  "time"
  "github.com/johann8384/libbeat/common"
  "github.com/johann8384/libbeat/logp"
  "github.com/turn/turnbeat/inputs"
)

const procfsdir = "/proc"

type ProcfsInput struct {
  Sleep		 int
}

func (l *ProcfsInput) InputType() string {
  return "ProcfsInput"
}

func (l *ProcfsInput) InputVersion() string {
  return "0.0.1"
}

func (l *ProcfsInput) Init(config inputs.MothershipConfig) error {

  if config.Sleep_interval == 0 || config.Sleep_interval < 15 {
    l.Sleep = 15 /* default to 15s */
  } else {
    l.Sleep = config.Sleep_interval
  }

  logp.Info("[ProcfsInput] Initialized, using sleep interval " + strconv.Itoa(l.Sleep))

  return nil
}

type Process struct {
	PID	int
        Cmdline string
        Cwd     string
	Root	string
        Status  string
	// Fds
	// Threads
}

func getProcDetail(PID string) (*Process) {
  pdir := path.Join(procfsdir, PID)

  p := new(Process)
  p.PID, _ = strconv.Atoi(PID)

  cl, _ := ioutil.ReadFile(path.Join(pdir, "cmdline"))
  p.Cwd, _ = os.Readlink(path.Join(pdir, "cwd"))
  p.Root, _ = os.Readlink(path.Join(pdir, "root"))
  status, _ := ioutil.ReadFile(path.Join(pdir, "status"))

  p.Cmdline = byteTransform(cl)
  p.Status = byteTransform(status)

  return p
}

func getCmdline(PID string) string {
  pdir := path.Join(procfsdir, PID)

  cl, _ := ioutil.ReadFile(path.Join(pdir, "cmdline"))
  cmdline := byteTransform(cl)
  return cmdline
}

func scanProcs(output chan common.MapStr) {
  now := func() time.Time {
    t := time.Now()
    return t
  }

  if !pathExists(procfsdir) {
    return
  }
  ds, err := ioutil.ReadDir(procfsdir)
  if err != nil {
    return
  }

  event := common.MapStr{}
  processes := common.MapStr{}
  proc_detail := common.MapStr{}

  // get all numeric entries
  for _, d := range ds {
    n := d.Name()
    if isNumeric(n) {
      processes[n] = getCmdline(n)
      proc_detail[n] = getProcDetail(n)
    }
  }

  text := "process report"
  event["message"] = &text
  event["data"] = processes
  event["data_detail"] = proc_detail
  event["type"] = "report"

  event.EnsureTimestampField(now)
  event.EnsureCountField()
  output <- event
}

func (l *ProcfsInput) periodic(output chan common.MapStr) {
  logp.Debug("[procfsinput]", "Running..")

  scanProcs(output)
}

func (l *ProcfsInput) Run(output chan common.MapStr) error {
  logp.Debug("[procfsinput]", "Starting up Procfs Input")

  go func(output chan common.MapStr) {
    for {
      l.periodic(output)
      time.Sleep(time.Duration(l.Sleep) * time.Second)
    }
  }(output)

  return nil
}
