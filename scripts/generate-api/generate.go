//go:build generateapi
// +build generateapi

package tools

import (
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../../api/spec/config.server.yaml ../../api/spec/spec.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../../api/spec/config.types.yaml ../../api/spec/spec.yaml
