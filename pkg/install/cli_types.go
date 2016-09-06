package install

// CliOpts cli options
type CliOpts struct {
	PlanFilename             string
	CaCsr                    string
	CaConfigFile             string
	CaSigningProfile         string
	CertsDestination         string
	SkipCAGeneration         bool
	RestartEtcdService       bool
	RestartKubernetesService bool
	RestartCalicoService     bool
	RestartDockerService     bool
}
