/*
 * SPDX-FileCopyrightText: 2025 Pagefault Games
 * SPDX-FileContributor: SirzBenjie
 *
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package commands

import "github.com/amatsagu/tempest"

var NO_PING_OPTION = tempest.CommandOption{
	Type:        tempest.BOOLEAN_OPTION_TYPE,
	Name:        "no-ping",
	Description: "Do not ping the user associated with this ticket. Defaults to false (ping the user).",
	Required:    false,
}
