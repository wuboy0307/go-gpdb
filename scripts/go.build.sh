#!/usr/bin/env bash
#!/bin/bash

set -e
source /vagrant/scripts/functions.h
source <(parse_yaml /vagrant/gpdb/config.yml)
set +e

## abort information
abort() {
	log "$FAIL Return Code: [$1]"
	exit $1
}

## cleanup
cleanup() {

	if [[ ${#src} -gt 0 ]]; then
		banner "Cleanup"
		for package in "${src[@]}"
		do
			{ rm -rf $package & } &>/dev/null
			spinner $! "Removing Temporary File: $package"
		done
	fi

	if [[ -z "${build_complete}" ]]; then
		banner "Failed"
		log "$FAIL /vagrant/scripts/go.build.sh"
	else
	    banner "Completed"
		log "$PASS /vagrant/scripts/go.build.sh"
	fi
	echo
}

trap cleanup EXIT

env_go_lang() {
    # Update Environment Variables if it doesn't exist
	if ! ( grep -q "# GOLANG" /etc/profile.d/gpdb.profile.sh &>/dev/null ); then
	    {
	        sudo echo '# GOLANG'
	        sudo echo 'export GOROOT=/usr/local/go'
	        sudo echo 'export GOPATH=$BASE_DIR'
	        sudo echo 'export GOBIN=/usr/local/bin'
	        sudo echo 'export PATH=$PATH:$GOROOT/bin:$GOBIN'
	    } >> /etc/profile.d/gpdb.profile.sh
		spinner $! "Update Environment Variables"
		if [[ $? -ne 0 ]]; then wait $!; abort $?; fi

	fi
}

## Install the go binaries if requested
go_install() {
	# Download GO
	{ wget -q https://storage.googleapis.com/golang/go$GO_BUILD.$OS-$ARCH.tar.gz -O $BASE_DIR/go.tar.gz & } &>/dev/null
	spinner $! "Downloading GO Binary: $GO_BUILD"
	if [[ $? -ne 0 ]]; then wait $!; abort $?; fi

	# Extract
	{ tar -C "/usr/local" -xzf $BASE_DIR/go.tar.gz & } &>/dev/null
	spinner $! "Extracting: $GO_BUILD"
	if [[ $? -ne 0 ]]; then wait $!; abort $?; fi

	{ rm -rf "$BASE_DIR/go.tar.gz" & } &>/dev/null
	spinner $! "Removing Temporary File: $BASE_DIR/go.tar.gz"

	# Notify
	log "$PASS GO Binary Version Installed: $GO_BUILD"

}

## Download the latest version of the gpdb cli from the github release link
download_latest_gpdb_cli() {
    env_go_lang

    # Set the environment
    source /etc/profile.d/gpdb.profile.sh

    # cleanup old release if found any
    { rm -rf $GOBIN/gpdb & } &>/dev/null
    spinner $! "Cleaning up old gpdb releases"

    # Download the latest build or release of gpdb cli
    { curl -s https://api.github.com/repos/pivotal-gss/go-gpdb/releases/latest \
      | grep "browser_download_url.*gpdb" \
      | grep -v "browser_download_url.*datalab" \
      | cut -d : -f 2,3 \
      | tr -d \" \
      | wget -qi - -O $GOBIN/gpdb &
    } &> /dev/null
    spinner $! "Downloading the latest release of the gpdb cli"

    # Setting up the permission to execute of the gpdb CLI
    {
      chmod +x $GOBIN/gpdb &
    } &> /dev/null
    spinner $! "Setting the execute permission to the gpdb cli"

    # Copy the config file to the home directory
    {
      cp /vagrant/gpdb/config.yml /home/gpadmin/
      chown gpadmin:gpadmin /home/gpadmin/config.yml
    } &> /dev/null
    spinner $! "Copying the config file to the directory: /home/gpadmin"
}

banner "Configuration"

##  Internet Connetivity
{ wget -q --tries=2 --timeout=5 --spider http://google.com & } &>/dev/null
spinner $! "Internet Connection"
if [[ $? -ne 0 ]]; then wait $!; abort $?; fi

## Install all the golang packages if the developer mode is one
if [[ $1 == "true" ]]; then
    # YAML: BASE_DIR
    { mkdir -p "$BASE_DIR" && test -w "$BASE_DIR" & } &>/dev/null
    spinner $! "YAML: BASE_DIR: $BASE_DIR"
    if [[ $? -ne 0 ]]; then wait $!; abort $?; fi

    # Hostname
    { ping -c 1 `hostname` & } &>/dev/null
    spinner $! "YAML: HOSTNAME: `hostname`"
    if [[ $? -ne 0 ]]; then wait $!; abort $?; fi

    banner "GO Binary Installation"
    if ! [[ -d "/usr/local/go" ]]; then
        go_install
    else
        # Compare 1 : 2 [EQ 0; GT 1; LT 2)
        compare_versions $(go_version) $GO_BUILD

        if [[ $? -lt 2 ]]; then
            log "$PASS GO Binary Version Required: $GO_BUILD (Installed: $(go_version))"
        else
            log "$FAIL GO Binary Version Required: $GO_BUILD (Installed: $(go_version))"

            # Backup Exisitng Build
            { mv /usr/local/go /usr/local/go.$(go_version) & } &>/dev/null
            spinner $! "Backing Up Existing Build"
            if [[ $? -ne 0 ]]; then wait $!; abort $?; fi

            # Call Installer
            go_install
        fi
    fi

		# env_go_lang # Should already be sourced from profile
		banner "GO Developer Packages"
		{
			if ! [ -x "$(command -v dep)" ]; then
				sudo curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh &
			fi
    } &>> /tmp/go.build.out
    spinner $! "GO: Package Manager"

		{ go get github.com/pivotal-gss/go-gpdb/gpdb & } &>/dev/null
		spinner $! "GO: GPDB Project Files"

		{
			cd $GOPATH/src/github.com/pivotal-gss/go-gpdb/gpdb
			dep init
		} &>> /tmp/go.build.out
		spinner $! "GO: Downloading Dependencies"

else
	banner "Download gpdb CLI"
	download_latest_gpdb_cli
fi

##  Download the latest build of the gpdb cli
# source /etc/profile.d/gpdb.profile.sh
build_complete=true
exit 0
