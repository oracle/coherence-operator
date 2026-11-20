<doc-view>

<h2 id="_coherence_ipmonitor">Coherence IPMonitor</h2>
<div class="section">
<p>The Coherence IPMonitor is a failure detection mechanism used by Coherence to detect machine failures.
It does this by pinging the echo port, (port 7) on remote hosts that other cluster members are running on.
When running in Kubernetes, every Pod has its own IP address, so it looks to Coherence like every member is on a different host.
Failure detection using IPMonitor is less useful in Kubernetes than it is on physical machines or VMs, so the Operator disables
the IPMonitor by default. This is configurable though and if it is felt that using IPMonitor is useful to an application,
it can be re-enabled.</p>


<h3 id="_coherence_warning_message">Coherence Warning Message</h3>
<div class="section">
<p>Disabling IP Monitor causes Coherence to print a warning in the logs similar to the one shown below.
This can be ignored when using the Operator.</p>

<markup


>2024-07-01 14:43:55.410/3.785 Oracle Coherence GE 14.1.1.2206.10 (dev-jonathanknight) &lt;Warning&gt; (thread=Coherence, member=n/a): IPMonitor has been explicitly disabled, this is not a recommended practice and will result in a minimum death detection time of 300 seconds for failed machines or networks.</markup>

</div>

<h3 id="_re_enable_the_ip_monitor">Re-Enable the IP Monitor</h3>
<div class="section">
<p>To re-enable IPMonitor set the boolean flag <code>enableIpMonitor</code> in the <code>coherence</code> section of the Coherence resource yaml.</p>

<div class="admonition caution">
<p class="admonition-textlabel">Caution</p>
<p ><p>The Coherence IP Monitor works by using Java&#8217;s <code>INetAddress.isReachable()</code> method to "ping" another cluster member&#8217;s IP address.
Under the covers the JDK will use an ICMP echo request to port 7 of the server. This can fail if port 7 is blocked,
for example using firewalls, or in Kubernetes using Network Policies or tools such as Istio.
In particular when using Network Policies it is impossible to open a port for ICMP as currently Network Policies
only support TCP or UDP and not ICMP.</p>

<p>If the Coherence IP Monitor is enabled in a Kubernetes cluster where port 7 is blocked then the cluster will fail to start.
Typically, the issue will be seen as one member will start and become the senior member. None of the other cluster members
will be abe to get IP Monitor to connect to the senior member, so they wil fail to start.</p>
</p>
</div>
<p>The yaml below shows an example of re-enabling the IP Monitor.</p>

<markup
lang="yaml"
title="coherence-storage.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  coherence:
    enableIpMonitor: true</markup>

<p>Setting <code>enableIpMonitor</code> field to <code>false</code> will disable the IPMonitor, which is the default behaviour when <code>enableIpMonitor</code> is
not specified in the yaml.</p>

</div>
</div>
</doc-view>
