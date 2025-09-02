/*
 * SPDX-FileCopyrightText: 2025 Pagefault Games
 * SPDX-FileContributor: SirzBenjie
 * SPDX-FileContributor: patapancakes
 *
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pagefaultgames/ticketune/constants"
	"github.com/pagefaultgames/ticketune/types"

	"github.com/amatsagu/tempest"
)

// Fetch a channel object from its ID
func GetChannelFromID(client *tempest.Client, cid tempest.Snowflake) (types.Channel, error) {
	response, err := client.Rest.Request(http.MethodGet, fmt.Sprintf("/channels/%d", cid), nil)
	if err != nil {
		return types.Channel{}, err
	}

	var channel types.Channel
	err = json.Unmarshal(response, &channel)
	if err != nil {
		return types.Channel{}, err
	}

	return channel, nil
}

func CheckIfPasswordTicketChannel(channel types.Channel) bool {
	// Check if this is a thread in the ticket channel
	if channel.ParentID != constants.TICKET_CHANNEL_ID || channel.Type != tempest.GUILD_PRIVATE_THREAD_CHANNEL_TYPE {
		return false
	}

	return true
}
