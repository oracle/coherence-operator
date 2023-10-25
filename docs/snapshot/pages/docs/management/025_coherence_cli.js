<doc-view>

<h2 id="_the_coherence_cli">The Coherence CLI</h2>
<div class="section">
<p>If Coherence Management over REST is enabled, it is possible to use the
<a id="" title="" target="_blank" href="https://github.com/oracle/coherence-cli">Coherence CLI</a>
to access management information. The Operator enables Coherence Management over REST by default, so unless it
has specifically been disabled, the CLI can be used.</p>

<p>See the <a id="" title="" target="_blank" href="https://oracle.github.io/coherence-cli/docs/latest">Coherence CLI Documentation</a>
for more information on how to use the CLI.</p>

<p>The Coherence CLI is automatically added to Coherence Pods by the Operator, so it is available as an executable
that can be run using <code>kubectl exec</code>.
At start-up of a Coherence container a default Coherence CLI configuration is created so that the CLI
knows about the local cluster member.</p>


<h3 id="_using_the_cli_in_pods">Using the CLI in Pods</h3>
<div class="section">
<p>The Operator installs the CLI at the location <code>/coherence-operator/utils/cohctl</code>.
Most official Coherence images are distroless images so they do not have a shell that can be used to create a session and execute commands. Each <code>cohctl</code> command will need to be executed as a separate <code>kubectl exec</code> command.</p>

<p>Once a Pod is running is it simple to use the CLI.
For example, the yaml below will create a simple three member cluster.</p>

<markup

title="minimal.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
    replicas: 3</markup>

<p>The cluster name is <code>storage</code> and there will be three Pods created, <code>storage-0</code>, <code>storage-1</code> and <code>storage-2</code>.</p>

<p>To list the services running in the <code>storage-0</code> Pod the following command can be run:</p>

<markup
lang="bash"

>kubectl exec storage-0 -c coherence -- /coherence-operator/utils/cohctl get services</markup>

<p>The <code>-c coherence</code> option tells <code>kubectl</code> to exec the command in the <code>coherence</code> container.
By default this is the only container that will be running in the Pod, so the option could be omitted.
If the option is omitted, <code>kubectl</code> will display a warning to say it assumes you mean the <code>coherence</code> container.</p>

<p>Everything after the <code>--</code> is the command to run in the Pod. In this case we execute:</p>

<markup
lang="bash"

>/coherence-operator/utils/cohctl get services</markup>

<p>which runs the Coherence CLI binary at <code>/coherence-operator/utils/cohctl</code> with the command <code>get services</code>.</p>

<p>The output displayed by the command will look something like this:</p>

<markup
lang="bash"

>Using cluster connection 'default' from current context.

SERVICE NAME            TYPE              MEMBERS  STATUS HA  STORAGE  PARTITIONS
"$GRPC:GrpcProxy"       Proxy                   3  n/a             -1          -1
"$SYS:Concurrent"       DistributedCache        3  NODE-SAFE        3         257
"$SYS:ConcurrentProxy"  Proxy                   3  n/a             -1          -1
"$SYS:Config"           DistributedCache        3  NODE-SAFE        3         257
"$SYS:HealthHttpProxy"  Proxy                   3  n/a             -1          -1
"$SYS:SystemProxy"      Proxy                   3  n/a             -1          -1
ManagementHttpProxy     Proxy                   3  n/a             -1          -1
MetricsHttpProxy        Proxy                   3  n/a             -1          -1
PartitionedCache        DistributedCache        3  NODE-SAFE        3         257
PartitionedTopic        PagedTopic              3  NODE-SAFE        3         257
Proxy                   Proxy                   3  n/a             -1          -1</markup>

<p>The exact output will vary depending on the version of Coherence and the configurations being used.</p>

<p>More CLI commands can be run by changing the CLI commands specified after <code>/coherence-operator/utils/cohctl</code>.</p>

<p>For example, to list all the members of the cluster:</p>

<markup
lang="bash"

>kubectl exec storage-0 -c coherence -- /coherence-operator/utils/cohctl get members</markup>

</div>

<h3 id="_disabling_cli_access">Disabling CLI Access</h3>
<div class="section">
<p>There may be certain circumstances in which you wish to disable the use of the CLI in your cluster.
To do this, add the <code>CLI_DISABLED</code> env variable to you config and set to <code>true</code>.</p>

<markup

title="minimal.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
    replicas: 3
    env:
     - name: "CLI_DISABLED"
       value: "true"</markup>

<p>If you try to run the CLI you will get the following message:</p>

<markup


>cohctl has been disabled from running in the Coherence Operator</markup>

</div>
</div>
</doc-view>
