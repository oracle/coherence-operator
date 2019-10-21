<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Management over ReST</dt>
<dd slot="desc"><p>Since version 12.2.1.4 Coherence has had functionality to expose a management API over ReST.
This API is disabled by default in Coherence clusters but can be enabled and configured by setting the relevant fields
in the <code>CoherenceCluster</code> resource.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_management_over_rest">Management over ReST</h2>
<div class="section">
<p>This example shows how to enable and access Coherence MBeans using Management over ReST.</p>

<p>Once the Management port is exposed via a load balancer or port-forward command the ReEST
endpoint is available at <code><a id="" title="" target="_blank" href="http://host:port/management/coherence/cluster">http://host:port/management/coherence/cluster</a></code> and the Swagger JSON document is available at <code><a id="" title="" target="_blank" href="http://host:port/management/coherence/cluster/metadata-catalog">http://host:port/management/coherence/cluster/metadata-catalog</a></code>.</p>

<p>See <a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/fusion-middleware/coherence/12.2.1.4/rest-reference/index.html">REST API for Managing Oracle Coherence</a> for
full details on each of the endpoints.</p>

<p>For more details on enabling Management over ReST including enabling SSL, please see the
<router-link to="/clusters/058_coherence_management">Coherence Operator documentation</router-link>.</p>

<p>See the <a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/fusion-middleware/coherence/12.2.1.4/manage/using-jmx-manage-oracle-coherence.html">Coherence Management</a> documentation for more information.</p>

<div class="admonition note">
<p class="admonition-inline">Note: Use of Management over ReST is available only when using the operator with clusters running
Coherence 12.2.1.4 or later version.</p>
</div>

<h3 id="_1_install_a_coherence_cluster_with_management_over_rest_enabled">1. Install a Coherence cluster with Management over ReST enabled</h3>
<div class="section">
<p>Deploy a simple management enabled <code>CoherenceCluster</code> resource with a single role like this:</p>

<markup
lang="yaml"
title="management-cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: management-cluster
spec:
  role: storage
  replicas: 3
  coherence:
    management:
      enabled: true            <span class="conum" data-value="1" />
  ports:
    - name: management
      port: 30000              <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">Indicates to enable Management over ReST</li>
<li data-value="2">The management port must be added to the additional <code>ports</code> list so that it is exposed on a service</li>
</ul>
<p>The yaml above can be installed into Kubernetes using <code>kubectl</code>:</p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; create -f management-cluster.yaml

coherencecluster.coherence.oracle.com/management-cluster created

kubectl -n &lt;namespace&gt; get pod -l coherenceCluster=management-cluster

NAME                           READY   STATUS    RESTARTS   AGE
management-cluster-storage-0   1/1     Running   0          36s
management-cluster-storage-1   1/1     Running   0          36s
management-cluster-storage-2   1/1     Running   0          36s</markup>

</div>

<h3 id="_2_port_forward_the_management_over_rest_port">2. Port-forward the Management over ReST port</h3>
<div class="section">
<markup
lang="bash"

>kubectl -n coherence-example port-forward management-cluster-storage-0 30000:30000

Forwarding from [::1]:30000 -&gt; 30000
Forwarding from 127.0.0.1:30000 -&gt; 30000</markup>

</div>

<h3 id="_3_access_the_rest_endpoint">3. Access the REST endpoint</h3>
<div class="section">
<p>Issue the following to access the ReST endpoint:</p>

<markup
lang="bash"

>curl http://127.0.0.1:30000/management/coherence/cluster/ | jq</markup>

<markup
lang="json"

>{
  "links": [
    {
      "rel": "parent",
      "href": "http://127.0.0.1:30000/management/coherence"
    },
    {
      "rel": "self",
      "href": "http://127.0.0.1:30000/management/coherence/cluster/"
    },
    {
      "rel": "canonical",
      "href": "http://127.0.0.1:30000/management/coherence/cluster/"
    },
    {
      "rel": "services",
      "href": "http://127.0.0.1:30000/management/coherence/cluster/services"
    },
    {
      "rel": "caches",
      "href": "http://127.0.0.1:30000/management/coherence/cluster/caches"
    },
    {
      "rel": "members",
      "href": "http://127.0.0.1:30000/management/coherence/cluster/members"
    },
    {
      "rel": "management",
      "href": "http://127.0.0.1:30000/management/coherence/cluster/management"
    },
    {
      "rel": "journal",
      "href": "http://127.0.0.1:30000/management/coherence/cluster/journal"
    },
    {
      "rel": "hotcache",
      "href": "http://127.0.0.1:30000/management/coherence/cluster/hotcache"
    },
    {
      "rel": "reporters",
      "href": "http://127.0.0.1:30000/management/coherence/cluster/reporters"
    },
    {
      "rel": "webApplications",
      "href": "http://127.0.0.1:30000/management/coherence/cluster/webApplications"
    }
  ],
  "clusterSize": 3,
  "membersDeparted": [],
  "memberIds": [
    1,
    2,
    3
  ],
  "oldestMemberId": 1,
  "refreshTime": "2019-10-15T03:55:46.461Z",
  "licenseMode": "Development",
  "localMemberId": 1,
  "version": "12.2.1.4.0",
  "running": true,
  "clusterName": "management-cluster",
  "membersDepartureCount": 0,
  "members": [
    "Member(Id=1, Timestamp=2019-10-15 03:46:15.848, Address=10.1.2.184:36531, MachineId=49519, Location=site:coherence.coherence-example.svc.cluster.local,machine:docker-desktop,process:1,member:management-cluster-storage-1, Role=storage)",
    "Member(Id=2, Timestamp=2019-10-15 03:46:19.405, Address=10.1.2.183:40341, MachineId=49519, Location=site:coherence.coherence-example.svc.cluster.local,machine:docker-desktop,process:1,member:management-cluster-storage-2, Role=storage)",
    "Member(Id=3, Timestamp=2019-10-15 03:46:19.455, Address=10.1.2.185:38719, MachineId=49519, Location=site:coherence.coherence-example.svc.cluster.local,machine:docker-desktop,process:1,member:management-cluster-storage-0, Role=storage)"
  ],
  "type": "Cluster"
}</markup>

<div class="admonition note">
<p class="admonition-inline">The <code>jq</code> utility is used to format the JSON, and may not be available on all platforms.</p>
</div>
</div>

<h3 id="_3_access_the_swagger_endpoint">3. Access the Swagger endpoint</h3>
<div class="section">
<p>Issue the following to access the Sagger endpoint which documents all the API&#8217;s available.</p>

<markup
lang="bash"

>curl http://127.0.0.1:30000/management/coherence/cluster/metadata-catalog | jq</markup>

<markup
lang="json"

>{
  "swagger": "2.0",
  "info": {
    "title": "RESTful Management Interface for Oracle Coherence MBeans",
    "description": "RESTful Management Interface for Oracle Coherence MBeans",
    "version": "12.2.1.4.0"
  },
  "schemes": [
    "http",
    "https"
  ],
...</markup>

<div class="admonition note">
<p class="admonition-inline">The above output has been truncated due to the large size.</p>
</div>
</div>

<h3 id="_4_other_resources">4. Other Resources</h3>
<div class="section">
<p>Management over ReST can be used for all management functions, as one would with
standard MBean access over JMX.</p>

<p>Please see the <a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/fusion-middleware/coherence/12.2.1.4/rest-reference/index.html">Coherence REST API</a> for more information on these features.</p>

<ul class="ulist">
<li>
<p><a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/fusion-middleware/coherence/12.2.1.4/manage/using-jmx-manage-oracle-coherence.html#GUID-D160B16B-7C1B-4641-AE94-3310DF8082EC">Connecting JVisualVM to Management over ReST</a></p>

</li>
<li>
<p><router-link to="/clusters/058_coherence_management">Enabling SSL</router-link></p>

</li>
<li>
<p><a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/fusion-middleware/coherence/12.2.1.4/rest-reference/op-management-coherence-cluster-members-memberidentifier-diagnostic-cmd-jfrcmd-post.html">Produce and extract a Java Flight Recorder (JFR) file</a></p>

</li>
<li>
<p><a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/fusion-middleware/coherence/12.2.1.4/rest-reference/api-reporter.html">Access the Reporter</a></p>

</li>
</ul>
</div>

<h3 id="_5_clean_up">5. Clean Up</h3>
<div class="section">
<p>After running the above the Coherence cluster can be removed using <code>kubectl</code>:</p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; delete -f management-cluster.yaml</markup>

<p>Stop the port-forward command using <code>CTRL-C</code>.</p>

</div>
</div>
</doc-view>
