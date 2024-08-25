package metastorageprovider

import (
	"context"
	"github.com/AccessibleAI/cnvrg-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EventLogger is a struct that records events, sets messages on objects, and writes logs.
type EventLogger struct {
	client   client.Client
	logger   logr.Logger
	recorder record.EventRecorder
	message  string
	err      error
	Status   v1alpha1.MetaStorageProvisionerStatus
	Object   *v1alpha1.MetaStorageProvisioner
}

// WithMessage sets the message for the EventLogger and returns the updated EventLogger.
func (e EventLogger) WithMessage(message string) EventLogger {
	e.message = message
	return e
}

// WithError sets the error for the EventLogger and returns the updated EventLogger.
func (e EventLogger) WithError(err error) EventLogger {
	e.err = err
	return e
}

// WithStatus sets the status for the EventLogger and returns the updated EventLogger.
func (e EventLogger) WithStatus(status v1alpha1.MetaStorageProvisionerStatus) EventLogger {
	e.Status = status
	return e
}

// Log writes the log message, records the event, and updates the status of the object if applicable.
func (e EventLogger) Log(ctx context.Context) {
	if e.err != nil {
		e.logger.Error(e.err, e.message)
		if e.Object != nil {
			e.recorder.Event(e.Object, "Warning", "Error", e.err.Error())
		}
		return
	}

	if e.Status != "" && e.Object != nil {
		e.Object.Status.Status = e.Status
		if err := e.client.Status().Update(ctx, e.Object); err != nil {
			e.logger.Error(err, "failed to update status")
			e.recorder.Event(e.Object, "Warning", "Error", err.Error())
		}
	}

	e.logger.Info(e.message)
	if e.Object != nil {
		e.recorder.Event(e.Object, "Normal", "Success", e.message)
	}
}
