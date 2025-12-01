{{/* Shortened name suffixed with upgrade-crd */}}
{{- define "resource-state-metrics.crd.upgradeJob.name" -}}
{{- print (include "resource-state-metrics.fullname" .) "-upgrade" -}}
{{- end -}}

{{- define "resource-state-metrics.crd.upgradeJob.labels" -}}
{{- include "resource-state-metrics.labels" . }}
app: {{ template "resource-state-metrics.name" . }}
app.kubernetes.io/name: {{ template "resource-state-metrics.name" . }}
app.kubernetes.io/component: crds-upgrade
{{- end -}}

{{/* Create the name of crd.upgradeJob service account to use */}}
{{- define "resource-state-metrics.crd.upgradeJob.serviceAccountName" -}}
{{- if .Values.upgradeJob.serviceAccount.create -}}
    {{ default (include "resource-state-metrics.crd.upgradeJob.name" .) .Values.upgradeJob.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.upgradeJob.serviceAccount.name }}
{{- end -}}
{{- end -}}
