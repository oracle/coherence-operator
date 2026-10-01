<doc-view>

<h2 id="_storage_enabled_or_disabled_deployments">Storage Enabled or Disabled Deployments</h2>
<div class="section">
<p>Partitioned cache services that manage Coherence caches are configured as storage enabled or storage disabled.
Whilst it is possible to configure individual services to be storage enabled or disabled in the cache configuration file
and have a mixture of modes in a single JVM, typically all the services in a JVM share the same mode by setting the
<code>coherence.distributed.localstorage</code> system property to <code>true</code> for storage enabled members and to <code>false</code> for
storage disabled members. The <code>Coherence</code> CRD allows this property to be set by specifying the
<code>spec.coherence.storageEnabled</code> field to either true or false. The default value when nothing is specified is <code>true</code>.</p>

<markup
lang="yaml"
title="storage enabled"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  coherence:
    storageEnabled: true  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The <code>Coherence</code> resource specifically sets <code>coherence.distributed.localstorage</code> to <code>true</code></li>
</ul>
<markup
lang="yaml"
title="storage disabled"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  coherence:
    storageEnabled: false  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The <code>Coherence</code> resource specifically sets <code>coherence.distributed.localstorage</code> to <code>false</code></li>
</ul>
</div>
</doc-view>
