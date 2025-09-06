/*
 * SPDX-FileCopyrightText: 2025 Pagefault Games
 * SPDX-FileContributor: SirzBenjie
 * SPDX-FileContributor: patapancakes
 *
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package commands

import (
	"github.com/amatsagu/tempest"
)

var PingCommand = tempest.Command{
	Name:                "ping",
	Description:         "Check if the bot is alive",
	RequiredPermissions: tempest.ADMINISTRATOR_PERMISSION_FLAG,
	Options: []tempest.CommandOption{{
		Type:        tempest.BOOLEAN_OPTION_TYPE,
		Required:    false,
		Name:        "ephemeral",
		Description: "Whether the reply should be ephemeral (only visible to you, default false)",
	}},
	SlashCommandHandler: func(itx *tempest.CommandInteraction) {
		ephemeral := false
		if len(itx.Data.Options) > 0 {
			if v, ok := itx.Data.Options[0].Value.(bool); ok {
				ephemeral = v
			}
		}
		itx.SendLinearReply("I'm still alive!", ephemeral)
		return
	},
}
