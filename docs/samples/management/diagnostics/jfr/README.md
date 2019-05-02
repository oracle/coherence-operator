# Produce and extract a Java Flight Recorder (JFR) file

Java Flight Recorder (JFR) is a tool for collecting diagnostic and profiling data 
about a running Java application. It is integrated into the Java Virtual Machine (JVM) 
and causes almost no performance overhead, so it can be used even in heavily loaded production environments.

TIn Coherence 12.2.1.4.0 and above, the [Management over REST](../../rest) functionality provides 
the ability to create and managed JFR recordings.

This sample shows how to execute a JFR operation across all nodes of a cluster.

Valid JFR commands are jfrStart, jfrStop, jfrDump, and jfrCheck.
     
See [documentation](https://docs.oracle.com/javacomponents/jmc-5-4/jfr-runtime-guide/run.htm#JFRUH176) for more details on commands.  
 
> Note, use of Management over REST is only available when using the
> operator with Coherence 12.2.1.4.    

[Return to Diagnostics Tools](../) / [Return to Management samples](../../) / [Return to samples](../../../README.md#list-of-samples)

