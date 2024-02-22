package kube

import (
	"context"
	"io"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func FromEnv() (*Ctl, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	c, err := client.New(config, client.Options{})
	if err != nil {
		return nil, err
	}

	return &Ctl{
		config:    config,
		client:    c,
		namespace: metav1.NamespaceAll,
	}, nil
}

type Ctl struct {
	config    *rest.Config
	client    client.Client
	namespace string
}

// RestConfig returns a deep copy of the *[rest.Config] used by this [Ctl].
func (c *Ctl) RestConfig() *rest.Config {
	return rest.CopyConfig(c.config)
}

type ExecOptions struct {
	Container string
	Command   []string
	Stdin     io.Reader
	Stdout    io.Writer
	Stderr    io.Writer
}

func (c *Ctl) Exec(ctx context.Context, pod *corev1.Pod, opts ExecOptions) error {
	if opts.Container == "" {
		opts.Container = pod.Spec.Containers[0].Name
	}

	// Apparently, nothing in the k8s SDK, except exec'ing, uses RESTClientFor.
	// RESTClientFor checks for GroupVersion and NegotiatedSerializer which are
	// never set by the config loading tool chain.
	// The .APIPath setting was a random shot in the dark that happened to work...
	// Pulled from https://github.com/kubernetes/kubectl/blob/master/pkg/cmd/util/kubectl_match_version.go#L115
	cfg := c.RestConfig()
	cfg.APIPath = "/api"
	cfg.GroupVersion = &schema.GroupVersion{Version: "v1"}
	cfg.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	restClient, err := rest.RESTClientFor(cfg)
	if err != nil {
		return err
	}

	// Inspired by https://github.com/kubernetes/kubectl/blob/acf4a09f2daede8fdbf65514ade9426db0367ed3/pkg/cmd/exec/exec.go#L388
	req := restClient.Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec")

	req.VersionedParams(&corev1.PodExecOptions{
		Container: opts.Container,
		Command:   opts.Command,
		Stdin:     opts.Stdin != nil,
		Stdout:    opts.Stdout != nil,
		Stderr:    opts.Stderr != nil,
		TTY:       false,
	}, runtime.NewParameterCodec(c.client.Scheme()))

	// TODO(chrisseto): SPDY is reported to be deprecated but
	// NewWebSocketExecutor doesn't appear to work in our version of KinD.
	exec, err := remotecommand.NewSPDYExecutor(c.config, "POST", req.URL())
	// exec, err := remotecommand.NewWebSocketExecutor(c.config, "GET", req.URL().String())
	if err != nil {
		return err
	}

	return exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stderr: opts.Stderr,
		Stdout: opts.Stdout,
		Stdin:  opts.Stdin,
	})
}

type Object = client.Object

type ObjectList[T any] interface {
	client.ObjectList
	*T
}

type AddrofObject[T any] interface {
	*T
	Object
}

func List[T any, L ObjectList[T]](ctx context.Context, ctl *Ctl, opts ...client.ListOption) (*T, error) {
	var list T
	if err := ctl.client.List(ctx, L(&list), opts...); err != nil {
		return nil, err
	}
	return &list, nil
}

func Get[T any, PT AddrofObject[T]](ctx context.Context, ctl *Ctl, meta metav1.ObjectMeta) (*T, error) {
	key := client.ObjectKey{
		Name:      meta.Name,
		Namespace: meta.Namespace,
	}
	var obj T
	if err := ctl.client.Get(ctx, key, PT(&obj)); err != nil {
		return nil, err
	}
	return &obj, nil
}

func Create[T any, PT AddrofObject[T]](ctx context.Context, ctl *Ctl, obj T) (*T, error) {
	if err := ctl.client.Create(ctx, PT(&obj)); err != nil {
		return nil, err
	}

	return &obj, nil
}

func Delete[T any, PT AddrofObject[T]](ctx context.Context, ctl *Ctl, meta metav1.ObjectMeta) error {
	obj := PT(new(T))
	obj.SetName(meta.Name)
	obj.SetNamespace(meta.Namespace)

	return ctl.client.Delete(ctx, obj)
}
