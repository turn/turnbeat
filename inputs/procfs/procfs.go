package procfs

import (
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

  l.Sleep = 10 /* hard code for now */

  logp.Info("[ProcfsInput] Initialized, using sleep interval %s", l.Sleep)

  return nil
}

func scanProc(PID string) (common.MapStr) {
  now := func() time.Time {
    t := time.Now()
    return t
  }

  pdir := path.Join(procfsdir, PID)

  cl, _ := ioutil.ReadFile(path.Join(pdir, "cmdline"))
  cmdline := string(cl[:])

  cwd, _ := os.Readlink(path.Join(pdir, "cwd"))

  event := common.MapStr{}
  event["message"] = cmdline
  event["cwd"] = cwd
  event["type"] = "process"

  event.EnsureTimestampField(now)
  event.EnsureCountField()
  return event
}

func getProcInfo(PID string) (common.MapStr) {
   pdir := path.Join(procfsdir, PID)

  cl, _ := ioutil.ReadFile(path.Join(pdir, "cmdline"))
  cmdline := string(cl[:])

  cwd, _ := os.Readlink(path.Join(pdir, "cwd"))

  retval := common.MapStr{}
  retval["cmdline"] = cmdline
  retval["cwd"] = cwd

  return retval
}

func scanProcs(output chan common.MapStr) {
//  event := scanProc("self")
//  output <- event
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

  // get all numeric entries
  for _, d := range ds {
    n := d.Name()
    if isNumeric(n) {
      processes[n] = getProcInfo(n)
    }
  }

  text := "process report"
  event["message"] = &text
  event["data"] = processes
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
