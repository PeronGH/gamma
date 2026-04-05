//go:build !darwin && !linux && !aix && !windows
// +build !darwin,!linux,!aix,!windows

package gamma

func supportsBackspace(uint64) bool {
	return false
}
