package validation

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/metalkube/kni-installer/pkg/types"
	"github.com/metalkube/kni-installer/pkg/types/aws"
	awsvalidation "github.com/metalkube/kni-installer/pkg/types/aws/validation"
	"github.com/metalkube/kni-installer/pkg/types/baremetal"
	baremetalvalidation "github.com/metalkube/kni-installer/pkg/types/baremetal/validation"
	"github.com/metalkube/kni-installer/pkg/types/libvirt"
	libvirtvalidation "github.com/metalkube/kni-installer/pkg/types/libvirt/validation"
	"github.com/metalkube/kni-installer/pkg/types/openstack"
	openstackvalidation "github.com/metalkube/kni-installer/pkg/types/openstack/validation"
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(p *types.MachinePool, fldPath *field.Path, platform string) field.ErrorList {
	allErrs := field.ErrorList{}
	if p.Replicas != nil {
		if *p.Replicas < 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("replicas"), p.Replicas, "number of replicas must not be negative"))
		}
	} else {
		allErrs = append(allErrs, field.Required(fldPath.Child("replicas"), "replicas is required"))
	}
	allErrs = append(allErrs, validateMachinePoolPlatform(&p.Platform, fldPath.Child("platform"), platform)...)
	return allErrs
}

func validateMachinePoolPlatform(p *types.MachinePoolPlatform, fldPath *field.Path, platform string) field.ErrorList {
	allErrs := field.ErrorList{}
	validate := func(n string, value interface{}, validation func(*field.Path) field.ErrorList) {
		f := fldPath.Child(n)
		if platform == n {
			allErrs = append(allErrs, validation(f)...)
		} else {
			allErrs = append(allErrs, field.Invalid(f, value, fmt.Sprintf("cannot specify %q for machine pool when cluster is using %q", n, platform)))
		}
	}
	if p.AWS != nil {
		validate(aws.Name, p.AWS, func(f *field.Path) field.ErrorList { return awsvalidation.ValidateMachinePool(p.AWS, f) })
	}
	if p.Libvirt != nil {
		validate(libvirt.Name, p.Libvirt, func(f *field.Path) field.ErrorList { return libvirtvalidation.ValidateMachinePool(p.Libvirt, f) })
	}
	if p.OpenStack != nil {
		validate(openstack.Name, p.OpenStack, func(f *field.Path) field.ErrorList { return openstackvalidation.ValidateMachinePool(p.OpenStack, f) })
	}
	if p.BareMetal != nil {
		validate(baremetal.Name, p.BareMetal, func(f *field.Path) field.ErrorList { return baremetalvalidation.ValidateMachinePool(p.BareMetal, f) })
	}
	return allErrs
}
