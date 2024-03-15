package gin

import (
	"context"
	"github.com/NeptuneYeh/simplerecommend/init/auth"
	"github.com/NeptuneYeh/simplerecommend/internal/app/controllers"
	"github.com/NeptuneYeh/simplerecommend/internal/app/middlewares"
	myValidator "github.com/NeptuneYeh/simplerecommend/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"os"
	"time"
)

type Module struct {
	Router *gin.Engine
	Server *http.Server
}

func NewModule() *Module {

	r := gin.Default()
	ginModule := &Module{
		Router: r,
	}
	gin.ForceConsoleColor()
	ginModule.setupRoute()

	return ginModule
}

// setup route
func (module *Module) setupRoute() {
	// binding validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("passwd", myValidator.ValidPassword)
		if err != nil {
			return
		}
	}
	// init controller
	accountController := controllers.NewAccountController()
	recommendProductController := controllers.NewRecommendProductController()
	// add routes to router
	module.Router.POST("/accounts", accountController.CreateAccount)
	module.Router.POST("/account/login", accountController.LoginAccount)
	module.Router.GET("/account/verify/email", accountController.VerifyEmail)
	// create route group
	authRoutes := module.Router.Group("/").Use(middlewares.AuthMiddleware(*auth.MyAuth))
	authRoutes.GET("/recommendation", recommendProductController.ListProducts)
}

// Run gin
func (module *Module) Run(address string, osChannel chan os.Signal) {
	module.Server = &http.Server{
		Addr:    address,
		Handler: module.Router,
	}

	go func() {
		if err := module.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()
}

func (module *Module) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := module.Server.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to run Gin shutdown: %v", err)
	}
	return nil
}
