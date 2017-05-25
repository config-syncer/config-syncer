package env

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Environment string

const (
	Dev    Environment = "dev"
	QA     Environment = "qa"
	Prod   Environment = "prod"
	Onebox Environment = "onebox"
)

const (
	Key           = "APPSCODE_ENV"
	ProdApiServer = "https://api.appscode.com:3443"
	QAApiServer   = "https://api.appscode.info:3443"
)

func (e Environment) IsPublic() bool {
	switch e {
	case Prod, Onebox:
		return true
	default:
		return false
	}
}

func (e Environment) IsHosted() bool {
	switch e {
	case Dev, QA, Prod:
		return true
	default:
		return false
	}
}

func (e Environment) DebugEnabled() bool {
	switch e {
	case Dev, QA:
		return true
	default:
		return false
	}
}

func (e Environment) DevMode() bool {
	return e == Dev
}

func (e Environment) APIServer() (string, error) {
	switch e {
	case Prod:
		return ProdApiServer, nil
	case QA:
		return QAApiServer, nil
	default:
		return "", fmt.Errorf("apiserver unknown for env %v", e)
	}
}

func (e Environment) String() string {
	return string(e)
}

func (e *Environment) MarshalJSON() ([]byte, error) {
	return []byte(`"` + *e + `"`), nil
}

func (e *Environment) UnmarshalJSON(data []byte) error {
	if e == nil {
		return errors.New("jsontypes.ArrayOrInt: UnmarshalJSON on nil pointer")
	}
	*e = FromString(string(data[1 : len(data)-1]))
	return nil
}

func FromHost() Environment {
	return FromString(strings.ToLower(strings.TrimSpace(os.Getenv(Key))))
}

func FromString(e string) Environment {
	switch e {
	case "prod":
		return Prod
	case "onebox":
		return Onebox
	case "qa":
		return QA
	case "dev":
		return Dev
	default:
		if InCluster() {
			return Prod
		} else {
			return Dev
		}
	}
}

// Possible returns true if loading an inside-kubernetes-cluster is possible.
// ref: https://goo.gl/mrlLyr
func InCluster() bool {
	fi, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount/token")
	return os.Getenv("KUBERNETES_SERVICE_HOST") != "" &&
		os.Getenv("KUBERNETES_SERVICE_PORT") != "" &&
		err == nil && !fi.IsDir()
}
