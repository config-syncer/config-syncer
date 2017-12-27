package util

import (
	"strconv"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	MaxSyncInterval = 5 * time.Minute
)

func IsRecent(t metav1.Time) bool {
	return time.Now().Sub(t.Time) < MaxSyncInterval
}

func GetBool(m map[string]string, key string) (bool, error) {
	if m == nil {
		return false, nil
	}
	v, ok := m[key]
	if !ok || v == "" {
		return false, nil
	}
	return strconv.ParseBool(v)
}
