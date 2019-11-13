<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Building the Docs</dt>
<dd slot="desc"><p>The Coherence Operator documentation can be built directly from <code>make</code> commands.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_building_the_coherence_operator_documentation">Building the Coherence Operator Documentation</h2>
<div class="section">
<p>The Coherence Operator documentation is written in <a id="" title="" target="_blank" href="https://asciidoctor.org">Ascii Doc</a> format and is built with tools
provided by our friends over in the <a id="" title="" target="_blank" href="http://helidon.io">Helidon</a> team.</p>

<p>The documentation source is under the <code>docs/</code> directory.</p>


<h3 id="_build">Build</h3>
<div class="section">
<p>To build the documentation run</p>

<markup
lang="bash"

>make docs</markup>

<p>This will build the documentation into the directory <code>build/_output/docs</code></p>

</div>

<h3 id="_view">View</h3>
<div class="section">
<p>To see the results of local changes to the documentation it is possible to run a local web-server that will allow the docs
to be viewed in a browser.</p>

<markup
lang="bash"

>make server-docs</markup>

<p>This will start a local web-server on <a id="" title="" target="_blank" href="http://localhost:8080">http://localhost:8080</a>
This is useful to see changes in real time as documentation is edited and re-built.
The server does no need to be restarted between documentation builds.</p>

<div class="admonition note">
<p class="admonition-inline">The local web-server requires Python</p>
</div>
</div>
</div>
</doc-view>
