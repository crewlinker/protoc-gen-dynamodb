package ddbpath

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLexer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ddb/ddbpath")
}

var _ = DescribeTable("parse path", func(s string, expParts []PathElement, expErr error) {
	parts, err := ParsePath(s)
	if expErr != nil {
		Expect(err).To(Equal(err))
	} else {
		Expect(err).To(BeNil())
	}

	Expect(parts).To(Equal(expParts))
	// @TODO assert error
},
	Entry("1", ".foo.bar.dar", []PathElement{{"foo", -1}, {"bar", -1}, {"dar", -1}}, nil),
	Entry("2", "[1][2][3]", []PathElement{{"1", 1}, {"2", 2}, {"3", 3}}, nil),
)

var _ = DescribeTable("select values", func(av types.AttributeValue, paths []string, expVals map[string]types.AttributeValue, expErr error) {
	vals, err := SelectValues(av, paths...)
	if expErr != nil {
		Expect(err).To(Equal(err))
	} else {
		Expect(err).To(BeNil())
	}

	Expect(vals).To(Equal(expVals))
},
	Entry("1",
		&types.AttributeValueMemberM{},
		[]string{".foo"},
		map[string]types.AttributeValue{},
		nil,
	),
	Entry("2",
		&types.AttributeValueMemberM{Value: map[string]types.AttributeValue{"foo": &types.AttributeValueMemberN{Value: "100"}}},
		[]string{".foo"},
		map[string]types.AttributeValue{
			".foo": &types.AttributeValueMemberN{Value: "100"},
		},
		nil,
	),
)

var res []PathElement

func BenchmarkPathSelect(b *testing.B) {
	r := make([]PathElement, 100)
	p := []string{
		".bar.dar.sar.x.foo.d",
		"[100][5].dar[4].foo[0]",
		".v.bar",
		".bar.dar.sar.x.1.d",
	}

	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		for i := 0; i < len(p); i++ {
			v, err := AppendParsePath(p[i], r)
			if err != nil {
				b.Fatalf("append failed: %v", err)
			}
			res = v
			r = r[:0]
		}
	}
}

var v1 map[string]types.AttributeValue

func BenchmarkSelectValues(b *testing.B) {
	inner := &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
		"xin": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
			"dir": &types.AttributeValueMemberL{Value: []types.AttributeValue{
				&types.AttributeValueMemberN{Value: "100"},
				&types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
					"deep": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
						"deep": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
							"final": &types.AttributeValueMemberS{Value: "foo"},
						}},
					}},
				}},
			}},
		}},
	}}

	outer := &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
		"bar": inner,
		"foo": inner,
	}}
	p := []string{
		".bar.xin.dir[0]",
		".foo.xin.dir[0]",
		".foo",
		".bar.xin.dir[1].deep.deep.final",
	}

	b.ResetTimer()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		r, err := SelectValues(outer, p...)
		if err != nil || len(r) != len(p) {
			b.Fatalf("failed to select: %v/%v", r, err)
		}
		v1 = r
	}
}
