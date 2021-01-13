package log

// Field represents a logging field.
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

// MakeField creates a new logging field.
func MakeField(name string, value interface{}) Field {
	return Field{name, value}
}

func FieldsToFieldsArray(fields Fields) (arr []Field) {
	for name, value := range fields {
		arr = append(arr, MakeField(name, value))
	}
	return
}

func makeFieldStackItem(fields []Field) *fieldStackItem {
	return &fieldStackItem{fields}
}

func makeFieldStack() *fieldStack {
	return &fieldStack{}
}

func (s *fieldStack) push(fields []Field) *fieldStack {
	item := makeFieldStackItem(fields)
	s.items = append(s.items, item)
	return s
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
