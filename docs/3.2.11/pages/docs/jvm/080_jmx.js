<doc-view>

<h2 id="_using_jmx">Using JMX</h2>
<div class="section">
<p>The Java Management Extensions (JMX) are a common way to connect to a JVM and retrieve information from MBeans
attributes or trigger operations by calling MBean methods. By default, the JVM uses RMI as the transport layer for
JMX but RMI can be notoriously tricky to make work in a container environment due to the address and port NAT&#8217;ing
that is typical with containers or clouds. For this reason the Operator supports an alternative transport called JMXMP.
The difference is that JMXMP only requires a single port for communications and this port is simple to configure.</p>

<p>JMXMP is configured using the fields in the <code>jvm.jmxmp</code> section of the <code>Coherence</code> CRD spec.
Enabling JMXMP support adds the <code>opendmk_jmxremote_optional_jar.jar</code> JMXMP library to the classpath and sets the
the Coherence MBean server factory to produce a JMXMP MBean server. By default, the JMXMP server will bind
to port 9099 in the container but this can be configured to bind to a different port.</p>

<div class="admonition note">
<p class="admonition-inline">Using a custom transport for JMX, such as JMXMP, requires any JMX client that will connect to the JMX server to
also have a JMXMP library on its classpath.</p>
</div>
<div class="admonition warning">
<p class="admonition-textlabel">Warning</p>
<p ><p>JMXMP does not support secure transports such as TLS so cannot be recommended for production use.
If used in production clusters, then the JMXMP ports should be secured behind TLS enabled ingress or
with suitable network policies.</p>

<p>Coherence has other mechanisms to access management APIs and metrics that do support TLS.</p>
</p>
</div>
<p>See the <router-link to="/docs/management/030_visualvm">VisualVM Example</router-link> for a detailed example of how to configure
JMX and connect to a server in a <code>Coherence</code> resource.</p>

<p>Example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  jvm:
    jmxmp:
      enabled: true  <span class="conum" data-value="1" />
      port: 9099     <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">JMXMP is enabled so that a JMXMP server will be started in the Coherence container&#8217;s JVM</li>
<li data-value="2">The port that the JMX server will bind to in the container is 9099</li>
</ul>
</div>
</doc-view>
