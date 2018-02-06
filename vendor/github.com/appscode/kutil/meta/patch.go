package meta

import (
	"encoding/json"

	apps "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/jsonmergepatch"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
)

func CreateStrategicPatch(cur runtime.Object, transform func(runtime.Object) runtime.Object) ([]byte, error) {
	curJson, err := json.Marshal(cur)
	if err != nil {
		return nil, err
	}

	modJson, err := json.Marshal(transform(cur.DeepCopyObject()))
	if err != nil {
		return nil, err
	}

	return strategicpatch.CreateTwoWayMergePatch(curJson, modJson, apps.DaemonSet{})
}

func CreateJSONMergePatch(cur runtime.Object, transform func(runtime.Object) runtime.Object) ([]byte, error) {
	curJson, err := json.Marshal(cur)
	if err != nil {
		return nil, err
	}

	modJson, err := json.Marshal(transform(cur.DeepCopyObject()))
	if err != nil {
		return nil, err
	}

	return jsonmergepatch.CreateThreeWayJSONMergePatch(curJson, modJson, curJson)
}
