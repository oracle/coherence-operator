# Coherence K8s Operator

This module builds the Coherence Kubernetes Operator and Helm Charts.


## Build the workspace
- Check that helm 2.11 or greater is in classpath and helm init has been run.
- Build the workspace without Docker image
    ```
    mvn clean install
    ```
- Build the workspace with Docker image
    ```
    mvn -P docker clean install
    ```

## Install with Helm Chart
- The built Helm chart can be found in `operator/target/coherence-operator-${version}-helm`
    ```
    cd operator/target/coherence-operator-*-helm
    ```
- Install Coherence Kubernetes Operator
Update `coherence-operator/values.yaml` if it is necessary.
    ```
    helm install coherence-operator
    ```
- Install Coherence
Update `coherence/values.yaml` to point to the desired Coherence Docker image.
    ```
    cd ../coherence-*-helm
    helm install coherence
    ```
    
## Accessing the Pre-Loaded Grafana Dashboards

Note: These notes are for the namespace `custom`, please change as required.
- Get the name of the Grafana pod and port forward.
```bash
export GRAFANA_POD=$(kubectl get pods --namespace custom | grep grafana | awk '{print $1}')
kubectl --namespace custom port-forward $GRAFANA_POD 3000:3000
```

- Get the Prometheus endpoint so we can add this to Grafana as a Datasource. (**This should be automated in the chart**)

```bash
export PROM_POD=$(kubectl get pod --namespace custom | grep prometheusope-prometheus | awk '{print $1}')

kubectl get pod --namespace custom $PROM_POD -o json | egrep 'hostname|subdomain'
```

Note the values of `hostname` and `subdomain`.

The following is only required once after helm install. (This should all be eventually automated)
- Login to Grafana at http://127.0.0.1:3000. Username is admin/prom-operator. Don't change the password if prompted.

- Click on `Add DataSource`.

- Name the datasource `Prometheus` and change the type the `Prometheus`.

- Make sure you mark this as `The Default`

- Enter the url using the `hostname` and `subdomain` from above.

```bash
http://hostname.subdomain:9090/
```

E.g. Example: `http://prometheus-hissing-seal-prometheusope-prometheus-0.prometheus-operated:9090/`

- Click on `Save and Test` button down the bottom and the message `Datasource is Working` should be displayed if you have correctly entered the details.

- Go to the following URL to access the main dashboard - `http://127.0.0.1:3000/d/coh-main/` 

## Known Issues

Port-forwarding seems very unreliable, at least on a VPN connection and seems to 
drop out every minute or so.  You may need to restart the port forward. 

This should be resolved when Grafana is accessed through a Load Balancer.

If you get the following error after deleting a helm chart adn then trying to install again: 
   ```bash
   Error: object is being deleted: customresourcedefinitions.apiextensions.k8s.io "prometheuses.monitoring.coreos.com" already exists
   ```
   
Use the script below to remove.  
   ```bash
   kubectl get customresourcedefinition | sed 1d | awk '{print $1}' | xargs kubectl delete customresourcedefinition
   ```
