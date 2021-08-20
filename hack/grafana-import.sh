#!/bin/bash
#
# Copyright (c) 2021, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

#
# add the "-x" option to the shebang line if you want a more verbose output
#
#

#set -x 
#set +x
OPTSPEC=":hd:t:u:w:"

show_help() {
cat << EOF
Usage: $0 [-u Grafana_User] [-w Grafana_Password] [-p PATH] [-t TARGET_HOST]
Script to import dashboards into Grafana
    -u      Required. Grafana User
    -w      Required. Grafana Password
    -d      Required. Root path containing JSON dashboard files you want imported.
    -t      Required. The full URL of the target host

    -h      Display this help and exit.
EOF
}

###### Check script invocation options ######
while getopts "$OPTSPEC" optchar; do
    case "$optchar" in
        h)
            show_help
            exit
            ;;
        d)
            DASH_DIR="$OPTARG";;
        t)
            TARGET_HOST="$OPTARG";;
        u)
            GRAFANA_USER="$OPTARG";;
        w)
            GRAFANA_PASSWORD="$OPTARG";;
        \?)
          echo "Invalid option: -$OPTARG" >&2
          exit 1
          ;;
        :)
          echo "Option -$OPTARG requires an argument." >&2
          exit 1
          ;;
    esac
done

if [ -z "$DASH_DIR" ] || [ -z "$TARGET_HOST" ] || [ -z "$GRAFANA_USER" ] || [ -z "$GRAFANA_PASSWORD" ]; then
    show_help
    exit 1
fi

# set some colors for status OK, FAIL and titles
SETCOLOR_SUCCESS="echo -en \\033[0;32m"
SETCOLOR_FAILURE="echo -en \\033[1;31m"
SETCOLOR_NORMAL="echo -en \\033[0;39m"
SETCOLOR_TITLE_PURPLE="echo -en \\033[0;35m" # purple

# usage log "string to log" "color option"
function log_success() {
   if [ $# -lt 1 ]; then
       ${SETCOLOR_FAILURE}
       echo "Not enough arguments for log function! Expecting 1 argument got $#"
       exit 1
   fi

   timestamp=$(date "+%Y-%m-%d %H:%M:%S %Z")

   ${SETCOLOR_SUCCESS}
   printf "[%s] $1\n" "$timestamp"
   ${SETCOLOR_NORMAL}
}

function log_failure() {
   if [ $# -lt 1 ]; then
       ${SETCOLOR_FAILURE}
       echo "Not enough arguments for log function! Expecting 1 argument got $#"
       exit 1
   fi

   timestamp=$(date "+%Y-%m-%d %H:%M:%S %Z")

   ${SETCOLOR_FAILURE}
   printf "[%s] $1\n" "$timestamp"
   ${SETCOLOR_NORMAL}
}

function log_title() {
   if [ $# -lt 1 ]; then
       ${SETCOLOR_FAILURE}
       log_failure "Not enough arguments for log function! Expecting 1 argument got $#"
       exit 1
   fi

   ${SETCOLOR_TITLE_PURPLE}
   printf "|------------------------------------------------------------------------------------------------------------------------|\n"
   printf "|%s|\n" "$1";
   printf "|------------------------------------------------------------------------------------------------------------------------|\n"
   ${SETCOLOR_NORMAL}
}

### API KEY GENERATION

KEYNAME=$(head /dev/urandom | LC_ALL=C tr -dc A-Za-z0-9 | head -c 13 | cut -c -7)
KEYLENGTH=70
GENERATE_POST_DATA="{\"name\": \"${KEYNAME}\", \"role\": \"Admin\", \"secondsToLive\": 3600 }"

if [ -n "$GRAFANA_USER" ] || [ -n "$GRAFANA_PASSWORD" ] || [ -n "$TARGET_HOST" ]; then
    KEY=$(curl -X POST -H "Content-Type: application/json" -d "${GENERATE_POST_DATA}" http://${GRAFANA_USER}:${GRAFANA_PASSWORD}@${TARGET_HOST}/api/auth/keys | jq -r '.key')
    if [ ${#KEY} -ge $KEYLENGTH ]; then
        log_title "---- API Key Generated successfully, correct character number generated in API Key, we're going into the next step -----"
    else
        log_title "------------------------------------------------- API Key is not valid ! ----------------------------------------------"
        log_failure "$KEY does not contain the correct information to do API actions. Please Check, and try again."
        exit 1

    fi
else
    log_title "----------------- One of the parameters is not correct -----------------"
    log_failure "Set correct parameters and try again."
    exit 1
fi

if [ -d "$DASH_DIR" ]; then
    DASH_LIST=$(find "$DASH_DIR" -mindepth 1 -name \*.json)
    if [ -z "$DASH_LIST" ]; then
        log_title "----------------- $DASH_DIR contains no JSON files! -----------------"
        log_failure "Directory $DASH_DIR does not appear to contain any JSON files for import. Check your path and try again."
        exit 1
    else
        FILESTOTAL=$(echo "$DASH_LIST" | wc -l)
        log_title "------------------------------------------- Starting import of $FILESTOTAL dashboards -------------------------------------------"
    fi
else
    log_title "----------------- $DASH_DIR directory not found! -----------------"
    log_failure "Directory $DASH_DIR does not exist. Check your path and try again."
    exit 1
fi

NUMSUCCESS=0
NUMFAILURE=0
COUNTER=0

for DASH_FILE in $DASH_LIST; do
    COUNTER=$((COUNTER + 1))
    echo "Import $COUNTER/$FILESTOTAL: $DASH_FILE..."
    echo '{ "overwrite": true, "dashboard":' > tmp.json
    cat $DASH_FILE >> tmp.json
    echo '}' >> tmp.json
    RESULT=$(cat tmp.json | jq '.dashboard.id = null' | curl -s -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $KEY" "http://$TARGET_HOST/api/dashboards/import" -d @-)
    rm tmp.json
    if [[ "$RESULT" == *"\"imported\":true"* ]]; then
        log_success "$RESULT"
        NUMSUCCESS=$((NUMSUCCESS + 1))
    else
        log_failure "$RESULT"
        NUMFAILURE=$((NUMFAILURE + 1))
    fi
done

log_title "------------ Import complete. $NUMSUCCESS dashboards were successfully imported. $NUMFAILURE dashboard imports failed.------------";
log_title "-------------------------------------------------------- FINISHED ------------------------------------------------------";
