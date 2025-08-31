package sender

import (
	"fmt"
	"ulxng/blueprintbot/app/resolver"
	"ulxng/blueprintbot/lib/messages"

	tele "gopkg.in/telebot.v4"
)

type DefaultSender struct {
	resolver resolver.MessageResolver
}

func NewDefaultSender(r resolver.MessageResolver) *DefaultSender {
	return &DefaultSender{resolver: r}
}

func (b *DefaultSender) Send(c tele.Context, messageKey string) error {
	message, markup, err := b.resolver.Get(messageKey)
	if err != nil {
		return fmt.Errorf("get: %w", err)
	}
	return c.Send(message, markup)
}

func (b *DefaultSender) SendRaw(c tele.Context, msg messages.Message) error {
	message, markup, err := b.resolver.Convert(msg)
	if err != nil {
		return fmt.Errorf("convert: %w", err)
	}
	return c.Send(message, markup)
}

func (b *DefaultSender) Edit(c tele.Context, messageKey string) error {
	message, markup, err := b.resolver.Get(messageKey)
	if err != nil {
		return fmt.Errorf("get: %w", err)
	}
	if c.Message().Text == "" {
		// для предотвращения ошибки telegram: Bad Request: there is no text in the message to edit (400)
		// происходит, когда пытаешься сделать edit на сообщениях с кнопками и файлом/картинкой
		return c.Send(message, markup)
	}
	return c.Edit(message, markup)
}
