package proto_to_entity

import (
	"strings"
	"testing"
)

func TestGenerateMapper(t *testing.T) {
	fdp, _ := generateFileDescriptorProto("digitalgrowth.proto")
	response := fdp.GetMessageType()[2]

	w := &strings.Builder{}
	err := generateMapToEntity(response, w, "digitalgrowth")
	if err != nil {
		t.Error(err)
	}

	//t.Error(w.String())
}
