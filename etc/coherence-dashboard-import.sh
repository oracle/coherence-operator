#!/usr/bin/env bash
#
# kibana dashboard import script
#

cd /usr/share/kibana/data/coherence/dashboards

echo "Waiting up to 60 seconds for Kibana to get in green overall state..."
for i in {1..60}; do
  curl -s localhost:5601/api/status | python -c 'import sys, json; print json.load(sys.stdin)["status"]["overall"]["state"]' 2> /dev/null | grep green > /dev/null && break || sleep 1
done

for DASHBOARD_FILE in *; do
  echo -e "Importing ${DASHBOARD_FILE} dashboard..."
  if ! python -c 'import sys, json; print json.load(sys.stdin)' < "${DASHBOARD_FILE}" &> /dev/null ; then
    echo "${DASHBOARD_FILE} is not valid JSON, assuming it's an URL..."
    TMP_FILE="$(mktemp)"
    curl -s $(cat ${DASHBOARD_FILE}) > ${TMP_FILE}
    curl -v -s --connect-timeout 60 --max-time 60 -XPOST localhost:5601/api/kibana/dashboards/import?force=true -H 'kbn-xsrf:true' -H 'Content-type:application/json' -d @${TMP_FILE}
    rm ${TMP_FILE}
  else
    echo "Valid JSON found in ${DASHBOARD_FILE}, importing..."
    TMP_FILE="$(mktemp)"
    # Following to allow import via API of exported data from "saved" objects - https://discuss.elastic.co/t/saved-objects-api-use-example/168742/4
    echo "{ \"version\": \"6.5.4\", \"objects\": " > ${TMP_FILE}
    cat ${DASHBOARD_FILE} | sed -e 's/"_id":/"id":/g' -e 's/"_type":/"type":/g' -e 's/"_source":/"attributes":/g' >> ${TMP_FILE}
    echo "}" >> ${TMP_FILE}
    curl -v -s --connect-timeout 60 --max-time 60 -XPOST localhost:5601/api/kibana/dashboards/import?force=true -H 'kbn-xsrf:true' -H 'Content-type:application/json' -d @${TMP_FILE}
    rm ${TMP_FILE}
    # Set the default index Pattern
    curl -v -s --connect-timeout 60 --max-time 60 -XPOST localhost:5601/api/kibana/settings/defaultIndex -H 'kbn-xsrf:true' -H 'Content-type:application/json' -d '{"value": "'$DEFAULT_INDEX_PATTERN'"}'

  fi
  if [ "$?" != "0" ]; then
    echo -e "\nImport of ${DASHBOARD_FILE} dashboard failed... Exiting..."
    exit 1
  else
    echo -e "\nImport of ${DASHBOARD_FILE} dashboard finished :-)"
  fi
done
