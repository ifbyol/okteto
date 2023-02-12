// Copyright 2022 The Okteto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package istio

import (
	"reflect"
	"testing"

	"github.com/okteto/okteto/pkg/model"
	"github.com/stretchr/testify/assert"
	istioNetworkingV1beta1 "istio.io/api/networking/v1beta1"
	istioV1beta1 "istio.io/client-go/pkg/apis/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_translateDeveloperVirtualService(t *testing.T) {
	tests := []struct {
		name     string
		vs       *istioV1beta1.VirtualService
		expected *istioV1beta1.VirtualService
	}{
		{
			name: "match",
			vs: &istioV1beta1.VirtualService{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "service-a",
					Namespace:   "staging",
					Labels:      map[string]string{"l1": "v1"},
					Annotations: map[string]string{"a1": "v1"},
				},
				Spec: istioNetworkingV1beta1.VirtualService{
					Gateways: []string{"ingress-http"},
					Hosts: []string{
						"service-a.staging.svc.cluster.local",
						"service-a.staging.com",
					},
					Http: []*istioNetworkingV1beta1.HTTPRoute{
						{
							Name: "ingress-gateway-http-app-service",
							Match: []*istioNetworkingV1beta1.HTTPMatchRequest{
								{
									Gateways: []string{"ingress-http"},
									Port:     80,
								},
							},
							Route: []*istioNetworkingV1beta1.HTTPRouteDestination{
								{
									Destination: &istioNetworkingV1beta1.Destination{
										Host: "service-a.staging.svc.cluster.local",
										Port: &istioNetworkingV1beta1.PortSelector{
											Number: 80,
										},
										Subset: "stable",
									},
									Weight: 100,
								},
							},
						},
					},
				},
			},
			expected: &istioV1beta1.VirtualService{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "service-a",
					Namespace:   "staging",
					Labels:      map[string]string{"l1": "v1"},
					Annotations: map[string]string{"a1": "v1"},
				},
				Spec: istioNetworkingV1beta1.VirtualService{
					Gateways: []string{"ingress-http"},
					Hosts: []string{
						"service-a.staging.svc.cluster.local",
						"service-a.staging.com",
					},
					Http: []*istioNetworkingV1beta1.HTTPRoute{
						{
							Name: "okteto-divert-cindy-ingress-gateway-http-app-service",
							Match: []*istioNetworkingV1beta1.HTTPMatchRequest{
								{
									Gateways: []string{"ingress-http"},
									Headers: map[string]*istioNetworkingV1beta1.StringMatch{
										"x-okteto-divert": &istioNetworkingV1beta1.StringMatch{
											MatchType: &istioNetworkingV1beta1.StringMatch_Exact{Exact: "cindy"},
										},
									},
									Port: 80,
								},
							},
							Route: []*istioNetworkingV1beta1.HTTPRouteDestination{
								{
									Destination: &istioNetworkingV1beta1.Destination{
										Host: "service-a.cindy.svc.cluster.local",
										Port: &istioNetworkingV1beta1.PortSelector{
											Number: 80,
										},
										Subset: "stable",
									},
									Weight: 100,
								},
							},
						},
						{
							Name: "ingress-gateway-http-app-service",
							Match: []*istioNetworkingV1beta1.HTTPMatchRequest{
								{
									Gateways: []string{"ingress-http"},
									Port:     80,
								},
							},
							Route: []*istioNetworkingV1beta1.HTTPRouteDestination{
								{
									Destination: &istioNetworkingV1beta1.Destination{
										Host: "service-a.staging.svc.cluster.local",
										Port: &istioNetworkingV1beta1.PortSelector{
											Number: 80,
										},
										Subset: "stable",
									},
									Weight: 100,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "no-match",
			vs: &istioV1beta1.VirtualService{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "service-b",
					Namespace:   "staging",
					Labels:      map[string]string{"l1": "v1"},
					Annotations: map[string]string{"a1": "v1"},
				},
				Spec: istioNetworkingV1beta1.VirtualService{
					Gateways: []string{"ingress-http"},
					Hosts: []string{
						"service-b.staging.svc.cluster.local",
						"service-b.staging.com",
					},
					Http: []*istioNetworkingV1beta1.HTTPRoute{
						{
							Name: "ingress-gateway-http-app-service",
							Match: []*istioNetworkingV1beta1.HTTPMatchRequest{
								{
									Gateways: []string{"ingress-http"},
									Port:     80,
								},
							},
							Route: []*istioNetworkingV1beta1.HTTPRouteDestination{
								{
									Destination: &istioNetworkingV1beta1.Destination{
										Host: "service-b.staging.svc.cluster.local",
										Port: &istioNetworkingV1beta1.PortSelector{
											Number: 80,
										},
										Subset: "stable",
									},
									Weight: 100,
								},
							},
						},
					},
				},
			},
			expected: &istioV1beta1.VirtualService{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "service-b",
					Namespace:   "staging",
					Labels:      map[string]string{"l1": "v1"},
					Annotations: map[string]string{"a1": "v1"},
				},
				Spec: istioNetworkingV1beta1.VirtualService{
					Gateways: []string{"ingress-http"},
					Hosts: []string{
						"service-b.staging.svc.cluster.local",
						"service-b.staging.com",
					},
					Http: []*istioNetworkingV1beta1.HTTPRoute{
						{
							Name: "ingress-gateway-http-app-service",
							Match: []*istioNetworkingV1beta1.HTTPMatchRequest{
								{
									Gateways: []string{"ingress-http"},
									Port:     80,
								},
							},
							Route: []*istioNetworkingV1beta1.HTTPRouteDestination{
								{
									Destination: &istioNetworkingV1beta1.Destination{
										Host: "service-b.staging.svc.cluster.local",
										Port: &istioNetworkingV1beta1.PortSelector{
											Number: 80,
										},
										Subset: "stable",
									},
									Weight: 100,
								},
							},
						},
					},
				},
			},
		},
	}

	m := &model.Manifest{
		Name:      "test",
		Namespace: "cindy",
		Deploy: &model.DeployInfo{
			Divert: &model.DivertDeploy{
				Namespace: "staging",
				Service:   "service-a",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := translateDivertVirtualService(m, tt.vs)
			assert.True(t, reflect.DeepEqual(result, tt.expected))
		})
	}
}

func Test_translateDivertVirtualService(t *testing.T) {
	tests := []struct {
		name     string
		vs       *istioV1beta1.VirtualService
		expected *istioV1beta1.VirtualService
	}{
		{
			name: "match",
			vs: &istioV1beta1.VirtualService{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "service-a",
					Namespace:   "staging",
					Labels:      map[string]string{"l1": "v1"},
					Annotations: map[string]string{"a1": "v1"},
				},
				Spec: istioNetworkingV1beta1.VirtualService{
					Gateways: []string{"ingress-http"},
					Hosts: []string{
						"service-a.staging.svc.cluster.local",
						"service-a.staging.com",
					},
					Http: []*istioNetworkingV1beta1.HTTPRoute{
						{
							Name: "ingress-gateway-http-app-service",
							Match: []*istioNetworkingV1beta1.HTTPMatchRequest{
								{
									Gateways: []string{"ingress-http"},
									Port:     80,
								},
							},
							Route: []*istioNetworkingV1beta1.HTTPRouteDestination{
								{
									Destination: &istioNetworkingV1beta1.Destination{
										Host: "service-a.staging.svc.cluster.local",
										Port: &istioNetworkingV1beta1.PortSelector{
											Number: 80,
										},
										Subset: "stable",
									},
									Weight: 100,
								},
							},
						},
					},
				},
			},
			expected: &istioV1beta1.VirtualService{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "service-a",
					Namespace:   "staging",
					Labels:      map[string]string{"l1": "v1"},
					Annotations: map[string]string{"a1": "v1"},
				},
				Spec: istioNetworkingV1beta1.VirtualService{
					Gateways: []string{"ingress-http"},
					Hosts: []string{
						"service-a.staging.svc.cluster.local",
						"service-a.staging.com",
					},
					Http: []*istioNetworkingV1beta1.HTTPRoute{
						{
							Name: "okteto-divert-cindy-ingress-gateway-http-app-service",
							Match: []*istioNetworkingV1beta1.HTTPMatchRequest{
								{
									Gateways: []string{"ingress-http"},
									Headers: map[string]*istioNetworkingV1beta1.StringMatch{
										"x-okteto-divert": &istioNetworkingV1beta1.StringMatch{
											MatchType: &istioNetworkingV1beta1.StringMatch_Exact{Exact: "cindy"},
										},
									},
									Port: 80,
								},
							},
							Route: []*istioNetworkingV1beta1.HTTPRouteDestination{
								{
									Destination: &istioNetworkingV1beta1.Destination{
										Host: "service-a.cindy.svc.cluster.local",
										Port: &istioNetworkingV1beta1.PortSelector{
											Number: 80,
										},
										Subset: "stable",
									},
									Weight: 100,
								},
							},
						},
						{
							Name: "ingress-gateway-http-app-service",
							Match: []*istioNetworkingV1beta1.HTTPMatchRequest{
								{
									Gateways: []string{"ingress-http"},
									Port:     80,
								},
							},
							Route: []*istioNetworkingV1beta1.HTTPRouteDestination{
								{
									Destination: &istioNetworkingV1beta1.Destination{
										Host: "service-a.staging.svc.cluster.local",
										Port: &istioNetworkingV1beta1.PortSelector{
											Number: 80,
										},
										Subset: "stable",
									},
									Weight: 100,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "no-match",
			vs: &istioV1beta1.VirtualService{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "service-b",
					Namespace:   "staging",
					Labels:      map[string]string{"l1": "v1"},
					Annotations: map[string]string{"a1": "v1"},
				},
				Spec: istioNetworkingV1beta1.VirtualService{
					Gateways: []string{"ingress-http"},
					Hosts: []string{
						"service-b.staging.svc.cluster.local",
						"service-b.staging.com",
					},
					Http: []*istioNetworkingV1beta1.HTTPRoute{
						{
							Name: "ingress-gateway-http-app-service",
							Match: []*istioNetworkingV1beta1.HTTPMatchRequest{
								{
									Gateways: []string{"ingress-http"},
									Port:     80,
								},
							},
							Route: []*istioNetworkingV1beta1.HTTPRouteDestination{
								{
									Destination: &istioNetworkingV1beta1.Destination{
										Host: "service-b.staging.svc.cluster.local",
										Port: &istioNetworkingV1beta1.PortSelector{
											Number: 80,
										},
										Subset: "stable",
									},
									Weight: 100,
								},
							},
						},
					},
				},
			},
			expected: &istioV1beta1.VirtualService{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "service-b",
					Namespace:   "staging",
					Labels:      map[string]string{"l1": "v1"},
					Annotations: map[string]string{"a1": "v1"},
				},
				Spec: istioNetworkingV1beta1.VirtualService{
					Gateways: []string{"ingress-http"},
					Hosts: []string{
						"service-b.staging.svc.cluster.local",
						"service-b.staging.com",
					},
					Http: []*istioNetworkingV1beta1.HTTPRoute{
						{
							Name: "ingress-gateway-http-app-service",
							Match: []*istioNetworkingV1beta1.HTTPMatchRequest{
								{
									Gateways: []string{"ingress-http"},
									Port:     80,
								},
							},
							Route: []*istioNetworkingV1beta1.HTTPRouteDestination{
								{
									Destination: &istioNetworkingV1beta1.Destination{
										Host: "service-b.staging.svc.cluster.local",
										Port: &istioNetworkingV1beta1.PortSelector{
											Number: 80,
										},
										Subset: "stable",
									},
									Weight: 100,
								},
							},
						},
					},
				},
			},
		},
	}

	m := &model.Manifest{
		Name:      "test",
		Namespace: "cindy",
		Deploy: &model.DeployInfo{
			Divert: &model.DivertDeploy{
				Namespace: "staging",
				Service:   "service-a",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := translateDivertVirtualService(m, tt.vs)
			assert.True(t, reflect.DeepEqual(result, tt.expected))
		})
	}
}

func Test_restoreDivertVirtualService(t *testing.T) {
	tests := []struct {
		name     string
		vs       *istioV1beta1.VirtualService
		expected *istioV1beta1.VirtualService
	}{
		{
			name: "clean-divert-routes",
			vs: &istioV1beta1.VirtualService{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "service-a",
					Namespace:   "staging",
					Labels:      map[string]string{"l1": "v1"},
					Annotations: map[string]string{"a1": "v1"},
				},
				Spec: istioNetworkingV1beta1.VirtualService{
					Gateways: []string{"ingress-http"},
					Hosts: []string{
						"service-a.staging.svc.cluster.local",
						"service-a.staging.com",
					},
					Http: []*istioNetworkingV1beta1.HTTPRoute{
						{
							Name: "okteto-divert-cindy-ingress-gateway-http-app-service",
							Match: []*istioNetworkingV1beta1.HTTPMatchRequest{
								{
									Gateways: []string{"ingress-http"},
									Headers: map[string]*istioNetworkingV1beta1.StringMatch{
										"x-okteto-divert": &istioNetworkingV1beta1.StringMatch{
											MatchType: &istioNetworkingV1beta1.StringMatch_Exact{Exact: "cindy"},
										},
									},
									Port: 80,
								},
							},
							Route: []*istioNetworkingV1beta1.HTTPRouteDestination{
								{
									Destination: &istioNetworkingV1beta1.Destination{
										Host: "service-a.cindy.svc.cluster.local",
										Port: &istioNetworkingV1beta1.PortSelector{
											Number: 80,
										},
										Subset: "stable",
									},
									Weight: 100,
								},
							},
						},
						{
							Name: "ingress-gateway-http-app-service",
							Match: []*istioNetworkingV1beta1.HTTPMatchRequest{
								{
									Gateways: []string{"ingress-http"},
									Port:     80,
								},
							},
							Route: []*istioNetworkingV1beta1.HTTPRouteDestination{
								{
									Destination: &istioNetworkingV1beta1.Destination{
										Host: "service-a.staging.svc.cluster.local",
										Port: &istioNetworkingV1beta1.PortSelector{
											Number: 80,
										},
										Subset: "stable",
									},
									Weight: 100,
								},
							},
						},
					},
				},
			},
			expected: &istioV1beta1.VirtualService{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "service-a",
					Namespace:   "staging",
					Labels:      map[string]string{"l1": "v1"},
					Annotations: map[string]string{"a1": "v1"},
				},
				Spec: istioNetworkingV1beta1.VirtualService{
					Gateways: []string{"ingress-http"},
					Hosts: []string{
						"service-a.staging.svc.cluster.local",
						"service-a.staging.com",
					},
					Http: []*istioNetworkingV1beta1.HTTPRoute{
						{
							Name: "ingress-gateway-http-app-service",
							Match: []*istioNetworkingV1beta1.HTTPMatchRequest{
								{
									Gateways: []string{"ingress-http"},
									Port:     80,
								},
							},
							Route: []*istioNetworkingV1beta1.HTTPRouteDestination{
								{
									Destination: &istioNetworkingV1beta1.Destination{
										Host: "service-a.staging.svc.cluster.local",
										Port: &istioNetworkingV1beta1.PortSelector{
											Number: 80,
										},
										Subset: "stable",
									},
									Weight: 100,
								},
							},
						},
					},
				},
			},
		},
	}

	m := &model.Manifest{
		Name:      "test",
		Namespace: "cindy",
		Deploy: &model.DeployInfo{
			Divert: &model.DivertDeploy{
				Namespace: "staging",
				Service:   "service-a",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := restoreDivertVirtualService(m, tt.vs)
			assert.True(t, reflect.DeepEqual(result, tt.expected))
		})
	}
}
