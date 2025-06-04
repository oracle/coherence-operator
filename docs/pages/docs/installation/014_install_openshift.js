<doc-view>

<h2 id="_install_on_openshift">Install on OpenShift</h2>
<div class="section">
<p>The Coherence Operator can be installed into OpenShift using either the web console or
manually using yaml manifests.</p>

<p>The <a id="" title="" target="_blank" href="https://docs.redhat.com/en/documentation/openshift_container_platform/4.18/html/operators/index">OpenShift documentation</a>
covers operators in great detail along with how to install them.
It is advisable to check this documentation for the version of OpenShift being used.</p>

<p>There are two methods to install an operator in OpenShift.</p>

<ul class="ulist">
<li>
<p><router-link to="#manual" @click.native="this.scrollFix('#manual')">Manually</router-link> using a subscription yaml</p>

</li>
<li>
<p><router-link to="#console" @click.native="this.scrollFix('#console')">Automatically</router-link> from the OpenShift web console</p>

</li>
</ul>

<h3 id="manual">Manual Installation</h3>
<div class="section">
<p>Create a subscription yaml</p>

<markup

title="coherence-operator-subscription.yaml"
>apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: coherence-operator
  namespace: openshift-operators
spec:
  channel: stable
  installPlanApproval: Automatic
  name: coherence-operator
  source: coherence-operator-catalog
  sourceNamespace: openshift-marketplace
  startingCSV: coherence-operator.v3.5.1</markup>

<p>Apply the subscription yaml:</p>

<markup
lang="bash"

>oc apply -f coherence-operator-subscription.yaml</markup>

<p>The Coherence Operator will be installed into the <code>openshift-operators</code> namespace.</p>

<p>The Coherence operator pods have a label <code>app.kubernetes.io/name=coherence-operator</code> and can be listed
with the following command:</p>

<markup
lang="bash"

>oc -n openshift-operators get pod -l app.kubernetes.io/name=coherence-operator</markup>

<markup


>NAME                                                     READY   STATUS    RESTARTS      AGE
coherence-operator-controller-manager-859675d947-llmvd   1/1     Running   1 (14h ago)   14h
coherence-operator-controller-manager-859675d947-mk765   1/1     Running   3 (14h ago)   14h
coherence-operator-controller-manager-859675d947-z5m2x   1/1     Running   2 (14h ago)   14h</markup>

</div>

<h3 id="console">Install from the Web Console</h3>
<div class="section">
<p>Using the OpenShift console, the Coherence operator can be installed in a few clicks.</p>

<p>In the web-console, expand "Operators" on the left-hand menu, select "OperatorHub" and then type "coherence"
into the search box. The Coherence Operator panel should be displayed.</p>



<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/openshift-operatorhub-coherence.png" alt="OpenShift OperatorHub"width="1024" />
</v-card-text>
</v-card>

<p>Click on the Coherence Operator panel, which will display the Coherence Operator install page.</p>



<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/openshift-coherence.png" alt="OpenShift Coherence Operator"width="1024" />
</v-card-text>
</v-card>

<p>Typically, the latest version will be installed so click on the "Install" button which will display the installation
options panel.</p>



<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/openshift-coherence-install.png" alt="OpenShift Coherence Operator Install"width="1024" />
</v-card-text>
</v-card>

<p>Click on the "Install" button to start the installation.
The installation progress will be displayed.</p>



<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/openshift-coherence-install-progress.png" alt="OpenShift Coherence Install Progress"width="1024" />
</v-card-text>
</v-card>

<p>The display will change to show when installation is complete.</p>



<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/openshift-coherence-install-done.png" alt="OpenShift Coherence Install Complete"width="1024" />
</v-card-text>
</v-card>

<p>Click on the "View Operator" button to see the details page for the Coherence Operator installation</p>



<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/openshift-coherence-operator-details.png" alt="OpenShift Coherence Details"width="1024" />
</v-card-text>
</v-card>

<p>The Coherence Operator is now installed and ready to manage Coherence workloads.</p>

</div>
</div>
</doc-view>
