<doc-view>

<h2 id="_overview">Overview</h2>
<div class="section">
<p>A typical Coherence deployment contains custom application code that runs with Coherence.
To run custom application code in a <code>Coherence</code> resource that code needs to be packaged into an image that the
deployment will use.</p>


<h3 id="_building_and_deploying_applications">Building and Deploying Applications</h3>
<div class="section">
<v-layout row wrap class="mb-5">
<v-flex xs12>
<v-container fluid grid-list-md class="pa-0">
<v-layout row wrap class="pillars">
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/docs/applications/020_build_application"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Build Custom Application Images</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Building custom Coherence application images for use with the Coherence Operator.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/docs/applications/030_deploy_application"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Deploy Custom Application Images</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Deploying custom application images using the Coherence Operator.</p>
</v-card-text>
</v-card>
</v-flex>
</v-layout>
</v-container>
</v-flex>
</v-layout>
</div>

<h3 id="_configuring_applications">Configuring Applications</h3>
<div class="section">
<p>There are many settings in a <code>Coherence</code> resource that control the behaviour of Coherence, the JVM and
the application code. Some of the application specific settings are shown below:</p>

<v-layout row wrap class="mb-5">
<v-flex xs12>
<v-container fluid grid-list-md class="pa-0">
<v-layout row wrap class="pillars">
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/docs/jvm/020_classpath"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Setting the Classpath</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Setting a custom classpath for the application.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/docs/applications/040_application_main"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Setting a Main Class</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Setting a custom main class to run.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/docs/applications/050_application_args"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Setting Application Arguments</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Setting arguments to pass to the main class.</p>
</v-card-text>
</v-card>
</v-flex>
<v-flex xs12 sm4 lg3>
<v-card>
<router-link to="/docs/applications/060_application_working_dir"><div class="card__link-hover"/>
</router-link>
<v-card-title primary class="headline layout justify-center">
<span style="text-align:center">Working Directory</span>
</v-card-title>
<v-card-text class="caption">
<p></p>
<p>Setting the application&#8217;s working directory.</p>
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
