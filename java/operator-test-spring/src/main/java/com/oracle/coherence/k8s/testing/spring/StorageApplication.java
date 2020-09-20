/*
 * Copyright (c) 2020, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s.testing.spring;

import java.util.Arrays;
import java.util.logging.Logger;

import com.tangosol.net.DefaultCacheServer;
import com.tangosol.net.Session;

import org.springframework.boot.ApplicationRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
import org.springframework.core.annotation.Order;

/**
 * @author Jonathan Knight  2020.09.10
 */
@SpringBootApplication
public class StorageApplication {
	public static final Logger logger = Logger.getLogger(StorageApplication.class.getName());

	private static String[] arguments;

    public static void main(String[] args) {
    	arguments = args;
   		SpringApplication.run(StorageApplication.class, args);
   	}

   	@Bean(name = "commandLineArguments")
	public String[] commandLineArguments() {
    	return arguments;
	}

	@Bean
	@Order(1)
 	public ApplicationRunner runCoherence() {
		return (args) -> {
			logger.info("Starting DefaultCacheServer with args " + Arrays.toString(args.getSourceArgs()));
			DefaultCacheServer.main(args.getSourceArgs());
		};
	}

	@Bean
	 public Session createCoherenceSession() {
		 return Session.create();
	 }
}
