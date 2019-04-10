/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package custom;

import java.io.IOException;

import com.tangosol.io.pof.PofReader;
import com.tangosol.io.pof.PofWriter;
import com.tangosol.io.pof.PortableObject;

/**
 * @author cp
 */
public class PortableCloud implements Cloud, PortableObject
    {

    private String cloudName;

    public PortableCloud()
        {
        }

    /**
     * 
     */
    public PortableCloud(String cname)
        {
        this.cloudName = cname;
        }

    @Override
    public String getName()
        {
        return cloudName;
        }

    @Override
    public int hashCode()
        {
        return 31 + cloudName.hashCode();
        }

    @Override
    public boolean equals(Object o)
        {
        if (o == null || !(o instanceof Cloud))
            {
            return false;
            }
        return getName().equals(((Cloud) o).getName());
        }

    /*
     * (non-Javadoc)
     * @see com.tangosol.io.pof.PortableObject#readExternal(com.tangosol.io.pof.
     * PofReader)
     */
    @Override
    public void readExternal(PofReader in) throws IOException
        {
        cloudName = in.readString(0);
        }

    /*
     * (non-Javadoc)
     * @see com.tangosol.io.pof.PortableObject#writeExternal(com.tangosol.io.pof.
     * PofWriter)
     */
    @Override
    public void writeExternal(PofWriter out) throws IOException
        {
        out.writeString(0, cloudName);
        }
    }
