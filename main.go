package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/abice/go-enum/generator"
	"github.com/mkideal/cli"
)

type rootT struct {
	cli.Helper
	FileNames []string `cli:"*f,file" usage:"The file(s) to generate enums.  Use more than one flag for more files."`
	NoPrefix  bool     `cli:"noprefix" usage:"Prevents the constants generated from having the Enum as a prefix."`
	Lowercase bool     `cli:"lower" usage:"Adds lowercase variants of the enum strings for lookup."`
	Marshal   bool     `cli:"marshal" usage:"Adds text marshalling functions."`
	Prefix    string   `cli:"prefix" usage:"Replaces the prefix with a user one."`
}

func main() {
	cli.Run(new(rootT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*rootT)

		for _, fileName := range argv.FileNames {

			g := generator.NewGenerator()

			if argv.NoPrefix {
				g.WithNoPrefix()
			}
			if argv.Lowercase {
				g.WithLowercaseVariant()
			}
			if argv.Marshal {
				g.WithMarshal()
			}

			originalName := fileName

			ctx.String("go-enum started. file: %s\n", ctx.Color().Cyan(originalName))
			fileName, _ = filepath.Abs(fileName)
			outFilePath := fmt.Sprintf("%s_enum.go", strings.TrimSuffix(fileName, filepath.Ext(fileName)))

			// Parse the file given in arguments
			raw, err := g.GenerateFromFile(fileName)
			if err != nil {
				return fmt.Errorf("Error while generating enums\nInputFile=%s\nError=%s\n", ctx.Color().Cyan(fileName), ctx.Color().RedBg(err))
			}

			mode := int(0644)
			err = ioutil.WriteFile(outFilePath, raw, os.FileMode(mode))
			if err != nil {
				return fmt.Errorf("Error while writing to file %s: %s\n", ctx.Color().Cyan(outFilePath), ctx.Color().Red(err))
			}
			ctx.String("go-enum finished. file: %s\n", ctx.Color().Cyan(originalName))
		}

		return nil
	})
}
