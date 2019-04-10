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
public class OCI implements Cloud, PortableObject
    {

    @Override
    public String getName()
        {
        return OCI.class.getSimpleName();
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
        int result = super.hashCode();
        return result * 31 + getName().hashCode();
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
    }
