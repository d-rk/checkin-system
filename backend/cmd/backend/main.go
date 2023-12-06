package main

import "github.com/d-rk/checkin-system/internal/server"

//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=../../open-api-conf.yaml ../../open-api-spec.yaml

func main() {
	server.Run()
}
