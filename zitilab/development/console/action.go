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

package console

import (
	"github.com/netfoundry/fablab/kernel/model"
	"net/http"
	"path/filepath"
)

func Console() model.Action {
	return &console{}
}

func (consoleAction *console) Execute(m *model.Model) error {
	server := NewServer()
	go server.Listen()

	http.Handle("/", http.FileServer(http.Dir(filepath.Join(model.FablabRoot(), "zitilab/console/webroot"))))
	return http.ListenAndServe(":8080", nil)
}

type console struct{}