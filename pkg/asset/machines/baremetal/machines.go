// Package baremetal generates Machine objects for bare metal.
package baremetal

import (
	"fmt"

	// FIXME: baremetal
	libvirtprovider "github.com/openshift/cluster-api-provider-libvirt/pkg/apis/libvirtproviderconfig/v1alpha1"

	machineapi "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/metalkube/kni-installer/pkg/types"
	"github.com/metalkube/kni-installer/pkg/types/baremetal"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, role, userDataSecret string) ([]machineapi.Machine, error) {
	if configPlatform := config.Platform.Name(); configPlatform != baremetal.Name {
		return nil, fmt.Errorf("non bare metal configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != baremetal.Name {
		return nil, fmt.Errorf("non bare metal machine-pool: %q", poolPlatform)
	}
	clustername := config.ObjectMeta.Name
	platform := config.Platform.BareMetal

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	provider := provider(clustername, config.Networking.MachineCIDR.String(), platform, userDataSecret)
	var machines []machineapi.Machine
	for idx := int64(0); idx < total; idx++ {
		machine := machineapi.Machine{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "machine.openshift.io/v1beta1",
				Kind:       "Machine",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "openshift-machine-api",
				Name:      fmt.Sprintf("%s-%s-%d", clustername, pool.Name, idx),
				Labels: map[string]string{
					"sigs.k8s.io/cluster-api-cluster":      clustername,
					"sigs.k8s.io/cluster-api-machine-role": role,
					"sigs.k8s.io/cluster-api-machine-type": role,
				},
			},
			Spec: machineapi.MachineSpec{
				ProviderSpec: machineapi.ProviderSpec{
					Value: &runtime.RawExtension{Object: provider},
				},
				// we don't need to set Versions, because we control those via cluster operators.
			},
		}
		machines = append(machines, machine)
	}

	return machines, nil
}

func provider(clusterName string, networkInterfaceAddress string, platform *baremetal.Platform, userDataSecret string) *libvirtprovider.LibvirtMachineProviderConfig {
	// FIXME: baremetal
	return &libvirtprovider.LibvirtMachineProviderConfig{}
}
