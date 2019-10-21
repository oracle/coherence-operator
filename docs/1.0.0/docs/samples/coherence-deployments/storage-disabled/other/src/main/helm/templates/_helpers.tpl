{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "coherence.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "coherence.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "coherence.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Define the Coherence role name for this install.
*/}}
{{- define "coherence.role" -}}
{{ .Values.role | default "WebServer" }}
{{- end -}}

{{/*
Create the release labels.
These are a common set of labels applied to all of the resources
generated from this chart.
*/}}
{{- define "coherence.release_labels" }}
coherenceDeployment: {{ template "coherence.fullname" . }}
coherenceCluster: {{ .Values.clusterName }}
heritage: {{ .Release.Service }}
release: {{ .Release.Name }}
chart: {{ template "coherence.chart" . }}
app: {{ template "coherence.name" . }}
{{- end }}
