// +build linux

package libcontainer

import (
	"testing"

	"github.com/opencontainers/runc/libcontainer/configs"
)

func TestCheckMountDestOnProc(t *testing.T) {
	dest := "/rootfs/proc/"
	err := checkMountDestination("/rootfs", dest)
	if err == nil {
		t.Fatal("destination inside proc should return an error")
	}
}

func TestCheckMountDestInSys(t *testing.T) {
	dest := "/rootfs//sys/fs/cgroup"
	err := checkMountDestination("/rootfs", dest)
	if err != nil {
		t.Fatal("destination inside /sys should not return an error")
	}
}

func TestCheckMountDestFalsePositive(t *testing.T) {
	dest := "/rootfs/sysfiles/fs/cgroup"
	err := checkMountDestination("/rootfs", dest)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNeedsSetupDev(t *testing.T) {
	config := &configs.Config{
		Mounts: []*configs.Mount{
			{
				Device:      "bind",
				Source:      "/dev",
				Destination: "/dev",
			},
		},
	}
	if needsSetupDev(config) {
		t.Fatal("expected needsSetupDev to be false, got true")
	}
}

func TestNeedsSetupDevStrangeSource(t *testing.T) {
	config := &configs.Config{
		Mounts: []*configs.Mount{
			{
				Device:      "bind",
				Source:      "/devx",
				Destination: "/dev",
			},
		},
	}
	if needsSetupDev(config) {
		t.Fatal("expected needsSetupDev to be false, got true")
	}
}

func TestNeedsSetupDevStrangeDest(t *testing.T) {
	config := &configs.Config{
		Mounts: []*configs.Mount{
			{
				Device:      "bind",
				Source:      "/dev",
				Destination: "/devx",
			},
		},
	}
	if !needsSetupDev(config) {
		t.Fatal("expected needsSetupDev to be true, got false")
	}
}

func TestNeedsSetupDevStrangeSourceDest(t *testing.T) {
	config := &configs.Config{
		Mounts: []*configs.Mount{
			{
				Device:      "bind",
				Source:      "/devx",
				Destination: "/devx",
			},
		},
	}
	if !needsSetupDev(config) {
		t.Fatal("expected needsSetupDev to be true, got false")
	}
}

func TestDetectShiftFS(t *testing.T) {
	procFilesystems := `nodev	sysfs
nodev	rootfs
nodev	ramfs
nodev	bdev
nodev	proc
nodev	cpuset
nodev	cgroup
nodev	cgroup2
nodev	tmpfs
nodev	devtmpfs
nodev	debugfs
nodev	tracefs
nodev	securityfs
nodev	sockfs
nodev	bpf
nodev	pipefs
nodev	hugetlbfs
nodev	devpts
	ext3
	ext2
	ext4
	squashfs
	vfat
nodev	ecryptfs
	fuseblk
nodev	fuse
nodev	fusectl
nodev	pstore
nodev	mqueue
	btrfs
nodev	autofs
nodev	zfs
nodev	binfmt_misc
nodev	aufs
	xfs
	jfs
	msdos
	ntfs
	minix
	hfs
	hfsplus
	qnx4
	ufs
`
	if !detectShiftFS(procFilesystems) {
		t.Fatal("didn't detect shiftfs")
	}
}
