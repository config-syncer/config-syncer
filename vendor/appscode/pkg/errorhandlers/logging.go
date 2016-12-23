package errorhandlers

import (
	"github.com/appscode/errors"
)

type loggingHandler struct{}

func (loggingHandler) Handle(e error) {
	errors.Log(e)
}
