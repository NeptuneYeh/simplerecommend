package auth

import (
	"github.com/NeptuneYeh/simplerecommend/init/config"
	myToken "github.com/NeptuneYeh/simplerecommend/pkg/token"
	"log"
)

var MyAuth *myToken.Maker

type Module struct {
	TokenMaker myToken.Maker
}

func NewModule() *Module {
	maker, err := myToken.NewPasetoMaker(config.MyConfig.TokenSymmetricKey)
	if err != nil {
		log.Fatalf("Failed to create token maker: %v", err)
	}

	MyAuth = &maker
	authModule := &Module{
		TokenMaker: maker,
	}

	return authModule
}
