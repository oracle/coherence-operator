/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package custom;


import com.tangosol.net.DefaultCacheServer;

/**
 * @author lh
 */
public class ServerMain
    {
    public static void main(String[] argv)
        {
        DefaultCacheServer server =DefaultCacheServer.startServerDaemon();
        server.waitForServiceStart();

        for (int i = 0; i < 1200; i++)
            {
            try
                {
                Thread.sleep(100);
                }
            catch (Exception e)
                {}
            }
        }
    }
