<doc-view>

<h2 id="_java_util_logging_configuration">Java Util Logging Configuration</h2>
<div class="section">
<p>Coherence is configured to use Java util logging as its logger in the <code>Pods</code> that are started from the <code>Coherence</code>
resource. A default Java util logging configuration file will be injected into the <code>Pod</code> by the operator that configures
the Coherence logger to log to rolling log files in a <code>/logs</code> directory.</p>

<p>The Java util logging configuration file can be overridden using the <code>jvm.loggingConfig</code> field in the CRD.
The value should point to a file either in the image being run or loaded from a volume such as a <code>ConfigMap</code>.</p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  jvm:
    loggingConfig: storage-logging.properties <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The logging configuration file is set to <code>storage-logging.properties</code>.</li>
</ul>
</div>
</doc-view>
