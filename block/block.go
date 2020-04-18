package block

import (
	"encoding/json"
	"reflect"
	"strings"
)

type Content interface{}

type Block struct {
	Type    string `json:"_type"`
	Content interface{}
}

func New(style string, opts ...BlockOption) Block {
	bc := &BlockContent{
		Style:    style,
		MarkDefs: make([]MarkDef, 0),
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
			tagVals := strings.Split(field.Tag.Get("json"), ",")
			if len(tagVals) > 1 && tagVals[1] == "omitempty" {
				if val.FieldByName(name).IsZero() {
					continue
				}
			}
			m[tagVals[0]] = val.FieldByName(name).Interface()
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
	ListItem string    `json:"listItem,omitempty"`
	Level    int       `json:"level,omitempty"`
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

func (md MarkDef) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"_type": md.Type,
		"_key":  md.Key,
	}

	if data, ok := md.Data.(map[string]interface{}); ok {
		for k, v := range data {
			m[k] = v
		}
	} else if t := reflect.TypeOf(md.Data).Elem(); t.Kind() == reflect.Struct {
		val := reflect.ValueOf(md.Data).Elem()
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			name := t.Field(i).Name
			tagVals := strings.Split(field.Tag.Get("json"), ",")
			if len(tagVals) > 1 && tagVals[1] == "omitempty" {
				if val.FieldByName(name).IsZero() {
					continue
				}
			}
			m[tagVals[0]] = val.FieldByName(name).Interface()
		}
	}

	return json.Marshal(m)
}

type LinkData struct {
	Href string `json:"href"`
}

type SpanContent struct {
	Text  string   `json:"text"`
	Marks []string `json:"marks,omitempty"`
}

func ToPlainText(blocks []Block) string {
	var s strings.Builder
	for _, b := range blocks {
		if bc, ok := b.Content.(*BlockContent); ok {
			if s.Len() > 0 {
				s.WriteString("\n\n")
			}
			for _, span := range bc.Children {
				if sc, ok := span.Content.(*SpanContent); ok {
					s.WriteString(sc.Text)
				}
			}
		}
	}
	return s.String()
}
