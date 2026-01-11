package build

// Platform represents the target platform with OS and architecture.
type Platform struct {
	OS   string
	Arch string
}

// String returns the platform as a hyphen-separated string (e.g., "linux-amd64").
func (p *Platform) String() string {
	return p.OS + "-" + p.Arch
}

// GetPlatform returns the current platform based on OS and Arch.
func GetPlatform() *Platform {
	return &Platform{
		OS:   OS(),
		Arch: Arch(),
	}
}
