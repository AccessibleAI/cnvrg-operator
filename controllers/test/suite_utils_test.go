package test

import (
	"context"
	"errors"
	"github.com/onsi/gomega"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"
	v1 "k8s.io/api/apps/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	ktypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

// ExpectDryRunCreate runs a dry-run creation.
func ExpectDryRunCreate(obj client.Object) {
	ExpectWithOffset(1, Create(obj, client.DryRunAll)).To(Succeed())
}

// ExpectCreate creates an object as part of an "It".
// It's just a helper to make setup actions easy to find.
func ExpectCreate(obj client.Object) {
	ExpectWithOffset(1, Create(obj)).To(Succeed())
}

// Create wraps the k8s client create.
func Create(obj client.Object, opts ...client.CreateOption) error {
	return tc.client.Create(context.Background(), obj, opts...)
}

// GenName generates a random name from a given base.
func GenName(base string) string {
	const (
		charset   = "abcdefghijklmnopqrstuvwxyz"
		suffixLen = 8
	)
	sb := strings.Builder{}
	sb.WriteString(base)
	for i := 0; i < suffixLen; i++ {
		sb.WriteByte(charset[tc.rnd.Intn(len(charset))])
	}
	return sb.String()
}

// Key builds a namespace-scoped object key.
func Key(namespace string, name string) ktypes.NamespacedName {
	return ktypes.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}
}

// EventuallyGet gets an object with an eventually wrapper.
func EventuallyGet(key ktypes.NamespacedName, obj client.Object, intervals ...interface{}) {
	// this matches the return value, so any non error result is fine
	EventuallyWithOffset(1, GetF(key, obj), intervals...).ShouldNot(BeNil())
}

// EventuallyUpdate helps with updates that tend to race with controller updates.
func EventuallyUpdate(key ktypes.NamespacedName, obj client.Object, modify func(o client.Object), opts ...client.UpdateOption) {
	EventuallyWithOffset(1, func() error {
		if err := Get(key, obj); err != nil {
			return err
		}
		modify(obj)
		return Update(obj, opts...)
	}).Should(Succeed())
}

// GetF returns a function that wraps the k8s client get.
// It helps with the use of async matchers.
func GetF(key ktypes.NamespacedName, obj client.Object) func() (interface{}, error) {
	return func() (interface{}, error) {
		err := Get(key, obj)
		return obj, err
	}
}

func GetTemplateAnnotationsF(key ktypes.NamespacedName, obj v1.Deployment) func() (interface{}, error) {
	return func() (interface{}, error) {
		err := Get(key, &obj)
		return obj.Spec.Template.Annotations, err
	}
}

// Get wraps the k8s client get.
func Get(key ktypes.NamespacedName, obj client.Object) error {
	return tc.client.Get(tc.ctx, key, obj)
}

// MatchSubMap succeeds if the provided map are a subset of actual.
func MatchSubMap(kv map[string]string) types.GomegaMatcher {
	kk := gstruct.Keys{}
	for k, v := range kv {
		kk[k] = gomega.Equal(v)
	}
	return gstruct.MatchKeys(gstruct.IgnoreExtras, kk)
}

// Delete wraps the k8s client delete.
func Delete(obj client.Object, opts ...client.DeleteOption) error {
	return tc.client.Delete(tc.ctx, obj, opts...)
}

// DeleteF returns a function that wraps the k8s client delete.
// It helps with the use of async matchers.
func DeleteF(obj client.Object, opts ...client.DeleteOption) func() error {
	return func() error {
		return Delete(obj, opts...)
	}
}

// ExpectDelete deletes an object.
// For objects without finalizers, an easier option is ExceptDeleteGone.
func ExpectDelete(key ktypes.NamespacedName, obj client.Object) {
	ExpectWithOffset(1, Get(key, obj)).To(Succeed())
	ExpectWithOffset(1, Delete(obj)).To(Succeed())
}

// EventuallyGone expects an object to eventually disappear.
func EventuallyGone(key ktypes.NamespacedName, obj client.Object, intervals ...interface{}) {
	EventuallyWithOffset(1, func() bool {
		err := Get(key, obj)
		return kerrors.IsNotFound(err)
	}, intervals...).Should(BeTrue())
}

// ExpectDeleteGone deletes an object and waits for it to disappear.
// It is like calling ExpectDelete followed by EventuallyGone.
// Do not use this method for objects with finalizers.
func ExpectDeleteGone(key ktypes.NamespacedName, obj client.Object, intervals ...interface{}) {
	ExpectWithOffset(1, Get(key, obj)).To(Succeed())
	ExpectWithOffset(1, Delete(obj)).To(Succeed())
	Eventually(func() bool {
		err := Get(key, obj)
		return kerrors.IsNotFound(err)
	}, intervals...).Should(BeTrue())
}

// EventuallyFinalize waits for an object to have a deletion timestamp and removes all finalizers from it.
// It then waits for the finalized object to disappear.
func EventuallyFinalize(key ktypes.NamespacedName, obj client.Object, intervals ...interface{}) {
	ctx := tc.ctx
	EventuallyWithOffset(1, func() error {
		err := tc.client.Get(ctx, key, obj)
		if err != nil {
			if kerrors.IsNotFound(err) {
				return nil
			}
			return err
		}
		if !obj.GetDeletionTimestamp().IsZero() {
			return errors.New("not deleted")
		}
		obj.SetFinalizers(nil)
		return Update(obj)
	}, intervals...).Should(BeNil())
	EventuallyWithOffset(1, func() bool {
		err := tc.client.Get(ctx, key, obj)
		return kerrors.IsNotFound(err)
	}, intervals...).Should(BeTrue())
}

// Update wraps the k8s client update.
func Update(obj client.Object, opts ...client.UpdateOption) error {
	return tc.client.Update(tc.ctx, obj, opts...)
}

// EventuallyTemplateAnnotationIsPresent helps with the asserting on the annotation filed in deployment template, assert if annotation key is present
func EventuallyTemplateAnnotationIsPresent(key ktypes.NamespacedName, obj v1.Deployment, annotation string, intervals ...interface{}) {
	EventuallyWithOffset(1, GetTemplateAnnotationsF(key, obj), intervals...).Should(gomega.HaveKey(annotation))
}
