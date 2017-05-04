package library

import (
	log "../../core/logger"
	"../../core/arguments"
	"../../core/methods"
	"../objects"
	"strconv"
)

// Create environment file of this installation
func CreateEnvFile(t string) error {

	// Environment file fully qualified path
	EnvFile := arguments.EnvFileDir + "env_" + arguments.RequestedInstallVersion + "_" + t
	log.Println("Creating environment file for this installation at: " + EnvFile)

	// Create the file
	err := methods.CreateFile(EnvFile)
	if err != nil { return err }

	// Build arguments to write
	var EnvFileContents []string
	EnvFileContents = append(EnvFileContents, "source " + objects.BinaryInstallLocation + "/greenplum_path.sh")
	EnvFileContents = append(EnvFileContents, "export MASTER_DATA_DIRECTORY=" + objects.GpInitSystemConfig.MasterDir + objects.GpInitSystemConfig.ArrayName + "-1")
	EnvFileContents = append(EnvFileContents, "export PGPORT=" + strconv.Itoa(objects.GpInitSystemConfig.MasterPort))
	EnvFileContents = append(EnvFileContents, "export PGDATABASE=" + objects.GpInitSystemConfig.DBName)

	// Write to EnvFile
	err = methods.WriteFile(EnvFile, EnvFileContents)
	if err != nil { return err }

	return nil
}