package client

import "github.com/appscode/api/dtypes"

type Response interface {
	Reset()
	String() string
	ProtoMessage()
	Descriptor() ([]byte, []int)

	GetStatus() *dtypes.Status
}
