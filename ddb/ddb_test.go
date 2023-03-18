package ddb

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDdb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ddb")
}

var _ = Describe("encoding options", func() {
	It("should apply options", func() {
		eo1 := ApplyEncodingOptions(WithMask("foo"))
		Expect(eo1.mask).To(Equal(map[string][]string{"foo": {"foo"}}))

		eo2 := ApplyEncodingOptions(WithEncodingOptions(eo1))
		Expect(&eo1 == &eo2).To(BeFalse()) // different memory
		Expect(eo1).To(Equal(eo2))         // same data
	})

	It("should remove duplicate mask", func() {
		eo1 := ApplyEncodingOptions(WithMask("foo", "bar.dar", "", "foo"))
		Expect(eo1.mask).To(Equal(map[string][]string{"foo": {"foo"}, "bar.dar": {"bar", "dar"}}))
	})

	It("should mask", func() {
		eo1 := ApplyEncodingOptions(WithMask("sk.dar", "pk", "0", "1"))
		Expect(eo1.IsMasked("3")).To(BeFalse())
		Expect(eo1.IsMasked("sk")).To(BeTrue())
	})

	It("should select with empty mask", func() {
		eo2 := ApplyEncodingOptions(WithMask())
		Expect(eo2.IsMasked("sk")).To(BeTrue())
	})

	Describe("sub mask", func() {
		It("should return empty mask with nothing matching", func() {
			eo1 := ApplyEncodingOptions(WithMask("foo", "bar.dar", "", "foo"))
			eo2 := eo1.SubMask("bogus")
			Expect(eo2.mask).ToNot(Equal(eo1))
			Expect(eo2.mask).To(Equal(map[string][]string(nil)))
		})

		It("should return copy of encoding options", func() {
			eo1 := ApplyEncodingOptions(WithMask("foo", "bar.dar", "", "foo"))
			eo2 := eo1.SubMask("bar")
			Expect(eo2.mask).ToNot(Equal(eo1.mask))
			Expect(eo2.mask).To(Equal(map[string][]string{"dar": {"dar"}}))
		})

		It("should return empty submask", func() {
			eo1 := ApplyEncodingOptions(WithMask("foo", "bar.dar", "", "foo"))
			eo2 := eo1.SubMask("foo")
			Expect(eo2.mask).To(Equal(map[string][]string{}))
		})

		It("should correctly pass submask as option", func() {
			eo1 := ApplyEncodingOptions(WithMask("foo", "bar.dar", "", "foo"))
			eo2 := ApplyEncodingOptions(WithEncodingOptions(eo1.SubMask("bar")))
			Expect(eo1.mask).ToNot(Equal(eo2))
			Expect(eo2.mask).To(Equal(map[string][]string{"dar": {"dar"}}))
		})
	})
})
