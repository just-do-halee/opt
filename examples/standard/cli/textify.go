package cli

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/just-do-halee/opt"
)

type file string

type Textify struct {
	Input   opt.Argument[file] `msg:"Input path"`
	Output  opt.Option[file]   `msg:"Output path" opt:"s,l"`
	Verbose opt.Option[int]    `msg:"Verbosity level -vv.." opt:"s,l,o"`
	Cat     opt.Command[cat]   `msg:"Print file contents"`

	Help        opt.Option[opt.Help]  `opt:"l,s"`
	HelpCommand opt.Command[opt.Help] `rename:"help"`
}

func (o *Textify) Before() error {
	opt.Set(&o.Output, file("./output.txt"))
	opt.Set(&o.Verbose, 2)
	return nil
}

func (o *Textify) After() error {
	var err error

	err = opt.Validate(&o.Input, opt.IsFile[file])
	if err != nil {
		return err
	}
	err = opt.Validate(&o.Verbose, opt.IsMinMax(0, 3))
	if err != nil {
		return err
	}

	return nil
}

func (o *Textify) Run() error {
	input := o.Input.Get()
	output := o.Output.Get()

	verbose := o.Verbose.Get()
	if verbose > 0 {
		log.Println("Copying file:", input, "to", output)
	}

	bytesRead, err := ioutil.ReadFile(string(input))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(string(output), bytesRead, 0644)
	if err != nil {
		return err
	}

	return nil
}

// ----------------------------------------------
type cat struct {
	Parent opt.Parent[Textify]
	File   opt.Argument[file] `msg:"File to print"`
	Len    opt.Option[uint]   `msg:"Length of output" opt:"l,s"`
}

func (o *cat) Before() error {
	opt.Set(&o.Len, uint(10))
	return nil
}

func (o *cat) After() error {
	return opt.Validate(&o.File, opt.IsFile[file])
}

func (o *cat) Run() error {
	// not checked the range
	verbose := o.Parent.Get().Verbose.Get()

	println := func(a ...any) {
		if verbose > 0 {
			log.Println(a...)
		}
	}

	file := o.File.Get()

	println("Opening file:", file)
	f, err := os.OpenFile(string(file), os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	println("<---------Reading contents--------->")
	// print file contents
	buf := bufio.NewReader(f)
	for i := uint(0); i < o.Len.Get(); i++ {
		line, err := buf.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Print(line)
	}
	println(">---------End of contents----------<")
	return nil
}
