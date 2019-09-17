<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Kubernetes on your Desktop</dt>
<dd slot="desc"><p>For development and testing of the Coherence Operator it&#8217;s often convenient to run Kubernetes on your desktop.
Two popular ways to do this are with
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/getting-started-guides/minikube/">Kubernetes Minikube</a>
or
<a id="" title="" target="_blank" href="https://docs.docker.com/docker-for-mac/kubernetes/">Kubernetes support in Docker for Desktop</a>.</p>

<p>In this guide we&#8217;ll use Docker for Desktop.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_install">Install</h2>
<div class="section">
<p>Install
<a id="" title="" target="_blank" href="https://docs.docker.com/docker-for-mac/install/">Docker for Mac</a> or
<a id="" title="" target="_blank" href="https://docs.docker.com/docker-for-windows/install/">Docker for Windows</a>.</p>

<p>Starting with version 18.06 Docker for Desktop includes Kubernetes support.</p>

</div>

<h2 id="_enable_kubernetes_support">Enable Kubernetes Support</h2>
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
</doc-view>
