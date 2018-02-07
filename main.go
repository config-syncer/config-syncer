package main

import (
	"os"
	"runtime"

	logs "github.com/appscode/go/log/golog"
	"github.com/appscode/kubed/pkg/cmds"
	_ "k8s.io/client-go/kubernetes/fake"
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
	os.Exit(0)
}
