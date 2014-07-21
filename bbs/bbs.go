package bbs

import (
	runbbs "github.com/cloudfoundry-incubator/runtime-schema/bbs"
	runmodels "github.com/cloudfoundry-incubator/runtime-schema/models"
	gunktp "github.com/cloudfoundry/gunk/timeprovider"
	stor "github.com/cloudfoundry/storeadapter"
	storetcd "github.com/cloudfoundry/storeadapter/etcdstoreadapter"
	storwpool "github.com/cloudfoundry/storeadapter/workerpool"
	lager "github.com/pivotal-golang/lager"
)

var LRPNotFound = stor.ErrorKeyNotFound

type BBS interface {
	DesireLRP(runmodels.DesiredLRP) error
	GetDesiredLRPByProcessGuid(string) (runmodels.DesiredLRP, error)
	RemoveDesiredLRPByProcessGuid(string) error

	GetActualLRPsByProcessGuid(string) ([]runmodels.ActualLRP, error)
}

func NewBBS(etcdConnectAddresses []string) (BBS, error) {
	workerPool := storwpool.NewWorkerPool(10)
	adapter := storetcd.NewETCDStoreAdapter(etcdConnectAddresses, workerPool)

	err := adapter.Connect()
	if err != nil {
		return nil, err
	}

	timeProvider := gunktp.NewTimeProvider()
	bbsLogger := lager.NewLogger("bbs")

	return runbbs.NewBBS(adapter, timeProvider, bbsLogger), nil
}
