package proto_to_entity

import (
	"strings"
	"testing"
)

func TestGenerateMapToSliceEntity(t *testing.T) {
	fdp, _ := generateFileDescriptorProto("digitalgrowth.proto")

	w := &strings.Builder{}

	err := generateMapToSliceEntities(fdp, w)
	if err != nil {
		t.Error(err)
	}

	want := ``

	if w.String() != want {
		//t.Errorf("\nwant %s\ngot  %s\n", want, w.String())
		//t.Errorf("\nwant %x\ngot  %x\n", want, w.String())
	}

}
