package letme

import (
	"fmt"
	"github.com/spf13/cobra"
	//"github.com/lockedinspace/letme/pkg"
	"os"
	"text/tabwriter"
)

var version = "0.1.0"
var rootCmd = &cobra.Command{
	Use:     "letme",
	Short:   "Obtain AWS credentials from another account",
	Long: `letme will query the DynamoDB table or cache file for the specified account and
load the temporal credentials onto your aws files.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		versionFlag, _ := cmd.Flags().GetBool("version")
		if versionFlag {
			getVersions()
			os.Exit(0)
		}
		fmt.Println("letme: try 'letme --help' or 'letme -h' for more information")
		os.Exit(0)
	},
}
func getVersions() (string) {
	w := tabwriter.NewWriter(os.Stdout, 20, 20, 10, ' ', 0)
	fmt.Fprintln(w, "CURRENT VERSION:")
	fmt.Fprintln(w, "---------------")
	fmt.Fprintln(w, version)
	w.Flush()
	return " "
}
func init() {
	var Version bool
	rootCmd.PersistentFlags().BoolVarP(&Version, "version", "v", false, "list current, development and latest versions for letme")
}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
	
}