/*
Copyright 2018 The Kubernetes Authors.

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

package main

import (
	"context"
	"encoding/json"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/mutate-v1-secret,mutating=true,failurePolicy=fail,groups="",resources=secrets,verbs=create;update,versions=v1,name=msecret.conjur.org

// secretInjector injects Conjur secrets into Kubernetes secrets.
type secretInjector struct {
	Client  client.Client
	decoder *admission.Decoder
}

// secretInjector injects Conjur secrets into every incoming secret in an appropriately labelled namespace.
func (a *secretInjector) Handle(ctx context.Context, req admission.Request) admission.Response {
	secret := &corev1.Secret{}

	err := a.decoder.Decode(req, secret)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	if secret.Annotations == nil {
		secret.Annotations = map[string]string{}
	}
	secret.Annotations["example-mutating-admission-webhook"] = "foo"

	marshaledPod, err := json.Marshal(secret)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

// podAnnotator implements admission.DecoderInjector.
// A decoder will be automatically injected.

// InjectDecoder injects the decoder.
func (a *secretInjector) InjectDecoder(d *admission.Decoder) error {
	a.decoder = d
	return nil
}
