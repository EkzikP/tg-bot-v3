package main

import (
	"context"
	"database/sql"
	"log"
	"sync"

	"github.com/EkzikP/sdk_andromeda_go_v2"
	"github.com/EkzikP/tg-bot-v3/internal/config"
	"github.com/EkzikP/tg-bot-v3/internal/handlers"
	"github.com/EkzikP/tg-bot-v3/internal/services"
	"github.com/EkzikP/tg-bot-v3/internal/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "modernc.org/sqlite"
)

func main() {
	ctx := context.Background()
	// Загрузка конфигурации
	cfg := config.Read()

	// Инициализация БД
	db, err := sql.Open("sqlite", "users.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := storage.New(db)
	if err := store.Initialize(); err != nil {
		log.Fatal(err)
	}

	// Инициализация Andromeda
	andromedaCfg := andromeda.Config{
		ApiKey: cfg.ApiKey,
		Host:   cfg.Host,
	}
	andromedaService := services.New(andromedaCfg)

	// Создание бота
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	// Инициализация мапы для хранения кэша пользователей
	var TgUsers sync.Map

	// Инициализация обработчиков
	handlerMessage := handlers.NewMessageHandler(bot, cfg, store, andromedaService)
	handlerCallback := handlers.NewCallbackHandler(bot, cfg.PhoneEngineer, andromedaService)

	// Настройка обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	// Обработка сообщений
	for update := range updates {
		if update.Message != nil {
			handlerMessage.HandleCommand(ctx, update, &TgUsers)
			return
		}

		if update.CallbackQuery != nil && update.CallbackQuery.Message != nil {
			handlerCallback.HandleCallback(ctx, update, &TgUsers)
			return
		}

	}
}
