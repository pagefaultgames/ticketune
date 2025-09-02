/*
 * SPDX-FileCopyrightText: 2025 Pagefault Games
 * SPDX-FileContributor: SirzBenjie
 * SPDX-FileContributor: patapancakes
 *
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package commands

import (
	"database/sql"
	"log"

	utils "github.com/pagefaultgames/ticketune/utils"

	"github.com/amatsagu/tempest"
)

var tryDiscordCommandDescription = "Ping the user associated with this ticket and ask them to log into discord"

var TryDiscordCommand = tempest.Command{
	Name:                "try-discord",
	Description:         tryDiscordCommandDescription,
	RequiredPermissions: tempest.ADMINISTRATOR_PERMISSION_FLAG,
	SlashCommandHandler: tryDiscordCommandImpl,
	Contexts:            []tempest.InteractionContextType{tempest.GUILD_CONTEXT_TYPE},
}

var tryWithDiscordMessage = "Please try to log in with Discord now, and let us know here if it works!\n\n" +
	"Alternatively, you can also try this:\n" +
	"Open Discord on your web browser. Connect to your Discord account (this one you are using). " +
	"Open PokéRogue in another tab, while keeping the Discord one open. " +
	"On the login page, click on the Discord button to try to log in with Discord."

func tryDiscordCommandImpl(itx *tempest.CommandInteraction) {

	// Get the user associated with this thread (this handles responding to the interaction on error)
	userID, err := utils.GetUserFromThread(itx)
	if err != sql.ErrNoRows && err != nil {
		return
	}

	threadID := itx.ChannelID

	// The message to send publicaly to the thread
	var msg string
	// The message to use to respond to the interaction
	var responseMsg string

	if err == nil {
		msg = "Hi <@" + userID.String() + ">!\n" + tryWithDiscordMessage
		responseMsg = "The user has been requested to attempt a login."
	} else {
		log.Println("Error fetching user for thread:", err)
		msg = tryWithDiscordMessage
		responseMsg = "I couldn't find a user associated with this thread in my database, so I can't ping them." +
			"However, I've sent the login message to the thread."
	}

	// Send the user a message
	_, err = itx.Client.SendLinearMessage(
		threadID,
		msg,
	)

	if err != nil {
		itx.SendLinearReply("Something went wrong trying to send the message: "+err.Error(), true)
		return
	}

	itx.SendLinearReply(responseMsg, true)

}
