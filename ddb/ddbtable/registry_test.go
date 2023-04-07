package ddbtable_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDdbtable(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ddb/ddbtable")
}
