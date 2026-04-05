//go:build !darwin && !linux && !freebsd && !solaris && !aix && !windows
// +build !darwin,!linux,!freebsd,!solaris,!aix,!windows

package gamma

func supportsHardTabs(uint64) bool {
	return false
}
