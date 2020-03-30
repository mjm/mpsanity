package patch

type Description struct {
	ID             string                 `json:"id,omitempty"`
	Query          string                 `json:"query,omitempty"`
	Set            map[string]interface{} `json:"set,omitempty"`
	SetIfMissing   map[string]interface{} `json:"setIfMissing,omitempty"`
	Unset          []string               `json:"unset,omitempty"`
	Insert         *insertion             `json:"insert,omitempty"`
	Inc            map[string]interface{} `json:"inc,omitempty"`
	Dec            map[string]interface{} `json:"dec,omitempty"`
	DiffMatchPatch map[string]string      `json:"diffMatchPatch,omitempty"`
}

type Patch interface {
	Apply(p *Description)
}

type patchFn func(p *Description)

func (fn patchFn) Apply(p *Description) {
	fn(p)
}

func Set(key string, val interface{}) Patch {
	return patchFn(func(p *Description) {
		if p.Set == nil {
			p.Set = make(map[string]interface{})
		}
		p.Set[key] = val
	})
}

func SetIfMissing(key string, val interface{}) Patch {
	return patchFn(func(p *Description) {
		if p.SetIfMissing == nil {
			p.SetIfMissing = make(map[string]interface{})
		}
		p.SetIfMissing[key] = val
	})
}

func Unset(keys ...string) Patch {
	return patchFn(func(p *Description) {
		p.Unset = append(p.Unset, keys...)
	})
}

type insertion struct {
	Before  string        `json:"before,omitempty"`
	After   string        `json:"after,omitempty"`
	Replace string        `json:"replace,omitempty"`
	Items   []interface{} `json:"items"`
}

func InsertBefore(match string, items ...interface{}) Patch {
	return patchFn(func(p *Description) {
		p.Insert = &insertion{
			Before: match,
			Items:  items,
		}
	})
}

func InsertAfter(match string, items ...interface{}) Patch {
	return patchFn(func(p *Description) {
		p.Insert = &insertion{
			After: match,
			Items: items,
		}
	})
}

func Replace(match string, items ...interface{}) Patch {
	return patchFn(func(p *Description) {
		p.Insert = &insertion{
			Replace: match,
			Items:   items,
		}
	})
}

func Inc(key string, val interface{}) Patch {
	return patchFn(func(p *Description) {
		if p.Inc == nil {
			p.Inc = make(map[string]interface{})
		}
		p.Inc[key] = val
	})
}

func Dec(key string, val interface{}) Patch {
	return patchFn(func(p *Description) {
		if p.Dec == nil {
			p.Dec = make(map[string]interface{})
		}
		p.Dec[key] = val
	})
}

func DiffMatchPatch(key string, patch string) Patch {
	return patchFn(func(p *Description) {
		if p.DiffMatchPatch == nil {
			p.DiffMatchPatch = make(map[string]string)
		}
		p.DiffMatchPatch[key] = patch
	})
}
