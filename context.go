package log

import "context"

const ContextKeyLogFields = "nrfta/go-log/Fields"

func WithContext(parent context.Context) context.Context {
	return context.WithValue(parent, ContextKeyLogFields, make(Fields))
}

func AddContextField(ctx context.Context, name string, value interface{}) {
	fields := GetContextFields(ctx)
	fields[name] = value
}

func GetContextFields(ctx context.Context) Fields {
	fields := ctx.Value(ContextKeyLogFields)
	if fields == nil {
		panic("logging fields have not been added to context yet; call WithContext")
	}
	f, ok := fields.(Fields)
	if !ok {
		panic("logging fields are not of the correct type")
	}
	return f
}
