// Extends the tempest library to correct the behavior of some functions using improper types

package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/amatsagu/tempest"
	"github.com/pagefaultgames/ticketune/types"
)

var ErrMissingRequiredField = errors.New("at least one of content, embeds, components, or files to be present")

// A replacement for `tempest.SendMessage` that accepts `types.CreateMessageParams` instead of `tempest.Message`
// Necessary, as tempest does not include support for fields like AllowedMentions
// At the moment, does not support files.
// Also adds an additional parameter, `discardResponse`, for when the message response is not needed
func SendDiscordMessage[
	t tempest.Message,
](
	client *tempest.Client,
	channelID tempest.Snowflake,
	message types.CreateMessageParams,
	files []tempest.File,
	discardResponse bool,
) (tempest.Message, error) {
	// Discord requires at least one of content, embeds, sticker_ids, components, files[n], or poll to be present.
	if message.Content == "" && len(message.Embeds) == 0 && len(message.Components) == 0 && len(files) == 0 {
		return tempest.Message{}, ErrMissingRequiredField
	}

	raw, err := client.Rest.RequestWithFiles(http.MethodPost, "/channels/"+channelID.String()+"/messages", message, files)
	if err != nil {
		return tempest.Message{}, err
	}

	if discardResponse {
		return tempest.Message{}, nil
	}

	res := tempest.Message{}
	err = json.Unmarshal(raw, &res)
	if err != nil {
		return tempest.Message{}, err
	}

	return res, nil

}
