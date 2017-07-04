package util

import (
	"io/ioutil"
	"os"
)

func MountedSecretToMap(path string) (map[string]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var m map[string]string = make(map[string]string)
	for _, f := range files {
		if !f.Mode().IsRegular() {
			continue
		}
		fi, err := os.Open(path + "/" + f.Name())
		if err != nil {
			return nil, err
		}
		cnt, err := ioutil.ReadAll(fi)
		if err != nil {
			return nil, err
		}
		m[f.Name()] = string(cnt)
	}
	return m, nil
}

func MapsToFiles(mp map[string]string, path string) {
	for k, v := range mp {
		ioutil.WriteFile(path+"/"+k, []byte(v), 0777)
	}
}
