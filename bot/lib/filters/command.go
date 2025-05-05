package filters

import (
	"ImpiFilesBot/bot/lib"
	"context"
	"fmt"
	"strings"

	"github.com/mymmrac/telego"
)

var _ lib.Filter = &CommandFilter{}

type CommandFilter struct {
	command string
}

func NewCommandFilter(command string) *CommandFilter {
	command = strings.TrimSpace(command)
	if !strings.HasPrefix(command, "/") {
		command = "/" + command
	}
	return &CommandFilter{command: command}
}

func (f *CommandFilter) Filter(ctx context.Context, bot *lib.BotWrapper, update *telego.Update) bool {
	if update.Message == nil {
		return false
	}

	mess := strings.TrimSpace(update.Message.Text)
	fmt.Println(mess, f.command)
	return mess == f.command
}
