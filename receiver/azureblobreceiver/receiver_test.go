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

var (
	logsJSON = []byte(`{"resourceLogs":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"dotnet"}}]},"scopeLogs":[{"scope":{},"logRecords":[{"timeUnixNano":"1643240673066096200","severityText":"Information","body":{"stringValue":"Message Body"},"flags":1,"traceId":"7b20d1349ef9b6d6f9d4d1d4a3ac2e82","spanId":"0c2ad924e1771630"}]}]}]}`)
)

func TestNewReceiver(t *testing.T) {
	receiver, err := getBlobReceiver(t)

	require.NoError(t, err)

	assert.NotNil(t, receiver)
}

func TestConsumeLogsJSON(t *testing.T) {
	receiver, _ := getBlobReceiver(t)

	logsSink := new(consumertest.LogsSink)
	logsConsumer, ok := receiver.(logsDataConsumer)
	require.True(t, ok)

	logsConsumer.setNextLogsConsumer(logsSink)

	err := logsConsumer.consumeLogsJSON(context.Background(), logsJSON)
	require.NoError(t, err)
	assert.Equal(t, 1, logsSink.LogRecordCount())
}

func getBlobReceiver(t *testing.T) (component.Component, error) {
	set := receivertest.NewNopSettings(metadata.Type)

	blobClient := newMockBlobClient()
	blobEventHandler := getBlobEventHandler(t, blobClient)

	getBlobEventHandler(t, blobClient)
	return newReceiver(set, blobEventHandler)
}
