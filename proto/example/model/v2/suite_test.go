package modelv2_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestV2(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "proto/example/model/v2")
}
