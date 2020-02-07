package handler

import (
	"fmt"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (e *Env) HandleCancel(m *tb.Message) {
	Send(e.Bot, m.Sender,"There is no active command to cancel. I wasn't doing anything anyway. Zzzzz..." )
}

func (e *Env) HandleConvoCancel(c Conversation, m *tb.Message) {
	e.CM.StopConversation(c)

	var msg []string
	msg = append(msg, fmt.Sprintf("The '*%s*' command was cancelled. Anything else I can do for you?", c.Name()))
	msg = append(msg, "")
	msg = append(msg, "Send /help for a list of commands.")
	Send(e.Bot, m.Sender, strings.Join(msg, "\n"))
}