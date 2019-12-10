/*
	Copyright 2019 Netfoundry, Inc.

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

package lib

import (
	"fmt"
	"github.com/netfoundry/fablab/kernel"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

func TemplateFuncMap(m *kernel.Model) template.FuncMap {
	return template.FuncMap{
		"publicIp": func(regionTag, hostTag string) string {
			host := m.GetHostByTags(regionTag, hostTag)
			if host != nil {
				return host.PublicIp
			}
			return ""
		},
	}
}

func RenderTemplate(src, dst string, m *kernel.Model, data interface{}) error {
	tData, err := ioutil.ReadFile(src)
	if err != nil {
		return fmt.Errorf("error reading template [%s] (%w)", src, err)
	}

	t, err := template.New("config").Funcs(TemplateFuncMap(m)).Parse(string(tData))
	if err != nil {
		return fmt.Errorf("error parsing template [%s] (%w)", src, err)
	}

	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return fmt.Errorf("error creating output parent directories [%s] (%w)", dst, err)
	}

	dstF, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating output [%s] (%w)", dst, err)
	}
	defer func() { _ = dstF.Close() }()

	err = t.Execute(dstF, data)
	if err != nil {
		return fmt.Errorf("error rendering template [%s] (%w)", src, err)
	}

	logrus.Infof("[%s] => [%s]", src, dst)

	return nil
}
