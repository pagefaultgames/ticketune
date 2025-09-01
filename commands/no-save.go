package commands

import (
	"database/sql"
	"log"

	utils "github.com/pagefaultgames/ticketune/utils"

	"github.com/amatsagu/tempest"
)

var noSaveCommandDescription = "Ping the user associated with this ticket and ask them to try to use a different browser"

var NoSaveCommmand = tempest.Command{
	Name:                "no-save",
	Description:         noSaveCommandDescription,
	RequiredPermissions: tempest.ADMINISTRATOR_PERMISSION_FLAG,
	SlashCommandHandler: noSaveCommmandImpl,
	Contexts:            []tempest.InteractionContextType{tempest.GUILD_CONTEXT_TYPE},
}

var tryDifferentBrowserMessage = "If there is another device or browser you've played on before, please use the gear there.\n" +
	"Otherwise, please provide your username, as well as the date of account creation or the date you last played on this account " +
	"(last played meaning the date you last started any kind of run)."

func noSaveCommmandImpl(itx *tempest.CommandInteraction) {
	// Get the user associated with this thread (this handles responding to the interaction on error)
	userID, err := utils.GetUserFromThread(itx)
	if err != sql.ErrNoRows && err != nil {
		return
	}

	// The message to send publicly to the thread
	msg := "Hi <@" + userID.String() + ">!\n" + tryDifferentBrowserMessage

	// The message to use to respond to the interaction
	responseMsg := "The user has been requested to try a different device/browser."

	if err != nil {
		log.Println("Error fetching user for thread:", err)
		msg = tryWithDiscordMessage
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
