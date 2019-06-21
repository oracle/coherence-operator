/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s.operator;

import io.kubernetes.client.ApiCallback;
import io.kubernetes.client.ApiClient;
import io.kubernetes.client.ApiException;
import io.kubernetes.client.Configuration;
import io.kubernetes.client.apis.CoreV1Api;
import io.kubernetes.client.models.V1Namespace;
import io.kubernetes.client.models.V1ObjectMeta;
import io.kubernetes.client.models.V1Secret;
import io.kubernetes.client.util.Config;
import io.kubernetes.client.util.Watch;

import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.concurrent.Executors;
import java.util.concurrent.ThreadFactory;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.function.Consumer;
import java.util.logging.Level;
import java.util.logging.Logger;

/**
 * A Kubernetes Operator for Coherence.
 *
 * @author sc
 */
public class CoherenceOperator
    {
    /**
     * Entry point of Coherence operator.
     *
     * @param args none, ignored
     */
    public static void main(String[] args)
        {
        String sNamespace = System.getenv(OPERATOR_NAMESPACE);
        if (sNamespace == null)
            {
            sNamespace = DEFAULT_OPERATOR_NAMESPACE;
            }

        String sTargetNamespaces = System.getenv(TARGET_NAMESPACES);

        String[] asNamespaces = (sTargetNamespaces == null || sTargetNamespaces.isEmpty()) ?
                new String[] { null } : sTargetNamespaces.trim().split("\\s*,\\s*");

        String sExcludedNamespaces = System.getenv(EXCLUDED_NAMESPACES);

        String[] asExcludedNamespaces = (sExcludedNamespaces == null || sExcludedNamespaces.isEmpty()) ?
                DEFAULT_EXCLUDED_NAMESPACES : sExcludedNamespaces.trim().split("\\s*,\\s*");

        assertNamespaces(asNamespaces, asExcludedNamespaces);

        List<AbstractWatcher<?, ?>> listWatchers = new ArrayList<>();

        try
            {
            // set Kubernetes default Api Client
            ApiClient client = Config.defaultClient();
            Configuration.setDefaultApiClient(client);

            new KubernetesInfoServer(K8S_INFO_SERVER_PORT).start();

            if (Boolean.parseBoolean(System.getenv(EFK_INTEGRATION_ENABLED)))
                {
                AbstractWatcher<?, ?> namespaceWatcher = createNamespaceWatcher(sNamespace, asNamespaces,
                        asExcludedNamespaces, fStopping);
                listWatchers.add(namespaceWatcher);
                namespaceWatcher.start(threadFactory);
                }
            }
        catch(Throwable t)
            {
            fStopping.set(false);

            LOGGER.log(Level.SEVERE, "Cannot start watchers: " + t, t);
            }

        LOGGER.info("Started Coherence Kubernetes Operator ...");
        // wait for end
        listWatchers.stream().forEach(w -> w.waitForDeath());
        }

    // ---- helper methods ---------------------------------------------------

    /**
     * Creates a NamespaceWatcher with namespaces and AtomicBoolean indicating stopping.
     *
     * @param sNamespace            the namespace of the operator
     * @param asNamespaces          an array of Namespace
     * @param asExcludedNamespaces  an array of excluded Namspace
     * @param fStop                 the AtomicBoolean indicating stopping
     *
     * @return a DeploymentWatcher
     *
     * @exception IllegalAccessException  if the class or its nullary constructor is not accessible
     * @exception InstantiationException  if this Class represents an abstract class, an interface, an array class,
     *     a primitive type, or void; or if the class has no nullary constructor; or if the instantiation fails for
     *     some other reason
     */
    private static NamespaceWatcher createNamespaceWatcher(String sNamespace, String[] asNamespaces,
                                                           String[] asExcludedNamespaces, AtomicBoolean fStop)
            throws IllegalAccessException, InstantiationException
        {
        return new NamespaceWatcher(fStop, new NamespaceProcessor(sNamespace, asNamespaces, asExcludedNamespaces));
        }


    /**
     * Assert that the included namespaces should not be excluded.
     *
     * @param asIncludedNamespaces  namespaces to be included
     * @param asExcludedNamespaces  namespaces to be excluded
     */
    static void assertNamespaces(String[] asIncludedNamespaces, String[] asExcludedNamespaces)
        {
        Set<String> setExcludedNamespaces = new HashSet<>();
            for (String sNamesp : asExcludedNamespaces)
            {
                setExcludedNamespaces.add(sNamesp);
            }

        for (String includedNamespace : asIncludedNamespaces)
            {
            if (setExcludedNamespaces.contains(includedNamespace))
                {
                throw new IllegalArgumentException(includedNamespace + " is in the excluded namespace list: " + setExcludedNamespaces);
                }
            }
        }

    /**
     * The Namespace consumer of Watch.Response&lt;V1Namespace&gt; for processing Kubernetes namespace events.
     */
    static class NamespaceProcessor implements Consumer<Watch.Response<V1Namespace>>
        {
        /**
         * Construct a Namesprocessor for processing Namespace Watch events.
         *
         * @param sNamespace            the namespace of the operator
         * @param asIncludedNamespaces  an array of namespaces included for processing
         * @param asExcludedNamespaces  an array of namespaces excluded for processing
         */
        NamespaceProcessor(String sNamespace, String[] asIncludedNamespaces, String[] asExcludedNamespaces)
            {
            m_sElasticsearchHost = Env.get(ES_HOST, "elasticsearch." + sNamespace + ".svc.cluster.local");
            m_sElasticsearchPort = Env.get(ES_PORT, DEFAULT_ES_PORT);

            assertInteger(m_sElasticsearchPort);

            m_sOperatorHost      = "coherence-operator-service." + sNamespace + ".svc.cluster.local";

            for (String sNamesp : asExcludedNamespaces)
                {
                m_setExcludedNamespaces.add(sNamesp);
                }

            // null means all which corresponds to empty set here
            if (asIncludedNamespaces.length > 0 && asIncludedNamespaces[0] != null)
                {
                for (String sNamesp : asIncludedNamespaces)
                    {
                    m_setIncludedNamespaces.add(sNamesp);
                    }
                }
            }

        /**
         * Consume the item by creating the Coherence internal ConfigMap in each appropriate namespace.
         *
         * @param item  the input argument
         */
        @Override
        public void accept(Watch.Response<V1Namespace> item)
            {
            if (item.object != null)
                {
                final V1ObjectMeta objectMeta = item.object.getMetadata();

                if (objectMeta != null && "ADDED".equals(item.type))
                    {
                    String sNamesp = objectMeta.getName();

                    if (isAcceptableNamespace(sNamesp))
                        {
                        try {
                            V1Secret secret      = null;
                            boolean  secretExist = false;

                            try
                                {
                                secret      = m_coreV1Api.readNamespacedSecret(COHERENCE_MONITORING_CONFIG,
                                                  sNamesp, null, Boolean.TRUE, Boolean.TRUE);
                                secretExist = (secret != null);
                                }
                            catch(Throwable ignore)
                                {
                                // Not Found
                                }

                            if (!secretExist)
                                {
                                V1ObjectMeta secretMeta = new V1ObjectMeta();
                                secretMeta.setName(COHERENCE_MONITORING_CONFIG);

                                secret = new V1Secret().metadata(secretMeta);
                                }

                            secret.putStringDataItem("operatorhost", m_sOperatorHost);
                            if (m_sElasticsearchHost != null)
                                {
                                secret.putStringDataItem("elasticsearchhost", m_sElasticsearchHost);
                                }
                            if (m_sElasticsearchPort != null)
                                {
                                secret.putStringDataItem("elasticsearchport", m_sElasticsearchPort);
                                }

                            if (secretExist)
                                {
                                m_coreV1Api.replaceNamespacedSecret(COHERENCE_MONITORING_CONFIG, sNamesp, secret, null);
                                LOGGER.info("Updated '" + COHERENCE_MONITORING_CONFIG +
                                        "' Secret in namespace[" + sNamesp + "]");
                                }
                            else
                                {
                                ApiCallback<V1Secret> callback = new ApiCallback<>() {
                                    @Override
                                    public void onFailure(ApiException e, int i, Map<String, List<String>> map)
                                    {
                                        LOGGER.warning("Failed in creating '" + COHERENCE_MONITORING_CONFIG +
                                                "' Secret in namespace[" + sNamesp + "]" + e.toString());
                                    }

                                    @Override
                                    public void onSuccess(V1Secret v1Secret, int i, Map<String, List<String>> map)
                                    {
                                        LOGGER.info("Created '" + COHERENCE_MONITORING_CONFIG +
                                                "' Secret in namespace[" + sNamesp + "]");
                                    }

                                    @Override
                                    public void onUploadProgress(long l, long l1, boolean b)
                                    {
                                    }

                                    @Override
                                    public void onDownloadProgress(long l, long l1, boolean b)
                                    {
                                    }
                                };
                                m_coreV1Api.createNamespacedSecretAsync(sNamesp, secret, null, callback);
                                }
                        }
                        catch(Throwable t)
                            {
                            LOGGER.warning("Exception in creating secret in namespace[" + sNamesp + "]: " + t);
                            }
                        }
                    }
                }
            }

        // ----- helpers  ----------------------------------------------------

        /**
         * The Setter of CoreV1Api.
         * @param coreV1Api
         */
        void setCoreV1Api(CoreV1Api coreV1Api)
            {
            m_coreV1Api = coreV1Api;
            }

        /**
         * Returns a boolean indicating whether the given namespace is not in excluded namespace list,
         *
         * @param sNamespace  the namepace
         * @return true if the namespace is accepted for processing.
         */
        private boolean isAcceptableNamespace(String sNamespace)
            {
            return !m_setExcludedNamespaces.contains(sNamespace) &&
                    (m_setIncludedNamespaces.size() == 0 || m_setIncludedNamespaces.contains(sNamespace));
            }

        /**
         * Assert whether the value is an integer.
         *
         * @param value  the value
         */
        private void assertInteger(String value) {
            Integer.parseInt(value);
        }

        // ----- constants ---------------------------------------------------

        /**
         * The Elasticsearch host.
         */
        private final String m_sElasticsearchHost;

        /**
         * The Elasticsearch port.
         */
        private final String m_sElasticsearchPort;

        /**
         * The Operator host.
         */
        private final String m_sOperatorHost;

        /**
         * The CoreV1Api.
         */
        private CoreV1Api m_coreV1Api = new CoreV1Api();

        /**
         * The set of excluded namespaces.
         */
        private final Set<String> m_setExcludedNamespaces = new HashSet<>();

        /**
         * The set of included namespaces.
         */
        private final Set<String> m_setIncludedNamespaces = new HashSet<>();
        }

    // ----- constants -------------------------------------------------------

    /**
     * The environment property name for Coherence operator namespace.
     */
    private static final String OPERATOR_NAMESPACE = "OPERATOR_NAMESPACE";

    /**
     * The default of Coherence operator namespace environment property.
     */
    private static final String DEFAULT_OPERATOR_NAMESPACE = "default";

    /**
     * The environment property name for Coherence operator target namespaces.
     */
    private static final String TARGET_NAMESPACES = "TARGET_NAMESPACES";

    /**
     * The environment property name for Coherence operator excluded namespaces.
     */
    private static final String EXCLUDED_NAMESPACES = "EXCLUDED_NAMESPACES";

    /**
     * The default of Coherence operator excluded namespaces.
     */
    private static final String[] DEFAULT_EXCLUDED_NAMESPACES = new String[] {"docker", "kube-public", "kube-system"};

    /**
     * The environment property name for EFK integration enabled.
     */
    private static final String EFK_INTEGRATION_ENABLED = "EFK_INTEGRATION_ENABLED";

    /**
     * The environment property name for Elasticsearch host.
     */
    private static final String ES_HOST = "ES_HOST";

    /**
     * The environment property name for Elasticsearch port.
     */
    private static final String ES_PORT = "ES_PORT";

    /**
     * The default of Elasticsearch port, 9200.
     */
    private static final String DEFAULT_ES_PORT = "9200";

    private static final int K8S_INFO_SERVER_PORT = 8000;

    /**
     * The name of the Coherence monitoring config secret created by operator.
     */
    private static final String COHERENCE_MONITORING_CONFIG = "coherence-monitoring-config";

    // ----- data members ----------------------------------------------------

    /**
     * The default ThreadFactory.
     */
    private static final ThreadFactory defaultThreadFactory = Executors.defaultThreadFactory();

    /**
     * ThreadFactory with daemon threads.
     */
    private static final ThreadFactory threadFactory = (r) ->
        {
        Thread t = defaultThreadFactory.newThread(r);
        if (!t.isDaemon())
            {
            t.setDaemon(true);
            }
        return t;
        };

    /**
     * The AtomicBoolean indicates whether the operator is stopping.
     */
    private static final AtomicBoolean fStopping = new AtomicBoolean();

    /**
     * Class Logger.
     */
    private static final Logger LOGGER = Logger.getLogger("Operator");
    }
