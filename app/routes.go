package main

import (
	tele "gopkg.in/telebot.v4"
)

func (a *App) registerRoutes() {
	a.bot.Handle("/start", func(c tele.Context) error {
		return a.sender.Send(c, "root")
	})

	//добавляем динамические эндпоинты для reply buttons (приходят как текст)
	for _, message := range a.loader.All() {
		for _, answer := range message.Answers {
			replyButton := answer
			a.bot.Handle(replyButton.Text, func(c tele.Context) error {
				return a.sender.Route(c, c.Text())
			})
		}
	}
	a.bot.Handle(tele.OnCallback, func(c tele.Context) error {
		if c.Callback() == nil {
			return nil
		}
		return a.sender.Send(c, c.Callback().Data)
	})
}

func (a *App) registerFlows() {

}
