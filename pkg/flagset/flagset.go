package flagset

import (
	"fmt"
	"log/slog"

	"github.com/k8s-2025-pschoeppner/ctf/pkg/conditions"
	"github.com/k8s-2025-pschoeppner/ctf/pkg/flags"
	"k8s.io/client-go/kubernetes"
)

const (
	configmapName = "ctf-configmap"
	volumeName    = "ctf-config"
	configmapPath = "/etc/ctf/config"
	secretName    = "ctf-secret"
	secretVolume  = "ctf-secret-config"
	secretPath    = "/etc/ctf/secret"
)

type FlagSet map[string]*flags.Flag

func NewFlagSet(client kubernetes.Interface, logger *slog.Logger) FlagSet {
	return FlagSet{
		"RunningInPod": flags.NewFlag("RunningInPod", client, logger, flags.WithValidators(conditions.PodValidators())),
		"ConfigMap": flags.NewFlag("ConfigMap", client, logger,
			flags.WithValidators(
				conditions.PodValidators(
					conditions.WithConfigmap(volumeName, configmapName),
				),
			),
			flags.WithFulfillers(
				conditions.ReadFromMountedConfigMap("/etc/ctf/config"),
			)),
		"Secret": flags.NewFlag("Secret", client, logger,
			flags.WithValidators(
				conditions.PodValidators(conditions.WithSecret(secretVolume, secretName)),
			),
			flags.WithFulfillers(
				conditions.ReadFromMountedSecret(secretPath),
			)),
		"ServiceAccount": flags.NewFlag("ServiceAccount", client, logger,
			flags.WithValidators(
				conditions.PodValidators(conditions.WithServiceAccount()),
			),
			flags.WithFulfillers(
				conditions.ReadFromExternalConfigMap("ctf-configmap"),
			)),
		"SecurityContext": flags.NewFlag("SecurityContext", client, logger,
			flags.WithValidators(
				conditions.PodValidators(conditions.WithSecurityContext()),
			)),
		"FromTwoPods": flags.NewFlag("FromTwoPods", client, logger,
			flags.WithValidators(
				conditions.WithStatefulConditions(flags.NewStore(), conditions.WithDifferentPodCount(2, 6)),
			),
			flags.WithContinuous()),
		"FromTwoPodsOnce": flags.NewFlag("FromTwoPodsOnce", client, logger,
			flags.WithValidators(
				conditions.WithStatefulConditions(flags.NewStore(), conditions.WithDifferenPodOnce(2)),
			)),
		"FromOnePodTwice": flags.NewFlag("FromOnePodTwice", client, logger,
			flags.WithValidators(
				conditions.WithStatefulConditions(flags.NewStore(), conditions.WithOnePodTwice()),
			)),
		"FromEveryNode": flags.NewFlag("FromEveryNode", client, logger,
			flags.WithValidators(
				conditions.WithStatefulConditions(flags.NewStore(), conditions.WithEveryNode()),
			),
			flags.WithContinuous()),
	}
}

func (fs FlagSet) LogValue() slog.Value {
	keys := make([]slog.Attr, 0, len(fs))
	i := 0
	for k := range fs {
		keys = append(keys, slog.String(fmt.Sprintf("flag %d", i), k))
		i++
	}
	return slog.GroupValue(keys...)
}
