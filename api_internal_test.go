package openai

import (
	"context"
	"testing"
)

func TestOpenAIFullURL(t *testing.T) {
	cases := []struct {
		Name   string
		Suffix string
		Expect string
	}{
		{
			"ChatCompletionsURL",
			"/chat/completions",
			"https://api.openai.com/v1/chat/completions",
		},
		{
			"CompletionsURL",
			"/completions",
			"https://api.openai.com/v1/completions",
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			cli := NewClient("dummy")
			actual := cli.fullURL(c.Suffix)
			if actual != c.Expect {
				t.Errorf("Expected %s, got %s", c.Expect, actual)
			}
			t.Logf("Full URL: %s", actual)
		})
	}
}

func TestRequestAuthHeader(t *testing.T) {
	cases := []struct {
		Name      string
		APIType   APIType
		HeaderKey string
		Token     string
		OrgID     string
		Expect    string
	}{
		{
			"OpenAIDefault",
			"",
			"Authorization",
			"dummy-token-openai",
			"",
			"Bearer dummy-token-openai",
		},
		{
			"OpenAIOrg",
			APITypeOpenAI,
			"Authorization",
			"dummy-token-openai",
			"dummy-org-openai",
			"Bearer dummy-token-openai",
		},
		{
			"OpenAI",
			APITypeOpenAI,
			"Authorization",
			"dummy-token-openai",
			"",
			"Bearer dummy-token-openai",
		},
		{
			"AzureAD",
			APITypeAzureAD,
			"Authorization",
			"dummy-token-azure",
			"",
			"Bearer dummy-token-azure",
		},
		{
			"Azure",
			APITypeAzure,
			AzureAPIKeyHeader,
			"dummy-api-key-here",
			"",
			"dummy-api-key-here",
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			cli := NewAzureClient(c.Token, "", "", WithOrganizationID(c.OrgID), WithSpecificAPIType(c.APIType))
			req, err := cli.newStreamRequest(context.Background(), "POST", "/chat/completions", nil)
			if err != nil {
				t.Errorf("Failed to create request: %v", err)
			}
			actual := req.Header.Get(c.HeaderKey)
			if actual != c.Expect {
				t.Errorf("Expected %s, got %s", c.Expect, actual)
			}
			t.Logf("%s: %s", c.HeaderKey, actual)
		})
	}
}

func TestAzureFullURL(t *testing.T) {
	cases := []struct {
		Name    string
		BaseURL string
		Engine  string
		Expect  string
	}{
		{
			"AzureBaseURLWithSlashAutoStrip",
			"https://httpbin.org/",
			"chatgpt-demo",
			"https://httpbin.org/" +
				"openai/deployments/chatgpt-demo" +
				"/chat/completions?api-version=2023-03-15-preview",
		},
		{
			"AzureBaseURLWithoutSlashOK",
			"https://httpbin.org",
			"chatgpt-demo",
			"https://httpbin.org/" +
				"openai/deployments/chatgpt-demo" +
				"/chat/completions?api-version=2023-03-15-preview",
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			cli := NewAzureClient("dummy", c.BaseURL, c.Engine)
			// /openai/deployments/{engine}/chat/completions?api-version={api_version}
			actual := cli.fullURL("/chat/completions")
			if actual != c.Expect {
				t.Errorf("Expected %s, got %s", c.Expect, actual)
			}
			t.Logf("Full URL: %s", actual)
		})
	}
}
