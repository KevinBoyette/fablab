/*
	Copyright 2019 NetFoundry, Inc.

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
	"fmt"
	"github.com/openziti/fablab/kernel/fablib"
	"github.com/openziti/fablab/kernel/model"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	cmd := newSshExecCmd()
	RootCmd.AddCommand(cmd.cobraCmd)
}

type sshExecCmd struct {
	cobraCmd    *cobra.Command
	concurrency int
}

func newSshExecCmd() *sshExecCmd {
	cmd := &sshExecCmd{
		cobraCmd: &cobra.Command{
			Use:   "sshexec <hostSpec> <cmd>",
			Short: "establish an ssh connection to the model and runs the given command on the selected hosts",
			Args:  cobra.ExactArgs(2),
		},
	}

	cmd.cobraCmd.Run = cmd.run
	cmd.cobraCmd.Flags().IntVarP(&cmd.concurrency, "concurrency", "c", 1, "Number of hosts to run in parallel")
	return cmd
}

func (cmd *sshExecCmd) run(_ *cobra.Command, args []string) {
	if err := model.Bootstrap(); err != nil {
		logrus.Fatalf("unable to bootstrap (%s)", err)
	}

	label := model.GetLabel()
	if label == nil {
		logrus.Fatalf("no label for instance [%s]", model.ActiveInstancePath())
	}

	if label != nil {
		m, found := model.GetModel(label.Model)
		if !found {
			logrus.Fatalf("no such model [%s]", label.Model)
		}

		if !m.IsBound() {
			logrus.Fatalf("model not bound")
		}

		logrus.Infof("executing %v with concurrency %v", args[1], cmd.concurrency)
		err := m.ForEachHost(args[0], cmd.concurrency, func(h *model.Host) error {
			sshConfigFactory := fablib.NewSshConfigFactoryImpl(m, h.PublicIp)
			o, err := fablib.RemoteExecAll(sshConfigFactory, args[1])
			if err != nil {
				logrus.Errorf("output [%s]", o)
				return fmt.Errorf("error executing process on [%s] (%s)", h.PublicIp, err)
			}
			logrus.Infof("[%v] output:\n%s", h.PublicIp, o)
			return nil
		})

		if err != nil {
			logrus.Fatalf("error executing remote shell (%v)", err)
		}
	}
}
