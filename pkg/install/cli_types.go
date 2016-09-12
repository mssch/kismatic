package install

// CliOpts cli options
type CliOpts struct {
	PlanFilename     string
	CaCsr            string
	CaConfigFile     string
	CaSigningProfile string
	CertsDestination string
	RestartServices  bool
	Verbose          bool
}
