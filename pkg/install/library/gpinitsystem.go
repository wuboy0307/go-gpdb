package library

import (
	"../objects"
	"../../core/arguments"
	"../../core/methods"
	log "../../core/logger"
	"strconv"
	"os/exec"
	"os"
	"strings"
)


func BuildGpInitSystemConfig(t string) error {

	log.Println("Building gpinitsystem config file for this installation")

	// Set the values of the below parameters
	objects.GpInitSystemConfig.ArrayName = "gp_" + arguments.RequestedInstallVersion + "_" + t
	objects.GpInitSystemConfig.SegPrefix = "gp_" + arguments.RequestedInstallVersion + "_" + t
	objects.GpInitSystemConfig.DBName = objects.DBName
	objects.GpInitSystemConfig.MasterDir = strings.TrimSuffix(arguments.EnvYAML.Install.MasterDataDirectory, "/")
	arguments.EnvYAML.Install.SegmentDataDirectory = strings.TrimSuffix(arguments.EnvYAML.Install.SegmentDataDirectory, "/")

	// Store the hostname of the hostfile
	objects.MasterHostFileList = arguments.TempDir + objects.MachineFileList
	_ = methods.DeleteFile(objects.MasterHostFileList)
	err := methods.CreateFile(objects.MasterHostFileList)
	if err != nil { return err }
	var host []string
	host = append(host, objects.GpInitSystemConfig.MasterHostName)
	err = methods.WriteFile(objects.MasterHostFileList, host)
	if err != nil { return err }

	// Check if we have the last used port base file
	log.Println("Obtaining ports to be set for PORT_BASE")
	PORT_BASE , _ := DoWeHavePortBase(arguments.FutureRefDir, objects.PortBaseFileName, "PORT_BASE")
	if PORT_BASE != "" {
		PORT_BASE = string(PORT_BASE)[11:]
	} else {
		log.Warn("Didn't find PORT_BASE in the file, setting it to default value: " + strconv.Itoa(objects.PORT_BASE))
		PORT_BASE = strconv.Itoa(objects.PORT_BASE)
	}

	// Check if the port is available or not
	pb, err := strconv.Atoi(PORT_BASE)
	if err != nil { return err }
	pb, err = IsPortUsed(pb, arguments.EnvYAML.Install.TotalSegments)
	if err != nil { return err }
	log.Println("Setting the PORT_BASE has: "+ strconv.Itoa(pb))
	objects.GpInitSystemConfig.PortBase = pb

	// Check if we have the last used master port file
	log.Println("Obtaining ports to be set for MASTER_PORT")
	MASTER_PORT , _ := DoWeHavePortBase(arguments.FutureRefDir, objects.PortBaseFileName, "MASTER_PORT")
	if MASTER_PORT != "" {
		MASTER_PORT = string(MASTER_PORT)[13:]
	} else {
		log.Warn("Didn't find MASTER_PORT in the file, setting it to default value: " + strconv.Itoa(objects.MASTER_PORT))
		MASTER_PORT = strconv.Itoa(objects.MASTER_PORT)
	}

	// MASTER PORT
	mp, err := strconv.Atoi(MASTER_PORT)
	if err != nil { return err }
	mp, err = IsPortUsed(mp, 1)
	if err != nil { return err }
	log.Println("Setting the MASTER_PORT has: "+ strconv.Itoa(mp))
	objects.GpInitSystemConfig.MasterPort = mp

	// Build gpinitsystem config file
	objects.GpInitSystemConfigDir = arguments.TempDir + objects.GpInitSystemConfigFile + "_" + arguments.RequestedInstallVersion + "_" + t
	log.Println("Creating the gpinitsystem config file at: " + objects.GpInitSystemConfigDir)
	_ = methods.DeleteFile(objects.GpInitSystemConfigDir)
	err = methods.CreateFile(objects.GpInitSystemConfigDir)
	if err != nil { return err }

	// Write the below content to config file
	var ToWrite []string
	var primaryDir string
	ToWrite = append(ToWrite,"ARRAY_NAME=" + objects.GpInitSystemConfig.ArrayName )
	ToWrite = append(ToWrite,"MACHINE_LIST_FILE=" + objects.MasterHostFileList)
	ToWrite = append(ToWrite,"SEG_PREFIX=" + objects.GpInitSystemConfig.SegPrefix )
	ToWrite = append(ToWrite,"MASTER_HOSTNAME=" + objects.GpInitSystemConfig.MasterHostName )
	ToWrite = append(ToWrite,"MASTER_DIRECTORY=" + objects.GpInitSystemConfig.MasterDir )
	ToWrite = append(ToWrite,"PORT_BASE=" + strconv.Itoa(objects.GpInitSystemConfig.PortBase) )
	ToWrite = append(ToWrite,"MASTER_PORT=" + strconv.Itoa(objects.GpInitSystemConfig.MasterPort) )
	ToWrite = append(ToWrite,"DATABASE_NAME=" + objects.GpInitSystemConfig.DBName)
	for i:=1; i<=arguments.EnvYAML.Install.TotalSegments; i++ {
		primaryDir = primaryDir + " "+ arguments.EnvYAML.Install.SegmentDataDirectory
	}
	ToWrite = append(ToWrite,"declare -a DATA_DIRECTORY=(" + primaryDir + " )")
	err = methods.WriteFile(objects.GpInitSystemConfigDir, ToWrite)
	if err != nil { return err }

	return nil
}


func ExecuteGpInitSystem() error {

	log.Println("Executing gpinitsystem to install GPDB software")

	// Initalize the command
	cmdOut := exec.Command("gpinitsystem", "-c", objects.GpInitSystemConfigDir , "-h", objects.MasterHostFileList, "-a")

	// Attach the os output from the screen
	cmdOut.Stdout = os.Stdout

	// Start the program
	err := cmdOut.Start()
	if err != nil { return err }

	// wait till it ends
	err = cmdOut.Wait()
	if err != nil { return err }

	return nil
}