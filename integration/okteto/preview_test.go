//go:build integration
// +build integration

// Copyright 2022 The Okteto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package okteto

import (
	"fmt"
	"testing"

	"github.com/okteto/okteto/integration"
	"github.com/okteto/okteto/integration/commands"
	"github.com/stretchr/testify/require"
)

func TestPreviewCommand(t *testing.T) {
	integration.SkipIfNotOktetoCluster(t)
	t.Parallel()
	oktetoPath, err := integration.GetOktetoPath()
	require.NoError(t, err)

	testNamespace := integration.GetTestNamespace("TestPreview", user)

	previewOptions := &commands.DeployPreviewOptions{
		Namespace:  testNamespace,
		Repository: fmt.Sprintf("%s/%s", githubHTTPSURL, pipelineRepo),
		Branch:     "cli-e2e",
		Wait:       true,
	}
	require.NoError(t, commands.RunOktetoDeployPreview(oktetoPath, previewOptions))

	contentURL := fmt.Sprintf("https://movies-%s.%s", testNamespace, appsSubdomain)
	require.NotEmpty(t, integration.GetContentFromURL(contentURL, timeout))

	previewDestroyOptions := &commands.DestroyPreviewOptions{
		Namespace: testNamespace,
	}
	require.NoError(t, commands.RunOktetoPreviewDestroy(oktetoPath, previewDestroyOptions))
}
