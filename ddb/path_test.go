package ddb_test

import (
	"testing"

	messagev1 "github.com/crewlinker/protoc-gen-dynamodb/proto/example/message/v1"
)

func BenchmarkPathBuilding(b *testing.B) {
	b.ReportAllocs()
	k1 := messagev1.Kitchen{
		ExtraKitchen: &messagev1.Kitchen{
			ExtraKitchen: &messagev1.Kitchen{
				ApplianceEngines: []*messagev1.Engine{{Brand: "my-engine"}},
			}}}

	for n := 0; n < b.N; n++ {
		p1 := k1.DynamoPath().ExtraKitchen().ExtraKitchen().ApplianceEngines().Index(0).Brand().String()
		if p1 != ".16.16.19[0].1" {
			b.Fatalf("not valid path: %s", p1)
		}
	}
}
