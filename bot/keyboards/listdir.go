package keyboards

import (
	"ImpiFilesBot/internal/domain"

	"github.com/mymmrac/telego"
)

func NewListDirKeyboard(dir *domain.Directory) telego.ReplyMarkup {
	keyboard := make([][]telego.KeyboardButton, 0)
	for _, file := range dir.Files {
		keyboard = append(keyboard, []telego.KeyboardButton{
			{
				Text: file.Name,
			},
		})
	}
	for _, dir := range dir.Dirs {
		keyboard = append(keyboard, []telego.KeyboardButton{
			{
				Text: dir.Name + "/",
			},
		})
	}

	keyboard = append(keyboard, []telego.KeyboardButton{
		{
			Text: "..",
		},
	})

	return &telego.ReplyKeyboardMarkup{
		Keyboard:       keyboard,
		ResizeKeyboard: true,
	}
}
