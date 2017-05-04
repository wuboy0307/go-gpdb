package library

import (
	log "../../core/logger"
	"../../core/methods"
	"../../core/arguments"
	"os/exec"
	"os"
	"bytes"
)

func StopDB() error {

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

	// Write it a file
	StopScriptLoc := arguments.TempDir + "stop_db.sh"
	methods.DeleteFile(StopScriptLoc)
	err := methods.CreateFile(StopScriptLoc)
	if err != nil { return err }
	var StopScript []string
	StopScript = append(StopScript, cleanupScript)
	err = methods.WriteFile(StopScriptLoc, StopScript)
	if err != nil { return err }

	// Execute the script
	cmd := exec.Command("/bin/sh", StopScriptLoc)
	cmd.Stdout = os.Stdout
	err = cmd.Start()
	if err != nil { return err }
	err = cmd.Wait()
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
