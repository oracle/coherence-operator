/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import com.oracle.bedrock.runtime.LocalPlatform;
import com.tangosol.util.Resources;
import org.junit.Before;
import org.junit.BeforeClass;
import org.junit.ClassRule;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.ExpectedException;

import javax.ws.rs.core.MediaType;
import java.net.URL;

import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.collection.IsEmptyIterable.emptyIterable;
import static org.hamcrest.collection.IsIterableContainingInAnyOrder.containsInAnyOrder;

/**
 * @author jk  2019.02.23
 */
public class RestProbeIT
    {
    @BeforeClass
    public static void setupClass()
        {
        s_client = new ProbeHttpClient(s_httpServer.getBoundPort());
        }

    @Before
    public void setup()
        {
        s_httpServer.reset()
                .onGet(RestProbe.PATH_CLUSTER, URL_CLUSTER, MediaType.APPLICATION_JSON);
        }

    @Test
    public void shouldBeAvailable()
        {
        RestProbe probe = new RestProbe(s_client);

        assertThat(probe.isAvailable(), is(true));
        }

    @Test
    public void shouldNotBeAvailableIfStatusCodeIsNot200()
        {
        s_httpServer.reset();
        s_httpServer.onGet(RestProbe.PATH_CLUSTER, p -> HttpServerStub.NOT_FOUND);

        RestProbe probe = new RestProbe(s_client);

        assertThat(probe.isAvailable(), is(false));
        }

    @Test
    public void shouldNotBeAvailableIfNoManagementServerListening()
        {
        int nPort = LocalPlatform.get().getAvailablePorts().next();

        try (ProbeHttpClient client = new ProbeHttpClient(nPort))
            {
            RestProbe probe = new RestProbe(client);

            assertThat(probe.isAvailable(), is(false));
            }
        }

    @Test
    public void shouldBeClusterMemberIfSingleMember()
        {
        s_httpServer.onGet(RestProbe.PATH_MEMBERS, URL_CLUSTER_MEMBERS_SINGLE, MediaType.APPLICATION_JSON);

        RestProbe probe = new RestProbe(s_client);

        assertThat(probe.isClusterMember(), is(true));
        }

    @Test
    public void shouldBeClusterMemberIfMultipleMembers()
        {
        s_httpServer.onGet(RestProbe.PATH_MEMBERS, URL_CLUSTER_MEMBERS_MULTI, MediaType.APPLICATION_JSON);

        RestProbe probe = new RestProbe(s_client);

        assertThat(probe.isClusterMember(), is(true));
        }

    @Test
    public void shouldBeNotClusterMemberIfZeroMembers()
        {
        s_httpServer.onGet(RestProbe.PATH_MEMBERS, URL_CLUSTER_MEMBERS_EMPTY, MediaType.APPLICATION_JSON);

        RestProbe probe = new RestProbe(s_client);

        assertThat(probe.isClusterMember(), is(false));
        }

    @Test
    public void shouldBeNotClusterMemberIfNullMembers()
        {
        s_httpServer.onGet(RestProbe.PATH_MEMBERS, URL_CLUSTER_MEMBERS_NULL, MediaType.APPLICATION_JSON);

        RestProbe probe = new RestProbe(s_client);

        assertThat(probe.isClusterMember(), is(false));
        }

    @Test
    public void shouldGetClusterMemberIfResponseNot200()
        {
        s_httpServer.onGet(RestProbe.PATH_MEMBERS, p -> HttpServerStub.NOT_FOUND);

        RestProbe probe = new RestProbe(s_client);

        m_executionException.expect(RuntimeException.class);

        probe.isClusterMember();
        }

    @Test
    public void shouldGetServiceNamesForSingleService()
        {
        s_httpServer.onGet(RestProbe.PATH_SERVICES, URL_CLUSTER_SERVICES_SINGLE, MediaType.APPLICATION_JSON);

        RestProbe probe = new RestProbe(s_client);

        assertThat(probe.getPartitionAssignmentMBeans(), containsInAnyOrder(PARTITION_LINK_FOO));
        }

    @Test
    public void shouldGetServiceNamesForMultipleService()
        {
        s_httpServer.onGet(RestProbe.PATH_SERVICES, URL_CLUSTER_SERVICES, MediaType.APPLICATION_JSON);

        RestProbe probe = new RestProbe(s_client);

        assertThat(probe.getPartitionAssignmentMBeans(), containsInAnyOrder(PARTITION_LINK_FOO, PARTITION_LINK_BAR));
        }

    @Test
    public void shouldGetServiceNamesForNonPartitionedService()
        {
        s_httpServer.onGet(RestProbe.PATH_SERVICES, URL_CLUSTER_SERVICES_NOT_PARTITIONED, MediaType.APPLICATION_JSON);

        RestProbe probe = new RestProbe(s_client);

        assertThat(probe.getPartitionAssignmentMBeans(), is(emptyIterable()));
        }

    @Test
    public void shouldGetServiceNamesForEmptyService()
        {
        s_httpServer.onGet(RestProbe.PATH_SERVICES, URL_CLUSTER_SERVICES_EMPTY, MediaType.APPLICATION_JSON);

        RestProbe probe = new RestProbe(s_client);

        assertThat(probe.getPartitionAssignmentMBeans(), is(emptyIterable()));
        }

    @Test
    public void shouldGetServiceNamesForNullService()
        {
        s_httpServer.onGet(RestProbe.PATH_SERVICES, URL_CLUSTER_SERVICES_NULL, MediaType.APPLICATION_JSON);

        RestProbe probe = new RestProbe(s_client);

        assertThat(probe.getPartitionAssignmentMBeans(), is(emptyIterable()));
        }

    @Test
    public void shouldGetServiceNamesWhenResponseIsNot200()
        {
        s_httpServer.onGet(RestProbe.PATH_MEMBERS, p -> HttpServerStub.NOT_FOUND);

        RestProbe probe = new RestProbe(s_client);

        m_executionException.expect(RuntimeException.class);

        probe.getPartitionAssignmentMBeans();
        }


    // ----- helper methods -------------------------------------------------

    /**
     * Obtain the {@link URL} of the specified file or resource.
     *
     * @param sResource  the path to the file or resource
     *
     * @return  the {@link URL} of the file or resource or {@code null}
     *          if the resource does not exist
     */
    private static URL findURL(String sResource)
        {
        return Resources.findFileOrResource(sResource, null);
        }

    // ----- constants ------------------------------------------------------

    public static final URL URL_CLUSTER =  findURL("json/cluster.json");

    public static final URL URL_CLUSTER_MEMBERS_MULTI =  findURL("json/cluster-members-multiple-members.json");

    public static final URL URL_CLUSTER_MEMBERS_SINGLE =  findURL("json/cluster-members-single-member.json");

    public static final URL URL_CLUSTER_MEMBERS_EMPTY =  findURL("json/cluster-members-empty.json");

    public static final URL URL_CLUSTER_MEMBERS_NULL =  findURL("json/cluster-members-null.json");

    public static final URL URL_CLUSTER_SERVICES =  findURL("json/cluster-services.json");

    public static final URL URL_CLUSTER_SERVICES_SINGLE =  findURL("json/cluster-services-single.json");

    public static final URL URL_CLUSTER_SERVICES_EMPTY =  findURL("json/cluster-services-empty.json");

    public static final URL URL_CLUSTER_SERVICES_NULL =  findURL("json/cluster-services-null.json");

    public static final URL URL_CLUSTER_SERVICES_NOT_PARTITIONED  =  findURL("json/cluster-services-no-partition.json");

    public static final String PARTITION_LINK_FOO = "http://0%3A0%3A0%3A0%3A0%3A0%3A0%3A1:30000/management/coherence/cluster/services/Foo/partition";

    public static final String PARTITION_LINK_BAR = "http://0%3A0%3A0%3A0%3A0%3A0%3A0%3A1:30000/management/coherence/cluster/services/Bar/partition";

    // ----- data members ---------------------------------------------------

    @ClassRule
    public static HttpServerStub s_httpServer = new HttpServerStub();

    private static ProbeHttpClient s_client;

    @Rule
    public ExpectedException m_executionException = ExpectedException.none();
    }
