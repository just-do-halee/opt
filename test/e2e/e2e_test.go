package e2e_test

import (
	"testing"

	"github.com/just-do-halee/lum"
	"github.com/just-do-halee/opt"
)

type basis struct {
	Input   opt.Argument[string] `msg:"Input path"`
	Output  opt.Option[string]   `msg:"Output path" opt:"s,l"`
	Verbose opt.Option[int]      `msg:"Verbosity level -vv.." opt:"s,l,o"`

	Help        opt.Option[opt.Help]  `msg:"Show help message" opt:"l,s"`
	HelpCommand opt.Command[opt.Help] `rename:"help"`
}

type Case int

const (
	TestRequiredArgument Case = 1 + iota
)

func TestE2E(t *testing.T) {
	type Arg struct {
		args []string
		info struct {
			tree   opt.Tree
			caseId Case
		}
	}
	type Ret = string
	type Ctx = *lum.Context[*Arg, Ret]
	lum.Batch[*Arg, Ret]{
		{
			Name: "Putting Help from -h",
			Args: &Arg{args: []string{"-h"}},
			Pass: func(c Ctx) {
				c.AssertResultEqual(c.Arguments.info.tree.ToHelp())
			},
			Parallel: true,
		},
		{
			Name: "Putting Help from --help",
			Args: &Arg{args: []string{"--help"}},
			Pass: func(c Ctx) {
				c.AssertResultEqual(c.Arguments.info.tree.ToHelp())
			},
			Parallel: true,
		},
		{
			Name: "Putting Help from Command",
			Args: &Arg{args: []string{"help"}},
			Pass: func(c Ctx) {
				c.AssertResultEqual(c.Arguments.info.tree.ToHelp())
			},
			Parallel: true,
		},
		{
			Name: "Required Argument",
			Args: &Arg{args: []string{"./test.txt"}},
			Pass: func(c Ctx) {
				if c.Arguments.info.caseId == TestRequiredArgument {
					c.AssertResultEqual("")
				}
			},
			Parallel: true,
		},
		{
			Name: "Default Argument",
			Args: &Arg{args: []string{""}},
			Pass: func(c Ctx) {
				c.AssertResultEqual("")
			},
			Parallel: true,
		},
		{
			Name: "Option without Value (short)",
			Args: &Arg{args: []string{"-o"}},
			Pass: func(c Ctx) {

			},
			Parallel: true,
		},
		{
			Name: "Option without Value (long)",
			Args: &Arg{args: []string{"--output"}},
			Pass: func(c Ctx) {

			},
			Parallel: true,
		},
		{
			Name: "Required short Option (string)",
			Args: &Arg{args: []string{"-o", "./test.txt"}},
			Pass: func(c Ctx) {

			},
			Parallel: true,
		},
		{
			Name: "Required long Option (string)",
			Args: &Arg{args: []string{"--output", "./test.txt"}},
			Pass: func(c Ctx) {

			},
			Parallel: true,
		},
		{
			Name: "Default short Option (string)",
			Args: &Arg{args: []string{""}},
			Pass: func(c Ctx) {
				c.AssertResultEqual("")
			},
			Parallel: true,
		},
		{
			Name: "Default long Option (string)",
			Args: &Arg{args: []string{""}},
			Pass: func(c Ctx) {
				c.AssertResultEqual("")
			},
			Parallel: true,
		},
		{
			Name: "Required Occurrence short Option (int)",
			Args: &Arg{args: []string{"-vvv"}},
			Pass: func(c Ctx) {

			},
			Parallel: true,
		},
	}.Run(t, "End To End",
		func(a *Arg) Ret {
			// before each

			var op basis

			// set defaults
			switch a.info.caseId {
			default:
				op.Input.Set("./test.txt")
				op.Output.Set("./test.txt")
				op.Verbose.Set(3)
			}

			err, tree := opt.Args(a.args...).BuildForDebugging(&op)
			a.info.tree = tree
			if err == nil {
				return ""
			}
			return err.Error()
		},
		func(c Ctx) {
			// after each
		})
}
