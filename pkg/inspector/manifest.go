package inspector

// BuildManifest returns a manifest that is specific
// for the provided Linux distribution
func BuildManifest(d Distro) (*Manifest, error) {
	// Special case for testing on darwin
	switch d {
	case Darwin:
		return &Manifest{
			BinaryDependencies: []BinaryDependencyCheck{
				BinaryDependencyCheck{
					BinaryName: "ls",
				},
			},
		}, nil
	}
	return &Manifest{}, nil
}

// The Manifest lists all the checks that will be performed
// by the inspector.
type Manifest struct {
	BinaryDependencies           []BinaryDependencyCheck
	InstalledPackageDependencies []PackageInstalledCheck
	AvailablePackageDependencies []PackageAvailableCheck
	OpenTCPPorts                 []int
}
