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
			ctx := WithContext(context.Background(), MakeField("foo", 5))
			fields := GetContextFields(ctx)
			g.Expect(len(fields)).To(g.Equal(1))
			g.Expect(fields["foo"]).To(g.Equal(5))
		})

		It("should push and pop fields", func() {
			ctx := WithContext(context.Background(), MakeField("foo", 5))
			fields := GetContextFields(ctx)
			g.Expect(len(fields)).To(g.Equal(1))
			g.Expect(fields["foo"]).To(g.Equal(5))

			func1(ctx)

			fields = GetContextFields(ctx)
			g.Expect(len(fields)).To(g.Equal(1))
			g.Expect(fields["foo"]).To(g.Equal(5))

			// pop should be safe to call even beyond actual stack height
			PopContextFields(ctx)
			PopContextFields(ctx)
			PopContextFields(ctx)
		})

		It("should safely ignore when context is not initialized", func() {
			ctx := context.Background()

			fields := GetContextFields(ctx)
			g.Expect(len(fields)).To(g.Equal(0))

			PushContextFields(ctx, MakeField("foo", 0))
			PopContextFields(ctx)

			ctx = context.WithValue(ctx, ContextKeyLogFields, "wrong")
			PopContextFields(ctx)
			PopContextFields(ctx)
		})
	})
})

func func1(ctx context.Context) {
	PushContextFields(ctx, MakeField("bar", 6))
	defer PopContextFields(ctx)

	fields := GetContextFields(ctx)
	g.Expect(len(fields)).To(g.Equal(2))
	g.Expect(fields["foo"]).To(g.Equal(5))
	g.Expect(fields["bar"]).To(g.Equal(6))

	func2(ctx)

	fields = GetContextFields(ctx)
	g.Expect(len(fields)).To(g.Equal(2))
	g.Expect(fields["foo"]).To(g.Equal(5))
	g.Expect(fields["bar"]).To(g.Equal(6))
}

func func2(ctx context.Context) {
	WithContext(ctx, MakeField("taz", 7))
	defer PopContextFields(ctx)

	fields := GetContextFields(ctx)
	g.Expect(len(fields)).To(g.Equal(3))
	g.Expect(fields["foo"]).To(g.Equal(5))
	g.Expect(fields["bar"]).To(g.Equal(6))
	g.Expect(fields["taz"]).To(g.Equal(7))
}
