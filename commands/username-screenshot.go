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

var usernameScreenshotDescription = "Ping and ask the user to check for any screenshot or .prsv file where their username can appear."
var usernameScreenshotMessage = "By any chance, maybe you have some screenshot with your username visible, or even a PokéRogue save file (.prsv)?\n" +
	"In some device, Discord server, DMs, etc...?\n" +
	"__The username can apprear on screenshots taken from:__\n" +
	"- The first page of a Pokémon Sumamry, as OT\n" +
	"- Game stats screen *(since August 23rd 2025)*\n" +
	"- Title screen *(since October 31st 2025)*"

var UsernameScreenshotCommmand = tempest.Command{
	Name:                "username-screenshot",
	Description:         usernameScreenshotCommandDescription,
	RequiredPermissions: tempest.ADMINISTRATOR_PERMISSION_FLAG,
	SlashCommandHandler: func(itx *tempest.CommandInteraction) {
		utils.SayCommandTemplate(itx, usernameScreenshotMessage, "The user has been requested to check for any screenshot or .prsv file where their username can appear.")
	},
	Contexts: []tempest.InteractionContextType{tempest.GUILD_CONTEXT_TYPE},
}
