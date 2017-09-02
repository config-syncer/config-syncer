package main

import (
	"os"

	logs "github.com/appscode/go/log/golog"
	"github.com/appscode/kubed/pkg/cmds"
	_ "k8s.io/client-go/kubernetes/fake"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	if err := cmds.NewCmdKubed(Version).Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
