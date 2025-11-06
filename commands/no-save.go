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

var noSaveCommandDescription = "Ping and ask the user to try to login on a different browser or device they may have also played on"
var tryDifferentBrowserMessage = "If there is another device or browser you've played on before, please __try to use the gear there__.\n" +
	"Otherwise, please provide:\n" +
	"- The username of the account you want to recover\n" +
	"- To the best of your memory, the date of account creation and the date you played for the last time on this account " +
	"(played for the last time = when you most recently started any kind of run)" +
	"- Any information regarding game stats and/or the progress of your Pokédex."

var NoSaveCommmand = tempest.Command{
	Name:                "no-save",
	Description:         noSaveCommandDescription,
	RequiredPermissions: tempest.ADMINISTRATOR_PERMISSION_FLAG,
	SlashCommandHandler: func(itx *tempest.CommandInteraction) {
		utils.SayCommandTemplate(itx, tryDifferentBrowserMessage, "The user has been requested to try to login on a different browser or device.")
	},
	Contexts: []tempest.InteractionContextType{tempest.GUILD_CONTEXT_TYPE},
}
