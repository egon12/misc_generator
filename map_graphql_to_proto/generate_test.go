package map_graphql_to_proto

import (
	"strings"
	"testing"
)

const graphQLContractInput = `

enum SomePrefixOrderStatus {
        Unknown
        Pending
        Success
        Failed
}

input SomePrefixOrderFilter {
        date_from: Date
        date_until: Date
        status: SomePrefixOrderStatus
}

type SomePrefixOrder {
        date: Time!
        date_end: Time
        product: String!
		price: Int!
        status: SomePrefixOrderStatus!
}

type SomePrefixOrderList {
        list: [SomePrefixOrder!]!
		count: Int!
}`

func TestGenerateMapper(t *testing.T) {
	w := &strings.Builder{}
	SetConfig(Config{
		RemovePrefixes: []string{"SomePrefix"},
		ProtoPackage:   "proto",
	})
	err := generateMapper(graphQLContractInput, w)
	if err != nil {
		t.Error(err)
	}
	t.Error(w)
}
