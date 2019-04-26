package atcoder

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type Client struct {
	baseURL  string
	username string
	password string

	ctx    context.Context
	cancel context.CancelFunc

	outStream io.Writer
	errStream io.Writer
}

func NewClient(baseURL, username, password string, timeout time.Duration, outStream, errStream io.Writer) *Client {
	ctx, cancel := chromedp.NewContext(context.Background())
	ctx, cancel = context.WithTimeout(ctx, timeout)

	return &Client{
		baseURL:   baseURL,
		username:  username,
		password:  password,
		ctx:       ctx,
		cancel:    cancel,
		outStream: outStream,
		errStream: errStream,
	}
}

func (c *Client) Submit(contest, problem, language, file string) error {
	defer c.cancel()

	code, err := readCode(file)
	if err != nil {
		return fmt.Errorf("failed to read code: %s", err)
	}

	if err := c.login(); err != nil {
		return fmt.Errorf("failed to login: %s", err)
	}

	if err := c.accessSubmitPage(contest); err != nil {
		return err
	}

	if err := c.chooseLanguage(language); err != nil {
		return err
	}

	if err := c.chooseProblem(problem); err != nil {
		return err
	}

	if err := c.typeCode(code, file); err != nil {
		return fmt.Errorf("failed to submit: %s", err)
	}

	if err := c.screenshot(`#main-div`, "before_submit.png"); err != nil {
		return err
	}

	if err := c.submit(); err != nil {
		return fmt.Errorf("failed to submit: %s", err)
	}

	if err := c.screenshot(`#main-div`, "after_submit.png"); err != nil {
		return err
	}

	return nil
}

func readCode(file string) (string, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (c *Client) login() error {
	loginURL := c.baseURL + "/login"
	_, _ = fmt.Fprintf(c.outStream, "logging in with username %s: %s\n", c.username, loginURL)

	return chromedp.Run(
		c.ctx,
		chromedp.Navigate(loginURL),
		chromedp.SendKeys(`#username`, c.username),
		chromedp.SendKeys(`#password`, c.password),
		chromedp.Click(`submit`),
		chromedp.WaitVisible(`.alert-success`),
	)
}

func (c *Client) screenshot(selector, filepath string) error {
	var buf []byte
	err := chromedp.Run(
		c.ctx,
		chromedp.Screenshot(selector, &buf, chromedp.NodeVisible),
	)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath, buf, 0644)
}

func (c *Client) accessSubmitPage(contest string) error {
	submitURL := fmt.Sprintf("%s/contests/%s/submit", c.baseURL, contest)
	_, _ = fmt.Fprintf(c.outStream, "accessing to submit page: %s\n", submitURL)

	return chromedp.Run(
		c.ctx,
		chromedp.Navigate(submitURL),
		chromedp.WaitVisible(`#main-div`),
	)
}

func (c *Client) chooseProblem(problem string) error {
	_, _ = fmt.Fprintf(c.outStream, "choosing problem: %s\n", problem)

	return chromedp.Run(
		c.ctx,
		chromedp.WaitVisible(`//span[@id="select2-select-task-container"]`),
		chromedp.Click(`//span[@id="select2-select-task-container"]`),
		chromedp.Click(fmt.Sprintf(`//ul[@id="select2-select-task-results"]/li[starts-with(text(), "%s - ")]`, strings.ToUpper(problem))),
	)
}

func (c *Client) chooseLanguage(language string) error {
	_, _ = fmt.Fprintf(c.outStream, "choosing language: %s\n", language)

	return chromedp.Run(
		c.ctx,
		chromedp.WaitVisible(`//span[starts-with(@id, "select2-dataLanguageId")]`),
		chromedp.Click(`//span[starts-with(@id, "select2-dataLanguageId")]`),
		chromedp.KeyAction(language+"\n"),
	)
}

func (c *Client) typeCode(code, file string) error {
	_, _ = fmt.Fprintf(c.outStream, "typing in source code: %s\n", file)

	return chromedp.Run(
		c.ctx,
		chromedp.WaitVisible(`div.CodeMirror-scroll`),
		chromedp.Click(`div.CodeMirror-scroll`),
		chromedp.KeyAction(code),
	)
}

func (c *Client) submit() error {
	_, _ = fmt.Fprintf(c.outStream, "submit!\n")

	return chromedp.Run(
		c.ctx,
		chromedp.WaitVisible(`#submit`),
		chromedp.Click(`#submit`),
	)
}
