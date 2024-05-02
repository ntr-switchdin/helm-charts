package redpanda

import (
	"fmt"

	"github.com/redpanda-data/helm-charts/pkg/gotohelm/helmette"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func PostUpgrade(dot *helmette.Dot) *batchv1.Job {
	values := helmette.Unwrap[Values](dot.Values)

	if !values.PostUpgradeJob.Enabled {
		return nil
	}

	labels := helmette.Default(map[string]string{}, values.PostUpgradeJob.Labels)
	annotations := helmette.Default(map[string]string{}, values.PostUpgradeJob.Annotations)

	labels = helmette.Merge(FullLabels(dot), labels)
	annotations = helmette.Merge(map[string]string{
		"helm.sh/hook":               "post-upgrade",
		"helm.sh/hook-delete-policy": "before-hook-creation",
		"helm.sh/hook-weight":        "-10",
	}, labels)

	securityContext := helmette.Merge(*ContainerSecurityContext(dot), ptr.Deref(values.PostUpgradeJob.SecurityContext, corev1.SecurityContext{}))

	return &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "batch/v1",
			Kind:       "Job",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf("%s-post-upgrade", Name(dot)),
			Namespace:   dot.Release.Namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: values.PostUpgradeJob.BackoffLimit,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: dot.Release.Name,
					Labels: helmette.Merge(map[string]string{
						"app.kubernetes.io/name":      Name(dot),
						"app.kubernetes.io/instance":  dot.Release.Name,
						"app.kubernetes.io/component": fmt.Sprintf("%s-post-upgrade", helmette.Trunc(50, Name(dot))),
					}, values.CommonLabels),
				},
				Spec: corev1.PodSpec{
					NodeSelector: values.NodeSelector,
					// Affinity:           ptr.Deref(&values.PostUpgradeJob.Affinity, &values.Affinity),
					Tolerations:        values.Tolerations,
					RestartPolicy:      corev1.RestartPolicyNever,
					SecurityContext:    PodSecurityContext(dot),
					ServiceAccountName: ServiceAccountName(dot),
					ImagePullSecrets:   values.ImagePullSecrets,
					Containers: []corev1.Container{
						{
							Name:    fmt.Sprintf("%s-post-upgrade", Name(dot)),
							Image:   fmt.Sprintf("%s:%s", values.Image.Repository, Tag(dot)),
							Command: []string{"/bin/bash", "-c"},
							Args:    []string{PostUpgradeJobScript(dot)},
							// TODO Env and EnvFrom could have previously been dynamic template references. Drop support for doing so??
							Env:             values.PostUpgradeJob.ExtraEnv,
							EnvFrom:         values.PostUpgradeJob.ExtraEnvFrom,
							SecurityContext: &securityContext,
							Resources:       values.PostUpgradeJob.Resources,
							VolumeMounts:    DefaultMounts(dot),
						},
					},
					Volumes: DefaultVolumes(dot),
				},
			},
		},
	}
}

func PostUpgradeJobScript(dot *helmette.Dot) string {
	values := helmette.Unwrap[Values](dot.Values)

	// {{- $service := .Values.listeners.admin -}}
	service := values.Listeners.Admin

	// {{- $cert := get .Values.tls.certs $service.tls.cert -}}
	cert := values.TLS.Certs[service.TLS.Cert]

	script := "set -e\n"
	for key, value := range values.Config.Cluster {
		if key == "default_topic_replications" {
			//     {{/* "sub (add $i (mod $i 2)) 1" calculates the closest odd number less than or equal to $i: 1=1, 2=1, 3=3, ... */}}
			//     {{- $r := $.Values.statefulset.replicas }}
			r := values.Statefulset.Replicas
			r = (r + (r % 2)) - 1
			//     {{- $value = min $value (sub (add $r (mod $r 2)) 1) }}
			value = helmette.Min(value.(int), r)
		}

		if asFloat, ok := value.(float64); ok {
			value = int64(asFloat)
		}

		//   {{- if and (typeIs "float64" $value) (eq (floor $value) $value) }}
		//     {{- $value = int64 $value }}
		//   {{- end }}
		//   {{- if and (typeIs "bool" $value ) ( not ( empty $value ) ) }}
		//             rpk cluster config set {{ $key }} {{ $value }}
		//   {{- else if and (typeIs "[]interface {}" $value ) ( not ( empty $value ) ) }}
		//             rpk cluster config set {{ $key }} "[ {{ join "," $value }} ]"
		//   {{- else if $value }}
		//             rpk cluster config set {{ $key }} {{ $value }}
		//   {{- end }}
		script = fmt.Sprintf("%srpk cluster config set %s %v\n", script, key, value)
	}

	if _, ok := values.Config.Cluster["storage_min_free_bytes"]; !ok {
		//             rpk cluster config set storage_min_free_bytes {{ include "storage-min-free-bytes" . }}
		script = fmt.Sprintf("%srpk cluster config set storage_min_free_bytes %d\n", script, int64(StorageMinFreeBytes(dot)))
	}

	// {{- if (include "redpanda-atleast-23-2-1" . | fromJson).bool }}
	if RedpandaAtLeast_23_2_1(dot) {
		caCert := ""
		//   {{- if and $cert ( dig "caEnabled" false $cert ) }}
		if cert.CAEnabled {
			caCert = fmt.Sprintf("--cacert /etc/tls/certs/%s/ca.crt", service.TLS.Cert)
		} //   {{- end }}

		//   {{- if (include "admin-internal-tls-enabled" . | fromJson).bool }}
		scheme := "http"
		if values.TLS.Enabled && service.TLS.Enabled && service.TLS.Cert != "" {
			scheme = "https"
		}

		url := fmt.Sprintf("%s://%s:%d/v1/debug/restart_service?service=schema-registry", scheme, InternalDomain(dot), int64(service.Port))

		script = fmt.Sprintf(`%s
  if [ -d "/etc/secrets/users/" ]; then
  IFS=":" read -r USER_NAME PASSWORD MECHANISM < <(grep "" $(find /etc/secrets/users/* -print))
  curl -svm3 --fail --retry "120" --retry-max-time "120" --retry-all-errors --ssl-reqd \
  -X PUT -u ${USER_NAME}:${PASSWORD} \
  %s \
  %s || true
fi
`, script, caCert, url)

	}

	return script
}
