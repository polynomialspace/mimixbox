//
// mimixbox/internal/applets/textutils/sha512sum/sha512sum.go
//
// Copyright 2021 Naohiro CHIKAMATSU
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package sha512sum

import (
	"crypto/sha512"
	"os"
	"strings"

	mb "github.com/nao1215/mimixbox/internal/lib"

	"github.com/jessevdk/go-flags"
)

const cmdName string = "sha512sum"

const version = "1.0.0"

var osExit = os.Exit

// Exit code
const (
	ExitSuccess int = iota // 0
	ExitFailuer
)

type options struct {
	Check   bool `short:"c" long:"check" description:"Check if the SHA1 value matches the file"`
	Version bool `short:"v" long:"version" description:"Show sha512sum command version"`
}

func Run() (int, error) {
	var opts options
	var args []string
	var err error = nil
	hash := sha512.New()

	if args, err = parseArgs(&opts); err != nil {
		return ExitFailuer, nil
	}

	if mb.HasPipeData() {
		err = mb.ChecksumOutput(hash, strings.NewReader(args[0]), "-")
		if err != nil {
			return ExitFailuer, err
		}
		return ExitSuccess, nil
	}

	if len(args) == 0 || mb.Contains(args, "-") {
		err = mb.ChecksumOutput(hash, os.Stdin, "-")
		if err != nil {
			return ExitSuccess, nil
		}
		return ExitSuccess, nil
	}

	if opts.Check {
		err = mb.CompareChecksum(hash, args)
		if err != nil {
			return ExitFailuer, err
		}
		return ExitSuccess, nil
	}

	for _, path := range args {
		p := os.ExpandEnv(path)
		r, err := os.Open(p)
		if err != nil {
			return ExitFailuer, err
		}
		defer r.Close()

		if err := mb.ChecksumOutput(hash, r, p); err != nil {
			return ExitSuccess, nil
		}
		hash.Reset()
	}
	return ExitSuccess, nil
}

func parseArgs(opts *options) ([]string, error) {
	p := initParser(opts)

	args, err := p.Parse()
	if err != nil {
		return nil, err
	}

	if mb.HasPipeData() {
		stdin, err := mb.FromPIPE()
		if err != nil {
			return nil, err
		}
		return []string{stdin}, nil
	}

	if opts.Version {
		mb.ShowVersion(cmdName, version)
		osExit(ExitSuccess)
	}

	return args, nil
}

func initParser(opts *options) *flags.Parser {
	parser := flags.NewParser(opts, flags.Default)
	parser.Name = cmdName
	parser.Usage = "[OPTIONS] FILE_PATH"

	return parser
}
