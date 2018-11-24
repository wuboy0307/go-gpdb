package main

// Build uninstall script
func (i *Installation) createUninstallScript() error {

	// Uninstall script location
	uninstallFile := Config.INSTALL.UNINTSALLDIR + "uninstall_" + cmdOptions.Version + "_" + i.timestamp
	Infof("Creating Uninstall file for this installation at: " + uninstallFile)

	// Query
	// TODO: mutli enviornment ssh
	queryString := "select $$ps -ef|grep postgres|grep -v grep|grep $$ ||  port " +
		"|| $$ | awk '{print $2}'| xargs -n1 /bin/kill -11 &>/dev/null $$ " +
		"from gp_segment_configuration union select $$rm -rf /tmp/.s.PGSQL.$$ || port || $$*$$ " +
		"from gp_segment_configuration union select $$rm -rf $$ || fselocation from pg_filespace_entry"

	// Execute the query
	cmdOut, err := executeOsCommandOutput("psql", "-p", i.GPInitSystem.MasterPort, "-d", "template1", "-Atc", queryString)
	if err != nil {
		Fatalf("Error in running uninstall command on database, err: %v", err)
	}

	// Create the file
	createFile(uninstallFile)
	writeFile(uninstallFile, []string{
		string(cmdOut),
		"rm -rf "+ Config.INSTALL.ENVDIR +"env_" + cmdOptions.Version + "_"+ i.timestamp,
		"rm -rf " + uninstallFile,

	})
	return nil
}
