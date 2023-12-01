<doc-view>

<h2 id="_secure_coherence_extend_with_tls">Secure Coherence Extend with TLS</h2>
<div class="section">
<p>This example shows how to deploy a simple Coherence cluster in Kubernetes manually, and secure the Extend endpoint using TLS.
This example expands on the <code>StatefulSet</code> used in the first simple deployment example.</p>

<div class="admonition tip">
<p class="admonition-textlabel">Tip</p>
<p ><p><img src="./images/GitHub-Mark-32px.png" alt="GitHub Mark 32px" />
 The complete source code for this example is in the <a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/tree/main/examples/no-operator/03_extend_tls">Coherence Operator GitHub</a> repository.</p>
</p>
</div>
<p><strong>Prerequisites</strong></p>

<p>This example assumes that you have already built the example server image.</p>

<p>There are a number of ways to use TLS to secure ingress in Kubernetes. We could use a load balancer <code>Service</code> and terminate TLS at the load balance, or we could use an add-on such as Istio to manage TLS ingress. Both of those approaches would require no changes to the Coherence server, as the server would not know TLS was being used.
The <a id="" title="" target="_blank" href="https://oracle.github.io/coherence-operator/docs/latest/#/examples/010_overview">Coherence Operator Examples</a>
contains examples of using TLS with Coherence and using Istio. The TLS example also shows how to use Kubernetes built in certificate management to create keys and certificates.</p>

<p>In this example we are going to actually change the server to use TLS for its Extend endpoints.</p>

</div>

<h2 id="_create_certs_and_java_keystores">Create Certs and Java Keystores</h2>
<div class="section">
<p>To use TLS we will need some certificates and Java keystore files. For testing and examples, self-signed certs are fine. The source code for this example contains some keystores.
* <code>server.jks</code> contains the server key and certificate files
* <code>trust.jks</code> contains the CA certificate used to create the client and server certificates</p>

<p>The keystores are password protected, the passwords are stored in files with the example source.
We will use these files to securely provide the passwords to the client and server instead of hard coding or providing credentials via system properties or environment variables.
* <code>server-password.txt</code> is the password to open the <code>server.jks</code> keystore
* <code>server-key-password.txt</code> is the password for the key file stored in the <code>server.jks</code> keystore
* <code>trust-password.txt</code> is the password to open the trust.jks` keystore.</p>

</div>

<h2 id="_configure_coherence_extend_tls">Configure Coherence Extend TLS</h2>
<div class="section">
<p>The Coherence documentation explains how to
<a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/standalone/coherence/14.1.1.0/secure/using-ssl-secure-communication.html#GUID-90E20139-3945-4993-9048-7FBC93B243A3">Use TLS Secure Communication</a>.
This example is going to use a standard approach to securing Extend with TLS. To provide the keystores and credentials the example will make use of Kubernetes <code>Secrets</code> to mount those files as <code>Volumes</code> in the <code>StatefulSet</code>. This is much more flexible and secure than baking them into an application&#8217;s code or image.</p>


<h3 id="_configure_the_extend_proxy">Configure the Extend Proxy</h3>
<div class="section">
<p>If we look at the <code>test-cache-config.xml</code> file in the <code>simple-server</code> example project, we can see the configuration for the Extend proxy.</p>

<markup
lang="xml"
title="test-cache-config.xml"
>    &lt;proxy-scheme&gt;
      &lt;service-name&gt;Proxy&lt;/service-name&gt;
      &lt;acceptor-config&gt;
        &lt;tcp-acceptor&gt;
          &lt;socket-provider system-property="coherence.extend.socket.provider"/&gt;
          &lt;local-address&gt;
            &lt;address system-property="coherence.extend.address"&gt;0.0.0.0&lt;/address&gt;
            &lt;port system-property="coherence.extend.port"&gt;20000&lt;/port&gt;
          &lt;/local-address&gt;
        &lt;/tcp-acceptor&gt;
      &lt;/acceptor-config&gt;
      &lt;autostart&gt;true&lt;/autostart&gt;
    &lt;/proxy-scheme&gt;</markup>

<p>The important item to note above is the <code>socket-provider</code> element, which is empty, but can be set using the <code>coherence.extend.socket.provider</code> system property (or the <code>COHERENCE_EXTEND_SOCKET_PROVIDER</code> environment variable). By default, a plain TCP socket will be used, but by setting the specified property a different socket can be used, in this case we&#8217;ll use one configured for TLS.</p>

</div>

<h3 id="_socket_providers">Socket Providers</h3>
<div class="section">
<p>In Coherence, socket providers can be configured in the operational configuration file, typically named <code>tangosol-coherence-override.xml</code>. The source code for the <code>simple-server</code> module contains this file with the TLS socket provider already configured.</p>

<p>We need to configure two things in the operational configuration file, the socket provider and some password providers to supply the keystore credentials.</p>

<p>The <code>socket-provider</code> section looks like this:</p>

<markup
lang="xml"
title="tangosol-coherence-override.xml"
>&lt;socket-providers&gt;
    &lt;socket-provider id="extend-tls"&gt;
        &lt;ssl&gt;
            &lt;protocol&gt;TLS&lt;/protocol&gt;
            &lt;identity-manager&gt;
                &lt;algorithm&gt;SunX509&lt;/algorithm&gt;
                &lt;key-store&gt;
                    &lt;url system-property="coherence.extend.keystore"&gt;file:server.jks&lt;/url&gt;
                    &lt;password-provider&gt;
                        &lt;name&gt;identity-password-provider&lt;/name&gt;
                    &lt;/password-provider&gt;
                    &lt;type&gt;JKS&lt;/type&gt;
                &lt;/key-store&gt;
                &lt;password-provider&gt;
                    &lt;name&gt;key-password-provider&lt;/name&gt;
                &lt;/password-provider&gt;
            &lt;/identity-manager&gt;
            &lt;trust-manager&gt;
                &lt;algorithm&gt;SunX509&lt;/algorithm&gt;
                &lt;key-store&gt;
                    &lt;url system-property="coherence.extend.truststore"&gt;file:trust.jks&lt;/url&gt;
                    &lt;password-provider&gt;
                        &lt;name&gt;trust-password-provider&lt;/name&gt;
                    &lt;/password-provider&gt;
                    &lt;type&gt;JKS&lt;/type&gt;
                &lt;/key-store&gt;
            &lt;/trust-manager&gt;
            &lt;socket-provider&gt;tcp&lt;/socket-provider&gt;
        &lt;/ssl&gt;
    &lt;/socket-provider&gt;
&lt;/socket-providers&gt;</markup>

<p>There is a <code>socket-provider</code> with the id of <code>extend-tls</code>. This id is the value that must be used to tell the Extend proxy which socket provider to use, i.e. using the system property <code>-Dcoherence.extend.socket.provider=extend-tls</code></p>

<p>The <code>&lt;identity-manager&gt;</code> element specifies the keystore containing the key and certificate file that the proxy should use. This is set to <code>file:server.jks</code> but can be overridden using the <code>coherence.extend.keystore</code> system property, or corresponding environment variable. The password for the <code>&lt;identity-manager&gt;</code> keystore is configured to be provided by the <code>password-provider</code> named <code>identity-password-provider</code>. The password for the key file in the identity keystore is configured to be provided by the <code>password-provider</code> named <code>key-password-provider</code>.</p>

<p>The <code>&lt;trust-manager&gt;</code> element contains the configuration for the trust keystore containing the CA certs used to validate client certificates. By default, the keystore name is <code>file:trust.jks</code> but this can be overridden using the <code>coherence.extend.truststore</code> system property or corresponding environment variable. The password for the trust keystore is configured to be provided by the <code>password-provider</code> named <code>trust-password-provider</code>.</p>

<p>There are three <code>&lt;password-provider&gt;</code> elements in the configuration above, so we need to also configure these three password providers in the operational configuration file.</p>

<markup
lang="xml"
title="tangosol-coherence-override.xml"
>&lt;password-providers&gt;
    &lt;password-provider id="trust-password-provider"&gt;
        &lt;class-name&gt;com.oracle.coherence.examples.tls.FileBasedPasswordProvider&lt;/class-name&gt;
        &lt;init-params&gt;
            &lt;init-param&gt;
                &lt;param-name&gt;fileName&lt;/param-name&gt;
                &lt;param-value system-property="coherence.trust.password.file"&gt;trust-password.txt&lt;/param-value&gt;
            &lt;/init-param&gt;
        &lt;/init-params&gt;
    &lt;/password-provider&gt;
    &lt;password-provider id="identity-password-provider"&gt;
        &lt;class-name&gt;com.oracle.coherence.examples.tls.FileBasedPasswordProvider&lt;/class-name&gt;
        &lt;init-params&gt;
            &lt;init-param&gt;
                &lt;param-name&gt;fileName&lt;/param-name&gt;
                &lt;param-value system-property="coherence.identity.password.file"&gt;server-password.txt&lt;/param-value&gt;
            &lt;/init-param&gt;
        &lt;/init-params&gt;
    &lt;/password-provider&gt;
    &lt;password-provider id="key-password-provider"&gt;
        &lt;class-name&gt;com.oracle.coherence.examples.tls.FileBasedPasswordProvider&lt;/class-name&gt;
        &lt;init-params&gt;
            &lt;init-param&gt;
                &lt;param-name&gt;fileName&lt;/param-name&gt;
                &lt;param-value system-property="coherence.key.password.file"&gt;server-key-password.txt&lt;/param-value&gt;
            &lt;/init-param&gt;
        &lt;/init-params&gt;
    &lt;/password-provider&gt;
&lt;/password-providers&gt;</markup>

<p>There are three password providers declared above, each with an 'id' attribute corresponding to the names used in the socket provider configuration. Each password provider is identical, they just have a different password file name.</p>

<p>The <code>class-name</code> element refers to a class named <code>com.oracle.coherence.examples.tls.FileBasedPasswordProvider</code>, which is in the source code for both the server and client. This is an implementation of the <code>com.tangosol.net.PasswordProvider</code> interface which can read a password from a file.</p>

<p>Each password provider&#8217;s password file name can be set using the relevant system property or environment variable. The name of the trust keystore password file is set using the <code>coherence.trust.password.file</code> system property. The name of the identity keystore is set using the <code>coherence.identity.password.file</code> system property. The nam eof the identity key file password file is set using the <code>coherence.key.password.file</code> system property.</p>

<p>The simple server image has all the configuration above built in so there is nothing additional to do to use TLS other than set the system properties or environment variables. The test client uses the same configurations, so it can also be run using TLS by setting the relevant system properties.</p>

</div>
</div>

<h2 id="_create_the_kubernetes_resources">Create the Kubernetes Resources</h2>
<div class="section">
<p>We can now create the resources we need to run the Cluster with TLS enabled.</p>


<h3 id="_keystore_secret">Keystore Secret</h3>
<div class="section">
<p>We first need to supply the keystores and credentials to the Coherence <code>Pods</code>. The secure way to do this in Kubernetes is to use a <code>Secret</code>. We can create a <code>Secret</code> from the command line using <code>kubectl</code>. From the <code>03_extend_tls/</code> directory containing the keystores and password file srun the following command:</p>

<markup
lang="bash"

>kubectl create secret generic coherence-tls \
    --from-file=./server.jks \
    --from-file=./server-password.txt \
    --from-file=./server-key-password.txt \
    --from-file=./trust.jks \
    --from-file=./trust-password.txt</markup>

<p>The command above will create a <code>Secret</code> named <code>coherence-tls</code> containing the files specified. We can now use the <code>Secret</code> in the cluster&#8217;s <code>StatefulSet</code></p>

</div>

<h3 id="_statefulset">StatefulSet</h3>
<div class="section">
<p>We will expand on the <code>StatefulSet</code> created in the simple server example and add TLS.</p>

<markup
lang="yaml"
title="coherence-tls.yaml"
>apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: storage
  labels:
    coherence.oracle.com/cluster: test-cluster
    coherence.oracle.com/deployment: storage
    coherence.oracle.com/component: statefulset
spec:
  selector:
    matchLabels:
      coherence.oracle.com/cluster: test-cluster
      coherence.oracle.com/deployment: storage
  serviceName: storage-sts
  replicas: 3
  template:
    metadata:
      labels:
        coherence.oracle.com/cluster: test-cluster
        coherence.oracle.com/deployment: storage
    spec:
      volumes:
        - name: tls
          secret:
            secretName: coherence-tls
      containers:
        - name: coherence
          image: simple-coherence:1.0.0
          volumeMounts:
            - mountPath: /certs
              name: tls
          command:
            - java
          args:
            - -cp
            - "@/app/jib-classpath-file"
            - -Xms1800m
            - -Xmx1800m
            - "@/app/jib-main-class-file"
          env:
            - name: COHERENCE_CLUSTER
              value: storage
            - name: COHERENCE_WKA
              value: storage-wka
            - name: COHERENCE_CACHECONFIG
              value: test-cache-config.xml
            - name: COHERENCE_EXTEND_SOCKET_PROVIDER
              value: extend-tls
            - name: COHERENCE_EXTEND_KEYSTORE
              value: file:/certs/server.jks
            - name: COHERENCE_IDENTITY_PASSWORD_FILE
              value: /certs/server-password.txt
            - name: COHERENCE_KEY_PASSWORD_FILE
              value: /certs/server-key-password.txt
            - name: COHERENCE_EXTEND_TRUSTSTORE
              value: file:/certs/trust.jks
            - name: COHERENCE_TRUST_PASSWORD_FILE
              value: /certs/trust-password.txt
          ports:
            - name: extend
              containerPort: 20000</markup>

<p>The yaml above is identical to the simple server example with the following additions:</p>

<ul class="ulist">
<li>
<p>A <code>Volume</code> has been added to the <code>spec</code> section.</p>

</li>
</ul>
<div class="listing">
<pre>volumes:
- name: tls
  secret:
    secretName: coherence-tls</pre>
</div>

<p>The volume name is <code>tls</code> and the files to mount to the file system in the Pod come from the <code>coherence-tls</code> secret we created above.</p>

<ul class="ulist">
<li>
<p>A <code>volumeMount</code> has been added to the Coherence container to map the <code>tls</code> volume to the mount point <code>/certs</code>.</p>

</li>
</ul>
<div class="listing">
<pre>volumeMounts:
  - mountPath: /certs
    name: tls</pre>
</div>

<ul class="ulist">
<li>
<p>A number of environment variables have been added to configure Coherence to use the <code>extend-tls</code> socket provider and the locations of the keystores and password files.</p>

</li>
</ul>
<div class="listing">
<pre>- name: COHERENCE_EXTEND_SOCKET_PROVIDER
  value: extend-tls
- name: COHERENCE_EXTEND_KEYSTORE
  value: file:/certs/server.jks
- name: COHERENCE_IDENTITY_PASSWORD_FILE
  value: /certs/server-password.txt
- name: COHERENCE_KEY_PASSWORD_FILE
  value: /certs/server-key-password.txt
- name: COHERENCE_EXTEND_TRUSTSTORE
  value: file:/certs/trust.jks
- name: COHERENCE_TRUST_PASSWORD_FILE
  value: /certs/trust-password.txt</pre>
</div>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>The <code>COHERENCE_EXTEND_KEYSTORE</code> and <code>COHERENCE_EXTEND_TRUSTSTORE</code> values must be URLs. In this case we refer to files usinf the <code>file:</code> prefix.</p>
</p>
</div>
</div>
</div>

<h2 id="_deploy_to_kubernetes">Deploy to Kubernetes</h2>
<div class="section">
<p>The source code for this example contains a file named <code>coherence-tls.yaml</code> containing all the configuration above as well as the <code>Services</code> required to run Coherence and expose the Extend port.</p>

<p>We can deploy it with the following command:</p>

<markup
lang="bash"

>kubectl apply -f coherence-tls.yaml</markup>

<p>We can see all the resources created in Kubernetes are the same as for the simple server example.</p>

<markup
lang="bash"

>kubectl get all</markup>

<p>Which will display something like the following:</p>

<markup


>NAME            READY   STATUS    RESTARTS   AGE
pod/storage-0   1/1     Running   0          19s
pod/storage-1   1/1     Running   0          17s
pod/storage-2   1/1     Running   0          16s

NAME                     TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)     AGE
service/storage-extend   ClusterIP   10.105.78.34   &lt;none&gt;        20000/TCP   19s
service/storage-sts      ClusterIP   None           &lt;none&gt;        7/TCP       19s
service/storage-wka      ClusterIP   None           &lt;none&gt;        7/TCP       19s

NAME                       READY   AGE
statefulset.apps/storage   3/3     19s</markup>

</div>

<h2 id="_run_the_client">Run the Client</h2>
<div class="section">
<p>If we run the test client using the same instructions as the simple server example, we will run an interactive Coherence console.</p>

<markup
lang="bash"

>cd test-client/
mvn exec:java</markup>

<p>When the <code>Map (?):</code> prompt is displayed we can try to create a cache.</p>

<markup


>Map (?): cache test</markup>

<p>This will not throw an exception because the client is not using TLS so the server rejected the connection.</p>

<markup


>2021-09-17 18:19:39.182/12.090 Oracle Coherence CE 21.12.1 &lt;Error&gt; (thread=com.tangosol.net.CacheFactory.main(), member=1): Error while starting service "RemoteCache": com.tangosol.net.messaging.ConnectionException: could not establish a connection to one of the following addresses: [127.0.0.1:20000]
	at com.tangosol.coherence.component.util.daemon.queueProcessor.service.peer.initiator.TcpInitiator.openConnection(TcpInitiator.CDB:139)
	at com.tangosol.coherence.component.util.daemon.queueProcessor.service.peer.Initiator.ensureConnection(Initiator.CDB:11)
	at com.tangosol.coherence.component.net.extend.remoteService.RemoteCacheService.openChannel(RemoteCacheService.CDB:7)
	at com.tangosol.coherence.component.net.extend.RemoteService.ensureChannel(RemoteService.CDB:6)
	at com.tangosol.coherence.component.net.extend.RemoteService.doStart(RemoteService.CDB:11)</markup>


<h3 id="_enable_client_tls">Enable Client TLS</h3>
<div class="section">
<p>Just like the server, the example test client contains the same operational configuration to configure a socket provider and password providers. The test client directory also contains copies of the keystores and password files. We can therefore run the client with the relevant system properties to enable it to use TLS and connect to the server.</p>

<p>We just need to run the client from the <code>test-client/</code> directory setting the socket provider system property.</p>

<markup
lang="bash"

>cd test-client/
mvn exec:java -Dcoherence.extend.socket.provider=extend-tls</markup>

<p>After the client starts we can run the <code>cache</code> command, which should complete without
an error.</p>

<markup


>Map (?): cache test</markup>

<p>We can see from the output below that the client connected and created a remote cache.</p>

<markup


>Cache Configuration: test
  SchemeName: remote
  ServiceName: RemoteCache
  ServiceDependencies: DefaultRemoteCacheServiceDependencies{RemoteCluster=null, RemoteService=Proxy, InitiatorDependencies=DefaultTcpInitiatorDependencies{EventDispatcherThreadPriority=10, RequestTimeoutMillis=30000, SerializerFactory=null, TaskHungThresholdMillis=0, TaskTimeoutMillis=0, ThreadPriority=10, WorkerThreadCount=0, WorkerThreadCountMax=2147483647, WorkerThreadCountMin=0, WorkerThreadPriority=5}{Codec=null, FilterList=[], PingIntervalMillis=0, PingTimeoutMillis=30000, MaxIncomingMessageSize=0, MaxOutgoingMessageSize=0}{ConnectTimeoutMillis=30000, RequestSendTimeoutMillis=30000}{LocalAddress=null, RemoteAddressProviderBldr=com.tangosol.coherence.config.builder.WrapperSocketAddressProviderBuilder@5431b4b4, SocketOptions=SocketOptions{LingerTimeout=0, KeepAlive=true, TcpNoDelay=true}, SocketProvideBuilderr=com.tangosol.coherence.config.builder.SocketProviderBuilder@52c85af7, isNameServiceAddressProvider=false}}{DeferKeyAssociationCheck=false}

Map (test):</markup>

<p>Now the client is connected using TLS, we could do puts and gets, or other operations on the cache.</p>

<p>To exit from the client press ctrl-C, and uninstall the cluster</p>

<markup
lang="bash"

>kubectl delete -f coherence-tls.yaml</markup>

</div>
</div>
</doc-view>
