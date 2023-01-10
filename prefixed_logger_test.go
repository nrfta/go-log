package log

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	g "github.com/onsi/gomega"
)

var _ = Describe("PrefixedLogger", func() {
	It("should prefix the logger instance", func() {
		buf := bytes.Buffer{}
		logger := New(true, "info")
		logger.Out = &buf

		pl := NewPrefixedLogger("Test", logger)
		pl2 := NewPrefixedLogger("Foo", logger)

		pl.Info("info log")
		pl2.Warnf("a warning %s", "log")
		logger.Info("normal log")

		//
		b := ByteLogs{} // buffer can be passed in
		b.Parse(&buf)   // or will be stored if passed to this function
		g.Expect(len(b.Parsed)).To(g.Equal(3))
		g.Expect(b.LogInLogs("msg", "Test: info log")).To(g.BeTrue())
		g.Expect(b.LogInLogs("msg", "Foo: a warning log")).To(g.BeTrue())
		g.Expect(b.LogInLogs("msg", "normal log")).To(g.BeTrue())

		// log not yet in log
		lateLog := "another log!"
		g.Expect(b.LogInLogs("msg", lateLog)).To(g.BeFalse())
		logger.Info(lateLog)
		// log will be in buffer ref, but not parsed yet
		g.Expect(b.LogInLogs("msg", lateLog)).To(g.BeFalse())
		b.Parse(nil) // parse stored pointer to buffer
		g.Expect(b.LogInLogs("msg", lateLog)).To(g.BeTrue())
	})

	It("should create an info log by default if one is not given", func() {
		buf := bytes.Buffer{}
		pl := NewPrefixedLogger("test", nil)
		pl.LoggerInstance.Out = &buf
		pl.Info("goldfish")

		b := ByteLogs{Log: &buf}
		b.Parse(nil)
		g.Expect(b.Parsed[0]["level"]).To(g.Equal("info"))
		g.Expect(b.Parsed[0]["msg"]).To(g.Equal("test: goldfish"))
	})
})
