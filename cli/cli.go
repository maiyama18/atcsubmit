package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/mui87/atcsubmit/atcoder"
)

const baseURL = "https://atcoder.jp"

type Cli struct {
	client *atcoder.Client

	contest  string
	problem  string
	file     string
	language string

	outStream io.Writer
	errStream io.Writer
}

func New(args []string, outStream, errStream io.Writer) (*Cli, error) {
	flags := flag.NewFlagSet("atcsubmit", flag.ContinueOnError)
	flags.SetOutput(errStream)
	flags.Usage = func() {
		_, _ = fmt.Fprintf(errStream, helpMsg)
		flags.PrintDefaults()
		_, _ = fmt.Fprintln(errStream, "")
	}

	// TODO: add silent option
	var (
		contest  string
		problem  string
		file     string
		language string
	)
	flags.StringVar(&contest, "contest", "", "contest you are challenging. e.g.) ABC051")
	flags.StringVar(&problem, "problem", "", "problem you are solving. e.g.) C")
	flags.StringVar(&file, "file", "", "filepath whose content will run. if not set, the content is got from standard input")
	flags.StringVar(&language, "language", "", "if set, the animation run from right to left")
	if err := flags.Parse(args[1:]); err != nil {
		return nil, fmt.Errorf("failed to parse command line options: %s", strings.Join(args[1:], " "))
	}

	if contest == "" {
		flags.Usage()
		return nil, errors.New("specify the contest you are challenging. e.g.) ABC051")
	}
	if problem == "" {
		flags.Usage()
		return nil, errors.New("specify the problem you are solving. e.g.) C")
	}
	if file == "" {
		flags.Usage()
		return nil, errors.New("specify the file to submit. e.g.) 'abc051-c.py'")
	}
	// TODO: accept only valid programming languages
	if language == "" {
		flags.Usage()
		return nil, errors.New("specify the programing language you use. e.g.) 'python3'")
	}

	// TODO: add 'atcsubmit login' command which persist the encoded username/password
	username := os.Getenv("ATCODER_USERNAME")
	if username == "" {
		flags.Usage()
		return nil, errors.New("specify your atcoder username to ATCODER_USERNAME env var. e.g.) export ATCODER_USERNAME=mui87")
	}
	password := os.Getenv("ATCODER_PASSWORD")
	if password == "" {
		flags.Usage()
		return nil, errors.New("specify your atcoder password to ATCODER_PASSWORD env var. e.g.) export ATCODER_PASSWORD=password")
	}

	// TODO: get timeout from command line
	client := atcoder.NewClient(baseURL, username, password, 45*time.Second, outStream, errStream)

	return &Cli{
		client: client,

		contest:  contest,
		problem:  problem,
		file:     file,
		language: language,

		outStream: outStream,
		errStream: errStream,
	}, nil
}

func (c *Cli) Run() error {
	return c.client.Submit(c.contest, c.problem, c.language, c.file)
}

const helpMsg = `
USAGE:
  atcsubmit is a command line tool to submit your code to AtCoder.

EXAMPLE:
  $ export ATCODER_USERNAME=mui87 ATCODER_PASSWORD=password
  $ atcsubmit -contest ABC124 -problem C -file abc124-c.py -language python

ENVIRONMENT_VARIABLES:
  ATCODER_USERNAME: your atcoder username.
  ATCODER_PASSWORD: your atcoder password.

OPTIONS:
`
