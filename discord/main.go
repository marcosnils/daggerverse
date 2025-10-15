// Dagger module whih allows calling the Discord API

package main

import (
	"context"
	"dagger/discord/internal/dagger"
	"io"
	"net/http"
	"strings"
)

type Discord struct {

	// +private
	BotToken *dagger.Secret
}

func New(ctx context.Context,

	botToken *dagger.Secret,
) *Discord {

	return &Discord{
		BotToken: botToken,
	}
}

// Calls the discord API and returns a JSON response
func (m *Discord) Call(
	ctx context.Context,
	// The HTTP method to use
	// +optional
	// +default="GET"
	method string,
	// The path to call after the "https://discord.com/api" base URl
	path string,

	// an optional body parameter for non GET requests
	// +optional
	body string,
) (string, error) {

	r, _ := http.NewRequest(strings.ToUpper(method), "https://discord.com/api"+path, strings.NewReader(body))
	token, err := m.BotToken.Plaintext(ctx)

	if err != nil {
		return "", nil
	}

	r.Header.Add("Authorization", "Bot "+token)
	r.Header.Add("Content-type", "application/json")

	res, err := http.DefaultClient.Do(r)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	ret, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(ret), nil

}
