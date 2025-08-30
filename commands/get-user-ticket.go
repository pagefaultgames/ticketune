package command

// import (
// 	ticketune_db "ticketune-bot/ticketune-db"

// 	"github.com/amatsagu/tempest"
// )

// var CreateSupportTicketCommand tempest.Command = tempest.Command{
// 	Name:                "get-user-ticket",
// 	Description:         "Get a link to the support ticket thread for a user, if it exists",
// 	SlashCommandHandler: getTicketForUser,
// 	// By default, only let admins use the command.
// 	RequiredPermissions: tempest.ADMINISTRATOR_PERMISSION_FLAG,
// 	Contexts:            []tempest.InteractionContextType{tempest.GUILD_CONTEXT_TYPE},
// 	Options: []tempest.CommandOption{{
// 		Type:        tempest.USER_OPTION_TYPE,
// 		Name:        "user",
// 		Description: "User to get the ticket for",
// 		Required:    true,
// 	}},
// }

// // func getTicketForUser(itx *tempest.CommandInteraction) {
// // 	if len(itx.Data.Options) == 0 {
// // 		itx.SendLinearReply("You must specify a user", true)
// // 		return
// // 	}
// // 	userID := itx.Data.Options[0].Value.
// // 	threadID, _, err = ticketune_db.GetDB().GetUserThread(userID)
// // 	return
// // }
