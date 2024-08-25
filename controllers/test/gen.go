package test

import (
	mlopsv1 "github.com/AccessibleAI/cnvrg-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func genApp(ns string, name string) mlopsv1.CnvrgApp {
	testSpec := mlopsv1.DefaultCnvrgAppSpec()
	testSpec.ClusterDomain = "test.local"
	return mlopsv1.CnvrgApp{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CnvrgApp",
			APIVersion: "mlops.cnvrg.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: testSpec,
	}
}

func genNS(name string) corev1.Namespace {
	testSpec := mlopsv1.DefaultCnvrgAppSpec()
	testSpec.ClusterDomain = "test.local"
	return corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind: "Namespace",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}
