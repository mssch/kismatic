package cli

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/apprenda/kismatic/pkg/install"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

const dashboardAdminKubeconfigFilename = "dashboard-admin-kubeconfig"

type dashboardOpts struct {
	generatedAssetsDir string
	planFilename       string
}

const url = "http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/#!/login"

// NewCmdDashboard opens or displays the dashboard URL
func NewCmdDashboard(in io.Reader, out io.Writer) *cobra.Command {
	opts := dashboardOpts{}

	cmd := &cobra.Command{
		Use:   "dashboard",
		Short: "Opens the kubernetes dashboard URL of the cluster",
		Long: `
  This command is a convenience command to open the dashbord.
  - Retrieves the token of the secret for ServiceAccount 'kubernetes-dashboard-admin':
  -----------------------------------------------------------------------------------------------------------------------------------------------------
  export SECRET="$(./kubectl get sa kubernetes-dashboard-admin -o 'jsonpath={.secrets[0].name}' -n kube-system --kubeconfig generated/kubeconfig)"
  ./kubectl describe secrets $SECRET -n kube-system --kubeconfig generated/kubeconfig | awk '$1=="token:"{print $2}'
  -----------------------------------------------------------------------------------------------------------------------------------------------------

  - Runs 'kubectl proxy':
  -----------------------------------------------------------------------------------------------------------------------------------------------------
  kubectl proxy --kubeconfig generated/kubeconfig
  -----------------------------------------------------------------------------------------------------------------------------------------------------`,
		Example: `  -----------------------------------------------------------------------------------------------------------------------------------------------------
  ./kismatic dashboard
  Opening kubernetes dashboard in default browser...
  Use the kubeconfig in "generated/dashboard-admin-kubeconfig"
  Starting to serve on 127.0.0.1:8001
  -----------------------------------------------------------------------------------------------------------------------------------------------------
  ./kismatic dashboard url
  http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/#!/login
  -----------------------------------------------------------------------------------------------------------------------------------------------------
  ./kismatic dashboard kubeconfig
  Generated kubeconfig in "generated/dashboard-admin-kubeconfig"
  -----------------------------------------------------------------------------------------------------------------------------------------------------
  ./kismatic dashboard token
  eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrdWJlcm5ldGVzLWRhc2hib2FyZC1hZG1pbi10b2tlbi1kd2o1eiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJrdWJlcm5ldGVzLWRhc2hib2FyZC1hZG1pbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImNjZGQyYmViLTZmMmUtMTFlOC1hNjM5LTBhMWQ1NzNmODZiOCIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTprdWJlcm5ldGVzLWRhc2hib2FyZC1hZG1pbiJ9.VxrM2p3lJHFK9U7Zg6wTjHC-bGKOvkR31_8KIbveQyUjlP7xnlwKm6PKgS-mjyGyZEvBjEIO3xsY-8YEW-1n091fNnPCkBmJLxlljhfpPu9JzgsG1KOe6Ha1-aOHO4PsH8fZYVylOdIP13zo9v5kgmpE7j5YecKY-6aWzyB8tverNJoMN8kvCUsrzVcfV3uOGBsdcn1aDtSSyiKfb5UdIKVkB-4i9VDR3xgAmDP1hTM50aXT1chpt69E-4Cl4qBwYR4mj47V1aTh0oK10Qv6XHd4zydHahlbSiM7LHMjTVekEIooDHoQuqIe9vnzVoPHp-PRWrNetRCSfNJfsRBD7Q
  -----------------------------------------------------------------------------------------------------------------------------------------------------`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("Unexpected args: %v", args)
			}
			return doDashboard(in, out, opts)
		},
	}

	cmd.PersistentFlags().StringVar(&opts.generatedAssetsDir, "generated-assets-dir", "generated", "path to the directory where assets generated during the installation process will be stored")
	addPlanFileFlag(cmd.PersistentFlags(), &opts.planFilename)

	cmd.AddCommand(NewCmdDashboardURL(out))
	cmd.AddCommand(NewCmdDashboardToken(out, opts))
	cmd.AddCommand(NewCmdDashboardKubeconfig(out, opts))

	return cmd
}

func NewCmdDashboardURL(out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "url",
		Short: "Display the kubernetes dashboard URL",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("Unexpected args: %v", args)
			}
			fmt.Fprintln(out, url)
			return nil
		},
	}
}

func NewCmdDashboardToken(out io.Writer, opts dashboardOpts) *cobra.Command {
	return &cobra.Command{
		Use:   "token",
		Short: "Print the ServiceAccount 'kubernetes-dashboard-admin' token",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("Unexpected args: %v", args)
			}
			kubeconfig := filepath.Join(opts.generatedAssetsDir, "kubeconfig")
			if stat, err := os.Stat(kubeconfig); os.IsNotExist(err) || stat.IsDir() {
				return fmt.Errorf("Did not find required kubeconfig file %q", kubeconfig)
			}
			token, err := getToken(opts.generatedAssetsDir)
			if err != nil {
				return fmt.Errorf("Error retrieving the token: %v", err)
			}
			fmt.Fprintln(out, token)
			return nil
		},
	}
}

func NewCmdDashboardKubeconfig(out io.Writer, opts dashboardOpts) *cobra.Command {
	return &cobra.Command{
		Use:   "kubeconfig",
		Short: "Generate a kubeconfig file with the ServiceAccount 'kubernetes-dashboard-admin' token",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("Unexpected args: %v", args)
			}
			kubeconfig := filepath.Join(opts.generatedAssetsDir, "kubeconfig")
			if stat, err := os.Stat(kubeconfig); os.IsNotExist(err) || stat.IsDir() {
				return fmt.Errorf("Did not find required kubeconfig file %q", kubeconfig)
			}
			adminKubeconfig := filepath.Join(opts.generatedAssetsDir, dashboardAdminKubeconfigFilename)
			// Generate dashboard admin certificate if it does not exist
			if _, err := os.Stat(adminKubeconfig); os.IsNotExist(err) {
				if err := generateKubeconfig(opts.planFilename, opts.generatedAssetsDir, adminKubeconfig); err != nil {
					return err
				}
				fmt.Fprintf(out, "Generated kubeconfig in %q\n", adminKubeconfig)
			} else {
				fmt.Fprintf(out, "Found kubeconfig in %q\n", adminKubeconfig)
			}
			return nil
		},
	}
}

func doDashboard(in io.Reader, out io.Writer, opts dashboardOpts) error {
	kubeconfig := filepath.Join(opts.generatedAssetsDir, "kubeconfig")
	if stat, err := os.Stat(kubeconfig); os.IsNotExist(err) || stat.IsDir() {
		return fmt.Errorf("Did not find required kubeconfig file %q", kubeconfig)
	}
	var generateErr error
	adminKubeconfig := filepath.Join(opts.generatedAssetsDir, dashboardAdminKubeconfigFilename)
	// Generate dashboard admin certificate if it does not exist
	if _, err := os.Stat(adminKubeconfig); os.IsNotExist(err) {
		generateErr = generateKubeconfig(opts.planFilename, opts.generatedAssetsDir, adminKubeconfig)
	}

	if generateErr != nil {
		fmt.Fprintf(out, "Error generating a kubeconfig file, you may still use the dashboard with your own ServiceAccount token\n\n")
	}

	fmt.Fprintf(out, "Opening kubernetes dashboard in default browser...\n")
	if generateErr == nil {
		fmt.Fprintf(out, "Use the kubeconfig in %q\n", adminKubeconfig)
	}
	if err := browser.OpenURL(url); err != nil {
		fmt.Fprintf(out, "Unexpected error opening the kubernetes dashboard: %v. You may access it at %q", err, url)
	}

	cmd := exec.Command("./kubectl", "proxy", "--kubeconfig", filepath.Join(opts.generatedAssetsDir, "kubeconfig"))
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = out
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Error running kubectl proxy: %v", err)
	}

	return nil
}

func generateKubeconfig(planeFile, generatedAssetsDir, outFile string) error {
	planner := &install.FilePlanner{File: planeFile}
	plan, err := planner.Read()
	if err != nil {
		return fmt.Errorf("Error reading plan file: %v", err)
	}
	kubeconfig := filepath.Join(generatedAssetsDir, "kubeconfig")
	if stat, err := os.Stat(kubeconfig); os.IsNotExist(err) || stat.IsDir() {
		return fmt.Errorf("Did not find required kubeconfig file %q", kubeconfig)
	}
	token, err := getToken(generatedAssetsDir)
	if err != nil {
		return fmt.Errorf("Error retrieving the token: %v", err)
	}
	if err := install.GenerateDashboardAdminKubeconfig(token, plan, generatedAssetsDir, outFile); err != nil {
		return fmt.Errorf("Error generating the kubeconfig file: %v", err)
	}
	return nil
}

func getToken(generatedAssetsDir string) (string, error) {
	// All of this is required because cannot set a label on the secret so no selectors
	cmd := exec.Command("./kubectl", "-n", "kube-system", "get", "sa", "kubernetes-dashboard-admin", "-o", "jsonpath={.secrets[0].name}", "--kubeconfig", filepath.Join(generatedAssetsDir, "kubeconfig"))
	sa, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting token secret: %v", err)
	}
	if len(sa) == 0 || !strings.Contains(string(sa), "kubernetes-dashboard-admin-token") {
		return "", fmt.Errorf("kubernetes-dashboard-admin-token secret not found")
	}

	cmd = exec.Command("./kubectl", "-n", "kube-system", "get", "secrets", string(sa), "-o", "jsonpath={.data.token}", "--kubeconfig", filepath.Join(generatedAssetsDir, "kubeconfig"))
	tokenBytes, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error getting the token: %v", err)
	}
	if len(tokenBytes) == 0 {
		return "", fmt.Errorf("got an empty token")
	}
	decodedBytes, err := base64.StdEncoding.DecodeString(string(tokenBytes))
	if err != nil {
		return "", fmt.Errorf("error decoding token: %v", err)
	}
	return string(decodedBytes), nil
}
