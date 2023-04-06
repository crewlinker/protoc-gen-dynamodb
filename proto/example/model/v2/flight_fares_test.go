package modelv2_test

import (
	modelv2 "github.com/crewlinker/protoc-gen-dynamodb/proto/example/model/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("flight fares", func() {

	It("should marshal entities", func() {
		ff1 := &modelv2.FlightFares{}
		ff1.FromDynamoEntity(&modelv2.FlightFares_Fare{&modelv2.Fare{}})
		fm1, err := ff1.MarshalDynamoItem()
		Expect(err).ToNot(HaveOccurred())
		_ = fm1
	})

})
