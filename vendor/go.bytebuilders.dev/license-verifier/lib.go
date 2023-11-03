/*
Copyright AppsCode Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package verifier

import (
	"crypto/x509"
	"fmt"
	"strings"

	"go.bytebuilders.dev/license-verifier/apis/licenses/v1alpha1"
	"go.bytebuilders.dev/license-verifier/info"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
)

type Options struct {
	ClusterUID string `json:"clusterUID"`
	Features   string `json:"features"`
	CACert     []byte `json:"caCert,omitempty"`
	License    []byte `json:"license"`
}

type ParserOptions struct {
	ClusterUID string
	CACert     *x509.Certificate
	License    []byte
}

type VerifyOptions struct {
	ParserOptions
	Features string
}

func ParseLicense(opts ParserOptions) (v1alpha1.License, error) {
	cert, err := info.ParseCertificate(opts.License)
	if err != nil {
		return BadLicense(err)
	}

	roots := x509.NewCertPool()
	roots.AddCert(opts.CACert)

	crtopts := x509.VerifyOptions{
		DNSName: opts.ClusterUID,
		Roots:   roots,
		KeyUsages: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
		},
	}

	// wildcard certificate
	if strings.HasPrefix(cert.Subject.CommonName, "*.") {
		if len(opts.CACert.Subject.Organization) > 0 {
			crtopts.DNSName = "*." + opts.CACert.Subject.Organization[0]
		}
	}

	license := v1alpha1.License{
		TypeMeta: metav1.TypeMeta{
			APIVersion: v1alpha1.SchemeGroupVersion.String(),
			Kind:       "License",
		},
		Data:      opts.License,
		Issuer:    "byte.builders",
		Clusters:  cert.DNSNames,
		NotBefore: &metav1.Time{Time: cert.NotBefore},
		NotAfter:  &metav1.Time{Time: cert.NotAfter},
		ID:        cert.SerialNumber.String(),
		Features:  cert.Subject.Organization,
	}
	if len(cert.Subject.OrganizationalUnit) > 0 {
		license.PlanName = cert.Subject.OrganizationalUnit[0]
	} else {
		// old certificate, so plan name auto detected from feature
		// ref: https://github.com/appscode/offline-license-server/blob/v0.0.20/pkg/server/constants.go#L50-L59
		features := sets.NewString(cert.Subject.Organization...)
		if features.Has("kubedb-enterprise") {
			license.PlanName = "kubedb-enterprise"
		} else if features.Has("kubedb-community") {
			license.PlanName = "kubedb-community"
		} else if features.Has("stash-enterprise") {
			license.PlanName = "stash-enterprise"
		} else if features.Has("stash-community") {
			license.PlanName = "stash-community"
		}
	}
	if len(cert.Subject.Country) > 0 {
		license.ProductLine = cert.Subject.Country[0]
	}
	if len(cert.Subject.Province) > 0 {
		license.TierName = cert.Subject.Province[0]
	}
	if license.ProductLine == "" || license.TierName == "" {
		parts := strings.SplitN(license.PlanName, "-", 2)
		if len(parts) > 0 {
			license.ProductLine = parts[0]
		}
		if len(parts) > 1 {
			license.TierName = parts[1]
		}
	}
	license.FeatureFlags = map[string]string{}
	for _, ff := range cert.Subject.Locality {
		parts := strings.SplitN(ff, "=", 2)
		if len(parts) == 2 {
			license.FeatureFlags[parts[0]] = parts[1]
		}
	}

	var user *v1alpha1.User
	for _, e := range cert.EmailAddresses {
		parts := strings.FieldsFunc(e, func(r rune) bool {
			return r == '<' || r == '>'
		})
		if len(parts) == 0 {
			continue
		}

		if len(parts) == 1 {
			email := strings.TrimSpace(parts[0])
			if user == nil {
				user = &v1alpha1.User{
					Name:  "",
					Email: email,
				}
			} else if user.Email != email {
				return BadLicense(fmt.Errorf("license issued to multiple emails %s", strings.Join(cert.EmailAddresses, ";")))
			}
		} else { // == 2
			email := strings.TrimSpace(parts[1])
			if user == nil {
				user = &v1alpha1.User{
					Name:  strings.TrimSpace(parts[0]),
					Email: email,
				}
			} else if user.Email != email {
				return BadLicense(fmt.Errorf("license issued to multiple emails %s", strings.Join(cert.EmailAddresses, ";")))
			}
		}
	}
	license.User = user

	// ref: https://github.com/appscode/gitea/blob/master/models/stripe_license.go#L117-L126
	if _, err := cert.Verify(crtopts); err != nil {
		e2 := errors.Wrap(err, "failed to verify certificate")
		license.Status = v1alpha1.LicenseInvalid
		license.Reason = e2.Error()
		return license, e2
	}
	license.Status = v1alpha1.LicenseActive
	return license, nil
}

func CheckLicense(opts VerifyOptions) (v1alpha1.License, error) {
	license, err := ParseLicense(opts.ParserOptions)
	if err != nil {
		return license, err
	}
	if !sets.NewString(license.Features...).HasAny(info.ParseFeatures(opts.Features)...) {
		e2 := fmt.Errorf("license was not issued for %s", opts.Features)
		license.Status = v1alpha1.LicenseInvalid
		license.Reason = e2.Error()
		return license, e2
	}
	license.Status = v1alpha1.LicenseActive
	return license, nil
}

func VerifyLicense(opts Options) (v1alpha1.License, error) {
	caCert, err := info.ParseCertificate(opts.CACert)
	if err != nil {
		return BadLicense(err)
	}

	return CheckLicense(VerifyOptions{
		ParserOptions: ParserOptions{
			ClusterUID: opts.ClusterUID,
			CACert:     caCert,
			License:    opts.License,
		},
		Features: opts.Features,
	})
}

func BadLicense(err error) (v1alpha1.License, error) {
	if err == nil {
		// This should never happen
		panic(err)
	}
	return v1alpha1.License{
		TypeMeta: metav1.TypeMeta{
			APIVersion: v1alpha1.SchemeGroupVersion.String(),
			Kind:       "License",
		},
		Status: v1alpha1.LicenseUnknown,
		Reason: err.Error(),
	}, err
}
