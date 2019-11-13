<doc-view>

<h2 id="_logging_configuration">Logging Configuration</h2>
<div class="section">
<p>When configuring a <code>Pod</code> in Kubernetes there are a number of settings related to networking and DNS and these can also
be configured for roles in a <code>CoherenceCluster</code>.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  network:
    hostname: "foo.com"                  <span class="conum" data-value="1" />
    hostNetwork: false                   <span class="conum" data-value="2" />
    hostAliases:
      - ip: "10.10.10.100"               <span class="conum" data-value="3" />
        hostnames:
          - "a.foo.com"
          - "b.foo.com"
    dnsPolicy: "ClusterFirstWithHostNet" <span class="conum" data-value="4" />
    dnsConfig:                           <span class="conum" data-value="5" />
      nameservers:
        - "dns.foo.com"
        - "dns.bar.com"
      searches:
        - "foo.com"
        - "bar.com"
      options:
        - name: "option-name"
          value: "option-value"</markup>

<ul class="colist">
<li data-value="1">the <code>network.hostname</code> field specifies the hostname of the Pod If not specified, the pod&#8217;s hostname will be set
to a system-defined value as per the Kubernetes defaults.</li>
<li data-value="2">the <code>hostNetwork</code> field setw whether host networking is requested for Pods in a role . If set to true Pods will use
the host&#8217;s network namespace. If this option is set, the ports that will be used must be specified. If set to <code>true</code>
care must be taken not to schedule multiple Pods for a role onto the same Kubernetes node. Default to false.</li>
<li data-value="3">the <code>hostAliases</code> field adds host aliases to the Pods in a roles.
See the Kubernetes documentation on
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/services-networking/add-entries-to-pod-etc-hosts-with-host-aliases/">Adding entries to Pod /etc/hosts with HostAliases</a></li>
<li data-value="4">the <code>dnsPolicy</code> field sets the DNS policy for the pod. Defaults to "ClusterFirst". Valid values are
'ClusterFirstWithHostNet', 'ClusterFirst', 'Default' or 'None'. DNS parameters given in DNSConfig will be merged with
the policy selected with DNSPolicy. To have DNS options set along with hostNetwork, you have to specify DNS policy
explicitly to 'ClusterFirstWithHostNet'.
See the Kubernetes documentation on
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/">DNS for Services and Pods</a></li>
<li data-value="5">the <code>dnsConfig</code> section sets other DNS configuration that will be applied to Pods for a role.
See the Kubernetes documentation on
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/">DNS for Services and Pods</a></li>
</ul>
<p>As with other configuration secions in the CRD <code>spec</code> the <code>network</code> section can be specified under the <code>spec</code> section
if configuring a single implied role or under individual roles if configuring explicit roles or a combination of both
if configuring explicit roles with defaults.</p>


<h3 id="_setting_network_defaults_with_explicit_roles">Setting Network Defaults with Explicit Roles</h3>
<div class="section">
<p>If configuring explicit roles with default values it is important to note how some fields are merged.</p>


<h4 id="_hostaliases">HostAliases</h4>
<div class="section">
<p>When using default values the <code>hostAliases</code> field is merged using the <code>ip</code> field to identify duplicate aliases.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  network:
    hostAliases:                <span class="conum" data-value="1" />
      - ip: "10.10.10.100"
        hostnames:
          - "a.foo.com"
          - "b.foo.com"
      - ip: "10.10.10.200"
        hostnames:
          - "a.bar.com"
          - "b.bar.com"
  roles:
    - role: data                <span class="conum" data-value="2" />
    - role: proxy
      network:
        hostAliases:            <span class="conum" data-value="3" />
          - ip: "10.10.10.100"
            hostnames:
              - "c.foo.com"
          - ip: "10.10.10.300"
            hostnames:
              - "acme.com"</markup>

<ul class="colist">
<li data-value="1">The default <code>hostAliases</code> list contains aliases for the ip addresses <code>10.10.10.100</code> and <code>10.10.10.200</code></li>
<li data-value="2">The <code>data</code> role does not specify any <code>hostAliases</code> so it will use the default aliases for the ip addresses
<code>10.10.10.100</code> and <code>10.10.10.200</code></li>
<li data-value="3">The <code>proxy</code> role specifies an alias for the ip addresses <code>10.10.10.100</code> and <code>10.10.10.300</code> so when the <code>proxy</code>
role&#8217;s alias list is merged with the defaults the alias for <code>10.10.10.200</code> will be inherited from the defaults, the
<code>proxy</code> role&#8217;s own alias for <code>10.10.10.100</code> will override the default alias for the same ip address, and the <code>proxy</code>
role&#8217;s alias for <code>10.10.10.300</code> will also be used. The <code>proxy</code> role&#8217;s effective alias list will be:</li>
</ul>
<markup
lang="yaml"

>hostAliases:
  - ip: "10.10.10.100"
    hostnames:
      - "c.foo.com"
  - ip: "10.10.10.200"
    hostnames:
      - "a.bar.com"
      - "b.bar.com"
  - ip: "10.10.10.300"
    hostnames:
      - "acme.com"</markup>

</div>

<h4 id="_dns_config_nameservers">DNS Config NameServers</h4>
<div class="section">
<p>The <code>dnsConfig.nameservers</code> field is a list of strings so the effective list of <code>nameservers</code> that applies for a role is
any <code>nameservers</code> set at the default level with any <code>nameservers</code> set for the role appended to it.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  network:
    dnsConfig:
      nameservers:
        - "dns.foo.com"    <span class="conum" data-value="1" />
  roles:
    - role: data           <span class="conum" data-value="2" />
    - role proxy
      network:
        dnsConfig:
          nameservers:
          - "dns.bar.com"  <span class="conum" data-value="3" /></markup>

<ul class="colist">
<li data-value="1">The default <code>dnsConfig.nameservers</code> list has a single entry for <code>dns.foo.com</code></li>
<li data-value="2">The <code>data</code> role does not specify a <code>nameservers</code> list so it will inherit just the default <code>dns.foo.com</code></li>
<li data-value="3">The <code>proxy</code> role does specify <code>nameservers</code> list so this will be merged with the defaults giving an effective
list of <code>dns.foo.com</code> and <code>dns.bar.com</code></li>
</ul>
<div class="admonition note">
<p class="admonition-inline">The operator will not attempt to remove duplicate values when merging <code>nameserver</code> lists so if a value appears in
the default list and in a role list then that value will appear twice in the effective list used to create Pods.</p>
</div>
</div>

<h4 id="_dns_config_searches">DNS Config Searches</h4>
<div class="section">
<p>The <code>dnsConfig.searches</code> field is a list of strings so the effective list of <code>searches</code> that applies for a role is
any <code>searches</code> set at the default level with any <code>searches</code> set for the role appended to it.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  network:
    dnsConfig:
      searches:
        - "foo.com"    <span class="conum" data-value="1" />
  roles:
    - role: data       <span class="conum" data-value="2" />
    - role proxy
      network:
        dnsConfig:
          searches:
          - "bar.com"  <span class="conum" data-value="3" /></markup>

<ul class="colist">
<li data-value="1">The default <code>dnsConfig.searches</code> list has a single entry for <code>foo.com</code></li>
<li data-value="2">The <code>data</code> role does not specify a <code>searches</code> list so it will inherit just the default <code>foo.com</code></li>
<li data-value="3">The <code>proxy</code> role does specify <code>searches</code> list so this will be merged with the defaults giving an effective
list of <code>foo.com</code> and <code>bar.com</code></li>
</ul>
<div class="admonition note">
<p class="admonition-inline">The operator will not attempt to remove duplicate values when merging <code>searches</code> lists so if a value appears in
the default list and in a role list then that value will appear twice in the effective list used to create Pods.</p>
</div>
</div>

<h4 id="_dns_config_options">DNS Config Options</h4>
<div class="section">
<p>The <code>network.dnsConfig.options</code> field is a list of name/value pairs. If <code>options</code> are set at both the default and role
level then the lists are merged using the <code>name</code> to identify options.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  network:
    dnsConfig:
      options:                        <span class="conum" data-value="1" />
        - name: "option-one"
          value: "value-one"
        - name: "option-two"
          value: "value-two"
  roles:
    - role: data                      <span class="conum" data-value="2" />
    - role: proxy
      network:
        dnsConfig:
          options:
            - name: "option-one"      <span class="conum" data-value="3" />
              value: "different-one"
            - name: "option-three"
              value: "value-three"</markup>

<ul class="colist">
<li data-value="1">The default <code>options</code> has a single value with <code>name: option-one</code>, <code>value: value-one</code> and
<code>name: option-two</code>, <code>value: value-two</code></li>
<li data-value="2">The <code>data</code> role does not specify any options so it will just inherit the defaults of <code>name: option-one</code>,
<code>value: value-one</code> and <code>name: option-two</code>, <code>value: value-two</code></li>
<li data-value="3">The <code>proxy</code> role specifies two <code>options</code>, one with the name <code>option-one</code> which will override the default option
named <code>option-one</code> and an additonal option named <code>option-three</code> so the effective list applied to the proxy role will be:</li>
</ul>
<markup
lang="yaml"

>options:
    - name: "option-one"
      value: "different-one"
    - name: "option-two"
      value: "value-two"
    - name: "option-three"
      value: "value-three"</markup>

</div>
</div>
</div>
</doc-view>
