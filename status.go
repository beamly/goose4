package goose4

import (
	"encoding/json"
	"runtime"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
)

// Status embeds Config and System to give a concise system status
type Status struct {
	Config
	System
}

// Marshal returns a status doc based on passed in config and up-to-date
// system details
func (s Status) Marshal(boot time.Time) ([]byte, error) {
	status := Status{
		Config: s.Config,
		System: NewSystem(boot),
	}

	return json.Marshal(status)
}

// System contains system specific data for status responses
type System struct {
	MachineName string `json:"machine_name"`
	OSArch      string `json:"os_arch"`
	OSLoad      string `json:"os_avgload"`
	OSName      string `json:"os_name"`
	OSProcs     string `json:"os_numprocessors"`
	OSVersion   string `json:"os_version"`
	UpDuration  string `json:"up_duration"`
	UpSince     string `json:"up_since"`
}

// NewSystem will generate a goose4.System and fill it with information
// taken from the system on which it is instantiated
func NewSystem(boot time.Time) System {
	h, _ := host.Info()
	l, _ := load.Avg()
	c, _ := cpu.Counts(true)

	return System{
		h.Hostname,
		runtime.GOARCH,
		strconv.Itoa(int(l.Load1)),
		h.OS,
		strconv.Itoa(c),
		h.PlatformVersion,
		time.Since(boot).String(),
		boot.String(),
	}
}
