package library

import (
	log "../../core/logger"
	"../objects"
	"../../core/arguments"
	"../../core/methods"
)


// Installed GPPERFMON software to the database
func InstallGpperfmon() error {

	log.Println("Installing GPPERFMON Database..")

	// Check if the password for gpmon is provided on config.yml, else set it to default "changeme"
	if methods.IsValueEmpty(arguments.EnvYAML.Install.GpMonPass) {
		log.Warn("The value of GPMON_PASS is missing, setting the password of gpperfmon & gpmon to \"changeme\"")
	}

	// Installation script
	var install_script []string
	install_script = append(install_script, "source " + objects.EnvFileName)
	install_script = append(install_script, "gpperfmon_install --enable --password "+ arguments.EnvYAML.Install.GpMonPass +" --port $PGPORT")
	install_script = append(install_script, "echo \"host      all     gpmon    0.0.0.0/0    md5\" >> $MASTER_DATA_DIRECTORY/pg_hba.conf")
	install_script = append(install_script, "echo \"host     all         gpmon         ::1/128       md5\" >> $MASTER_DATA_DIRECTORY/pg_hba.conf")

	// Execute the script
	err := ExecuteBash("install_gpperfmon.sh", install_script)
	if err != nil { return err }

	return nil
}
