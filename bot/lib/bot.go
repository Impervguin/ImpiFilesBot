package lib

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"ImpiFilesBot/internal/service"

	"github.com/mymmrac/telego"
)

const (
	defaultWorkers = 16
)

type Handler interface {
	Handle(ctx context.Context, bot *BotWrapper, update *telego.Update) error
}

type Filter interface {
	Filter(ctx context.Context, bot *BotWrapper, update *telego.Update) bool
}

type FilterHandler interface {
	Handler
	Filter
}

type BotWrapper struct {
	Bot     *telego.Bot
	Logger  Logger
	Service *service.Service

	handlers   []FilterHandler
	handlersMx sync.RWMutex
	workers    int
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewBotWrapper(botToken string, service *service.Service, logger Logger) (*BotWrapper, error) {
	if logger == nil {
		return nil, ErrNilLogger
	}
	if service == nil {
		return nil, ErrNilService
	}

	bot, err := telego.NewBot(botToken, telego.WithLogger(logger))
	if err != nil {
		return nil, fmt.Errorf("can't create bot: %w", err)
	}

	_, err = bot.GetMe(context.Background())
	if err != nil {
		return nil, fmt.Errorf("can't get bot user: %w", err)
	}
	return &BotWrapper{
		Bot:        bot,
		Service:    service,
		ctx:        context.Background(),
		handlers:   make([]FilterHandler, 0),
		handlersMx: sync.RWMutex{},
		workers:    defaultWorkers,
		Logger:     logger,
	}, nil
}

func (b *BotWrapper) RegisterHandler(handler FilterHandler) error {
	b.handlersMx.Lock()
	defer b.handlersMx.Unlock()
	if handler == nil {
		return ErrNilHandler
	}

	b.handlers = append(b.handlers, handler)
	return nil
}

func (b *BotWrapper) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	b.cancel = cancel
	b.ctx = ctx

	updates, err := b.Bot.UpdatesViaLongPolling(ctx, nil)
	if err != nil {
		return fmt.Errorf("can't start bot: %w", err)
	}

	b.Logger.Infof("Starting bot with %d workers", b.workers)
	wg := sync.WaitGroup{}

	wg.Add(b.workers)
	for i := 0; i < b.workers; i++ {
		go b.waitUpdates(ctx, &wg, updates)
	}
	wg.Wait()

	return nil
}

func (b *BotWrapper) Stop() {
	if b.cancel != nil {
		b.Logger.Infof("Stopping bot")
		b.cancel()
		b.cancel = nil
	}
}

func (b *BotWrapper) waitUpdates(ctx context.Context, wg *sync.WaitGroup, updates <-chan telego.Update) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			b.Logger.Infof("Got update: %+v", update)
			err := b.handleUpdate(ctx, &update)
			if err != nil {
				b.Logger.Errorf("Error on handle update: %s", err)
			}
		}
	}
}

func (b *BotWrapper) handleUpdate(ctx context.Context, update *telego.Update) error {
	b.handlersMx.RLock()
	defer b.handlersMx.RUnlock()

	for _, handler := range b.handlers {
		if handler.Filter(ctx, b, update) {
			fmt.Println(reflect.TypeOf(handler))
			err := handler.Handle(ctx, b, update)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
