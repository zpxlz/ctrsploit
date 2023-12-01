package seccomp

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/pkg/version"
	"golang.org/x/sys/unix"
	"strings"
	"syscall"
)

type Status struct {
	Version version.Version
	Enable  bool
}

type VersionRange [2]Status
type VersionRanges []VersionRange

func (v VersionRange) String() (s string) {
	s = fmt.Sprintf("[%s, %s)", v[0].Version, v[1].Version)
	return
}

func (v VersionRanges) String() (s string) {
	var ranges []string
	for _, r := range v {
		ranges = append(ranges, r.String())
	}
	s = strings.Join(ranges, " U ")
	return
}

type Syscall struct {
	Number           int
	KernelMinVersion string
	KernelMaxVersion string
	DockerMinVersion version.Docker
	DockerMaxVersion version.Docker
	DockerChangelog  []Status
}

func (s Syscall) Enabled() bool {
	_, _, errno := syscall.RawSyscall(
		uintptr(s.Number),
		0,
		0,
		0,
	)
	// return errno == unix.EFAULT
	return errno != unix.EPERM
}

func (s Syscall) Range(status bool) (r VersionRanges) {
	for i := 0; i < len(s.DockerChangelog); i++ {
		changelog := s.DockerChangelog[i]
		if changelog.Enable == status {
			var nextChangelog Status
			if i+1 < len(s.DockerChangelog) {
				nextChangelog = s.DockerChangelog[i+1]
			} else {
				nextChangelog = Status{
					Version: version.FurtherDockerVersion,
					Enable:  !changelog.Enable,
				}
			}
			r = append(r, VersionRange{changelog, nextChangelog})
		}
	}
	return
}

var (
	// IOURingSetup
	// kernel enable: https://elixir.bootlin.com/linux/v5.1-rc1/source/include/linux/syscalls.h#L318
	// seccomp disable: https://github.com/moby/moby/pull/46762
	// seccomp enable: https://github.com/moby/moby/pull/39415
	IOURingSetup = Syscall{
		Number:           unix.SYS_IO_URING_SETUP,
		KernelMinVersion: "v5.1-rc1",
		DockerMaxVersion: version.NewDocker("25.0.0-beta.1"),
		DockerChangelog: []Status{
			{
				Version: version.FirstDockerVersion,
				Enable:  false,
			},
			{
				Version: version.NewDocker("20.10.0-beta1"),
				Enable:  true,
			},
			{
				Version: version.NewDocker("25.0.0-beta.1"),
				Enable:  false,
			},
		},
	}
)
