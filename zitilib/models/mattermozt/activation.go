package mattermozt

import (
	"github.com/openziti/fablab/kernel/fablib/runlevel/4_activation/action"
	"github.com/openziti/fablab/kernel/model"
)

func newActivationFactory() model.Factory {
	return &activationFactory{}
}

func (f *activationFactory) Build(m *model.Model) error {
	m.Activation = model.ActivationBinders{
		func(m *model.Model) model.ActivationStage {
			return action.Activation("bootstrap", "start")
		},
	}
	return nil
}

type activationFactory struct{}
