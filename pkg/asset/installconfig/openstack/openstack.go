// Package openstack collects OpenStack-specific configuration.
package openstack

import (
	"sort"
	"strings"

	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/metalkube/kni-installer/pkg/types/openstack"
	openstackvalidation "github.com/metalkube/kni-installer/pkg/types/openstack/validation"
)

// Platform collects OpenStack-specific configuration.
func Platform() (*openstack.Platform, error) {
	validValuesFetcher := openstackvalidation.NewValidValuesFetcher()

	cloudNames, err := validValuesFetcher.GetCloudNames()
	if err != nil {
		return nil, err
	}
	// Sort cloudNames so we can use sort.SearchStrings
	sort.Strings(cloudNames)
	var cloud string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Cloud",
				Help:    "The OpenStack cloud name from clouds.yaml.",
				Options: cloudNames,
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				value := ans.(string)
				i := sort.SearchStrings(cloudNames, value)
				if i == len(cloudNames) || cloudNames[i] != value {
					return errors.Errorf("invalid cloud name %q, should be one of %+v", value, strings.Join(cloudNames, ", "))
				}
				return nil
			}),
		},
	}, &cloud)
	if err != nil {
		return nil, err
	}

	regionNames, err := validValuesFetcher.GetRegionNames(cloud)
	if err != nil {
		return nil, err
	}
	sort.Strings(regionNames)
	var region string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Region",
				Help:    "The OpenStack region to be used for installation.",
				Default: "regionOne",
				Options: regionNames,
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				value := ans.(string)
				i := sort.SearchStrings(regionNames, value)
				if i == len(regionNames) || regionNames[i] != value {
					return errors.Errorf("invalid region name %q, should be one of %+v", value, strings.Join(regionNames, ", "))
				}
				return nil
			}),
		},
	}, &region)
	if err != nil {
		return nil, err
	}

	networkNames, err := validValuesFetcher.GetNetworkNames(cloud)
	if err != nil {
		return nil, err
	}
	sort.Strings(networkNames)
	var extNet string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "ExternalNetwork",
				Help:    "The OpenStack external network name to be used for installation.",
				Options: networkNames,
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				value := ans.(string)
				i := sort.SearchStrings(networkNames, value)
				if i == len(networkNames) || networkNames[i] != value {
					return errors.Errorf("invalid network name %q, should be one of %+v", value, strings.Join(networkNames, ", "))
				}
				return nil
			}),
		},
	}, &extNet)
	if err != nil {
		return nil, err
	}

	flavorNames, err := validValuesFetcher.GetFlavorNames(cloud)
	if err != nil {
		return nil, err
	}
	sort.Strings(flavorNames)
	var flavor string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "FlavorName",
				Help:    "The OpenStack compute flavor to use for servers. A flavor with at least 4 GB RAM is recommended.",
				Options: flavorNames,
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				value := ans.(string)
				i := sort.SearchStrings(flavorNames, value)
				if i == len(flavorNames) || flavorNames[i] != value {
					return errors.Errorf("invalid flavor name %q, should be one of %+v", value, strings.Join(flavorNames, ", "))
				}
				return nil
			}),
		},
	}, &flavor)
	if err != nil {
		return nil, err
	}

	netExts, err := validValuesFetcher.GetNetworkExtensionsAliases(cloud)
	if err != nil {
		return nil, err
	}
	sort.Strings(netExts)
	i := sort.SearchStrings(netExts, "trunk")
	var trunkSupport string
	if i == len(netExts) || netExts[i] != "trunk" {
		trunkSupport = "0"
	} else {
		trunkSupport = "1"
	}

	return &openstack.Platform{
		Region:          region,
		Cloud:           cloud,
		ExternalNetwork: extNet,
		FlavorName:      flavor,
		TrunkSupport:    trunkSupport,
	}, nil
}
