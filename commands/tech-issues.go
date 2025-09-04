/*
 * SPDX-FileCopyrightText: 2025 Pagefault Games
 * SPDX-FileContributor: SirzBenjie
 * SPDX-FileContributor: patapancakes
 * SPDX-FileContributor: Lugiadrien
 *
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package commands

import (
	utils "github.com/pagefaultgames/ticketune/utils"

	"github.com/amatsagu/tempest"
)

var techIssuesCommandDescription = "Ping the user about technical issues prevening to help them"

var techIssuesMessage = "We are currently experiencing technical issues preventing us to help you as we speak. :jolteondead:\n" +
	"We apologize for the inconvenience, and we'll ping you as soon as possible once the issue is solved!\n\n" +
	"Also, if in the meantime you happen to remember your password or prefer to close your ticket for the time being, __please let us now__!"

var TechIssuesCommand = tempest.Command{
	Name:                "tech-issues",
	Description:         techIssuesCommandDescription,
	RequiredPermissions: tempest.ADMINISTRATOR_PERMISSION_FLAG,
	SlashCommandHandler: func(itx *tempest.CommandInteraction) {
		utils.SayCommandTemplate(itx, techIssuesMessage, "The user has been warned about technical issues prevening to help them.")
	},
	Contexts: []tempest.InteractionContextType{tempest.GUILD_CONTEXT_TYPE},
}
