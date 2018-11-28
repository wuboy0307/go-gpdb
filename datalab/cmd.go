package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

// Global Parameter
var (
	cmdOptions Command
)
const (
	defaultHostname = "gpdb"
	defaultSegments = 0
	defaultCpu = 2
	defaultMemory = 4096
	defaultOs = "bento/centos-7.5"
	defaultSubnet = "192.168.99.100"
)

// Command line options
type Command struct {
	Cpu 		int
	Memory 		int
	Standby 	bool
	Os 			string
	Subnet 		string
	Hostname    string
	Segments    int
	Debug       bool
	Token 		string
	GoGPDBPath  string
	GlobalStatus bool
}

// The create command.
var createCmd = &cobra.Command{
	Use:   "create",
	Aliases: []string{`c`},
	Short: "Create the vagrant environment",
	Long: "Create the vagrant environment and customize the environment",
	PostRun: func(cmd *cobra.Command, args []string) {
		saveConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		createVM()
	},
}

// All the usage flags of the create command
func createFlags() {
	createCmd.Flags().IntVarP(&cmdOptions.Cpu, "cpu", "c",defaultCpu,  "Customize the CPU of the vagrant VM that is going to created")
	createCmd.Flags().IntVarP(&cmdOptions.Memory, "memory", "m",defaultMemory,  "Customize the Memory of the vagrant VM that is going to created, units in MegaBytes")
	createCmd.Flags().IntVarP(&cmdOptions.Segments, "segments", "s",defaultSegments,  "Customize the number of segments host created during the vagrant provision")
	createCmd.Flags().BoolVar(&cmdOptions.Standby, "standby",false,  "Do you need a standby host vagrants to be created?")
	createCmd.Flags().StringVarP(&cmdOptions.Os, "os","o",defaultOs,"The vagrant OS to be used when being provisioned")
	createCmd.Flags().StringVarP(&cmdOptions.Subnet, "subnet","b", defaultSubnet,"The vagrant subnet to be used when being provisioned")
	createCmd.Flags().StringVarP(&cmdOptions.Hostname, "hostname","n",defaultHostname,"The name of the host that should be used when being provisioned")
}

// The ssh command.
var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "SSH to the vagrant environment",
	Long: "SSH the vagrant environment that is already created",
	Run: func(cmd *cobra.Command, args []string) {
		sshVM()
	},
}

// All the usage flags of the ssh command
func sshFlags() {
	sshCmd.Flags().StringVarP(&cmdOptions.Hostname, "hostname","n",defaultHostname,"The name of the host that should be used when being provisioned")
}

// The stop command.
var stopCmd = &cobra.Command{
	Use:   "stop",
	Aliases: []string{`s`},
	Short: "Stop the vagrant environment",
	Long: "Stop the vagrant environment that is already created",
	Run: func(cmd *cobra.Command, args []string) {
		stopVM()
	},
}

// All the usage flags of the stop command
func stopFlags() {
	stopCmd.Flags().StringVarP(&cmdOptions.Hostname, "hostname","n",defaultHostname,"The name of the host that should be used when being provisioned")
}

// The up command.
var upCmd = &cobra.Command{
	Use:   "up",
	Aliases: []string{`u`},
	Short: "Bring up the vagrant environment",
	Long: "Bring up the vagrant environment that is already created",
	Run: func(cmd *cobra.Command, args []string) {
		upVM()
	},
}

// All the usage flags of the up command
func upFlags() {
	upCmd.Flags().StringVarP(&cmdOptions.Hostname, "hostname","n",defaultHostname,"The name of the host that should be used when being provisioned")
}

// The status command.
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Status of the vagrant environment",
	Long: "Status the vagrant environment that is already created",
	Run: func(cmd *cobra.Command, args []string) {
		statusVM()
	},
}

// All the usage flags of the status command
func statusFlags() {
	statusCmd.Flags().StringVarP(&cmdOptions.Hostname, "hostname","n",defaultHostname,"The name of the host that should be used when being provisioned")
}

// The destroy command.
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy the vagrant environment",
	Long: "Destroy the vagrant environment that is already created",
	Run: func(cmd *cobra.Command, args []string) {
		destroyVM()
	},
}

// All the usage flags of the destroy command
func destroyFlags() {
	destroyCmd.Flags().StringVarP(&cmdOptions.Hostname, "hostname","n",defaultHostname,"The name of the host that should be used when being provisioned")
	destroyCmd.MarkFlagRequired("hostname")
}


// The update configuration command.
var updateCmd = &cobra.Command{
	Use:   "update-config",
	Aliases: []string{`uc`},
	Short: fmt.Sprintf("Update the configuration of the %s tool", programName),
	Long:  fmt.Sprintf("Update the configuration of the %s tool", programName),
	PostRun: func(cmd *cobra.Command, args []string) {
		saveConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		updateConfig()
	},
}

// All the usage flags of the update config command
func updateConfigFlags() {
	updateCmd.Flags().StringVarP(&cmdOptions.Token, "token","t","","Pivotal Network API Token")
	updateCmd.Flags().StringVarP(&cmdOptions.GoGPDBPath, "location","l","","Location of the go-gpdb software that should be used to provision the VM's")
}

// The update configuration command.
var deleteCmd = &cobra.Command{
	Use:   "delete-config",
	Aliases: []string{`dc`},
	Short: fmt.Sprintf("Delete the configuration from the %s config file", programName),
	Long: fmt.Sprintf("Delete the configuration from the %s config file", programName),
	Run: func(cmd *cobra.Command, args []string) {
		deleteConfigKey()
	},
}

// All the usage flags of the delete config command
func deleteConfigFlags() {
	deleteCmd.Flags().StringVarP(&cmdOptions.Hostname, "hostname","n",defaultHostname, "The name of the host that should be used when being provisioned")
	deleteCmd.MarkFlagRequired("hostname")
}

// The list command.
var listCmd = &cobra.Command{
	Use:   "list",
	Aliases: []string{`l`},
	Short: fmt.Sprintf("list all the configuration from the %s config file", programName),
	Long: fmt.Sprintf("list all the configuration from the %s config file", programName),
	Run: func(cmd *cobra.Command, args []string) {
		listVM()
	},
}

// All the usage flags of the list command
func listFlags() {
	listCmd.Flags().BoolVarP(&cmdOptions.GlobalStatus, "global-status","g",false, "Also show the vagrant global status of all the VM's")
}

// The root CLI.
var rootCmd = &cobra.Command{
	Use:   fmt.Sprintf("%s [command]", programName),
	Short: "Manages all the vagrant environments",
	Long: "The program manages like create / destroy / stop / list and helps to login to vagrant installation",
	Version: programVersion,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Before running any command Setup the logger log level
		initLogger(cmdOptions.Debug)
		// Load Configuration
		config()
		// Check if the token and location of vagrant is set before running any command
		if cmd.Use != "update-config" {
			if IsValueEmpty(Config.APIToken) {
				Fatalf(apiTokenMissing, programName)
			}
			if IsValueEmpty(Config.GoGPDBPath) {
				Fatalf(goGPDBLocationMissing, programName)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 { // if no argument specified throw the help menu on the screen
			cmd.Help()
		}
	},
}

// Initialize all the Cobra CLI
func init ()  {
	// root command flags
	rootCmd.PersistentFlags().BoolVarP(&cmdOptions.Debug, "verbose", "v", false, "Enable verbose or debug logging")
	// Attach the sub command to the root command.
	rootCmd.AddCommand(createCmd)
	createFlags()
	rootCmd.AddCommand(upCmd)
	upFlags()
	rootCmd.AddCommand(sshCmd)
	sshFlags()
	rootCmd.AddCommand(stopCmd)
	stopFlags()
	rootCmd.AddCommand(statusCmd)
	statusFlags()
	rootCmd.AddCommand(destroyCmd)
	destroyFlags()
	rootCmd.AddCommand(updateCmd)
	updateConfigFlags()
	rootCmd.AddCommand(deleteCmd)
	deleteConfigFlags()
	rootCmd.AddCommand(listCmd)
	listFlags()
}