<doc-view>

<h2 id="_polyglot_client_demo">Polyglot Client Demo</h2>
<div class="section">
<p>This example shows how to deploy simple Python, JavaScript or Go applications that connect to Coherence running in Kubernetes.
The Coherence Operator is used to deploy a Coherence cluster and the applications connect via gRPC using the gRPC proxy.</p>

<p>A basic REST serve is written in each language which exposes the following endpoints to create, get, update or remove JSON people objects in the Coherence cluster.</p>

<ul class="ulist">
<li>
<p><code>POST /api/people</code> - create a person</p>

</li>
<li>
<p><code>GET /api/people</code> - return all people</p>

</li>
<li>
<p><code>GET /api/people/{id}</code> - return a single person</p>

</li>
<li>
<p><code>DELETE /api/people/{id}</code> - delete a person</p>

</li>
</ul>
<p>The example shows how to connect to the Coherence cluster from any of the clients via two different methods:</p>

<ol style="margin-left: 15px;">
<li>
From an application deployed with Kubernetes (purple processes shown below, accessed via <code>1. REST Client</code>)

</li>
<li>
From an application outside of Kubernetes using simple port-forward or LBR (yellow gRPC Service, accessed via <code>2. Python, JS or Go Client</code>)

</li>
</ol>
<p><strong>The diagram below outlines the example components.</strong></p>

<div class="admonition note">
<p class="admonition-inline">We use <code>port-forward</code> for this example, but you would normally expose the services via load balancers.</p>
</div>


<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/images/example-overview.png" alt="Service Details"width="100%" />
</v-card-text>
</v-card>

<p>See below for information on the Coherence langauge clients used:</p>

<ul class="ulist">
<li>
<p><a id="" title="" target="_blank" href="https://github.com/oracle/coherence-py-client">Coherence Python Client</a></p>

</li>
<li>
<p><a id="" title="" target="_blank" href="https://github.com/oracle/coherence-js-client">Coherence JavaScript Client</a></p>

</li>
<li>
<p><a id="" title="" target="_blank" href="https://github.com/oracle/coherence-go-client">Coherence Go Client</a></p>

</li>
</ul>
</div>

<h2 id="_what_the_example_will_cover">What the Example Will Cover</h2>
<div class="section">
<ul class="ulist">
<li>
<p><router-link to="#pre" @click.native="this.scrollFix('#pre')">Prerequisites</router-link></p>
<ul class="ulist">
<li>
<p><router-link to="#pre-1" @click.native="this.scrollFix('#pre-1')">Clone the GitHub repository</router-link></p>

</li>
<li>
<p><router-link to="#pre-2" @click.native="this.scrollFix('#pre-2')">Create the coherence-demo namespace</router-link></p>

</li>
<li>
<p><router-link to="#pre-3" @click.native="this.scrollFix('#pre-3')">Install the Coherence Operator</router-link></p>

</li>
<li>
<p><router-link to="#pre-4" @click.native="this.scrollFix('#pre-4')">Download additional software (Optional)</router-link></p>

</li>
</ul>
</li>
<li>
<p><router-link to="#deploy" @click.native="this.scrollFix('#deploy')">Deploy the example</router-link></p>
<ul class="ulist">
<li>
<p><router-link to="#dep-1" @click.native="this.scrollFix('#dep-1')">Examine the Docker files</router-link></p>

</li>
<li>
<p><router-link to="#dep-2" @click.native="this.scrollFix('#dep-2')">Build the example images</router-link></p>

</li>
<li>
<p><router-link to="#dep-3" @click.native="this.scrollFix('#dep-3')">Push images</router-link></p>

</li>
<li>
<p><router-link to="#dep-4" @click.native="this.scrollFix('#dep-4')">Deploy the Coherence Cluster</router-link></p>

</li>
</ul>
</li>
<li>
<p><router-link to="#run-example" @click.native="this.scrollFix('#run-example')">Run the example</router-link></p>
<ul class="ulist">
<li>
<p><router-link to="#run-1" @click.native="this.scrollFix('#run-1')">Run from within Kubernetes</router-link></p>

</li>
<li>
<p><router-link to="#run-2" @click.native="this.scrollFix('#run-2')">View the cache information</router-link></p>

</li>
<li>
<p><router-link to="#run-3" @click.native="this.scrollFix('#run-3')">Run locally using native clients</router-link></p>

</li>
</ul>
</li>
<li>
<p><router-link to="#cleanup" @click.native="this.scrollFix('#cleanup')">Cleaning up</router-link></p>

</li>
</ul>

<h3 id="pre">PreRequisites</h3>
<div class="section">
<p>You must have:</p>

<ol style="margin-left: 15px;">
<li>
Docker running on your system. Either Docker Desktop or Rancher will work

</li>
<li>
Access to a Kubernetes cluster. You can use <code>kind</code> to create a local cluster on Mac. See <a id="" title="" target="_blank" href="https://kind.sigs.k8s.io/">Kind Documentation</a>

</li>
<li>
<code>kubectl</code> executable - See <a id="" title="" target="_blank" href="https://kubernetes.io/docs/tasks/tools/">Kubernetes Documentation</a>

</li>
</ol>
<div class="admonition note">
<p class="admonition-inline">We have <code>Makefile</code> targets that will make the building and running of the example easier.</p>
</div>

<h4 id="pre-1">Clone the Coherence Operator Repository:</h4>
<div class="section">
<p>To build the examples, you first need to clone the Operator GitHub repository to your development machine.</p>

<markup
lang="bash"

>git clone https://github.com/oracle/coherence-operator

cd coherence-operator/examples/910_polyglot_demo</markup>

</div>

<h4 id="pre-2">Create the coherence-demo Namespace</h4>
<div class="section">
<markup
lang="bash"

>make create-namespace</markup>

<markup
lang="bash"
title="Output"
>kubectl create namespace coherence-demo || true
namespace/coherence-demo created</markup>

</div>

<h4 id="pre-3">Install the Coherence Operator</h4>
<div class="section">
<div class="admonition tip">
<p class="admonition-inline">The Coherence Operator is installed into a namespace called <code>coherence</code>. To change this see the documentation below.</p>
</div>
<markup
lang="bash"

>make deploy-operator</markup>

<markup
lang="bash"
title="Output"
>kubectl apply -f https://github.com/oracle/coherence-operator/releases/download/v3.4.3/coherence-operator.yaml

namespace/coherence created
customresourcedefinition.apiextensions.k8s.io/coherence.coherence.oracle.com created
...
service/coherence-operator-webhook created
deployment.apps/coherence-operator-controller-manager created</markup>

<p>You can check to see the Coherence Operator in running by issuing the following command, which will wait for it to be available.</p>

<markup
lang="bash"

>kubectl -n coherence wait --for=condition=available deployment/coherence-operator-controller-manager --timeout 120s</markup>

<p>When you see the following message, you can continue.</p>

<markup
lang="bash"
title="Output"
>deployment.apps/coherence-operator-controller-manager condition met</markup>

<p>See the <router-link to="/docs/installation/001_installation">Installation Guide</router-link> for more information about installing the Coherence Operator.</p>

</div>

<h4 id="pre-4">Download Additional Software (Optional)</h4>
<div class="section">
<p>If you are planning on running the clients locally, E.g. <strong>not just from within the Kuberenetes cluster</strong>, and connect to the cluster, you must install the following
software for your operating system.</p>

<p>Otherwise, you can continue to the next step.</p>

<ol style="margin-left: 15px;">
<li>
Python 3.9 or Later - <a id="" title="" target="_blank" href="https://www.python.org/downloads/">https://www.python.org/downloads/</a>

</li>
<li>
Node 18.15.x or later and NPM 9.x or later - <a id="" title="" target="_blank" href="https://nodejs.org/en/download">https://nodejs.org/en/download</a>

</li>
<li>
Go 1.23 or later - <a id="" title="" target="_blank" href="https://go.dev/doc/install">https://go.dev/doc/install</a>

</li>
</ol>
<div class="admonition tip">
<p class="admonition-inline">If you are just going to run the example within your Kubernetes cluster then you do not need the software.</p>
</div>
</div>
</div>

<h3 id="deploy">Deploy the Example</h3>
<div class="section">

<h4 id="dep-1">Examine the Docker Files</h4>
<div class="section">
<p>Each of the clients has a <code>Dockerfile</code> to build for the specific language client. You can inspect each of them below:</p>

<p><strong>Python Client</strong></p>

<markup
lang="bash"

>FROM python:3.11-slim

RUN addgroup --system appgroup &amp;&amp; adduser --system --ingroup appgroup appuser

WORKDIR /app
COPY --chown=appuser:appgroup main.py .
RUN chmod 444 main.py

RUN pip install --no-cache-dir coherence-client==2.0.0 Quart

RUN chown -R appuser:appgroup /app

USER appuser

CMD ["python3", "./main.py"]</markup>

<p><strong>JavScript Client</strong></p>

<markup
lang="bash"

>FROM node:18-alpine

RUN addgroup -S appgroup &amp;&amp; adduser -S appuser -G appgroup

WORKDIR /usr/src/app

COPY --chown=appuser:appgroup package*.json ./
COPY --chown=appuser:appgroup main.js ./

RUN chmod 444 package*.json main.js

RUN npm install --ignore-scripts
RUN chown -R appuser:appgroup /usr/src/app

USER appuser
EXPOSE 8080

CMD [ "node", "main.js" ]</markup>

<p><strong>Go Client</strong></p>

<markup
lang="bash"

>FROM golang:1.24 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY main.go ./

ENV APP_USER_UID=1000
ENV APP_USER_GID=1000

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o runner .
RUN chown ${APP_USER_UID}:${APP_USER_GID} /app/runner

FROM scratch

COPY --from=builder /app/runner /files/runner
USER 1000:1000

EXPOSE 8080
ENTRYPOINT ["/files/runner"]
CMD ["-h"]</markup>

</div>

<h4 id="code">Examine the Code</h4>
<div class="section">
<p>You can view the source of each of the language clients here:</p>

<ul class="ulist">
<li>
<p><a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/blob/main/examples/910_polyglot_demo/py/main.py">py/main.py</a></p>

</li>
<li>
<p><a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/blob/main/examples/910_polyglot_demo/js/main.js">js/main.js</a></p>

</li>
<li>
<p><a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/blob/main/examples/910_polyglot_demo/go/main.go">go/main.go</a></p>

</li>
</ul>
</div>

<h4 id="dep-2">Build the Example images</h4>
<div class="section">
<div class="admonition tip">
<p class="admonition-inline">If you are deploying to a remote cluster or deploying to a different architecture in Kubernetes, please read the <strong>Important Notes</strong> below before you build.</p>
</div>
<p>Build each of the images using make, or run <code>make create-all-images</code>, as shown below.</p>

<ul class="ulist">
<li>
<p><code>make create-py-image</code></p>

</li>
<li>
<p><code>make create-js-image</code></p>

</li>
<li>
<p><code>make create-go-image</code></p>

</li>
</ul>
<markup
lang="bash"

>make create-all-images</markup>

<markup
lang="bash"
title="Output"
>cd go &amp;&amp; docker buildx build --platform linux/arm64 -t polyglot-client-go:1.0.0 .
[+] Building 27.8s (12/12) FINISHED
....
cd js &amp;&amp; docker buildx build --platform linux/arm64 -t polyglot-client-js:1.0.0 .
[+] Building 4.2s (10/10) FINISHED
...
cd py &amp;&amp; docker buildx build --platform linux/arm64 -t polyglot-client-py:1.0.0 .
[+] Building 3.0s (8/8) FINISHED</markup>

<p>You will end up with the following images:</p>

<ul class="ulist">
<li>
<p><code>polyglot-client-py:1.0.0</code> - Python Client</p>

</li>
<li>
<p><code>polyglot-client-js:1.0.0</code> - JavaScript Client</p>

</li>
<li>
<p><code>polyglot-client-go:1.0.0</code> - Go Client</p>

</li>
</ul>
<p><strong>Important Notes</strong></p>

<ol style="margin-left: 15px;">
<li>
The images are built by default for arm64, you can specify <code>PLATFORM=linux/amd64</code> before your make commands to change the architecture.

</li>
<li>
If you need to push the images to a remote repository, you will need to specify the IMAGE_PREFIX=.. before the make commands, e.g.:
<markup
lang="bash"

>IMAGE_PREFIX=ghcr.io/username/repo/ make create-py-image</markup>

<markup
lang="bash"
title="Output"
>cd py &amp;&amp; docker buildx build --platform linux/arm64 -t ghcr.io/username/repo/polyglot-client-py:1.0.0 .
...</markup>

<div class="admonition note">
<p class="admonition-inline">In this example the image will be <code>ghcr.io/username/repo/polyglot-client-py:1.0.0</code>. You will also need to update the deployment yaml files.</p>
</div>
</li>
</ol>
</div>

<h4 id="dep-3">Push Images</h4>
<div class="section">
<p>Choose one of the following methods, depending upon if you are using a local <strong>kind</strong> cluster or not.</p>

<p><strong>You are running a local cluster using kind</strong></p>

<p>Run the following to load the images you just created to the cluster.</p>

<markup
lang="bash"

>make kind-load-images</markup>

<p><strong>You need to push to a remote repository</strong></p>

<ol style="margin-left: 15px;">
<li>
Push the images to your local container repository. E.g. for the example above:
<markup
lang="bash"

>docker push ghcr.io/username/repo/polyglot-client-py:1.0.0
docker push ghcr.io/username/repo/polyglot-client-js:1.0.0
docker push ghcr.io/username/repo/polyglot-client-go:1.0.0</markup>

</li>
<li>
Modify the following files to change the image name accordingly in the following deployment yaml files:
<ul class="ulist">
<li>
<p><a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/blob/main/examples/910_polyglot_demo/yaml/py-client.yaml">py-client.yaml</a></p>

</li>
<li>
<p><a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/blob/main/examples/910_polyglot_demo/yaml/js-client.yaml">js-client.yaml</a></p>

</li>
<li>
<p><a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/blob/main/examples/910_polyglot_demo/yaml/go-client.yaml">go-client.yaml</a></p>

</li>
</ul>
</li>
<li>
Create a secret if your repository is not public:
<p>If the repository you are pushing to is not public, you will need to create a pull secret, and add this to the deployment yaml for each client.</p>

<markup
lang="bash"

>kubectl create secret docker-registry my-pull-secret \
    --docker-server=ghcr.io \
    --docker-username="&lt;username&gt;" --docker-password="&lt;password&gt;" \
    --docker-email="&lt;email&gt;" -n coherence-demo</markup>

<p>In each of the client deployment files, above add <code>imagePullSecrets</code> after the image. For example in the go-client:</p>

<markup
lang="yaml"

>        - name: go-client
          image: ghcr.io/username/repo/polyglot-client-go:1.0.0
          imagePullPolicy: IfNotPresent
          imagePullSecrets:
            - name: my-pull-secret</markup>

</li>
</ol>
</div>

<h4 id="dep-4">4. Deploy the Coherence Cluster</h4>
<div class="section">
<p>The following deployment file is used to deploy a 3 node Coherence cluster to your Kubernetes cluster.</p>

<markup
lang="bash"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: demo-cluster  <span class="conum" data-value="1" />
spec:
  jvm:
    memory:
      initialHeapSize: 1g
      maxHeapSize: 1g
  replicas: 3 <span class="conum" data-value="2" />
  image: "ghcr.io/oracle/coherence-ce:14.1.2-0-1-java17" <span class="conum" data-value="3" />
  coherence:
    management: <span class="conum" data-value="4" />
      enabled: true
  ports:
    - name: grpc <span class="conum" data-value="5" />
      port: 1408
    - name: management</markup>

<ul class="colist">
<li data-value="1">the cluster name</li>
<li data-value="2">the number of replicas to create</li>
<li data-value="3">Image to use to start Coherence</li>
<li data-value="4">Enable management</li>
<li data-value="5">Enable gRPC port</li>
</ul>
<div class="admonition note">
<p class="admonition-inline">When we deploy this yaml, each of the ports will become a service that we can use to connect to.</p>
</div>
<p>Deploy the Coherence Cluster using the following:</p>

<markup
lang="bash"

>make deploy-coherence</markup>

<markup
lang="bash"
title="Output"
>kubectl -n coherence-demo apply -f yaml/coherence-cluster.yaml
coherence.coherence.oracle.com/demo-cluster created
sleep 5
kubectl -n coherence-demo get pods
NAME             READY   STATUS    RESTARTS   AGE
demo-cluster-0   0/1     Running   0          5s
demo-cluster-1   0/1     Running   0          5s
demo-cluster-2   0/1     Running   0          5s</markup>

<p>Issue the following command to wait until the Coherence cluster is ready.</p>

<markup
lang="bash"

>kubectl -n coherence-demo wait --for=condition=ready coherence/demo-cluster --timeout 120s
coherence.coherence.oracle.com/demo-cluster condition met</markup>

</div>
</div>

<h3 id="run-example">Run the example</h3>
<div class="section">
<p>Choose one of the methods below to demo the application.</p>

<ol style="margin-left: 15px;">
<li>
<router-link to="#run-1" @click.native="this.scrollFix('#run-1')">Run from within Kubernetes</router-link>

</li>
<li>
<router-link to="#run-3" @click.native="this.scrollFix('#run-3')">Run locally using native clients</router-link>

</li>
</ol>

<h4 id="run-1">Run from within Kubernetes</h4>
<div class="section">
<p>In this section, we will deploy all the client pods to Kubernetes, and access the REST endpoints via their services, using the <code>port-forward</code> command.</p>

<p>See <strong>1. REST Client</strong> on the right of the diagram below:</p>



<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/images/example-overview.png" alt="Service Details"width="100%" />
</v-card-text>
</v-card>

<div class="admonition note">
<p class="admonition-inline">We are accessing the application HTTP endpoint via port-forward, but the client pods within Kubernetes
are directly accessing cluster within Kubernetes.</p>
</div>
<p>When we deploy the clients, the yaml used is shown below:</p>

<ul class="ulist">
<li>
<p><a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/blob/main/examples/910_polyglot_demo/yaml/py-client.yaml">Python Client</a></p>

</li>
<li>
<p><a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/blob/main/examples/910_polyglot_demo/yaml/js-client.yaml">JavaScript Client</a></p>

</li>
<li>
<p><a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/blob/main/examples/910_polyglot_demo/yaml/go-client.yaml">Go Client</a></p>

</li>
</ul>
<p>By default, the Python, JavaScript and Go clients connect to <code>localhost:1408</code> on startup, but you can specify the gRPC host and port to connect to in the
code or use the <code>COHERENCE_SERVER_ADDRESS</code> environment variable to specify this which is more flexible.</p>

<p>Each of the clients, (Python shown below for example), have this variable set to <code>demo-cluster-grpc:1408</code> where <code>demo-cluster-grpc</code>
is the service for the grpc port created when we deployed the Coherence cluster.</p>

<markup
lang="yaml"

>- name: py-client
  image: polyglot-client-py:1.0.0
  imagePullPolicy: IfNotPresent
  env:
    - name: COHERENCE_SERVER_ADDRESS
      value: "demo-cluster-grpc:1408"
    - name: COHERENCE_READY_TIMEOUT
      value: "60000"
  resources:
    requests:
      memory: "512Mi"
    limits:
      memory: "512Mi"
  ports:
    - containerPort: 8080
  securityContext:
    runAsNonRoot: true
    runAsUser: 10001
    capabilities:
      drop:
        - all
    readOnlyRootFilesystem: true</markup>

<p><strong>Deploy the clients</strong></p>

<p>Firstly, run the following to deploy all the clients. This yaml also deploys a service from which you can connect to the client.</p>

<markup
lang="bash"

>make deploy-all-clients</markup>

<markup
lang="bash"
title="Output"
>kubectl -n coherence-demo apply -f yaml/go-client.yaml
service/go-client created
deployment.apps/go-client created
kubectl -n coherence-demo apply -f yaml/js-client.yaml
service/js-client created
deployment.apps/js-client created
kubectl -n coherence-demo apply -f yaml/py-client.yaml
service/py-client created
deployment.apps/py-client created</markup>

<p>Issue the following to show the deployments and services.</p>

<markup
lang="bash"

>kubectl get deployments -n coherence-demo</markup>

<markup
lang="bash"
title="Output"
>NAME        READY   UP-TO-DATE   AVAILABLE   AGE
go-client   1/1     1            1           3s
js-client   1/1     1            1           3s
py-client   1/1     1            1           3s</markup>

<p><strong>Port forward to access the HTTP endpoint</strong></p>

<p>To port-forward the clients we will first need to view the services:</p>

<markup
lang="bash"

>kubectl get services -n coherence-demo</markup>

<markup
lang="bash"
title="Output"
>NAME                      TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                                               AGE
demo-cluster-grpc         ClusterIP   10.96.200.57    &lt;none&gt;        1408/TCP                                              8m2s
demo-cluster-management   ClusterIP   10.96.46.69     &lt;none&gt;        30000/TCP                                             8m2s
demo-cluster-sts          ClusterIP   None            &lt;none&gt;        7/TCP,7575/TCP,7574/TCP,6676/TCP,1408/TCP,30000/TCP   8m2s
demo-cluster-wka          ClusterIP   None            &lt;none&gt;        7/TCP,7575/TCP,7574/TCP,6676/TCP                      8m2s
go-client-http            ClusterIP   10.96.249.42    &lt;none&gt;        8080/TCP                                              32s
js-client-http            ClusterIP   10.96.114.88    &lt;none&gt;        8080/TCP                                              31s
py-client-http            ClusterIP   10.96.196.163   &lt;none&gt;        8080/TCP                                              31s</markup>

<p>The services we are interested in are the <code>go-client-http</code>, <code>js-client-http</code> or <code>py-client-http</code>.</p>

<p>As all the clients expose the same API, You can choose any of the clients to port-forward to. For this example will
choose the JavaScript client.</p>

<markup
lang="bash"

>kubectl port-forward -n coherence-demo svc/js-client-http 8080:8080</markup>

<markup
lang="bash"
title="Output"
>Forwarding from 127.0.0.1:8080 -&gt; 8080
Forwarding from [::1]:8080 -&gt; 8080</markup>

<p id="rest-endpoints"><strong>Exercise the REST endpoints</strong></p>

<p>Use the following commands to work with the REST endpoints.</p>

<p><strong>Create two People</strong></p>

<markup
lang="bash"

>curl -X POST -H 'Content-Type: application/json' http://localhost:8080/api/people -d '{"id": 1,"name":"Tim", "age": 25}'</markup>

<markup
lang="bash"

>curl -X POST -H 'Content-Type: application/json' http://localhost:8080/api/people -d '{"id": 2,"name":"John", "age": 35}'</markup>

<p><strong>List the people in the cache</strong></p>

<markup
lang="bash"

>curl http://localhost:8080/api/people</markup>

<markup
lang="bash"
title="Output"
>[{"id":1,"name":"Tim","age":25},{"id":2,"name":"John","age":35}]</markup>

<p><strong>Get a single Person</strong></p>

<markup
lang="bash"

>curl http://localhost:8080/api/people/1</markup>

<markup
lang="bash"
title="Output"
>{"id":1,"name":"Tim","age":25}</markup>

<p><strong>Remove a Person</strong></p>

<markup
lang="bash"

>curl -X DELETE http://localhost:8080/api/people/1</markup>

<markup
lang="bash"
title="Output"
>OK</markup>

<p><strong>Try to retrieve the deleted person</strong></p>

<markup
lang="bash"

>curl http://localhost:8080/api/people/1</markup>

<markup
lang="bash"
title="Output"
>Not Found</markup>

</div>

<h4 id="run-2">View the cache information</h4>
<div class="section">
<p>You can use the <a id="" title="" target="_blank" href="https://github.com/oracle/coherence-cli">Coherence CLI</a> to view cluster and cache information.</p>

<div class="admonition note">
<p class="admonition-inline">The Coherence CLI is automatically bundled with the Coherence Operator and can be accessed via <code>kuebctl exec</code>.</p>
</div>
<p><strong>Display the cluster members</strong></p>

<p>Use the following to show the cluster members:</p>

<markup
lang="bash"

>kubectl exec demo-cluster-0 -c coherence -n coherence-demo -- /coherence-operator/utils/cohctl get members</markup>

<markup
lang="bash"
title="Output"
>Total cluster members: 3
Storage enabled count: 3
Departure count:       0

Cluster Heap - Total: 8,964 MB Used: 150 MB Available: 8,814 MB (98.3%)
Storage Heap - Total: 8,964 MB Used: 150 MB Available: 8,814 MB (98.3%)

NODE ID  ADDRESS                                          PORT  PROCESS  MEMBER          ROLE          STORAGE  MAX HEAP  USED HEAP  AVAIL HEAP
      1  demo-cluster-wka.coherence-demo.svc/10.244.0.31  7575       56  demo-cluster-0  demo-cluster  true     2,988 MB      45 MB    2,943 MB
      2  demo-cluster-wka.coherence-demo.svc/10.244.0.32  7575       57  demo-cluster-2  demo-cluster  true     2,988 MB      36 MB    2,952 MB
      3  demo-cluster-wka.coherence-demo.svc/10.244.0.33  7575       57  demo-cluster-1  demo-cluster  true     2,988 MB      69 MB    2,919 MB</markup>

<p><strong>Display the cache information</strong></p>

<p>Use the following to show the cache information:</p>

<markup
lang="bash"

>kubectl exec demo-cluster-0 -c coherence -n coherence-demo -- /coherence-operator/utils/cohctl get caches -o wide</markup>

<div class="admonition note">
<p class="admonition-inline">You will see other system caches, but you should also see the <code>people</code> cache with one entry.</p>
</div>
<markup
lang="bash"
title="Output"
>Total Caches: 7, Total primary storage: 0 MB

SERVICE            CACHE                 COUNT  SIZE  AVG SIZE    PUTS    GETS  REMOVES  EVICTIONS    HITS   MISSES  HIT PROB
"$SYS:Concurrent"  executor-assignments      0  0 MB         0       0       0        0          0       0        0     0.00%
"$SYS:Concurrent"  executor-executors        2  0 MB     1,248  12,105  12,103        1          0  12,103        0   100.00%
"$SYS:Concurrent"  executor-tasks            0  0 MB         0       0       0        0          0       0        0     0.00%
"$SYS:Concurrent"  locks-exclusive           0  0 MB         0       0       0        0          0       0        0     0.00%
"$SYS:Concurrent"  locks-read-write          0  0 MB         0       0       0        0          0       0        0     0.00%
"$SYS:Concurrent"  semaphores                0  0 MB         0       0       0        0          0       0        0     0.00%
PartitionedCache   people                    1  0 MB       224       6      29        2          0      26        3    89.66%</markup>

</div>

<h4 id="run-3">Run locally using native clients</h4>
<div class="section">
<p>Depending upon the client you are wanting to run you need to ensure you have installed the relevant
client software as shown in the <router-link to="#pre-4" @click.native="this.scrollFix('#pre-4')">pre-requisites here</router-link>.</p>

<div class="admonition tip">
<p class="admonition-inline">Ensure you stopped the port-forward from the previous step if you ran this.</p>
</div>
<p>See <strong>2. Python, JS or GO Client</strong> on the bottom of the diagram below:</p>



<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/images/example-overview.png" alt="Service Details"width="100%" />
</v-card-text>
</v-card>

<p><strong>Run Port forward</strong></p>

<p>Firstly we have to run a <code>port-forward</code> command to port-forward the gRPC 1408 locally to the <code>demo-cluster-grpc:1408</code> port on the Kubernetes cluster.</p>

<p>Typically, if you want to access a service from outside the Kubernetes cluster, you would have a load balancer configured to access this directly, but in example we will use <code>port-forward</code>.</p>

<p>Run the following:</p>

<markup
lang="bash"

>kubectl -n coherence-demo port-forward svc/demo-cluster-grpc 1408:1408</markup>

<markup
lang="bash"
title="Output"
>Forwarding from 127.0.0.1:1408 -&gt; 1408
Forwarding from [::1]:1408 -&gt; 1408</markup>

<p>the <strong>DIAGARM</strong>&#8230;&#8203;</p>

<p>Then follow the instructions to start either of the clients below, which will list on port 8080 locally and connect to the Coherence cluster via gRPC on localhost:1408 which will be port-forwarded.</p>

<p><strong>Python Client</strong></p>

<div class="admonition note">
<p class="admonition-inline">We are install a python virtual environment for this example.</p>
</div>
<ol style="margin-left: 15px;">
<li>
Change to the <code>py</code> directory

</li>
<li>
Create a python virtual environment
<markup
lang="bash"

>python3 -m venv ./venv
. venv/bin/activate</markup>

</li>
<li>
Install the requirements
<markup
lang="bash"

>python3 -m pip install -r requirements.txt</markup>

</li>
<li>
Run the Python example
<markup
lang="bash"

>python3 main.py</markup>

<markup
lang="bash"
title="Output"
>2025-04-07 11:06:42,501 - coherence - INFO - Session [5d940a05-1cfc-4e6c-9ef8-52cc6e7705ba] connected to [localhost:1408].
2025-04-07 11:06:42,525 - coherence - INFO - Session(id=5d940a05-1cfc-4e6c-9ef8-52cc6e7705ba, connected to [localhost:1408] proxy-version=14.1.2.0.1, protocol-version=1 proxy-member-id=1)
[2025-04-07 11:06:42 +0800] [27645] [INFO] Running on http://0.0.0.0:8080 (CTRL + C to quit)</markup>

<div class="admonition note">
<p class="admonition-inline">This is now showing the HTTP server is running locally and connecting via port-forward to the Coherence Cluster.</p>
</div>
</li>
<li>
Exercise the REST end-points as per the instructions <router-link to="#rest-endpoints" @click.native="this.scrollFix('#rest-endpoints')">here</router-link>

</li>
</ol>
<p><strong>JavaScript Client</strong></p>

<ol style="margin-left: 15px;">
<li>
Change to the <code>js</code> directory

</li>
<li>
Install the modules
<markup
lang="bash"

>npm install</markup>

</li>
<li>
Run the JavaScript example
<markup
lang="bash"

>node main.js</markup>

</li>
<li>
Exercise the REST end-points as per the instructions <router-link to="#rest-endpoints" @click.native="this.scrollFix('#rest-endpoints')">here</router-link>

</li>
</ol>
<p><strong>Go Client</strong></p>

<ol style="margin-left: 15px;">
<li>
Change to the <code>go</code> directory

</li>
<li>
Ensure you have the latest Coherence Go client
<markup
lang="bash"

>npm install</markup>

</li>
<li>
Build the executable
<markup
lang="bash"

>go get github.com/oracle/coherence-go-client/v2@latest
go mod tidy</markup>

</li>
<li>
Run the Go example
<markup
lang="bash"

>/runner</markup>

<markup
lang="bash"
title="Output"
>2025/04/07 11:19:21 INFO: Session [2073aa45-68aa-426d-a0b8-99405dcaa942] connected to [localhost:1408] Coherence version: 14.1.2.0.1, serverProtocolVersion: 1, proxyMemberId: 1
Server running on port 8080</markup>

</li>
<li>
Exercise the REST end-points as per the instructions <router-link to="#rest-endpoints" @click.native="this.scrollFix('#rest-endpoints')">here</router-link>

</li>
</ol>
</div>
</div>

<h3 id="cleanup">Cleaning Up</h3>
<div class="section">
<ol style="margin-left: 15px;">
<li>
Undeploy the clients using:
<markup
lang="bash"

>make undeploy-all-clients</markup>

</li>
<li>
Undeploy the Coherence cluster
<markup
lang="bash"

>make undeploy-coherence</markup>

</li>
<li>
Undeploy the Coherence Operator
<markup
lang="bash"

>make undeploy-operator</markup>

</li>
<li>
Delete the namespace
<markup
lang="bash"

>make delete-namespace</markup>

</li>
</ol>
</div>
</div>
</doc-view>
