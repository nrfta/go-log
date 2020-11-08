package log

import "context"

// ContextKeyLogFields is the key for the logging fields context value.
const ContextKeyLogFields = "nrfta/go-log/Fields"

// WithContext adds logging `Fields` as a context value to the parent context and returns the new context.
func WithContext(parent context.Context) context.Context {
	if parent.Value(ContextKeyLogFields) != nil {
		return parent
	}
	return context.WithValue(parent, ContextKeyLogFields, makeFieldStack())
}

func PushContextFields(ctx context.Context, fields ...Field) {
	stack := getStack(ctx)
	stack.push(fields)
}

func PopContextFields(ctx context.Context) {
	stack := getStack(ctx)
	stack.pop()
}

// GetContextFields retrieves the logging `Fields` from context. GetContextFields panics if the logging fields
// could not be found in context or if the value cannot be cast to `Fields`.
func GetContextFields(ctx context.Context, additionalFields ...Field) Fields {
	stack := getStack(ctx)
	fields := stack.allFields()
	for _, f := range additionalFields {
		fields[f.Name] = f.Value
	}
	return fields
}

func getStack(ctx context.Context) *fieldStack {
	stackObj := ctx.Value(ContextKeyLogFields)
	if stackObj == nil {
		panic("logging fields have not been added to context yet; call WithContext")
	}
	stack, ok := stackObj.(*fieldStack)
	if !ok {
		panic("logging fields are not of the correct type")
	}
	return stack
}
