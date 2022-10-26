package server

import (
	"github.com/gin-gonic/gin"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"net"
)

type MuxOpt func(mux cmux.CMux) error

func Run(address string, opts ...MuxOpt) error {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	mux := cmux.New(listen)
	for idx := range opts {
		go func(idx int) {
			if err = opts[idx](mux); err != nil {
				panic(err)
			}
		}(idx)
	}
	return mux.Serve()
}

func GrpcServerMuxOpt(grpc *grpc.Server) MuxOpt {
	return func(mux cmux.CMux) error {
		return grpc.Serve(mux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc")))
	}
}

func GinServerMuxOpt(gin *gin.Engine) MuxOpt {
	return func(mux cmux.CMux) error {
		return gin.RunListener(mux.Match(cmux.HTTP1Fast()))
	}
}
