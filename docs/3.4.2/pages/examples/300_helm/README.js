<doc-view>

<h2 id="_manage_coherence_resources_using_helm">Manage Coherence Resources using Helm</h2>
<div class="section">
<p>Occasionally there is a requirement to manage Coherence resources using Helm instead of Kubernetes tools such as <code>kubectl</code>. There is no Helm chart for a Coherence resource as it is a single resource and any Helm chart and values file would need to replicate the entire Coherence CRD if it was to be of generic enough use for everyone. For this reason, anyone wanting to manage Coherence resource using Helm will need to create their own chart, which can then be specific to their needs.</p>

<p>This example shows some ways that Helm can be used to manage Coherence resources.</p>

<div class="admonition tip">
<p class="admonition-textlabel">Tip</p>
<p ><p><img src="./images/GitHub-Mark-32px.png" alt="GitHub Mark 32px" />
 The complete source code for this example is in the <a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/tree/main/examples/300_helm">Coherence Operator GitHub</a> repository.</p>
</p>
</div>

<h3 id="_a_simple_generic_helm_chart">A Simple Generic Helm Chart</h3>
<div class="section">
<p>This example contains the most basic Helm chart possible to support managing a Coherence resource locate in the <code>chart/</code> directory. The chart is actually completely generic and would support any configuration of Coherence resource.</p>

<p>The values file contains a single value <code>spec</code>, which will contain the entire spec of the Coherence resource.</p>

<markup
lang="yaml"
title="chart/values.yaml"
>spec:</markup>

<p>There is a single template file, as we only create a single Coherence resource.</p>

<span v-pre><div class="markup-container"><div class="block-title"><span>test-cluster.yaml</span></div><div data-lang="yaml" class="markup"><pre><code class="yaml hljs makefile">apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: {{ .Release.Name }}
  namespace: {{ .Release.Namespace }}
  labels:
{{- include "coherence-labels" . | indent 4 }}
spec:
{{- if .Values.spec }}
{{ toYaml .Values.spec | indent 2 }}
{{- end }}</code></pre><div class="markup__copy"><i aria-hidden="true" class="material-icons icon">content_copy</i><span class="markup__copy__message">Copied</span></div></div></div></span>

<p>The first part of the template is fairly standard for a Helm chart, we configure the resource name, namespace and add some labels.</p>

<p>The generic nature of the chart comes from the fact that the template then just takes the whole <code>spec</code> value from the values file, and if it is not <code>null</code> converts it to yaml under the <code>spec:</code> section of the template. This means that any yaml that is valid in a Coherence CRD <code>spec</code> section can be used in a values file (or with <code>--set</code> arguments) when installing the chart.</p>

</div>

<h3 id="_installing_the_chart">Installing the Chart</h3>
<div class="section">
<p>Installing the example Helm chart is as simple as any other chart. One difference here being that the chart is not installed into a chart repository so has to be installed from the <code>char/</code> directory. If you wanted to you could</p>

<div class="admonition note">
<p class="admonition-inline">The following commands are all run from the <code>examples/helm</code> directory so that the chart location is specified as <code>./chart</code>. You can run the commands from anywhere, but you would need to specify the full path to the example chart directory.</p>
</div>

<h4 id="_a_simple_dry_run">A Simple Dry Run</h4>
<div class="section">
<p>To start with we will do a simple dry-run install that will display the yaml Helm would have created if the install command had been real.</p>

<markup
lang="bash"

>helm  install test ./chart --dry-run</markup>

<p>The above command should result in the following output</p>

<markup


>NAME: test
LAST DEPLOYED: Sat Aug 28 16:30:53 2021
NAMESPACE: default
STATUS: pending-install
REVISION: 1
TEST SUITE: None
HOOKS:
MANIFEST:
---
# Source: coherence-example/templates/coherence.yaml
apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
  namespace: default
  labels:
    app.kubernetes.io/managed-by: Helm
    app.kubernetes.io/name: test
    app.kubernetes.io/version: "1.0.0"
spec:</markup>

<p>We can see at the bottom of the output the simple Coherence resource that would have been created by helm.
This is a valid Coherence resource because every field in the spec section is optional. If the install had been real this would have resulted in a Coherence cluster named "test" with three storage enabled cluster members, as the default replica count is three.</p>

</div>

<h4 id="_setting_values">Setting Values</h4>
<div class="section">
<p>But how do we set other values in the Coherence resouce. That is simple because Helm does not validate what we enter as values we can either create a values file with anything we like under the <code>spec</code> secion or we can specify values using the <code>--set</code> Helm argument.</p>

<p>For example, if we wanted to set the replica count to six in a Coherence resource we would need to set the <code>spec.replicas</code> field to <code>6</code>, and we do exactly the same in the helm chart.</p>

<p>We could create a values file like this:</p>

<markup

title="test-values.yaml"
>spec:
  replicas: 6</markup>

<p>Which we can install with</p>

<markup
lang="bash"

>helm  install test ./chart -f test-values.yaml</markup>

<p>Which would produce a Coherence resource like this:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
  namespace: default
  labels:
    app.kubernetes.io/managed-by: Helm
    app.kubernetes.io/name: test
    app.kubernetes.io/version: "1.0.0"
spec:
  replicas: 6</markup>

<p>We could have done the same thing using <code>--set</code>, for example:</p>

<markup
lang="bash"

>helm  install test ./chart -f test-values.yaml --set spec.replicas=6</markup>

<p>We can even set more deeply nested values, for example the Coherence log level is set in the <code>spec.coherence.logLevel</code> field of the Coherence CRD so we can use the same value in the Helm install command or values file:</p>

<markup
lang="bash"

>helm  install test ./chart -f test-values.yaml --set spec.coherence.logLevel=9</markup>

<p>Which would produce the following Coherence resource:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
  namespace: default
  labels:
    app.kubernetes.io/managed-by: Helm
    app.kubernetes.io/name: test
    app.kubernetes.io/version: "1.0.0"
spec:
  coherence:
    logLevel: 9</markup>

<p>Just like any Helm chart, whether you use <code>--set</code> arguments or use a values file depends on how complex the resulting yaml will be. Some fields of the Coherence CRD spec would be impractical to try to configure on the command line with <code>--set</code> and would be much simpler in the values file.</p>

</div>
</div>

<h3 id="_helm_wait_waiting_for_the_install_to_complete">Helm Wait - Waiting for the Install to Complete</h3>
<div class="section">
<p>The Helm <code>install</code> command (and <code>update</code> command) have a <code>--wait</code> argument that tells Helm to wait until the installed resources are ready. This can be very useful if you want to ensure that everything is created and running correctly after and install or upgrade. If you read the help test for the <code>--wait</code> argument you will see the following:</p>


<p>The limitation should be obvious, Helm can only wait for a sub-set of al the possible resources that you can create from a Helm chart. It has no idea how to wait for a <code>Coherence</code> resource to be ready. To work around this limitation we can use a <a id="" title="" target="_blank" href="https://helm.sh/docs/topics/charts_hooks/">Helm chart hook</a>, mre specifically a post-install and post-upgrade hook.</p>

<p>A hook is typically a k8s Job that Helm will execute, you create the Job spec as part of the Helm chart templates.</p>


<h4 id="_the_coherence_operator_utils_runner">The Coherence Operator Utils Runner</h4>
<div class="section">
<p>The Coherence Operator has two images, the Operator itself and a second image containing an executable named <code>runner</code> which the Operator uses to run Coherence servers in the Pods it is managing. One of the other commands that the runner can execute is a <code>status</code> command, which queries the Operator for the current status of a Coherence resource. If you pull the image and execute it you can see the help text for the runner CLI.</p>

<p>The following commands will pull the Operator utils image and run it to display the help fot eh status command:</p>

<markup
lang="bash"

>docker pull ghcr.io/oracle/coherence-operator:3.4.2
docker run -it --rm ghcr.io/oracle/coherence-operator:3.4.2 status -h</markup>

<p>By creating a K8s Job that runs the status command we can query the Operator for the status of the Coherence resource we installed from the Helm chart. Of course, we could have written something similar that used kubectl in the Job or similar to query k8s for the state of the Coherence resource, but this becomes more complex in RBAC enabled cluster. Querying the simple REST endpoint of the Coherence Operator does not require RBAC rules for the Job to execute.</p>

<p>To run a simple status check we are only interested in the following parameters for the status command:</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 50%;">
<col style="width: 50%;">
</colgroup>
<thead>
<tr>
<th>Argument</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>--operator-url</code></td>
<td class="">The Coherence Operator URL, typically the operator&#8217;s REST service (default "http://coherence-operator-rest.coherence.svc.local:8000"</td>
</tr>
<tr>
<td class=""><code>--namespace</code></td>
<td class="">The namespace the Coherence resource is deployed into. This will be the namespace our Helm chart was installed into.</td>
</tr>
<tr>
<td class=""><code>--name</code></td>
<td class="">The name of the Coherence resource. This will be the name from the Helm chart install</td>
</tr>
<tr>
<td class=""><code>--timeout</code></td>
<td class="">The maximum amount of time to wait for the Coherence resource to reach the required condition (default 5m0s)</td>
</tr>
<tr>
<td class=""><code>--interval</code></td>
<td class="">The status check re-try interval (default 10s)</td>
</tr>
</tbody>
</table>
</div>
<p>First we can add a few additional default values to our Helm chart values file that will be sensible defaults to pass to the hook Job.</p>

<markup
lang="yaml"
title="chart/values.yaml"
>spec:

operator:
  namespace: coherence
  service: coherence-operator-rest
  port: 8000
  image: ghcr.io/oracle/coherence-operator-utils:3.4.2
  condition: Ready
  timeout: 5m
  interval: 10s</markup>

<p>We have added an <code>operator</code> section to isolate the values for the hook from the <code>spec</code> values used in our Coherence resource.</p>

<p>We can now create the hook template in our Helm chart using the new values in the values file.</p>

<span v-pre><div class="markup-container"><div class="block-title"><span>chart/templates/hook.yaml</span></div><div data-lang="yaml" class="markup"><pre><code class="yaml hljs makefile">apiVersion: batch/v1
kind: Job
metadata:
  name: "{{ .Release.Name }}-helm-hook"
  namespace: {{ .Release.Namespace }}
  annotations:
    "helm.sh/hook": post-install,post-upgrade
    "helm.sh/hook-delete-policy": hook-succeeded
spec:
  template:
    metadata:
      name: "{{ .Release.Name }}-helm-hook"
    spec:
      restartPolicy: Never
      containers:
      - name: post-install-job
        image: {{ .Values.operator.image }}
        command:
          - "/files/runner"
          - "status"
          - "--namespace"
          -  {{ .Release.Namespace | quote }}
          - "--name"
          - {{ .Release.Name | quote }}
          - "--operator-url"
          - "http://{{ .Values.operator.service | default "coherence-operator-rest" }}.{{ .Values.operator.namespace | default "coherence" }}.svc:{{ .Values.operator.port | default 8000 }}"
          - "--condition"
          - {{ .Values.operator.condition | default "Ready" | quote }}
          - "--timeout"
          - {{ .Values.operator.timeout | default "5m" | quote }}
          - "--interval"
          - {{ .Values.operator.interval | default "10s" | quote }}</code></pre><div class="markup__copy"><i aria-hidden="true" class="material-icons icon">content_copy</i><span class="markup__copy__message">Copied</span></div></div></div></span>

<p>The annotations section is what tells Helm that this is a hook resource:</p>

<markup
lang="yaml"

>  annotations:
    "helm.sh/hook": post-install,post-upgrade
    "helm.sh/hook-delete-policy": hook-succeeded</markup>

<p>We define the hook as a <code>post-install</code> and <code>post-update</code> hook, so that it runs on both <code>install</code> and <code>update</code> of the Coherence resource.
The hook job will also be deleted once it has successfully run. It will not be deleted if it fails, so we can look at the output of the failure in the Jon Pod logs.</p>

</div>

<h4 id="_installing_with_the_hook">Installing with the Hook</h4>
<div class="section">
<p>If we repeat the Helm install command to install a Coherence resource with the hook in the chart we should see Helm wait and not complete until the Coherence resource (and by inference the StatefulSet and Pods) are all ready.</p>

<markup
lang="bash"

>helm  install test ./chart</markup>

<p>If we were installing a large Coherence cluster, or doing a Helm upgrade, which results in a rolling upgrade of the Coherence cluster, this could take a lot longer than the default timeout we used of 5 minutes. We can alter the timeout and re-try interval using <code>--set</code> arguments.</p>

<markup
lang="bash"

>helm  install test ./chart --set operator.timeout=20m --set operator.interval=1m</markup>

<p>In the above command the timeout is now 20 minutes and the status check will re-try every one minute.</p>

</div>

<h4 id="_skipping_hooks">Skipping Hooks</h4>
<div class="section">
<p>Sometime we might want to install the chart and not wait for everything to be ready. We can use the Helm <code>--no-hooks</code> argument to skip hook execution.</p>

<markup
lang="bash"

>helm  install test ./chart --no-hooks</markup>

<p>Now the Helm install command will return as soon as the Coherence resource has been created.</p>

</div>

<h4 id="_other_helm_hooks">Other Helm Hooks</h4>
<div class="section">
<p>We saw above how a custom post-install and post-update hook could be used to work aroud the restrictions of Helm&#8217;s <code>--wait</code> argument. Of course there are other hooks available in Helm that the method above could be used in. For example, say I had a front end application to be deployed using a Helm chart, but I did not want Helm to start the deployment until the Coherence back-end was ready, I could use the same method above in a pre-install hook.</p>

</div>
</div>
</div>
</doc-view>
