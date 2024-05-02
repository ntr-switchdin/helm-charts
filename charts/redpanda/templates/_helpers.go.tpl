{{- /* Generated from "helpers.go" */ -}}

{{- define "redpanda.Chart" -}}
{{- $dot := (index .a 0) -}}
{{- range $_ := (list 1) -}}
{{- (dict "r" (get (fromJson (include "redpanda.cleanForK8s" (dict "a" (list (replace "+" "_" (printf "%s-%s" $dot.Chart.Name $dot.Chart.Version))) ))) "r")) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.Name" -}}
{{- $dot := (index .a 0) -}}
{{- range $_ := (list 1) -}}
{{- $tmp_tuple_1 := (get (fromJson (include "_shims.compact" (dict "a" (list (get (fromJson (include "_shims.typetest" (dict "a" (list "string" (index $dot.Values "nameOverride") "") ))) "r")) ))) "r") -}}
{{- $ok_2 := $tmp_tuple_1.T2 -}}
{{- $override_1 := $tmp_tuple_1.T1 -}}
{{- if (and $ok_2 (ne $override_1 "")) -}}
{{- (dict "r" (get (fromJson (include "redpanda.cleanForK8s" (dict "a" (list $override_1) ))) "r")) | toJson -}}
{{- break -}}
{{- end -}}
{{- (dict "r" (get (fromJson (include "redpanda.cleanForK8s" (dict "a" (list $dot.Chart.Name) ))) "r")) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.Fullname" -}}
{{- $dot := (index .a 0) -}}
{{- range $_ := (list 1) -}}
{{- $tmp_tuple_2 := (get (fromJson (include "_shims.compact" (dict "a" (list (get (fromJson (include "_shims.typetest" (dict "a" (list "string" (index $dot.Values "fullnameOverride") "") ))) "r")) ))) "r") -}}
{{- $ok_4 := $tmp_tuple_2.T2 -}}
{{- $override_3 := $tmp_tuple_2.T1 -}}
{{- if (and $ok_4 (ne $override_3 "")) -}}
{{- (dict "r" (get (fromJson (include "redpanda.cleanForK8s" (dict "a" (list $override_3) ))) "r")) | toJson -}}
{{- break -}}
{{- end -}}
{{- (dict "r" (get (fromJson (include "redpanda.cleanForK8s" (dict "a" (list (printf "%s" $dot.Release.Name)) ))) "r")) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.FullLabels" -}}
{{- $dot := (index .a 0) -}}
{{- range $_ := (list 1) -}}
{{- $values := $dot.Values.AsMap -}}
{{- $labels := (dict ) -}}
{{- if (ne $values.commonLabels (coalesce nil)) -}}
{{- $labels = $values.commonLabels -}}
{{- end -}}
{{- $defaults := (dict "helm.sh/chart" (get (fromJson (include "redpanda.Chart" (dict "a" (list $dot) ))) "r") "app.kubernetes.io/name" (get (fromJson (include "redpanda.Name" (dict "a" (list $dot) ))) "r") "app.kubernetes.io/instance" $dot.Release.Name "app.kubernetes.io/managed-by" $dot.Release.Service "app.kubernetes.io/component" (get (fromJson (include "redpanda.Name" (dict "a" (list $dot) ))) "r") ) -}}
{{- (dict "r" (merge (dict ) $labels $defaults)) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.StatefulSetPodLabelsSelector" -}}
{{- $dot := (index .a 0) -}}
{{- $existing := (index .a 1) -}}
{{- range $_ := (list 1) -}}
{{- if (and $dot.Release.IsUpgrade (ne $existing (coalesce nil))) -}}
{{- if (gt (int (get (fromJson (include "_shims.len" (dict "a" (list $existing.spec.selector.matchLabels) ))) "r")) 0) -}}
{{- (dict "r" $existing.spec.selector.matchLabels) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}
{{- $values := $dot.Values.AsMap -}}
{{- $additionalSelectorLabels := (dict ) -}}
{{- if (ne $values.statefulset.additionalSelectorLabels (coalesce nil)) -}}
{{- $additionalSelectorLabels = $values.statefulset.additionalSelectorLabels -}}
{{- end -}}
{{- $component := (printf "%s-statefulset" (trimSuffix "-" (trunc 51 (get (fromJson (include "redpanda.Name" (dict "a" (list $dot) ))) "r")))) -}}
{{- $defaults := (dict "app.kubernetes.io/component" $component "app.kubernetes.io/instance" $dot.Release.Name "app.kubernetes.io/name" (get (fromJson (include "redpanda.Name" (dict "a" (list $dot) ))) "r") ) -}}
{{- (dict "r" (merge (dict ) $additionalSelectorLabels $defaults)) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.StatefulSetPodLabels" -}}
{{- $dot := (index .a 0) -}}
{{- $existing := (index .a 1) -}}
{{- range $_ := (list 1) -}}
{{- if (and $dot.Release.IsUpgrade (ne $existing (coalesce nil))) -}}
{{- if (gt (int (get (fromJson (include "_shims.len" (dict "a" (list $existing.spec.template.metadata.labels) ))) "r")) 0) -}}
{{- (dict "r" $existing.spec.template.metadata.labels) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}
{{- $values := $dot.Values.AsMap -}}
{{- $statefulSetLabels := (dict ) -}}
{{- if (ne $values.statefulset.podTemplate.labels (coalesce nil)) -}}
{{- $statefulSetLabels = $values.statefulset.podTemplate.labels -}}
{{- end -}}
{{- $defaults := (dict "redpanda.com/poddisruptionbudget" (get (fromJson (include "redpanda.Fullname" (dict "a" (list $dot) ))) "r") ) -}}
{{- (dict "r" (merge (dict ) $statefulSetLabels (get (fromJson (include "redpanda.StatefulSetPodLabelsSelector" (dict "a" (list $dot $existing) ))) "r") $defaults)) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.StatefulSetPodAnnotations" -}}
{{- $dot := (index .a 0) -}}
{{- $configMapChecksum := (index .a 1) -}}
{{- range $_ := (list 1) -}}
{{- $values := $dot.Values.AsMap -}}
{{- $configMapChecksumAnnotation := (dict "config.redpanda.com/checksum" $configMapChecksum ) -}}
{{- if (ne $values.statefulset.podTemplate.annotations (coalesce nil)) -}}
{{- (dict "r" (merge (dict ) $values.statefulset.podTemplate.annotations $configMapChecksumAnnotation)) | toJson -}}
{{- break -}}
{{- end -}}
{{- (dict "r" (merge (dict ) $values.statefulset.annotations $configMapChecksumAnnotation)) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.ServiceAccountName" -}}
{{- $dot := (index .a 0) -}}
{{- range $_ := (list 1) -}}
{{- $values := $dot.Values.AsMap -}}
{{- $serviceAccount := $values.serviceAccount -}}
{{- if (and $serviceAccount.create (ne $serviceAccount.name "")) -}}
{{- (dict "r" $serviceAccount.name) | toJson -}}
{{- break -}}
{{- else -}}{{- if $serviceAccount.create -}}
{{- (dict "r" (get (fromJson (include "redpanda.Fullname" (dict "a" (list $dot) ))) "r")) | toJson -}}
{{- break -}}
{{- else -}}{{- if (ne $serviceAccount.name "") -}}
{{- (dict "r" $serviceAccount.name) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}
{{- end -}}
{{- (dict "r" "default") | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.Tag" -}}
{{- $dot := (index .a 0) -}}
{{- range $_ := (list 1) -}}
{{- $values := $dot.Values.AsMap -}}
{{- $tag := (toString $values.image.tag) -}}
{{- if (eq $tag "") -}}
{{- $tag = $dot.Chart.AppVersion -}}
{{- end -}}
{{- $pattern := "^v(0|[1-9]\\d*)\\.(0|[1-9]\\d*)\\.(0|[1-9]\\d*)(?:-((?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+([0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?$" -}}
{{- if (not (regexMatch $pattern $tag)) -}}
{{- $_ := (fail "image.tag must start with a 'v' and be a valid semver") -}}
{{- end -}}
{{- (dict "r" $tag) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.ServiceName" -}}
{{- $dot := (index .a 0) -}}
{{- range $_ := (list 1) -}}
{{- $values := $dot.Values.AsMap -}}
{{- if (and (ne $values.service (coalesce nil)) (ne $values.service.name (coalesce nil))) -}}
{{- (dict "r" (get (fromJson (include "redpanda.cleanForK8s" (dict "a" (list $values.service.name) ))) "r")) | toJson -}}
{{- break -}}
{{- end -}}
{{- (dict "r" (get (fromJson (include "redpanda.Fullname" (dict "a" (list $dot) ))) "r")) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.InternalDomain" -}}
{{- $dot := (index .a 0) -}}
{{- range $_ := (list 1) -}}
{{- $values := $dot.Values.AsMap -}}
{{- $service := (get (fromJson (include "redpanda.ServiceName" (dict "a" (list $dot) ))) "r") -}}
{{- $ns := $dot.Release.Namespace -}}
{{- $domain := (trimSuffix "." $values.clusterDomain) -}}
{{- (dict "r" (printf "%s.%s.svc.%s." $service $ns $domain)) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.TLSEnabled" -}}
{{- $dot := (index .a 0) -}}
{{- range $_ := (list 1) -}}
{{- $values := $dot.Values.AsMap -}}
{{- if $values.tls.enabled -}}
{{- (dict "r" true) | toJson -}}
{{- break -}}
{{- end -}}
{{- $listeners := (list "kafka" "admin" "schemaRegistry" "rpc" "http") -}}
{{- range $_, $listener := $listeners -}}
{{- $tlsCert := (dig "listeners" $listener "tls" "cert" false $dot.Values.AsMap) -}}
{{- $tlsEnabled := (dig "listeners" $listener "tls" "enabled" false $dot.Values.AsMap) -}}
{{- if (and (not (empty $tlsEnabled)) (not (empty $tlsCert))) -}}
{{- (dict "r" true) | toJson -}}
{{- break -}}
{{- end -}}
{{- $external := (dig "listeners" $listener "external" false $dot.Values.AsMap) -}}
{{- if (empty $external) -}}
{{- continue -}}
{{- end -}}
{{- $keys := (keys (get (fromJson (include "_shims.typeassertion" (dict "a" (list (printf "map[%s]%s" "string" "interface {}") $external) ))) "r")) -}}
{{- range $_, $key := $keys -}}
{{- $enabled := (dig "listeners" $listener "external" $key "enabled" false $dot.Values.AsMap) -}}
{{- $tlsCert := (dig "listeners" $listener "external" $key "tls" "cert" false $dot.Values.AsMap) -}}
{{- $tlsEnabled := (dig "listeners" $listener "external" $key "tls" "enabled" false $dot.Values.AsMap) -}}
{{- if (and (and (not (empty $enabled)) (not (empty $tlsCert))) (not (empty $tlsEnabled))) -}}
{{- (dict "r" true) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}
{{- end -}}
{{- (dict "r" false) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.ClientAuthRequired" -}}
{{- $dot := (index .a 0) -}}
{{- range $_ := (list 1) -}}
{{- $listeners := (list "kafka" "admin" "schemaRegistry" "rpc" "http") -}}
{{- range $_, $listener := $listeners -}}
{{- $required := (dig $listener "tls" "requireClientAuth" false $dot.Values.AsMap) -}}
{{- if (not (empty $required)) -}}
{{- (dict "r" true) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}
{{- (dict "r" false) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.DefaultMounts" -}}
{{- $dot := (index .a 0) -}}
{{- range $_ := (list 1) -}}
{{- (dict "r" (concat (list (mustMergeOverwrite (dict "name" "" "mountPath" "" ) (dict "name" "config" "mountPath" "/etc/redpanda" ))) (get (fromJson (include "redpanda.CommonMounts" (dict "a" (list $dot) ))) "r"))) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.CommonMounts" -}}
{{- $dot := (index .a 0) -}}
{{- range $_ := (list 1) -}}
{{- $values := $dot.Values.AsMap -}}
{{- $mounts := (list ) -}}
{{- $sasl_5 := $values.auth.sasl -}}
{{- if (and $sasl_5.enabled (ne $sasl_5.secretRef "")) -}}
{{- $mounts = (mustAppend $mounts (mustMergeOverwrite (dict "name" "" "mountPath" "" ) (dict "name" "users" "mountPath" "/etc/secrets/users" "readOnly" true ))) -}}
{{- end -}}
{{- if (get (fromJson (include "redpanda.TLSEnabled" (dict "a" (list $dot) ))) "r") -}}
{{- $certNames := (keys $values.tls.certs) -}}
{{- $_ := (sortAlpha $certNames) -}}
{{- range $_, $name := $certNames -}}
{{- $mounts = (mustAppend $mounts (mustMergeOverwrite (dict "name" "" "mountPath" "" ) (dict "name" (printf "redpanda-%s-cert" $name) "mountPath" (printf "/etc/tls/certs/%s" $name) ))) -}}
{{- end -}}
{{- if (get (fromJson (include "redpanda.ClientAuthRequired" (dict "a" (list $dot) ))) "r") -}}
{{- $mounts = (mustAppend $mounts (mustMergeOverwrite (dict "name" "" "mountPath" "" ) (dict "name" "mtls-client" "mountPath" (printf "/etc/tls/certs/%s-client" (get (fromJson (include "redpanda.Fullname" (dict "a" (list $dot) ))) "r")) ))) -}}
{{- end -}}
{{- end -}}
{{- (dict "r" $mounts) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.DefaultVolumes" -}}
{{- $dot := (index .a 0) -}}
{{- range $_ := (list 1) -}}
{{- (dict "r" (concat (list (mustMergeOverwrite (mustMergeOverwrite (dict ) (dict "name" "" )) (mustMergeOverwrite (dict ) (dict "configMap" (mustMergeOverwrite (mustMergeOverwrite (dict ) (dict )) (mustMergeOverwrite (dict ) (dict "name" (get (fromJson (include "redpanda.Fullname" (dict "a" (list $dot) ))) "r") )) (dict )) )) (dict "name" "config" ))) (get (fromJson (include "redpanda.CommonVolumes" (dict "a" (list $dot) ))) "r"))) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.CommonVolumes" -}}
{{- $dot := (index .a 0) -}}
{{- range $_ := (list 1) -}}
{{- $volumes := (list ) -}}
{{- $values := $dot.Values.AsMap -}}
{{- if (get (fromJson (include "redpanda.TLSEnabled" (dict "a" (list $dot) ))) "r") -}}
{{- $certNames := (keys $values.tls.certs) -}}
{{- $_ := (sortAlpha $certNames) -}}
{{- $mode := (int 0o440) -}}
{{- range $_, $name := $certNames -}}
{{- $cert := (index $values.tls.certs $name) -}}
{{- $volumes = (mustAppend $volumes (mustMergeOverwrite (mustMergeOverwrite (dict ) (dict "name" "" )) (mustMergeOverwrite (dict ) (dict "secret" (mustMergeOverwrite (dict ) (dict "secretName" (get (fromJson (include "redpanda.CertSecretName" (dict "a" (list $dot $name $cert) ))) "r") "defaultMode" $mode )) )) (dict "name" (printf "redpanda-%s-cert" $name) ))) -}}
{{- end -}}
{{- if (get (fromJson (include "redpanda.ClientAuthRequired" (dict "a" (list $dot) ))) "r") -}}
{{- $volumes = (mustAppend $volumes (mustMergeOverwrite (mustMergeOverwrite (dict ) (dict "name" "" )) (mustMergeOverwrite (dict ) (dict "secret" (mustMergeOverwrite (dict ) (dict "secretName" (printf "%s-client" (get (fromJson (include "redpanda.Fullname" (dict "a" (list $dot) ))) "r")) "defaultMode" $mode )) )) (dict "name" "mtls-client" ))) -}}
{{- end -}}
{{- end -}}
{{- $sasl_6 := $values.auth.sasl -}}
{{- if (and $sasl_6.enabled (ne $sasl_6.secretRef "")) -}}
{{- $volumes = (mustAppend $volumes (mustMergeOverwrite (mustMergeOverwrite (dict ) (dict "name" "" )) (mustMergeOverwrite (dict ) (dict "secret" (mustMergeOverwrite (dict ) (dict "secretName" $sasl_6.secretRef )) )) (dict "name" "users" ))) -}}
{{- end -}}
{{- (dict "r" $volumes) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.CertSecretName" -}}
{{- $dot := (index .a 0) -}}
{{- $certName := (index .a 1) -}}
{{- $cert := (index .a 2) -}}
{{- range $_ := (list 1) -}}
{{- if (ne $cert.secretRef (coalesce nil)) -}}
{{- (dict "r" $cert.secretRef.name) | toJson -}}
{{- break -}}
{{- end -}}
{{- (dict "r" (printf "%s-%s-cert" (get (fromJson (include "redpanda.Fullname" (dict "a" (list $dot) ))) "r") $certName)) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.PodSecurityContext" -}}
{{- $dot := (index .a 0) -}}
{{- range $_ := (list 1) -}}
{{- $values := $dot.Values.AsMap -}}
{{- $sc := (get (fromJson (include "_shims.ptr_Deref" (dict "a" (list $values.statefulset.podSecurityContext $values.statefulset.securityContext) ))) "r") -}}
{{- (dict "r" (mustMergeOverwrite (dict ) (dict "fsGroup" $sc.fsGroup "fsGroupChangePolicy" $sc.fsGroupChangePolicy ))) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.ContainerSecurityContext" -}}
{{- $dot := (index .a 0) -}}
{{- range $_ := (list 1) -}}
{{- $values := $dot.Values.AsMap -}}
{{- $sc := (get (fromJson (include "_shims.ptr_Deref" (dict "a" (list $values.statefulset.podSecurityContext $values.statefulset.securityContext) ))) "r") -}}
{{- (dict "r" (mustMergeOverwrite (dict ) (dict "runAsUser" $sc.runAsUser "runAsGroup" (coalesce $sc.runAsGroup $sc.fsGroup) "allowPrivilegeEscalation" $sc.allowPriviledgeEscalation "runAsNonRoot" $sc.runAsNonRoot ))) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.RedpandaAtLeast_22_2_0" -}}
{{- $dot := (index .a 0) -}}
{{- range $_ := (list 1) -}}
{{- (dict "r" (get (fromJson (include "redpanda.redpandaAtLeast" (dict "a" (list $dot ">=22.2.0-0 || <0.0.1-0") ))) "r")) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.RedpandaAtLeast_22_3_0" -}}
{{- $dot := (index .a 0) -}}
{{- range $_ := (list 1) -}}
{{- (dict "r" (get (fromJson (include "redpanda.redpandaAtLeast" (dict "a" (list $dot ">=23.3.0-0 || <0.0.1-0") ))) "r")) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.RedpandaAtLeast_23_1_1" -}}
{{- $dot := (index .a 0) -}}
{{- range $_ := (list 1) -}}
{{- (dict "r" (get (fromJson (include "redpanda.redpandaAtLeast" (dict "a" (list $dot ">=23.1.1-0 || <0.0.1-0") ))) "r")) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.RedpandaAtLeast_23_1_2" -}}
{{- $dot := (index .a 0) -}}
{{- range $_ := (list 1) -}}
{{- (dict "r" (get (fromJson (include "redpanda.redpandaAtLeast" (dict "a" (list $dot ">=23.1.2-0 || <0.0.1-0") ))) "r")) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.RedpandaAtLeast_22_3_atleast_22_3_13" -}}
{{- $dot := (index .a 0) -}}
{{- range $_ := (list 1) -}}
{{- (dict "r" (get (fromJson (include "redpanda.redpandaAtLeast" (dict "a" (list $dot ">=22.3.13-0,<22.4") ))) "r")) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.RedpandaAtLeast_22_2_atleast_22_2_10" -}}
{{- $dot := (index .a 0) -}}
{{- range $_ := (list 1) -}}
{{- (dict "r" (get (fromJson (include "redpanda.redpandaAtLeast" (dict "a" (list $dot ">=22.2.10-0,<22.3") ))) "r")) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.RedpandaAtLeast_23_2_1" -}}
{{- $dot := (index .a 0) -}}
{{- range $_ := (list 1) -}}
{{- (dict "r" (get (fromJson (include "redpanda.redpandaAtLeast" (dict "a" (list $dot ">=23.2.1-0 || <0.0.1-0") ))) "r")) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.RedpandaAtLeast_23_3_0" -}}
{{- $dot := (index .a 0) -}}
{{- range $_ := (list 1) -}}
{{- (dict "r" (get (fromJson (include "redpanda.redpandaAtLeast" (dict "a" (list $dot ">=23.3.0-0 || <0.0.1-0") ))) "r")) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.redpandaAtLeast" -}}
{{- $dot := (index .a 0) -}}
{{- $constraint := (index .a 1) -}}
{{- range $_ := (list 1) -}}
{{- $values := $dot.Values.AsMap -}}
{{- if (ne $values.image.repository "docker.redpanda.com/redpandadata/redpanda") -}}
{{- (dict "r" true) | toJson -}}
{{- break -}}
{{- end -}}
{{- $version := (trimPrefix "v" (get (fromJson (include "redpanda.Tag" (dict "a" (list $dot) ))) "r")) -}}
{{- (dict "r" (semverCompare $constraint $version)) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.StorageMinFreeBytes" -}}
{{- $dot := (index .a 0) -}}
{{- range $_ := (list 1) -}}
{{- $values := $dot.Values.AsMap -}}
{{- $fiveGiB := (int64 5368709120) -}}
{{- if (not $values.storage.persistentVolume.enabled) -}}
{{- (dict "r" $fiveGiB) | toJson -}}
{{- break -}}
{{- end -}}
{{- $minSize := (int64 (float64 (mulf (float64 $values.storage.persistentVolume.size.IntValue) 0.05))) -}}
{{- (dict "r" (min $minSize $fiveGiB)) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

{{- define "redpanda.cleanForK8s" -}}
{{- $in := (index .a 0) -}}
{{- range $_ := (list 1) -}}
{{- (dict "r" (trimSuffix "-" (trunc 63 $in))) | toJson -}}
{{- break -}}
{{- end -}}
{{- end -}}

