package block

import (
	"encoding/json"
	"reflect"
)

type Content interface{}

type Block struct {
	Type    string `json:"_type"`
	Content interface{}
}

func New(style string, opts ...BlockOption) Block {
	bc := &BlockContent{
		Style: style,
	}

	for _, o := range opts {
		o.Apply(bc)
	}

	return Block{
		Type:    "block",
		Content: bc,
	}
}

func (b Block) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"_type": b.Type,
	}

	if content, ok := b.Content.(map[string]interface{}); ok {
		for k, v := range content {
			m[k] = v
		}
	} else if t := reflect.TypeOf(b.Content).Elem(); t.Kind() == reflect.Struct {
		val := reflect.ValueOf(b.Content).Elem()
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			name := t.Field(i).Name
			m[field.Tag.Get("json")] = val.FieldByName(name).Interface()
		}
	}

	return json.Marshal(m)
}

func (b *Block) UnmarshalJSON(data []byte) error {
	var typeVal struct {
		Type string `json:"_type"`
	}
	if err := json.Unmarshal(data, &typeVal); err != nil {
		return err
	}

	b.Type = typeVal.Type

	switch b.Type {
	case "block":
		var bc BlockContent
		if err := json.Unmarshal(data, &bc); err != nil {
			return err
		}
		b.Content = &bc
		return nil
	case "span":
		var sc SpanContent
		if err := json.Unmarshal(data, &sc); err != nil {
			return err
		}
		b.Content = &sc
		return nil
	default:
		m := map[string]interface{}{}
		if err := json.Unmarshal(data, &m); err != nil {
			return err
		}
		delete(m, "_type")
		b.Content = m
		return nil
	}
}

type BlockContent struct {
	Style    string    `json:"style"`
	Children []Block   `json:"children"`
	MarkDefs []MarkDef `json:"markDefs"`
}

type BlockOption interface {
	Apply(bc *BlockContent)
}

type blockOptionFn func(bc *BlockContent)

func (fn blockOptionFn) Apply(bc *BlockContent) {
	fn(bc)
}

func Text(s string, marks ...string) BlockOption {
	if marks == nil {
		marks = []string{}
	}

	return blockOptionFn(func(bc *BlockContent) {
		bc.Children = append(bc.Children, Block{
			Type: "span",
			Content: &SpanContent{
				Text:  s,
				Marks: marks,
			},
		})
	})
}

type Data interface{}

type MarkDef struct {
	Type string `json:"_type"`
	Key  string `json:"_key"`
	Data
}

type SpanContent struct {
	Text  string   `json:"text"`
	Marks []string `json:"marks"`
}