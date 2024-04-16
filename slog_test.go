package log

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	. "github.com/onsi/ginkgo"
	g "github.com/onsi/gomega"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

var _ = Describe("Logger", func() {
	Describe("NewSlogHTTPMiddleware", func() {
		It("should log http request information", func() {
			var (
				buf    bytes.Buffer
				logger = slog.New(slog.NewJSONHandler(&buf, nil))
				mw     = NewSLogChiMiddleware(logger)
				res    = "test"
				hf     = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(res))
				})
			)

			ts := httptest.NewServer(mw(hf))
			defer ts.Close()

			_, err := http.Get(ts.URL)
			g.Expect(err).To(g.Succeed())

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
			g.Expect(err).To(g.Succeed())

			g.Expect(lo.Msg).To(g.Equal("HTTP Request Served"))
			g.Expect(lo.Proto).To(g.Equal("HTTP/1.1"))
			g.Expect(lo.Path).To(g.Equal("/"))
			g.Expect(lo.Duration).To(g.BeNumerically(">", 0))
			g.Expect(lo.Status).To(g.Equal(200))
			g.Expect(lo.Size).To(g.Equal(len(res)))
			g.Expect(strings.Split(lo.IP, ":")[0]).To(g.Equal("127.0.0.1"))
		})
	})

	Describe("noOpVariablesScrubber#Scrub", func() {
		It("should return nil", func() {
			s := noopVariablesScrubber{}

			res := s.Scrub(map[string]any{"test": struct{}{}})

			g.Expect(res).To(g.BeNil())
		})
	})

	Describe("NewGraphQLResponseMiddleware", func() {
		It("should log GraphQL response and request info", func() {
			var (
				buf    bytes.Buffer
				logger = slog.New(slog.NewJSONHandler(&buf, nil))
				query  = "query"
				vars   = map[string]any{"token": "super secrect stuff"}
				oc     = &graphql.OperationContext{
					RawQuery:  query,
					Variables: vars,
				}
				errMsg = "Testing Errors"
				errors = gqlerror.List{
					&gqlerror.Error{
						Err:     errors.New("error"),
						Message: errMsg,
					},
				}
				handler = func(ctx context.Context) *graphql.Response {
					return &graphql.Response{
						Errors: errors,
					}
				}
				subject = NewSLogGraphQLResponseMiddleware(logger, nil)
			)

			subject(graphql.WithOperationContext(context.Background(), oc), handler)

			type logOutput struct {
				Msg     string `json:"msg"`
				Graphql struct {
					Req struct {
						Query     string         `json:"query"`
						Variables map[string]any `json:"variables"`
					} `json:"req"`
					Res struct {
						Errors string `json:"errors"`
					} `json:"res"`
					Duration int `json:"duration"`
				} `json:"graphql"`
			}

			var lo logOutput
			err := json.Unmarshal(buf.Bytes(), &lo)
			g.Expect(err).To(g.Succeed())

			g.Expect(lo.Msg).To(g.Equal("GraphQL Request Served"))
			g.Expect(lo.Graphql.Req.Query).To(g.Equal(query))
			g.Expect(lo.Graphql.Req.Variables).To(g.BeNil())
			g.Expect(lo.Graphql.Res.Errors).To(g.Equal(fmt.Sprintf("input: %s\n", errMsg)))
			g.Expect(lo.Graphql.Duration).To(g.BeNumerically(">", 0))
		})
	})
})
