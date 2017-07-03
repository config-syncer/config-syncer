package common

import (
	"errors"
	"fmt"
)

func EnsureRequiredKeys(mp map[string]string, keys []string) error {
	for _, k := range keys {
		if _, found := mp[k]; !found {
			return errors.New(fmt.Sprintf("%v not found", k))
		}
	}
	return nil
}
