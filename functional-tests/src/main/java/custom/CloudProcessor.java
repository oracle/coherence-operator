/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package custom;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.Reader;
import java.net.URL;
import java.util.Objects;

import com.tangosol.io.pof.PofReader;
import com.tangosol.io.pof.PofWriter;
import com.tangosol.io.pof.PortableObject;
import com.tangosol.util.InvocableMap.Entry;
import com.tangosol.util.Resources;
import com.tangosol.util.processor.AbstractProcessor;


/**
 * @author cp
 */
public class CloudProcessor extends AbstractProcessor<String, Cloud, Cloud> implements PortableObject
    {
    private Cloud newOCICloud()
        {
        return new PortableCloud(OCI.getName());
        }

    private Cloud newGCPCloud()
        {
        return new PortableCloud(GCP.getName());
        }

    private String getVersion()
        {
        String version = "v1";

        URL url = Resources.findFileOrResource("version.txt", null);

        if (url != null)
            {
            try (BufferedReader reader = new BufferedReader(new InputStreamReader(url.openStream())))
                {
                version = reader.readLine().trim();
                }
            catch (Exception e)
                {
                e.printStackTrace();
                }
            }

            return version;
        }

    @Override
    public Cloud process(Entry<String, Cloud> entry)
        {
        System.out.println("Before Cloud Processor: " + entry.getValue().getName());

        String sVersion = getVersion();
        Cloud  cloud    = "v1".equals(sVersion) ? newGCPCloud() : newOCICloud();

        entry.setValue(cloud);

        Cloud newCloud = entry.getValue();

        System.out.println("After Cloud Processor: " + newCloud.getName());

        return newCloud;
        }

    @Override
    public void readExternal(PofReader in) throws IOException
        {
        }

    @Override
    public void writeExternal(PofWriter out) throws IOException
        {
        }

    @Override
    public int hashCode()
        {
        return Objects.hashCode(this);
        }

    @Override
    public boolean equals(Object o)
        {
        return Objects.equals(this, o);
        }

    // ----- data members ---------------------------------------------------

    private Cloud OCI = new PortableCloud(new OCI().getName());

    private Cloud GCP = new PortableCloud(new GCP().getName());
    }
