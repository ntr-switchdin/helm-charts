package redpanda_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"maps"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/redpanda-data/helm-charts/charts/redpanda"
	"github.com/redpanda-data/helm-charts/pkg/gotohelm/helmette"
	"github.com/redpanda-data/helm-charts/pkg/helm"
	"github.com/redpanda-data/helm-charts/pkg/helm/helmtest"
	"github.com/redpanda-data/helm-charts/pkg/kube"
	"github.com/redpanda-data/helm-charts/pkg/testutil"
	"github.com/redpanda-data/helm-charts/pkg/valuesutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zapio"
	"go.uber.org/zap/zaptest"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func kafkaListenerTest(t *testing.T, ctx context.Context, rpk Client) {
	input := "test-input"
	topicName := "testTopic"
	_, err := rpk.CreateTopic(ctx, topicName)
	require.NoError(t, err)

	_, err = rpk.KafkaProduce(ctx, input, topicName)
	require.NoError(t, err)

	consumeOutput, err := rpk.KafkaConsume(ctx, topicName)
	require.NoError(t, err)
	require.Equal(t, input, consumeOutput["value"])
}

func adminListenerTest(t *testing.T, ctx context.Context, rpk Client) {
	require.Eventually(t, func() bool {
		out, err := rpk.GetClusterHealth(ctx)
		if err != nil {
			return false
		}
		return out["is_healthy"].(bool)
	}, time.Minute, 5*time.Second)
}

func schemaRegistryListenerTest(t *testing.T, ctx context.Context, rpk Client) {
	// Test schema registry
	// Based on https://docs.redpanda.com/current/manage/schema-reg/schema-reg-api/
	formats, err := rpk.QuerySupportedFormats(ctx)
	require.NoError(t, err)
	require.Len(t, formats, 2)

	schema := map[string]any{
		"type": "record",
		"name": "sensor_sample",
		"fields": []map[string]any{
			{
				"name":        "timestamp",
				"type":        "long",
				"logicalType": "timestamp-millis",
			},
			{
				"name":        "identifier",
				"type":        "string",
				"logicalType": "uuid",
			},
			{
				"name": "value",
				"type": "long",
			},
		},
	}

	registeredID, err := rpk.RegisterSchema(ctx, schema)
	require.NoError(t, err)

	var id float64
	if idForSchema, ok := registeredID["id"]; ok {
		id = idForSchema.(float64)
	}

	schemaBytes, err := json.Marshal(schema)
	require.NoError(t, err)

	retrievedSchema, err := rpk.RetrieveSchema(ctx, int(id))
	require.NoError(t, err)
	require.JSONEq(t, string(schemaBytes), retrievedSchema)

	resp, err := rpk.ListRegistrySubjects(ctx)
	require.NoError(t, err)
	require.Equal(t, "sensor-value", resp[0])

	_, err = rpk.SoftDeleteSchema(ctx, resp[0], int(id))
	require.NoError(t, err)

	_, err = rpk.HardDeleteSchema(ctx, resp[0], int(id))
	require.NoError(t, err)
}

func httpProxyListenerTest(t *testing.T, ctx context.Context, rpk Client) {
	// Test http proxy
	// Based on https://docs.redpanda.com/current/develop/http-proxy/
	_, err := rpk.ListTopics(ctx)
	require.NoError(t, err)

	records := map[string]any{
		"records": []map[string]any{
			{
				"value":     "Redpanda",
				"partition": 0,
			},
			{
				"value":     "HTTP proxy",
				"partition": 1,
			},
			{
				"value":     "Test event",
				"partition": 2,
			},
		},
	}

	httpTestTopic := "httpTestTopic"
	_, err = rpk.CreateTopic(ctx, httpTestTopic)
	require.NoError(t, err)

	_, err = rpk.SendEventToTopic(ctx, records, httpTestTopic)
	require.NoError(t, err)

	time.Sleep(time.Second * 5)

	record, err := rpk.RetrieveEventFromTopic(ctx, httpTestTopic, 0)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("[{\"topic\":\"%s\",\"key\":null,\"value\":\"Redpanda\",\"partition\":0,\"offset\":0}]", httpTestTopic), record)

	record, err = rpk.RetrieveEventFromTopic(ctx, httpTestTopic, 1)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("[{\"topic\":\"%s\",\"key\":null,\"value\":\"HTTP proxy\",\"partition\":1,\"offset\":0}]", httpTestTopic), record)

	record, err = rpk.RetrieveEventFromTopic(ctx, httpTestTopic, 2)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("[{\"topic\":\"%s\",\"key\":null,\"value\":\"Test event\",\"partition\":2,\"offset\":0}]", httpTestTopic), record)
}

func mTLSValuesUsingCertManager() redpanda.PartialValues {
	return redpanda.PartialValues{
		ClusterDomain: ptr.To("cluster.local"),
		Listeners: &redpanda.PartialListeners{
			Admin: &redpanda.PartialAdminListeners{
				TLS: &redpanda.PartialInternalTLS{
					RequireClientAuth: ptr.To(true),
				},
			},
			HTTP: &redpanda.PartialHTTPListeners{
				TLS: &redpanda.PartialInternalTLS{
					RequireClientAuth: ptr.To(true),
				},
			},
			Kafka: &redpanda.PartialKafkaListeners{
				TLS: &redpanda.PartialInternalTLS{
					RequireClientAuth: ptr.To(true),
				},
			},
			SchemaRegistry: &redpanda.PartialSchemaRegistryListeners{
				TLS: &redpanda.PartialInternalTLS{
					RequireClientAuth: ptr.To(true),
				},
			},
			RPC: &struct {
				Port *int32                       `json:"port,omitempty" jsonschema:"required"`
				TLS  *redpanda.PartialInternalTLS `json:"tls,omitempty" jsonschema:"required"`
			}{
				TLS: &redpanda.PartialInternalTLS{
					RequireClientAuth: ptr.To(true),
				},
			},
		},
	}
}

func mTLSValuesWithProvidedCerts(serverTLSSecretName, clientTLSSecretName string) redpanda.PartialValues {
	return redpanda.PartialValues{
		ClusterDomain: ptr.To("cluster.local"),
		TLS: &redpanda.PartialTLS{
			Enabled: ptr.To(true),
			Certs: redpanda.PartialTLSCertMap{
				"provided": redpanda.PartialTLSCert{
					Enabled:         ptr.To(true),
					CAEnabled:       ptr.To(true),
					SecretRef:       &corev1.LocalObjectReference{Name: serverTLSSecretName},
					ClientSecretRef: &corev1.LocalObjectReference{Name: clientTLSSecretName},
				},
				"default": redpanda.PartialTLSCert{Enabled: ptr.To(false)},
			},
		},
		Listeners: &redpanda.PartialListeners{
			Admin: &redpanda.PartialAdminListeners{
				TLS: &redpanda.PartialInternalTLS{
					RequireClientAuth: ptr.To(true),
					Cert:              ptr.To("provided"),
				},
			},
			HTTP: &redpanda.PartialHTTPListeners{
				TLS: &redpanda.PartialInternalTLS{
					RequireClientAuth: ptr.To(true),
					Cert:              ptr.To("provided"),
				},
			},
			Kafka: &redpanda.PartialKafkaListeners{
				TLS: &redpanda.PartialInternalTLS{
					RequireClientAuth: ptr.To(true),
					Cert:              ptr.To("provided"),
				},
			},
			SchemaRegistry: &redpanda.PartialSchemaRegistryListeners{
				TLS: &redpanda.PartialInternalTLS{
					RequireClientAuth: ptr.To(true),
					Cert:              ptr.To("provided"),
				},
			},
			RPC: &struct {
				Port *int32                       `json:"port,omitempty" jsonschema:"required"`
				TLS  *redpanda.PartialInternalTLS `json:"tls,omitempty" jsonschema:"required"`
			}{
				TLS: &redpanda.PartialInternalTLS{
					RequireClientAuth: ptr.To(true),
					Cert:              ptr.To("provided"),
				},
			},
		},
	}
}

func TieredStorageStatic(t *testing.T) redpanda.PartialValues {
	license := os.Getenv("REDPANDA_LICENSE")
	if license == "" {
		t.Skipf("$REDPANDA_LICENSE is not set")
	}

	return redpanda.PartialValues{
		Config: &redpanda.PartialConfig{
			Node: redpanda.PartialNodeConfig{
				"developer_mode": true,
			},
		},
		Enterprise: &redpanda.PartialEnterprise{
			License: &license,
		},
		Storage: &redpanda.PartialStorage{
			Tiered: &redpanda.PartialTiered{
				Config: redpanda.PartialTieredStorageConfig{
					"cloud_storage_enabled":    true,
					"cloud_storage_region":     "static-region",
					"cloud_storage_bucket":     "static-bucket",
					"cloud_storage_access_key": "static-access-key",
					"cloud_storage_secret_key": "static-secret-key",
				},
			},
		},
	}
}

func TieredStorageSecret(namespace string) corev1.Secret {
	return corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "tiered-storage-",
			Namespace:    namespace,
		},
		Data: map[string][]byte{
			"access": []byte("from-secret-access-key"),
			"secret": []byte("from-secret-secret-key"),
		},
	}
}

func TieredStorageSecretRefs(t *testing.T, secret *corev1.Secret) redpanda.PartialValues {
	license := os.Getenv("REDPANDA_LICENSE")
	if license == "" {
		t.Skipf("$REDPANDA_LICENSE is not set")
	}

	access := "access"
	secretKey := "secret"
	return redpanda.PartialValues{
		Config: &redpanda.PartialConfig{
			Node: redpanda.PartialNodeConfig{
				"developer_mode": true,
			},
		},
		Enterprise: &redpanda.PartialEnterprise{
			License: &license,
		},
		Storage: &redpanda.PartialStorage{
			Tiered: &redpanda.PartialTiered{
				CredentialsSecretRef: &redpanda.PartialTieredStorageCredentials{
					AccessKey: &redpanda.PartialSecretRef{Name: &secret.Name, Key: &access},
					SecretKey: &redpanda.PartialSecretRef{Name: &secret.Name, Key: &secretKey},
				},
				Config: redpanda.PartialTieredStorageConfig{
					"cloud_storage_enabled": true,
					"cloud_storage_region":  "a-region",
					"cloud_storage_bucket":  "a-bucket",
				},
			},
		},
	}
}

func TestChart(t *testing.T) {
	if testing.Short() {
		t.Skipf("Skipping log running test...")
	}

	log := zaptest.NewLogger(t)
	w := &zapio.Writer{Log: log, Level: zapcore.InfoLevel}
	wErr := &zapio.Writer{Log: log, Level: zapcore.ErrorLevel}

	redpandaChart := "."

	h := helmtest.Setup(t)

	t.Run("tiered-storage-secrets", func(t *testing.T) {
		ctx := testutil.Context(t)

		env := h.Namespaced(t)

		credsSecret, err := kube.Create(ctx, env.Ctl(), TieredStorageSecret(env.Namespace()))
		require.NoError(t, err)

		rpRelease := env.Install(ctx, redpandaChart, helm.InstallOptions{
			Values: redpanda.PartialValues{
				Config: &redpanda.PartialConfig{
					Node: redpanda.PartialNodeConfig{
						"developer_mode": true,
					},
				},
			},
		})

		rpk := Client{Ctl: env.Ctl(), Release: &rpRelease}

		config, err := rpk.ClusterConfig(ctx)
		require.NoError(t, err)
		require.Equal(t, false, config["cloud_storage_enabled"])

		rpRelease = env.Upgrade(redpandaChart, rpRelease, helm.UpgradeOptions{Values: TieredStorageStatic(t)})

		config, err = rpk.ClusterConfig(ctx)
		require.NoError(t, err)
		require.Equal(t, true, config["cloud_storage_enabled"])
		require.Equal(t, "static-access-key", config["cloud_storage_access_key"])

		rpRelease = env.Upgrade(redpandaChart, rpRelease, helm.UpgradeOptions{Values: TieredStorageSecretRefs(t, credsSecret)})

		config, err = rpk.ClusterConfig(ctx)
		require.NoError(t, err)
		require.Equal(t, true, config["cloud_storage_enabled"])
		require.Equal(t, "from-secret-access-key", config["cloud_storage_access_key"])
	})

	t.Run("mtls-using-cert-manager", func(t *testing.T) {
		ctx := testutil.Context(t)

		env := h.Namespaced(t)

		partial := mTLSValuesUsingCertManager()

		rpRelease := env.Install(ctx, redpandaChart, helm.InstallOptions{
			Values: partial,
		})

		rpk := Client{Ctl: env.Ctl(), Release: &rpRelease}

		v, err := convertPartialToValues(&partial)
		require.NoError(t, err)

		dot := &helmette.Dot{
			Values:  *v,
			Release: helmette.Release{Name: rpRelease.Name, Namespace: rpRelease.Namespace},
			Chart: helmette.Chart{
				Name: "redpanda",
			},
		}

		cleanup, err := rpk.ExposeRedpandaCluster(ctx, dot, w, wErr)
		if cleanup != nil {
			t.Cleanup(cleanup)
		}
		require.NoError(t, err)

		t.Run("kafka-listener", func(t *testing.T) {
			kafkaListenerTest(t, ctx, rpk)
		})

		t.Run("admin-listener", func(t *testing.T) {
			adminListenerTest(t, ctx, rpk)
		})

		t.Run("schema-registry-listener", func(t *testing.T) {
			schemaRegistryListenerTest(t, ctx, rpk)
		})

		t.Run("http-proxy-listener", func(t *testing.T) {
			httpProxyListenerTest(t, ctx, rpk)
		})
	})

	t.Run("mtls-using-self-created-certificates", func(t *testing.T) {
		ctx := testutil.Context(t)

		env := h.Namespaced(t)

		serverTLSSecretName := "server-tls-secret"
		clientTLSSecretName := "client-tls-secret"

		partial := mTLSValuesWithProvidedCerts(serverTLSSecretName, clientTLSSecretName)

		r, err := rand.Int(rand.Reader, new(big.Int).SetInt64(1799999999))
		require.NoError(t, err)

		chartReleaseName := fmt.Sprintf("chart-%d", r.Int64())
		ca, sPublic, sPrivate, cPublic, cPrivate, err := createCertificates(chartReleaseName, env.Namespace())
		require.NoError(t, err)

		s := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      serverTLSSecretName,
				Namespace: env.Namespace(),
			},
			Data: map[string][]byte{
				"ca.crt":  ca,
				"tls.crt": sPublic,
				"tls.key": sPrivate,
			},
		}
		_, err = kube.Create[corev1.Secret](ctx, env.Ctl(), s)
		require.NoError(t, err)

		c := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      clientTLSSecretName,
				Namespace: env.Namespace(),
			},
			Data: map[string][]byte{
				"ca.crt":  ca,
				"tls.crt": cPublic,
				"tls.key": cPrivate,
			},
		}
		_, err = kube.Create[corev1.Secret](ctx, env.Ctl(), c)
		require.NoError(t, err)

		rpRelease := env.Install(ctx, redpandaChart, helm.InstallOptions{
			Values:    partial,
			Name:      chartReleaseName,
			Namespace: env.Namespace(),
		})

		rpk := Client{Ctl: env.Ctl(), Release: &rpRelease}

		v, err := convertPartialToValues(&partial)
		require.NoError(t, err)

		dot := &helmette.Dot{
			Values:  *v,
			Release: helmette.Release{Name: rpRelease.Name, Namespace: rpRelease.Namespace},
			Chart: helmette.Chart{
				Name: "redpanda",
			},
		}

		cleanup, err := rpk.ExposeRedpandaCluster(ctx, dot, w, wErr)
		if cleanup != nil {
			t.Cleanup(cleanup)
		}
		require.NoError(t, err)

		t.Run("kafka-listener", func(t *testing.T) {
			kafkaListenerTest(t, ctx, rpk)
		})

		t.Run("admin-listener", func(t *testing.T) {
			adminListenerTest(t, ctx, rpk)
		})

		t.Run("schema-registry-listener", func(t *testing.T) {
			schemaRegistryListenerTest(t, ctx, rpk)
		})

		t.Run("http-proxy-listener", func(t *testing.T) {
			httpProxyListenerTest(t, ctx, rpk)
		})
	})
}

func convertPartialToValues(partial *redpanda.PartialValues) (*helmette.Values, error) {
	b, err := json.Marshal(partial)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	v := &helmette.Values{}
	err = json.Unmarshal(b, v)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return v, nil
}

// getConfigMaps is parsing all manifests (resources) created by helm template
// execution. Redpanda helm chart creates 3 distinct files in ConfigMap:
// redpanda.yaml (node, tunable and cluster configuration), bootstrap.yaml
// (only cluster configuration) and profile (external connectivity rpk profile
// which is in different ConfigMap than other two).
func getConfigMaps(manifests []byte) (r *corev1.ConfigMap, rpk *corev1.ConfigMap, err error) {
	objs, err := kube.DecodeYAML(manifests, redpanda.Scheme)
	if err != nil {
		return nil, nil, err
	}

	for _, obj := range objs {
		switch obj := obj.(type) {
		case *corev1.ConfigMap:
			switch obj.Name {
			case "redpanda":
				r = obj
			case "redpanda-rpk":
				rpk = obj
			}
		}
	}

	return r, rpk, nil
}

func TestLabels(t *testing.T) {
	ctx := testutil.Context(t)
	client, err := helm.New(helm.Options{ConfigHome: testutil.TempDir(t)})
	require.NoError(t, err)

	for _, labels := range []map[string]string{
		{"foo": "bar"},
		{"baz": "1", "quux": "2"},
		// TODO: Add a test for asserting the behavior of adding a commonLabel
		// overriding a builtin value (app.kubernetes.io/name) once the
		// expected behavior is decided.
	} {
		values := &redpanda.PartialValues{
			CommonLabels: labels,
		}

		helmValues, err := valuesutil.UnmarshalInto[helmette.Values](values)
		require.NoError(t, err)

		dot := &helmette.Dot{
			Values: helmValues,
			Chart:  redpanda.ChartMeta(),
			Release: helmette.Release{
				Name:      "redpanda",
				Namespace: "redpanda",
				Service:   "Helm",
			},
		}

		manifests, err := client.Template(ctx, ".", helm.TemplateOptions{
			Name:      dot.Release.Name,
			Namespace: dot.Release.Namespace,
			// This guarantee does not currently extend to console.
			Set: []string{"console.enabled=false"},
			// Nor does it extend to tests.
			SkipTests: true,
			Values:    values,
		})
		require.NoError(t, err)

		objs, err := kube.DecodeYAML(manifests, redpanda.Scheme)
		require.NoError(t, err)

		expectedLabels := redpanda.FullLabels(dot)
		require.Subset(t, expectedLabels, values.CommonLabels, "FullLabels does not contain CommonLabels")

		for _, obj := range objs {
			// Assert that CommonLabels is included on all top level objects.
			require.Subset(t, obj.GetLabels(), expectedLabels, "%T %q", obj, obj.GetName())

			// For other objects (replication controllers) we want to assert
			// that common labels are also included on whatever object (Pod)
			// they generate/contain a template of.
			switch obj := obj.(type) {
			case *appsv1.StatefulSet:
				expectedLabels := maps.Clone(expectedLabels)
				expectedLabels["app.kubernetes.io/component"] += "-statefulset"
				require.Subset(t, obj.Spec.Template.GetLabels(), expectedLabels, "%T/%s's %T", obj, obj.Name, obj.Spec.Template)
			}
		}
	}
}

func createCertificates(chartReleaseName, chartReleaseNamespace string) (ca, serverPublic, serverPrivate, clientPublic, clientPrivate []byte, err error) {
	now := time.Now()

	rootCASubject := pkix.Name{
		CommonName:   "test.example.com",
		Organization: []string{"Σ Acme Co"},
		Country:      []string{"US"},
	}
	root := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      rootCASubject,
		NotBefore:    now.Add(-time.Hour),
		NotAfter:     now.Add(time.Hour),

		KeyUsage:    x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},

		BasicConstraintsValid: true,
		IsCA:                  true,

		DNSNames:       []string{"test.example.com"},
		EmailAddresses: []string{"gopher@golang.org"},
	}

	priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, root, root, priv.Public(), priv)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	root, err = x509.ParseCertificate(derBytes)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	ca = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	commonName := fmt.Sprintf("%s.%s.svc.cluster.local", chartReleaseName, chartReleaseNamespace)
	shortTestName := fmt.Sprintf("%s.%s", chartReleaseName, chartReleaseNamespace)
	serverTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"Σ Acme Co"},
			Country:      []string{"US"},
		},
		Issuer:    rootCASubject,
		NotBefore: now.Add(-time.Hour),
		NotAfter:  now.Add(time.Hour),

		SignatureAlgorithm: x509.ECDSAWithSHA384,

		KeyUsage:    x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},

		BasicConstraintsValid: true,
		IsCA:                  false,

		DNSNames: []string{
			shortTestName,
			commonName,
			fmt.Sprintf("%s.", commonName),
			fmt.Sprintf("*.%s", commonName),
			fmt.Sprintf("*.%s.", commonName),
		},
		EmailAddresses: []string{"gopher@golang.org"},
	}

	privServer, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	derServerBytes, err := x509.CreateCertificate(rand.Reader, &serverTemplate, root, privServer.Public(), priv)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(privServer)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	serverPublic = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derServerBytes})
	serverPrivate = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})

	clientTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "redpanda-client",
			Organization: []string{"Σ Acme Co"},
			Country:      []string{"US"},
		},
		Issuer:    rootCASubject,
		NotBefore: now.Add(-time.Hour),
		NotAfter:  now.Add(time.Hour),

		SignatureAlgorithm: x509.ECDSAWithSHA384,

		KeyUsage:    x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},

		BasicConstraintsValid: true,
		IsCA:                  false,

		DNSNames:       []string{"redpanda-client"},
		EmailAddresses: []string{"gopher@golang.org"},
	}

	privClient, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	derClientBytes, err := x509.CreateCertificate(rand.Reader, &clientTemplate, root, privClient.Public(), priv)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	privBytes, err = x509.MarshalPKCS8PrivateKey(privClient)

	clientPublic = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derClientBytes})
	clientPrivate = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})

	return
}
