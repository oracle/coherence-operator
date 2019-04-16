# Accessing the EFK stack for viewing logs

The Oracle Coherence Kubernetes Operator manages Oracle Coherence on Kubernetes.
It manages monitoring data through Prometheus, and logging data through the EFK
(ElasticSearch, Fluentd and Kibana) stack.

To access and configure the Kibana UI, please follow the instructions below.

## 1. Installing the Charts

When you install the `coherence-operator` and `coherence` charts, you must specify the following
option for `helm` , for both charts, to ensure the EFK stack (Elasitcsearch, Fluentd and Kibana) 
is installed and correctly configured.

```bash
--set logCaptureEnabled=true 
```

## 2. Port Forward Kibana

Once you have installed both charts, use the following script to port forward the Kibana port 5601.

Note: If your chart is installed in a namespace other than `default`
then include `--namespace` option for both `kubectl` commands.

```bash
#!/bin/bash
  
export KIBANA_POD=$(kubectl get pods | grep kibana | awk '{print $1}')
echo $KIBANA_POD

while :
do
 kubectl port-forward $KIBANA_POD 5601:5601
done

```

## 3. Default Dashboards

There are a number of dashboard created via the import process.

**Coherence Operator**

* *Coherence Operator - All Messages* - Shows all Coherence Operator messages

**Coherence Cluster**

* *Coherence Cluster - All Messages* - Shows all messages

* *Coherence Cluster - Errors and Warnings* - Shows only errors and warnings

* *Coherence Cluster - Persistence* - Shows persistence related messages 

* *Coherence Cluster - Partitions* - Shows Ppartition related messages 

* *Coherence Cluster - Message Sources* - Allows visualization of messages via the message source (Thread)

* *Coherence Cluster - Configuration Messages* - Shows configuration related messages

* *Coherence Cluster - Network* - Shows network related messages such as communication delays and TCP ring disconnects 

## 4. Default Queries

There are many queries related to common Coherence messages, warnings and errors that are 
loaded and cam be accessed via the discover side-bar.

## 5. WIP: Application Log Capture


### Application sidecar 

#### Resources

##### 1. custom-log.properties

Make a copy of default coherence logging properties and 
add the following for application level logging.

```
# sample extension to logging
sample.handlers=java.util.logging.ConsoleHandler,sample.CustomFileHandler
sample.level=INFO

sample.CustomFileHandler.pattern=/logs/sample-%g.log
sample.CustomFileHandler.limit=10485760
sample.CustomFileHandler.count=10
sample.CustomFileHandler.formatter=java.util.logging.SimpleFormatter
```

Java logging to logger "sample*" will go to stdout and to log file /logs/sample-%g.log.

##### 2. fluentd conf fluentd-sample.conf

Create fluentd events tagged `sample` from `sample-*.log` source.

```
# Application Logs
    <source>
      @type tail
      path /logs/sample-*.log
      pos_file /tmp/sample.log.pos
      read_from_head true
      tag sample
      <parse>
        @type regexp
        expression /^(?<logtime>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}) (?<product>.+) <(?<level>[^\s]+)> \(thread=(?<thread>.+), member=(?<member>.+)\): (?<log>.*)$/
        time_key logtime
      </parse>
    </source>
```

#### Java Code 

##### 1. Java log code

Added CustomFileHandler to enable logging application log messages to different log 
from coherence logs. CustomFileHandler just extends java.util.logging.FileHandler.

```
Logger logger = new java.util.logging.Logger.getLogger("sample");

logger.info("....");
```

### Install Coherence Operator Helm Chart 

##### 1. set values
```
--set logCaptureEnabled=true
```

### Install Coherence Helm Chart with sidecar image

##### 1. set values to configure fluentd for application log events
```
--set logCaptureEnabled=true,store.logging.configFile=custom-logging.properties
--set fluentd.application.configFile=/conf/fluentd-sample.conf,fluentd.application.tag=sample
```

#### 2. Port forward extend port


### Run extend client 


### View application log messages from Kibana

##### 1. Port forward Kibana as described above.

##### 2. create index pattern for sample-*


## 6. Troubleshooting

### No Default Index Pattern

There are two index patterns created via the import process and the `coherence-cluster-*` pattern
will be set as the default. If for some reason this is not the case, then carry out the following
to set `coherence-cluster-*` as the default:

* Open a the following URL `http://127.0.0.1:5601/`

* Click on `Management` side-bar.

* Click on `Index Patterns`, then select `coherence-cluster-*' pattern.

* Click the `Star` on the top right to set as default.

