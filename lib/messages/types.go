package messages

import "gopkg.in/yaml.v3"

type Message struct {
	Text    string       `yaml:"text"`
	Buttons Rows[Button] `yaml:"buttons"`
	Answers Rows[Answer] `yaml:"answers"`
	Image   string       `yaml:"image"`
	File    *File        `yaml:"file"`
}

type Button struct {
	Text string `yaml:"text"`
	Code string `yaml:"code"`
	Link string `yaml:"link"`
}

type File struct {
	Path string `yaml:"path"`
	Name string `yaml:"name"`
}

type Answer struct {
	Text    string `yaml:"text"`
	Contact bool   `yaml:"request_contact"`
	Link    string `yaml:"link"`
}

type Rows[T any] [][]T

func (r *Rows[T]) UnmarshalYAML(value *yaml.Node) error {
	// Разрешаем пустое значение
	if value == nil {
		*r = nil
		return nil
	}
	// Ждём последовательность (answers: [...])
	if value.Kind != yaml.SequenceNode {
		*r = nil
		return nil
	}

	var out [][]T
	for _, item := range value.Content {
		switch item.Kind {
		case yaml.SequenceNode:
			var row []T
			if err := item.Decode(&row); err != nil {
				return err
			}
			out = append(out, row)
		default:
			var cell T
			if err := item.Decode(&cell); err != nil {
				return err
			}
			out = append(out, []T{cell})
		}
	}

	*r = out
	return nil
}

func (r Rows[T]) MarshalYAML() (any, error) {
	return [][]T(r), nil
}
