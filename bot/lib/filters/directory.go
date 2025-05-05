package filters

import (
	"ImpiFilesBot/bot/lib"
	"context"
	"strings"

	"github.com/mymmrac/telego"
)

var _ lib.Filter = &DirectoryFilter{}

type DirectoryFilter struct{}

func NewDirectoryFilter() *DirectoryFilter {
	return &DirectoryFilter{}
}

func (f *DirectoryFilter) Filter(ctx context.Context, bot *lib.BotWrapper, update *telego.Update) bool {
	if update.Message == nil {
		return false
	}
	mess := strings.TrimSpace(update.Message.Text)
	return strings.HasSuffix(mess, "/") || mess == ".." || mess == "."
}
