//  Copyright (c) 2018 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dockerfile

import "strings"

// EntrypointDirective represents the "ENTRYPOINT" dockerfile command.
type EntrypointDirective struct {
	*baseDirective
	Entrypoint []string
}

// Variables:
//   Replaced from ARGs and ENVs from within our stage.
// Formats:
//   ENTRYPOINT ["<executable>", "<param>"...]
//   ENTRYPOINT <command>
func newEntrypointDirective(base *baseDirective, state *parsingState) (Directive, error) {
	if err := base.replaceVarsCurrStage(state); err != nil {
		return nil, err
	}

	if entrypoint, ok := parseJSONArray(base.Args); ok {
		return &EntrypointDirective{base, entrypoint}, nil
	}

	// This is the Shell form (https://docs.docker.com/engine/reference/builder/#shell-form-entrypoint-example)
	// It is expected to wrap the whole entrypoint into a sh -c command)
	args, err := splitArgs(base.Args, true)
	if err != nil {
		return nil, base.err(err)
	}

	cmd := append([]string{"/bin/sh", "-c"}, strings.Join(args, " "))
	return &EntrypointDirective{base, cmd}, nil
}

// Add this command to the build stage.
func (d *EntrypointDirective) update(state *parsingState) error {
	return state.addToCurrStage(d)
}
