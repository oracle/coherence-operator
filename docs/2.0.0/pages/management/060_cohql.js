<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Accessing CohQL</dt>
<dd slot="desc"><p>You can use Coherence Query Language (CohQL) to interact with Coherence caches.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_accessing_the_cohql_client">Accessing the CohQL client</h2>
<div class="section">
<p>CohQL is a light-weight syntax (in the tradition of SQL) that is used to
perform cache operations on a Coherence cluster. The language can be used
either programmatically or from a command-line tool.</p>

<p>The example shows how to access the Coherence CohQL client in a running cluster.</p>

<p>See the <a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/fusion-middleware/coherence/12.2.1.4/develop-applications/using-coherence-query-language.html">Coherence CohQL documentation</a> for more information.</p>


<h3 id="_1_install_a_coherence_cluster">1. Install a Coherence Cluster</h3>
<div class="section">
<p>Deploy a simple <code>CoherenceCluster</code> resource with a single role like this:</p>

<markup
lang="yaml"
title="example-cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: example-cluster
spec:
  role: storage
  replicas: 3</markup>

<div class="admonition note">
<p class="admonition-inline">Add an <code>imagePullSecrets</code> entry if required to pull images from a private repository.</p>
</div>
<markup
lang="bash"

>kubectl create -n &lt;namespace&gt; -f  example-cluster.yaml

coherencecluster.coherence.oracle.com/example-cluster created

kubectl -n &lt;namespace&gt; get pod -l coherenceCluster=example-cluster

NAME                        READY   STATUS    RESTARTS   AGE
example-cluster-storage-0   1/1     Running   0          59s
example-cluster-storage-1   1/1     Running   0          59s
example-cluster-storage-2   1/1     Running   0          59s</markup>

</div>

<h3 id="_2_connect_to_cohql_client_to_add_data">2. Connect to CohQL client to add data</h3>
<div class="section">
<markup
lang="bash"

>kubectl exec -it -n &lt;namespace&gt; example-cluster-storage-0 bash /scripts/startCoherence.sh queryplus</markup>

<p>Run the following <code>CohQL</code> commands to view and insert data into the cluster.</p>

<markup
lang="sql"

>CohQL&gt; select count() from 'test';

Results
0

CohQL&gt; insert into 'test' key('key-1') value('value-1');

CohQL&gt; select key(), value() from 'test';
Results
["key-1", "value-1"]</markup>

<p>You can issue the <code>help</code> command to get details help information in each command or
<code>commands</code> command to get a brief view of all commands.</p>

<p>Issue the command <code>bye</code> to exit CohQL.</p>

<markup
lang="sql"

>CohQL&gt; commands

java com.tangosol.coherence.dslquery.QueryPlus [-t] [-c] [-s] [-e] [-l &lt;cmd&gt;]*
    [-f &lt;file&gt;]* [-g &lt;garFile&gt;] [-a &lt;appName&gt;] [-dp &lt;parition-list&gt;] [-timeout &lt;value&gt;]

Command Line Arguments:
-a               the application name to use. Used in combination with the -g
                 argument.
-c               exit when command line processing is finished
-e               or -extend
                 extended language mode.  Allows object literals in update and
                 insert statements.
                 elements between '[' and']'denote an ArrayList.
                 elements between '{' and'}'denote a HashSet.
                 elements between '{' and'}'with key/value pairs separated by
                 ':' denotes a HashMap. A literal HashMap  preceded by a class
                 name are processed by calling a zero argument constructor then
                 followed by each pair key being turned into a setter and
                 invoked with the value.
-f &lt;value&gt;       Each instance of -f followed by a filename load one file of
                 statements.
-g &lt;value&gt;       An optional GAR file to load before running QueryPlus.
                 If the -a argument is not used the application name will be the
                 GAR file name without the parent directory name.
-l &lt;value&gt;       Each instance of -l followed by a statement will execute one
                 statement.
-s               silent mode. Suppress prompts and result headings, read from
                 stdin and write to stdout. Useful for use in pipes or filters
-t               or -trace
                 turn on tracing. This shows information useful for debugging
-dp &lt;list&gt;       A comma delimited list of domain partition names to use.
                 On start-up the first domain partition in the list will be the
                 current partition. The -dp argument is only applicable in
                 combination with the -g argument.
-timeout &lt;value&gt; Specifies the timeout value for CohQL statements in
                 milli-seconds.
BYE |  QUIT
(ENSURE | CREATE) CACHE 'cache-name'
(ENSURE | CREATE) INDEX [ON] 'cache-name' value-extractor-list
DROP CACHE 'cache-name'
TRUNCATE CACHE 'cache-name'
DROP INDEX [ON] 'cache-name' value-extractor-list
BACKUP CACHE 'cache-name' [TO] [FILE] 'filename'
RESTORE CACHE 'cache-name' [FROM] [FILE] 'filename'
INSERT INTO 'cache-name' [KEY (literal | new java-constructor | static method)]
        VALUE (literal |  new java-constructor | static method)
DELETE FROM 'cache-name'[[AS] alias] [WHERE conditional-expression]
UPDATE 'cache-name' [[AS] alias] SET update-statement {, update-statement}*
        [WHERE conditional-expression]
SELECT (properties* aggregators* | * | alias) FROM 'cache-name' [[AS] alias]
        [WHERE conditional-expression] [GROUP [BY] properties+]
SOURCE FROM [FILE] 'filename'
@ 'filename'
. filename
SHOW PLAN 'CohQL command' | EXPLAIN PLAN for 'CohQL command'
TRACE 'CohQL command'
LIST SERVICES [ENVIRONMENT]
LIST [ARCHIVED] SNAPSHOTS ['service']
LIST ARCHIVER 'service'
CREATE SNAPSHOT 'snapshot-name' 'service'
RECOVER SNAPSHOT 'snapshot-name' 'service'
REMOVE [ARCHIVED] SNAPSHOT 'snapshot-name' 'service'
VALIDATE SNAPSHOT 'snapshot-directory' [VERBOSE]
VALIDATE SNAPSHOT 'snapshot-name' 'service-name' [VERBOSE]
VALIDATE ARCHIVED SNAPSHOT 'snapshot-name' 'service-name' [VERBOSE]
ARCHIVE SNAPSHOT 'snapshot-name' 'service'
RETRIEVE ARCHIVED SNAPSHOT 'snapshot-name' 'service' [OVERWIRTE]
RESUME SERVICE 'service'
SUSPEND SERVICE 'service'
FORCE RECOVERY 'service'
COMMANDS
EXTENDED LANGUAGE (ON | OFF)
HELP
SANITY [CHECK] (ON | OFF)
SERVICES INFO
TRACE (ON | OFF)
WHENEVER COHQLERROR THEN (CONTINUE | EXIT)
ALTER SESSION SET DOMAIN PARTITION &lt;partition-name&gt;
ALTER SESSION SET TIMEOUT &lt;milli-seconds&gt;</markup>

</div>

<h3 id="_3_clean_up">3. Clean Up</h3>
<div class="section">
<p>After running the above the Coherence cluster can be removed using <code>kubectl</code>:</p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; delete -f example-cluster.yaml</markup>

</div>
</div>
</doc-view>
