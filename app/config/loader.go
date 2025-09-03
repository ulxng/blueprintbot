package config

import (
	"fmt"
	"ulxng/blueprintbot/lib/config"
	"ulxng/blueprintbot/lib/messages"
)

type NavigationLoader interface {
	config.Loader[messages.Message]
	messages.Navigator
}

type NavigableLoader struct {
	// dictionary: answer.text => message.code
	replyButtonsCodes map[string]string
	config.Loader[messages.Message]
}

func NewNavigableLoader(source string) (*NavigableLoader, error) {
	l, err := messages.NewLoader(source)
	if err != nil {
		return nil, fmt.Errorf("config.NewNavigableLoader: %w", err)
	}
	loader := &NavigableLoader{
		replyButtonsCodes: make(map[string]string),
		Loader:            l,
	}
	loader.buildTextToNextMap()
	return loader, nil
}

func (l *NavigableLoader) buildTextToNextMap() {
	for _, msg := range l.All() {
		for _, row := range msg.Answers {
			for _, btn := range row {
				if btn.Link != "" {
					l.replyButtonsCodes[btn.Text] = btn.Link
				}
			}
		}
	}
}

func (l *NavigableLoader) GetNextByText(text string) (messages.Message, error) {
	if val, ok := l.replyButtonsCodes[text]; ok {
		return l.GetByKey(val), nil
	} else {
		return messages.Message{}, messages.ErrMessageNotFound
	}
}
