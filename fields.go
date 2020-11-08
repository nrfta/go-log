package log

type Field struct {
	Name string
	Value interface{}
}

type fieldStackItem struct {
	fields []Field
}
type fieldStack struct {
	items []*fieldStackItem
}

func MakeField(name string, value interface{}) Field {
	return Field{name, value}
}

func makeFieldStackItem(fields []Field) *fieldStackItem {
	return &fieldStackItem{fields}
}

func makeFieldStack() *fieldStack {
	return &fieldStack{}
}

func (s *fieldStack) push(fields []Field) {
	item := makeFieldStackItem(fields)
	s.items = append(s.items, item)
}

func (s *fieldStack) pop() {
	if len(s.items) > 0 {
		s.items = s.items[:len(s.items)-1]
	}
}

func (s *fieldStack) allFields() Fields {
	fields := make(Fields)
	for _, i := range s.items {
		for _, f := range i.fields {
			fields[f.Name] = f.Value
		}
	}
	return fields
}
