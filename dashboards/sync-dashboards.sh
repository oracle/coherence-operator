#!/bin/bash

#
# Copyright (c) 2021, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

# This script synchronizes the dashboards in the grafana-micrometer
# with the dashboard in the grafana directory.
#
# The process for modifying any Grafana dahsboard should be done in the grafana directory
# and then this script should be run to update the grafana-micrometer with
# the dashboards.
#
# THe dashboards in the grafana-microprofile directory are significantly different and need
# to be modified separated
#
# Note: The dashboards in the grafana-micrometer overwritten.

DB_BASE=grafana
DB_MM=grafana-micrometer

if [ ! -d $DB_BASE -o ! -d $DB_MM ] ; then
    echo "You must run this script from the dashboards directory"
    exit 1
fi

echo -n "Are you sure you want to contine? (y/n) "
read ans

if [ "$ans" != "y" ] ; then
    echo "No dashboards were changed"
    exit 2
fi

DIR=`pwd`
cd $DB_BASE

for file in *.json
do
    echo "Processing dashboard $file ..."

    # Micrometer changes are simple
    cat $file | sed 's/vendor:coherence_/coherence_/g' > $DIR/$DB_MM/$file
done
