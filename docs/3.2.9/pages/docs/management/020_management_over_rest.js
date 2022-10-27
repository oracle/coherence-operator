<doc-view>

<h2 id="_management_over_rest">Management over REST</h2>
<div class="section">
<p>Since version 12.2.1.4 Coherence has had functionality to expose a management API over REST.</p>

<div class="admonition note">
<p class="admonition-inline">The Management over REST  API is <strong>disabled</strong> by default in Coherence clusters but can be enabled and configured by
setting the relevant fields in the <code>Coherence</code> CRD.</p>
</div>
<p>The example below shows how to enable and access Coherence MBeans using Management over REST.</p>

<p>Once the Management port has been exposed, for example via a load balancer or port-forward command, the REST
endpoint is available at <code><a id="" title="" target="_blank" href="http://host:port/management/coherence/cluster">http://host:port/management/coherence/cluster</a></code>.
The Swagger JSON document for the API is available at <code><a id="" title="" target="_blank" href="http://host:port/management/coherence/cluster/metadata-catalog">http://host:port/management/coherence/cluster/metadata-catalog</a></code>.</p>

<p>See the <a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/standalone/coherence/14.1.1.0/rest-reference/">REST API for Managing Oracle Coherence</a>
documentation for full details on each of the endpoints.</p>

<div class="admonition note">
<p class="admonition-inline">Note: Use of Management over REST is available only when using the operator with clusters running
Coherence 12.2.1.4 or later version.</p>
</div>

<h3 id="_deploy_coherence_with_management_over_rest_enabled">Deploy Coherence with Management over REST Enabled</h3>
<div class="section">
<p>To deploy a <code>Coherence</code> resource with management over REST enabled and exposed on a port, the simplest yaml
would look like this:</p>

<markup
lang="yaml"
title="management-cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: management-cluster
spec:
  coherence:
    management:
      enabled: true     <span class="conum" data-value="1" />
  ports:
    - name: management  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">Setting the <code>coherence.management.enabled</code> field to <code>true</code> will enable Management over REST</li>
<li data-value="2">To expose Management over REST via a <code>Service</code> it is added to the <code>ports</code> list.
The <code>management</code> port is a special case where the <code>port</code> number is optional so in this case Management over REST
will bind to the default port <code>30000</code>.
(see <router-link to="/docs/ports/020_container_ports">Exposing Ports</router-link> for details)</li>
</ul>
<p>To expose Management over REST on a different port the alternative port value can be set in the <code>coherence.management</code>
section, for example:</p>

<markup
lang="yaml"
title="management-cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: management-cluster
spec:
  coherence:
    management:
      enabled: true
      port: 8080      <span class="conum" data-value="1" />
  ports:
    - name: management</markup>

<ul class="colist">
<li data-value="1">Management over REST will now be exposed on port <code>8080</code></li>
</ul>
</div>

<h3 id="_port_forward_the_management_over_rest_port">Port-forward the Management over REST Port</h3>
<div class="section">
<p>After installing the basic <code>management-cluster.yaml</code> from the first example above there would be a three member
Coherence cluster installed into Kubernetes.</p>

<p>For example, the cluster can be installed with <code>kubectl</code></p>

<markup
lang="bash"

>kubectl -n coherence-test create -f management-cluster.yaml

coherence.coherence.oracle.com/management-cluster created</markup>

<p>The <code>kubectl</code> CLI can be used to list <code>Pods</code> for the cluster:</p>

<markup
lang="bash"

>kubectl -n coherence-test get pod -l coherenceCluster=management-cluster

NAME                   READY   STATUS    RESTARTS   AGE
management-cluster-0   1/1     Running   0          36s
management-cluster-1   1/1     Running   0          36s
management-cluster-2   1/1     Running   0          36s</markup>

<p>In a test or development environment the simplest way to reach an exposed port is to use the <code>kubectl port-forward</code> command.
For example to connect to the first <code>Pod</code> in the deployment:</p>

<markup
lang="bash"

>kubectl -n coherence-test port-forward management-cluster-0 30000:30000

Forwarding from [::1]:30000 -&gt; 30000
Forwarding from 127.0.0.1:30000 -&gt; 30000</markup>

</div>

<h3 id="_access_the_rest_endpoint">Access the REST Endpoint</h3>
<div class="section">
<p>Now that a port is being forwarded from localhost to a <code>Pod</code> in the cluster the Management over REST endpoints can be accessed.</p>

<p>Issue the following <code>curl</code> command to access the REST endpoint:</p>

<markup
lang="bash"

>curl http://127.0.0.1:30000/management/coherence/cluster/</markup>

<p>Which should result in a response similar to the following:</p>

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
  "version": "14.1.1.0.0",
  "running": true,
  "clusterName": "management-cluster",
  "membersDepartureCount": 0,
  "members": [
    "Member(Id=1, Timestamp=2019-10-15 03:46:15.848, Address=10.1.2.184:36531, MachineId=49519, Location=site:coherence.coherence-test.svc.cluster.local,machine:docker-desktop,process:1,member:management-cluster-1, Role=storage)",
    "Member(Id=2, Timestamp=2019-10-15 03:46:19.405, Address=10.1.2.183:40341, MachineId=49519, Location=site:coherence.coherence-test.svc.cluster.local,machine:docker-desktop,process:1,member:management-cluster-2, Role=storage)",
    "Member(Id=3, Timestamp=2019-10-15 03:46:19.455, Address=10.1.2.185:38719, MachineId=49519, Location=site:coherence.coherence-test.svc.cluster.local,machine:docker-desktop,process:1,member:management-cluster-0, Role=storage)"
  ],
  "type": "Cluster"
}</markup>

</div>

<h3 id="_access_the_swagger_endpoint">Access the Swagger Endpoint</h3>
<div class="section">
<p>Issue the following <code>curl</code> command to access the Sagger endpoint, which documents all the REST API&#8217;s available.</p>

<markup
lang="bash"

>curl http://127.0.0.1:30000/management/coherence/cluster/metadata-catalog</markup>

<p>Which should result in a response like the following:</p>

<markup
lang="json"

>{
  "swagger": "2.0",
  "info": {
    "title": "RESTful Management Interface for Oracle Coherence MBeans",
    "description": "RESTful Management Interface for Oracle Coherence MBeans",
    "version": "14.1.1.0.0"
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

<h3 id="_other_rest_resources">Other REST Resources</h3>
<div class="section">
<p>Management over REST can be used for all Coherence management functions, the same as would be available when using
standard MBean access over JMX.</p>

<p>Please see the
<a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/standalone/coherence/14.1.1.0/rest-reference/">Coherence REST API</a> for more information on these features.</p>

<ul class="ulist">
<li>
<p><a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/standalone/coherence/14.1.1.0/manage/using-jmx-manage-oracle-coherence.html#GUID-D160B16B-7C1B-4641-AE94-3310DF8082EC">Connecting JVisualVM to Management over REST</a></p>

</li>
<li>
<p><router-link to="#docs/clusters/058_coherence_management.adoc" @click.native="this.scrollFix('#docs/clusters/058_coherence_management.adoc')">Enabling SSL</router-link></p>

</li>
<li>
<p><a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/standalone/coherence/14.1.1.0/rest-reference/op-management-coherence-cluster-members-memberidentifier-diagnostic-cmd-jfrcmd-post.html">Produce and extract a Java Flight Recorder (JFR) file</a></p>

</li>
<li>
<p><a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/standalone/coherence/14.1.1.0/rest-reference/api-reporter.html">Access the Reporter</a></p>

</li>
</ul>
</div>
</div>
</doc-view>
