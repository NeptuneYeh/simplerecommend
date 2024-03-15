package init

import (
	"github.com/NeptuneYeh/simplerecommend/init/asynq"
	"github.com/NeptuneYeh/simplerecommend/init/auth"
	"github.com/NeptuneYeh/simplerecommend/init/config"
	"github.com/NeptuneYeh/simplerecommend/init/gin"
	"github.com/NeptuneYeh/simplerecommend/init/logger"
	"github.com/NeptuneYeh/simplerecommend/init/redis"
	"github.com/NeptuneYeh/simplerecommend/init/store"
	"os"
	"os/signal"
	"syscall"
)

type MainInitProcess struct {
	ConfigModule *config.Module
	LogModule    *logger.Module
	AuthModule   *auth.Module
	StoreModule  *store.Module
	RedisModule  *redis.Module
	AsynqModule  *asynq.Module
	GinModule    *gin.Module
	OsChannel    chan os.Signal
}

func NewMainInitProcess(configPath string) *MainInitProcess {

	configModule := config.NewModule(configPath)
	logModule := logger.NewModule()
	authModule := auth.NewModule()
	storeModule := store.NewModule()
	redisModule := redis.NewModule()
	asynqModule := asynq.NewModule()
	ginModule := gin.NewModule()
	channel := make(chan os.Signal, 1)

	return &MainInitProcess{
		ConfigModule: configModule,
		LogModule:    logModule,
		AuthModule:   authModule,
		StoreModule:  storeModule,
		RedisModule:  redisModule,
		AsynqModule:  asynqModule,
		GinModule:    ginModule,
		OsChannel:    channel,
	}
}

// Run
func (m *MainInitProcess) Run() {
	go m.AsynqModule.Run(m.AsynqModule.RedisOpt, m.StoreModule.Store)
	m.GinModule.Run(m.ConfigModule.ServerAddress, m.OsChannel)
	// register os signal for graceful shutdown
	signal.Notify(m.OsChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	s := <-m.OsChannel
	//log.Fatalf("Received signal: " + s.String())
	m.LogModule.ZeroLogger.Fatal().Msg("Received signal: " + s.String())
	m.GinModule.Shutdown()
	m.AsynqModule.Shutdown()
}
