/*
	Copyright 2019 NetFoundry Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	https://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package subcmd

import (
	"github.com/openziti/fablab/kernel/model"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(disposeCmd)
}

var disposeCmd = &cobra.Command{
	Use:   "dispose",
	Short: "dispose of all model resources",
	Args:  cobra.ExactArgs(0),
	Run:   dispose,
}

func dispose(_ *cobra.Command, _ []string) {
	if err := model.Bootstrap(); err != nil {
		logrus.WithError(err).Fatal("unable to bootstrap")
	}

	ctx, err := model.NewRun()
	if err != nil {
		logrus.WithError(err).Fatal("error initializing run")
	}
	if err := ctx.GetModel().Dispose(ctx); err != nil {
		logrus.WithError(err).Fatal("error building configuration")
	}
}
