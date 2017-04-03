package cli

import (
	"fmt"

	"github.com/spf13/pflag"
)

func addPlanFileFlag(flagSet *pflag.FlagSet, p *string) {
	flagSet.StringVarP(p, "plan-file", "f", "kismatic-cluster.yaml", "path to the installation plan file")
}

type planFileNotFoundErr struct {
	filename string
}

func (e planFileNotFoundErr) Error() string {
	return fmt.Sprintf("Plan file not found at %q. If you don't have a plan file, you may generate one with 'kismatic install plan'", e.filename)
}
