package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
)

var configfile = flag.String("conf", "pmc.conf", "The config file location")

// Verify holds the configurations to be passed to a input plugin
type verify struct {
	Type    string   //Usually the name of the plugin
	Label   string   //A label for this verification
	Every   duration //A time.Duration parseble by time.ParseDuration
	Timeout duration //A time.Duration timeout value
	Config  string   //Any configuration that the plugin needs to be executade
}

// Register holds the configurations to be passed to a output plugin
type register struct {
	Type   string //Usually the name of the plugin
	Config string //Any configuration that the plugin needs to be executade
}

// Config is the main configuration format
type Config struct {
	Verify   []verify
	Register []register
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

const (
	statusOk int = iota
	statusFail
	statusTimeout
)

type record struct {
	Type      string
	Label     string
	Host      string
	Starttime time.Time
	Endtime   time.Time
	Status    int
}

func actionWrapper(ver verify, inp func(verify) record, outputs map[string]func(record) error) func() {
	return func() {
		for {
			rec := inp(ver)
			for _, reg := range outputs {
				err := reg(rec)
				if err != nil {
					log.Fatal(err)
				}
			}
			timer := time.NewTimer(ver.Every.Duration)
			<-timer.C
		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//flags de configuração inicial
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n  %s [options] \n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	var config Config
	if _, err := toml.DecodeFile(*configfile, &config); err != nil {
		log.Fatal(err)
	}

	//Valid verifiers list
	var availInputs = make(map[string]func(verify) record)
	availInputs["http"] = httpVerify
	availInputs["dns"] = dnsVerify
	availInputs["ping"] = pingVerify

	//Valid registers list
	var availOutputs = make(map[string]func(register) func(record) error)
	availOutputs["text"] = textRegister
	availOutputs["influxdb"] = influxdbRegister

	//Validating config registers
	var outputs = make(map[string]func(record) error)
	for _, reg := range config.Register {
		found := false
		for otype := range availOutputs {
			if otype == strings.ToLower(reg.Type) {
				outputs[otype] = availOutputs[otype](reg)
				found = true
				break
			}
		}
		if !found {
			log.Printf("Register type not available: %s", strings.ToLower(reg.Type))
		}
	}

	//Validating config verifiers
	var verifiers []func()
	for _, ver := range config.Verify {
		found := false
		for itype := range availInputs {
			if itype == strings.ToLower(ver.Type) {
				verifiers = append(verifiers, actionWrapper(ver, availInputs[itype], outputs))
				found = true
				break
			}
		}
		if !found {
			log.Printf("Verifier type not available: %s", strings.ToLower(ver.Type))
		}
	}

	//Main loop until break signal
	for _, ver := range verifiers {
		go ver()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
loop:
	for {
		select {
		case <-c:
			fmt.Println("get interrupted")
			break loop
		}
	}
	signal.Stop(c)
}
