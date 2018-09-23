package stub

import (
	"context"

	"github.com/joatmon08/hello-stateful-operator/pkg/apis/hello-stateful/v1alpha1"
	"github.com/joatmon08/hello-stateful-operator/pkg/hellostateful"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
)

// NewHandler creates a new handler.
func NewHandler() sdk.Handler {
	return &Handler{}
}

// Handler stub
type Handler struct {
}

// Handle processes the event triggered
func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.HelloStateful:
		hs := o
		err := hellostateful.CreateVolume(hs)
		if err != nil {
			return err
		}

		err = hellostateful.UpdateStatus(hs)
		if err != nil {
			return err
		}

		err = hellostateful.Restore(hs)
		if err != nil {
			return err
		}

		err = hellostateful.Create(hs)
		if err != nil {
			return err
		}

		err = hellostateful.Update(hs)
		if err != nil {
			return err
		}

		// err = hellostateful.Backup(hs)
		// if err != nil {
		// 	return err
		// }
	}
	return nil
}
