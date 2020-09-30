package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fgahr/mini"
)

// DBConf demonstrates a simple configuration with field tags.
type DBConf struct {
	Host     string `ini:"host"` // This field will show up as `host = ...`
	Port     int    `ini:"port"` // Likewise this one as `port = ...`
	User     string `ini:"user"`
	Password string `ini:"password"`
}

// NetConf demonstrates some of the advanced options.
type NetConf struct {
	Port int // no tag, field name will be used; same as `ini:"Port"`
	// Determined at runtime
	Hostname string `ini:"-"` // ignored by mini, will not be read or written
	UseSSL   bool   `ini:"use_ssl"`
}

type logLevel int

const (
	LogOff logLevel = iota
	LogError
	LogWarn
	LogInfo
	LogDebug
	LogTrace
)

// NOTE: plain value receiver
func (l logLevel) ToINI() string {
	switch l {
	case LogOff:
		return "off"
	case LogError:
		return "error"
	case LogWarn:
		return "warn"
	case LogInfo:
		return "info"
	case LogDebug:
		return "debug"
	case LogTrace:
		return "trace"
	default:
		panic("unknown log level")
	}
}

// NOTE: pointer receiver
func (l *logLevel) FromINI(s string) error {
	switch strings.ToLower(s) {
	case "off":
		*l = LogOff
	case "error":
		*l = LogError
	case "warn":
		*l = LogWarn
	case "info":
		*l = LogInfo
	case "debug":
		*l = LogDebug
	case "trace":
		*l = LogTrace
	default:
		panic("unknown log level specification: " + s)
	}
	return nil
}

type LogConf struct {
	Level logLevel `ini:"level"`
}

// Config shows a typical top-level configuration with several sections.
type Config struct {
	DB  DBConf  `ini:"db"`
	Net NetConf `ini:"net"`
	Log LogConf `ini:"logger"`
}

func defaults() Config {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}

	return Config{
		DB: DBConf{
			Host: "localhost",
			Port: 3306,
			User: "me",
			// Lines starting with semicolons are treated as comments.
			// They can appear anywhere else, though.
			// This password will be written and read correctly.
			Password: "super;secure",
		},
		Net: NetConf{
			Port:     8080,
			Hostname: hostname,
			UseSSL:   true,
		},
		Log: LogConf{
			Level: LogError,
		},
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [read|defaults]\n", os.Args[0])
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	// NOTE: ignoring possible errors here
	switch os.Args[1] {
	case "read":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "usage: %s read [file]\n", os.Args[0])
			os.Exit(2)
		}
		// This is how you can read from a config file
		file, err := os.Open(os.Args[2])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(3)
		}
		defer file.Close()
		conf := Config{}
		mini.Read(file, &conf)
	case "pipe":
		// Read a config and write it to output
		// Not useful outside a demonstration such as this
		conf := Config{}
		mini.Read(os.Stdin, &conf)
		mini.Write(os.Stdout, conf)
	case "defaults":
		// Write the default config
		mini.Write(os.Stdout, defaults())
	default:
		usage()
		os.Exit(1)
	}
}
