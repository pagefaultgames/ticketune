package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ticketune-bot/constants"
	ticketuneTypes "ticketune-bot/discord-types"

	"github.com/amatsagu/tempest"
)

// Fetch a channel object from its ID
func GetChannelFromID(client *tempest.Client, cid tempest.Snowflake) (channel ticketuneTypes.Channel, err error) {
	response, err := client.Rest.Request(
		http.MethodGet,
		fmt.Sprintf("/channels/%d", cid),
		nil,
	)

	if err != nil {
		return
	}

	err = json.Unmarshal(response, &channel)
	return
}

func CheckIfPasswordTicketChannel(channel ticketuneTypes.Channel) bool {
	// Check if this is a thread in the ticket channel
	if channel.ParentID != constants.TICKET_CHANNEL_ID || channel.Type != tempest.GUILD_PRIVATE_THREAD_CHANNEL_TYPE {
		return false
	}
	return true
}
