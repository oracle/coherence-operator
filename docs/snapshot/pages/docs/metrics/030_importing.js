<doc-view>

<h2 id="_importing_the_coherence_dashboards">Importing the Coherence Dashboards</h2>
<div class="section">
<p>The Coherence Operator provides a set of dashboards for Coherence that may be imported into Grafana.
The Coherence dashboards are explained in detail on the <router-link to="#040_dashboards.adoc" @click.native="this.scrollFix('#040_dashboards.adoc')">Coherence Grafana Dashboards</router-link> page.</p>

<p>There are two ways to obtain the dashboards:</p>

<p>1 - Download the <code>.tar.gz</code> dashboards package for the release you want to use.</p>

<markup
lang="bash"

>curl https://oracle.github.io/coherence-operator/dashboards/latest/coherence-dashboards.tar.gz \
    -o coherence-dashboards.tar.gz
tar -zxvf coherence-dashboards.tar.gz</markup>

<p>The above commands will download the <code>coherence-dashboards.tar.gz</code> file and unpack it resulting in a
directory named <code>dashboards/</code> in the current working directory. This <code>dashboards/</code> directory will contain
the various Coherence dashboard files.</p>

<p>2 - Clone the Coherence Operator GitHub repo, checkout the branch or tag for the version you want to use and
then obtain the dashboards from the <code>dashboards/</code> directory.</p>

<div class="admonition note">
<p class="admonition-inline">The recommended versions of Grafana to use are: <code>8.5.6</code> or <code>6.7.4</code>. It is not yet recommended to use the <code>9.x</code> versions of Grafana as there are a number of bugs that cause issues when using the dashboards.</p>
</div>
</div>

<h2 id="_import_the_dashboards_into_grafana">Import the Dashboards into Grafana.</h2>
<div class="section">
<p>This section shows you how to import the Grafana dashboards into your own Grafana instance.
Once you have obtained the dashboards using one of the methods above, the Grafana dashboard <code>.json</code> files will be in the <code>dashboards/grafana/</code> subdirectory</p>

<div class="admonition important">
<p class="admonition-textlabel">Important</p>
<p ><p>By default, the Coherence dashboards require a datasource in Grafana named <code>prometheus</code> (which is case-sensitive).
This datasource usually exists in an out-of-the-box Prometheus Operator installation.
If your Grafana environment does not have this datasource, then there are two choices.</p>

<ul class="ulist">
<li>
<p>Create a Prometheus datasource named <code>prometheus</code> as described in the <a id="" title="" target="_blank" href="https://grafana.com/docs/grafana/latest/datasources/add-a-data-source/">Grafana Add a Datasource</a> documentation and make this the default datasource.</p>

</li>
<li>
<p>If you have an existing Prometheus datasource with a different name then you will need to edit the dashboard json
files to change all occurrences of <code>"datasource": "prometheus"</code> to have the name of your Prometheus datasource.
For example, running the script below in the directory containing the datasource <code>.json</code> files to be imported will
change the datasource name from <code>prometheus</code> to <code>Coherence-Prometheus</code>.</p>

</li>
</ul>
<div class="listing">
<pre>for file in *.json
do
    sed -i '' -e 's/"datasource": "prometheus"/"datasource": "Coherence-Prometheus"/g' $file;
done</pre>
</div>
</p>
</div>

<h3 id="_manually_import_grafana_dashboards">Manually Import Grafana Dashboards</h3>
<div class="section">
<p>The dashboard <code>.json</code> files can be manually imported into Grafana using the Grafana UI following the instructions
in the
<a id="" title="" target="_blank" href="https://grafana.com/docs/grafana/latest/dashboards/manage-dashboards/#import-a-dashboard">Grafana Import Dashboard</a>
documentation.</p>

</div>

<h3 id="_bulk_import_grafana_dashboards">Bulk Import Grafana Dashboards</h3>
<div class="section">
<p>At the time of writing, for whatever reason, Grafana does not provide a simple way to bulk import a set of dashboard files.
There are many examples and scripts on available in the community that show how to do this.
The Coherence Operator source contains a script that can be used for this purpose
<a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/raw/main/hack/grafana-import.sh">grafana-import.sh</a></p>

<div class="admonition note">
<p class="admonition-inline">The <code>grafana-import.sh</code> script requires the <a id="" title="" target="_blank" href="https://stedolan.github.io/jq/">JQ</a> utility to parse json.</p>
</div>
<p>The commands below will download and run the shell script to import the dashboards.
Change the <code>&lt;GRAFANA-USER&gt;</code> and <code>&lt;GRAFANA_PWD&gt;</code> to the Grafana credentials for your environment.
For example if using the default Prometheus Operator installation they are as specified on the
<a id="" title="" target="_blank" href="https://prometheus-operator.dev/docs/prologue/quick-start/#access-grafana">Access Grafana section of the Quick Start</a> page.
We do not document the credentials here as the default values have been known to change between Prometheus Operator and Grafana versions.</p>

<markup
lang="bash"

>curl -Lo grafana-import.sh https://github.com/oracle/coherence-operator/raw/main/hack/grafana-import.sh
chmod +x grafana-import.sh</markup>

<markup
lang="bash"

>./grafana-import.sh -u &lt;GRAFANA-USER&gt; -w &lt;GRAFANA_PWD&gt; -d dashboards/grafana -t localhost:3000</markup>

<p>Note: the command above assumes you can reach Grafana on <code>localhost:3000</code> (for example, if you have a kubectl port forward process
running to forward localhost:3000 to the Grafana service in Kubernetes). You may need to change the host and port to match however
you are exposing your Grafana instance.</p>

<p>Coherence clusters can now be created as described in the <router-link to="/docs/metrics/020_metrics">Publish Metrics</router-link>
page, and metrics will eventually appear in Prometheus and Grafana. It can sometimes take a minute or so for
Prometheus to start scraping metrics and for them to appear in Grafana.</p>

</div>
</div>
</doc-view>
