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

		pl := NewPrefixedLogger("Test")
		pl2 := NewPrefixedLogger("Foo")

		pl.Info("info log")
		pl2.Warnf("a warning %s", "log")
		logger.Info("normal log")

		//
		g.Expect(1).To(g.Equal(1))
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
})
