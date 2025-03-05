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
