#
# Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
#!/usr/bin/env sh

#!/bin/sh -e -x -u

trap "echo TRAPed signal" HUP INT QUIT KILL TERM

# ---------------------------------------------------------------------------
# Main entry point.
# ---------------------------------------------------------------------------
main()
    {
    SCRIPT_NAME=$(basename "${0}")
    SCRIPT_DIR=$(dirname "$0")
    COMMAND=${1};
    shift
    MAIN_CLASS="com.tangosol.net.DefaultCacheServer"
    MAIN_ARGS=""
    GET_SITE=true

    case "${COMMAND}" in
        server) server ;;
        console) console ;;
        probe) probe ${@} ;;
        queryplus) queryplus ${@} ;;
        mbeanserver) mbeanserver ${@} ;;
        *) usage ;;
    esac

    start
    }


# ---------------------------------------------------------------------------
# Add the configuration for running a DefaultCacheServer
# ---------------------------------------------------------------------------
usage()
    {
    echo "Invalid command '${COMMAND}', must be one of server, console, probe, queryplus, or mbeanserver"
    exit 1
    }

# ---------------------------------------------------------------------------
# Add the configuration for running a cache server
# ---------------------------------------------------------------------------
server()
    {
    MAIN_CLASS="com.oracle.coherence.k8s.Main"

    if [[ -n "${COH_MAIN_CLASS}" ]]
    then
        MAIN_ARGS=${COH_MAIN_CLASS}
    else
        MAIN_ARGS="com.tangosol.net.DefaultCacheServer"
    fi

    echo "Configuring cache server '${MAIN_ARGS}'"

    if [[ -n "${COH_MAIN_ARGS}" ]]
    then
        MAIN_ARGS="${MAIN_ARGS} ${COH_MAIN_ARGS}"
    fi

    CLASSPATH="${CLASSPATH}:${COH_UTIL_DIR}/lib/coherence-utils.jar"

#   We must have this to allow the JMX readiness probe to connect.
    PROPS="${PROPS} -Dcom.sun.management.jmxremote=true -Dcoherence.operator.server=true"

#   Configure the Coherence member's role
    if [[ -n "${COH_ROLE}" ]]
    then
        PROPS="${PROPS} -Dcoherence.role=${COH_ROLE}"
    else
        PROPS="${PROPS} -Dcoherence.role=storage"
    fi

#   Configure whether this member is storage enabled
    if [[ "${COH_STORAGE_ENABLED}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.distributed.localstorage=${COH_STORAGE_ENABLED}"
    fi

#   By default we use the G1 collector
    if [[ "${JVM_GC_OPTS}" == "" ]]
    then
        JVM_GC_OPTS="-XX:+UseG1GC"
    fi

#   If the MAX_HEAP variable is set use it to set the -Xms and -Xmx arguments
    if [[ "${MAX_HEAP}" != "" ]]
    then
        MEM_OPTS="-Xms${MAX_HEAP} -Xmx${MAX_HEAP}"
    fi

#   Configure whether management is added to the classpath
    if [[ "${COH_MGMT_ENABLED}" == "true" ]]
    then
        CLASSPATH="${CLASSPATH}:${COHERENCE_HOME}/lib/coherence-management.jar"
        PROPS="${PROPS} -Dcoherence.management.http=all"
    fi

#   Configure whether metrics is added to the classpath
    if [[ "${COH_METRICS_ENABLED}" == "true" ]]
    then
        CLASSPATH="${CLASSPATH}:${COHERENCE_HOME}/lib/coherence-metrics.jar"
        PROPS="${PROPS} -Dcoherence.metrics.http.enabled=true"
    fi

#   Configure whether to add third-party modules to the classpath
    if [[ "${COH_MGMT_ENABLED}" == "true" || "${COH_METRICS_ENABLED}" == "true" ]]
    then
      if [[ "${DEPENDENCY_MODULES}" != "" ]]
      then
          CLASSPATH="${CLASSPATH}:${DEPENDENCY_MODULES}/*"
      fi
    fi
    }


# ---------------------------------------------------------------------------
# Add the configuration for running a CacheFactory console.
# ---------------------------------------------------------------------------
console()
    {
    echo "Configuring Coherence console"

    APP_TYPE="java"
    COH_ROLE="console"
    MAX_HEAP=""
    CLASSPATH="${CLASSPATH}:${COHERENCE_HOME}/lib/jline.jar$:${COH_UTIL_DIR}/lib/coherence-utils.jar"
    MAIN_CLASS="com.tangosol.net.CacheFactory"
    PROPS="${PROPS} -Dcoherence.distributed.localstorage=false"
    }


# ---------------------------------------------------------------------------
# Add the configuration for running a QueryPlus console.
# ---------------------------------------------------------------------------
queryplus()
    {
    echo "Configuring QueryPlus"

    APP_TYPE="java"
    COH_ROLE="queryPlus"
    MAX_HEAP=""
    CLASSPATH="${CLASSPATH}:${COHERENCE_HOME}/lib/jline.jar:${COH_UTIL_DIR}/lib/coherence-utils.jar"
    MAIN_CLASS="com.tangosol.coherence.dslquery.QueryPlus"
    PROPS="${PROPS} -Dcoherence.distributed.localstorage=false"
    MAIN_ARGS=${@}
    }


# ---------------------------------------------------------------------------
# Add the configuration for running a Coherence MBeanConnector.
# ---------------------------------------------------------------------------
mbeanserver()
    {
    echo "Configuring MBeanConnector Server"

    APP_TYPE="java"
    CLASSPATH="${CLASSPATH}:${COH_UTIL_DIR}/lib/*"
    MAIN_CLASS="com.oracle.coherence.k8s.JmxmpServer"
    MAIN_ARGS=""
    PROPS="${PROPS} -Dcoherence.distributed.localstorage=false \
         -Dcoherence.management.serverfactory=com.oracle.coherence.k8s.JmxmpServer \
         -Dcoherence.jmxmp.port=9099 \
         -Dcoherence.management=all \
         -Dcoherence.management.remote=true \
         -Dcom.sun.management.jmxremote.ssl=false \
         -Dcom.sun.management.jmxremote.authenticate=false"

#   If the MAX_HEAP variable is set use it to set the -Xms and -Xmx arguments
    if [[ "${MAX_HEAP}" != "" ]]
    then
        MEM_OPTS="-Xms${MAX_HEAP} -Xmx${MAX_HEAP}"
    fi

#   Configure the Coherence member's role
    if [[ -n "${COH_ROLE}" ]]
    then
        PROPS="${PROPS} -Dcoherence.role=${COH_ROLE}"
    else
        PROPS="${PROPS} -Dcoherence.role=MBeanServer"
    fi
    }


# ---------------------------------------------------------------------------
# Add the configuration for running a readiness/liveness probe.
# ---------------------------------------------------------------------------
probe()
    {
    echo "Configuring readiness/liveness probe"

    APP_TYPE="java"
    COH_ROLE="probe"
    MAX_HEAP=""
    DEPENDENCY_MODULES=""
    CLASSPATH="${CLASSPATH}:${COH_UTIL_DIR}/lib/*:${JAVA_HOME}/lib/tools.jar"
    MAIN_CLASS=${1}
    shift
    MAIN_ARGS=${@}
    PROPS="${PROPS} -Dcoherence.distributed.localstorage=false"
#   The probe does not require the site
    GET_SITE=""
    }


# ---------------------------------------------------------------------------
# Configure and start the application
# ---------------------------------------------------------------------------
start()
    {
    echo "Starting ${COMMAND}"

#   Set the common configuration
    commonConfiguration

#   Check whether the Coherence version is >= 12.2.1.4.0
    checkVersion "12.2.1.4.0"
    IS_12_2_1_4=$?

    echo "IS_12_2_1_4 = ${IS_12_2_1_4}"

    if [[ ${IS_12_2_1_4} == 0 ]]
    then
        configure_12_2_1_4
    else
        configure_pre_12_2_1_4
    fi

#   If debug is enabled configure the JVM debug arguments
    if [ "${DEBUG_ENABLED}" != "" ]
    then
      if [ "${DEBUG_ATTACH}" == "" ]
      then
#         The debugger should listen for connections on port 5055
        JVM_DEBUG_ARGS="-agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=*:5005"
      else
#         The debugger should connect back to the socket addess in DEBUG_ATTACH
        JVM_DEBUG_ARGS="-agentlib:jdwp=transport=dt_socket,server=n,address=${DEBUG_ATTACH},suspend=n,timeout=5000"
      fi
    else
      JVM_DEBUG_ARGS=""
    fi

#   Create the full classpath to use
    CLASSPATH="${COH_EXTRA_CLASSPATH}:${CLASSPATH}:${COHERENCE_HOME}/conf:${COHERENCE_HOME}/lib/coherence.jar"

#   Create the command line to use to start the JVM
    CMD="-cp ${CLASSPATH} ${JVM_GC_OPTS} ${MEM_OPTS} ${JVM_DEBUG_ARGS} \
        -XX:+HeapDumpOnOutOfMemoryError -XX:+ExitOnOutOfMemoryError \
        -XX:+UnlockDiagnosticVMOptions -XX:+UnlockExperimentalVMOptions \
        -Dcoherence.ttl=0 \
        ${PROPS} ${JVM_ARGS}"

#   Dump the full set of environment variables to the console for logging/debugging
    echo "---------------------------------"
    echo "Environment:"
    echo "---------------------------------"
    env
    echo "---------------------------------"

    if [[ "${COH_APP_DIR}" != "" ]]
    then
      echo "Changing working directory to ${COH_APP_DIR}"
      cd ${COH_APP_DIR}
    fi

    if [[ "${APP_TYPE}" == "" ]]
    then
      APP_TYPE="java"
    fi

    if [[ "${APP_TYPE}" == "java" ]]
    then
      runJava
    else
      runGraal
    fi
    }

# ---------------------------------------------------------------------------
# Executes the command as a plain Java command line.
# ---------------------------------------------------------------------------
runJava()
    {
    CMD="${JAVA_HOME}/bin/java ${CMD} ${MAIN_CLASS} ${MAIN_ARGS}"

    echo "---------------------------------"
    echo "Starting the Coherence ${COMMAND} using:"
    echo "${CMD}"
    echo "---------------------------------"

    exec ${CMD}
    }

# ---------------------------------------------------------------------------
# Executes the command as a Graal command line.
# ---------------------------------------------------------------------------
runGraal()
    {
    CMD=$(echo ${CMD} | sed -e "s/\-cp /--vm.cp /g")
    CMD=$(echo ${CMD} | sed -e "s/\ -D/ --vm.D/g")
    CMD=$(echo ${CMD} | sed -e "s/\ -XX/ --vm.XX/g")
    CMD=$(echo ${CMD} | sed -e "s/\ -Xms/ --vm.Xms/g")
    CMD=$(echo ${CMD} | sed -e "s/\ -Xmx/ --vm.Xmx/g")

    CMD="${APP_TYPE} --polyglot --jvm ${CMD} ${COH_MAIN_CLASS} ${COH_MAIN_ARGS}"

    echo "--------------------------------------------------------------------"
    echo "Starting the Coherence Graal ${APP_TYPE} ${COMMAND} using:"
    echo "${CMD}"
    echo "--------------------------------------------------------------------"

    exec ${CMD}
    }


# ---------------------------------------------------------------------------
# Execute a version check.
# This function takes a single parameter that is a version number and returns
# 0 if the current Coherence version is >= the specified version or 1 if the
# current Coherence version is lower than the specified version.
# ---------------------------------------------------------------------------
checkVersion()
    {
    ${JAVA_HOME}/bin/java -cp ${COH_UTIL_DIR}/lib/*:${COHERENCE_HOME}/lib/coherence.jar \
        com.oracle.coherence.k8s.CoherenceVersion $1

    return $?
    }


# ---------------------------------------------------------------------------
# Add the configuration common to all versions
# ---------------------------------------------------------------------------
commonConfiguration()
    {
#   Configure Coherence WKA
    if [ "${COH_WKA}" != "" ]
    then
       PROPS="${PROPS} -Dcoherence.wka=${COH_WKA}"
    fi

#   Configure the Coherence member properties
    PROPS="${PROPS} -Dcoherence.machine=${COH_MACHINE_NAME}"
    PROPS="${PROPS} -Dcoherence.member=${COH_MEMBER_NAME}"
    PROPS="${PROPS} -Dcoherence.cluster=${COH_CLUSTER_NAME}"

#   Configure the cache configuration file to use
    if [[ -n "${COH_CACHE_CONFIG}" ]]
    then
        PROPS="${PROPS} -Dcoherence.cacheconfig=${COH_CACHE_CONFIG}"
    else
        PROPS="${PROPS} -Dcoherence.cacheconfig=coherence-cache-config.xml"
    fi

#   Configure the port to publish metrics on
    if [[ "${COH_METRICS_PORT}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.metrics.http.port=${COH_METRICS_PORT}"
    fi

#   Configure the port to use for management over rest
    if [[ "${COH_MGMT_HTTP_PORT}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.management.http.port=${COH_MGMT_HTTP_PORT}"
    fi

#   Configure the Coherence member's site and rack
    if [[ "${GET_SITE}" != "" ]]
    then
      if [[ "${COH_SITE_INFO_LOCATION}" != "" ]]
      then
          case "${COH_SITE_INFO_LOCATION}" in
              http://\$*)
                  SITE=""
                  ;;
              http://*)
                  if [[ "${OPERATOR_REQUEST_TIMEOUT}" != "" ]]
                  then
                    TIMEOUT=${OPERATOR_REQUEST_TIMEOUT}
                  else
                    TIMEOUT=120
                  fi

                  SITE=$(curl --silent -m ${TIMEOUT} -X GET ${COH_SITE_INFO_LOCATION})
                  if [[ $? != 0 ]]
                  then
                      SITE=""
                  else
                      echo "Site value: ${SITE}"
                  fi
                  ;;
              *)
                  if [[ -f "${COH_SITE_INFO_LOCATION}" ]]
                  then
                      SITE=`cat ${COH_SITE_INFO_LOCATION}`
                  fi
          esac
      fi

      if [[ "${COH_RACK_INFO_LOCATION}" != "" ]]
      then
          case "${COH_RACK_INFO_LOCATION}" in
              http://\$*)
                  RACK=""
                  ;;
              http://*)
                  if [[ "${OPERATOR_REQUEST_TIMEOUT}" != "" ]]
                  then
                    TIMEOUT=${OPERATOR_REQUEST_TIMEOUT}
                  else
                    TIMEOUT=30
                  fi

                  RACK=$(curl --silent -m ${TIMEOUT} ${COH_RACK_INFO_LOCATION})
                  if [[ $? != 0 ]]
                  then
                      RACK=""
                  else
                      echo "Rack value: ${RACK}"
                  fi
                  ;;
              *)
                  if [[ -f "${COH_RACK_INFO_LOCATION}" ]]
                  then
                      RACK=`cat ${COH_RACK_INFO_LOCATION}`
                  fi
          esac
      fi

      if [[ -n "${SITE}" ]]
      then
          PROPS="${PROPS} -Dcoherence.site=${SITE}"
      fi

      if [[ -n "${RACK}" ]]
      then
          PROPS="${PROPS} -Dcoherence.rack=${RACK}"
      else
          if [[ -n "${SITE}" ]]
          then
              PROPS="${PROPS} -Dcoherence.rack=${SITE}"
          fi
      fi
    fi

#   Configure Coherence persistence
    if [[ "${COH_PERSISTENCE_ENABLED}" == "true" ]]
    then
        PROPS="${PROPS} -Dcoherence.distributed.persistence-mode=active -Dcoherence.distributed.persistence.base.dir=/persistence"
    else
        PROPS="${PROPS} -Dcoherence.distributed.persistence-mode=on-demand"
    fi

#   Configure the Coherence snapshot location
    if [[ "${COH_SNAPSHOT_ENABLED}" == "true" ]]
    then
        PROPS="${PROPS} -Dcoherence.distributed.persistence.snapshot.dir=/snapshot"
    fi

#   Configure logging
    if [[ "${COH_LOGGING_CONFIG}" != "" ]]
    then
        if [[ -f "${COH_LOGGING_CONFIG}" ]]
        then
        PROPS="${PROPS} -Dcoherence.log=jdk -Dcoherence.log.logger=com.oracle.coherence \
               -Djava.util.logging.config.file=${COH_LOGGING_CONFIG}"
        else
            echo "Logging configuration file ${COH_LOGGING_CONFIG} does not exists."
        fi
    fi

#   Configure the Coherence log level
    if [[ "${COH_LOG_LEVEL}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.log.level=${COH_LOG_LEVEL}"
    fi
    }


# ---------------------------------------------------------------------------
# Add configuration applicable to versions earlier than 12.2.1.4
# ---------------------------------------------------------------------------
configure_pre_12_2_1_4()
    {
    echo "Adding configuration for pre-12.2.1.4.0 version"

    #   This is pre Coherence 12.2.1.4.0
    #   copy our AddressProvider overrides file onto the classpath
    cp ${SCRIPT_DIR}/k8s-coherence-nossl-override.xml ${COHERENCE_HOME}/conf/k8s-coherence-nossl-override.xml

    #   use our AddressProvider overrides file
    PROPS="${PROPS} -Dcoherence.override=k8s-coherence-nossl-override.xml"

    if [[ -n "${COH_OVERRIDE_CONFIG}" ]]
    then
        PROPS="${PROPS} -Dcoherence.k8s.override=${COH_OVERRIDE_CONFIG}"
    fi
    }


# ---------------------------------------------------------------------------
# Add configuration applicable to 12.2.1.4 and above
# ---------------------------------------------------------------------------
configure_12_2_1_4()
    {
    echo "Adding configuration for 12.2.1.4.0 and above"

#   This is Coherence 12.2.1.4.0 or above so we support SSL management and metrics
#   copy our AddressProvider and SSL overrides file onto the classpath
    mkdir -p ${COHERENCE_HOME}/conf || true
    cp ${SCRIPT_DIR}/k8s-coherence-override.xml ${COHERENCE_HOME}/conf/k8s-coherence-override.xml

#   use our AP and SSL overrides file
    PROPS="${PROPS} -Dcoherence.override=k8s-coherence-override.xml"

    if [[ -n "${COH_OVERRIDE_CONFIG}" ]]
    then
        PROPS="${PROPS} -Dcoherence.k8s.override=${COH_OVERRIDE_CONFIG}"
    fi

    addManagementSSL
    addMetricsSSL
    }


# ---------------------------------------------------------------------------
# Add the configuration properties to enable SSL
# on the management over ReST endpoint.
# ---------------------------------------------------------------------------
addManagementSSL()
    {
#   Configure SSL for the management endpoint
    if [[ "${COH_MGMT_SSL_CERTS}" != "" ]]
    then
        if [[ "${COH_MGMT_SSL_CERTS:${#COH_MGMT_SSL_CERTS}-1}" != "/" ]]
        then
            COH_MGMT_SSL_CERTS="${COH_MGMT_SSL_CERTS}/"
        fi

        if [[ "${COH_MGMT_SSL_CERTS:1:5}" == "file:" ]]
        then
            COH_MGMT_SSL_URL_PREFIX="${COH_MGMT_SSL_CERTS}"
        else
            COH_MGMT_SSL_URL_PREFIX="file:${COH_MGMT_SSL_CERTS}"
        fi
    else
        COH_MGMT_SSL_URL_PREFIX="file:"
    fi

    if [[ "${COH_MGMT_SSL_ENABLED}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.management.http.provider=ManagementSSLProvider"
    fi
    if [[ "${COH_MGMT_SSL_KEYSTORE}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.management.security.keystore=${COH_MGMT_SSL_URL_PREFIX}${COH_MGMT_SSL_KEYSTORE}"
    fi
    if [[ "${COH_MGMT_SSL_KEYSTORE_PASSWORD_FILE}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.management.security.keystore.password=${COH_MGMT_SSL_URL_PREFIX}${COH_MGMT_SSL_KEYSTORE_PASSWORD_FILE}"
    fi
    if [[ "${COH_MGMT_SSL_KEY_PASSWORD_FILE}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.management.security.key.password=${COH_MGMT_SSL_URL_PREFIX}${COH_MGMT_SSL_KEY_PASSWORD_FILE}"
    fi
    if [[ "${COH_MGMT_SSL_KEYSTORE_ALGORITHM}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.management.security.keystore.algorithm=${COH_MGMT_SSL_KEYSTORE_ALGORITHM}"
    fi
    if [[ "${COH_MGMT_SSL_KEYSTORE_PROVIDER}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.management.security.keystore.provider=${COH_MGMT_SSL_KEYSTORE_PROVIDER}"
    fi
    if [[ "${COH_MGMT_SSL_KEYSTORE_TYPE}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.management.security.keystore.type=${COH_MGMT_SSL_KEYSTORE_TYPE}"
    fi
    if [[ "${COH_MGMT_SSL_TRUSTSTORE}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.management.security.truststore=${COH_MGMT_SSL_URL_PREFIX}${COH_MGMT_SSL_TRUSTSTORE}"
    fi
    if [[ "${COH_MGMT_SSL_TRUSTSTORE_PASSWORD_FILE}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.management.security.truststore.password=${COH_MGMT_SSL_URL_PREFIX}${COH_MGMT_SSL_TRUSTSTORE_PASSWORD_FILE}"
    fi
    if [[ "${COH_MGMT_SSL_TRUSTSTORE_ALGORITHM}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.management.security.truststore.algorithm=${COH_MGMT_SSL_TRUSTSTORE_ALGORITHM}"
    fi
    if [[ "${COH_MGMT_SSL_TRUSTSTORE_PROVIDER}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.management.security.truststore.provider=${COH_MGMT_SSL_TRUSTSTORE_PROVIDER}"
    fi
    if [[ "${COH_MGMT_SSL_TRUSTSTORE_TYPE}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.management.security.truststore.type=${COH_MGMT_SSL_TRUSTSTORE_TYPE}"
    fi
    if [[ "${COH_MGMT_SSL_REQUIRE_CLIENT_CERT}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.management.http.auth=cert"
    fi

#   Dump the contents of the Management certs folder for debugging
    if [[ "${COH_MGMT_SSL_CERTS}" != "" ]]
    then
        echo "-----------------------------------------------------------------------------"
        echo "Contents of Coherence Management Certs folder ${COH_MGMT_SSL_CERTS}"
        echo "-----------------------------------------------------------------------------"
        ls -alr ${COH_MGMT_SSL_CERTS}
        echo "-----------------------------------------------------------------------------"
    fi
    }


# ---------------------------------------------------------------------------
# Add the configuration properties to enable SSL
# on the metrics endpoint.
# ---------------------------------------------------------------------------
addMetricsSSL()
    {
#   Configure SSL for the metrics endpoint
    if [[ "${COH_METRICS_SSL_CERTS}" != "" ]]
    then
        if [[ "${COH_METRICS_SSL_CERTS:${#COH_METRICS_SSL_CERTS}-1}" != "/" ]]
        then
            COH_METRICS_SSL_CERTS="${COH_METRICS_SSL_CERTS}/"
        fi

        if [[ "${COH_METRICS_SSL_CERTS:1:5}" == "file:" ]]
        then
            COH_METRICS_SSL_URL_PREFIX="${COH_METRICS_SSL_CERTS}"
        else
            COH_METRICS_SSL_URL_PREFIX="file:${COH_METRICS_SSL_CERTS}"
        fi
    else
        COH_METRICS_SSL_URL_PREFIX="file:"
    fi

    if [[ "${COH_METRICS_SSL_ENABLED}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.metrics.http.provider=MetricsSSLProvider"
    fi
    if [[ "${COH_METRICS_SSL_KEYSTORE}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.metrics.security.keystore=${COH_METRICS_SSL_URL_PREFIX}${COH_METRICS_SSL_KEYSTORE}"
    fi
    if [[ "${COH_METRICS_SSL_KEYSTORE_PASSWORD_FILE}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.metrics.security.keystore.password=${COH_METRICS_SSL_URL_PREFIX}${COH_METRICS_SSL_KEYSTORE_PASSWORD_FILE}"
    fi
    if [[ "${COH_METRICS_SSL_KEY_PASSWORD_FILE}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.metrics.security.key.password=${COH_METRICS_SSL_URL_PREFIX}${COH_METRICS_SSL_KEY_PASSWORD_FILE}"
    fi
    if [[ "${COH_METRICS_SSL_KEYSTORE_ALGORITHM}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.metrics.security.keystore.algorithm=${COH_METRICS_SSL_KEYSTORE_ALGORITHM}"
    fi
    if [[ "${COH_METRICS_SSL_KEYSTORE_PROVIDER}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.metrics.security.keystore.provider=${COH_METRICS_SSL_KEYSTORE_PROVIDER}"
    fi
    if [[ "${COH_METRICS_SSL_KEYSTORE_TYPE}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.metrics.security.keystore.type=${COH_METRICS_SSL_KEYSTORE_TYPE}"
    fi
    if [[ "${COH_METRICS_SSL_TRUSTSTORE}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.metrics.security.truststore=${COH_METRICS_SSL_URL_PREFIX}${COH_METRICS_SSL_TRUSTSTORE}"
    fi
    if [[ "${COH_METRICS_SSL_TRUSTSTORE_PASSWORD_FILE}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.metrics.security.truststore.password=${COH_METRICS_SSL_URL_PREFIX}${COH_METRICS_SSL_TRUSTSTORE_PASSWORD_FILE}"
    fi
    if [[ "${COH_METRICS_SSL_TRUSTSTORE_ALGORITHM}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.metrics.security.truststore.algorithm=${COH_METRICS_SSL_TRUSTSTORE_ALGORITHM}"
    fi
    if [[ "${COH_METRICS_SSL_TRUSTSTORE_PROVIDER}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.metrics.security.truststore.provider=${COH_METRICS_SSL_TRUSTSTORE_PROVIDER}"
    fi
    if [[ "${COH_METRICS_SSL_TRUSTSTORE_TYPE}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.metrics.security.truststore.type=${COH_METRICS_SSL_TRUSTSTORE_TYPE}"
    fi
    if [[ "${COH_METRICS_SSL_REQUIRE_CLIENT_CERT}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.metrics.http.auth=cert"
    fi

#   Dump the contents of the Metrics certs folder for debugging
    if [[ "${COH_METRICS_SSL_CERTS}" != "" ]]
    then
        echo "-----------------------------------------------------------------------------"
        echo "Contents of Coherence Metrics Certs folder ${COH_METRICS_SSL_CERTS}"
        echo "-----------------------------------------------------------------------------"
        ls -alr ${COH_METRICS_SSL_CERTS}
        echo "-----------------------------------------------------------------------------"
    fi
    }

# ---------------------------------------------------------------------------
main "$@"
