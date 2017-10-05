package framework

import (
	"log"
	"strings"
	"flag"

	"path/filepath"
	"k8s.io/client-go/util/homedir"
	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/go/flags"
)

type E2EConfig struct {
	Master            string
	KubeConfig        string
	CloudProviderName string
	HAProxyImageName  string
	TestNamespace     string
	// IngressClass      string
	InCluster         bool
	Cleanup           bool
	DaemonHostName    string
	LBPersistIP       string
	// RBACEnabled       bool
	// TestCertificate   bool
}

func init() {
	enableLogging()
}


var testConfigs E2EConfig

func enableLogging() {
	flag.Set("logtostderr", "true")
	logLevelFlag := flag.Lookup("v")
	if logLevelFlag != nil {
		if len(logLevelFlag.Value.String()) > 0 && logLevelFlag.Value.String() != "0" {
			return
		}
	}
	flags.SetLogLevel(2)
}

func (c *E2EConfig) validate()  {
	if len(c.KubeConfig) == 0 {
		c.KubeConfig = filepath.Join(homedir.HomeDir(), ".kube/config")
	}

	if len(c.TestNamespace) == 0 {
		c.TestNamespace = rand.WithUniqSuffix("test-kubed")
	}

	if !strings.HasPrefix(c.TestNamespace, "test-") {
		log.Fatal("Namespace is not a Test namespace")
	}
}
