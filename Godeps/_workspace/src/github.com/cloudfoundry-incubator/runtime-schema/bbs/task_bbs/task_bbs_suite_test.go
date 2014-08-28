package task_bbs_test

import (
	"github.com/cloudfoundry-incubator/runtime-schema/models"
	"github.com/cloudfoundry/storeadapter"
	"github.com/cloudfoundry/storeadapter/storerunner/etcdstorerunner"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"

	"testing"
	"time"
)

var etcdRunner *etcdstorerunner.ETCDClusterRunner
var etcdClient storeadapter.StoreAdapter

var dummyActions = []models.ExecutorAction{
	{
		Action: models.RunAction{
			Path: "cat",
			Args: []string{"/tmp/file"},
		},
	},
}

func TestTaskBbs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Task BBS Suite")
}

var _ = BeforeSuite(func() {
	etcdRunner = etcdstorerunner.NewETCDClusterRunner(5001+config.GinkgoConfig.ParallelNode, 1)
	etcdClient = etcdRunner.Adapter()
})

var _ = AfterSuite(func() {
	etcdRunner.Stop()
})

var _ = BeforeEach(func() {
	etcdRunner.Stop()
	etcdRunner.Start()
})

func itRetriesUntilStoreComesBack(action func() error) {
	It("should keep trying until the store comes back", func(done Done) {
		etcdRunner.GoAway()

		runResult := make(chan error)
		go func() {
			err := action()
			runResult <- err
		}()

		time.Sleep(200 * time.Millisecond)

		etcdRunner.ComeBack()

		Ω(<-runResult).ShouldNot(HaveOccurred())

		close(done)
	}, 5)
}
