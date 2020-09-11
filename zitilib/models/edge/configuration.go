/*
	Copyright 2020 NetFoundry, Inc.

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

package edge

import (
	"github.com/openziti/fablab/kernel/fablib/runlevel/1_configuration/config"
	"github.com/openziti/fablab/kernel/fablib/runlevel/2_kitting/devkit"
	"github.com/openziti/fablab/kernel/model"
	zitilib_bootstrap "github.com/openziti/fablab/zitilib"
	"github.com/openziti/fablab/zitilib/runlevel/1_configuration"
)

func newConfigurationFactory() model.Factory {
	return &configurationFactory{}
}

func (self *configurationFactory) Build(m *model.Model) error {
	configs := []config.StaticConfig{
		{Src: "loop/10-ambient.loop2.yml", Name: "10-ambient.loop2.yml"},
		{Src: "loop/4k-chatter.loop2.yml", Name: "4k-chatter.loop2.yml"},
	}

	zitiBinaries := []string{
		"ziti-controller",
		"ziti-fabric",
		"ziti-fabric-test",
		"ziti-router",
	}

	m.Configuration = model.ConfigurationStages{
		zitilib_runlevel_1_configuration.IfNoPki(zitilib_runlevel_1_configuration.Fabric(), zitilib_runlevel_1_configuration.DotZiti()),
		config.Component(),
		config.Static(configs),
		devkit.DevKit(zitilib_bootstrap.ZitiDistBinaries(), zitiBinaries),
	}

	return nil
}

type configurationFactory struct{}