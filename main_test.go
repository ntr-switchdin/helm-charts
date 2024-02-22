package main

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/redpanda-data/helm-charts/pkg/helm"
	"github.com/redpanda-data/helm-charts/pkg/kube"
	"github.com/redpanda-data/helm-charts/pkg/testutil"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func EnsureCertManager(ctx context.Context, client *helm.Client) error {
	release, err := client.Get(ctx, "cert-manager", "cert-manager")
	if err == nil && release.Status == "deployed" {
		return nil
	}
	fmt.Printf("%#v\n", release)
	fmt.Printf("%#v\n", err)

	if err := client.RepoAdd(ctx, "jetstack", "https://charts.jetstack.io"); err != nil {
		return err
	}

	_, err = client.Install(ctx, "jetstack/cert-manager", helm.InstallOptions{
		Name:            "cert-manager",
		Version:         "v1.11.0",
		Namespace:       "cert-manager",
		CreateNamespace: true,
		Values: map[string]any{
			"installCRDs": true,
		},
	})
	return err
}

func TestRedpandaChart(t *testing.T) {
	ctx := testutil.Context(t)

	ctl, err := kube.FromEnv()
	require.NoError(t, err)

	client, err := helm.New(helm.Options{
		KubeConfig: ctl.RestConfig(),
		ConfigHome: t.TempDir(),
	})
	require.NoError(t, err)

	require.NoError(t, EnsureCertManager(ctx, client))

	ns, err := kube.Create(ctx, ctl, corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "test-",
		},
		Spec: corev1.NamespaceSpec{},
	})
	require.NoError(t, err)

	testutil.MaybeCleanup(t, func() {
		require.NoError(t, kube.Delete[corev1.Namespace](ctx, ctl, ns.ObjectMeta))
	})

	// TODO(chrisseto): This is a bit kludgey as we're relying on the directory
	// that `go test` is being run from.
	release, err := client.Install(ctx, "./charts/redpanda", helm.InstallOptions{
		Namespace: ns.Name,
		Values: map[string]any{
			"config": map[string]any{
				"cluster": map[string]any{
					"disable_metrics": true,
				},
			},
		},
	})
	require.NoError(t, err)
	t.Logf("%#v\n", release)

	pod, err := kube.Get[corev1.Pod](ctx, ctl, metav1.ObjectMeta{
		Name:      fmt.Sprintf("%s-0", release.Name),
		Namespace: ns.Name,
	})
	require.NoError(t, err)

	var out bytes.Buffer
	require.NoError(t, ctl.Exec(ctx, pod, kube.ExecOptions{
		Command: []string{"rpk", "cluster", "health"},
		Stdout:  &out,
	}))

	require.Regexp(t, `Healthy:\s+true`, out.String())

	// Assert that disable_metrics is set as expected.
	out.Reset()
	require.NoError(t, ctl.Exec(ctx, pod, kube.ExecOptions{
		Command: []string{"rpk", "cluster", "config", "get", "disable_metrics"},
		Stdout:  &out,
	}))
	require.Equal(t, "true\n", out.String())

	_, err = client.Upgrade(ctx, release.Name, "./charts/redpanda", helm.UpgradeOptions{
		Namespace:   ns.Name,
		ReuseValues: true,
		Values: map[string]any{
			"config": map[string]any{
				"cluster": map[string]any{
					"disable_metrics": false,
				},
			},
		},
	})
	require.NoError(t, err)

	// Assert that disable_metrics is set as expected.
	out.Reset()
	require.NoError(t, ctl.Exec(ctx, pod, kube.ExecOptions{
		Command: []string{"rpk", "cluster", "config", "get", "disable_metrics"},
		Stdout:  &out,
	}))
	require.Equal(t, "false\n", out.String())
}
