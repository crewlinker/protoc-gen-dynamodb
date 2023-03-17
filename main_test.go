package main_test

import (
	"bytes"
	"context"
	"io"
	"os/exec"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestProtocGenDynamodb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "main")
}

var _ = Describe("e2e", func() {

	DescribeTable("generator errors", func(ctx context.Context, name, errExp string) {
		errb := bytes.NewBuffer(nil)
		cmd := exec.CommandContext(ctx, "buf", "generate", "--path", filepath.Join("example", "wrong", "v1", name))
		cmd.Stderr = io.MultiWriter(GinkgoWriter, errb)
		Expect(cmd.Run()).ToNot(Succeed())
		Expect(errb.String()).To(MatchRegexp(errExp))
	},
		Entry("sort key only", "sort_key_only.proto", `has a sort key, but not a partition key`),
		Entry("field is both pk and sk", "pk_sk_same_field.proto", `both marked as PK and as SK`),
		Entry("multiple fields as pk", "multiple_fields_pk.proto", `field 'One' is already marked as PK`),
		Entry("multiple fields as sk", "multiple_fields_sk.proto", `field 'One' is already marked as SK`),
		Entry("invalid type for pk", "pk_invalid_type.proto", `field 'Pk' must be a basic type that marshals to Number,String or Bytes to be a PK`),
		Entry("invalid type for sk", "sk_invalid_type.proto", `field 'Sk' must be a basic type that marshals to Number,String or Bytes to be a SK`),
	)
})
