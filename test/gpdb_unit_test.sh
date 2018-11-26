#!/usr/bin/env bash

#Directories
export script=$0
export script_basename=`basename $script`
export script_dir=`dirname $script`/..
cd $script_dir
export base_dir=`pwd`
source /vagrant/scripts/functions.h

#Parameters
endpoint="https://network.pivotal.io"
slug="pivotal-gpdb"
token=`env | grep UAA_API_TOKEN | cut -d'=' -f2`

#Abort the script if found a failure
abort() {
	log "$FAIL Return Code: [$1]"
	exit $1
}

#Download the version and test
function download_gpdb_version() {
    local ver=`echo $1|sed 's/"//g'`
    cd ${base_dir}/gpdb
    { go run *.go d -v ${ver} & } &> /tmp/download_${ver}.out
    spinner $! "Downloading the version: ${ver}"
    if [[ $? -ne 0 ]]; then wait $!; abort $?; fi
}

#Install the GPDB version and test
function install_gpdb_version() {
    local ver=`echo $1|sed 's/"//g'`
    cd ${base_dir}/gpdb
    { go run *.go i -v ${ver} --standby & } &> /tmp/install_${ver}.out
    spinner $! "Install the version: ${ver}"
    if [[ $? -ne 0 ]]; then wait $!; abort $?; fi
}

#Set env of the version and test
function env_version() {
    local ver=`echo $1|sed 's/"//g'`
    cd ${base_dir}/gpdb
    { go run *.go e -v ${ver} & } &> /tmp/env_${ver}.out
    spinner $! "Setting the environment of version: ${ver}"
    if [[ $? -ne 0 ]]; then wait $!; abort $?; fi
}

#Remove the version and test
function remove_version() {
    local ver=`echo $1|sed 's/"//g'`
    cd ${base_dir}/gpdb
    { go run *.go r -v ${ver} & } &> /tmp/remove_${ver}.out
    spinner $! "Remove the version: ${ver}"
    if [[ $? -ne 0 ]]; then wait $!; abort $?; fi
}

#Network login
access_token=`curl -X POST ${endpoint}/api/v2/authentication/access_tokens -d '{"refresh_token":"'${token}'"}' | jq .access_token`

#Unit test for all Greenplum Releases
releases=`curl ${endpoint}/api/v2/products/${slug}/releases -H "Authorization: Bearer ${access_token}"`
for i in `echo ${releases} | jq .releases[] | jq 'select(.version != "Pivotal Greenplum Text")' | jq .version`
do
    banner "TEST: GO GPDB Version: {$i}"
    download_gpdb_version ${i}
    install_gpdb_version  ${i}
    env_version      ${i}
    remove_version   ${i}
done