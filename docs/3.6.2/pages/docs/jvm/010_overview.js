<doc-view>

<h2 id="_overview">Overview</h2>
<div class="section">
<p>The Coherence Operator allows full control over the configuration of the JVM used to run the Coherence application.
The <code>jvm</code> section of the <code>Coherence</code> CRD spec has a number of fields to easily configure specific aspects of the
JVM as well as a catch-all <code>jvm.args</code> list that allows any arbitrary argument to be passed to the JVM.</p>

<p>Whilst every configuration setting could, in theory, be set only by specifying JVM arguments in the <code>jvm.args</code>
field of the <code>Coherence</code> CRD, the other configuration fields provide simpler means to set configuration
without having to remember specific JVM argument names or system property names to set.
You are, of course, free to use whichever approach best suits your requirements;
but obviously it is better to choose one approach and be consistent.</p>


<h3 id="_guides_to_jvm_settings">Guides to JVM Settings</h3>
<div class="section">
<v-layout row wrap class="mb-5">
<v-flex xs12>
<v-container fluid grid-list-md class="pa-0">
<v-layout row wrap class="pillars">
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/docs/jvm/020_classpath"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Classpath</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Default classpath settings and options for setting a custom classpath.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/docs/jvm/030_jvm_args"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">JVM Arguments</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Adding arbitrary JVM arguments and system properties.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/docs/jvm/040_gc"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Garbage Collection</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Configuring the garbage collector.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/docs/jvm/050_memory"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Heap & Memory Settings</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Configuring the heap size and other memory settings.</p>
</v-card-text>
</v-card>
</v-flex>
</v-layout>
</v-container>
</v-flex>
</v-layout>
<v-layout row wrap class="mb-5">
<v-flex xs12>
<v-container fluid grid-list-md class="pa-0">
<v-layout row wrap class="pillars">
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/docs/jvm/070_debugger"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Debugger</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Using debugger settings.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/docs/jvm/090_container_limits"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Use Container Limits</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Configuring the JVM to respect container resource limits.</p>
</v-card-text>
</v-card>
</v-flex>
</v-layout>
</v-container>
</v-flex>
</v-layout>
</div>
</div>
</doc-view>
