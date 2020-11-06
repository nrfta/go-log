package log

import (
	"testing"

	. "github.com/onsi/ginkgo"
	g "github.com/onsi/gomega"
)

func TestGoLog(t *testing.T) {
	g.RegisterFailHandler(Fail)
	RunSpecs(t, "GoLog Suite")
}
