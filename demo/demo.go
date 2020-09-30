package main

import (
	"fmt"
	"os"

	"github.com/fgahr/mini"
)

type DBConf struct {
	Host     string `ini:"host"`
	Port     int    `ini:"port"`
	User     string `ini:"user"`
	Password string `ini:"password"`
}

type NetConf struct {
	Port int // no tag, field name will be used; same as `ini:"Port"`
	// Determined at runtime
	Hostname string `ini:"-"` // ignored by mini, will not be read or written
}

type Config struct {
	DB  DBConf  `ini:"db"`
	Net NetConf `ini:"net"`
}

func defaults() Config {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}

	return Config{
		DB: DBConf{
			Host:     "localhost",
			Port:     3306,
			User:     "me",
			Password: "super!secure",
		},
		Net: NetConf{
			Port:     8080,
			Hostname: hostname,
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
		file, err := os.Open(os.Args[2])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(3)
		}
		defer file.Close()
		conf := Config{}
		mini.Read(file, &conf)
	case "pipe":
		conf := Config{}
		mini.Read(os.Stdin, &conf)
		mini.Write(os.Stdout, conf)
	case "defaults":
		mini.Write(os.Stdout, defaults())
	default:
		usage()
		os.Exit(1)
	}
}
