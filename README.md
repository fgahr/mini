# mini - a library for INI (config) file handling

## What is it?

It's a simple library you can use in your project to serialize configuration
in a simple INI format.

## What is it not?

An attempt to cover all the possible intricacies of this (underspecified)
format. See [go-ini](https://github.com/go-ini/ini) for something more complete.

## Why, then?

Because I wanted something very simple with no dependencies. And I had been
looking for an opportunity to use reflection in Go.

## Where can I try it?

There is a `demo/demo.go` file in this repository that you can build with
`go build`. Then, for instance you can try

```sh
# build
$ go build demo.go
# write the default config, then read it back and write it again
# relax, I'm aware this makes no sense, it's just a demo
$ ./demo defaults | ./demo pipe
[db]
host = localhost
...
```

or whatever sheenanigans you may come up with. Inspect the source code for a
quick tour of the capabilities.

## OK, so how do I use it?

With the go toolchain installed, run

```sh
go get -u github.com/fgahr/mini
```

then import it into your projects. Usage is similar to standard library packages
such as `encoding/json` but less featureful.

This is a portion from the `demo/demo.go` file that demonstrates basic usage.

```go
import "github.com/fgahr/mini"

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

// ...

// Config shows a typical top-level configuration with several sections.
type Config struct {
	DB  DBConf  `ini:"db"`
	Net NetConf `ini:"net"`
	// ...
}
```

Organization into sections and then entries is assumed. If you need more than
that, look for YAML or JSON. However, if you are brave, you can do it with
custom types.

Supported field types are strings, floats (32 and 64), and all kinds of built-in
integers (int8 up to uint64), as well as booleans and `time.Duration` instances.
If a type has a `String()` method, that will be used for serialization.
Deserialization however, is a different matter.

## Custom types

Another portion from `demo/demo.go`

```go
type logLevel int

const (
	LogOff logLevel = iota
	LogError
	// ...
)

// NOTE: plain value receiver
func (l logLevel) ToINI() string {
	switch l {
	case LogOff:
		return "off"
	case LogError:
		return "error"
	// ...
	default:
		panic("unknown log level")
	}
}

// NOTE: pointer receiver
func (l *logLevel) FromINI(s string) error {
	switch s {
	case "off":
		*l = LogOff
	case "error":
		*l = LogError
	// ...
	default:
		panic("unknown log level specification: " + s)
	}
	return nil
}

type LogConf struct {
	Level logLevel `ini:"level"`
}

```

Note the different types of receiver on the `ToINI` and `FromINI` methods. This
is NOT an interface! DISCLAIMER: as I get more familiar with reflection in
Golang I may find a better way to do this.

## Where are the tests?

So far I've only done manual testing. I do plan to add tests as soon as I figure
out an efficient way to do that.
