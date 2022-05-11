<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Coherence Grafana Dashboards</dt>
<dd slot="desc"><p>The Coherence Operator provides detailed Grafana dashboards to provide insight into your running Coherence Clusters.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_coherence_grafana_dashboards">Coherence Grafana Dashboards</h2>
<div class="section">
<div class="admonition note">
<p class="admonition-inline">Note: The Grafana dashboards require Coherence metrics, which is available only when using Coherence version 12.2.1.4.</p>
</div>
</div>

<h2 id="_table_of_contents">Table of Contents</h2>
<div class="section">
<ol style="margin-left: 15px;">
<li>
<router-link to="#navigation" @click.native="this.scrollFix('#navigation')">Navigation</router-link>

</li>
<li>
<router-link to="#dashboards" @click.native="this.scrollFix('#dashboards')">Dashboards</router-link>
<ol style="margin-left: 15px;">
<li>
<router-link to="#main" @click.native="this.scrollFix('#main')">Coherence Dashboard Main</router-link>

</li>
<li>
<router-link to="#members" @click.native="this.scrollFix('#members')">Members Summary &amp; Details Dashboards</router-link>

</li>
<li>
<router-link to="#services" @click.native="this.scrollFix('#services')">Services Summary &amp; Details Dashboards</router-link>

</li>
<li>
<router-link to="#caches" @click.native="this.scrollFix('#caches')">Caches Summary &amp; Detail Dashboards</router-link>

</li>
<li>
<router-link to="#proxies" @click.native="this.scrollFix('#proxies')">Proxy Servers Summary &amp; Detail Dashboards</router-link>

</li>
<li>
<router-link to="#persistence" @click.native="this.scrollFix('#persistence')">Persistence Summary Dashboard</router-link>

</li>
<li>
<router-link to="#federation" @click.native="this.scrollFix('#federation')">Federation Summary &amp; Details Dashboards</router-link>

</li>
<li>
<router-link to="#machines" @click.native="this.scrollFix('#machines')">Machines Summary Dashboard</router-link>

</li>
<li>
<router-link to="#http" @click.native="this.scrollFix('#http')">HTTP Servers Summary Dashboard</router-link>

</li>
<li>
<router-link to="#ed" @click.native="this.scrollFix('#ed')">Elastic Data Summary Dashboard</router-link>

</li>
<li>
<router-link to="#executors" @click.native="this.scrollFix('#executors')">Executors Summary &amp; Details Dashboards</router-link>

</li>
</ol>
</li>
</ol>
</div>

<h2 id="navigation">Navigation</h2>
<div class="section">
<p>The pre-loaded Coherence Dashboards provide a number of common features and
navigation capabilities that appear at the top of most dashboards.</p>


<h3 id="_variables">Variables</h3>
<div class="section">
<p><img src="./images/grafana-variables.png" alt="Variables"width="250" />
</p>

<p>Allows for selection of information to be displayed where there is more than one item.</p>

<ol style="margin-left: 15px;">
<li>
Cluster Name - Allows selection of the cluster to view metrics for

</li>
<li>
Top N Limit - Limits the display of <code>Top</code> values for tables that support it

</li>
<li>
Service Name, Member Name, Cache Name - These will appear on various dashboards

</li>
</ol>
<p>See the <a id="" title="" target="_blank" href="https://grafana.com/docs/reference/templating/">Grafana Documentation</a> for more information on Variables.</p>

</div>

<h3 id="_annotations">Annotations</h3>
<div class="section">
<p><img src="./images/grafana-annotations.png" alt="Annotations"width="250" />
</p>

<p>Vertical red lines on a graph to indicate a change in a key markers such as:</p>

<ol style="margin-left: 15px;">
<li>
Show Cluster Size Changes - Displays when the cluster size has changed

</li>
<li>
Show Partition Transfers - Displays when partition transfers have occurred

</li>
</ol>
<p>See the <a id="" title="" target="_blank" href="https://grafana.com/docs/reference/annotations/">Grafana Documentation</a> for more information on Annotations.</p>

</div>

<h3 id="_navigation">Navigation</h3>
<div class="section">
<p><img src="./images/grafana-navigation.png" alt="Navigation"width="250" />
</p>

<ol style="margin-left: 15px;">
<li>
Select Dashboard - In the top right a drop down list of dashboards is available selection

</li>
<li>
Drill Through - Ability to drill through based upon service, member, node, etc.

</li>
</ol>
</div>
</div>

<h2 id="dashboards">Dashboards</h2>
<div class="section">

<h3 id="main">1. Coherence Dashboard Main</h3>
<div class="section">
<p>Shows a high-level overview of the selected Coherence cluster including metrics such as:</p>

<ul class="ulist">
<li>
<p>Cluster member count, services, memory and health</p>

</li>
<li>
<p>Top N loaded members, Top N heap usage and GC activity</p>

</li>
<li>
<p>Service backlogs and endangered or vulnerable services</p>

</li>
<li>
<p>Top query times, non-optimized queries</p>

</li>
<li>
<p>Guardian recoveries and terminations</p>

</li>
</ul>


<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/grafana-main.png" alt="Dashboard Main"width="950" />
</v-card-text>
</v-card>

</div>

<h3 id="members">2. Members Summary &amp; Details Dashboards</h3>
<div class="section">
<p>Shows an overview of all cluster members that are enabled for metrics capture including metrics such as:</p>

<ul class="ulist">
<li>
<p>Member list include heap usage</p>

</li>
<li>
<p>Top N members for GC time and count</p>

</li>
<li>
<p>Total GC collection count and time by Member</p>

</li>
<li>
<p>Publisher and Receiver success rates</p>

</li>
<li>
<p>Guardian recoveries and send queue size</p>

</li>
</ul>

<h4 id="_members_summary">Members Summary</h4>
<div class="section">


<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/grafana-members.png" alt="Members"width="950" />
</v-card-text>
</v-card>

</div>

<h4 id="_member_details">Member Details</h4>
<div class="section">


<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/grafana-members.png" alt="Member Details"width="950" />
</v-card-text>
</v-card>

</div>
</div>

<h3 id="services">3. Services Summary &amp; Details Dashboards</h3>
<div class="section">
<p>Shows an overview of all cluster services including metrics such as:</p>

<ul class="ulist">
<li>
<p>Service members for storage and non-storage services</p>

</li>
<li>
<p>Service task count</p>

</li>
<li>
<p>StatusHA values as well as endangered, vulnerable and unbalanced partitions</p>

</li>
<li>
<p>Top N services by task count and backlog</p>

</li>
<li>
<p>Task rates, request pending counts and task and request averages</p>

</li>
</ul>

<h4 id="_services_summary">Services Summary</h4>
<div class="section">


<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/grafana-services.png" alt="Services"width="950" />
</v-card-text>
</v-card>

</div>

<h4 id="_service_details">Service Details</h4>
<div class="section">


<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/grafana-service.png" alt="Service Details"width="950" />
</v-card-text>
</v-card>

</div>
</div>

<h3 id="caches">4. Caches Summary &amp; Detail Dashboards</h3>
<div class="section">
<p>Shows an overview of all caches including metrics such as:</p>

<ul class="ulist">
<li>
<p>Cache entries, memory and index usage</p>

</li>
<li>
<p>Cache access counts including gets, puts and removed,  max query times</p>

</li>
<li>
<p>Front cache hit and miss rates</p>

</li>
</ul>

<h4 id="_caches_summary">Caches Summary</h4>
<div class="section">


<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/grafana-caches.png" alt="Caches"width="950" />
</v-card-text>
</v-card>

</div>

<h4 id="_cache_details">Cache Details</h4>
<div class="section">


<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/grafana-cache.png" alt="Cache Details"width="950" />
</v-card-text>
</v-card>

</div>
</div>

<h3 id="proxies">5. Proxy Servers Summary &amp; Detail Dashboards</h3>
<div class="section">
<p>Shows and overview of Proxy servers including metrics such as:</p>

<ul class="ulist">
<li>
<p>Active connection count and service member count</p>

</li>
<li>
<p>Total messages sent/ received</p>

</li>
<li>
<p>Proxy server data rates</p>

</li>
<li>
<p>Individual connection details abd byte backlogs</p>

</li>
</ul>

<h4 id="_proxy_servers_summary">Proxy Servers Summary</h4>
<div class="section">


<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/grafana-proxies.png" alt="Proxy Servers"width="950" />
</v-card-text>
</v-card>

</div>

<h4 id="_proxy_servers_detail">Proxy Servers Detail</h4>
<div class="section">


<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/grafana-proxy.png" alt="Proxy Server Details"width="950" />
</v-card-text>
</v-card>

</div>
</div>

<h3 id="persistence">6. Persistence Summary Dashboard</h3>
<div class="section">
<p>Shows and overview of Persistence including metrics such as:</p>

<ul class="ulist">
<li>
<p>Persistence enabled services</p>

</li>
<li>
<p>Maximum active persistence latency</p>

</li>
<li>
<p>Active space total usage and by service</p>

</li>
</ul>


<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/grafana-persistence.png" alt="Persistence"width="950" />
</v-card-text>
</v-card>

</div>

<h3 id="federation">7. Federation Summary &amp; Details Dashboards</h3>
<div class="section">
<p>Shows overview of Federation including metrics such as:</p>

<ul class="ulist">
<li>
<p>Destination and Origins details</p>

</li>
<li>
<p>Entries, records and bytes send and received</p>

</li>
</ul>

<h4 id="_federation_summary">Federation Summary</h4>
<div class="section">


<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/grafana-federation-summary.png" alt="Federation Summary"width="950" />
</v-card-text>
</v-card>

</div>

<h4 id="_federation_details">Federation Details</h4>
<div class="section">


<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/grafana-federation-detail.png" alt="Federation Details"width="950" />
</v-card-text>
</v-card>

</div>
</div>

<h3 id="machines">8. Machines Summary Dashboard</h3>
<div class="section">
<p>Shows an overview of all machines that make up the Kubernetes cluster underlying the Coherence cluster including metrics such as:</p>

<ul class="ulist">
<li>
<p>Machine processors, free swap space and physical memory</p>

</li>
<li>
<p>Load averages</p>

</li>
</ul>


<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/grafana-machines.png" alt="Machines"width="950" />
</v-card-text>
</v-card>

</div>

<h3 id="http">9. HTTP Servers Summary Dashboard</h3>
<div class="section">
<p>Shows an overview of all HTTP Servers running in the cluster including metrics such as:</p>

<ul class="ulist">
<li>
<p>Service member count, requests, error count and average request time</p>

</li>
<li>
<p>HTTP Request rates and response codes</p>

</li>
</ul>


<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/grafana-http.png" alt="HTTP Servers"width="950" />
</v-card-text>
</v-card>

</div>

<h3 id="ed">10. Elastic Data Summary Dashboard</h3>
<div class="section">
<p>Shows an overview of all HTTP Servers running in the cluster including metrics such as:</p>

<ul class="ulist">
<li>
<p>RAM and Flash journal files in use</p>

</li>
<li>
<p>RAM and Flash compactions</p>

</li>
</ul>


<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/grafana-elastic-data.png" alt="Elastic Data"width="950" />
</v-card-text>
</v-card>

</div>

<h3 id="executors">11. Executors Summary &amp; Details Dashboards</h3>
<div class="section">
<p>Shows an overview of all Executors running in the cluster including metrics such as:</p>

<ul class="ulist">
<li>
<p>Tasks in Progress</p>

</li>
<li>
<p>Completed and Rejected Tasks</p>

</li>
<li>
<p>Individual Executor status</p>

</li>
</ul>

<h4 id="_executors_summary">Executors Summary</h4>
<div class="section">


<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/grafana-executors-summary.png" alt="Executors Summary"width="950" />
</v-card-text>
</v-card>

</div>

<h4 id="_executor_details">Executor Details</h4>
<div class="section">


<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/grafana-executor-detail.png" alt="Executor Detail"width="950" />
</v-card-text>
</v-card>

</div>
</div>
</div>
</doc-view>
