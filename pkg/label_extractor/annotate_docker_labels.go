package label_extractor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/heroku/docker-registry-client/registry"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// getAllSecrets() takes imagePullSecrets and return the list of secret names as an array of
// string
func getAllSecrets(imagePullSecrets []corev1.LocalObjectReference) []string {
	secretNames := []string{}
	for _, secretName := range imagePullSecrets {
		secretNames = append(secretNames, secretName.Name)
	}

	return secretNames
}

// This method add <prefix> as prefix to the keys of <labels>
func addPrefixToLabels(labels map[string]string, prefix string) {
	for key, val := range labels {
		delete(labels, key)
		labels[prefix+key] = val
	}
}

// This method removes the annotations which have key with <prefix> as prefix
func removeOldAnnotations(annotations map[string]string, prefix string) {
	if annotations == nil {
		return
	}

	for key := range annotations {
		if strings.HasPrefix(key, prefix) {
			delete(annotations, key)
		}
	}
}

// This method takes namespace_name <namespace> of provided secrets <secretNames> and a docker image
// name <image>. For each secret it reads the config data of secret and store it to registrySecrets
// (map[string]RegistrySecret) where the api url is the key and value is the credentials. Then it tries
// to extract labels of the <image> for all secrets' content. If found then returns labels otherwise
// returns corresponding error. If <image> is not found with the secret info, then it tries with the
// public docker url="https://registry-1.docker.io/"
func (l *ExtractDockerLabel) GetLabels(namespace, repoName, tag string, secretNames []string) (map[string]string, error) {
	var err error
	image := repoName + ":" + tag

	if l.twoQCache.Contains(image) {
		val, _ := l.twoQCache.Get(image)
		return val.(map[string]string), nil
	}

	for _, item := range secretNames {
		secret, err := l.kubeClient.CoreV1().Secrets(namespace).Get(item, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("couldn't get secret(%s): %v", item, err)
		}

		configData := []byte{}
		for _, val := range secret.Data {
			configData = append(configData, val...)
			break
		}
		//log.Infoln("config.json =", string(configData))

		var registrySecrets map[string]RegistrySecret
		err = json.NewDecoder(bytes.NewReader(configData)).Decode(&registrySecrets)
		if err != nil {
			return nil, fmt.Errorf("couldn't decode the configData for secret(%s): %v", item, err)
		}

		for key, val := range registrySecrets {
			labels, extracErr, errStatus := l.ExtractLabelsForThisCred(key, val.Username, val.Password, repoName, tag)

			if errStatus != 0 {
				err = fmt.Errorf("%v\n%v", err, extracErr)
				continue
			}

			l.twoQCache.Add(image, labels)

			return labels, err
		}
	}

	url := "https://registry-1.docker.io/"
	username := "" // anonymous
	pass := ""     // anonymous

	labels, extractErr, errStatus := l.ExtractLabelsForThisCred(url, username, pass, repoName, tag)
	if errStatus != 0 {
		err = fmt.Errorf("%v\n%v", err, extractErr)
		return nil, fmt.Errorf("couldn't find image(%s:%s): %v", repoName, tag, err)
	}

	l.twoQCache.Add(image, labels)

	return labels, err
}

// This method returns the labels of docker image. The image name is <reopName/tag> and the api url
// is <url>. The essential credentials are <username> and <pass>. If image is found it returns tuple
// {labels, err=nil, status=0}, otherwise it returns tuple {label=nil, err, status}
func (l *ExtractDockerLabel) ExtractLabelsForThisCred(
	url, username, pass string,
	repoName, tag string) (map[string]string, error, int) {

	hub := &registry.Registry{
		URL: url,
		Client: &http.Client{
			Transport: registry.WrapTransport(http.DefaultTransport, url, username, pass),
		},
		Logf: registry.Quiet,
	}

	manifest, err := hub.ManifestV2(repoName, tag)
	if err != nil {
		return nil,
			fmt.Errorf("couldn't get the manifest for credential(url->%s, username->%s, pass->%s): %v",
				url, username, pass, err),
			1
	}

	reader, err := hub.DownloadLayer(repoName, manifest.Config.Digest)
	if err != nil {
		return nil,
			fmt.Errorf("couldn't get encoded imageInspect for credential(url->%s, username->%s, pass->%s): %v",
				url, username, pass, err),
			2
	}

	var cfg types.ImageInspect
	defer reader.Close()
	err = json.NewDecoder(reader).Decode(&cfg)
	if err != nil {
		return nil,
			fmt.Errorf("couldn't get decode imageInspect for credential(url->%s, username->%s, pass->%s): %v",
				url, username, pass, err),
			3
	}

	return cfg.Config.Labels, nil, 0
}
