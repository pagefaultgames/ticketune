/*
 * SPDX-FileCopyrightText: 2025 Pagefault Games
 * SPDX-FileContributor: Lugiadrien
 *
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package commands

import (
	utils "github.com/pagefaultgames/ticketune/utils"

	"github.com/amatsagu/tempest"
)

var pingSpamCommandDescription = "Ping the user tell them to stop ping abuse"

var pingSpamMessage = "**Please keep in mind that Helpers are real people volunteering on their free time, so please don't ping them over and over.**\n" +
	"**When someone is free, they'll reach out to help you, but until then, please be patient and wait until someone gets back to you.**\n\n" +
	"To be more precise, there is currently a colossal amount of 3 people helping on these tickets, all of them also working on other things for the project as well as having real life on the side.\n" +
	"Please understand that the 20 tickets they get every day take a lot of time to deal with, they might be able to check them only once or twice per day and it can take them more than half an hour.\n" +
	"They will absolutely help you, **just give them time please**!"

var PingSpamCommand = tempest.Command{
	Name:                "ping-spam",
	Description:         pingSpamCommandDescription,
	RequiredPermissions: tempest.ADMINISTRATOR_PERMISSION_FLAG,
	SlashCommandHandler: func(itx *tempest.CommandInteraction) {
		utils.SayCommandTemplate(itx, pingSpamMessage, "The user has been asked to stop ping abuse.")
	},
	Contexts: []tempest.InteractionContextType{tempest.GUILD_CONTEXT_TYPE},
}
