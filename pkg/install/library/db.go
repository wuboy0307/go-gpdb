package library

import (
	"os/exec"
	log "../../core/logger"
	"../objects"
	"strconv"
	"../../core/arguments"
)

// Execute DB Query
func execute_db_query(query_string string, store_to_file bool, file_name string) (error) {

	master_port := strconv.Itoa(objects.GpInitSystemConfig.MasterPort)
	master_port = strconv.Itoa(3000)

	if store_to_file {
		cmd := exec.Command("psql", "-p", master_port , "-d", "template1", "-Atc", query_string)
		err := cmd.Run()
		if err != nil { return err }
	} else {
		cmd := exec.Command("psql", "-p", master_port , "-d", "template1", "-Atc", query_string)
		err := cmd.Run()
		if err != nil { return err }

	}

	return nil
}

// Check if the database is healthy
func IsDBHealthy() error {

	log.Println("Ensuring the database is healthy...")
	query_string := "select 1"
	err := execute_db_query(query_string, false, "")
	if err != nil { return err }

	return nil
}

// Build uninstall script
func CreateUnistallScript(t string) error {

	// Uninstall script location
	UninstallFileDir := arguments.UninstallDir + "uninstall_" + arguments.RequestedInstallVersion + "_" + t
	log.Println("Creating Uninstall file for this installation at: " + UninstallFileDir)
	query_string := "select \\$\\$ps -ef|grep postgres|grep -v grep|grep \\$\\$ ||  port " +
			"|| \\$\\$ | awk '{print \\$2}'| xargs -n1 /bin/kill -11 &>/dev/null\\$\\$ " +
			"from gp_segment_configuration union select \\$\\$rm -rf /tmp/.s.PGSQL.\\$\\$ || port || \\$\\$*\\$\\$ " +
		        "from gp_segment_configuration union select \\$\\$rm -rf \\$\\$ || fselocation from pg_filespace_entry"

	query_string = "select 1"

	err := execute_db_query(query_string, true, UninstallFileDir)
	if err != nil { return err }

	return nil
}
