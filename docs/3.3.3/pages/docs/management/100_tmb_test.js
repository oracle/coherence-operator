<doc-view>

<h2 id="_coherence_network_testing">Coherence Network Testing</h2>
<div class="section">
<p>Coherence provides utilities that can be used to test network performance, which obviously has a big impact on
a distributed system such as Coherence. The documentation for these utilities can be found in the official
<a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/standalone/coherence/14.1.1.2206/administer/performing-network-performance-test.html#GUID-7267AB06-6353-416E-B9FD-A75F7FBFE523">Coherence Documentation</a>.</p>

<p>Whilst generally these tests would be run on server hardware, with more and more Coherence deployments moving into the
cloud and into Kubernetes these tests can also be performed in <code>Pods</code> to measure inter-Pod network performance.
This test can be used to see the impact of running <code>Pods</code> across different zones, or on different types of Kubernetes
networks, with different <code>Pod</code> resource settings, etc.</p>

</div>

<h2 id="_run_the_message_bus_test_in_pods">Run the Message Bus Test in Pods</h2>
<div class="section">
<p>The message bus test can easily be run using <code>Pods</code> in Kubernetes.
Using the example from the Coherence documentation there will need to be two <code>Pods</code>, a listener and a sender.
This example will create a <code>Service</code> for the listener so that the sender <code>Pod</code> can use the <code>Service</code> name
to resolve the listener <code>Pod</code> address.</p>


<h3 id="_run_the_listener_pod">Run the Listener Pod</h3>
<div class="section">
<p>Create a <code>yaml</code> file that will create the <code>Service</code> and <code>Pod</code> for the listener:</p>

<markup
lang="yaml"
title="message-bus-listener.yaml"
>apiVersion: v1
kind: Service
metadata:
  name: message-bus-listener
spec:
  selector:
    app: message-bus-listener
  ports:
  - protocol: TCP
    port: 8000
    targetPort: mbus
---
apiVersion: v1
kind: Pod
metadata:
  name: message-bus-listener
  labels:
    app: message-bus-listener
spec:
  restartPolicy: Never
  containers:
    - name: coherence
      image: ghcr.io/oracle/coherence-ce:22.06  <span class="conum" data-value="1" />
      ports:
        - name: mbus
          containerPort: 8000
          protocol: TCP
      command:
        - java                                                   <span class="conum" data-value="2" />
        - -cp
        - /u01/oracle/oracle_home/coherence/lib/coherence.jar
        - com.oracle.common.net.exabus.util.MessageBusTest
        - -bind
        - tmb://0.0.0.0:8000</markup>

<ul class="colist">
<li data-value="1">This example uses a Coherence CE image, but any image with <code>coherence.jar</code> in it could be used.</li>
<li data-value="2">The command line that the container will execute is exactly the same as that for the listener process in the
<a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/standalone/coherence/14.1.1.2206/administer/performing-network-performance-test.html#GUID-7267AB06-6353-416E-B9FD-A75F7FBFE523">Coherence Documentation</a>.</li>
</ul>
<p>Start the listener <code>Pod</code>:</p>

<markup
lang="bash"

>kubectl create -f message-bus-listener.yaml</markup>

<p>Retrieving the logs for the listener <code>Pod</code> the messages should show that the <code>Pod</code> has started:</p>

<markup
lang="bash"

>kubectl logs pod/message-bus-listener
OPEN event for tmb://message-bus-listener:8000</markup>

</div>

<h3 id="_run_the_sender_pod">Run the Sender Pod</h3>
<div class="section">
<markup
lang="yaml"
title="message-bus-sender.yaml"
>apiVersion: v1
kind: Pod
metadata:
  name: message-bus-sender
  labels:
    app: message-bus-sender
spec:
  restartPolicy: Never
  containers:
    - name: coherence
      image: ghcr.io/oracle/coherence-ce:22.06
      command:
        - java                         <span class="conum" data-value="1" />
        - -cp
        - /u01/oracle/oracle_home/coherence/lib/coherence.jar
        - com.oracle.common.net.exabus.util.MessageBusTest
        - -bind
        - tmb://0.0.0.0:8000
        - -peer
        - tmb://message-bus-listener:8000  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">Again, the command line is the same as that for the sender process in the
<a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/standalone/coherence/14.1.1.2206/administer/performing-network-performance-test.html#GUID-7267AB06-6353-416E-B9FD-A75F7FBFE523">Coherence Documentation</a>.</li>
<li data-value="2">The <code>peer</code> address uses the <code>Service</code> name <code>message-bus-listener</code> from the sender <code>yaml</code>.</li>
</ul>
<p>Start the sender <code>Pod</code>:</p>

<markup
lang="bash"

>kubectl create -f message-bus-sender.yaml</markup>

<p>Retrieving the logs for the sender <code>Pod</code> the messages should show that the <code>Pod</code> has started and show the test results:</p>

<markup
lang="bash"

>kubectl logs pod/message-bus-sender
OPEN event for tmb://message-bus-sender:8000
CONNECT event for tmb://message-bus-listener:8000 on tmb://message-bus-sender:8000
now:  throughput(out 34805msg/s 1.14gb/s, in 348msg/s 11.3mb/s), latency(response(avg 25.31ms, effective 110.03ms, min 374.70us, max 158.10ms), receipt 25.47ms), backlog(out 77% 83/s 308KB, in 0% 0/s 0B), connections 1, errors 0
now:  throughput(out 34805msg/s 1.14gb/s, in 348msg/s 11.3mb/s), latency(response(avg 25.31ms, effective 110.03ms, min 374.70us, max 158.10ms), receipt 25.47ms), backlog(out 77% 83/s 308KB, in 0% 0/s 0B), connections 1, errors 0</markup>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>Don&#8217;t forget to stop the <code>Pods</code> after obtaining the results:</p>

<markup
lang="bash"

>kubectl delete -f message-bus-sender.yaml
kubectl delete -f message-bus-listener.yaml</markup>
</p>
</div>
</div>

<h3 id="_run_pods_on_specific_nodes">Run Pods on Specific Nodes</h3>
<div class="section">
<p>In the example above the <code>Pods</code> will be scheduled wherever Kubernetes decides to put them. This could have a big impact
on the test result for different test runs. For example in a Kubernetes cluster that spans zones and data centres, if
the two <code>Pods</code> get scheduled in different data centres this will have worse results than if the two <code>Pods</code> get scheduled
onto the same node.</p>

<p>To get consistent results add node selectors, taints, tolerations etc, as covered in the Kubernetes
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/">assign Pods to Nodes</a> documentation.</p>

</div>
</div>
</doc-view>
