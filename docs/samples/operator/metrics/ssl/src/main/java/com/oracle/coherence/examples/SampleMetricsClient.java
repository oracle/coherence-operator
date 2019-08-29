/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.examples;


import javax.ws.rs.client.Client;
import javax.ws.rs.client.WebTarget;
import javax.ws.rs.core.Response;

/**
 * Sample client to connect to metrics  via SSL.
 *
 * @author tam 2019-06-16
 */
public class SampleMetricsClient
    {
    public static void main(String ... args)
        {
        try {
            HttpSSLHelper clientHelper = new HttpSSLHelper(9612);
            Client clientStarLord = clientHelper.getClient(HttpSSLHelper.CERT_STAR_LORD,
                                                             HttpSSLHelper.STORE_PASSWORD,
                                                             HttpSSLHelper.KEY_PASSWORD,
                                                             HttpSSLHelper.TRUSTSTORE_GUARDIANS,
                                                             HttpSSLHelper.TRUST_PASSWORD);
            WebTarget webTarget = clientHelper.getHttpsWebTarget(clientStarLord, "/metrics");
            Response   response  = webTarget.request().get();
            int        status    = response.getStatus();
            if (status == 200)
                {
                System.out.println("\nSuccess, HTTP Response code is " + status + "\n");
                System.out.println(response.readEntity(String.class));
                }
            else
                {
                System.out.println("\nFailed, HTTP Response code is " + status + "\n");
                }
            }
        catch (Exception e)
            {
            System.out.println("Failed, Error = " + e.getMessage());
            }
        }
    }
