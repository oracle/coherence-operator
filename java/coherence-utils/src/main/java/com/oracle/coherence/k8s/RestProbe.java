/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import javax.ws.rs.NotFoundException;
import javax.ws.rs.core.MediaType;
import java.net.URI;
import java.util.Collections;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.Set;
import java.util.stream.Collectors;

/**
 * A readiness/liveness probe that uses Coherence management over ReST.
 *
 * @author jk
 */
public class RestProbe
        extends BaseProbe
    {
    // ----- constructors ---------------------------------------------------

    /**
     * Create a {@link RestProbe}.
     */
    public RestProbe()
        {
        this(new ProbeHttpClient());
        }

    /**
     * Create a {@link RestProbe}.
     *
     * @param client  the http client to use
     */
    RestProbe(ProbeHttpClient client)
        {
        f_client = client;
        }

    // ----- Probe methods --------------------------------------------------

    @Override
    public boolean isAvailable()
        {
        try
            {
            jsonQuery(PATH_CLUSTER);
            return true;
            }
        catch (Throwable e)
            {
            return false;
            }
        }

    @Override
    public void close()
        {
        f_client.close();
        }

    @Override
    @SuppressWarnings("unchecked")
    protected Set<String> getPartitionAssignmentMBeans()
        {
        return queryItems(PATH_SERVICES)
                    .stream()
                    .flatMap(m -> ((List<Map>)  m.getOrDefault(JSON_ATTRIBUTE_LINKS, Collections.emptyList())).stream())
                    .filter(m -> JSON_ATTRIBUTE_LINK_PARTITION.equals(m.get(JSON_ATTRIBUTE_LINKS_REL)))
                    .map(m -> (String) m.get(JSON_ATTRIBUTE_LINKS_HREF))
                    .filter(Objects::nonNull)
                    .collect(Collectors.toSet());
        }

    @Override
    @SuppressWarnings("unchecked")
    protected Map<String, Object> getMBeanAttributes(String sMBean, String[] asAttributes)
        {
        // Future optimization: only get attributes specified in asAttributes
        return jsonQuery(sMBean, true);
        }

    @Override
    public boolean isClusterMember()
        {
        return !queryItems(PATH_MEMBERS).isEmpty();
        }

    // ----- helper methods -------------------------------------------------

    /**
     * Perform a ReST query and return the json response as a {@link Map}.
     * <p>
     * The query can be just the part part of a query or a full URI. If the
     * query string is a full URI only the path part is used.
     * 
     * @param sQuery  the query to execute
     *
     * @return  the json response as a {@link Map} or null if the response
     *          code was not 200.
     */
    private Map jsonQuery(String sQuery)
        {
        return jsonQuery(sQuery, false);
        }

    /**
     * Perform a ReST query and return the json response as a {@link Map}.
     * <p>
     * The query can be just the part part of a query or a full URI. If the
     * query string is a full URI only the path part is used.
     *
     * @param sQuery  the query to execute
     *
     * @return  the json response as a {@link Map} or null if the response
     *          code was not 200.
     */
    private Map jsonQuery(String sQuery, boolean fAllowNotFound)
        {
        try
            {
            URI uri = URI.create(sQuery);

            return f_client.getWebTarget(uri.getSchemeSpecificPart())
                                            .request(MediaType.APPLICATION_JSON)
                                            .get(Map.class);
            }
        catch (NotFoundException e)
            {
            if (fAllowNotFound)
                {
                return null;
                }
            else
                {
                throw e;
                }
            }
        }

    /**
     * Perform a query for a json value and obtain the list
     * mapped to the {@code items} attribute of the json.
     *
     * @param sPath  the path to use for the json query
     *
     * @return  the {@link List} mapped to the {@code items} attribute
     *          of the json returned by the query
     */
    @SuppressWarnings("unchecked")
    private List<Map> queryItems(String sPath)
        {
        return (List) jsonQuery(sPath)
                .getOrDefault(JSON_ATTRIBUTE_ITEMS, Collections.emptyList());
        }

    // ----- constants ------------------------------------------------------

    /**
     * The json attribute name for the cluster members list.
     */
    public static final String JSON_ATTRIBUTE_ITEMS = "items";

    /**
     * The json attribute name for the links attribute.
     */
    public static final String JSON_ATTRIBUTE_LINKS = "links";

    /**
     * The json attribute name for the links rel attribute.
     */
    public static final String JSON_ATTRIBUTE_LINKS_REL = "rel";

    /**
     * The json attribute name for the links rel attribute.
     */
    public static final String JSON_ATTRIBUTE_LINKS_HREF = "href";

    /**
     * The json attribute name for the partition link.
     */
    public static final String JSON_ATTRIBUTE_LINK_PARTITION = "partition";

    /**
     * The base path to use for queries.
     */
    public static final String PATH_BASE = "/management/coherence";

    /**
     * The path to use to query whether management over ReST is available.
     */
    public static final String PATH_CLUSTER = PATH_BASE + "/cluster";

    /**
     * The path to use to query for services.
     */
    public static final String PATH_SERVICES = PATH_CLUSTER + "/services";

    /**
     * The path to use to query the member count.
     */
    public static final String PATH_MEMBERS = "/management/coherence/cluster/members?fields=memberName";

    // ----- data members ---------------------------------------------------

    /**
     * The http client to use.
     */
    private final ProbeHttpClient f_client;
    }
