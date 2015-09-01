package procfs

import (
        "os"
        "regexp"
        "syscall"
)

func pathExists(pathname string) bool {
        _, err := os.Stat(pathname)
        return err != syscall.ENOENT
}

func isNumeric(s string) bool {
        a, _ := regexp.MatchString("^[0-9]+$", s)
        return a
}

