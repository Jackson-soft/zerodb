package etcdtool_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Etcdtool", func() {
	var err error
	Context("put data", func() {
		key := "myTest"
		value := "test_xxx"
		It("put", func() {
			err = store.PutData(key, value)
			Expect(err).NotTo(HaveOccurred())
		})

		It("get", func() {
			data, err := store.GetData(key)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(data)).To(Equal(value))
		})

		It("delete", func() {
			err = store.Delete(key)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
