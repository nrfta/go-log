package log

import (
	"bytes"

	asyncTest "github.com/nrfta/go-asynq-helpers/v2/pkg/tests"
	. "github.com/onsi/ginkgo"
	g "github.com/onsi/gomega"
)

var _ = Describe("PrefixedLogger", func() {
	It("should prefix the logger instance", func() {
		buf := bytes.Buffer{}
		logger := New(true, "info")
		logger.Out = &buf

		pl := NewPrefixedLogger("Test")
		pl2 := NewPrefixedLogger("Foo")

		pl.Info("info log")
		pl2.Warnf("a warning %s", "log")
		logger.Info("normal log")

		logs := asyncTest.ParseLogs(buf)
		g.Expect(len(logs)).To(g.Equal(3))
		g.Expect(asyncTest.LogInLogs("msg", "Test: info log", logs)).To(g.BeTrue())
		g.Expect(asyncTest.LogInLogs("msg", "Foo: a warning log", logs)).To(g.BeTrue())
		g.Expect(asyncTest.LogInLogs("msg", "normal log", logs)).To(g.BeTrue())
	})
})
