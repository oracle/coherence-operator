/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package custom;

import java.io.IOException;
import java.util.logging.FileHandler;

public class CustomFileHandler extends FileHandler
    {
    public CustomFileHandler()
        throws IOException, SecurityException
        {
        super();
        }
    }
