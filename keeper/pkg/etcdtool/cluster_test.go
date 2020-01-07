package etcdtool_test

import (
	"fmt"
	"testing"

	"git.2dfire.net/zerodb/keeper/pkg/etcdtool"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cluster", func() {
	Context("cluster", func() {
		It("get", func() {
			data, err := store.GetClusters(myName)
			Expect(err).NotTo(HaveOccurred())
			Expect(data).ShouldNot(BeEmpty())
			fmt.Fprintln(GinkgoWriter, data)
		})
	})
})

func TestMy(t *testing.T) {
	vv := []etcdtool.ClusterInfo{{"ddd", 7, 5}, {"xxx", 4, 5}}
	for k, v := range vv {
		if v.Host == "xxx" && v.Port == 5 {
			vv = append(vv[:k], vv[k+1:]...)
			break
		}
	}
	t.Log(vv)
}
