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

const saveAccessDescription = "Informs the user about what Helpers can check about their saves"
const saveAccessMessage = "Here is the range of things we can check or not about a username:\n" +
	"- When it has been created\n" +
	"- When it has been saved for the last time\n" +
	"- Game stats screen\n" +
	"- Pokédex progress\n" +
	"- We **can't** check the content of any of your runs\n" +
	"- We **can't** check your Run History\n\n"

var SaveAccessCommmand = tempest.Command{
	Name:                "save-access",
	Description:         saveAccessDescription,
	RequiredPermissions: tempest.ADMINISTRATOR_PERMISSION_FLAG,
	Options:             []tempest.CommandOption{NO_PING_OPTION},
	SlashCommandHandler: func(itx *tempest.CommandInteraction) {
		utils.SayCommandTemplate(itx, saveAccessMessage, "The user has been informed about what Helpers can check about their saves.")
	},
	Contexts: []tempest.InteractionContextType{tempest.GUILD_CONTEXT_TYPE},
}
