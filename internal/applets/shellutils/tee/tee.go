package tee

import (
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/jessevdk/go-flags"
	mb "github.com/nao1215/mimixbox/internal/lib"
)

type options struct {
	Version      bool `short:"v" long:"version" description:"print version and exit"`
	Append       bool `short:"a" long:"append" description:"append to files, do not overwrite"`
	IgnoreSIGINT bool `short:"i" long:"ignore-interrupts" description:"ignore SIGINT"`
}

const cmdName string = "tee"
const version = "1.0.1"

func initParser(opts *options) *flags.Parser {
	parser := flags.NewParser(opts, flags.Default)
	parser.Name = cmdName
	parser.Usage = "[OPTIONS]"

	return parser
}

func Run() (int, error) {
	var opts options
	args, err := initParser(&opts).Parse()
	if err != nil {
		return mb.ExitFailure, err
	}

	if opts.Version {
		mb.ShowVersion(cmdName, version)
		os.Exit(mb.ExitSuccess)
	}

	openFlags := os.O_WRONLY | os.O_CREATE
	if opts.Append {
		openFlags |= os.O_APPEND
	}

	if opts.IgnoreSIGINT {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT)
		go func() {
			for range sigs {
			}
		}()

	}

	files := make([]io.Writer, 0, 1+len(args))
	files = append(files, os.Stdout)
	for _, filename := range args {
		f, err := os.OpenFile(filename, openFlags, 0644)
		if err != nil {
			return mb.ExitFailure, err
		}
		files = append(files, f)
		defer f.Close()
	}

	multiwriter := io.MultiWriter(files...)
	_, err = io.Copy(multiwriter, os.Stdin)
	if err != nil {
		return mb.ExitFailure, err
	}

	return mb.ExitSuccess, nil
}
