# Updating Grafana Dashboards

**Ensure you read and follow the following when updating Grafana dashboards**

There are three flavours of Grafana Dashboards available in the sub directories here:

1. `grafana` - dashboards that are compatible with the metric names produced by the Coherence metrics publisher
2. `grafana-micrometer` - dashboards that are compatible with the metric names produced by the Coherence Micrometer Metrics module
3.  `grafana-microprofile` - dashboards that are compatible with the metric names produced by the Coherence MP Metrics module

These are described in the Coherence Operator documentation: [https://oracle.github.io/coherence-operator/docs/latest/#/metrics/030_importing](https://oracle.github.io/coherence-operator/docs/latest/#/metrics/030_importing).

If you need to update any Grafana dashboards, please use the following guidelines.

1. Any updates to the default dashboards in the `grafana` directory can be synced to the `grafana-micrometer` directory by running the script `sync-dashboards.sh`.

2. Updates to the `grafana-microprofile` dashboards need to be done manually as the metric
   names are significantly different and there is additional information displayed on the Microprofile dashboards. This means they cannot, and will not be updated using the script.

See the comments in the `sync-dashboards.sh` for more details.
