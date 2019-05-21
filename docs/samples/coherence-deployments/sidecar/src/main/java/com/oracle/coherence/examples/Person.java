/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.examples;

import com.tangosol.io.pof.annotation.Portable;
import com.tangosol.io.pof.annotation.PortableProperty;

import java.util.Objects;

/**
 * Sample server side class to represent a person.
 * 
 * @author tam 2019-04-16
 */
@Portable
public class Person
    {

    // ----- constructors ---------------------------------------------------

    /**
     * No-args constructor for POF.
     */
    public Person()
        {
        }

    /**
     * Create a Person.
     *
     * @param nId       id of the Person
     * @param sName     name of the Person
     * @param sAddress  address of the Person
     */
    public Person(int nId, String sName, String sAddress)
        {
        m_id       = nId;
        m_sName    = sName;
        m_sAddress = sAddress;
        }

    // ----- accessors ---------------------------------------------------

    /**
     * Returns the id of the Person.
     *
     * @return the id of the Person
     */
    public int getId()
        {
        return m_id;
        }

    /**
     * Sets the id of the Person.
     *
     * @param nId  the id of the Person
     */
    public void setId(int nId)
        {
        m_id = nId;
        }

    /**
     * Returns the name of the Person.
     *
     * @return the name of the Person
     */
    public String getName()
        {
        return m_sName;
        }

    /**
     * Sets the name of the Person.
     *
     * @param sName the name of the Person
     */
    public void setName(String sName)
        {
        m_sName = sName;
        }

    /**
     * Returns the address of the Person.
     *
     * @return the address of the Person
     */
    public String getAddress()
        {
        return m_sAddress;
        }

    /**
     * Sets the address of the Person.
     *
     * @param sAddress the address of the Person
     */
    public void setAddress(String sAddress)
        {
        m_sAddress = sAddress;
        }


    // ----- Object methods -------------------------------------------------
    
    @Override
    public boolean equals(Object o)
        {
        if (this == o)
            {
            return true;
            }
        if (o == null || getClass() != o.getClass())
            {
            return false;
            }
        Person person = (Person) o;
        return m_id == person.m_id &&
               Objects.equals(m_sName, person.m_sName) &&
               Objects.equals(m_sAddress, person.m_sAddress);
        }

    @Override
    public int hashCode()
        {
        return Objects.hash(m_id, m_sName, m_sAddress);
        }

    @Override
    public String toString()
        {
        return "Person{" +
               "Id=" + m_id +
               ", Name='" + m_sName + '\'' +
               ", Address='" + m_sAddress + '\'' +
               '}';
        }

    // ----- data members ------------------------------------------------

    @PortableProperty(0)
    private int m_id;

    @PortableProperty(1)
    private String m_sName;

    @PortableProperty(2)
    private String m_sAddress;
    }
