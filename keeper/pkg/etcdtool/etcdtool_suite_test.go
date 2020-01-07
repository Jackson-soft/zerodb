package etcdtool_test

import (
	"io/ioutil"
	"testing"

	"git.2dfire.net/zerodb/common/config"
	"git.2dfire.net/zerodb/common/utils"
	"git.2dfire.net/zerodb/keeper/pkg/etcdtool"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	myName   = "cluster_test"
	myConfig = new(config.Config)

	store                *etcdtool.Store
	fileData, configData []byte
)

func TestEtcdtool(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Etcdtool Suite")
}

var _ = BeforeSuite(func() {
	endpoint := []string{"10.1.22.0:2379"}

	var err error
	store, err = etcdtool.NewStore(endpoint)
	Expect(err).NotTo(HaveOccurred())

	fileData, err = ioutil.ReadFile("../../../proxy/proxy/test-conf/proxy_conf_test.yaml")
	Expect(err).NotTo(HaveOccurred())

	err = utils.LoadYaml(fileData, myConfig)
	Expect(err).NotTo(HaveOccurred())

	configData, err = utils.UnLoadYaml(myConfig)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	err := store.Close()
	Expect(err).NotTo(HaveOccurred())
})
