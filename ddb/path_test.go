package ddb_test

import (
	"testing"

	messagev1 "github.com/crewlinker/protoc-gen-dynamodb/proto/example/message/v1"
)

var p2 string

func BenchmarkBasicListPathBuilding(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		p2 = (messagev1.KitchenP{}).AnotherKitchen().AnotherKitchen().ApplianceBrands().At(1)
		if p2 != ".2.2.4[1]" {
			b.Fatalf("failed to build: %v", p2)
		}
	}
}

func BenchmarkMessageListPathBuilding(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		p2 = (messagev1.KitchenP{}).AnotherKitchen().AnotherKitchen().ApplianceEngines().At(1).Brand()
		if p2 != ".2.2.3[1].5" {
			b.Fatalf("failed to build: %v", p2)
		}
	}
}
