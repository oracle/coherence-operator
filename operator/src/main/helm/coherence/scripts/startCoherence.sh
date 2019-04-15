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
    COMMAND=${1};
    shift
    MAIN_CLASS="com.tangosol.net.DefaultCacheServer"
    COH_MAIN_ARGS=""

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
    echo "Invalid command '${COMMAND}', must be one of server, console, probe or queryplus"
    exit 1
    }

# ---------------------------------------------------------------------------
# Add the configuration for running a DefaultCacheServer
# ---------------------------------------------------------------------------
server()
    {
    echo "Configuring DefaultCacheServer"

    MAIN_CLASS="com.tangosol.net.DefaultCacheServer"
    CLASSPATH="${CLASSPATH}:${COH_UTIL_DIR}/lib/coherence-utils.jar"

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
    if [[ "${JVM_ARGS}" == "" ]]
    then
        JVM_ARGS="-XX:+UseG1GC"
    fi

#   If the MAX_HEAP variable is set use it to set the -Xms and -Xmx arguments
    if [[ "${MAX_HEAP}" != "" ]]
    then
        MEM_OPTS="-Xms${MAX_HEAP} -Xmx${MAX_HEAP}"
    fi

#   Configure whether coherence-rest is added to the classpath
    if [[ "${COH_USE_REST}" != "" ]]
    then
        CLASSPATH="${CLASSPATH}:${COHERENCE_HOME}/lib/coherence-rest.jar"
        PROPS="${PROPS} -Dcoherence.management.http=all"
    fi

#   Configure whether to add third-party modules to the classpath
    if [[ "${DEPENDENCY_MODULES}" != "" ]]
    then
        CLASSPATH="${CLASSPATH}:${DEPENDENCY_MODULES}/*"
    fi
    }


# ---------------------------------------------------------------------------
# Add the configuration for running a CacheFactory console.
# ---------------------------------------------------------------------------
console()
    {
    echo "Configuring Coherence console"

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

    COH_ROLE="queryPlus"
    MAX_HEAP=""
    CLASSPATH="${CLASSPATH}:${COHERENCE_HOME}/lib/jline.jar:${COH_UTIL_DIR}/lib/coherence-utils.jar"
    MAIN_CLASS="com.tangosol.coherence.dslquery.QueryPlus"
    PROPS="${PROPS} -Dcoherence.distributed.localstorage=false"
    COH_MAIN_ARGS=${@}
    }


# ---------------------------------------------------------------------------
# Add the configuration for running a Coherence MBeanConnector.
# ---------------------------------------------------------------------------
mbeanserver()
    {
    echo "Configuring MBeanConnector Server"

    CLASSPATH="${CLASSPATH}:${COH_UTIL_DIR}/lib/*"
    MAIN_CLASS="com.oracle.coherence.k8s.JmxmpServer"
    COH_MAIN_ARGS=""
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

    COH_ROLE="probe"
    MAX_HEAP=""
    DEPENDENCY_MODULES=""
    CLASSPATH="${CLASSPATH}:${COH_UTIL_DIR}/lib/*"
    MAIN_CLASS=${1}
    shift
    COH_MAIN_ARGS=${@}
    PROPS="${PROPS} -Dcoherence.distributed.localstorage=false"
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

echo "IS_12_2_1_4 ${IS_12_2_1_4}"

    if [[ ${IS_12_2_1_4} == 0 ]]
    then
        # when 12.2.1.4 released, enable following line
        # configure_12_2_1_4
        # and delete next line
        configure_pre_12_2_1_4
    else
        configure_pre_12_2_1_4
    fi

#   Create the full classpath to use
    CLASSPATH="${COH_EXTRA_CLASSPATH}:${CLASSPATH}:${COHERENCE_HOME}/conf:${COHERENCE_HOME}/lib/coherence.jar"

#   Create the command line to use to start the JVM
    CMD="${JAVA_HOME}/bin/java -cp ${CLASSPATH} ${JVM_ARGS} ${MEM_OPTS} \
        -XX:+HeapDumpOnOutOfMemoryError -XX:+ExitOnOutOfMemoryError \
        -XX:+UnlockDiagnosticVMOptions -XX:+UnlockExperimentalVMOptions \
        -Dcoherence.ttl=0 \
        ${PROPS} ${JAVA_OPTS} ${MAIN_CLASS} ${COH_MAIN_ARGS}"

#   Dump the full set of environment variables to the console for logging/debugging
    echo "---------------------------------"
    echo "Environment:"
    echo "---------------------------------"
    env
    echo "---------------------------------"

    echo "---------------------------------"
    echo "Starting the Coherence ${COMMAND} using:"
    echo "${CMD}"
    echo "---------------------------------"

#   Start the JVM
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

#   Configure the extend port property (and if a cache config is not set then set the extend cache configuration)
    if [[ "${COH_EXTEND_PORT}" != "" ]]
    then
        PROPS="${PROPS} -Dcoherence.extend.port=${COH_EXTEND_PORT}"
        if [[ -z "${COH_CACHE_CONFIG}" ]]
        then
            PROPS="${PROPS} -Dcoherence.cacheconfig=extend-cache-config.xml"
        fi
    fi

#   Configure the cache configuration file to use
    if [[ -n "${COH_CACHE_CONFIG}" ]]
    then
        PROPS="${PROPS} -Dcoherence.cacheconfig=${COH_CACHE_CONFIG}"
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
    if [[ "${COH_SITE_INFO_LOCATION}" != "" ]]
    then
        case "${COH_SITE_INFO_LOCATION}" in
            http://\$*)
                SITE=""
                break;;
            http://*)
                SITE=$(curl ${COH_SITE_INFO_LOCATION})
                if [[ $? != 0 ]]
                then
                    SITE=""
                fi
                break;;
            *)
                if [[ -f "${COH_SITE_INFO_LOCATION}" ]]
                then
                    SITE=`cat ${COH_SITE_INFO_LOCATION}`
                fi
        esac

        if [[ -n "${SITE}" ]]
        then
            PROPS="${PROPS} -Dcoherence.site=${SITE} -Dcoherence.rack=${SITE}"
        fi
    fi

#   Configure the POF configuration file to use
    if [[ -n "${COH_POF_CONFIG}" ]]
    then
        PROPS="${PROPS} -Dcoherence.pof.config=${COH_POF_CONFIG}"
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
    cp /scripts/k8s-coherence-nossl-override.xml ${COHERENCE_HOME}/conf/k8s-coherence-nossl-override.xml

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
    cp /scripts/k8s-coherence-override.xml ${COHERENCE_HOME}/conf/k8s-coherence-override.xml

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
