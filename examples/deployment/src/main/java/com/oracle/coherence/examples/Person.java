/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.examples;

import java.util.Objects;

import com.tangosol.io.pof.annotation.Portable;
import com.tangosol.io.pof.annotation.PortableProperty;


/**
 * Sample server side class to represent a person.
 *
 * @author tam 2019-04-16
 */
@Portable
public class Person {

    // ----- constructors ---------------------------------------------------

    /**
     * No-args constructor for POF.
     */
    public Person() {
        }

    /**
     * Create a Person.
     *
     * @param nId       id of the Person
     * @param sName     name of the Person
     * @param sAddress  address of the Person
     */
    public Person(int nId, String sName, String sAddress) {
        id       = nId;
        name    = sName;
        address = sAddress;
        }

    // ----- accessors ---------------------------------------------------

    /**
     * Returns the id of the Person.
     *
     * @return the id of the Person
     */
    public int getId() {
        return id;
        }

    /**
     * Sets the id of the Person.
     *
     * @param nId  the id of the Person
     */
    public void setId(int nId) {
        id = nId;
        }

    /**
     * Returns the name of the Person.
     *
     * @return the name of the Person
     */
    public String getName() {
        return name;
        }

    /**
     * Sets the name of the Person.
     *
     * @param sName the name of the Person
     */
    public void setName(String sName) {
        name = sName;
        }

    /**
     * Returns the address of the Person.
     *
     * @return the address of the Person
     */
    public String getAddress() {
        return address;
        }

    /**
     * Sets the address of the Person.
     *
     * @param sAddress the address of the Person
     */
    public void setAddress(String sAddress) {
        address = sAddress;
        }

    // ----- Object methods -------------------------------------------------

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
            }
        if (o == null || getClass() != o.getClass()) {
            return false;
            }
        Person person = (Person) o;
        return id == person.id
               && Objects.equals(name, person.name)
               && Objects.equals(address, person.address);
        }

    @Override
    public int hashCode() {
        return Objects.hash(id, name, address);
        }

    @Override
    public String toString() {
        return "Person{"
               + "Id=" + id
               + ", Name='" + name + '\''
               + ", Address='" + address + '\''
               + '}';
        }

    // ----- data members ------------------------------------------------

    @PortableProperty(0)
    private int id;

    @PortableProperty(1)
    private String name;

    @PortableProperty(2)
    private String address;
    }
