package filters

import (
	"ImpiFilesBot/bot/lib"
	"context"

	"github.com/mymmrac/telego"
)

type FileFilter struct {
}

func NewFileFilter() *FileFilter {
	return &FileFilter{}
}

func (f *FileFilter) Filter(ctx context.Context, bot *lib.BotWrapper, update *telego.Update) bool {
	if update.Message == nil {
		return false
	}
	return update.Message.Document != nil
}
