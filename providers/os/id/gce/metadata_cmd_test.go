package gce

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mondoo.com/cnquery/providers/os/connection/mock"
	"go.mondoo.com/cnquery/providers/os/detector"
)

func TestCommandProviderLinux(t *testing.T) {
	conn, err := mock.New("./testdata/metadata_linux.toml")
	require.NoError(t, err)
	platform, ok := detector.DetectOS(conn)
	require.True(t, ok)

	metadata := &commandInstanceMetadata{conn, platform}
	ident, err := metadata.Identify()

	assert.Nil(t, err)
	assert.Equal(t, "//platformid.api.mondoo.app/runtime/gcp/projects/mondoo-dev-262313", ident.ProjectID)
	assert.Equal(t, "//platformid.api.mondoo.app/runtime/gcp/compute/v1/projects/mondoo-dev-262313/zones/us-central1-a/instances/6001244637815193808", ident.InstanceID)
	assert.Equal(t, "//platformid.api.mondoo.app/runtime/gcp/compute/v1/projects/mondoo-dev-262313/zones/us-central1-a/instances/instance-name", ident.PlatformMrn)
}

func TestCommandProviderWindows(t *testing.T) {
	conn, err := mock.New("./testdata/metadata_windows.toml")
	require.NoError(t, err)
	platform, ok := detector.DetectOS(conn)
	require.True(t, ok)

	metadata := &commandInstanceMetadata{conn, platform}
	ident, err := metadata.Identify()

	assert.Nil(t, err)
	assert.Equal(t, "//platformid.api.mondoo.app/runtime/gcp/compute/v1/projects/mondoo-dev-262313/zones/us-central1-a/instances/5275377306317132843", ident.InstanceID)
	assert.Equal(t, "//platformid.api.mondoo.app/runtime/gcp/projects/mondoo-dev-262313", ident.ProjectID)
	assert.Equal(t, "//platformid.api.mondoo.app/runtime/gcp/compute/v1/projects/mondoo-dev-262313/zones/us-central1-a/instances/instance-name", ident.PlatformMrn)
}
