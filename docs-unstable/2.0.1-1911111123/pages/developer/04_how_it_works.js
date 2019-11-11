<doc-view>

<h2 id="_how_the_operator_works">How The Operator Works</h2>
<div class="section">
<p>The high level operation of the Coherence Operator can be seen in the diagram below.</p>



<v-card>
<v-card-text class="overflow-y-hidden" >
<img src="./images/operator.png" alt="operator"width="1000" />
</v-card-text>
</v-card>

<p>The entry point to the operator is the`main()` function in the <code>cmd/manager/main.go</code> file. This function performs
the creation and initialisation of the three controllers and the ReST server. It also creates a configuration k8s
<code>secret</code> that is used by Coherence Pods. The Coherence Operator works in a single namespace, that is it manages CRDs
and hence Coherence clusters only in the same namespace that it is installed into.</p>


<h3 id="_controllers">Controllers</h3>
<div class="section">
<p>In the Operator SDK framework a controller is responsible for managing a specific CRD. A single controller could,
in theory, manage multiple CRDs but it is clearer and simpler to keep them separate. The Coherence Operator has three
controllers, two are part of the operator source code and one is provided by the Operator SDK framework.</p>

<p>All controllers have a <code>Reconcile</code> function that is triggered by events from Kubernetes for resources that the
controller is listening to.</p>

</div>

<h3 id="_coherencecluster_controller">CoherenceCluster Controller</h3>
<div class="section">
<p>The CoherenceCluster controller manages instances of the CoherenceCluster CRD. The source for this controller is
in the <code>pkg/controller/coherencecluster/coherencecluster_controller.go</code> file.
The CoherenceCluster controller listens for events related to CoherenceCluster CRDs created or modified in the
namespace that the operator is running in. It also listens to events for any CoherenceRole CRD that it owns. When
a CoherenceCluster resource is created or modified a CoherenceRole is created (or modified or deleted) for each role
in the CoherenceCluster spec. Each time a k8s event is raised for a CoherenceCluster or CoherenceRole resource the
<code>Reconcile</code> method on the CoherenceCluster controller is called.</p>

<ul class="ulist">
<li>
<p><strong>Create</strong> -
When a CoherenceCluster is created the controller will work out how many roles are present in the spec. For each role
that has a <code>Replica</code> count greater than zero a CoherenceRole is created in k8s. When a CoherenceRole is created it is
associated to the parent CoherenceCluster so that k8s can track ownership of related resources (this is used for
cascade delete - see below).</p>

</li>
<li>
<p><strong>Update</strong> -
When a CoherenceCluster is updated the controller will work out what the roles in the updated spec should be.
It then compares these roles to the currently deployed CoherenceRoles for that cluster. It then creates, updates or
deletes CoherenceRoles as required.</p>

</li>
<li>
<p><strong>Delete</strong> -
When a CoherenceCluster is deleted the controller does not currently need to do anything. This is because k8s has
cascade delete functionality that allows related resources to be deleted together (a little like cascade delete in
a database). When a CoherenceCluster is deleted then any related CoherenceRoles will be deleted and also any resources
that have those CoherenceRoles as owners (i.e. the corresponding CoherenceInternal resources)</p>

</li>
</ul>
</div>

<h3 id="_coherencerole_controller">CoherenceRole Controller</h3>
<div class="section">
<p>The CoherenceRole controller manages instances of the CoherenceRole CRD. The source for this controller is
in the <code>pkg/controller/coherencerole/coherencerole_controller.go</code> file.</p>

<p>The CoherenceRole controller listens for events related to CoherenceRole CRDs created or modified in the
namespace that the operator is running in. It also listens to events for any StatefulSet resources that were created
by the corresponding Helm install for the role.
When a CoherenceRole resource is created or modified a corresponding CoherenceInternal resource is created
(or modified or deleted) from the role&#8217;s spec. Creation of a CoherenceInternal resource will trigger a Helm install
of the Coherence Helm chart by the Helm Controller.
Each time a k8s event is raised for a CoherenceRole or for a StatefulSet resource related to the role the
<code>Reconcile</code> method on the CoherenceRole controller is called.</p>

<p>The StatefulSet resource is listened to as a way to keep track of the state fo the role, i.e how many replicas are actually
running and ready compared to the desired state. The StatefulSet is also used to obtain references to the Pods that make up
the role when performing a StatusHA check prior to scaling.</p>

<ul class="ulist">
<li>
<p><strong>Create</strong> -
When a CoherenceRole is created a corresponding CoherenceInternal resource will be created in k8s.</p>

</li>
<li>
<p><strong>Update</strong> -
When a CoherenceRole is updated one of three actions can take place.</p>
<ul class="ulist">
<li>
<p>Scale Up - If the update increases the role&#8217;s replica count then the role is being scaled up. The role&#8217;s spec is
first checked to determine whether anything else has changed, if it has a rolling upgrade is performed first to bring
the existing members up to the desired spec. After any possible the upgrade then the role&#8217;s member count is scaled up.</p>

</li>
<li>
<p>Scale Down - If the update decreases the role&#8217;s replica count then the role is being scaled down. The member count
of the role is scaled down and then the role&#8217;s spec is checked to determine whether anything else has changed, if it has
a rolling upgrade is performed to bring the remaining members up to the desired spec.</p>

</li>
<li>
<p>Update Only - If the changes to the role&#8217;s spec do not include a change to the replica count then a rolling upgrade
is performed of the existing cluster members.</p>

</li>
</ul>
</li>
<li>
<p><strong>Rolling Upgrade</strong> -
A rolling upgrade is actually performed out of the box by the StatefulSet associated to the role. To upgrade the
members of a role the CoherenceRole controller only has to update the CoherenceInternal spec. This will cause the Helm
controller to update the associated Helm install whivh in turn causes the StatefulSet to perform a rolling update of
the associated Pods.</p>

</li>
<li>
<p><strong>Scaling</strong> -
The CoherenceOperator supports safe scaling of the members of a role. This means that a scaling operation will not take
place unless the members of the role are Status HA. Safe scaling means that the number of replicas is scaled one at a time
untile the desired size is reached with a Status HA check being performed before each member is added or removed.
The exact action is controlled by a customer defined scaling policy that is part of the role&#8217;s spec.
There are three policy types:</p>
<ul class="ulist">
<li>
<p>SafeScaling - the safe scaling policy means that regardless of whether a role is being scaled up or down the size
is always scaled one at a time with a Status HA check before each member is added or removed.</p>

</li>
<li>
<p>ParallelScaling - with parallel scaling no Status HA check is performed, a role is scaled to the desired size by
adding or removing the required number of members at the same time. For a storage enabled role with this policy scaling
down could result in data loss. Ths policy is intended for storage disabled roles where it allows for fatser start and
scaling times.</p>

</li>
<li>
<p>ParallelUpSafeDownScaling - this policy is the default scaling policy. It means that when scaling up the required number
of members is added all at once but when scaling down members are removed one at a time with a Status HA check before each
removal. This policy allows clusters to start and scale up fatser whilst protecting from data loss when scaling down.</p>

</li>
</ul>
</li>
<li>
<p><strong>Delete</strong> -
As with a CoherenceCluster, when a CoherenceRole is deleted its corresponding CoherenceInternal resource is also deleted
by a cascading delete in k8s. The CoherenceRole controller does not need to take any action on deletion.</p>

</li>
</ul>
</div>

<h3 id="_helm_controller">Helm Controller</h3>
<div class="section">
<p>The final controller in the Coherence Operator is the Helm controller. This controller is actually part of the Operator SDK
and the source is not in the Coherence Operator&#8217;s source code tree. The Helm controller is configured to watch for a
particular CRD and performs Helm install, delete and upgrades as resources based on that CRD are created, deleted or updated.</p>

<p>In the case of the Coherence Operator the Helm controller is watching for instances of the CoherenceInternal CRD that are
created, updated or deleted by the CoherenceRole controller. When this occurs the Helm controller uses the spec of the
CoherenceInternal resource as the values file to install or upgrade the Coherence Helm chart.</p>

<p>The Coherence Helm chart used by the operator is actually embedded in the Coherence Operator Docker image so there is no
requirement for the customer to have access to a chart repository.</p>

<p>The Helm operator also uses an embedded helm and tiller so there is no requirement for the customer to install Helm in
their k8s cluster. A customer can have Helm installed but it will never be used by the operator so there is no version
conflict. If a customer were to perform a <code>helm ls</code> operation in their cluster they would not see the installs controlled
by the Coherence Operator.</p>

</div>
</div>
</doc-view>
