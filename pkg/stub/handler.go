package stub

import (
	"context"
	"fmt"

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
		err := hellostateful.Create(hs)
		if err != nil {
			return fmt.Errorf("Failed to generate hello stateful custom resource: %v", err)
		}
	}
	return nil
}
