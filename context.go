package log

import "context"

const ContextKeyLogFields = "nrfta/go-log/Fields"

// WithContext adds logging `Fields` as a context value to the parent context and returns the new context.
func WithContext(parent context.Context) context.Context {
	return context.WithValue(parent, ContextKeyLogFields, make(Fields))
}

// AddContextField adds a new field to the `Fields` held in the logging fields context value.
func AddContextField(ctx context.Context, name string, value interface{}) {
	fields := GetContextFields(ctx)
	fields[name] = value
}

// GetContextFields retrieves the logging `Fields` from context. GetContextFields panics if the logging fields
// could not be found in context or if the value cannot be cast to `Fields`.
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
