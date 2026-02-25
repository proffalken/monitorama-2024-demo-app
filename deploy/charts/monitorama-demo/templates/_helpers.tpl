{{- define "monitorama-demo.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "monitorama-demo.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name (include "monitorama-demo.name" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "monitorama-demo.labels" -}}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "monitorama-demo.cnpgClusterName" -}}
{{- printf "%s-cnpg" (include "monitorama-demo.fullname" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "monitorama-demo.cnpgAppSecretName" -}}
{{- if .Values.cnpg.database.existingSecret -}}
{{- .Values.cnpg.database.existingSecret -}}
{{- else -}}
{{- printf "%s-app" (include "monitorama-demo.cnpgClusterName" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
