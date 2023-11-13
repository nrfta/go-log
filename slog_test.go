package log_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/nrfta/go-log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logger", func() {
	Describe("NewSlogHTTPMiddleware", func() {
		It("should log http request information", func() {
			var (
				buf    bytes.Buffer
				logger = slog.New(slog.NewJSONHandler(&buf, nil))
				mw     = log.NewSLogChiMiddleware(logger)
				res    = "test"
				hf     = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(res))
				})
			)

			ts := httptest.NewServer(mw(hf))
			defer ts.Close()

			_, err := http.Get(ts.URL)
			Expect(err).To(Succeed())

			type logOutput struct {
				Msg      string `json:"msg"`
				Proto    string `json:"proto"`
				Path     string `json:"path"`
				Duration int    `json:"duration"`
				Status   int    `json:"status"`
				Size     int    `json:"size"`
				IP       string `json:"ip"`
			}

			var lo logOutput
			err = json.Unmarshal(buf.Bytes(), &lo)
			Expect(err).To(Succeed())

			Expect(lo.Msg).To(Equal("HTTP Request Served"))
			Expect(lo.Proto).To(Equal("HTTP/1.1"))
			Expect(lo.Path).To(Equal("/"))
			Expect(lo.Duration).To(BeNumerically(">", 0))
			Expect(lo.Status).To(Equal(200))
			Expect(lo.Size).To(Equal(len(res)))
			Expect(strings.Split(lo.IP, ":")[0]).To(Equal("127.0.0.1"))
		})
	})
})
