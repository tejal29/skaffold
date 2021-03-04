/*
Copyright 2021 The Skaffold Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package error

import (
	"fmt"
	"testing"

	sErrors "github.com/GoogleContainerTools/skaffold/pkg/skaffold/errors"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/runner/runcontext"
	"github.com/GoogleContainerTools/skaffold/proto/v1"
	"github.com/GoogleContainerTools/skaffold/testutil"
)

func TestUserError(t *testing.T) {
	tests := []struct {
		description string
		expected    string
		err         error
		isMinikube  bool
		shdInternal bool
	}{
		{
			description: "internal system error for minikube cluster",
			isMinikube:  true,
			expected:    "Deploy Failed. Error: (Internal Server Error: the server is currently unable to handle the request). Run minikube status -p test to check if minikube is running. Try again.\nIf this keeps happening please open an issue https://github.com/GoogleContainerTools/skaffold/issues/new.",
			err:         fmt.Errorf("Error: (Internal Server Error: the server is currently unable to handle the request)"),
			shdInternal: true,
		},
		{
			description: "internal system error for k8s cluster",
			isMinikube:  false,
			expected:    "Deploy Failed. Error: (Internal Server Error: the server is currently unable to handle the request). Something went wrong with your cluster \"test\". Try again.\nIf this keeps happening please open an issue https://github.com/GoogleContainerTools/skaffold/issues/new.",
			err:         fmt.Errorf("Error: (Internal Server Error: the server is currently unable to handle the request)"),
			shdInternal: true,
		},
		{
			description: "random tool error",
			isMinikube:  true,
			expected:    "helm --namepsace wrong flag",
			err:         fmt.Errorf("helm --namepsace wrong flag"),
		},
	}
	for _, test := range tests {
		testutil.Run(t, test.description, func(t *testutil.T) {
			t.Override(&isMinikube, func(string) bool {
				return test.isMinikube
			})
			sErrors.SetRunContext(runcontext.RunContext{KubeContext: "test"})
			actual := UserError(test.err, proto.StatusCode_DEPLOY_HELM_USER_ERR)
			t.CheckDeepEqual(test.expected, sErrors.ShowAIError(actual).Error())
		})
	}
}
