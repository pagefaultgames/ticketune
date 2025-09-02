/*
 * SPDX-FileCopyrightText: 2025 Pagefault Games
 * SPDX-FileContributor: SirzBenjie
 * SPDX-FileContributor: patapancakes
 *
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package commands

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/pagefaultgames/ticketune/utils"

	"github.com/amatsagu/tempest"
)

const oldAccountCommandDescription = "Notify the user that they are likely misremembering their username due to account inactivity"
const oldAccountDefaultDescription = "Notify the user the account has not been played for a long time"
const oldAccountSpecificDescription = "Notify the user the account has not been played for a specified amount of time"

var OldAccountCommandGroup = tempest.Command{
	Name:        "old-account",
	Description: oldAccountCommandDescription,
}

var OldAccountDefault = tempest.Command{
	Name:                "default",
	Description:         oldAccountDefaultDescription,
	RequiredPermissions: tempest.ADMINISTRATOR_PERMISSION_FLAG,
	SlashCommandHandler: func(itx *tempest.CommandInteraction) { oldAccountCommandImpl(itx, true) },
	Options: []tempest.CommandOption{{
		Type:        tempest.STRING_OPTION_TYPE,
		Name:        "username",
		Description: "The username of the old account",
		Required:    false,
		MinLength:   1,
		MaxLength:   64,
	}},
	Contexts: []tempest.InteractionContextType{tempest.GUILD_CONTEXT_TYPE},
}

var OldAccountSpecific = tempest.Command{
	Name:                "specific",
	Description:         oldAccountSpecificDescription,
	RequiredPermissions: tempest.ADMINISTRATOR_PERMISSION_FLAG,
	SlashCommandHandler: func(itx *tempest.CommandInteraction) { oldAccountCommandImpl(itx, false) },
	Contexts:            []tempest.InteractionContextType{tempest.GUILD_CONTEXT_TYPE},
	Options: []tempest.CommandOption{
		{
			Type:        tempest.STRING_OPTION_TYPE,
			Name:        "username",
			Description: "The username of the old account",
			Required:    true,
			MinLength:   1,
			MaxLength:   64,
		},
		{
			Type:        tempest.STRING_OPTION_TYPE,
			Name:        "unit",
			Description: "The unit of time",
			Required:    true,
			Choices: []tempest.CommandOptionChoice{
				{Name: "week", Value: "week"},
				{Name: "month", Value: "month"},
				{Name: "year", Value: "year"},
			},
		},
		{
			Type:        tempest.INTEGER_OPTION_TYPE,
			Name:        "amount",
			Description: "The number of time units. If omitted, will just say \"in a few [units]\"",
			Required:    false,
		},
	},
}

// Given an amount and a unit, return a string like "since last [unit]" or "for [amount] [units]"
func computeTimeString(amount int, unit string) string {
	// if unit is 1, then say, "since last [unit]"
	switch amount {
	case 1:
		return "since last " + unit
	case 0:
		return "for a few " + unit + "s"
	default:
		return "for " + strconv.Itoa(amount) + " " + unit + "s"
	}
}

// Build the message to send to the user
// Notably, `username` is *not* the discord username, but the old account username
func buildMessage(username string, amount int, unit string) string {
	return "The account ``" + username + "`` has not played " + computeTimeString(amount, unit) + ". " +
		"It is likely that you are misremembering your username."
}

const defaultMessage string = " has not played for a long time." +
	" It is likely that you are misremembering your username."

func defaultMessageWithUsername(username string) string {
	// If username was empty, use the generic "The account you provided"
	if username == "" {
		return "The account you provided" + defaultMessage
	}

	return "The account ``" + username + "``" + defaultMessage
}

func oldAccountCommandImpl(itx *tempest.CommandInteraction, isDefault bool) {
	// if no arguments, then use a default message
	var msg string
	if !isDefault {
		username, parseError := utils.GetOption[string](itx, "username", true)
		if parseError != nil {
			return
		}

		unit, parseError := utils.GetOption[string](itx, "unit", true)
		if parseError != nil {
			return
		}

		amount, _ := utils.GetNumericOption[int](itx, "amount", false)
		// If there was a parse error, treat it as if they didn't provide the option, (default to "in a few [units]")
		msg = buildMessage(username, amount, unit)
	} else {
		// if the option was not provided, username will be the zero value, which is fine
		username, _ := utils.GetOption[string](itx, "username", false)
		msg = defaultMessageWithUsername(username)
	}

	// Get the user associated with this thread (this handles responding to the	 interaction on error)
	userID, err := utils.GetUserFromThread(itx)
	if err != sql.ErrNoRows && err != nil {
		return
	}

	// The message to use to respond to the interaction
	responseMsg := "The user has been notified."

	msg = "Hi <@" + userID.String() + ">!\n" + msg

	if err != nil {
		log.Println("Error fetching user for thread:", err)
		responseMsg = "I couldn't find a user associated with this thread in my database, so I can't ping them." +
			"However, I've sent the message to the thread."
	}

	// Send the user a message
	_, err = itx.Client.SendLinearMessage(
		itx.ChannelID,
		msg,
	)
	if err != nil {
		itx.SendLinearReply("Something went wrong trying to send the message: "+err.Error(), true)
		return
	}

	itx.SendLinearReply(responseMsg, true)
}
