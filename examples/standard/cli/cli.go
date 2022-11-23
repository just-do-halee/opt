package cli

import (
	"fmt"

	"github.com/just-do-halee/opt"
)

func Execute() {
	err :=
		opt.Args().
			Version("v0.1.0").
			Author("just-do-halee <just.do.halee@gmail.com>").
			About("This is a test cli app program.").
			Build(new(Textify))

	if err != nil {
		fmt.Print(err)
	}
}
