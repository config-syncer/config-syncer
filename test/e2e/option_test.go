/*
Copyright The Config Syncer Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package e2e_test

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"gomodules.xyz/logs"
	"k8s.io/client-go/util/homedir"
)

type E2EOptions struct {
	KubeContext string
	KubeConfig  string
}

var options = &E2EOptions{
	KubeConfig: func() string {
		kubecfg := os.Getenv("KUBECONFIG")
		if kubecfg != "" {
			return kubecfg
		}
		return filepath.Join(homedir.HomeDir(), ".kube", "config")
	}(),
}

// xref: https://github.com/onsi/ginkgo/issues/602#issuecomment-559421839
func TestMain(m *testing.M) {
	flag.StringVar(&options.KubeConfig, "kubeconfig", options.KubeConfig, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	flag.StringVar(&options.KubeContext, "kube-context", "", "Name of kube context")
	enableLogging()
	flag.Parse()
	os.Exit(m.Run())
}

func enableLogging() {
	defer func() {
		logs.InitLogs()
		defer logs.FlushLogs()
	}()
	logLevelFlag := flag.Lookup("v")
	if logLevelFlag != nil {
		if len(logLevelFlag.Value.String()) > 0 && logLevelFlag.Value.String() != "0" {
			return
		}
	}
	_ = flag.Set("v", "3")
}
