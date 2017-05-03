package library

import (
	log "../../core/logger"
	"os/exec"
)


import (
	"../../core/arguments"
	"../../core/methods"
)


// Build executable shell script
func BuildExecutableBashScript (filename string, executable_filename string, args []string) error {

	log.Println("Starting program to build executable script to install the binaries")

	// Create the file
	log.Println("Creating temp bash executable file: "+ filename)
	err := methods.CreateFile(filename)
	if err != nil { return err }

	// Create arguments on what needs to be written to the file
	log.Println("Generating installer arguments to be passed to the file: "+ filename)
	executable_line := "/bin/sh " +  executable_filename + " &>/dev/null << EOF"
	var passArgs []string
	passArgs = append(passArgs, executable_line)
	for _ , v := range args {
		passArgs = append(passArgs, v)
	}
	passArgs = append(passArgs, "EOF")

	// Write to the file
	log.Println("Building the bash script file: "+ filename)
	err = methods.WriteFile(filename, passArgs)
	if err != nil { return err }

	return nil
}


// Execute shell script when called
func ExecuteBinaries(binary_file string, bashfilename string, script_options []string) error {

	// Build a quick shell script to install binaries
	// Filename name
	filename := arguments.TempDir + bashfilename

	// Cleanup the file if already exists (ignore error if any)
	_ = CleanupTempFile(filename)

	// Create the shell script
	err := BuildExecutableBashScript(filename, binary_file, script_options)
	if err != nil { return err }

	// Execute the installer script
	log.Println("Executing the bash script: "+ filename )
	_, err = exec.Command("/bin/sh", filename).Output()
	if err != nil { return err}

	// Cleanup the tempFile
	err = CleanupTempFile(filename)
	if err != nil { return err}

	return nil
}


// Cleanup the temp files once done
func CleanupTempFile (filename string) error {

	// Remove any tempfile
	log.Println("Cleaning up the file \""+ filename + "\", if found")
	err := methods.DeleteFile(filename)
	if err != nil { return err }

	return nil
}