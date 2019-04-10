/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package custom;

import java.util.Date;
import java.text.SimpleDateFormat;

import com.tangosol.net.CacheFactory;
import com.tangosol.net.NamedCache;


/**
 * @author cp
 */
public class CloudClient
    {
    public static void main(String[] argv) throws Exception
        {
        AWS                       awsCloud       = new AWS();
        CloudProcessor            cloudProcessor = new CloudProcessor();
        NamedCache<String, Cloud> cache          = null;
        long                      cSize          = 0;

        while (true)
            {
            try
                {
                cache = CacheFactory.getCache("cloud-cache");
                cSize = cache.size();
                }
            catch (Exception e)
                {
                System.out.println("Storage Cluster is *NOT* available due to temporary condition. " +
                                           "Will resume once the cluster is available!");
                }

            if (cache != null)
                {
                try
                    {
                    cache.put("mynewcloud-" + ++cSize, awsCloud);

                    SimpleDateFormat dateFormat = new SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss.SSS");

                    System.out.println("[Time Key Added   : "
                            + dateFormat.format(new Date(System.currentTimeMillis()))
                            + "] [Key: mynewcloud-" + cSize + "]" + " [Cache Value Before Cloud EntryProcessor: "
                            + cache.get("mynewcloud-" + cSize).getName() + "]");

                    cache.invoke("mynewcloud-" + cSize, cloudProcessor);

                    System.out.println("[Time Key Updated : "
                            + dateFormat.format(new Date(System.currentTimeMillis()))
                            + "] [Key: mynewcloud-" + cSize + "]" + " [Cache Value After Cloud EntryProcessor: "
                            + cache.get("mynewcloud-" + cSize).getName() + "]");

                    System.out.println("[CacheSize: " + cache.size() + "]");
                    }
                catch (Exception ignored)
                    {
                    // ignored - expected exception when server shuts down
                    }
                }
            Thread.sleep(10000);
            }
        }
    }
