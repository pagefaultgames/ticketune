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

var noSaveCommandDescription = "Ping and ask the user to try to log on a different browser or device they may have also played on"
var tryDifferentBrowserMessage = "If there is another device or browser you've played on before, please use the gear there.\n" +
	"Otherwise, please provide:\n" +
	"- The username of the account you want to recover\n" +
	"- The date of account creation and/or the date you last played on this account " +
	"(last played meaning the date you last started any kind of run)."

var NoSaveCommmand = tempest.Command{
	Name:                "no-save",
	Description:         noSaveCommandDescription,
	RequiredPermissions: tempest.ADMINISTRATOR_PERMISSION_FLAG,
	SlashCommandHandler: func(itx *tempest.CommandInteraction) {
		utils.SayCommandTemplate(itx, tryDifferentBrowserMessage, "The user has been requested to try to log on a different browser or device.")
	},
	Contexts: []tempest.InteractionContextType{tempest.GUILD_CONTEXT_TYPE},
}
