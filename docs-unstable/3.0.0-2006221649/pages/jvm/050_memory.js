<doc-view>

<h2 id="_heap_memory_settings">Heap &amp; Memory Settings</h2>
<div class="section">
<p>The JVM has a number of arguments that set the sizes of different memory regions; the most commonly set is the heap
size but there are a number of others. The <code>Coherence</code> CRD spec has fields that allow some of these to sizes to be
set.</p>

<p>The <code>Coherence</code> CRD also has settings to control the behaviour of the JVM if an out of memory error occurs.
For example, killing the container, creating a heap dump etc.</p>


<h3 id="_max_ram">Max RAM</h3>
<div class="section">
<p>The JVM has an option <code>-XX:MaxRAM=N</code> the maximum amount of memory used by the JVM to <code>n</code>, where <code>n</code> is expressed in
terms of megabytes (for example, <code>100m</code>) or gigabytes (for example <code>2g</code>).</p>

<p>When using resource limited containers it is useful to set the max RAM option to avoid the JVM exceeding the
container limits.</p>

<p>The <code>Coherence</code> CRD allows the max RAM option to be set using the <code>jvm.memory.maxRAM</code> field, for example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  jvm:
    memory:
      maxRAM: 10g</markup>

</div>

<h3 id="_heap_size_as_a_percentage_of_container_memory">Heap Size as a Percentage of Container Memory</h3>
<div class="section">
<p>There are three JVM options that can be used to control the JVM heap as a percentage of the available memory.
These options can be useful when controlling memory as a percentage of container memory in combination
with rescource limits on containers.</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 50%;">
<col style="width: 50%;">
</colgroup>
<thead>
<tr>
<th>JVM Option</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>-XX:InitialRAMPercentage=N</code></td>
<td class="">Sets the initial amount of memory that the JVM will use for the Java heap before applying ergonomics heuristics as a
percentage of the maximum amount determined as described in the -XX:MaxRAM option. The default value is 1.5625 percent.</td>
</tr>
<tr>
<td class="">'-XX:MaxRAMPercentage=N'</td>
<td class="">Sets the maximum amount of memory that the JVM may use for the Java heap before applying ergonomics heuristics as a
percentage of the maximum amount determined as described in the <code>-XX:MaxRAM</code> option.
The default value is 25 percent.</td><td class="">Specifying this option disables automatic use of compressed oops if the combined result of this and other options
influencing the maximum amount of memory is larger than the range of memory addressable by compressed oops.</td>
</tr>
<tr>
<td class="">'-XX:MinRAMPercentage=N'</td>
<td class="">Sets the maximum amount of memory that the JVM may use for the Java heap before applying ergonomics heuristics as a
percentage of the maximum amount determined as described in the -XX:MaxRAM option for small heaps. A small heap is a heap of approximately 125 MB.
The default value is 50 percent.</td>
</tr>
</tbody>
</table>
</div>
<p>Where <code>N</code> is a decimal value between 0 and 100. For example, 12.3456.</p>

<p>When running in a container, and the <code>-XX:+UseContainerSupport</code> is set (which it is by default for the Coherence
container), both the default heap size for containers, the <code>-XX:InitialRAMPercentage</code> option, the <code>-XX:MaxRAMPercentage</code>
option, and the <code>-XX:MaxRAMPercentage</code> option, will be based on the available container memory.</p>

<div class="admonition note">
<p class="admonition-inline">Some JVMs may not support these options.</p>
</div>
<p>The <code>Coherence</code> CRD allows these options to be set with the <code>jvm.memory.initialRAMPercentage</code>, <code>jvm.memory.minRAMPercentage</code>,
and <code>jvm.memory.maxRAMPercentage</code> fields.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  jvm:
    memory:
      initialRAMPercentage: 10
      minRAMPercentage: 5.75
      maxRAMPercentage: 75</markup>

<div class="admonition note">
<p class="admonition-inline">Setting the <code>jvm.memory.heapSize</code> field will cause RAM percentage fields to be ignored.</p>
</div>
<div class="admonition note">
<p class="admonition-inline">Due to CRDs not supporting decimal fields the RAM percentage fields are of type resource.Quantity,
see the Kubernetes <a id="" title="" target="_blank" href="https://godoc.org/k8s.io/apimachinery/pkg/api/resource#Quantity">Quantity</a> API docs for details
of the different number formats allowed.</p>
</div>
</div>

<h3 id="_heap_size">Heap Size</h3>
<div class="section">
<p>To set the JVM heap size set the <code>jvm.memory.heapSize</code> field.
The value of the field can be any value that can be used with the JVM <code>-Xmx</code> and <code>-Xms</code> arguments.
The value of the <code>jvm.memory.heapSize</code> field will be used to set both the <code>-Xms</code> and <code>-Xmx</code> arguments,
so the heap will be a fixed size. For example setting <code>jvm.memory.heapSize</code> to <code>5g</code> will effectively pass
<code>-Xms5g -Xmx5g</code> to the JVM.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  jvm:
    memory:
      heapSize: 10g  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">This example sets the heap size to <code>10g</code>.</li>
</ul>
</div>

<h3 id="_direct_memory_size_nio_memory">Direct Memory Size (NIO Memory)</h3>
<div class="section">
<p>Direct memory size is used to limit on memory that can be reserved for all Direct Byte Buffers.
If a value is set for this option, the sum of all Direct Byte Buffer sizes cannot exceed the limit.
After the limit is reached, a new Direct Byte Buffer can be allocated only when enough old buffers are freed to provide
enough space to allocate the new buffer.</p>

<p>By default, the VM limits the amount of heap memory used for Direct Byte Buffers to approximately 85% of the maximum heap size.</p>

<p>To set a value for the direct memory size use the <code>jvm.memory.directMemorySize</code> field. This wil set the value of the
<code>-XX:MaxDirectMemorySize</code> JVM option.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  jvm:
    memory:
      directMemorySize: 10g  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The direct memory size is set to <code>10g</code> which will pass <code>-XX:MaxDirectMemorySize=10g</code> to the JVM.</li>
</ul>
</div>

<h3 id="_metaspace_size">Metaspace Size</h3>
<div class="section">
<p>Metaspace is memory the VM uses to store class metadata.
Class metadata are the runtime representation of java classes within a JVM process - basically any information the JVM
needs to work with a Java class. That includes, but is not limited to, runtime representation of data from the JVM
class file format.</p>

<p>To set the size of the metaspace use the <code>jvm.memory.metaspaceSize</code> field in the <code>Coherence</code> CRD.
Setting this field sets both the <code>-XX:MetaspaceSize</code> and <code>-XX:MaxMetaspaceSize</code> JVM options to this value giving a
fixed size metaspace.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  jvm:
    memory:
      metaspaceSize: 100m  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">Set the metaspace size to <code>100m</code> which will pass <code>-XX:MetaspaceSize=100m -XX:MaxMetaspaceSize=100m</code>
to the JVM.</li>
</ul>
</div>

<h3 id="_stack_size">Stack Size</h3>
<div class="section">
<p>Thread stacks are memory areas allocated for each Java thread for their internal use.
This is where the thread stores its local execution state.
The current default size for a linux JVM is 1MB.</p>

<p>To set the stack size use the <code>jvm.memory.stackSize</code> field in the <code>Coherence</code> CRD.
Setting this value sets the <code>-Xss</code> JVM option.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  jvm:
    memory:
      stackSize: 500k  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The stack size will be set to <code>500k</code>, passing <code>-Xss500k</code> to the JVM.</li>
</ul>
</div>

<h3 id="_out_of_memory_behaviour">Out Of Memory Behaviour</h3>
<div class="section">
<p>The <code>Coherence</code> CRD allows two optional behaviours to be specified if the JVM throws an out of memory error.</p>

<p>The <code>jvm.memory.onOutOfMemory.heapDump</code> is a bool field that when set to true will pass the
<code>-XX:+HeapDumpOnOutOfMemoryError</code> option to the JVM. The default value of the field when not specified is <code>true</code>,
hence to turn off heap dumps on OOM set the specifically field to be <code>false</code>.</p>

<p>The <code>jvm.memory.onOutOfMemory.exit</code> is a bool field that when set to true will pass the
<code>-XX:+ExitOnOutOfMemoryError</code> option to the JVM. The default value of the field when not specified is <code>true</code>,
hence to turn off killing the JVM on OOM set the specifically field to be <code>false</code>.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  jvm:
    memory:
      onOutOfMemory:
        heapDump: true   <span class="conum" data-value="1" />
        exit: true       <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The JVM will create a heap dump on OOM</li>
<li data-value="2">The JVM will exit on OOM</li>
</ul>
</div>

<h3 id="_native_memory_tracking">Native Memory Tracking</h3>
<div class="section">
<p>The Native Memory Tracking (NMT) is a Java VM feature that tracks internal memory usage for a JVM.
The <code>Coherence</code> CRD allows native memory tracking to be configured using the <code>jvm.memory.nativeMemoryTracking</code> field.
Setting this field sets the <code>-XX:NativeMemoryTracking</code> JVM option. There are three valid values, <code>off</code>, <code>summary</code> or <code>detail</code>.
If not specified the default value used by the operator is <code>summary</code></p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  jvm:
    memory:
      nativeMemoryTracking: detail <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">Native memory tracking is set to <code>detail</code> which will pass the <code>-XX:NativeMemoryTracking=detail</code> option to the JVM.</li>
</ul>
</div>
</div>
</doc-view>
