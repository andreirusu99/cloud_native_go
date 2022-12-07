package api

import (
	"cloud_native_go/core"
	"fmt"
	"os"
)

type FrontEnd interface {
	Start(service *core.KVStore) error
}

func NewFrontEnd() (FrontEnd, error) {
	frontEndType := os.Getenv("FRONTEND_TYPE")

	switch frontEndType {
	case "rest":
		return NewRestFrontEnd()

	case "grpc":
		return nil, fmt.Errorf("gRPC frontend not implemented")

	case "":
		return nil, fmt.Errorf("no frontend type defined")

	default:
		return nil, fmt.Errorf("frontend type %s not defined", frontEndType)

	}
}
