/*
 * File: MultiplexingApplicationConsole.java
 *
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS HEADER.
 *
 * The contents of this file are subject to the terms and conditions of 
 * the Common Development and Distribution License 1.0 (the "License").
 *
 * You may not use this file except in compliance with the License.
 *
 * You can obtain a copy of the License by consulting the LICENSE.txt file
 * distributed with this file, or by consulting https://oss.oracle.com/licenses/CDDL
 *
 * See the License for the specific language governing permissions
 * and limitations under the License.
 *
 * When distributing the software, include this License Header Notice in each
 * file and include the License file LICENSE.txt.
 *
 * MODIFICATIONS:
 * If applicable, add the following below the License Header, with the fields
 * enclosed by brackets [] replaced by your own identifying information:
 * "Portions Copyright [year] [name of copyright owner]"
 */

package util;

import com.oracle.bedrock.runtime.ApplicationConsole;
import com.oracle.bedrock.runtime.ApplicationConsoleBuilder;
import com.oracle.bedrock.runtime.console.AbstractPipedApplicationConsole;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.PrintWriter;
import java.util.Arrays;
import java.util.List;
import java.util.stream.Collectors;

/**
 * An implementation of an {@link ApplicationConsole} that
 * multiplexes output to other {@link ApplicationConsole}s.
 * <p>
 * Copyright (c) 2019. All Rights Reserved. Oracle Corporation.<br>
 * Oracle is a registered trademark of Oracle Corporation and/or its affiliates.
 *
 * @author Jonathan Knight
 */
public class MultiplexingApplicationConsole
        extends AbstractPipedApplicationConsole
{
    /**
     * The {@link ApplicationConsole}s to multiplex.
     */
    private final List<ApplicationConsole> consoles;

    /**
     * The {@link Thread} capturing StdOut lines
     */
    private final Thread stdoutThread;

    /**
     * The {@link Thread} capturing StdErr lines
     */
    private final Thread stderrThread;


    /**
     * Constructs {@link MultiplexingApplicationConsole}.
     * <p>
     * This constructor will set the maximum number of lines to capture to {@link Integer#MAX_VALUE}.
     */
    public MultiplexingApplicationConsole(ApplicationConsole... consoles)
    {
        super(DEFAULT_PIPE_SIZE, false);

        List<PrintWriter> outWriters = Arrays.stream(consoles)
                        .map(ApplicationConsole::getOutputWriter)
                        .collect(Collectors.toList());

        List<PrintWriter> errWriters = Arrays.stream(consoles)
                .map(ApplicationConsole::getErrorWriter)
                .collect(Collectors.toList());

        this.consoles     = Arrays.asList(consoles);
        this.stdoutThread = new Thread(new OutputCaptor(stdoutReader, outWriters));
        this.stderrThread = new Thread(new OutputCaptor(stderrReader, errWriters));

        this.stdoutThread.start();
        this.stderrThread.start();
    }


    @Override
    public void close()
    {
        super.close();

        consoles.forEach(this::closeSafely);

        try
        {
            stdoutThread.join();
            stderrThread.join();
        }
        catch (InterruptedException e)
        {
            // Ignored
        }
    }


    /**
     * Obtain an {@link ApplicationConsoleBuilder} that builds a
     * {@link MultiplexingApplicationConsole} wrapping all of the
     * {@link ApplicationConsole}s produced by the builders.
     *
     * @param builders  The {@link ApplicationConsoleBuilder}s to build
     *                  the wrapped {@link ApplicationConsole}s
     *
     * @return  An {@link ApplicationConsoleBuilder} that builds a
     *          {@link MultiplexingApplicationConsole}
     */
    public static ApplicationConsoleBuilder builder(ApplicationConsoleBuilder... builders)
    {
        return new Builder(Arrays.asList(builders));
    }

    private void closeSafely(ApplicationConsole console)
    {
        try
        {
            console.close();
        }
        catch (Throwable t)
        {
            // Ignored
        }
    }

    /**
     * Obtains a {@link PrintWriter} that can be used to write to the stdin
     * of an {@link ApplicationConsole}.
     *
     * @return a {@link PrintWriter}
     */
    public PrintWriter getInputWriter()
    {
        return stdinWriter;
    }


    private static class Builder
            implements ApplicationConsoleBuilder
    {
        private final List<ApplicationConsoleBuilder> builders;

        private Builder(List<ApplicationConsoleBuilder> builders)
        {
            this.builders = builders;
        }

        @Override
        public ApplicationConsole build(String applicationName)
        {
            ApplicationConsole[] consoles = builders.stream()
                    .map(builder -> builder.build(applicationName))
                    .toArray(ApplicationConsole[]::new);

            return new MultiplexingApplicationConsole(consoles);
        }
    }

    /**
     * The {@link Runnable} used to capture lines of output.
     */
    class OutputCaptor implements Runnable
    {
        /**
         * The {@link BufferedReader} to capture output from.
         */
        private BufferedReader reader;

        /**
         * The {@link PrintWriter}s to write to.
         */
        private final List<PrintWriter> writers;


        /**
         * Create an {@link OutputCaptor}.
         *
         * @param reader   The {@link BufferedReader} to capture output from
         * @param writers  The {@link PrintWriter}s to write to
         */
        private OutputCaptor(BufferedReader reader, List<PrintWriter> writers)
        {
            this.reader  = reader;
            this.writers = writers;
        }


        /**
         * The {@link Runnable#run()} method for this {@link OutputCaptor}
         * that will capture output.
         */
        @Override
        public void run()
        {
            try
            {
                String line = reader.readLine();

                while (line != null)
                {
                    for (PrintWriter writer : writers)
                    {
                        writer.println(line);
                    }

                    line = reader.readLine();
                }
            }
            catch (IOException e)
            {
                // Skip: Likely caused by application termination
            }
        }
    }
}
