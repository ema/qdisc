package qdisc

import (
	"syscall"
	"testing"

	"github.com/mdlayher/netlink"
	"github.com/mdlayher/netlink/nltest"
)

func TestGetAndParseFail(t *testing.T) {
	c := nltest.Dial(func(req netlink.Message) ([]netlink.Message, error) {
		return nltest.Error(int(syscall.ENOENT), req)
	})

	msgs, err := getAndParse(c)
	if msgs != nil {
		t.Fatalf("msgs should be nil, got '%v' instead", msgs)
	}

	if err == nil {
		t.Fatalf("err should not be nil")
	}
}

func TestGetAndParseShort(t *testing.T) {
	msg := netlink.Message{Data: []byte{0xff, 0xff, 0xff, 0xff}}

	c := nltest.Dial(func(req netlink.Message) ([]netlink.Message, error) {
		msg.Header.Sequence = req.Header.Sequence
		return []netlink.Message{msg}, nil
	})

	msgs, err := getAndParse(c)
	if msgs != nil {
		t.Fatalf("msgs should be nil, got '%v' instead", msgs)
	}

	if err.Error() != "Short message, len=4" {
		t.Fatalf("expected short message error, got '%v' instead", err)
	}
}

func TestGetAndParseUnmarshal(t *testing.T) {
	msg := netlink.Message{Data: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}}

	c := nltest.Dial(func(req netlink.Message) ([]netlink.Message, error) {
		msg.Header.Sequence = req.Header.Sequence
		return []netlink.Message{msg}, nil
	})

	msgs, err := getAndParse(c)
	if msgs != nil {
		t.Fatalf("msgs should be nil, got '%v' instead", msgs)
	}

	if err.Error() != "failed to unmarshal attributes: invalid attribute; length too short or too large" {
		t.Fatalf("expected unmarshal failure error, got '%v' instead", err)
	}
}

func TestGetAndParseOK(t *testing.T) {
	d := []byte{0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 128, 255, 255, 255, 255, 2, 0, 0, 0, 7, 0, 1, 0, 102, 113, 0, 0, 84, 0, 2, 0, 8, 0, 1, 0, 16, 39, 0, 0, 8, 0, 2, 0, 100, 0, 0, 0, 8, 0, 3, 0, 212, 11, 0, 0, 8, 0, 4, 0, 36, 59, 0, 0, 8, 0, 5, 0, 1, 0, 0, 0, 8, 0, 7, 0, 255, 255, 255, 255, 8, 0, 9, 0, 64, 156, 0, 0, 8, 0, 10, 0, 255, 3, 0, 0, 8, 0, 11, 0, 142, 12, 1, 0, 8, 0, 8, 0, 10, 0, 0, 0, 132, 0, 7, 0, 84, 0, 4, 0, 125, 97, 0, 0, 0, 0, 0, 0, 86, 82, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 203, 52, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 38, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 114, 61, 37, 209, 53, 184, 251, 255, 255, 7, 0, 0, 255, 7, 0, 0, 0, 0, 0, 0, 227, 183, 1, 0, 20, 0, 1, 0, 139, 42, 111, 15, 0, 0, 0, 0, 159, 76, 30, 0, 0, 0, 0, 0, 24, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 22, 0, 0, 0, 0, 0, 0, 0, 44, 0, 3, 0, 139, 42, 111, 15, 0, 0, 0, 0, 159, 76, 30, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 84, 0, 4, 0, 125, 97, 0, 0, 0, 0, 0, 0, 86, 82, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 203, 52, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 38, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 114, 61, 37, 209, 53, 184, 251, 255, 255, 7, 0, 0, 255, 7, 0, 0, 0, 0, 0, 0, 227, 183, 1, 0}
	msg := netlink.Message{Data: d}

	c := nltest.Dial(func(req netlink.Message) ([]netlink.Message, error) {
		msg.Header.Sequence = req.Header.Sequence
		return []netlink.Message{msg}, nil
	})

	msgs, err := getAndParse(c)
	if err != nil {
		t.Fatalf("err should be nil, got %v instead", err)
	}

	expect := []QdiscInfo{
		{
			IfaceName:   "lo",
			Parent:      0,
			Handle:      2147549184,
			Kind:        "fq",
			Bytes:       258943627,
			Packets:     1985695,
			Drops:       0,
			Requeues:    22,
			Overlimits:  0,
			GcFlows:     24957,
			Throttled:   13515,
			FlowsPlimit: 0,
		},
	}

	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %v instead", len(msgs))
	}

	if msgs[0] != expect[0] {
		t.Fatalf("messages not as expected: %v", msgs)
	}
}
