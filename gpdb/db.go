package main

import (
	"bytes"
	"fmt"
)

func stopAllDb() {

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
	areAllProcessDown()

	// Cleanup temp files.
	deleteFile(StopScriptLoc)
}

// Check if the process are all down
func areAllProcessDown() {

	Debugf("Checking if the all the database process are down")
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
func isDbHealthy(sourcePath, port string) bool {

	Debugf("Checking if the database is healthy")
	// Query string
	queryString := "select 1"
	if sourcePath == "" {
		_, err := executeOsCommandOutput("psql", "-p", port, "-d", "template1", "-Atc", queryString)
		if err != nil {
			Fatalf("Error in checking database health, err: %v", err)
		}
	} else {
		content := readFile(sourcePath)
		c := contentExtractor(content, fmt.Sprintf("/%s/ {print $2}", "export PGPORT="), []string{"FS", "="})
		p := Config.CORE.TEMPDIR + "db_health.sh"
		createFile(p)
		writeFile(p, []string{
			"source " + sourcePath,
			fmt.Sprintf("psql -d template1 -p %s -Atc \"%s\"", removeBlanks(c.String()), queryString),
		})
		_, err := executeOsCommandOutput("/bin/sh", p) // if command execute its healthy
		if err != nil {
			deleteFile(p)
			return false
		} else {
			deleteFile(p)
			return true
		}
	}

	return true
}

// Start the database if not started
func startDBifNotStarted(envFile string)  {

	// is the database running , then return
	if isDbHealthy(envFile, "") { // Database is started and running
		Infof("Database seems to be running, contining...")
	} else { // database is not running, lets start it up

		Warnf("Database is not started, attempting to start the database...")
		// Stop all database is not stopped unless asked not to stop it
		if !cmdOptions.Stop {
			stopAllDb()
		}

		// Start the database of concern
		startDB(envFile)

		// Check again if the database is healthy
		if !isDbHealthy(envFile, "") {
			Fatalf("Can't seems to start the database in the environment file \"%s\"exiting...", envFile)
		}
	}
}

// Start database
func startDB(envFile string) error {

	Infof("Attempting to start the database as per the environment file: %s", envFile)

	// BashScript
	tempFile := Config.CORE.TEMPDIR + "start.sh"
	createFile(tempFile)
	writeFile(tempFile, []string{
		"source "+ envFile,
		"gpstart -a",
	})
	executeOsCommand("/bin/sh", tempFile)
	if !isDbHealthy(envFile, "") { // Check if the database is healthy after start
		deleteFile(tempFile)
		Fatalf("Can't seems to start the database in the environment file \"%s\" exiting...", envFile)
	}
	deleteFile(tempFile)

	// Start Command Center WEB UI if available on this environment
	//TODO: GPCC start
	//if !core.IsValueEmpty(GPPERFMONHOME) {
	//	StartGPCC(GPCC_INSTANCE_NAME, GPPERFMONHOME)
	//}

	return nil
}

// Stop Database
func stopDB(envFile string) {

	Infof("Attempting to stop the database as per the environment file: %s", envFile)
	tempFile := Config.CORE.TEMPDIR + "stop.sh"
	createFile(tempFile)
	writeFile(tempFile, []string{
		"source " + envFile,
		"gpstop -af",
	})
	executeOsCommand("/bin/sh", tempFile)
	deleteFile(tempFile)

	// Start Command Center WEB UI if available on this environment
	// TODO: GPCC Stop
	//if !core.IsValueEmpty(GPPERFMONHOME) {
	//	StopGPCC(GPCC_INSTANCE_NAME, GPPERFMONHOME)
	//}
}