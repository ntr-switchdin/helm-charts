//go:build !generate

// +gotohelm:ignore=true
//
// !DO NOT EDIT! Generated by genpartial
package redpanda

import (
	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
)

type PartialValues struct {
	NameOverride     *string           `json:"nameOverride,omitempty"`
	FullnameOverride *string           `json:"fullnameOverride,omitempty"`
	ClusterDomain    *string           `json:"clusterDomain,omitempty"`
	CommonLabels     map[string]string `json:"commonLabels,omitempty"`
	NodeSelector     map[string]string `json:"nodeSelector,omitempty"`
	Affinity         *PartialAffinity  `json:"affinity,omitempty" jsonschema:"required"`
	Tolerations      []map[string]any  `json:"tolerations,omitempty"`
	Image            *PartialImage     `json:"image,omitempty" jsonschema:"required,description=Values used to define the container image to be used for Redpanda"`
	Service          *PartialService   `json:"service,omitempty"`

	LicenseKey       *string                  `json:"license_key,omitempty" jsonschema:"deprecated,pattern=^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?\\.(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$|^$"`
	LicenseSecretRef *PartialLicenseSecretRef `json:"license_secret_ref,omitempty" jsonschema:"deprecated"`
	AuditLogging     *PartialAuditLogging     `json:"auditLogging,omitempty"`
	Enterprise       *PartialEnterprise       `json:"enterprise,omitempty"`
	RackAwareness    *PartialRackAwareness    `json:"rackAwareness,omitempty"`

	Auth           *PartialAuth              `json:"auth,omitempty"`
	TLS            *PartialTLS               `json:"tls,omitempty"`
	External       *PartialExternalConfig    `json:"external,omitempty"`
	Logging        *PartialLogging           `json:"logging,omitempty"`
	Monitoring     *PartialMonitoring        `json:"monitoring,omitempty"`
	Resources      *PartialRedpandaResources `json:"resources,omitempty"`
	Storage        *PartialStorage           `json:"storage,omitempty"`
	PostInstallJob *PartialPostInstallJob    `json:"post_install_job,omitempty"`
	PostUpgradeJob *PartialPostUpgradeJob    `json:"post_upgrade_job,omitempty"`
	Statefulset    *PartialStatefulset       `json:"statefulset,omitempty"`
	ServiceAccount *PartialServiceAccount    `json:"serviceAccount,omitempty"`
	RBAC           *PartialRBAC              `json:"rbac,omitempty"`
	Tuning         *PartialTuning            `json:"tuning,omitempty"`
	Listeners      *PartialListeners         `json:"listeners,omitempty"`
	Config         *PartialConfig            `json:"config,omitempty"`
	Tests          *struct {
		Enabled *bool `json:"enabled,omitempty"`
	} `json:"tests,omitempty"`
}

type PartialAffinity struct {
	NodeAffinity    map[string]any `json:"nodeAffinity,omitempty"`
	PodAffinity     map[string]any `json:"podAffinity,omitempty"`
	PodAntiAffinity map[string]any `json:"podAntiAffinity,omitempty"`
}

type PartialSecurityContext struct {
	RunAsUser                 *int64                         `json:"runAsUser,omitempty"`
	RunAsGroup                *int64                         `json:"runAsGroup,omitempty"`
	AllowPriviledgeEscalation *bool                          `json:"allowPriviledgeEscalation,omitempty"`
	RunAsNonRoot              *bool                          `json:"runAsNonRoot,omitempty"`
	FSGroup                   *int64                         `json:"fsGroup,omitempty"`
	FSGroupChangePolicy       *corev1.PodFSGroupChangePolicy `json:"fsGroupChangePolicy,omitempty"`
}

type PartialImage struct {
	Repository *string   `json:"repository,omitempty" jsonschema:"required,default=docker.redpanda.com/redpandadata/redpanda"`
	Tag        *ImageTag `json:"tag,omitempty" jsonschema:"default=Chart.appVersion"`
	PullPolicy *string   `json:"pullPolicy,omitempty" jsonschema:"required,pattern=^(Always|Never|IfNotPresent)$,description=The Kubernetes Pod image pull policy."`
}

type PartialService struct {
	Name     *string `json:"name,omitempty"`
	Internal struct {
		Annotations map[string]string `json:"annotations,omitempty"`
	} `json:"internal,omitempty"`
}

type PartialLicenseSecretRef struct {
	SecretName *string `json:"secret_name,omitempty"`
	SecretKey  *string `json:"secret_key,omitempty"`
}

type PartialAuditLogging struct {
	Enabled                    *bool    `json:"enabled,omitempty"`
	Listener                   *string  `json:"listener,omitempty"`
	Partitions                 *int     `json:"partitions,omitempty"`
	EnabledEventTypes          []string `json:"enabledEventTypes,omitempty"`
	ExcludedTopics             []string `json:"excludedTopics,omitempty"`
	ExcludedPrincipals         []string `json:"excludedPrincipals,omitempty"`
	ClientMaxBufferSize        *int     `json:"clientMaxBufferSize,omitempty"`
	QueueDrainIntervalMS       *int     `json:"queueDrainIntervalMs,omitempty"`
	QueueMaxBufferSizeperShard *int     `json:"queueMaxBufferSizePerShard,omitempty"`
	ReplicationFactor          *int     `json:"replicationFactor,omitempty"`
}

type PartialEnterprise struct {
	License          *string `json:"license,omitempty"`
	LicenseSecretRef *struct {
		Key  *string `json:"key,omitempty"`
		Name *string `json:"name,omitempty"`
	} `json:"licenseSecretRef,omitempty"`
}

type PartialRackAwareness struct {
	Enabled        *bool   `json:"enabled,omitempty" jsonschema:"required"`
	NodeAnnotation *string `json:"nodeAnnotation,omitempty" jsonschema:"required"`
}

type PartialAuth struct {
	SASL *PartialSASLAuth `json:"sasl,omitempty" jsonschema:"required"`
}

type PartialTLS struct {
	Enabled *bool              `json:"enabled,omitempty" jsonschema:"required"`
	Certs   *PartialTLSCertMap `json:"certs,omitempty"`
}

type PartialExternalConfig struct {
	Addresses      []string          `json:"addresses,omitempty"`
	Annotations    map[string]string `json:"annotations,omitempty"`
	Domain         *string           `json:"domain,omitempty"`
	Enabled        *bool             `json:"enabled,omitempty" jsonschema:"required"`
	Type           *string           `json:"type,omitempty" jsonschema:"pattern=^(LoadBalancer|NodePort)$"`
	PrefixTemplate *string           `json:"prefixTemplate,omitempty"`
	SourceRanges   []string          `json:"sourceRanges,omitempty"`
	Service        *struct {
		Enabled *bool `json:"enabled,omitempty"`
	} `json:"service,omitempty"`
	ExternalDNS *struct {
		Enabled *bool `json:"enabled,omitempty" jsonschema:"required"`
	} `json:"externalDns,omitempty"`
}

type PartialLogging struct {
	LogLevel    *string `json:"logLevel,omitempty" jsonschema:"required,pattern=^(error|warn|info|debug|trace)$"`
	UseageStats struct {
		Enabled      *bool   `json:"enabled,omitempty" jsonschema:"required"`
		Organization *string `json:"organization,omitempty"`
		ClusterID    *string `json:"clusterId,omitempty"`
	} `json:"usageStats,omitempty" jsonschema:"required"`
}

type PartialMonitoring struct {
	Enabled        *bool                   `json:"enabled,omitempty" jsonschema:"required"`
	ScrapeInterval monitoringv1.Duration   `json:"scrapeInterval,omitempty" jsonschema:"required"`
	Labels         map[string]string       `json:"labels,omitempty"`
	TLSConfig      *monitoringv1.TLSConfig `json:"tlsConfig,omitempty"`
	EnableHttp2    *bool                   `json:"enableHttp2,omitempty"`
}

type PartialRedpandaResources struct {
	CPU struct {
		Cores           any   `json:"cores,omitempty" jsonschema:"required,oneof_type=integer;string"`
		Overprovisioned *bool `json:"overprovisioned,omitempty"`
	} `json:"cpu,omitempty" jsonschema:"required"`

	Memory struct {
		EnableMemoryLocking *bool `json:"enable_memory_locking,omitempty"`

		Container struct {
			Min *MemoryAmount `json:"min,omitempty"`

			Max *MemoryAmount `json:"max,omitempty" jsonschema:"required"`
		} `json:"container,omitempty" jsonschema:"required"`

		Redpanda *struct {
			Memory *MemoryAmount `json:"memory,omitempty" jsonschema:"oneof_type=integer;string"`

			ReserveMemory *MemoryAmount `json:"reserveMemory,omitempty" jsonschema:"oneof_type=integer;string"`
		} `json:"redpanda,omitempty"`
	} `json:"memory,omitempty" jsonschema:"required"`
}

type PartialStorage struct {
	HostPath         *string        `json:"hostPath,omitempty" jsonschema:"required"`
	Tiered           *PartialTiered `json:"tiered,omitempty" jsonschema:"required"`
	PersistentVolume *struct {
		Annotations  map[string]string `json:"annotations,omitempty" jsonschema:"required"`
		Enabled      *bool             `json:"enabled,omitempty" jsonschema:"required"`
		Labels       map[string]string `json:"labels,omitempty" jsonschema:"required"`
		Size         *MemoryAmount     `json:"size,omitempty" jsonschema:"required"`
		StorageClass *string           `json:"storageClass,omitempty" jsonschema:"required"`
	} `json:"persistentVolume,omitempty" jsonschema:"required,deprecated"`
	TieredConfig                  *PartialTieredStorageConfig `json:"tieredConfig,omitempty" jsonschema:"deprecated"`
	TieredStorageHostPath         *string                     `json:"tieredStorageHostPath,omitempty" jsonschema:"deprecated"`
	TieredStoragePersistentVolume *struct {
		Annotations  map[string]string `json:"annotations,omitempty" jsonschema:"required"`
		Enabled      *bool             `json:"enabled,omitempty" jsonschema:"required"`
		Labels       map[string]string `json:"labels,omitempty" jsonschema:"required"`
		StorageClass *string           `json:"storageClass,omitempty" jsonschema:"required"`
	} `json:"tieredStoragePersistentVolume,omitempty" jsonschema:"deprecated"`
}

type PartialPostInstallJob struct {
	Resources *PartialJobResources `json:"resources,omitempty"`
	Affinity  map[string]any       `json:"affinity,omitempty"`
}

type PartialPostUpgradeJob struct {
	Resources    *PartialJobResources `json:"resources,omitempty"`
	Affinity     map[string]any       `json:"affinity,omitempty"`
	ExtraEnv     any                  `json:"extraEnv,omitempty" jsonschema:"oneof_type=array;string"`
	ExtraEnvFrom any                  `json:"extraEnvFrom,omitempty" jsonschema:"oneof_type=array;string"`
}

type PartialContainer struct {
	Name ContainerName   `json:"name,omitempty" jsonschema:"required"`
	Env  []corev1.EnvVar `json:"env,omitempty" jsonschema:"required"`
}

type PartialPodSpec struct {
	Containers []PartialContainer `json:"containers,omitempty" jsonschema:"required"`
}

type PartialPodTemplate struct {
	Labels      map[string]string `json:"labels,omitempty" jsonschema:"required"`
	Annotations map[string]string `json:"annotations,omitempty" jsonschema:"required"`
	Spec        *PartialPodSpec   `json:"spec,omitempty" jsonschema:"required"`
}

type PartialStatefulset struct {
	AdditionalSelectorLabels map[string]string `json:"additionalSelectorLabels,omitempty" jsonschema:"required"`
	NodeAffinity             map[string]any    `json:"nodeAffinity,omitempty"`
	Replicas                 *int              `json:"replicas,omitempty" jsonschema:"required"`
	UpdateStrategy           struct {
		Type *string `json:"type,omitempty" jsonschema:"required,pattern=^(RollingUpdate|OnDelete)$"`
	} `json:"updateStrategy,omitempty" jsonschema:"required"`
	AdditionalRedpandaCmdFlags []string `json:"additionalRedpandaCmdFlags,omitempty"`

	Annotations map[string]string   `json:"annotations,omitempty" jsonschema:"deprecated"`
	PodTemplate *PartialPodTemplate `json:"podTemplate,omitempty" jsonschema:"required"`
	Budget      struct {
		MaxUnavailable *int `json:"maxUnavailable,omitempty" jsonschema:"required"`
	} `json:"budget,omitempty" jsonschema:"required"`
	StartupProbe struct {
		InitialDelaySeconds *int `json:"initialDelaySeconds,omitempty" jsonschema:"required"`
		FailureThreshold    *int `json:"failureThreshold,omitempty" jsonschema:"required"`
		PeriodSeconds       *int `json:"periodSeconds,omitempty" jsonschema:"required"`
	} `json:"startupProbe,omitempty" jsonschema:"required"`
	LivenessProbe struct {
		InitialDelaySeconds *int `json:"initialDelaySeconds,omitempty" jsonschema:"required"`
		FailureThreshold    *int `json:"failureThreshold,omitempty" jsonschema:"required"`
		PeriodSeconds       *int `json:"periodSeconds,omitempty" jsonschema:"required"`
	} `json:"livenessProbe,omitempty" jsonschema:"required"`
	ReadinessProbe struct {
		InitialDelaySeconds *int `json:"initialDelaySeconds,omitempty" jsonschema:"required"`
		FailureThreshold    *int `json:"failureThreshold,omitempty" jsonschema:"required"`
		PeriodSeconds       *int `json:"periodSeconds,omitempty" jsonschema:"required"`
	} `json:"readinessProbe,omitempty" jsonschema:"required"`
	PodAffinity     map[string]any `json:"podAffinity,omitempty" jsonschema:"required"`
	PodAntiAffinity struct {
		TopologyKey *string        `json:"topologyKey,omitempty" jsonschema:"required"`
		Type        *string        `json:"type,omitempty" jsonschema:"required,pattern=^(hard|soft|custom)$"`
		Weight      *int           `json:"weight,omitempty" jsonschema:"required"`
		Custom      map[string]any `json:"custom,omitempty"`
	} `json:"podAntiAffinity,omitempty" jsonschema:"required"`
	NodeSelector                  map[string]string `json:"nodeSelector,omitempty" jsonschema:"required"`
	PriorityClassName             *string           `json:"priorityClassName,omitempty" jsonschema:"required"`
	TerminationGracePeriodSeconds *int              `json:"terminationGracePeriodSeconds,omitempty"`
	TopologySpreadConstraints     []struct {
		MaxSkew           *int    `json:"maxSkew,omitempty"`
		TopologyKey       *string `json:"topologyKey,omitempty"`
		WhenUnsatisfiable *string `json:"whenUnsatisfiable,omitempty" jsonschema:"pattern=^(ScheduleAnyway|DoNotSchedule)$"`
	} `json:"topologySpreadConstraints,omitempty" jsonschema:"required,minItems=1"`
	Tolerations []any `json:"tolerations,omitempty" jsonschema:"required"`

	PodSecurityContext *PartialSecurityContext `json:"podSecurityContext,omitempty"`
	SecurityContext    *PartialSecurityContext `json:"securityContext,omitempty" jsonschema:"required"`
	SideCars           struct {
		ConfigWatcher struct {
			Enabled           *bool                   `json:"enabled,omitempty"`
			ExtraVolumeMounts *string                 `json:"extraVolumeMounts,omitempty"`
			Resources         map[string]any          `json:"resources,omitempty"`
			SecurityContext   *corev1.SecurityContext `json:"securityContext,omitempty"`
		} `json:"configWatcher,omitempty"`
		Controllers struct {
			Image struct {
				Tag        *ImageTag `json:"tag,omitempty" jsonschema:"required,default=Chart.appVersion"`
				Repository *string   `json:"repository,omitempty" jsonschema:"required,default=docker.redpanda.com/redpandadata/redpanda-operator"`
			} `json:"image,omitempty"`
			Enabled         *bool                   `json:"enabled,omitempty"`
			Resources       any                     `json:"resources,omitempty"`
			SecurityContext *corev1.SecurityContext `json:"securityContext,omitempty"`
		} `json:"controllers,omitempty"`
	} `json:"sideCars,omitempty" jsonschema:"required"`
	ExtraVolumes      *string `json:"extraVolumes,omitempty"`
	ExtraVolumeMounts *string `json:"extraVolumeMounts,omitempty"`
	InitContainers    struct {
		Configurator struct {
			ExtraVolumeMounts *string        `json:"extraVolumeMounts,omitempty"`
			Resources         map[string]any `json:"resources,omitempty"`
		} `json:"configurator,omitempty"`
		FSValidator struct {
			Enabled           *bool          `json:"enabled,omitempty"`
			Resources         map[string]any `json:"resources,omitempty"`
			ExtraVolumeMounts *string        `json:"extraVolumeMounts,omitempty"`
			ExpectedFS        *string        `json:"expectedFS,omitempty"`
		} `json:"fsValidator,omitempty"`
		SetDataDirOwnership struct {
			Enabled           *bool          `json:"enabled,omitempty"`
			Resources         map[string]any `json:"resources,omitempty"`
			ExtraVolumeMounts *string        `json:"extraVolumeMounts,omitempty"`
		} `json:"setDataDirOwnership,omitempty"`
		SetTieredStorageCacheDirOwnership struct {
			Resources         map[string]any `json:"resources,omitempty"`
			ExtraVolumeMounts *string        `json:"extraVolumeMounts,omitempty"`
		} `json:"setTieredStorageCacheDirOwnership,omitempty"`
		Tuning struct {
			Resources         map[string]any `json:"resources,omitempty"`
			ExtraVolumeMounts *string        `json:"extraVolumeMounts,omitempty"`
		} `json:"tuning,omitempty"`
		ExtraInitContainers *string `json:"extraInitContainers,omitempty"`
	} `json:"initContainers,omitempty"`
}

type PartialServiceAccount struct {
	Create      *bool             `json:"create,omitempty" jsonschema:"required"`
	Name        *string           `json:"name,omitempty" jsonschema:"required"`
	Annotations map[string]string `json:"annotations,omitempty" jsonschema:"required"`
}

type PartialRBAC struct {
	Enabled     *bool             `json:"enabled,omitempty" jsonschema:"required"`
	Annotations map[string]string `json:"annotations,omitempty" jsonschema:"required"`
}

type PartialTuning struct {
	TuneAIOEvents   *bool   `json:"tune_aio_events,omitempty"`
	TuneClocksource *bool   `json:"tune_clocksource,omitempty"`
	TuneBallastFile *bool   `json:"tune_ballast_file,omitempty"`
	BallastFilePath *string `json:"ballast_file_path,omitempty"`
	BallastFileSize *string `json:"ballast_file_size,omitempty"`
	WellKnownIO     *string `json:"well_known_io,omitempty"`
}

type PartialListeners struct {
	Admin          *PartialAdminListeners          `json:"admin,omitempty" jsonschema:"required"`
	HTTP           *PartialHTTPListeners           `json:"http,omitempty" jsonschema:"required"`
	Kafka          *PartialKafkaListeners          `json:"kafka,omitempty" jsonschema:"required"`
	SchemaRegistry *PartialSchemaRegistryListeners `json:"schemaRegistry,omitempty" jsonschema:"required"`
	RPC            struct {
		Port *int                `json:"port,omitempty" jsonschema:"required"`
		TLS  *PartialInternalTLS `json:"tls,omitempty" jsonschema:"required"`
	} `json:"rpc,omitempty" jsonschema:"required"`
}

type PartialConfig struct {
	Cluster              *PartialClusterConfig        `json:"cluster,omitempty" jsonschema:"required"`
	Node                 *PartialNodeConfig           `json:"node,omitempty" jsonschema:"required"`
	RPK                  map[string]any               `json:"rpk,omitempty"`
	SchemaRegistryClient *PartialSchemaRegistryClient `json:"schema_registry_client,omitempty"`
	PandaProxyClient     *PartialPandaProxyClient     `json:"pandaproxy_client,omitempty"`
	Tunable              *PartialTunableConfig        `json:"tunable,omitempty" jsonschema:"required"`
}

type PartialJobResources struct {
	Limits struct {
		CPU    any           `json:"cpu,omitempty" jsonschema:"oneof_type=integer;string"`
		Memory *MemoryAmount `json:"memory,omitempty"`
	} `json:"limits,omitempty"`
	Requests struct {
		CPU    any           `json:"cpu,omitempty" jsonschema:"oneof_type=integer;string"`
		Memory *MemoryAmount `json:"memory,omitempty"`
	} `json:"requests,omitempty"`
}

type PartialSchemaRegistryClient struct {
	Retries                     *int `json:"retries,omitempty"`
	RetryBaseBackoffMS          *int `json:"retry_base_backoff_ms,omitempty"`
	ProduceBatchRecordCount     *int `json:"produce_batch_record_count,omitempty"`
	ProduceBatchSizeBytes       *int `json:"produce_batch_size_bytes,omitempty"`
	ProduceBatchDelayMS         *int `json:"produce_batch_delay_ms,omitempty"`
	ConsumerRequestTimeoutMS    *int `json:"consumer_request_timeout_ms,omitempty"`
	ConsumerRequestMaxBytes     *int `json:"consumer_request_max_bytes,omitempty"`
	ConsumerSessionTimeoutMS    *int `json:"consumer_session_timeout_ms,omitempty"`
	ConsumerRebalanceTimeoutMS  *int `json:"consumer_rebalance_timeout_ms,omitempty"`
	ConsumerHeartbeatIntervalMS *int `json:"consumer_heartbeat_interval_ms,omitempty"`
}

type PartialPandaProxyClient struct {
	Retries                     *int `json:"retries,omitempty"`
	RetryBaseBackoffMS          *int `json:"retry_base_backoff_ms,omitempty"`
	ProduceBatchRecordCount     *int `json:"produce_batch_record_count,omitempty"`
	ProduceBatchSizeBytes       *int `json:"produce_batch_size_bytes,omitempty"`
	ProduceBatchDelayMS         *int `json:"produce_batch_delay_ms,omitempty"`
	ConsumerRequestTimeoutMS    *int `json:"consumer_request_timeout_ms,omitempty"`
	ConsumerRequestMaxBytes     *int `json:"consumer_request_max_bytes,omitempty"`
	ConsumerSessionTimeoutMS    *int `json:"consumer_session_timeout_ms,omitempty"`
	ConsumerRebalanceTimeoutMS  *int `json:"consumer_rebalance_timeout_ms,omitempty"`
	ConsumerHeartbeatIntervalMS *int `json:"consumer_heartbeat_interval_ms,omitempty"`
}

type PartialTLSCert struct {
	Enabled               *bool                     `json:"enabled,omitempty"`
	CAEnabled             *bool                     `json:"caEnabled,omitempty" jsonschema:"required"`
	ApplyInternalDNSNames *bool                     `json:"applyInternalDNSNames,omitempty"`
	Duration              *string                   `json:"duration,omitempty" jsonschema:"pattern=.*[smh]$"`
	IssuerRef             *cmmeta.ObjectReference   `json:"issuerRef,omitempty"`
	SecretRef             *PartialNameOnlySecretRef `json:"secretRef,omitempty"`
}

type PartialNameOnlySecretRef struct {
	Name *string `json:"name,omitempty"`
}

type PartialTLSCertMap map[string]PartialTLSCert

type PartialSASLUser struct {
	Name      *string `json:"name,omitempty"`
	Password  *string `json:"password,omitempty"`
	Mechanism *string `json:"mechanism,omitempty" jsonschema:"pattern=^(SCRAM-SHA-512|SCRAM-SHA-256)$"`
}

type PartialSASLAuth struct {
	Enabled   *bool             `json:"enabled,omitempty" jsonschema:"required"`
	Mechanism *string           `json:"mechanism,omitempty"`
	SecretRef *string           `json:"secretRef,omitempty"`
	Users     []PartialSASLUser `json:"users,omitempty"`
}

type PartialInternalTLS struct {
	Cert              *string `json:"cert,omitempty" jsonschema:"required"`
	Enabled           *bool   `json:"enabled,omitempty"`
	RequireClientAuth *bool   `json:"requireClientAuth,omitempty" jsonschema:"required"`
}

type PartialExternalTLS struct {
	Cert              *string `json:"cert,omitempty"`
	Enabled           *bool   `json:"enabled,omitempty"`
	RequireClientAuth *bool   `json:"requireClientAuth,omitempty"`
}

type PartialAdminListeners struct {
	External PartialExternalListeners[PartialAdminExternal] `json:"external,omitempty"`
	Port     *int                                           `json:"port,omitempty" jsonschema:"required"`
	TLS      *PartialInternalTLS                            `json:"tls,omitempty" jsonschema:"required"`
}

type PartialAdminExternal struct {
	AdvertisedPorts []int32             `json:"advertisedPorts,omitempty" jsonschema:"minItems=1"`
	Enabled         *bool               `json:"enabled,omitempty"`
	Port            *int32              `json:"port,omitempty" jsonschema:"required"`
	TLS             *PartialExternalTLS `json:"tls,omitempty"`
}

type PartialHTTPListeners struct {
	Enabled              *bool                                         `json:"enabled,omitempty" jsonschema:"required"`
	External             PartialExternalListeners[PartialHTTPExternal] `json:"external,omitempty"`
	AuthenticationMethod *HTTPAuthenticationMethod                     `json:"authenticationMethod,omitempty"`
	TLS                  *PartialInternalTLS                           `json:"tls,omitempty" jsonschema:"required"`
	KafkaEndpoint        *string                                       `json:"kafkaEndpoint,omitempty" jsonschema:"required,pattern=^[A-Za-z_-][A-Za-z0-9_-]*$"`
	Port                 *int                                          `json:"port,omitempty" jsonschema:"required"`
}

type PartialHTTPExternal struct {
	AdvertisedPorts      []int32                   `json:"advertisedPorts,omitempty" jsonschema:"minItems=1"`
	Enabled              *bool                     `json:"enabled,omitempty"`
	Port                 *int32                    `json:"port,omitempty" jsonschema:"required"`
	AuthenticationMethod *HTTPAuthenticationMethod `json:"authenticationMethod,omitempty"`
	PrefixTemplate       *string                   `json:"prefixTemplate,omitempty"`
	TLS                  *PartialExternalTLS       `json:"tls,omitempty" jsonschema:"required"`
}

type PartialKafkaListeners struct {
	AuthenticationMethod *KafkaAuthenticationMethod                     `json:"authenticationMethod,omitempty"`
	External             PartialExternalListeners[PartialKafkaExternal] `json:"external,omitempty"`
	TLS                  *PartialInternalTLS                            `json:"tls,omitempty" jsonschema:"required"`
	Port                 *int                                           `json:"port,omitempty" jsonschema:"required"`
}

type PartialKafkaExternal struct {
	AdvertisedPorts      []int32                    `json:"advertisedPorts,omitempty" jsonschema:"minItems=1"`
	Enabled              *bool                      `json:"enabled,omitempty"`
	Port                 *int32                     `json:"port,omitempty" jsonschema:"required"`
	AuthenticationMethod *KafkaAuthenticationMethod `json:"authenticationMethod,omitempty"`
	PrefixTemplate       *string                    `json:"prefixTemplate,omitempty"`
	TLS                  *PartialExternalTLS        `json:"tls,omitempty"`
}

type PartialSchemaRegistryListeners struct {
	Enabled              *bool                                                   `json:"enabled,omitempty" jsonschema:"required"`
	External             PartialExternalListeners[PartialSchemaRegistryExternal] `json:"external,omitempty"`
	AuthenticationMethod *HTTPAuthenticationMethod                               `json:"authenticationMethod,omitempty"`
	KafkaEndpoint        *string                                                 `json:"kafkaEndpoint,omitempty" jsonschema:"required,pattern=^[A-Za-z_-][A-Za-z0-9_-]*$"`
	Port                 *int                                                    `json:"port,omitempty" jsonschema:"required"`
	TLS                  *PartialInternalTLS                                     `json:"tls,omitempty" jsonschema:"required"`
}

type PartialSchemaRegistryExternal struct {
	AdvertisedPorts      []int32                   `json:"advertisedPorts,omitempty" jsonschema:"minItems=1"`
	Enabled              *bool                     `json:"enabled,omitempty"`
	Port                 *int32                    `json:"port,omitempty"`
	AuthenticationMethod *HTTPAuthenticationMethod `json:"authenticationMethod,omitempty"`
	TLS                  *PartialExternalTLS       `json:"tls,omitempty"`
}

type PartialTunableConfig map[string]any

type PartialNodeConfig map[string]any

type PartialClusterConfig map[string]any

type PartialSecretRef struct {
	ConfigurationKey *string `json:"configurationKey,omitempty"`
	Key              *string `json:"key,omitempty"`
	Name             *string `json:"name,omitempty"`
}

type PartialTieredStorageCredentials struct {
	ConfigurationKey *string           `json:"configurationKey,omitempty" jsonschema:"deprecated"`
	Key              *string           `json:"key,omitempty" jsonschema:"deprecated"`
	Name             *string           `json:"name,omitempty" jsonschema:"deprecated"`
	AccessKey        *PartialSecretRef `json:"accessKey,omitempty"`
	SecretKey        *PartialSecretRef `json:"secretKey,omitempty"`
}

type PartialTieredStorageConfig map[string]any

type PartialTiered struct {
	CredentialsSecretRef *PartialTieredStorageCredentials `json:"credentialsSecretRef,omitempty"`
	Config               *PartialTieredStorageConfig      `json:"config,omitempty"`
	HostPath             *string                          `json:"hostPath,omitempty"`
	MountType            *string                          `json:"mountType,omitempty" jsonschema:"required,pattern=^(none|hostPath|emptyDir|persistentVolume)$"`
	PersistentVolume     struct {
		Annotations   map[string]string `json:"annotations,omitempty" jsonschema:"required"`
		Enabled       *bool             `json:"enabled,omitempty"`
		Labels        map[string]string `json:"labels,omitempty" jsonschema:"required"`
		NameOverwrite *string           `json:"nameOverwrite,omitempty"`
		Size          *string           `json:"size,omitempty"`
		StorageClass  *string           `json:"storageClass,omitempty" jsonschema:"required"`
	} `json:"persistentVolume,omitempty"`
}

type PartialExternalListeners[T any] map[string]T
