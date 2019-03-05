package main

import (
	"os"
	"runtime"

	"github.com/appscode/kubed/pkg/cmds"
	_ "k8s.io/client-go/kubernetes/fake"
	"kmodules.xyz/client-go/logs"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	if err := cmds.NewCmdKubed(Version).Execute(); err != nil {
		os.Exit(1)
	}
}
