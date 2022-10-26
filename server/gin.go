package server

import "github.com/gin-gonic/gin"

type GinOpt func(gin *gin.Engine) error

func NewGin(mode string, opts ...GinOpt) (*gin.Engine, error) {
	gin.SetMode(mode)
	g := gin.New()
	for _, opt := range opts {
		if err := opt(g); err != nil {
			return nil, err
		}
	}
	return g, nil
}
