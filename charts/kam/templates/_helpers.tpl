{{/*
Expand the name of the chart.
*/}}
{{- define "charts.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "charts.fullname" -}}
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

{{- define "charts.fullnameBackend" -}}
  {{- $nameLength := int (sub 63 8) -}}
  {{- if .Values.fullnameOverride }}
    {{- printf "%s-backend" (.Values.fullnameOverride | trunc $nameLength | trimSuffix "-") -}}
  {{- else }}
    {{- $name := default .Chart.Name .Values.nameOverride -}}
    {{- if contains $name .Release.Name }}
      {{- printf "%s-backend" (.Release.Name | trunc $nameLength | trimSuffix "-") -}}
    {{- else }}
      {{- printf "%s-%s-backend" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
    {{- end }}
  {{- end }}
{{- end }}


{{- define "charts.fullnameFrontend" -}}
  {{- $nameLength := int (sub 63 8) -}}
  {{- if .Values.fullnameOverride }}
    {{- printf "%s-frontend" (.Values.fullnameOverride | trunc $nameLength | trimSuffix "-") -}}
  {{- else }}
    {{- $name := default .Chart.Name .Values.nameOverride -}}
    {{- if contains $name .Release.Name }}
      {{- printf "%s-frontend" (.Release.Name | trunc $nameLength | trimSuffix "-") -}}
    {{- else }}
      {{- printf "%s-%s-frontend" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
    {{- end }}
  {{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "charts.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "charts.labels" -}}
helm.sh/chart: {{ include "charts.chart" . }}
{{ include "charts.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{- define "charts.labelsBackend" -}}
helm.sh/chart: {{ include "charts.chart" . }}
{{ include "charts.selectorLabelsBackend" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{- define "charts.labelsFrontend" -}}
helm.sh/chart: {{ include "charts.chart" . }}
{{ include "charts.selectorLabelsFrontend" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "charts.selectorLabels" -}}
app.kubernetes.io/name: {{ include "charts.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "charts.selectorLabelsBackend" -}}
app.kubernetes.io/name: {{ include "charts.name" . }}-backend
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "charts.selectorLabelsFrontend" -}}
app.kubernetes.io/name: {{ include "charts.name" . }}-frontend
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "charts.serviceAccountName" -}}
{{- if .Values.backend.serviceAccount.create }}
{{- default (include "charts.fullname" .) .Values.backend.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.backend.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Define backend enviromental values
*/}}
{{- define "charts.backendEnvVars" }}
{{- if .Values.global.env.KEYCLOAK_JWKS_URL }}
- name: KEYCLOAK_JWKS_URL
  value: "{{ .Values.global.env.KEYCLOAK_JWKS_URL }}"
{{- end }}
{{- if .Values.global.env.VITE_KEYCLOAK_URL }}
- name: VITE_KEYCLOAK_URL
  value: "{{ .Values.global.env.VITE_KEYCLOAK_URL }}"
{{- end }}
{{- if .Values.global.env.VITE_KEYCLOAK_REALMNAME }}
- name: VITE_KEYCLOAK_REALMNAME
  value: "{{ .Values.global.env.VITE_KEYCLOAK_REALMNAME }}"
{{- end }}
{{- if .Values.global.env.VITE_KEYCLOAK_CLIENTNAME }}
- name: VITE_KEYCLOAK_CLIENTNAME
  value: "{{ .Values.global.env.VITE_KEYCLOAK_CLIENTNAME }}"
{{- end }}
{{- if .Values.backend.healthPort }}
- name: HEALTH_PORT
  value: "{{ .Values.backend.healthPort }}"
{{- end }}
{{- if .Values.backend.service.port }}
- name: BACKEND_PORT
  value: "{{ .Values.backend.service.port }}"
{{- end }}
{{- if .Values.global.env.ROLEMAP_NAMESPACE }}
- name: ROLEMAP_NAMESPACE
  value: "{{ .Values.global.env.ROLEMAP_NAMESPACE }}"
{{- end }}
{{- if .Values.global.env.ROLEMAP_NAME }}
- name: ROLEMAP_NAME
  value: "{{ .Values.global.env.ROLEMAP_NAME }}"
{{- end }}
- name: IN_CLUSTER_MODE
  value: "true"
{{- end }}

{{/*
Define frontend environmental values
*/}}
{{- define "charts.frontendEnvVars" }}
- name: CUM_API_URL
  value: "{{ include "charts.fullnameBackend" . }}:{{ .Values.backend.service.port }}"
{{- if .Values.global.env.VITE_KEYCLOAK_URL }}
- name: KAM_KEYCLOAK_URL
  value: "{{ .Values.global.env.VITE_KEYCLOAK_URL }}"
{{- end }}
{{- if .Values.global.env.VITE_KEYCLOAK_CLIENTNAME }}
- name: KAM_KEYCLOAK_CLIENTNAME
  value: "{{ .Values.global.env.VITE_KEYCLOAK_CLIENTNAME }}"
{{- end }}
{{- if .Values.global.env.VITE_KEYCLOAK_REALMNAME }}
- name: KAM_KEYCLOAK_REALMNAME
  value: "{{ .Values.global.env.VITE_KEYCLOAK_REALMNAME }}"
{{- end }}
{{- if .Values.global.env.KAM_KEYCLOAK_LOGIN_URL }}
- name: KAM_KEYCLOAK_LOGIN_URL
  value: "{{ .Values.global.env.KAM_KEYCLOAK_LOGIN_URL }}"
{{- end }}
{{- if .Values.global.env.KAM_KEYCLOAK_LOGOUT_URL }}
- name: KAM_KEYCLOAK_LOGOUT_URL
  value: "{{ .Values.global.env.KAM_KEYCLOAK_LOGOUT_URL }}"
{{- end }}
{{- if .Values.global.env.KAM_KEYCLOAK_TOKEN_URL }}
- name: KAM_KEYCLOAK_TOKEN_URL
  value: "{{ .Values.global.env.KAM_KEYCLOAK_TOKEN_URL }}"
{{- end }}
{{- end }}
