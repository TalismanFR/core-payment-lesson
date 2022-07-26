//go:build tools

package auth

import (
	_ "github.com/cweill/gotests"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang/mock/mockgen"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
