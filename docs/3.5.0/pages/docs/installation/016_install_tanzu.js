<doc-view>

<h2 id="_install_on_tanzu">Install On Tanzu</h2>
<div class="section">
<p>If using <a id="" title="" target="_blank" href="https://www.vmware.com/products/app-platform/tanzu">VMWare Tanzu</a> the Coherence Operator can be installed as a package.
Under the covers, Tanzu uses the <a id="" title="" target="_blank" href="https://carvel.dev">Carvel</a> tool set to deploy packages.
The Carvel tools can be used outside Tanzu, so the Coherence Operator repo and package images could also be deployed
using a standalone Carvel <a id="" title="" target="_blank" href="https://carvel.dev/kapp-controller/">kapp-controller</a>.</p>

<p>The Coherence Operator release published two images required to deploy the Operator as a Tanzu package.</p>

<ul class="ulist">
<li>
<p><code>ghcr.io/oracle/coherence-operator-package:3.5.0</code> - the Coherence Operator package</p>

</li>
<li>
<p><code>ghcr.io/oracle/coherence-operator-repo:3.5.0</code> - the Coherence Operator repository</p>

</li>
</ul>

<h3 id="_install_the_coherence_repository">Install the Coherence Repository</h3>
<div class="section">
<p>The first step to deploy the Coherence Operator package in Tanzu is to add the repository.
This can be done using the Tanzu CLI.</p>

<markup
lang="bash"

>tanzu package repository add coherence-repo \
    --url ghcr.io/oracle/coherence-operator-repo:3.5.0 \
    --namespace coherence \
    --create-namespace</markup>

<p>The installed repositories can be listed using the CLI:</p>

<markup
lang="bash"

>tanzu package repository list --namespace coherence</markup>

<p>which should display something like the following</p>

<markup
lang="bash"

>NAME            REPOSITORY                              TAG  STATUS               DETAILS
coherence-repo  ghcr.io/oracle/coherence-operator-repo  1h   Reconcile succeeded</markup>

<p>The available packages in the Coherence repository can also be displayed using the CLI</p>

<markup
lang="bash"

>tanzu package available list --namespace coherence</markup>

<p>which should include the Operator package, <code>coherence-operator.oracle.github.com</code> something like the following</p>

<markup
lang="bash"

>NAME                                  DISPLAY-NAME               SHORT-DESCRIPTION                                             LATEST-VERSION
coherence-operator.oracle.github.com  Oracle Coherence Operator  A Kubernetes operator for managing Oracle Coherence clusters  3.5.0</markup>

</div>

<h3 id="_install_the_coherence_operator_package">Install the Coherence Operator Package</h3>
<div class="section">
<p>Once the Coherence Operator repository has been installed, the <code>coherence-operator.oracle.github.com</code> package can be installed, which will install the Coherence Operator itself.</p>

<markup
lang="bash"

>tanzu package install coherence \
    --package-name coherence-operator.oracle.github.com \
    --version 3.5.0 \
    --namespace coherence</markup>

<p>The Tanzu CLI will display the various steps it is going through to install the package and if all goes well, finally display <code>Added installed package 'coherence'</code>
The packages installed in the <code>coherence</code> namespace can be displayed using the CLI.</p>

<markup
lang="bash"

>tanzu package installed list --namespace coherence</markup>

<p>which should display the Coherence Operator package.</p>

<markup
lang="bash"

>NAME       PACKAGE-NAME                          PACKAGE-VERSION  STATUS
coherence  coherence-operator.oracle.github.com  3.5.0            Reconcile succeeded</markup>

<p>The Operator is now installed and ready to mage Coherence clusters.</p>

</div>
</div>
</doc-view>
