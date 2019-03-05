package exec

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	core "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type Options struct {
	core.PodExecOptions
	remotecommand.StreamOptions
}

func Container(container string) func(*Options) {
	return func(opts *Options) {
		opts.Container = container
	}
}

func Command(cmd ...string) func(*Options) {
	return func(opts *Options) {
		opts.Command = cmd
	}
}

func Input(in string) func(*Options) {
	return func(opts *Options) {
		opts.PodExecOptions.Stdin = true
		opts.StreamOptions.Stdin = strings.NewReader(in)
	}
}

func TTY(enable bool) func(*Options) {
	return func(opts *Options) {
		opts.PodExecOptions.TTY = enable
	}
}

func ExecIntoPod(config *rest.Config, pod *core.Pod, options ...func(*Options)) (string, error) {
	var (
		execOut bytes.Buffer
		execErr bytes.Buffer
		opts    = &Options{
			PodExecOptions: core.PodExecOptions{
				Container: pod.Spec.Containers[0].Name,
				Stdout:    true,
				Stderr:    true,
			},
			StreamOptions: remotecommand.StreamOptions{
				Stdout: &execOut,
				Stderr: &execErr,
			},
		}
	)
	for _, option := range options {
		option(opts)
	}

	kc := kubernetes.NewForConfigOrDie(config)

	req := kc.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec")
	req.VersionedParams(&opts.PodExecOptions, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, http.MethodPost, req.URL())
	if err != nil {
		return "", fmt.Errorf("failed to init executor: %v", err)
	}

	err = exec.Stream(opts.StreamOptions)

	if err != nil {
		return "", fmt.Errorf("could not execute: %v", err)
	}

	if execErr.Len() > 0 {
		return "", fmt.Errorf("stderr: %v", execErr.String())
	}

	return execOut.String(), nil
}
