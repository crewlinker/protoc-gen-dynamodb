package ddbtable_test

import (
	"testing"

	"github.com/crewlinker/protoc-gen-dynamodb/ddb/ddbtable"
	modelv1 "github.com/crewlinker/protoc-gen-dynamodb/proto/example/model/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDdbtable(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ddb/ddbtable")
}

var _ = Describe("registry", func() {
	var reg ddbtable.Registry
	BeforeEach(func() {
		reg = *ddbtable.NewRegistry()
	})

	It("should register", func() {
		Expect(reg.Register(&modelv1.Thread{}, &ddbtable.TablePlacement{})).To(Succeed())
	})

	// @TODO could fuzz random table placements defintions, put them into dynamodb local and check it
	// never errors from dynamodb, only from our own code.

})
