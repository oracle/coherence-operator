<doc-view>

<h2 id="_coherence_ipmonitor">Coherence IPMonitor</h2>
<div class="section">
<p>The Coherence IPMonitor is a failure detection mechanism used by Coherence to detect machine failures. It does this by pinging the echo port, (port 7) on remote hosts that other cluster members are running on. When running in Kubernetes, every Pod has its own IP address, so it looks to Coherence like every member is on a different host. Failure detection using IPMonitor is less useful in Kubernetes than it is on physical machines or VMs, so the Operator disables the IPMonitor by default. This is configurable though and if it is felt that using IPMonitor is useful to an application, it can be re-enabled.</p>

<p>To re-enable IPMonitor set the boolean flag <code>enableIpMonitor</code> in the <code>coherence</code> section of the Coherence resource yaml:</p>

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

<p>Setting <code>enableIpMonitor</code> will disable the IPMonitor, which is the default behaviour when <code>enableIpMonitor</code> is not specified in the yaml.</p>

</div>
</doc-view>
