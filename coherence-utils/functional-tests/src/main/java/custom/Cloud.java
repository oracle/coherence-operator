/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package custom;


import com.tangosol.io.pof.PortableObject;


/**
 * @author cp
 */
public interface Cloud extends PortableObject
    {
    default public String getName()
        {
        return "Cloud";
        }
    }
