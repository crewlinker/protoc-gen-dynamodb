package ddb_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDdb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ddb")
}

var _ = Describe("mask", func() {

	It("should not panic", func() {

	})

})
