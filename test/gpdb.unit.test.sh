#!/usr/bin/env bash

#Directories
export script=$0
export script_basename=`basename $script`
export script_dir=`dirname $script`/..
cd $script_dir
export base_dir=`pwd`
source /vagrant/scripts/functions.h
source <(parse_yaml /vagrant/gpdb/config.yml)

#Parameters
endpoint="https://network.pivotal.io"
slug="pivotal-gpdb"
token=`env | grep UAA_API_TOKEN | cut -d'=' -f2`
access_token=""

#Abort the script if found a failure
abort() {
	log "$FAIL Return Code: [$1]"
}

#Check if the database is accessible
function check_db() {
    local ver=`echo $1|sed 's/"//g'`
    local env_loc=${BASE_DIR}/${APPLICATION_NAME}${ENV_DIR}
    {
        source ${env_loc}/env_${ver}_*; ${GPHOME}/bin/psql -d template1 -Atc "select 1" &
    } &> /dev/null
    spinner $! "Checking if the database with version ${ver} is healthy"
    if [[ $? -ne 0 ]]; then wait $!; abort $?; fi
}

#Download the version and test
function download_gpdb_version() {
    local ver=`echo $1|sed 's/"//g'`
    local filename="/tmp/download_${ver}.out"
    cd ${base_dir}/gpdb
    { go run *.go d -v ${ver} & } &> ${filename}
    spinner $! "Downloading the version ${ver}"
    if [[ $? -ne 0 ]]; then wait $!; abort $?; fi
}

#Install the GPDB version and test
function install_gpdb_version() {
    local ver=`echo $1|sed 's/"//g'`
    local filename="/tmp/install_${ver}.out"
    cd ${base_dir}/gpdb
    { go run *.go i -v ${ver} --standby & } &> ${filename}
    spinner $! "Install the version ${ver}"
    if [[ $? -ne 0 ]]; then wait $!; abort $?; fi
}

#Download and install command center
function download_n_install_command_center() {
    local ver=`echo $1|sed 's/"//g'`
    local releaseID=`echo ${releases} | jq .releases[] | jq 'select(.version == "'${ver}'")' | jq .id`
    local incr=0
    products=`curl -s ${endpoint}/api/v2/products/${slug}/releases/${releaseID} -H "Authorization: Bearer ${access_token}" | jq '.'`
    cc_products=`echo ${products} | jq .file_groups[] | jq 'select(.name == "Greenplum Command Center")' | jq -r '.product_files[] | "\(.file_version) \(.file_type),"'`
    OIFS=$IFS
    IFS=","
    echo ${cc_products} | while read -r line
    do
        type=`echo ${line}| awk '{print $2}'`
        cc_version=`echo ${line}| awk '{print $1}'`
        if [[ ${type} == "Documentation" ]]; then
            incr=$((${incr}+1))
        else
            incr=$((${incr}+1))
            download_gpcc_version ${ver} ${cc_version} ${incr}
            install_gpcc_version ${ver} ${cc_version}
            check_cc ${ver} ${cc_version}
        fi
    done
    IFS=$OIFS
}

#Download the gpcc version for the database version
function download_gpcc_version() {
    local ccver=`echo $2|sed 's/"//g'`
    local filename="/tmp/download_gpcc_$1_${ccver}.out"
    cd ${base_dir}/gpdb
    { TEST_PROMPT_CHOICE=$3 go run *.go d -p gpcc -v $1 & } &> ${filename}
    spinner $! "Downloading the cc version ${ccver} for gpdb version $1"
    if [[ $? -ne 0 ]]; then wait $!; abort $?; fi
}

# Install the gpcc on the database version
function install_gpcc_version() {
    local ccver=`echo $2|sed 's/"//g'`
    local filename="/tmp/install_gpcc_$1_${ccver}.out"
    cd ${base_dir}/gpdb
    { TEST_YES_CONFIRMATION="y" go run *.go i -p gpcc -v $1 -c ${ccver} & } &> ${filename}
    spinner $! "Install the cc version ${ccver} on gpdb Version $1"
    if [[ $? -ne 0 ]]; then wait $!; abort $?; fi
}

#Check if command center is working
function check_cc() {
    local ccver=`echo $2|sed 's/"//g'`
    local filename="/tmp/install_gpcc_$1_${ccver}.out"
    local cc_url=`tail -10 ${filename} | grep "GPCC Web URL" | grep -Eo '(http|https)://[^,/"]+'`
    { curl -s --head  --request GET ${cc_url} | grep 200 & } &>/dev/null
    spinner $! "Checking if the CC Url ${ccver} on gpdb Version $1 is working"
    if [[ $? -ne 0 ]]; then wait $!; abort $?; fi
}

#Set env of the version and test
function env_version() {
    local ver=`echo $1|sed 's/"//g'`
    local filename="/tmp/env_${ver}.out"
    cd ${base_dir}/gpdb
    { go run *.go e -v ${ver} & } &> ${filename}
    spinner $! "Setting the environment of version ${ver}"
    if [[ $? -ne 0 ]]; then wait $!; abort $?; fi
}

#Remove the version and test
function remove_version() {
    local ver=`echo $1|sed 's/"//g'`
    local filename="/tmp/remove_${ver}.out"
    cd ${base_dir}/gpdb
    { go run *.go r -v ${ver} & } &> ${filename}
    spinner $! "Remove the version ${ver}"
    if [[ $? -ne 0 ]]; then wait $!; abort $? ${filename}; fi
}

#Cleanup the files for the next run to proceed
function cleanup() {
   local download_loc=${BASE_DIR}/${APPLICATION_NAME}${DOWNLOAD_DIR}
   {
        rm -rf /usr/local/greenplum* &
        rm -rf ${download_loc}* &
   } &> /dev/null
   spinner $! "Cleaning the binaries in the /usr/local & ${download_loc}"
   if [[ $? -ne 0 ]]; then wait $!; abort $?; fi
}

#Network login
function authenticate() {
    access_token=`curl -s -X POST ${endpoint}/api/v2/authentication/access_tokens -d '{"refresh_token":"'${token}'"}' | jq .access_token`
}

#Unit test for all Greenplum Releases
authenticate
releases=`curl -s ${endpoint}/api/v2/products/${slug}/releases -H "Authorization: Bearer ${access_token}"`
for i in `echo ${releases} | jq .releases[] | jq 'select(.version != "Pivotal Greenplum Text")' | jq .version`
do
    banner "TEST: GO GPDB Version: ${i}"
    authenticate
    download_gpdb_version ${i}
    install_gpdb_version  ${i}
    check_db ${i}
    download_n_install_command_center ${i}
    env_version      ${i}
    remove_version   ${i}
    cleanup
done