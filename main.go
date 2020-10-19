package main

import (
	"fmt"

	"github.com/dhulihan/httpeeved/internal/selection"
	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

type Opts struct {
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information: -v for debug, -vv for trace."`

	Addr string `short:"a" long:"addr" description:"Address to bind too" default:":8080"`

	Codes []int `short:"c" long:"codes" default:"200" default:"206" default:"400" default:"404" default:"500" default:"502" description:"Repsonse status codes. Can be specified many times."`

	SelectionStrategy string `short:"s" long:"selection-strategy" default:"round-robin" choice:"round-robin" choice:"random" description:"response code selection strategy"`

	Responses map[int]string `short:"r" long:"responses" description:"use this to set a custom response message"`
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

	r := gin.Default()

	r.GET("/", codeHandler)
	r.HEAD("/", codeHandler)
	r.PUT("/", codeHandler)
	r.PATCH("/", codeHandler)
	r.DELETE("/", codeHandler)
	r.Run(opts.Addr) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func codeHandler(c *gin.Context) {
	code := sel.Code()
	log.WithField("code", code).Debug("generated code")
	c.JSON(code, gin.H{
		"message": fmt.Sprintf("%d", code),
	})
}
