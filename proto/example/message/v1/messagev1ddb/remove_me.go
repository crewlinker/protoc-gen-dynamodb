package messagev1ddb

import "github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"

func (p Kitchen) TryBasic() expression.NameBuilder {
	return p.AppendName(expression.Name("1"))
}

func (p Kitchen) TryEngine() Engine {
	return Engine{p.AppendName(expression.Name("2"))}
}
