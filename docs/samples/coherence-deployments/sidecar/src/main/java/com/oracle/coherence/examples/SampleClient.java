/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.examples;

import com.tangosol.net.NamedCache;
import com.tangosol.net.Session;

import static com.tangosol.net.cache.TypeAssertion.withTypes;

/**
 * Sample client to connect to cluster and issue Lambda entry processor.
 *
 * @author tam 2019-04-16
 */
public class SampleClient
    {

    public static void main(String ... args)
        {
        try (Session session = Session.create())
            {
            NamedCache<Integer, Person> cache = session.getCache("person", withTypes(Integer.class, Person.class));

            int nId = 1;
            Person person1 = new Person(nId, "Tom Jones", "123 Hollywood Ave, California, USA");
            System.out.println("\nNew Person is: " + person1);
            cache.put(nId, person1);

            // invoke server side Lambda Entry Processor
            cache.invoke(nId, (entry) -> {
                Person person = entry.getValue();
                person.setAddress(person.getAddress().toUpperCase());
                person.setName(person.getName().toUpperCase());
                entry.setValue(person);
                return null;
            });

            System.out.println("Person after entry processor is: " + cache.get(nId) + "\n");
            }
        catch (Exception e)
            {
            System.err.println("Exception: " + e.getMessage());
            }

        System.exit(1);
        }
    }
