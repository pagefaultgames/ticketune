/*
 * SPDX-FileCopyrightText: 2025 Pagefault Games
 * SPDX-FileContributor: SirzBenjie
 * SPDX-FileContributor: patapancakes
 *
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package commands

import (
	utils "github.com/pagefaultgames/ticketune/utils"

	"github.com/amatsagu/tempest"
)

var noSaveCommandDescription = "Ping the user associated with this ticket and ask them to try to use a different browser"
var tryDifferentBrowserMessage = "If there is another device or browser you've played on before, please use the gear there.\n" +
	"Otherwise, please provide your username, as well as the date of account creation or the date you last played on this account " +
	"(last played meaning the date you last started any kind of run)."

var NoSaveCommmand = tempest.Command{
	Name:                "no-save",
	Description:         noSaveCommandDescription,
	RequiredPermissions: tempest.ADMINISTRATOR_PERMISSION_FLAG,
	SlashCommandHandler: func(itx *tempest.CommandInteraction) {
		utils.SayCommandTemplate(itx, tryDifferentBrowserMessage, "The user has been requested to try a different device/browser.")
	},
	Contexts: []tempest.InteractionContextType{tempest.GUILD_CONTEXT_TYPE},
}
