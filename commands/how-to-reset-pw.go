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

var howResetPwDescription = "Ping the user to explain them where to change their password"

var howResetPwMessage = "You can change your password once logged in.\n" +
	"__Press Escape to open the menu__, then go to “Manage Data”, and finally choose “Change Password”. " +
	"This will log you out of all other devices.\n" +
	"Be sure to write down or remember this new password!"

var HowResetPwCommand = tempest.Command{
	Name:                "how-to-reset-pw",
	Description:         howResetPwCommandDescription,
	RequiredPermissions: tempest.ADMINISTRATOR_PERMISSION_FLAG,
	SlashCommandHandler: func(itx *tempest.CommandInteraction) {
		utils.SayCommandTemplate(itx, howResetPwMessage, "The user has been explain how to change their password.")
	},
	Contexts: []tempest.InteractionContextType{tempest.GUILD_CONTEXT_TYPE},
}
