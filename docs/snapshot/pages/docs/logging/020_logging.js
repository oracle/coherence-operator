<doc-view>

<h2 id="_log_capture_with_fluentd">Log Capture with Fluentd</h2>
<div class="section">
<p>There are many ways to capture container logs in Kubernetes, one possibility that this guide will cover is using
a Fluentd side-car container to ship log files to Elasticsearch.
This is a common pattern and one the the <code>Coherence</code> CRD makes simple by allowing easy injection of additional containers.</p>

<div class="admonition note">
<p class="admonition-inline">This guide is going to assume that the default logging related configurations provided by the operator will
be used. For example, Coherence will be configured to use Java util logging for logs, and the default logging configuration
file will be used. Whilst these things are not pre-requisites for shipping logs to Elasticsearch they are required
to make the examples below work.</p>
</div>
<p>To be able to send Coherence logs to Elasticsearch there are some steps that must be completed:</p>

<ul class="ulist">
<li>
<p>Configure Coherence to log to files</p>

</li>
<li>
<p>Add a <code>Volume</code> and <code>VolumeMount</code> to be used for log files</p>

</li>
<li>
<p>Add the Fluentd side-car container</p>

</li>
</ul>

<h3 id="_configure_coherence_to_log_to_files">Configure Coherence to Log to Files</h3>
<div class="section">
<p>Coherence will log to the console by default so to be able to ship logs to Elasticsearch it needs to be configured
to write to log files. One way to do this is to add a Java Util Logging configuration file and then to configure
Coherence to use the JDK logger.</p>

<p>In the <code>jvm.args</code> section of the <code>Coherence</code> CRD the system properties should be added to set the configuration file used by Java util logging and to configure Coherence logging.
See the Coherence <a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/standalone/coherence/14.1.1.0/develop-applications/operational-configuration-elements.html">Logging Config</a>
documentation for more details.</p>

<p>There are alternative ways to configure the Java util logger besides using a configuration file, just as there are
alternative logging frameworks that Coherence can be configured to use to produce log files.
This example is going to use Java util logging as that is the simplest to demonstrate without requiring any additional
logging libraries.</p>


<h4 id="_operator_provided_logging_configuration_file">Operator Provided Logging Configuration File</h4>
<div class="section">
<p>Whilst any valid Java util logging configuration file may be used, the Coherence Operator injects a default logging
configuration file into the <code>coherence</code> container that can be used to configure the logger to write
logs to files under the <code>/logs</code> directory. The log files will have the name <code>coherence-%g.log</code>, where <code>%g</code> is the
log file generation created as logs get rotated.</p>

<p>This file will be injected into the container at <code>/coherence-operator/utils/logging/logging.properties</code>
and will look something like this:</p>

<markup


>com.oracle.coherence.handlers=java.util.logging.ConsoleHandler,java.util.logging.FileHandler

com.oracle.coherence.level=FINEST

java.util.logging.ConsoleHandler.formatter=java.util.logging.SimpleFormatter
java.util.logging.ConsoleHandler.level=FINEST

java.util.logging.FileHandler.pattern=/logs/coherence-%g.log
java.util.logging.FileHandler.limit=10485760
java.util.logging.FileHandler.count=50
java.util.logging.FileHandler.formatter=java.util.logging.SimpleFormatter

java.util.logging.SimpleFormatter.format=%5$s%6$s%n</markup>

<p>To configure Cohrence and the logger some system properties need to be added to the <code>jvm.args</code> field
of the <code>Coherence</code> CRD spec:</p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: cluster-one
spec:
  jvm:
    args:
      - "-Dcoherence.log=jdk"                                                                   <span class="conum" data-value="1" />
      - "-Dcoherence.log.logger=com.oracle.coherence"                                           <span class="conum" data-value="2" />
      - "-Djava.util.logging.config.file=/coherence-operator/utils/logging/logging.properties"  <span class="conum" data-value="3" /></markup>

<ul class="colist">
<li data-value="1">Coherence has been configured to use the Java util logging.</li>
<li data-value="2">The Coherence logger name has been set to <code>com.oracle.coherence</code>, which matches the logging configuration file.</li>
<li data-value="3">The Java util logging configuration file is set to the file injected by the Operator.</li>
</ul>
</div>

<h4 id="_log_files_volume">Log Files Volume</h4>
<div class="section">
<p>The logging configuration above configures Coherence to write logs to the <code>/logs</code> directory.
For this location to be accessible to both the <code>coherence</code> container and to the <code>fluentd</code> container it needs to be
created as a <code>Volume</code> in the <code>Pod</code> and mounted to both containers.
As this <code>Volume</code> can be ephemeral and is typically not required to live longer than the <code>Pod</code> the simplest type of
<code>Volume</code> to use is an <code>emptyDir</code> volume source.</p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: cluster-one
spec:
  jvm:
    args:
      - "-Dcoherence.log=jdk"
      - "-Dcoherence.log.logger=com.oracle.coherence"
      - "-Djava.util.logging.config.file=/coherence-operator/utils/logging/logging.properties"
  volumes:
    - name: logs           <span class="conum" data-value="1" />
      emptyDir: {}
  volumeMounts:
    - name: logs           <span class="conum" data-value="2" />
      mountPath: /logs</markup>

<ul class="colist">
<li data-value="1">An additional empty-dir <code>Volume</code> named <code>logs</code> has been added to the <code>Coherence</code> spec.</li>
<li data-value="2">The <code>logs</code> volume will be mounted at <code>/logs</code> in all containers in the <code>Pod</code>.</li>
</ul>
</div>
</div>

<h3 id="_add_the_fluentd_side_car">Add the Fluentd Side-Car</h3>
<div class="section">
<p>With Coherence configured to write to log files, and those log files visible to other containers in the <code>Pod</code> the
Fluentd side-car container can be added.</p>

<p>The example yaml below shows a <code>Coherence</code> resource with the additional Fluentd container added.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: cluster-one
spec:
  jvm:
    args:
      - "-Dcoherence.log=jdk"
      - "-Dcoherence.log.logger=com.oracle.coherence"
      - "-Djava.util.logging.config.file=/coherence-operator/utils/logging/logging.properties"
  volumes:
    - name: logs
      emptyDir: {}
  volumeMounts:
    - name: logs
      mountPath: /logs
  sideCars:
    - name: fluentd                                     <span class="conum" data-value="1" />
      image: "fluent/fluentd-kubernetes-daemonset:v1.14-debian-elasticsearch7-1"
      args:
        - "-c"
        - "/etc/fluent.conf"
      env:
        - name: "FLUENTD_CONF"                          <span class="conum" data-value="2" />
          value: "fluentd-coherence.conf"
        - name: "FLUENT_ELASTICSEARCH_SED_DISABLE"      <span class="conum" data-value="3" />
          value: "true"
  configMapVolumes:
    - name: "efk-config"                                <span class="conum" data-value="4" />
      mountPath: "/fluentd/etc/fluentd-coherence.conf"
      subPath: "fluentd-coherence.conf"</markup>

<ul class="colist">
<li data-value="1">The <code>fluentd</code> container has been added to the <code>sideCars</code> list. This will create another container
in the <code>Pod</code> exactly as configured.</li>
<li data-value="2">The <code>FLUENTD_CONF</code> environment variable has been set to the name of the configuration file that Fluentd should use.
The standard Fluentd behaviour is to locate this file in the <code>/fluentd/etc/</code> directory.</li>
<li data-value="3">The <code>FLUENT_ELASTICSEARCH_SED_DISABLE</code> environment variable has been set to work around a known issue <a id="" title="" target="_blank" href="https://github.com/fluent/fluentd-kubernetes-daemonset#disable-sed-execution-on-elasticsearch-image">here</a>.</li>
<li data-value="4">An additional volume has been added from a <code>ConfigMap</code> named <code>efk-config</code>, that contains the Fluentd configuration to use.
This will be mounted to the <code>fluentd</code> container at <code>/fluentd/etc/fluentd-coherence.conf</code>, which corresponds to the
name of the file set in the <code>FLUENTD_CONF</code> environment variable.</li>
</ul>
<div class="admonition note">
<p class="admonition-inline">There is no need to add a <code>/logs</code> volume mount to the <code>fluentd</code> container. The operator will mount the <code>logs</code>
<code>Volume</code> to <strong>all</strong> containers in the <code>Pod</code>.</p>
</div>
<p>In the example above the Fluentd configuration has been provided from a <code>ConfigMap</code>. It could just as easily have come from a
<code>Secret</code> or some other external <code>Volume</code> mount, or it could have been baked into the Fluentd image to be used.</p>


<h4 id="_the_fluentd_configuration_file">The Fluentd Configuration File</h4>
<div class="section">
<p>The <code>ConfigMap</code> used to provide the Fluentd configuration might look something like this:</p>

<markup
lang="yaml"

>apiVersion: v1
kind: ConfigMap
metadata:
  name: efk-config                              <span class="conum" data-value="1" />
  labels:
    component: coherence-efk-config
data:
  fluentd-coherence.conf: |
    # Ignore fluentd messages
    &lt;match fluent.**&gt;
      @type null
    &lt;/match&gt;

    # Coherence Logs
    &lt;source&gt;                                    <span class="conum" data-value="2" />
      @type tail
      path /logs/coherence-*.log
      pos_file /tmp/cohrence.log.pos
      read_from_head true
      tag coherence-cluster
      multiline_flush_interval 20s
      &lt;parse&gt;
       @type multiline
       format_firstline /^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}.\d{3}/
       format1 /^(?&lt;time&gt;\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}.\d{3})\/(?&lt;uptime&gt;[0-9\.]+) (?&lt;product&gt;.+) &lt;(?&lt;level&gt;[^\s]+)&gt; \(thread=(?&lt;thread&gt;.+), member=(?&lt;member&gt;.+)\):[\S\s](?&lt;log&gt;.*)/
      &lt;/parse&gt;
    &lt;/source&gt;

    &lt;filter coherence-cluster&gt;                  <span class="conum" data-value="3" />
     @type record_transformer
     &lt;record&gt;
       cluster "#{ENV['COH_CLUSTER_NAME']}"
       role "#{ENV['COH_ROLE']}"
       host "#{ENV['HOSTNAME']}"
       pod-uid "#{ENV['COH_POD_UID']}"
     &lt;/record&gt;
    &lt;/filter&gt;

    &lt;match coherence-cluster&gt;                   <span class="conum" data-value="4" />
      @type elasticsearch
      hosts "http://elasticsearch-master:9200"
      logstash_format true
      logstash_prefix coherence-cluster
    &lt;/match&gt;</markup>

<ul class="colist">
<li data-value="1">The name of the <code>ConfigMap</code> is <code>efk-config</code> to match the name specified in the <code>Coherence</code> CRD spec.</li>
<li data-value="2">The <code>source</code> section is configured to match log files with the name <code>/logs/coherence-*.log</code>, which is the name that
Coherence logging has been configured to use. The pattern in the <code>source</code> section is a Fluentd pattern that matches the
standard Coherence log message format.</li>
<li data-value="3">A <code>filter</code> section will add additional fields to the log message. These come from the environment variables that
the Operator will inject into all containers in the Pod. In this case the Coherence cluster name, the Coherence role name,
the Pod host name and Pod UID.</li>
<li data-value="4">The final section tells Fluentd how to ship the logs to Elasticsearch, in this case to the endpoint <code><a id="" title="" target="_blank" href="http://elasticsearch-master:9200">http://elasticsearch-master:9200</a></code></li>
</ul>
<p>There are many ways to configure Fluentd, the example above is just one way and is in fact taken from one of the Operator&#8217;s functional tests.</p>

<p>With the <code>efk-config</code> <code>ConfigMap</code> created in the same namespace as the <code>Coherence</code> resource the Coherence logs from the
containers will now be shipped to Elasticsearch.</p>

</div>
</div>
</div>
</doc-view>
