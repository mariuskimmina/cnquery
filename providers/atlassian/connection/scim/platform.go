package scim

import (
	"go.mondoo.com/cnquery/v9/providers-sdk/v1/inventory"
)

func (a *ScimConnection) PlatformInfo() *inventory.Platform {
	return GetPlatformForObject("atlassian")
}

func GetPlatformForObject(platformName string) *inventory.Platform {
	if platformName != "atlassian" && platformName != "" {
		return &inventory.Platform{
			Name:    platformName,
			Title:   "atlassian",
			Kind:    "api",
			Runtime: "atlassian",
		}
	}
	return &inventory.Platform{
		Name:    "atlassian",
		Title:   "atlassian",
		Kind:    "api",
		Runtime: "atlassian",
	}
}

func (a *ScimConnection) PlatformID() string {
	return "//platformid.api.mondoo.app/runtime/atlassian/scim/" + a.Directory()
}
