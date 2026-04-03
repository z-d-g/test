package app

import (
	"flag"
	"fmt"
	"os"

	"github.com/z-d-g/md-cli/internal/utils"
)

type CLIArgs struct {
	PrintOnly bool
	Help      bool
	Version   bool
	Files     []string
}

func ParseCLIArgs(args []string) CLIArgs {
	var printFlag, helpFlag, versionFlag bool
	fs := flag.NewFlagSet("md-cli", flag.ContinueOnError)
	fs.BoolVar(&printFlag, "p", false, "Print rendered markdown to stdout")
	fs.BoolVar(&printFlag, "print", false, "Print rendered markdown to stdout")
	fs.BoolVar(&helpFlag, "h", false, "Show help")
	fs.BoolVar(&helpFlag, "help", false, "Show help")
	fs.BoolVar(&versionFlag, "v", false, "Show version")
	fs.BoolVar(&versionFlag, "version", false, "Show version")

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return CLIArgs{Help: true}
		}
		PrintUsage()
		os.Exit(1)
	}

	if versionFlag {
		return CLIArgs{Version: true}
	}

	return CLIArgs{
		PrintOnly: printFlag,
		Help:      helpFlag,
		Version:   false,
		Files:     utils.FilterMarkdownFiles(fs.Args()),
	}
}

func PrintUsage() {
	fmt.Println("Usage: md-cli [options] [file.md]")
	fmt.Println("Options:")
	fmt.Println("  -p, --print      Print rendered markdown to stdout")
	fmt.Println("  -v, --version    Show version information")
	fmt.Println("  -h, --help       Show help")
}
