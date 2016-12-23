package errorhandlers

import (
	"github.com/appscode/errors"
)

func init() {
	errors.Handlers.Add(loggingHandler{})
}
