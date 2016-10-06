package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/apprenda/kismatic-platform/pkg/inspector/rule"
	"github.com/spf13/cobra"
)

// NewCmdRules returns the "rules" command
func NewCmdRules(out io.Writer) *cobra.Command {
	var file string
	cmd := &cobra.Command{
		Use:   "rules",
		Short: "Manipulate the inspector's rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "inspector-rules.yaml", "file where inspector rules are to be written")
	cmd.AddCommand(NewCmdDumpRules(out, file))
	cmd.AddCommand(NewCmdValidateRules(out, file))
	return cmd
}

// NewCmdDumpRules returns the "dump" command
func NewCmdDumpRules(out io.Writer, file string) *cobra.Command {
	var overwrite bool
	cmd := &cobra.Command{
		Use:   "dump",
		Short: "Dump the inspector rules to a file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := os.Stat(file); err == nil && !overwrite {
				return fmt.Errorf("%q already exists. Use --overwrite to overwrite it", file)
			}
			f, err := os.Create(file)
			if err != nil {
				return fmt.Errorf("error creating %q: %v", file, err)
			}
			if err := rule.DumpDefaultRules(f); err != nil {
				return fmt.Errorf("error dumping rules: %v", err)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite the destination file if it exists")
	return cmd
}

func NewCmdValidateRules(out io.Writer, file string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate the inspector rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := os.Stat(file); os.IsNotExist(err) {
				return fmt.Errorf("%q does not exist", file)
			}
			rules, err := rule.ReadFromFile(file)
			if err != nil {
				return err
			}
			if !validateRules(out, rules) {
				return fmt.Errorf("invalid rules found in %q", file)
			}
			fmt.Fprintf(out, "Rules are valid\n")
			return nil
		},
	}
	return cmd
}

// validates rules, printing error messages to the console
func validateRules(out io.Writer, rules []rule.Rule) bool {
	allOK := true
	for i, r := range rules {
		errs := r.Validate()
		if len(errs) > 0 {
			allOK = false
			fmt.Fprintf(out, "%s (Rule #%d):\n", r.GetRuleMeta().Kind, i+1)
			for _, e := range errs {
				fmt.Fprintf(out, "- %v\n", e)
			}
			fmt.Fprintln(out, "")
		}
	}
	return allOK
}
