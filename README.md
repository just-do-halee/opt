<div align="center">
  <img src="./.github/logo.png" alt="go-arg" height="170px">
  <h1>OPT</h1>
  <p>Opt is a package for building an organized CLI app in Go.</p>
  <br/>
</div>

```shell
go get -u github.com/just-do-halee/opt@latest
```

[![CI][ci-badge]][ci-url]
[![Go Reference](https://pkg.go.dev/badge/github.com/just-do-halee/opt.svg)](https://pkg.go.dev/github.com/just-do-halee/opt)
[![Licensed][license-badge]][license-url]
[![Twitter][twitter-badge]][twitter-url]

[ci-badge]: https://github.com/just-do-halee/opt/actions/workflows/ci.yml/badge.svg
[license-badge]: https://img.shields.io/github/license/just-do-halee/opt?color=blue
[twitter-badge]: https://img.shields.io/twitter/follow/do_halee?style=flat&logo=twitter&color=4a4646&labelColor=333131&label=just-do-halee
[ci-url]: https://github.com/just-do-halee/opt/actions
[twitter-url]: https://twitter.com/do_halee
[license-url]: https://github.com/just-do-halee/opt

| [Examples](./examples/) | [Latest Note](./CHANGELOG.md) |

```go
type file string

type Textify struct {
	Input   opt.Argument[file] `msg:"Input path"`
	Output  opt.Option[file]   `msg:"Output path" opt:"s,l"`
	Verbose opt.Option[int]    `msg:"Verbosity level -vv.." opt:"s,l,o"`
	Silent  opt.Option[bool]   `msg:"Silent mode" opt:"s,l"`
	Cat     opt.Command[cat]   `msg:"Print file contents"`

	Help        opt.Option[opt.Help]  `opt:"s,l"`
	HelpCommand opt.Command[opt.Help] `rename:"help"`
}

func (o *Textify) Before() error {
	opt.Set(&o.Output, file("./output.txt"))
	opt.Set(&o.Verbose, 2)
	opt.Set(&o.Silent, false)
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
	silent := o.Silent.Get()

	if verbose > 0 && !silent {
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
	Len    opt.Option[uint]   `msg:"Length of output" opt:"s,l"`
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
	p := o.Parent.Get()
	verbose := p.Verbose.Get()
	silent := p.Silent.Get()

	println := func(a ...any) {
		if verbose > 0 && !silent {
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
```

```go
func main() {
	err :=
		opt.Args().
			Version("v0.1.0").
			Author("just-do-halee <just.do.halee@gmail.com>").
			About("This is a CLI app program.").
			Build(new(Textify))

	if err != nil {
		fmt.Print(err)
	}
}
```
