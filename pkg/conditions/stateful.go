package conditions

import (
	"context"
	"fmt"
	"reflect"

	"github.com/k8s-2025-pschoeppner/ctf/pkg/flags"
	"github.com/k8s-2025-pschoeppner/ctf/pkg/k8s"
	"github.com/k8s-2025-pschoeppner/ctf/pkg/types"
	"k8s.io/client-go/kubernetes"
)

type (
	StatefulValidator func(ctx context.Context, store *flags.Store, r types.Request, client kubernetes.Interface) error
)

func WithStatefulConditions(store *flags.Store, validators ...StatefulValidator) flags.Validator {
	return func(ctx context.Context, r types.Request, i kubernetes.Interface) error {
		for _, validator := range validators {
			if err := validator(ctx, store, r, i); err != nil {
				return err
			}
		}
		return nil
	}
}

func WithDifferentPodCount(minPods, total int) StatefulValidator {
	return func(ctx context.Context, store *flags.Store, r types.Request, client kubernetes.Interface) error {
		if store == nil {
			return fmt.Errorf("failed to write to nil store... again.")
		}
		podName := r.PodName
		requestID := r.ID
		pod, err := k8s.GetPod(ctx, client, podName, r.PodNamespace)
		if err != nil {
			return fmt.Errorf("failed to get pod %s/%s", podName, r.PodNamespace)
		}
		podTemplateHash := pod.Labels["pod-template-hash"]
		mapKey := fmt.Sprintf("%s-%s", requestID, podTemplateHash)
		podMap, exists := store.Get(mapKey)
		if !exists {
			m := make(map[string]int)
			m[podName] = 1
			store.Set(mapKey, m)
			return types.ErrStatefulValidatorContinue
		}
		m, ok := podMap.(map[string]int)
		if !ok {
			return fmt.Errorf("failed to cast %q to map[string]int", reflect.TypeOf(podMap))
		}
		m[podName]++
		store.Set(mapKey, m)
		t := 0
		for _, v := range m {
			t += v
		}
		if t >= total && len(m) >= minPods {
			return nil
		}
		return types.ErrStatefulValidatorContinue
	}
}

func WithDifferenPodOnce(minPods int) StatefulValidator {
	return func(ctx context.Context, store *flags.Store, r types.Request, client kubernetes.Interface) error {
		if store == nil {
			return fmt.Errorf("failed to write to nil store... again.")
		}
		podName := r.PodName
		requestID := r.ID
		pod, err := k8s.GetPod(ctx, client, podName, r.PodNamespace)
		if err != nil {
			return fmt.Errorf("failed to get pod %s/%s", podName, r.PodNamespace)
		}
		podTemplateHash := pod.Labels["pod-template-hash"]
		mapKey := fmt.Sprintf("%s-%s", requestID, podTemplateHash)
		podMap, exists := store.Get(mapKey)
		if !exists {
			m := make(map[string]int)
			m[podName] = 1
			store.Set(mapKey, m)
			return types.ErrStatefulValidatorContinue
		}
		m, ok := podMap.(map[string]int)
		if !ok {
			return fmt.Errorf("failed to cast %q to map[string]int", reflect.TypeOf(podMap))
		}
		m[podName]++
		store.Set(mapKey, m)
		once := true
		for _, v := range m {
			if v > 1 {
				once = false
				break
			}
		}
		if once && len(m) >= minPods {
			return nil
		}
		return types.ErrStatefulValidatorContinue
	}
}

func WithOnePodTwice() StatefulValidator {
	return func(ctx context.Context, store *flags.Store, r types.Request, client kubernetes.Interface) error {
		if store == nil {
			return fmt.Errorf("failed to write to nil store... again.")
		}
		podName := r.PodName
		requestID := r.ID
		pod, err := k8s.GetPod(ctx, client, podName, r.PodNamespace)
		if err != nil {
			return fmt.Errorf("failed to get pod %s/%s", podName, r.PodNamespace)
		}
		podTemplateHash := pod.Labels["pod-template-hash"]
		mapKey := fmt.Sprintf("%s-%s", requestID, podTemplateHash)
		podCounter, exists := store.Get(mapKey)
		if !exists {
			c := 1
			store.Set(mapKey, c)
			return types.ErrStatefulValidatorContinue
		}
		c, ok := podCounter.(int)
		if !ok {
			return fmt.Errorf("failed to cast %q to int", reflect.TypeOf(podCounter))
		}
		c++
		store.Set(mapKey, c)
		if c == 2 {
			return nil
		}
		return types.ErrStatefulValidatorContinue
	}
}

func WithEveryNode() StatefulValidator {
	return func(ctx context.Context, store *flags.Store, r types.Request, client kubernetes.Interface) error {
		if store == nil {
			return fmt.Errorf("failed to write to nil store... again.")
		}
		podName := r.PodName
		requestID := r.ID
		pod, err := k8s.GetPod(ctx, client, podName, r.PodNamespace)
		if err != nil {
			return fmt.Errorf("failed to get pod %s/%s", podName, r.PodNamespace)
		}
		node := pod.Spec.NodeName
		nodeMap, exists := store.Get(requestID)
		if !exists {
			m := make(map[string]int)
			m[node] = 1
			store.Set(requestID, m)
			return types.ErrStatefulValidatorContinue
		}
		m, ok := nodeMap.(map[string]int)
		if !ok {
			return fmt.Errorf("failed to cast %q to map[string]int", reflect.TypeOf(nodeMap))
		}
		m[node]++
		store.Set(requestID, m)
		nodeList, err := k8s.GetNodes(ctx, client)
		if err != nil {
			return fmt.Errorf("failed to get nodes")
		}
		if len(m) == len(nodeList.Items) {
			return nil
		}
		return types.ErrStatefulValidatorContinue
	}
}
