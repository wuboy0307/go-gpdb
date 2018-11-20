package main

import (
	"github.com/spf13/cobra"
	"fmt"
)

// Global Parameter
var (
	cmdOptions Command
	AcceptedDownloadProduct = []string{"gpdb", "gpcc", "gpextras"}
	AcceptedInstallProduct = []string{"gpdb", "gpcc"}
)

// Command line options
type Command struct {
	Product 	string
	Version 	string
	CCVersion 	string
	Debug		bool
	Install 	bool
}

// Sub Command: Download
// When this command is used it goes and download the product from pivnet
var downloadCmd = &cobra.Command{
	Use:   "download",
	Aliases: []string{`d`},
	Short: "Download the product from pivotal network",
	Long:  "Download sub-command helps to download the products that are greenplum related from pivotal network",
	Example: fmt.Sprintf(downloadExample(), programName),
	PreRun: func(cmd *cobra.Command, args []string) {
		// Accept only the options that we care about
		if !Contains(AcceptedDownloadProduct, cmdOptions.Product) {
			Fatalf("Invalid product option specified: %s, Accepted Options: %v", cmdOptions.Product, AcceptedDownloadProduct)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Run download to download the binaries
		Download()
	},
}

// All the usage flags of the download command
func downloadFlags() {
	downloadCmd.Flags().StringVarP(&cmdOptions.Product, "product", "p", "gpdb", "What product do you want to download? [OPTIONS: gpdb, gpcc, gpextras]")
	downloadCmd.Flags().StringVarP(&cmdOptions.Version, "version", "v", "", "OPTIONAL: Which GPDB version software do you want to download ?")
	downloadCmd.Flags().BoolVar(&cmdOptions.Install, "install", false, "OPTIONAL: Install after downloaded (Only works with \"gpdb\")?")
}

// download examples
func downloadExample() string {
	return `=> To download product interactively

%[1]s download

=> To download a specific version

%[1]s download -v <GPDB VERSION>

=> To download and install a specific version

%[1]s download -v <GPDB VERSION> --install

=> To download GPCC software in interactive mode.

%[1]s download -p gpcc

=> To download GPCC software of specific version.

%[1]s download -p gpcc -v <GPDB VERSION>

=> To download all GPDB products in interactive mode

%[1]s download -p gpextras

=> To download all products of specific version.

%[1]s download -p gpextras -v <GPDB VERSION>`
}

// Sub Command: Install
// When this command is used it goes and install the product that was downloaded from above
var installCmd = &cobra.Command{
	Use:   "install",
	Aliases: []string{`i`},
	Short: "Install the product downloaded from download command",
	Long:  "Install sub-command helps to install the products that was downloaded using the download command",
	Example: fmt.Sprintf(installExample(), programName),
	PreRun: func(cmd *cobra.Command, args []string) {
		// Accept only the options that we care about
		if !Contains(AcceptedInstallProduct, cmdOptions.Product) {
			Fatalf("Invalid product option specified: %s, Accepted Options: %v", cmdOptions.Product, AcceptedInstallProduct)
		}
		// If gpcc used then check if ccversion is set
		if cmdOptions.Product == "gpcc" && cmdOptions.CCVersion == "" {
			Fatalf("ccversion is not set, with product \"gpcc\" you will need to set ccversion")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Install the product that is downloaded
		install()
	},
}

// All the usage flags of the download command
func installFlags() {
	installCmd.Flags().StringVarP(&cmdOptions.Product, "product", "p", "gpdb", "What product do you want to Install? [OPTIONS: gpdb, gpcc, gpextras]")
	installCmd.Flags().StringVarP(&cmdOptions.Version, "version", "v", "", "OPTIONAL: Which GPDB version software do you want to install?")
	installCmd.MarkFlagRequired("version")
	installCmd.Flags().StringVarP(&cmdOptions.CCVersion, "ccversion", "c", "", "What is the version of GPCC that you can to install (for only -p gpcc)?")
}

// Install example
func installExample() string {
	return `=> To install gpdb

%[1]s install -v <GPDB VERSION>

=> To install gpcc

%[1]s install -p gpcc -v <GPDB VERSION> -c <GPCC VERSION>`
}

// Sub Command: Remove
// When this command is used it goes and remove the product that was installed by this program
var removeCmd = &cobra.Command{
	Use:   "remove",
	Aliases: []string{`r`},
	Short: "Removes the product installed using the install command",
	Long:  "Remove sub-command helps to remove the products that was installed using the install command",
	Example: fmt.Sprintf(removeExample(), programName),
	Run: func(cmd *cobra.Command, args []string) {
		// Search the tile from all the labs
		fmt.Println("will run remove one day")
	},
}

// All the usage flags of the download command
func removeFlags() {
	removeCmd.Flags().StringVarP(&cmdOptions.Version, "version", "v", "", "Which GPDB version software do you want to remove?")
	removeCmd.MarkFlagRequired("version")
}

// Remove example
func removeExample() string {
	return `=> To remove a particular installation.

%[1]s remove -v <GPDB VERSION>`
}

// Sub Command: Environment
// When this command is used it goes and remove the product that was installed by this program
var envCmd = &cobra.Command{
	Use:   "env",
	Aliases: []string{`e`},
	Short: "Show all the environment installed",
	Long:  "Env sub-command helps to show all the products version installed",
	Example: fmt.Sprintf(envExample(), programName),
	Run: func(cmd *cobra.Command, args []string) {
		// search the env directory for the environment files
		// and broadcast to the user
		envListing()
	},
}

// All the usage flags of the download command
func envFlags() {
	envCmd.Flags().StringVarP(&cmdOptions.Version, "version", "v", "", "OPTIONAL: Which GPDB version software do you want to list?")
}

// Env example
func envExample() string {
	return `=> To list all environment that has been installed and choose env in interactive mode.

%[1]s env

=> To start and use a specific installation.

%[1]s env -v <GPDB VERSION>`
}

// The root CLI.
var rootCmd = &cobra.Command{
	Use:   fmt.Sprintf("%s [command]", programName),
	Short: "Download / install / remove and manage the software of GPDB products",
	Long: "This repo helps to download / install / remove and manage the software of GPDB products",
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