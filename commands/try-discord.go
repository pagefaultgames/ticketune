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

var tryDiscordCommandDescription = "Ping the user associated with this ticket and ask them to log into discord"

var tryWithDiscordMessage = "Please try to log in with Discord now, and let us know here if it works!\n\n" +
	"Alternatively, you can also try this:\n" +
	"Open Discord on your web browser. Connect to your Discord account (this one you are using). " +
	"Open PokéRogue in another tab, while keeping the Discord one open. " +
	"On the login page, click on the Discord button to try to log in with Discord."

var TryDiscordCommand = tempest.Command{
	Name:                "try-discord",
	Description:         tryDiscordCommandDescription,
	RequiredPermissions: tempest.ADMINISTRATOR_PERMISSION_FLAG,
	SlashCommandHandler: func(itx *tempest.CommandInteraction) {
		utils.SayCommandTemplate(itx, tryWithDiscordMessage, "The user has been requested to attempt a login.")
	},
	Contexts: []tempest.InteractionContextType{tempest.GUILD_CONTEXT_TYPE},
}
