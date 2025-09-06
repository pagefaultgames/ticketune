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

var failDiscordCommandDescription = "Ping and ask the user to tell what happens when they attempt to login with Discord"

var failDiscordMessage = "Could you please tell me what happens when you try to log in by clicking on the Discord icon (on the right of the login page)?"

var FailDiscordCommand = tempest.Command{
	Name:                "tech-issues",
	Description:         failDiscordCommandDescription,
	RequiredPermissions: tempest.ADMINISTRATOR_PERMISSION_FLAG,
	SlashCommandHandler: func(itx *tempest.CommandInteraction) {
		utils.SayCommandTemplate(itx, failDiscordMessage, "The user has been asked to tell what happens when they attempt to login with Discord.")
	},
	Contexts: []tempest.InteractionContextType{tempest.GUILD_CONTEXT_TYPE},
}
