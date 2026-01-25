/*
 * SPDX-FileCopyrightText: 2025 Pagefault Games
 * SPDX-FileContributor: SirzBenjie
 *
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package utils

import (
	"database/sql"
	"log"

	"github.com/amatsagu/tempest"
	"github.com/pagefaultgames/ticketune/constants"
)

// Base say command functionality reusable by multiple command implementations
// Parameters:
// `itx“: The command interaction to respond to
// `content“: The message content to send to the thread
// `invokerResponse`: The message to send back to the command invoker on success. On error, a relevant error message will be sent instead.
func SayCommandTemplate(itx *tempest.CommandInteraction,
	content string,
	invokerResponse string,
) {
	// Get the user associated with this thread (this handles responding to the interaction on error)
	userID, err := GetUserFromThread(itx)
	if err != sql.ErrNoRows && err != nil {
		return
	}

	// Discard error; if the option is missing, we default to `false`
	noPing, _ := GetOption[bool](itx, "no-ping", false)

	// The message to send publicly to the thread
	if err == nil && !noPing {
		content = "Hi <@" + userID.String() + ">!\n" + content
	}

	if err != nil {
		log.Println("Error fetching user for thread:", err)
		invokerResponse = constants.COULD_NOT_FIND_USER_TO_PING
	}

	// Send the user a message
	_, err = itx.Client.SendLinearMessage(
		itx.ChannelID,
		content,
	)
	if err != nil {
		itx.SendLinearReply("Something went wrong trying to send the message: "+err.Error(), true)
		return
	}

	itx.SendLinearReply(invokerResponse, true)
}
