package proto_to_entity

import (
	"os"
	"strings"
	"testing"
)

func TestGenerateFileDescriptorSetProtobuf(t *testing.T) {
	fdp, err := generateFileDescriptorProto("ping.proto")
	if err != nil {
		t.Error(err)
	}

	if fdp == nil {
		t.Error("FileDescriptorProto is nil")
	}

	err = os.Remove("ping.protobuf")
	if err != nil {
		t.Error(err)
	}

}

func TestGenerateStructFromMessage(t *testing.T) {
	fdp, _ := generateFileDescriptorProto("ping.proto")

	w := &strings.Builder{}

	err := generateStructFromMessage(fdp, w)
	if err != nil {
		t.Error(err)
	}

	want := `type Request struct {
}
type Response struct {
Version string
Status Status
Time *time.Time
}
`

	if w.String() != want {
		t.Errorf("\nwant %s\ngot  %s\n", want, w.String())
		//t.Errorf("\nwant %x\ngot  %x\n", want, w.String())
	}

}

func TestGenerateMapper(t *testing.T) {
	fdp, _ := generateFileDescriptorProto("ping.proto")
	response := fdp.GetMessageType()[1]

	w := &strings.Builder{}
	err := generateMapToEntity(response, w)
	if err != nil {
		t.Error(err)
	}

	t.Error(w.String())
}
