package disk

import "testing"

func TestInferFormatFromVolume(t *testing.T) {
	t.Parallel()

	cases := []struct {
		in   string
		want string
		ok   bool
	}{
		{"vm-109-disk-0", "raw", true},                          // LVM-thin / ZFS: no extension
		{"vm-100-disk-1", "raw", true},                          // EFI disk on block storage
		{"109/vm-109-disk-0.qcow2", "qcow2", true},              // dir storage, qcow2
		{"109/vm-109-disk-0.raw", "raw", true},                  // dir storage, raw
		{"109/vm-109-disk-0.vmdk", "vmdk", true},                // dir storage, vmdk
		{"109/vm-109-disk-0.QCOW2", "qcow2", true},              // case-insensitive
		{"template/noble-server-cloudimg-amd64.img", "", false}, // unknown extension: ask the API
	}

	for _, c := range cases {
		got, ok := InferFormatFromVolume(c.in)
		if got != c.want || ok != c.ok {
			t.Errorf("InferFormatFromVolume(%q) = %q,%v; want %q,%v", c.in, got, ok, c.want, c.ok)
		}
	}
}
