package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	. "github.com/wushilin/portforwarder_go/logging"
	"github.com/wushilin/portforwarder_go/worker"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join([]string(*i), ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var binding arrayFlags
var logLevel int
var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])

	flag.PrintDefaults()
}

func main() {
	flag.Var(&binding, "b", "(multiple) Binding in [bind_address]:port::target_address:port format (e.g. :22::remote.host.com:22)")
	flag.IntVar(&logLevel, "loglevel", 0, "Log level (0 for debug, higher is less)")
	flag.Parse()
	if len(binding) == 0 {
		Usage()
		os.Exit(1)
	}

	SetLogLevel(logLevel)
	INFO("Setting logging level to %d", logLevel)
	wg := new(sync.WaitGroup)
	workers := make([]*worker.WorkerConfig, 0)
	for _, spec := range binding {
		config, err := parse(spec)
		if err != nil {
			FATAL("Error: %s", err)
			os.Exit(1)
		}
		wg.Add(1)
		workers = append(workers, config)
		go config.Start(wg)
	}
	go report(workers)

	wg.Wait()
}

func parse(input string) (config *worker.WorkerConfig, err error) {
	indexCC := strings.Index(input, "::")
	defaultError := errors.New(fmt.Sprintf("Invalid spec: %s", input))
	if indexCC == -1 {
		return nil, defaultError
	}

	listenSpec := input[:indexCC]
	targetSpec := input[indexCC+2:]

	listenIndexC := strings.Index(listenSpec, ":")
	if listenIndexC == -1 {
		return nil, defaultError
	}
	listen := listenSpec[:listenIndexC]
	listenPortStr := listenSpec[listenIndexC+1:]
	listenPort, err := strconv.Atoi(listenPortStr)
	if err != nil {
		return nil, defaultError
	}

	targetIndexC := strings.Index(targetSpec, ":")
	if targetIndexC == -1 {
		return nil, defaultError
	}
	target := targetSpec[:targetIndexC]
	targetPortStr := targetSpec[targetIndexC+1:]
	targetPort, err := strconv.Atoi(targetPortStr)
	if err != nil {
		return nil, defaultError
	}
	return &worker.WorkerConfig{
		BindAddress: listen,
		BindPort:    listenPort,
		TargetHost:  target,
		TargetPort:  targetPort,
	}, nil
}

func report(ws []*worker.WorkerConfig) {
	for {
		for _, next := range ws {
			INFO("* STATUS for %s:%d->%s:%d Up %d b; Down %d b; Active %d; Total %d",
				next.BindAddress, next.BindPort, next.TargetHost, next.TargetPort,
				next.Uploaded, next.Downloaded, next.Active, next.TotalHandled)
		}
		time.Sleep(30 * time.Second)
	}
}
