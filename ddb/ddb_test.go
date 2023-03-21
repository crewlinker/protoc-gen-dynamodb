package ddb_test

import (
	"fmt"
	"testing"

	messagev1 "github.com/crewlinker/protoc-gen-dynamodb/proto/example/message/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDdb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ddb")
}

var _ = Describe("expression building", func() {

	It("should build with list of basic types", func() {
		var p messagev1.KitchenP
		fmt.Println("AAAA", p.ApplianceBrands().At(1))

	})

	// It("should build with nested names", func() {
	// 	var p messagev1.KitchenP
	// 	fmt.Println(p.AnotherKitchen().AnotherKitchen().ApplianceEngines().Index(1).V)

	// })

})
