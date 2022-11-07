# ***`Opt`***

Opt is a package for building CLI apps in Go.

It is simple, and lightweight but also powerful.

It uses struct to parse everything.

```go

type Opt struct {
	Verbose  int                       `default:"2" description:"Verbosity" opt:"s,o"`
	Explorer opt.Command[*ExplorerOpt] `description:"Starts the Explorer"`
	Rest     opt.Command[*RestOpt]     `description:"Starts the REST API"`
	Help     opt.Help[*Opt]            `opt:"l,s"`
}


func main() {
 	var op *Opt
	err := opt.Args().Build(op)
	if err != nil {
		panic(err)
	}
}

```
