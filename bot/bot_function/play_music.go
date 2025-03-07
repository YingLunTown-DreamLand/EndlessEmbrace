package BotFunction

import (
	"fmt"
	"phoenixbuilder/fastbuilder/string_reader"
	"strings"
)

func ProcessPlayMusic(commandLine string) (
	shouldSendMessage bool,
	message string,
) {
	var commandName string = ""
	reader := string_reader.NewStringReader(&commandLine)

	reader.JumpSpace()
	if reader.Next(true) != "/" {
		return false, ""
	}

	reader.JumpSpace()
	for {
		token := reader.Next(true)
		if token == "" || token == " " || token == "\n" || token == "\r" {
			break
		}
		commandName += token
	}
	commandName = strings.ToLower(commandName)

	if commandName != "play_music" {
		return false, ""
	}

	reader.JumpSpace()
	musicLocation := reader.Sentence(
		max(
			0,
			len(reader.String())-reader.Pointer(),
		),
	)

	return true, fmt.Sprintf(`[CQ:record,file=%s]`, musicLocation)
}
