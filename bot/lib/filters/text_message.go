package filters

import (
	"ImpiFilesBot/bot/lib"
	"context"
	"strings"

	"github.com/mymmrac/telego"
)

var _ lib.Filter = &TextMessageFilter{}

type TextMessageFilter struct{}

func NewTextMessageFilter() *TextMessageFilter {
	return &TextMessageFilter{}
}

func (f *TextMessageFilter) Filter(ctx context.Context, bot *lib.BotWrapper, update *telego.Update) bool {
	if update.Message == nil {
		return false
	}
	mess := strings.TrimSpace(update.Message.Text)
	return mess != ""
}
