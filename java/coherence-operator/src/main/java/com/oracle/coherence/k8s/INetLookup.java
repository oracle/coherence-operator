package com.oracle.coherence.k8s;

import java.net.InetAddress;
import java.net.UnknownHostException;

/**
 * A utility to perform a DNS lookup in Java.
 */
public class INetLookup
    {
    public static void main(String[] args)
        {
        for (String sName: args)
            {
            System.out.println("Looking up: " + sName);
            try
                {
                InetAddress[] aAddress = InetAddress.getAllByName(sName);
                for (InetAddress address : aAddress)
                    {
                    System.out.println(address);
                    }
                System.out.println();
                }
            catch (UnknownHostException e)
                {
                System.out.println(e.getMessage());
                }
            }
        }
    }
