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

var tryDiscordCommandDescription = "Ping and ask the user to attempt to log with Discord"

var tryWithDiscordMessage = "Please try to log in with Discord now, and __let us know here if it works__! " +
	"Make sure you to use the __same Discord account__ you used to open this ticket!\n\n" +
	"Alternatively, you can also try this:\n" +
	"- Open Discord on your web browser\n" +
	"- Login with the Discord account you are currently using\n" +
	"- Open PokéRogue in another tab, while keeping the Discord one open\n" +
	"- On the login page, click on the Discord button to try to log in with Discord"

var TryDiscordCommand = tempest.Command{
	Name:                "try-discord",
	Description:         tryDiscordCommandDescription,
	RequiredPermissions: tempest.ADMINISTRATOR_PERMISSION_FLAG,
	SlashCommandHandler: func(itx *tempest.CommandInteraction) error {
		utils.SayCommandTemplate(itx, tryWithDiscordMessage, "The user has been requested to attempt a login with Discord.")
		return nil
	},
	Contexts: []tempest.InteractionContextType{tempest.GUILD_CONTEXT_TYPE},
}
