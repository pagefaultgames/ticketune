package command

import (
	"database/sql"
	"fmt"
	ticketune_db "ticketune-bot/ticketune-db"

	"github.com/amatsagu/tempest"
)

var GetUserTicketCommand tempest.Command = tempest.Command{
	Name:                "get-user-ticket",
	Description:         "Get a link to the support ticket thread for a user, if it exists",
	SlashCommandHandler: getUserTicketCommandImpl,
	// By default, only let admins use the command.
	RequiredPermissions: tempest.ADMINISTRATOR_PERMISSION_FLAG,
	Contexts:            []tempest.InteractionContextType{tempest.GUILD_CONTEXT_TYPE},
	Options: []tempest.CommandOption{{
		Type:        tempest.USER_OPTION_TYPE,
		Name:        "user",
		Description: "User to get the ticket for",
		Required:    true,
	}},
}

func getUserTicketCommandImpl(itx *tempest.CommandInteraction) {
	userIDStr, present := itx.GetOptionValue("user")
	if !present {
		itx.SendLinearReply("You must specify a user", true)
		return
	}
	userID, err := tempest.StringToSnowflake(userIDStr.(string))
	if err != nil {
		itx.SendLinearReply("Invalid user ID", true)
		return
	}
	tid, err := ticketune_db.GetDB().GetUserThread(userID)
	if err == sql.ErrNoRows {
		itx.SendLinearReply("This user does not have an open support ticket", true)
		return
	}
	itx.SendReply(tempest.ResponseMessageData{
		Content: fmt.Sprintf("Support ticket thread: <#%d>", tid),
	}, true, nil)
}
