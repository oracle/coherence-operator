<doc-view>

<h2 id="_capture_heap_dumps">Capture Heap Dumps</h2>
<div class="section">
<p>Heap dumps can be very useful when debugging but generating and downloading a heap dump from a container in Kubernetes can be tricky. When you are running minimal images without an O/S or full JDK (such as the distroless images used by JIB) this becomes even more tricky.</p>

</div>

<h2 id="_ephemeral_containers">Ephemeral Containers</h2>
<div class="section">
<p>Ephemeral containers were introduced in Kubernetes v1.16 and moved to beta in v1.23.
Ephemeral containers is a feature gate that must be enabled for your cluster.
If you have the <code>EphemeralContainers</code> feature gate enabled, then obtaining a heap dump is not so difficult.</p>


<h3 id="_enable_ephemeralcontainers_in_kind">Enable EphemeralContainers in KinD</h3>
<div class="section">
<p>We use <a id="" title="" target="_blank" href="https://kind.sigs.k8s.io">KinD</a> for a lot of our CI builds and testing, enabling the <code>EphemeralContainers</code> feature gate in KinD is very easy.</p>

<p>For example, this KinD configuration enables the <code>EphemeralContainers</code> feature gate</p>

<markup
lang="yaml"

>kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
featureGates:
  EphemeralContainers: true <span class="conum" data-value="1" />
nodes:
- role: control-plane
- role: worker
- role: worker</markup>

<ul class="colist">
<li data-value="1">The <code>EphemeralContainers</code> feature gate is set to <code>true</code></li>
</ul>
</div>

<h3 id="_shared_process_namespace">Shared Process Namespace</h3>
<div class="section">
<p>In this example we are going to use the <code>jps</code> and <code>jcmd</code> tools to generate the heap dump from an ephemeral container.
For this to work the ephemeral container must be able to see the processes running in the <code>coherence</code> container.
The <code>Coherence</code> CRD spec has a field named <code>ShareProcessNamespace</code>, which sets the corresponding field in the Coherence Pods that will be created for the deployment.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test-cluster
spec:
  shareProcessNamespace: true   <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The <code>shareProcessNamespace</code> must be set to <code>true</code>.</li>
</ul>
<p>If you have some other way to trigger a heap dump to a specific location without requiring the ephemeral container to see the Coherence container processes then the technique below can still be used without setting <code>shareProcessNamespace</code> to <code>true</code>.</p>

</div>

<h3 id="_create_an_ephemeral_container">Create an Ephemeral Container</h3>
<div class="section">
<p>Let&#8217;s say we have a Coherence cluster deployed named <code>test-cluster</code> in a namespace named <code>coherence-test</code>.
There will be a number of Pods created for this deployment, named <code>test-cluster-0</code>, <code>test-cluster-1</code> and so on.
For this example we will obtain a heap dump from Pod <code>test-cluster-1</code>.</p>

<p>The purpose of using an ephemeral container is because the Coherence container we are running does not contain any of the tools and programs we require for debugging, e.g. <code>jps</code>, <code>jcmd</code> etc.
The ephemeral container we run obviously needs to have all the required tools. You could create a custom image with what you need in it, but for this example we will use the <code>openjdk:11</code> image, as it has a full JDK and other tools we need in it.
You should obviously use a JDK version that matches the version in the Coherence container.</p>

<p>We can use the <code>kubectl debug</code> command that can be used to create an ephemeral containers.
For our purposes we cannot use this command as we will require volume mounts to share storage between the ephemeral container and the Coherence container so that the ephemeral container can see the heap dump file.</p>

<p>Instead of the <code>kubectl debug</code> command we can create ephemeral containers using the <code>kubectl --raw</code> API.
Ephemeral containers are a sub-resource of the Pod API.</p>

<ul class="ulist">
<li>
<p>First obtain the current ephemeral containers sub-resource for the Pod.
We do this using the <code>kubectl get --raw</code> command with the URL path in the format <code>/api/v1/namespaces/&lt;namespace&gt;&gt;/pods/&lt;pod&gt;/ephemeralcontainers</code>, where <code>&lt;namespace&gt;</code> is the namespace that the Pod is deployed into and <code>&lt;pod&gt;</code> is the name of the Pod.</p>

</li>
</ul>
<p>So in our example the command would be:</p>

<markup
lang="bash"

>kubectl get --raw /api/v1/namespaces/coherence-test/pods/test-cluster-1/ephemeralcontainers</markup>

<p>Which will output json similar to this, which we will save to a file named <code>ec.json</code>:</p>

<markup
lang="json"
title="ec.json"
>{
  "kind": "EphemeralContainers",
  "apiVersion": "v1",
  "metadata": {
    "name": "test-cluster-1",
    "namespace": "coherence-test",
    "selfLink": "/api/v1/namespaces/coherence-test/pods/test-cluster-1/ephemeralcontainers",
    "uid": "731ca9a9-332f-4999-821d-adfea2e1d2d4",
    "resourceVersion": "24921",
    "creationTimestamp": "2021-03-12T10:41:35Z"
  },
  "ephemeralContainers": []
}</markup>

<p>The <code>"ephemeralContainers"</code> field is an empty array as we have not created any previous containers.</p>

<p>We now need to edit this yaml to define the ephemeral container we want to create.
The Pod created by the Operator contains an empty directory volume with a volume mount at <code>/coherence-operator/jvm</code>, which is where the JVM is configured to dump debug information, such as heap dumps.
We will create an ephemeral container with the same mount so that the <code>/coherence-operator/jvm</code> directory will be shared between the Coherence container and the ephemeral container.</p>

<p>Another thing to note is that the default entrypoint in the <code>openjdk:11</code> image we are using in this example is JShell.
This is obviously not what we want, so we will make sure we specify <code>/bin/sh</code> as the entry point as we want a command line shell.</p>

<p>Our edited <code>ec.json</code> file looks like this:</p>

<markup
lang="json"
title="ec.json"
>{
  "kind": "EphemeralContainers",
  "apiVersion": "v1",
  "metadata": {
    "name": "test-cluster-1",
    "namespace": "coherence-test",
    "selfLink": "/api/v1/namespaces/coherence-test/pods/test-cluster-1/ephemeralcontainers",
    "uid": "731ca9a9-332f-4999-821d-adfea2e1d2d4",
    "resourceVersion": "24921",
    "creationTimestamp": "2021-03-12T10:41:35Z"
  },
  "ephemeralContainers": [
    {
      "name": "debug",                                 <span class="conum" data-value="1" />
      "image": "openjdk:11",                           <span class="conum" data-value="2" />
      "command": [
          "bin/sh"                                     <span class="conum" data-value="3" />
      ],
      "imagePullPolicy": "IfNotPresent",               <span class="conum" data-value="4" />
      "terminationMessagePolicy":"File",
      "stdin": true,                                   <span class="conum" data-value="5" />
      "tty": true,
      "volumeMounts": [
          {
              "mountPath": "/coherence-operator/jvm",  <span class="conum" data-value="6" />
              "name": "jvm"
          }
      ]
    }
  ]
}</markup>

<ul class="colist">
<li data-value="1">We add an ephemeral container named <code>debug</code>. The name can be anything as long as it is unique in the Pod.</li>
<li data-value="2">We specify that the image used for the container is <code>openjdk:11</code></li>
<li data-value="3">Specify <code>/bin/sh</code> as the container entry point so that we get a command line shell</li>
<li data-value="4">We must specify an image pull policy</li>
<li data-value="5">We want an interactive container, so we specify <code>stdin</code> and <code>tty</code></li>
<li data-value="6">We create the same volume mount to <code>/coherence-operator/jvm</code> that the Coherence container has.</li>
</ul>
<p>We can now re-apply the json to add the new ephemeral container using the <code>kubectl replace --raw</code> command to the same URL path we used for the <code>get</code> command above, this time using <code>-f ec.json</code> to specify the json we want to replace.</p>

<markup
lang="bash"

>kubectl replace --raw /api/v1/namespaces/coherence-test/pods/test-cluster-1/ephemeralcontainers -f ec.json</markup>

<p>After executing the above command the ephemeral container should have been created, we can now attach to it.</p>

</div>

<h3 id="_attach_to_the_ephemeral_container">Attach to the Ephemeral Container</h3>
<div class="section">
<p>We now have an ephemeral container named <code>debug</code> in the Pod <code>test-cluster-1</code>.
We need to attach to the container so that we can create the heap dump.</p>

<markup
lang="bash"

>kubectl attach test-cluster-1 -c debug -it -n coherence-test</markup>

<p>The command above will attach an interactive (<code>-it</code>) session to the <code>debug</code> container (specified with <code>-c debug</code>) in Pod <code>test-cluster-1</code>, in the namespace <code>coherence-test</code>.
Displaying something like this:</p>

<markup
lang="bash"

>If you don't see a command prompt, try pressing enter.

#</markup>

</div>

<h3 id="_trigger_the_heap_dump">Trigger the Heap Dump</h3>
<div class="section">
<p>We can now generate the heap dump for the Coherence process using <code>jcmd</code>, but first we need to find its PID using <code>jps</code>.</p>

<markup
lang="bash"

>jps -l</markup>

<p>Which will display something like this:</p>

<markup
lang="bash"

>117 jdk.jcmd/sun.tools.jps.Jps
55 com.oracle.coherence.k8s.Main</markup>

<p>The main class run by the Operator is <code>com.oracle.coherence.k8s.Main</code> so the PID of the Coherence process is <code>55</code>.
We can now use <code>jcmd</code> to generate the heap dump. We need to make sure that the heap dump is created in the <code>/coherence-operator/jvm/</code> directory, as this is shared between both containers.</p>

<markup
lang="bash"

>jcmd 55 GC.heap_dump /coherence-operator/jvm/heap-dump.hprof</markup>

<p>After running the command above, we will have a heap dump file that we can access from the ephemeral <code>Pod</code>.
We have a number of choices about how to get the file out of the Pod and somewhere that we can analyze it.
We could use <code>sftp</code> to ship it somewhere, or some tools to copy it to cloud storage or just simply use <code>kubectl cp</code> to copy it.</p>

<div class="admonition note">
<p class="admonition-inline">Do not exit out of the ephemeral container session until you have copied the heap dump.</p>
</div>
<p>The <code>kubectl cp</code> command is in the form <code>kubectl cp &lt;namespace&gt;/&lt;pod&gt;/&lt;file&gt; &lt;local-file&gt; -c &lt;container&gt;</code>.
So to use <code>kubectl cp</code> we can execute a command like the following:</p>

<markup
lang="bash"

>kubectl cp coherence-test/test-cluster-1:/coherence-operator/jvm/heap-dump.hprof \
    $(pwd)/heap-dump.hprof -c debug</markup>

<p>We will now have a file called <code>heap-dump.hprof</code> in the current directory.
We can now exit out of the ephemeral container.</p>

</div>
</div>
</doc-view>
