// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package azureblobreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azureblobreceiver"

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/receiverhelper"
	"go.uber.org/zap"
)

type logsDataConsumer interface {
	consumeLogsJSON(ctx context.Context, json []byte) error
	setNextLogsConsumer(nextLogsConsumer consumer.Logs)
}

type blobReceiver struct {
	blobEventHandler blobEventHandler
	logger           *zap.Logger
	logsUnmarshaler  plog.Unmarshaler
	nextLogsConsumer consumer.Logs
	obsrecv          *receiverhelper.ObsReport
}

func (b *blobReceiver) Start(ctx context.Context, _ component.Host) error {
	err := b.blobEventHandler.run(ctx)

	return err
}

func (b *blobReceiver) Shutdown(ctx context.Context) error {
	return b.blobEventHandler.close(ctx)
}

func (b *blobReceiver) setNextLogsConsumer(nextLogsConsumer consumer.Logs) {
	b.nextLogsConsumer = nextLogsConsumer
}

func (b *blobReceiver) consumeLogsJSON(ctx context.Context, buf []byte) error {
	if b.nextLogsConsumer == nil {
		return nil
	}

	logsContext := b.obsrecv.StartLogsOp(ctx)

	var blobBlock map[string]interface{}
	if err := json.Unmarshal(buf, &blobBlock); err != nil {
		return err
	}

	records, ok := blobBlock["records"].([]interface{})
	if !ok {
		return fmt.Errorf("invalid format for records")
	}

	logs := plog.NewLogs()
	resourceLogs := logs.ResourceLogs().AppendEmpty()
	scopeLogs := resourceLogs.ScopeLogs().AppendEmpty()

	for _, record := range records {
		recordMap, ok := record.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid record format")
		}

		recordBytes, err := json.Marshal(record)
		if err != nil {
			return err
		}

		logRecord := scopeLogs.LogRecords().AppendEmpty()
		attributes := logRecord.Attributes()
		attributes.PutStr("log.source", "AzureFlowLogs")
		logRecord.Body().SetStr(string(recordBytes))

		if timeStr, ok := recordMap["time"].(string); ok {
			if timestamp, err := time.Parse(time.RFC3339, timeStr); err == nil {
				logRecord.SetTimestamp(pcommon.NewTimestampFromTime(timestamp))
			} else {
				return fmt.Errorf("invalid time format: %v", err)
			}
		} else {
			return fmt.Errorf("time attribute missing or not a string")
		}
	}

	logRecordCount := logs.LogRecordCount()
	err := b.nextLogsConsumer.ConsumeLogs(logsContext, logs)

	b.obsrecv.EndLogsOp(logsContext, "json", logRecordCount, err)

	return err
}

// Returns a new instance of the log receiver
func newReceiver(set receiver.Settings, blobEventHandler blobEventHandler) (component.Component, error) {
	obsrecv, err := receiverhelper.NewObsReport(receiverhelper.ObsReportSettings{
		ReceiverID:             set.ID,
		Transport:              "event",
		ReceiverCreateSettings: set,
	})
	if err != nil {
		return nil, err
	}

	blobReceiver := &blobReceiver{
		blobEventHandler: blobEventHandler,
		logger:           set.Logger,
		logsUnmarshaler:  &plog.JSONUnmarshaler{},
		obsrecv:          obsrecv,
	}

	blobEventHandler.setLogsDataConsumer(blobReceiver)

	return blobReceiver, nil
}
