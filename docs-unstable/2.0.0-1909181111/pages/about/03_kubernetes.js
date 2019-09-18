<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Kubernetes on your Desktop</dt>
<dd slot="desc"><p>For development and testing of the Coherence Operator it&#8217;s often convenient to run Kubernetes on your desktop.</p>

<p>Some ways to do this are:</p>

<ul class="ulist">
<li>
<p><a id="" title="" target="_blank" href="https://docs.docker.com/docker-for-mac/kubernetes/">Kubernetes support in Docker for Desktop</a></p>

</li>
<li>
<p><a id="" title="" target="_blank" href="https://kind.sigs.k8s.io">Kind</a></p>

</li>
<li>
<p><a id="" title="" target="_blank" href="https://kubernetes.io/docs/getting-started-guides/minikube/">Kubernetes Minikube</a></p>

</li>
</ul></dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_docker_for_desktop">Docker for Desktop.</h2>
<div class="section">

<h3 id="_install">Install</h3>
<div class="section">
<p>Install
<a id="" title="" target="_blank" href="https://docs.docker.com/docker-for-mac/install/">Docker for Mac</a> or
<a id="" title="" target="_blank" href="https://docs.docker.com/docker-for-windows/install/">Docker for Windows</a>.</p>

<p>Starting with version 18.06 Docker for Desktop includes Kubernetes support.</p>

</div>

<h3 id="_enable_kubernetes_support">Enable Kubernetes Support</h3>
<div class="section">
<p>Enable
<a id="" title="" target="_blank" href="https://docs.docker.com/docker-for-mac/#kubernetes">Kubernetes Support for Mac</a>
or
<a id="" title="" target="_blank" href="https://docs.docker.com/docker-for-windows/#kubernetes">Kubernetes Support for Windows</a>.</p>

<p>Once Kubernetes installation is complete, make sure you have your context
set correctly to use docker-for-desktop.</p>

<markup
lang="bash"
title="Make sure K8s context is set to docker-for-desktop"
>kubectl config get-contexts
kubectl config use-context docker-for-desktop
kubectl cluster-info
kubectl version --short
kubectl get nodes</markup>

</div>
</div>

<h2 id="_kind">Kind</h2>
<div class="section">
<p>Install Kind as described in the <a id="" title="" target="_blank" href="https://kind.sigs.k8s.io/docs/user/quick-start/">Kind Quick Start</a></p>

<ul class="ulist">
<li>
<p>To create a Kubernetes three node cluster in Kind you need a configuration file</p>

</li>
</ul>
<markup
lang="yaml"
title="kind-config.yaml"
>kind: Cluster
apiVersion: kind.sigs.k8s.io/v1alpha3
nodes:
  - role: control-plane
  - role: worker
  - role: worker
  - role: worker</markup>

<ul class="ulist">
<li>
<p>Then create a Kind Kubernetes cluster with the following command</p>

</li>
</ul>
<markup
lang="bash"

>kind create cluster --config kind-config.yaml</markup>

<ul class="ulist">
<li>
<p>After a short while (depending on how long images take to download) there should be a Kubernetes cluster with a master
and three worker nodes running in Docker containers.</p>

</li>
</ul>
<markup
lang="bash"

>docker ps
d790d6b779ff  kindest/node:v1.15.3   "/usr/local/bin/entr…"   23 hours ago  Up 23 hours                                         kind-worker2
a096c8bf0c1a  kindest/node:v1.15.3   "/usr/local/bin/entr…"   23 hours ago  Up 23 hours                                         kind-worker3
4c01d94c29b7  kindest/node:v1.15.3   "/usr/local/bin/entr…"   23 hours ago  Up 23 hours   56603/tcp, 127.0.0.1:56603-&gt;6443/tcp  kind-control-plane
8f62284be151  kindest/node:v1.15.3   "/usr/local/bin/entr…"   23 hours ago  Up 23 hours                                         kind-worker</markup>

<ul class="ulist">
<li>
<p>As described in the Kind documentation now export <code>KUBECONFIG</code> for the Kind cluster</p>

</li>
</ul>
<markup
lang="bash"

>export KUBECONFIG="$(kind get kubeconfig-path --name="kind")"</markup>

<ul class="ulist">
<li>
<p>To be able to use this cluster with the Operator Helm chart Helm will need to be initialised.</p>

</li>
</ul>
<markup
lang="bash"

>helm init</markup>

<p>The Kind cluster has RBAC enabled so Helm&#8217;s Tiller will now need to be patched with a role:</p>

<markup
lang="bash"

>kubectl create serviceaccount \
  --namespace kube-system tiller

kubectl create clusterrolebinding tiller-cluster-rule \
  --clusterrole=cluster-admin --serviceaccount=kube-system:tiller

kubectl patch deploy --namespace kube-system \
  tiller-deploy -p '{"spec":{"template":{"spec":{"serviceAccount":"tiller"}}}}'</markup>


<h3 id="_a_word_about_kind_and_docker_images">A Word About Kind and Docker Images</h3>
<div class="section">
<p>If trying to use Kind to run the Coherence Operator using locally built images these images need to be added to the
Kubernetes cluster using the Kind CLI because the local images will obviously not be in a repository that the nodes
can pull from. Although a Kind cluster is running in Docker it does not appear to have access to any local Docker images
so all images either need to be pull-able or loaded via the Kind CLI.</p>

<p>For example if the Operator has been built with <code>make all</code> there will be the following local images</p>

<markup
lang="bash"

>docker images --format "table {{.Repository}}\t{{.Tag}}"

REPOSITORY                                                     TAG
iad.ocir.io/odx-stateservice/test/oracle/coherence-operator    2.0.0-ci
iad.ocir.io/odx-stateservice/test/oracle/operator-test-image   2.0.0-ci
iad.ocir.io/odx-stateservice/test/oracle/coherence-operator    2.0.0-ci-utils</markup>

<p>These images can be added to the Kind cluster with the commands:</p>

<markup
lang="bash"

>kind load docker-image iad.ocir.io/odx-stateservice/test/oracle/coherence-operator:2.0.0-ci
kind load docker-image iad.ocir.io/odx-stateservice/test/oracle/coherence-operator:2.0.0-ci-utils
kind load docker-image iad.ocir.io/odx-stateservice/test/oracle/operator-test-image:2.0.0-ci</markup>

</div>
</div>
</doc-view>
