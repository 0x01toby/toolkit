package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/taorzhang/toolkit/server"
	"testing"
)

func Test_Gin(t *testing.T) {
	gin, err := server.NewGin("debug")
	assert.NoError(t, err)
	grpc := server.NewGrpc()
	err = server.Run(":9002", server.GinServerMuxOpt(gin), server.GrpcServerMuxOpt(grpc))
	assert.NoError(t, err)
}
