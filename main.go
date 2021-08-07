package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dhulihan/httpeeved/internal/inspect"
	"github.com/dhulihan/httpeeved/internal/selection"
	"github.com/elazarl/goproxy"
	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

type Opts struct {
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information: -v for debug, -vv for trace."`

	Addr string `short:"a" long:"addr" description:"Bind address (eg: 0.0.0.0:80)" default:":8080"`

	Codes []int `short:"c" long:"codes" default:"200" default:"206" default:"400" default:"404" default:"500" default:"502" description:"Repsonse status codes. Can be specified many times."`

	SelectionStrategy string `short:"s" long:"selection-strategy" default:"round-robin" choice:"round-robin" choice:"random" description:"response code selection strategy"`

	Responses map[int]string `short:"r" long:"responses" description:"use this to set a custom response message"`

	Proxy bool `short:"x" long:"proxy" description:"Run as proxy. httpeeved will forward requests to destination and modify the response."`
}

var (
	sel selection.SelectionStrategy
)

func main() {
	opts := &Opts{}
	args, err := flags.Parse(opts)
	if err != nil {
		panic(err)
	}

	if len(args) > 0 {
		log.WithField("args", args).Errorf("unexpected args")
	}

	switch len(opts.Verbose) {
	case 1:
		log.SetLevel(log.DebugLevel)
	case 2:
		log.SetLevel(log.TraceLevel)
	}

	fmt.Printf("Serving response codes: [%v]\n", opts.Codes)

	// set up response selection strategy
	switch opts.SelectionStrategy {
	case "random":
		log.Info("using random selection strategy")
		sel = selection.NewRandomSelectionStrategy(opts.Codes)
	default:
		log.Info("using round-robin selection strategy")
		sel = selection.NewRoundRobinSelectionStrategy(opts.Codes)
	}

	// run in proxy mode
	if opts.Proxy {
		log.Info("running in proxy mode")
		proxy := goproxy.NewProxyHttpServer()
		proxy.Verbose = len(opts.Verbose) > 0

		respCond := goproxy.RespConditionFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) bool {
			return true
		})

		proxy.OnResponse(respCond).DoFunc(func(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			origStatus := r.StatusCode
			newStatus := sel.Code()
			log.WithFields(log.Fields{
				"Method":       ctx.Req.Method,
				"URL":          ctx.Req.URL.String(),
				"UpstreamCode": origStatus,
				"ResponseCode": newStatus,
			}).Info("backend request completed")
			r.StatusCode = newStatus
			return r
		})

		// TODO: hook this up to gin's RunListener
		log.Fatal(http.ListenAndServe(opts.Addr, proxy))
		return
	}

	r := gin.Default()

	// catch EVERYTHING
	r.NoRoute(codeHandler)

	r.Run(opts.Addr)
}

func codeHandler(c *gin.Context) {
	code := sel.Code()
	log.WithField("code", code).Debug("generated code")

	resp := gin.H{
		"Code":         fmt.Sprintf("%d", code),
		"Method":       c.Request.Method,
		"URL":          c.Request.URL.String(),
		"Content-Type": c.Request.Header.Get("Content-Type"),
	}

	inspect.MultipartForm(c, resp)
	inspect.PostForm(c, resp)

	// request body
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Fatal(err)
	}

	// request body
	reqBody := string(b)
	if reqBody != "" {
		log.Debug(reqBody)
		resp["body"] = reqBody
	}

	c.JSON(code, resp)
}
