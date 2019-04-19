package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Global Parameter
var (
	cmdOptions              Command
	AcceptedDownloadProduct = []string{"gpdb", "gpcc", "gpextras"}
	AcceptedInstallProduct  = []string{"gpdb", "gpcc"}
)

// Command line options
type Command struct {
	Product   string
	Version   string
	Username  string
	CCVersion string
	Debug     bool
	Install   bool
	//Stop        bool
	Force   bool
	Standby bool
	listEnv bool
	Always  bool
}

// Sub Command: Download
// When this command is used it goes and download the product from pivnet
var downloadCmd = &cobra.Command{
	Use:     "download",
	Aliases: []string{`d`},
	Short:   "Download the product from pivotal network",
	Long:    "Download sub-command helps to download the products that are greenplum related from pivotal network",
	Example: "For examples refer: https://github.com/pivotal-gss/go-gpdb/tree/master/gpdb#download",
	PreRun: func(cmd *cobra.Command, args []string) {
		// Accept only the options that we care about
		if !Contains(AcceptedDownloadProduct, cmdOptions.Product) {
			Fatalf("Invalid product option specified: %s, Accepted Options: %v", cmdOptions.Product, AcceptedDownloadProduct)
		}
		if IsValueEmpty(cmdOptions.Username) && cmdOptions.Install {
			Fatalf("No Username option supplied, please use the -u option and enter your Pivotal ID, this will be used to name your environment file during install" )
		}
		//if cmdOptions.Username == "-u " {
		//	Fatalf("Usage :- You have to use -u option ")
		//}

	},
	Run: func(cmd *cobra.Command, args []string) {

		if cmdOptions.Username == "sername"{

		Fatalf("Usage :- You have to use -u option ")

		} else
		{
			//Run download to download the binaries
			Download()
		}






	},
}

// All the usage flags of the download command
func downloadFlags() {
	downloadCmd.Flags().StringVarP(&cmdOptions.Product, "product", "p", "gpdb", "What product do you want to download? [OPTIONS: gpdb, gpcc, gpextras]")
	downloadCmd.Flags().StringVarP(&cmdOptions.Version, "version", "v", "", "OPTIONAL: Which GPDB version software do you want to download ?")
	downloadCmd.Flags().StringVarP(&cmdOptions.Username, "username", "u", "", "What is your Pivotal ID, this will be used to name the environmental file for this installation ?")
	downloadCmd.Flags().BoolVar(&cmdOptions.Install, "install", false, "OPTIONAL: Install after downloaded (Only works with \"gpdb\")?")
}

// Sub Command: Install
// When this command is used it goes and install the product that was downloaded from above
var installCmd = &cobra.Command{
	Use:   "install",
	Aliases: []string{`i`},
	Short: "Install the product downloaded from download command",
	Long:  "Install sub-command helps to install the products that was downloaded using the download command",
	Example: "For examples refer: https://github.com/pivotal-gss/go-gpdb/tree/master/gpdb#install",
	PreRun: func(cmd *cobra.Command, args []string) {
		// Accept only the options that we care about
		if !Contains(AcceptedInstallProduct, cmdOptions.Product) {
			Fatalf("Invalid product option specified: %s, Accepted Options: %v", cmdOptions.Product, AcceptedInstallProduct)
		}
		// If gpcc used then check if ccversion is set
		if cmdOptions.Product == "gpcc" && cmdOptions.CCVersion == "" {
			Fatalf("ccversion is not set, with product \"gpcc\" you will need to set ccversion")
		}
		// if product is not gpdb then standby flag should not be set
		if cmdOptions.Product != "gpdb" && cmdOptions.Standby {
			Fatalf("Cannot set standby flag with product flag \"%s\"", cmdOptions.Product)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Install the product that is downloaded
		install()
	},
}

// All the usage flags of the Install command
func installFlags() {
	installCmd.Flags().StringVarP(&cmdOptions.Product, "product", "p", "gpdb", "What product do you want to Install? [OPTIONS: gpdb, gpcc, gpextras]")
	installCmd.Flags().StringVarP(&cmdOptions.Version, "version", "v", "", "OPTIONAL: Which GPDB version software do you want to install?")
	installCmd.Flags().StringVarP(&cmdOptions.Username, "username", "u", "", "What is your PivotalID, this will be used to name your enviromental file?")
	installCmd.MarkFlagRequired("version")
	installCmd.MarkFlagRequired("username")
	installCmd.Flags().StringVarP(&cmdOptions.CCVersion, "ccversion", "c", "", "What is the version of GPCC that you can to install (for only -p gpcc)?")
	installCmd.Flags().BoolVar(&cmdOptions.Standby, "standby", false, "OPTIONAL: Install standby if the standby host is available")
}

// Sub Command: Remove
// When this command is used it goes and remove the product that was installed by this program
var removeCmd = &cobra.Command{
	Use:   "remove",
	Aliases: []string{`r`},
	Short: "Removes the product installed using the install command",
	Long:  "Remove sub-command helps to remove the products that was installed using the install command",
	Example: "For examples refer: https://github.com/pivotal-gss/go-gpdb/tree/master/gpdb#remove",
	Run: func(cmd *cobra.Command, args []string) {
		// Remove the installation
		remove()
	},
}

// All the usage flags of the remove command
func removeFlags() {
	removeCmd.Flags().StringVarP(&cmdOptions.Version, "version", "v", "", "Which GPDB version software do you want to remove?")
	removeCmd.MarkFlagRequired("version")
	removeCmd.Flags().BoolVarP(&cmdOptions.Force, "force", "f",false, "OPTIONAL: If the database start fails, use force to force manual removal")
}

// Sub Command: Environment
// When this command is used it goes and remove the product that was installed by this program
var envCmd = &cobra.Command{
	Use:   "env",
	Aliases: []string{`e`},
	Short: "Show all the environment installed",
	Long:  "Env sub-command helps to show all the products version installed",
	Example: "For examples refer: https://github.com/pivotal-gss/go-gpdb/tree/master/gpdb#env",
	////Adding an if statement here forcing the user to use the --dont-stop option
	//PreRun: func(cmd *cobra.Command, args []string) {
	//
	//	if (cmdOptions.Stop == false) {
	//		Fatalf("YOu need to specify the --dont-stop flag to avoid shutting down other peoples enviroments")
	//	}
	//},
	Run: func(cmd *cobra.Command, args []string) {
		// search the env directory for the environment files
		// and broadcast to the user
		env()
	},
}

// All the usage flags of the env command
func envFlags() {
	envCmd.Flags().StringVarP(&cmdOptions.Version, "version", "v", "", "OPTIONAL: Which GPDB version software do you want to list?")
	envCmd.Flags().BoolVarP(&cmdOptions.listEnv, "list", "l", false, "List all the environment installed")
	//envCmd.Flags().BoolVar(&cmdOptions.Stop, "dont-stop", false, "OPTIONAL: Don't stop other database when setting this environment")
}

// The root CLI.
var rootCmd = &cobra.Command{
	Use:   fmt.Sprintf("%s [command]", programName),
	Short: "Download / install / remove and manage the software of GPDB products",
	Long: "This program helps to download / install / remove and manage the software of GPDB products",
	Version: programVersion,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Before running any command Setup the logger log level
		initLogger(cmdOptions.Debug)
		// Load all the configuration to the memory
		config()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 { // if no argument specified throw the help menu on the screen
			cmd.Help()
		}
	},
}

func init ()  {
	// root command flags
	rootCmd.PersistentFlags().BoolVarP(&cmdOptions.Debug, "debug", "d", false, "Enable verbose or debug logging")

	// Attach the sub command to the root command.
	rootCmd.AddCommand(downloadCmd)
	downloadFlags()
	rootCmd.AddCommand(installCmd)
	installFlags()
	rootCmd.AddCommand(removeCmd)
	removeFlags()
	rootCmd.AddCommand(envCmd)
	envFlags()
}