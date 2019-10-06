<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Logging Configuration</dt>
<dd slot="desc"><p>There are various settings in a Coherence role that control different aspects of logging, including the Coherence
log level, configuration files and whether Fluentd log capture is enabled.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_logging_configuration">Logging Configuration</h2>
<div class="section">
<p>Logging configuration for a role is defined in the <code>logging</code> section of the role&#8217;s <code>spec</code>. There are a number of different
fields used to configure different logging features. The <code>logging</code> configuration can be set at different places depending
on whether the implicit role or explicit roles are being configured.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  coherence:
    logLevel: 9
  logging:
    configFile:    app-logging.properties
    configMapName: logging-cm
    fluentd:
      enabled: true
      image: fluent/fluentd-kubernetes-daemonset:v1.3.3-debian-elasticsearch-1.3
      imagePullPolicy: IfNotPresent
      configFile: fluentd-config.yaml
      tag: test-cluster</markup>

<p>The fields in the example above are described in detail in the following sections.</p>


<h3 id="_coherence_log_level">Coherence Log Level</h3>
<div class="section">
<p>The Coherence log level is set with the <code>coherence.logLevel</code> field. This field is an integer value between zero and nine
(see the Coherence documentation for a fuller explanation).</p>

<p>To set the Coherence log level when defining the implicit role:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  coherence:
    logLevel: 5 <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The implicit role will have a Coherence log level of <code>5</code></li>
</ul>
<p>To set the log level for explicit roles in the <code>roles</code> list:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  roles:
    - role: data
      coherence:
        logLevel: 9 <span class="conum" data-value="1" />
    - role: proxy
      coherence:
        logLevel: 5 <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>data</code> role will have a Coherence log level of <code>9</code></li>
<li data-value="2">The <code>proxy</code> role will have a Coherence log level of <code>5</code></li>
</ul>
<p>To set the log level for explicit roles in the <code>roles</code> list:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  coherence:
    logLevel: 9 <span class="conum" data-value="1" />
  roles:
    - role: data
    - role: proxy
    - role: web
      coherence:
        logLevel: 5 <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>data</code> and <code>proxy</code> roles will use the default Coherence log level of <code>9</code></li>
<li data-value="2">The <code>web</code> role overrides the default Coherence log level setting it to <code>5</code></li>
</ul>
</div>

<h3 id="_logging_config_file">Logging Config File</h3>
<div class="section">
<p>The default logging configuration for Coherence clusters started by the Coherence Operator is to set Coherence to used
JDK logging; the JDK logger is then configured with a configuration file. The default configuration file is embedded into
the Pod by the Coherence Operator but this default my be overridden; for example an application deployed into the cluster
may require different logging configurations. The name of the file is provided in the <code>logging.configFile</code> field.
The logging configuration file must be available to the JVM when it starts, either by providing it in
<router-link to="#clusters/065_application_image.adoc" @click.native="this.scrollFix('#clusters/065_application_image.adoc')">application code</router-link> or by <router-link to="/clusters/150_volumes">mounting a volume</router-link> containing
the file, or by using a <router-link to="#configmap" @click.native="this.scrollFix('#configmap')">ConfigMap</router-link>.</p>

<p>To set the logging configuration file when defining the implicit role:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  logging:
    configFile: app-logging.properties <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The implicit role will use the <code>app-logging.properties</code> logging configuration file</li>
</ul>
<p>To set the logging configuration file when defining explicit roles in the <code>roles</code> list:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  roles:
    - role: data
      logging:
        configFile: data-logging.properties <span class="conum" data-value="1" />
    - role: proxy
      logging:
        configFile: proxy-logging.properties <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>data</code> role will use the <code>data-logging.properties</code> logging configuration file</li>
<li data-value="2">The <code>proxy</code> role will use the <code>proxy-logging.properties</code> logging configuration file</li>
</ul>
<p>To set a default logging configuration file when defining explicit roles in the <code>roles</code> list:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  logging:
    configFile: app-logging.properties <span class="conum" data-value="1" />
  roles:
    - role: data
    - role: proxy
    - role: web
      logging:
        configFile: web-logging.properties <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>app-logging.properties</code> logging configuration file is set as the default ans will be used by the <code>data</code> and
<code>proxy</code> roles.</li>
<li data-value="2">The <code>web</code> role has a specific configuration file set and will use the <code>web-logging.properties</code> file</li>
</ul>
</div>

<h3 id="configmap">Logging ConfigMap</h3>
<div class="section">
<p>The <code>logging.ConfigMap</code> field can be used to specify the name of a <code>ConfigMap</code> that contains the logging configuration file
to use. The <code>ConfigMap</code> should exist in the same namespace as the Coherence cluster.</p>

<p>TBD&#8230;&#8203;</p>

</div>
</div>

<h2 id="_fluentd_logging_configuration">Fluentd Logging Configuration</h2>
<div class="section">
<p>The Coherence Operator allows Coherence cluster <code>Pods</code> to be configured with a Fluentd side-car container that will push
Coherence logs to Elasticsearch. The configuration for Fluentd is in the <code>logging.fluentd</code> section of the spec.</p>

<p>TBD&#8230;&#8203;</p>

</div>
</doc-view>
