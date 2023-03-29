package ddbtx_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDdbtx(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ddb/ddbtx")
}
