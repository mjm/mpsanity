package block

type Builder struct {
	bs      []Block
	current *Block
	curSpan *Block
}

func (b *Builder) StartBlock(style string) {
	if b.current != nil {
		b.EndBlock()
	}

	newBlock := New(style)
	b.current = &newBlock
}

func (b *Builder) EndBlock() {
	if b.current == nil {
		return
	}

	b.bs = append(b.bs, *b.current)
	b.current = nil
}

func (b *Builder) StartSpan(marks ...string) {
	if b.curSpan != nil {
		b.EndSpan()
	}
	b.curSpan = &Block{
		Type: "span",
		Content: &SpanContent{
			Marks: marks,
		},
	}
}

func (b *Builder) EndSpan() {
	if b.curSpan == nil {
		return
	}
	sc := b.curSpan.Content.(*SpanContent)
	if sc.Text == "" {
		b.curSpan = nil
		return
	}

	// TODO can there not be a current block?
	if b.current.Type == "block" {
		bc := b.current.Content.(*BlockContent)
		bc.Children = append(bc.Children, *b.curSpan)
		b.curSpan = nil
	} else {
		panic("trying to add a span to a custom block type")
	}
}

func (b *Builder) AppendText(text string) {
	if b.curSpan == nil {
		b.StartSpan()
		b.AppendText(text)
		b.EndSpan()
	} else {
		sc := b.curSpan.Content.(*SpanContent)
		sc.Text += text
	}
}

func (b *Builder) Blocks() []Block {
	b.EndBlock()
	return b.bs
}
