package main

import (
	"bytes"
)

func (i *Installation) stopAllDb() error {

	// Get all database running
	i.sourceGPDBPath()
	Infof("Stopping all the database running on this host, to free up semaphore for this installation")

	// Can't seems to find a simple way to stop all database, so we will built the below
	// simple shell script to execute the stop database command
	cleanupScript := "ps -ef | grep silent | grep master | while read line; " +
		"do " +
		"GPHOME=`echo $line|awk '{print $8}'| rev | cut -d'/' -f3- | rev`;" +
		"export MASTER_DATA_DIRECTORY=`echo $line|awk '{print $10}'`;" +
		"export PGPORT=`echo $line|awk '{print $12}'`;" +
		"export PGDATABASE=template1;" +
		"source $GPHOME/greenplum_path.sh;" +
		"gpstop -af;" +
		"done"

	// Execute the command
	StopScriptLoc := Config.CORE.TEMPDIR + "stop_all_db.sh"
	writeFile(StopScriptLoc, []string{
		cleanupScript,
		"ps -ef | egrep \"gpmon|gpmonws|lighttpd\" | grep -v grep | awk '{print $2}' | xargs -n1 /bin/kill -11 &>/dev/null; echo > /dev/null",
	})
	executeOsCommand("/bin/sh", StopScriptLoc)
	i.areAllProcessDown()

	// Cleanup temp files.
	deleteFile(StopScriptLoc)

	return nil
}

// Check if the process are all down
func (i *Installation) areAllProcessDown() {
	// Send a warning message if the process is not completely stopped.
	cmdOut, _ := executeOsCommandOutput("pgrep", "postgres")
	var EmptyBytes []byte
	if !bytes.Equal(cmdOut, EmptyBytes) {
		Warn("Can't stop all postgres process, seems like some are left behind")
	} else {
		Info("All Postgres process are stopped on this server")
	}
}

// Check if the database is healthy
func (i *Installation) isDbHealthy() {

	// Query string
	Infof("Ensuring the database is healthy...")
	queryString := "select 1"

	_, err := executeOsCommandOutput("psql", "-p", i.GPInitSystem.MasterPort, "-d", "template1", "-Atc", queryString)
	if err != nil {
		Fatalf("Error in checking database health, err: %v", err)
	}
}