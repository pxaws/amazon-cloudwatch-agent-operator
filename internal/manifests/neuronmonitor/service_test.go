// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package neuronmonitor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/aws/amazon-cloudwatch-agent-operator/apis/v1alpha1"
	"github.com/aws/amazon-cloudwatch-agent-operator/internal/config"
	"github.com/aws/amazon-cloudwatch-agent-operator/internal/manifests"
	"github.com/aws/amazon-cloudwatch-agent-operator/internal/manifests/manifestutils"
)

var logger = logf.Log.WithName("unit-tests")

func TestDesiredNeuronService(t *testing.T) {
	t.Run("should return the default service", func(t *testing.T) {
		params := manifests.Params{
			Config: config.Config{},
			Log:    logger,
			NeuronExp: v1alpha1.NeuronMonitor{
				Spec: v1alpha1.NeuronMonitorSpec{},
			},
		}
		trafficPolicy := v1.ServiceInternalTrafficPolicyLocal
		expected := v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:        fmt.Sprintf("%s-service", ComponentNeuronExporter),
				Namespace:   params.NeuronExp.Namespace,
				Labels:      map[string]string{},
				Annotations: map[string]string{},
			},
			Spec: v1.ServiceSpec{
				Type:                  v1.ServiceTypeClusterIP,
				InternalTrafficPolicy: &trafficPolicy,
				Selector:              manifestutils.SelectorLabels(params.NeuronExp.ObjectMeta, ComponentNeuronExporter),
				Ports: []v1.ServicePort{
					{
						Name:       "metrics",
						Port:       8000,
						TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8000},
						Protocol:   v1.ProtocolTCP,
					},
				},
			},
		}

		actual, err := Service(params)
		assert.Nil(t, err)
		assert.Equal(t, expected.ObjectMeta.Name, actual.ObjectMeta.Name)
		assert.Equal(t, expected.Spec.Type, actual.Spec.Type)
		assert.Equal(t, expected.Spec.InternalTrafficPolicy, actual.Spec.InternalTrafficPolicy)
		assert.Equal(t, expected.Spec.Ports, actual.Spec.Ports)
	})

	t.Run("should return a service object with overriden values", func(t *testing.T) {
		params := manifests.Params{
			Config: config.Config{},
			Log:    logger,
			NeuronExp: v1alpha1.NeuronMonitor{
				Spec: v1alpha1.NeuronMonitorSpec{
					Ports: []v1.ServicePort{
						{
							Name: "test",
							Port: 9999,
						},
					},
				},
			},
		}
		trafficPolicy := v1.ServiceInternalTrafficPolicyLocal
		expected := v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:        fmt.Sprintf("%s-service", ComponentNeuronExporter),
				Namespace:   params.NeuronExp.Namespace,
				Labels:      map[string]string{},
				Annotations: map[string]string{},
			},
			Spec: v1.ServiceSpec{
				Type:                  v1.ServiceTypeClusterIP,
				InternalTrafficPolicy: &trafficPolicy,
				Selector:              manifestutils.SelectorLabels(params.NeuronExp.ObjectMeta, ComponentNeuronExporter),
				Ports: []v1.ServicePort{
					{
						Name:       "test",
						Port:       9999,
						TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 9999},
						Protocol:   v1.ProtocolTCP,
					},
				},
			},
		}

		actual, err := Service(params)
		assert.Nil(t, err)
		assert.Equal(t, expected.ObjectMeta.Name, actual.ObjectMeta.Name)
		assert.Equal(t, expected.Spec.Type, actual.Spec.Type)
		assert.Equal(t, expected.Spec.InternalTrafficPolicy, actual.Spec.InternalTrafficPolicy)
		assert.Equal(t, expected.Spec.Ports, actual.Spec.Ports)
	})
}
