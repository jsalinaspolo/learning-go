package learn_httptest

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func githubSuccess() *httptest.Server {
	json := `
{
  "ref": "refs/tags/test-tag6",
  "url": "https://api.github.com/repos/an-owner/a-repo/git/refs/tags/test-tag6",
  "object": {
    "type": "commit",
    "sha": "2c88b23d2170e372979d1c606007a2a591d82d4d",
    "url": "https://api.github.com/repos/an-owner/a-repo/git/commits/2c88b23d2170e372979d1c606007a2a591d82d4d"
  },
  "node_id": "MDM6UmVmMjk2NTQ4ODk2OnJlZnMvdGFncy90ZXN0LXRhZzY="
}
`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, json)
	}))

	return ts
}

func TestCreateTag(t *testing.T) {
	t.Run("Create a Tag", func(t *testing.T) {
		server := githubSuccess()
		defer server.Close()

		ghClient := NewGithubClient(server.Client())
		ghClient.Client.BaseURL, _ = url.Parse(server.URL + "/")
		ref, err := ghClient.CreateTag("test-tag7", "2c88b23d2170e372979d1c606007a2a591d82d4d")

		data, _ := json.Marshal(ref)
		log.Printf("%s", data)

		require.NoError(t, err)
		require.NotNil(t, ref)
	})
}
