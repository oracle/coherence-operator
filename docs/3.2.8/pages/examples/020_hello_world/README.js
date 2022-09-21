<doc-view>

<h2 id="_a_hello_world_operator_example">A "Hello World" Operator Example</h2>
<div class="section">
<p>This is the most basic example of how to deploy a simple Coherence cluster to Kubernetes using the Coherence Operator.</p>

<div class="admonition tip">
<p class="admonition-textlabel">Tip</p>
<p ><p><img src="./images/GitHub-Mark-32px.png" alt="GitHub Mark 32px" />
 The complete source code for this example is in the <a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/tree/main/examples/020_hello_world">Coherence Operator GitHub</a> repository.</p>
</p>
</div>

<h3 id="_install_the_operator">Install the Operator</h3>
<div class="section">
<p>If you have not already done so, you need to install the Coherence Operator.
There are a few simple ways to do this as described in the <router-link to="/docs/installation/01_installation">Installation Guide</router-link></p>

</div>

<h3 id="_a_default_coherence_cluster">A Default Coherence Cluster</h3>
<div class="section">
<p>All the fields in the Coherence CRD spec are optional, the Operator will apply default values, if required, for fields not specified.</p>

<p>For example, this is the minimum required yaml to run a Coherence cluster:</p>

<markup
lang="yaml"
title="default-coherence.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test</markup>

<p>The yaml above could be installed into Kubernetes using kubectl:</p>

<markup
lang="bash"

>kubectl create -f default-coherence.yaml</markup>

<p>The command above will create a Coherence cluster named <code>test</code> in the <code>default</code> Kubernetes namespace.</p>

<p>Because no <code>spec</code> was specified in the yaml, the Operator will use its defaults for certain fields.</p>

<ul class="ulist">
<li>
<p>The <code>replicas</code> field, which controls the number of Pods in the cluster, will default to <code>3</code>.</p>

</li>
<li>
<p>The image used to run Coherence will be the default for this version of the Operator,
typically this is the latest Coherence CE image released at the time the Operator version was released.</p>

</li>
<li>
<p>No ports will be exposed on the container, and no additional services will be created.</p>

</li>
</ul>
<p>We can list the resources that have been created by the Operator.</p>

<markup
lang="bash"

>kubectl get all</markup>

<p>Which should display something like this:</p>

<markup
lang="bash"

>NAME         READY   STATUS    RESTARTS   AGE
pod/test-0   1/1     Running   0          81s
pod/test-1   1/1     Running   0          81s
pod/test-2   1/1     Running   0          81s

NAME                 TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
service/test-sts     ClusterIP   None         &lt;none&gt;        7/TCP     81s
service/test-wka     ClusterIP   None         &lt;none&gt;        7/TCP     81s

NAME                    READY   AGE
statefulset.apps/test   3/3     81s</markup>

<ul class="ulist">
<li>
<p>We can see that the Operator has created a <code>StatefulSet</code>, with three <code>Pods</code> and there are two <code>Services</code>.</p>

</li>
<li>
<p>The <code>test-sts</code> service is the headless service required for the <code>StatefulSet</code>.</p>

</li>
<li>
<p>The <code>test-wka</code> service is the headless service that Coherence will use for well known address cluster discovery.</p>

</li>
</ul>
<p>We can now undeploy the cluster:</p>

<markup
lang="bash"

>kubectl delete -f default-coherence.yaml</markup>

</div>

<h3 id="_deploy_the_simple_server_image">Deploy the Simple Server Image</h3>
<div class="section">
<p>We can deploy a specific image by setting the <code>spec.image</code> field in the yaml.
In this example we&#8217;ll deploy the <code>simple-coherence:1.0.0</code> image built in the
<router-link to="/examples/015_simple_image/README">Build a Coherence Server Image</router-link> example.</p>

<p>To deploy a specific image we just need to set the <code>spec.image</code> field.</p>

<markup
lang="yaml"
title="simple.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: simple
spec:
  image: simple-coherence:1.0.0  <span class="conum" data-value="1" />
  replicas: 6                    <span class="conum" data-value="2" />
  ports:
    - name: extend               <span class="conum" data-value="3" />
      port: 20000</markup>

<ul class="colist">
<li data-value="1">We have set the image to use to the <router-link to="/examples/015_simple_image/README">Build a Coherence Server Image</router-link> example <code>simple-coherence:1.0.0</code>.</li>
<li data-value="2">We have set the <code>replicas</code> field to <code>6</code>, so this time there should only be six Pods.</li>
<li data-value="3">The simple image starts a Coherence Extend proxy on port <code>20000</code>, so we expose this port in the <code>Coherence</code> spec. The Operator will then expose the port on the Coherence container and create a Service for the port.</li>
</ul>
<p>We can deploy the simple cluster into Kubernetes using kubectl:</p>

<markup
lang="bash"

>kubectl create -f simple.yaml</markup>

<p>Now list the resources the Operator has created.</p>

<markup
lang="bash"

>kubectl get all</markup>

<p>Which this time should look something like this:</p>

<markup
lang="bash"

>NAME         READY   STATUS    RESTARTS   AGE
pod/test-0   1/1     Running   0          4m49s
pod/test-1   1/1     Running   0          4m49s
pod/test-2   1/1     Running   0          4m49s
pod/test-3   1/1     Running   0          4m49s
pod/test-4   1/1     Running   0          4m49s
pod/test-5   1/1     Running   0          4m49s

NAME                  TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)     AGE
service/kubernetes    ClusterIP   10.96.0.1        &lt;none&gt;        443/TCP     164d
service/test-extend   ClusterIP   10.108.166.193   &lt;none&gt;        20000/TCP   4m49s
service/test-sts      ClusterIP   None             &lt;none&gt;        7/TCP       4m49s
service/test-wka      ClusterIP   None             &lt;none&gt;        7/TCP       4m49s

NAME                    READY   AGE
statefulset.apps/test   6/6     4m49s</markup>

<ul class="ulist">
<li>
<p>We can see that the Operator has created a <code>StatefulSet</code>, with six <code>Pods</code> and there are three <code>Services</code>.</p>

</li>
<li>
<p>The <code>simple-sts</code> service is the headless service required for the <code>StatefulSet</code>.</p>

</li>
<li>
<p>The <code>simple-wka</code> service is the headless service that Coherence will use for well known address cluster discovery.</p>

</li>
<li>
<p>The <code>simple-extend</code> service is the service that exposes the Extend port <code>20000</code>, and could be used by Extend clients to connect to the cluster.</p>

</li>
</ul>
<p>We can now delete the simple cluster:</p>

<markup
lang="bash"

>kubectl delete -f simple.yaml</markup>

</div>
</div>
</doc-view>
