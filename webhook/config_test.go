// Copyright 2026 LiveKit, Inc.
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

package webhook

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/livekit/protocol/auth"
)

func TestResourceURLNotifierHTTPClientConfig(t *testing.T) {
	conf := DefaultWebHookConfig
	err := yaml.Unmarshal([]byte(`
api_key: mykey
urls: [http://localhost:1234]
resource_url_notifier:
  max_age: 10s
  max_retries: 7
  retry_wait_min: 2s
  retry_wait_max: 1m
  client_timeout: 15s
`), &conf)
	require.NoError(t, err)

	rc := conf.ResourceURLNotifier
	require.Equal(t, 10*time.Second, rc.MaxAge)
	require.Equal(t, DefaultResourceURLNotifierConfig.MaxDepth, rc.MaxDepth)
	require.Equal(t, 7, rc.HTTPClient.MaxRetries)
	require.Equal(t, 2*time.Second, rc.HTTPClient.RetryWaitMin)
	require.Equal(t, time.Minute, rc.HTTPClient.RetryWaitMax)
	require.Equal(t, 15*time.Second, rc.HTTPClient.ClientTimeout)

	qn, err := NewDefaultNotifier(conf, auth.NewSimpleKeyProvider("mykey", "mysecret"))
	require.NoError(t, err)
	n := qn.(*DefaultNotifier)
	defer n.Stop(true)

	require.Len(t, n.notifiers, 1)
	for _, u := range []QueuedNotifier{n.notifiers[0], n.extraWebhookNotifier} {
		c := u.(*ResourceURLNotifier).client
		require.Equal(t, 7, c.RetryMax)
		require.Equal(t, 2*time.Second, c.RetryWaitMin)
		require.Equal(t, time.Minute, c.RetryWaitMax)
		require.Equal(t, 15*time.Second, c.HTTPClient.Timeout)
	}
}
