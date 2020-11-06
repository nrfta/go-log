package log

import (
	"context"
	. "github.com/onsi/ginkgo"
	g "github.com/onsi/gomega" // gomega.Panic function collides with log.Panic
)

var _ = Describe("Context", func() {
	Context("Context", func() {
		It("should provide context with value", func() {
			ctx := WithContext(context.Background())
			g.Expect(ctx.Value(ContextKeyLogFields)).ToNot(g.BeNil())
		})

		It("should save field value in context", func() {
			ctx := WithContext(context.Background())
			AddContextField(ctx, "foo", 5)
			fields := GetContextFields(ctx)
			g.Expect(fields["foo"]).To(g.Equal(5))
		})
	})
})
