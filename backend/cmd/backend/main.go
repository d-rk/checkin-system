package main

import "github.com/d-rk/checkin-system/pkg/server"

//go:generate go tool oapi-codegen --config=../../open-api-conf.yaml ../../open-api-spec.yaml

func main() {
	server.Run()
}
