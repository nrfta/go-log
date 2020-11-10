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

// PushContextFields pushes the given fields onto the logging fields stack.
func PushContextFields(ctx context.Context, fields ...Field) {
	stack := getStack(ctx)
	if stack == nil {
		return
	}
	stack.push(fields)
}

// PopContextFields pops the last entry off of the logging fields stack.
func PopContextFields(ctx context.Context) {
	stack := getStack(ctx)
	if stack == nil {
		return
	}
	stack.pop()
}

// GetContextFields retrieves the logging `Fields` from context. GetContextFields returns an empty Fields map
// if the context has not been initialized by calling WithContext.
func GetContextFields(ctx context.Context, additionalFields ...Field) Fields {
	stack := getStack(ctx)
	if stack == nil {
		return make(Fields)
	}
	fields := stack.allFields()
	for _, f := range additionalFields {
		fields[f.Name] = f.Value
	}
	return fields
}

func getStack(ctx context.Context) *fieldStack {
	stackObj := ctx.Value(ContextKeyLogFields)
	if stackObj == nil {
		Warn("context logging fields not initialized; call log.WithContext")
		return nil
	}
	stack, ok := stackObj.(*fieldStack)
	if !ok {
		Warn("context logging fields has incorrect type")
	}
	return stack
}
