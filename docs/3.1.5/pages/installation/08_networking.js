<doc-view>

<h2 id="_os_networking_configuration">O/S Networking Configuration</h2>
<div class="section">

<h3 id="_operating_system_library_requirements">Operating System Library Requirements</h3>
<div class="section">
<p>In order for Coherence clusters to form correctly, the <code>conntrack</code> library
must be installed.  Most Kubernetes distributions will do this for you.
If you have issues with clusters not forming, then you should check that
<code>conntrack</code> is installed using this command (or equivalent):</p>

<markup
lang="bash"

>rpm -qa | grep conntrack</markup>

<p>You should see output similar to that shown below.  If you do not, then you
should install <code>conntrack</code> using your operating system tools.</p>

<markup
lang="bash"

>libnetfilter_conntrack-1.0.6-1.el7_3.x86_64
conntrack-tools-1.4.4-4.el7.x86_64</markup>

</div>

<h3 id="_firewall_iptables_requirements">Firewall (iptables) Requirements</h3>
<div class="section">
<p>Some Kubernetes distributions create <code>iptables</code> rules that block some
types of traffic that Coherence requires to form clusters.  If you are
not able to form clusters, then you can check for this issue using the
following command:</p>

<markup
lang="bash"

>iptables -t nat -v  -L POST_public_allow -n</markup>

<p>If you see output similar to the example below:</p>

<markup
lang="bash"

>Chain POST_public_allow (1 references)
pkts bytes target     prot opt in     out     source               destination
164K   11M MASQUERADE  all  --  *      !lo     0.0.0.0/0            0.0.0.0/0
   0     0 MASQUERADE  all  --  *      !lo     0.0.0.0/0            0.0.0.0/0</markup>

<p>For example, if you see any entries in this chain, then you need to remove them.
You can remove the entries using this command:</p>

<markup
lang="bash"

>iptables -t nat -v -D POST_public_allow 1</markup>

<p>Note that you will need to run that command for each line. So in the example
above, you would need to run it twice.</p>

<p>After you are done, you can run the previous command again and verify that
the output is now an empty list.</p>

<p>After making this change, restart your domains and the Coherence cluster
should now form correctly.</p>


<h4 id="_make_iptables_updates_permanent_across_reboots">Make iptables Updates Permanent Across Reboots</h4>
<div class="section">
<p>The recommended way to make <code>iptables</code> updates permanent across reboots is
to create a <code>systemd</code> service that applies the necessary updates during
the startup process.</p>

<p>Here is an example; you may need to adjust this to suit your own
environment:</p>

<ul class="ulist">
<li>
<p>Create a <code>systemd</code> service:</p>

</li>
</ul>
<markup
lang="bash"

>echo 'Set up systemd service to fix iptables nat chain at each reboot (so Coherence will work)...'
mkdir -p /etc/systemd/system/
cat &gt; /etc/systemd/system/fix-iptables.service &lt;&lt; EOF
[Unit]
Description=Fix iptables
After=firewalld.service
After=docker.service

[Service]
ExecStart=/sbin/fix-iptables.sh

[Install]
WantedBy=multi-user.target
EOF</markup>

<ul class="ulist">
<li>
<p>Create the script to update <code>iptables</code>:</p>

</li>
</ul>
<markup
lang="bash"

>cat &gt; /sbin/fix-iptables.sh &lt;&lt; EOF
#!/bin/bash
echo 'Fixing iptables rules for Coherence issue...'
TIMES=$((`iptables -t nat -v -L POST_public_allow -n --line-number | wc -l` - 2))
COUNTER=1
while [ $COUNTER -le $TIMES ]; do
  iptables -t nat -v -D POST_public_allow 1
  ((COUNTER++))
done
EOF</markup>

<ul class="ulist">
<li>
<p>Start the service (or just reboot):</p>

</li>
</ul>
<markup
lang="bash"

>echo 'Start the systemd service to fix iptables nat chain...'
systemctl enable --now fix-iptables</markup>

</div>
</div>
</div>
</doc-view>
