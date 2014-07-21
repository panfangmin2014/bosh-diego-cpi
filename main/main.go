package main

import (
	"flag"
	"os"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshcmd "github.com/cloudfoundry/bosh-agent/platform/commands"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
	boshuuid "github.com/cloudfoundry/bosh-agent/uuid"

	bdcaction "github.com/cppforlife/bosh-diego-cpi/action"
	bdcdisp "github.com/cppforlife/bosh-diego-cpi/api/dispatcher"
	bdctrans "github.com/cppforlife/bosh-diego-cpi/api/transport"
	bdcbbs "github.com/cppforlife/bosh-diego-cpi/bbs"
)

const mainLogTag = "main"

var (
	configPathOpt = flag.String("configPath", "", "Path to configuration file")
)

func main() {
	logger, fs, cmdRunner, uuidGen := basicDeps()

	defer logger.HandlePanic("Main")

	flag.Parse()

	config, err := NewConfigFromPath(*configPathOpt, fs)
	if err != nil {
		logger.Error(mainLogTag, "Loading config: %s", err)
		os.Exit(1)
	}

	dispatcher, err := buildDispatcher(config, logger, fs, cmdRunner, uuidGen)
	if err != nil {
		logger.Error(mainLogTag, "Building dispatcher: %s", err)
		os.Exit(1)
	}

	cli := bdctrans.NewCLI(os.Stdin, os.Stdout, dispatcher, logger)

	err = cli.ServeOnce()
	if err != nil {
		logger.Error(mainLogTag, "Serving once: %s", err)
		os.Exit(1)
	}
}

func basicDeps() (boshlog.Logger, boshsys.FileSystem, boshsys.CmdRunner, boshuuid.Generator) {
	logger := boshlog.NewWriterLogger(boshlog.LevelDebug, os.Stderr, os.Stderr)

	fs := boshsys.NewOsFileSystem(logger)

	cmdRunner := boshsys.NewExecCmdRunner(logger)

	uuidGen := boshuuid.NewGenerator()

	return logger, fs, cmdRunner, uuidGen
}

func buildDispatcher(
	config Config,
	logger boshlog.Logger,
	fs boshsys.FileSystem,
	cmdRunner boshsys.CmdRunner,
	uuidGen boshuuid.Generator,
) (bdcdisp.Dispatcher, error) {
	store, err := bdcbbs.NewBBS(config.ETCD.ConnectAddresses)
	if err != nil {
		return nil, err
	}

	compressor := boshcmd.NewTarballCompressor(cmdRunner, fs)

	actionFactory := bdcaction.NewConcreteFactory(
		store,
		fs,
		uuidGen,
		compressor,
		config.Actions,
		logger,
	)

	caller := bdcdisp.NewJSONCaller()

	return bdcdisp.NewJSON(actionFactory, caller, logger), nil
}
