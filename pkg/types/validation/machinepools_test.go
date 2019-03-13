package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/pointer"

	"github.com/metalkube/kni-installer/pkg/types"
	"github.com/metalkube/kni-installer/pkg/types/aws"
	"github.com/metalkube/kni-installer/pkg/types/libvirt"
	"github.com/metalkube/kni-installer/pkg/types/openstack"
)

func validMachinePool() *types.MachinePool {
	return &types.MachinePool{
		Name:     "test-pool",
		Replicas: pointer.Int64Ptr(1),
	}
}

func TestValidateMachinePool(t *testing.T) {
	cases := []struct {
		name     string
		pool     *types.MachinePool
		platform string
		valid    bool
	}{
		{
			name:     "minimal",
			pool:     validMachinePool(),
			platform: "aws",
			valid:    true,
		},
		{
			name: "missing replicas",
			pool: func() *types.MachinePool {
				p := validMachinePool()
				p.Replicas = nil
				return p
			}(),
			platform: "aws",
			valid:    false,
		},
		{
			name: "invalid replicas",
			pool: func() *types.MachinePool {
				p := validMachinePool()
				p.Replicas = pointer.Int64Ptr(-1)
				return p
			}(),
			platform: "aws",
			valid:    false,
		},
		{
			name: "valid aws",
			pool: func() *types.MachinePool {
				p := validMachinePool()
				p.Platform = types.MachinePoolPlatform{
					AWS: &aws.MachinePool{},
				}
				return p
			}(),
			platform: "aws",
			valid:    true,
		},
		{
			name: "invalid aws",
			pool: func() *types.MachinePool {
				p := validMachinePool()
				p.Platform = types.MachinePoolPlatform{
					AWS: &aws.MachinePool{
						EC2RootVolume: aws.EC2RootVolume{
							IOPS: -10,
						},
					},
				}
				return p
			}(),
			platform: "aws",
			valid:    false,
		},
		{
			name: "valid libvirt",
			pool: func() *types.MachinePool {
				p := validMachinePool()
				p.Platform = types.MachinePoolPlatform{
					Libvirt: &libvirt.MachinePool{},
				}
				return p
			}(),
			platform: "libvirt",
			valid:    true,
		},
		{
			name: "valid openstack",
			pool: func() *types.MachinePool {
				p := validMachinePool()
				p.Platform = types.MachinePoolPlatform{
					OpenStack: &openstack.MachinePool{},
				}
				return p
			}(),
			platform: "openstack",
			valid:    true,
		},
		{
			name: "mis-matched platform",
			pool: func() *types.MachinePool {
				p := validMachinePool()
				p.Platform = types.MachinePoolPlatform{
					AWS: &aws.MachinePool{},
				}
				return p
			}(),
			platform: "libvirt",
			valid:    false,
		},
		{
			name: "multiple platforms",
			pool: func() *types.MachinePool {
				p := validMachinePool()
				p.Platform = types.MachinePoolPlatform{
					AWS:     &aws.MachinePool{},
					Libvirt: &libvirt.MachinePool{},
				}
				return p
			}(),
			platform: "aws",
			valid:    false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateMachinePool(tc.pool, field.NewPath("test-path"), tc.platform).ToAggregate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
