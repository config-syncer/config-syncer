package main

import (
	"os"

	"github.com/appscode/kubed/pkg/cmds"
	logs "github.com/appscode/log/golog"
	_ "k8s.io/client-go/kubernetes/fake"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	if err := cmds.NewCmdKubed().Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
