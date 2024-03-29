package transport

import (
	"io"
	"io/ioutil"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"

	bdcdisp "github.com/cppforlife/bosh-diego-cpi/api/dispatcher"
)

const cliLogTag = "CLI"

type CLI struct {
	in         io.Reader
	out        io.Writer
	dispatcher bdcdisp.Dispatcher
	logger     boshlog.Logger
}

func NewCLI(
	in io.Reader,
	out io.Writer,
	dispatcher bdcdisp.Dispatcher,
	logger boshlog.Logger,
) CLI {
	return CLI{
		in:         in,
		out:        out,
		dispatcher: dispatcher,
		logger:     logger,
	}
}

func (t CLI) ServeOnce() error {
	reqBytes, err := ioutil.ReadAll(t.in)
	if err != nil {
		t.logger.Error(cliLogTag, "Failed reading from IN: %s", err)
		return bosherr.WrapError(err, "Reading from IN")
	}

	respBytes := t.dispatcher.Dispatch(reqBytes)

	_, err = t.out.Write(respBytes)
	if err != nil {
		t.logger.Error(cliLogTag, "Failed writing to OUT: %s", err)
		return bosherr.WrapError(err, "Writing to OUT")
	}

	return nil
}
