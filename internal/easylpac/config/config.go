package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const AIDDefault = "A0000005591010FFFFFFFF8900000100"
const AID5BER = "A0000005591010FFFFFFFF8900050500"
const AIDESIMME = "A0000005591010000000008900000300"
const AIDXESIM = "A0000005591010FFFFFFFF8900000177"

type Config struct {
	LpacDir     string
	LpacAID     string
	EXEName     string
	DriverIFID  string
	DebugHTTP   bool
	DebugAPDU   bool
	LogDir      string
	LogFilename string
	LogFile     *os.File
	AutoMode    bool
}

var Instance Config

func Load() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return err
	}
	exeDir := filepath.Dir(exePath)
	Instance.LpacDir = exeDir

	switch platform := runtime.GOOS; platform {
	case "windows":
		Instance.EXEName = "lpac.exe"
		Instance.LogDir = filepath.Join(exeDir, "log")
	case "linux":
		Instance.EXEName = "lpac"
		Instance.LogDir = filepath.Join("/tmp", "EasyLPAC-log")
		_, err = os.Stat(filepath.Join(Instance.LpacDir, Instance.EXEName))
		if err != nil {
			Instance.LpacDir = "/usr/bin"
		}
	default:
		Instance.EXEName = "lpac"
		Instance.LogDir = filepath.Join("/tmp", "EasyLPAC-log")
	}
	Instance.AutoMode = true
	Instance.LpacAID = AIDDefault

	Instance.LogFilename = fmt.Sprintf("lpac-%s.txt", time.Now().Format("20060102-150405"))
	return nil
}
