package log

import "context"

// ContextKeyLogFields is the key for the logging fields context value.
const ContextKeyLogFields = "nrfta/go-log/Fields"

// WithContext initializes context with a logging fields stack with the given fields. If the given
// context has already bene initialized, then the fields are pushed onto the existing stack.
func WithContext(parent context.Context, fields ...Field) context.Context {
	if parent.Value(ContextKeyLogFields) != nil {
		PushContextFields(parent, fields...)
		return parent
	}
	return context.WithValue(parent, ContextKeyLogFields, makeFieldStack().push(fields))
}

// PushContextFields pushes the given fields onto the logging fields stack. PushContextFields panics if the
// context has not been initialized via WithContext.
func PushContextFields(ctx context.Context, fields ...Field) {
	stack := getStack(ctx)
	stack.push(fields)
}

// PopContextFields pops the last entry off of the logging fields stack. PopContextFields panics if the
// context has not been initialized via WithContext.
func PopContextFields(ctx context.Context) {
	stack := getStack(ctx)
	stack.pop()
}

// GetContextFields retrieves the logging `Fields` from context. GetContextFields panics if the
// context has not been initialized via WithContext.
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
