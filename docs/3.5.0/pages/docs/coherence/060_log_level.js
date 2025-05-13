<doc-view>

<h2 id="_set_the_coherence_log_level">Set the Coherence Log Level</h2>
<div class="section">
<p>Logging granularity in Coherence is controlled by a log level, that is a number between one and nine,
where the higher the number the more debug logging is produced. The <code>Coherence</code> CRD has a field
<code>spec.coherence.logLevel</code> that allows the log level to be configured by setting the <code>coherence.log.level</code>
system property.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  coherence:
    logLevel: 9  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The <code>Coherence</code> spec sets the log level to 9, effectively passing <code>-Dcoherence.log.level=9</code> to the Coherence
JVM&#8217;s command line.</li>
</ul>
</div>
</doc-view>
