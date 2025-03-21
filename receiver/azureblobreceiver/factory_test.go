// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package azureblobreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azureblobreceiver"

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/receiver/receivertest"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azureblobreceiver/internal/metadata"
)

func TestNewFactory(t *testing.T) {
	f := NewFactory()

	assert.NotNil(t, f)
}

func TestCreateLogs(t *testing.T) {
	f := NewFactory()
	ctx := context.Background()
	params := receivertest.NewNopSettings(metadata.Type)
	receiver, err := f.CreateLogs(ctx, params, getConfig(), consumertest.NewNop())

	require.NoError(t, err)
	assert.NotNil(t, receiver)
}

func getConfig() component.Config {
	return &Config{
		Authentication:   "connection_string",
		ConnectionString: goodConnectionString,
		Logs:             LogsConfig{ContainerName: logsContainerName},
	}
}
