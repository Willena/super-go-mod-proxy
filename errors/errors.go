package errors

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
)

var logger, _ = zap.NewDevelopment()

func GenerateError(err error) []byte {
	var b []byte
	if err != nil {
		b, _ = json.Marshal(map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		})

	} else {
		b, _ = json.Marshal(map[string]interface{}{
			"status": "error",
			"error":  fmt.Errorf("Unknown error").Error(),
		})
	}

	logger.Error(err.Error())

	return b
}
