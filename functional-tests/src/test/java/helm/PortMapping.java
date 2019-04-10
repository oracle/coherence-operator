/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package helm;

import com.oracle.bedrock.runtime.options.Ports;
import com.oracle.bedrock.util.Capture;

import java.util.Collections;
import java.util.Iterator;

/**
 * @author jk
 */
public class PortMapping
    {
    // ----- constructors ---------------------------------------------------

    public PortMapping(String sName, int nPortContainer)
        {
        this(sName, nPortContainer, nPortContainer);
        }

    public PortMapping(String sName, int nPortLocal, int nPortContainer)
        {
        this(sName, Collections.singleton(nPortLocal).iterator(), nPortContainer);
        }

    public PortMapping(String sName, Iterator<Integer> itPortLocal, int nPortContainer)
        {
        this(sName, new Capture<>(itPortLocal), nPortContainer);
        }

    public PortMapping(String sName, Capture<Integer> port, int nPortContainer)
        {
        m_sName         = sName;
        m_capture       = port;
        m_portContainer = nPortContainer;
        }

    // ----- PortMapping methods ------------------------------------------------

    @Override
    public String toString()
        {
        return m_capture.get() + ":" + m_portContainer;
        }

    public Ports.Port getPort()
        {
        return new Ports.Port(m_sName, m_capture.get(), m_portContainer);
        }

    // ----- data members ---------------------------------------------------

    private String m_sName;

    private Capture<Integer> m_capture;

    private int m_portContainer;
    }
