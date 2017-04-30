package argParser

import (
	"flag"
	"os"
	"fmt"
)

import (
	"../../pkg/core/arguments"
	"../../pkg/core/methods"
)

// OS Argument Parser
func ArgParser() {

	// Download Command Parser
	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)
	DownloadProductFlag := downloadCmd.String("p", "gpdb", "What product do you want to Install? [OPTIONS: gpdb, gpcc, gpextras]")
	DownloadVersionFlag := downloadCmd.String("v", "", "OPTIONAL: Which GPDB version software do you want to list ?")

	// Install Command Parser
	InstallCmd := flag.NewFlagSet("install", flag.ExitOnError)
	InstallProductFlag := InstallCmd.String("p", "gpdb", "What product do you want to Install ?")
	InstallVersionFlag := InstallCmd.String("v", "", "Which version do you want to dowload (only works with product equal to gpdb) ?")

	// Remove Command Parser
	RemoveCmd := flag.NewFlagSet("remove", flag.ExitOnError)
	RemoveVersionFlag := RemoveCmd.String("v", "", "Provide the version from the installed list to remove")

	// Environment Command Parser
	EnvCmd := flag.NewFlagSet("env", flag.ExitOnError)
	EnvListFlag := EnvCmd.String("l", "", "List all the version that is currently installed")
	EnvVersionFlag := EnvCmd.String("v", "", "Provide the version from the installed list to remove")

	// If no COMMAND keyword provided then show the help menu.
	if len(os.Args) == 1 {
		ShowHelp()
	}

	// If there is a command keyword provided then check to what is it and then parse the appropriate options
	switch os.Args[1] {
		case "download":                                                // Download command parser
			downloadCmd.Parse(os.Args[2:])
		case "install":                                                 // Install command parser
			InstallCmd.Parse(os.Args[2:])
		case "remove":                                                  // Remove command parser
			RemoveCmd.Parse(os.Args[2:])
		case "version":                                                 // Version of the software
			fmt.Printf("Version: %.1f\n", arguments.Version)
			os.Exit(0)
		case "help":                                                    // Help Menu
			ShowHelp()
		default:                                                        // Error if command is invalid
			fmt.Printf("ERROR: %q is not valid command.\n", os.Args[1])
			ShowHelp()
	}

	// If the command send is download, then parse the commandline arguments
	if downloadCmd.Parsed() {

		// If the product parameter is passed, then check if its valid value.
		if *DownloadProductFlag != "" {

			// If its valid then we are going to store it
			if methods.Contains(arguments.AcceptedDownloadProduct, *DownloadProductFlag) {
				arguments.RequestedDownloadProduct = *DownloadProductFlag
			} else { // Else print error to choose the right value
				fmt.Println("ERROR: Invalid options provided to the argument -p \n")
				fmt.Print("Usage of download: \n")
				downloadCmd.PrintDefaults()
				os.Exit(2)
			}
		}

		// If the version parameter is passed, then store the value
		if *DownloadVersionFlag != "" {
			arguments.RequestedDownloadVersion = *DownloadVersionFlag
		}
	}


	// All the below command line parse will be updated when the function is written
	if InstallCmd.Parsed() {
		if *InstallProductFlag == "" {
			fmt.Println("Please supply the recipient using -recipient option.")
			return
		}

		if *InstallVersionFlag == "" {
			fmt.Println("Please supply the message using -message option.")
			return
		}

		fmt.Printf("Your message is sent to %q.\n", *InstallProductFlag)
		fmt.Printf("Message: %q.\n", *InstallVersionFlag)
	}


	if RemoveCmd.Parsed() {

		if *RemoveVersionFlag == "" {
			fmt.Println("Please supply the message using -message option.")
			return
		}

		fmt.Printf("Your message is sent to %q.\n", *RemoveVersionFlag)
	}


	if EnvCmd.Parsed() {

		if *EnvListFlag == "" {
			fmt.Println("Please supply the message using -message option.")
			return
		}
		if *EnvVersionFlag == "" {
			fmt.Println("Please supply the message using -message option.")
			return
		}

		fmt.Printf("Your message is sent to %q.\n", *EnvListFlag)
		fmt.Printf("Message: %q.\n", *EnvVersionFlag)
	}

}