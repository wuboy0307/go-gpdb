package library

import (
	log "../../core/logger"
	"../../core/methods"
	"../../core/arguments"
	"../objects"
	"os/exec"
	"bytes"
	"errors"
)

func StopAllDB() error {

	SourceGPDBPath()

	// Get all database running
	log.Println("Stopping all the database running on this host, to free up semaphore for this installation")

	// Can't seems to find a simple way to stop all database, so we will built the below
	// simple shell script to execute the stop database command
	cleanupScript := "ps -ef | grep silent | grep master | while read line; " +
		         "do " +
				"GPHOME=`echo $line|awk '{print $8}'| rev | cut -d'/' -f3- | rev`;"+
				"export MASTER_DATA_DIRECTORY=`echo $line|awk '{print $10}'`;"+
				"export PGPORT=`echo $line|awk '{print $12}'`;"+
				"export PGDATABASE=template1;"+
				"source $GPHOME/greenplum_path.sh;"+
				"gpstop -af;" +
			 "done"

	// Execute the command
	StopScriptLoc := arguments.TempDir + "stop_all_db.sh"
	var StopScript []string
	StopScript = append(StopScript, cleanupScript)
	err := ExecuteBash(StopScriptLoc, StopScript)
	if err != nil { return err }

	// Send a warning message if the process is not completely stopped.
	cmdOut, _ := exec.Command("pgrep", "postgres").Output()
	var EmptyBytes []byte
	if !bytes.Equal(cmdOut, EmptyBytes) {
		log.Warn("Can't stop all postgres process, seems like some are left behind")
	} else {
		log.Println("All Postgres process are stopped on this server")
	}

	// Cleanup temp files.
	methods.DeleteFile(StopScriptLoc)

	return nil
}

// Start database
func StartDB() error {

	log.Println("Attempting to start the database as per the environment file: " + objects.EnvFileName)

	// BashScript
	var BashSricpt []string
	BashSricpt = append(BashSricpt, "source " + objects.EnvFileName)
	BashSricpt = append(BashSricpt, "gpstart -a")

	// Create the file
	temp_file := arguments.TempDir + "start.sh"

	// Execute the script
	err := ExecuteBash(temp_file, BashSricpt)
	if err != nil { return err }

	// Check if the database is healthy after start
	err = IsDBHealthy()
	if err != nil { return errors.New("Can't seems to start the database in the environment file \"" + objects.EnvFileName + "\"exiting...") }

	return nil
}

// Stop Database
func StopDB() error {

	log.Println("Attempting to stop the database as per the environment file: " + objects.EnvFileName)

	// BashScript
	var BashSricpt []string
	BashSricpt = append(BashSricpt, "source " + objects.EnvFileName)
	BashSricpt = append(BashSricpt, "gpstop -af")

	// Create the file
	temp_file := arguments.TempDir + "stop.sh"

	// Execute the script
	err := ExecuteBash(temp_file, BashSricpt)
	if err != nil { return err }

	return nil
}

// Start the database if not started
func StartDBifNotStarted() error {

	// is the database running , then return
	err := IsDBHealthy()
	if err == nil { // Database is started and running
		log.Println("Database seems to be running, contining...")
		return nil
	} else { // database is not running, lets start it up

		log.Warn("Database is not started, attempting to start the database...")

		// Stop all database is not stopped
		err := StopAllDB()
		if err != nil { return err }

		// Start the database of concern
		err = StartDB()
		if err != nil { return err }

		// Check again if the database is healthy
		err = IsDBHealthy()
		if err != nil { return errors.New("Can't seems to start the database in the environment file \"" + objects.EnvFileName + "\"exiting...") }
	}

	return nil
}