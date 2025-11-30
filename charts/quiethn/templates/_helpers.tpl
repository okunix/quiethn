{{- define "quiethn.configMapName" -}}
{{- include "quiethn.fullname" . -}}
{{- end -}}

{{- define "quiethn.secretName" -}}
{{- include "quiethn.fullname" . -}}
{{- end -}}

{{/*
Full environment variable config for quiethn
*/}}
{{- define "quiethn.env" -}}
- name: HN_REDIS_ADDR
  valueFrom:
    configMapKeyRef:
      key: redisAddr
      name: {{ include "quiethn.configMapName" . }}
- name: HN_REDIS_DB
  valueFrom:
    configMapKeyRef:
      key: redisDb
      name: {{ include "quiethn.configMapName" . }}
      optional: true
- name: HN_REDIS_PASSWORD
  valueFrom:
    secretKeyRef:
      key: redisPassword 
      name: {{ include "quiethn.secretName" . }}
- name: HN_BASE_URL
  valueFrom:
    configMapKeyRef:
      key: hnAddr 
      name: {{ include "quiethn.configMapName" . }}
      optional: true
- name: HN_SERVER_PORT
  valueFrom:
    configMapKeyRef:
      key: serverPort
      name: {{ include "quiethn.configMapName" . }}
- name: HN_SERVER_HOST
  valueFrom:
    configMapKeyRef:
      key: serverHost
      name: {{ include "quiethn.configMapName" . }}
{{- end -}}


{{/*
Expand the name of the chart.
*/}}
{{- define "quiethn.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "quiethn.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "quiethn.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "quiethn.labels" -}}
helm.sh/chart: {{ include "quiethn.chart" . }}
{{ include "quiethn.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "quiethn.selectorLabels" -}}
app.kubernetes.io/name: {{ include "quiethn.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "quiethn.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "quiethn.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}
