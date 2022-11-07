package main

import (
	"fmt"

	"github.com/just-do-halee/opt"
)

type Opt struct {
	Verbose  int                       `default:"3" description:"Verbosity" opt:"s,o"`
	Explorer opt.Command[*ExplorerOpt] `description:"Starts the explorer"`
	Rest     opt.Command[*RestOpt]     `description:"Starts the REST API"`
	Help     opt.Help[*Opt]            `opt:"l,s"`
}

type ExplorerOpt struct {
	Port int `short:"p" long:"port" description:"Port to listen on" default:"8080"`
}

func (o *ExplorerOpt) Run() error {
	fmt.Println(o.Port)
	return nil
}

type RestOpt struct {
	Port int `short:"p" long:"port" description:"Port to listen on" default:"4000"`
}

func (o *RestOpt) Run() error {
	fmt.Println(o.Port)
	return nil
}

func main() {

	err := opt.Args().Build(&Opt{})
	if err != nil {
		panic(err)
	}

}
