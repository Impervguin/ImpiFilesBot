package main

import (
	"ImpiFilesBot/bot/handlers"
	"ImpiFilesBot/bot/lib"
	"ImpiFilesBot/internal/auth"
	"ImpiFilesBot/internal/domain"
	"ImpiFilesBot/internal/infra/osfiles"
	"ImpiFilesBot/internal/infra/usermap"
	"ImpiFilesBot/internal/service"
	"context"
	"fmt"
	"os"
	"os/signal"

	"go.uber.org/zap"

	"github.com/spf13/viper"
)

type OsConfig struct {
	Root string `yaml:"root"`
}

type BotConfig struct {
	Token string `yaml:"token"`
}

type OwnerConfig struct {
	TelegramID int64  `yaml:"chat_id" mapstructure:"chat_id"`
	Username   string `yaml:"username" mapstructure:"username"`
	Cwd        string `yaml:"cwd" mapstructure:"cwd"`
}

type Config struct {
	Os     OsConfig      `yaml:"os"`
	Bot    BotConfig     `yaml:"bot"`
	Owners []OwnerConfig `yaml:"owners"`
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	config := Config{}
	viper.SetConfigName("bot")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println("Raw config:", viper.AllSettings())
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	fmt.Println("Config loaded: ", config)

	userStorage := usermap.NewTelegramUserMap()
	for _, owner := range config.Owners {
		user := domain.NewOwner(owner.Username, owner.TelegramID, owner.Cwd)
		userStorage.Set(owner.TelegramID, user)
		fmt.Println("Owner added:", user.TelegramID(), user.Username())
	}

	fileStorage := osfiles.NewOsFileSystem(config.Os.Root)

	auth := auth.NewAuthService(userStorage)

	serv := service.NewService(*auth, fileStorage)

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	sugar := logger.Sugar()

	bot, err := lib.NewBotWrapper(config.Bot.Token, serv, sugar)
	if err != nil {
		panic(err)
	}

	// Register handlers
	bot.RegisterHandler(handlers.NewRootHandler())
	bot.RegisterHandler(handlers.NewChDirHandler())
	bot.RegisterHandler(handlers.NewFileUploadHandler())
	bot.RegisterHandler(handlers.NewFileDownloadHandler())

	bot.Start(ctx)

	<-ctx.Done()
}
