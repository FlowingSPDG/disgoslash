package disgoslash

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	discord "github.com/wafer-bw/disgoslash/discord"
	"github.com/wafer-bw/disgoslash/errs"
)

const baseURL string = "https://discord.com/api"
const apiVersion string = "v8"

// client implements a `clientInterface` interface's properties
type client struct {
	apiURL    string
	authToken string
}

// clientInterface methods
type clientInterface interface {
	ListApplicationCommands(guildID string) ([]*discord.ApplicationCommand, error)
	CreateApplicationCommand(guildID string, command *discord.ApplicationCommand) error
	DeleteApplicationCommand(guildID string, commandID string) error
}

// NewClient creates a new `clientInterface` instance
func newClient(creds *discord.Credentials) clientInterface {
	return constructClient(creds, baseURL, apiVersion)
}

func constructClient(creds *discord.Credentials, baseURL string, apiVersion string) clientInterface {
	return &client{
		apiURL:    fmt.Sprintf("%s/%s/applications/%s", baseURL, apiVersion, creds.ClientID),
		authToken: fmt.Sprintf("Bot %s", creds.Token),
	}
}

// ListApplicationCommands // todo
func (client *client) ListApplicationCommands(guildID string) ([]*discord.ApplicationCommand, error) {
	var url string
	if guildID == "" {
		url = fmt.Sprintf("%s/commands", client.apiURL)
	} else {
		url = fmt.Sprintf("%s/guilds/%s/commands", client.apiURL, guildID)
	}
	return client.listApplicationCommands(url)
}

// CreateApplicationCommand // todo
func (client *client) CreateApplicationCommand(guildID string, command *discord.ApplicationCommand) error {
	var url string
	if guildID == "" {
		url = fmt.Sprintf("%s/commands", client.apiURL)
	} else {
		url = fmt.Sprintf("%s/guilds/%s/commands", client.apiURL, guildID)
	}
	return client.createApplicationCommand(url, command)
}

// DeleteApplicationCommand // todo
func (client *client) DeleteApplicationCommand(guildID string, commandID string) error {
	var url string
	if guildID == "" {
		url = fmt.Sprintf("%s/commands/%s", client.apiURL, commandID)
	} else {
		url = fmt.Sprintf("%s/guilds/%s/commands/%s", client.apiURL, guildID, commandID)
	}
	return client.deleteApplicationCommands(url)
}

func (client *client) listApplicationCommands(url string) ([]*discord.ApplicationCommand, error) {
	status, data, err := client.request(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	} else if status != http.StatusOK {
		return nil, fmt.Errorf("%d - %s", status, string(data))
	}
	commands := &[]*discord.ApplicationCommand{}
	if err := unmarshal(data, commands); err != nil {
		return nil, err
	}
	return *commands, nil
}

func (client *client) createApplicationCommand(url string, command *discord.ApplicationCommand) error {
	body, err := marshal(command)
	if err != nil {
		return err
	}
	if status, data, err := client.request(http.MethodPost, url, body); err != nil {
		return err
	} else if status == http.StatusOK {
		return errs.ErrAlreadyExists
	} else if status != http.StatusCreated {
		return fmt.Errorf("%d - %s", status, string(data))
	}
	return nil
}

func (client *client) deleteApplicationCommands(url string) error {
	if status, data, err := client.request(http.MethodDelete, url, nil); err != nil {
		return err
	} else if status != http.StatusNoContent {
		return fmt.Errorf("%d - %s", status, string(data))
	}
	return nil
}

func (client *client) request(method string, url string, body io.Reader) (int, []byte, error) {
	attempts := 0
	maxAttempts := 3

	for attempts < maxAttempts {
		attempts++

		httpClient := &http.Client{}
		request, err := http.NewRequest(method, url, body)
		if err != nil {
			return 0, nil, err
		}

		request.Header.Set("content-type", "application/json")
		request.Header.Set("authorization", client.authToken)

		response, err := httpClient.Do(request)
		if err != nil {
			return 0, nil, err
		}

		switch response.StatusCode {
		case http.StatusForbidden:
			return 0, nil, errs.ErrForbidden
		case http.StatusUnauthorized:
			return 0, nil, errs.ErrUnauthorized
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return 0, nil, err
		}

		waitTime, err := determineRetry(response.StatusCode, data)
		if err != nil {
			return 0, nil, err
		}

		if waitTime <= 0 {
			return response.StatusCode, data, nil
		}
		time.Sleep(waitTime)
	}
	return 0, nil, errs.ErrMaxRetries
}

func unmarshal(body []byte, v interface{}) error {
	if err := json.Unmarshal(body, v); err != nil {
		return err
	}
	return nil
}

func marshal(v interface{}) (io.Reader, error) {
	body, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(body), nil
}

func determineRetry(statusCode int, data []byte) (time.Duration, error) {
	if statusCode != http.StatusTooManyRequests {
		return 0, nil
	}
	responseErr := &discord.APIErrorResponse{}
	if err := unmarshal(data, responseErr); err != nil {
		return 0, err
	}
	return time.Duration(responseErr.RetryAfter) * time.Second, nil
}
