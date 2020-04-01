package block

type Builder struct {
	bs            []Block
	current       *Block
	curSpan       *Block
	listItemStack []string
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

	if bc, ok := b.current.Content.(*BlockContent); ok && len(bc.Children) == 0 {
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

func (b *Builder) StartList(listItem string) {
	b.listItemStack = append(b.listItemStack, listItem)
}

func (b *Builder) EndList() {
	b.listItemStack = b.listItemStack[:len(b.listItemStack)-1]
}

func (b *Builder) StartListItem() {
	if len(b.listItemStack) == 0 {
		return
	}

	b.StartBlock("normal")
	bc := b.current.Content.(*BlockContent)
	bc.ListItem = b.listItemStack[len(b.listItemStack)-1]
	bc.Level = len(b.listItemStack)
}

func (b *Builder) EndListItem() {
	b.EndBlock()
}

func (b *Builder) Blocks() []Block {
	b.EndBlock()
	return b.bs
}
