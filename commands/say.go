package commands

import (
	"log"

	"github.com/amatsagu/tempest"
	"github.com/pagefaultgames/ticketune/types"
	"github.com/pagefaultgames/ticketune/utils"
)

const sayCommandDescription = "Have Ticketune say a message in the current thread, optionally pinging the user."

var SayCommand = tempest.Command{
	Name:                "say",
	Description:         sayCommandDescription,
	SlashCommandHandler: sayCommandImpl,
	RequiredPermissions: tempest.ADMINISTRATOR_PERMISSION_FLAG,
	Options: []tempest.CommandOption{
		{
			Type:        tempest.STRING_OPTION_TYPE,
			Name:        "message",
			Description: "The message to send. Supports Discord markdown; pings to other users are intentionally suppressed",
			Required:    true,
			MinLength:   3,
			MaxLength:   1900,
		},
		{
			Type:        tempest.BOOLEAN_OPTION_TYPE,
			Name:        "no-ping",
			Description: "Do not ping the user associated with this ticket. Defaults to false (ping the user).",
			Required:    false,
		},
	},
}

func sayCommandImpl(itx *tempest.CommandInteraction) {
	// Error can be discarded, as the argument is optional, and we default to `false`
	noPing, _ := utils.GetOption[bool](itx, "no-ping", false)
	message, err := utils.GetOption[string](itx, "message", true)

	// GetOption already handles responding to the interaction on error
	if err != nil {
		return
	}
	var responseMsg string = "Your message has been sent to the thread."

	var messageParams types.CreateMessageParams = types.CreateMessageParams{}

	userID, err := utils.GetUserFromThread(itx)
	// These errors are already handled in GetUserFromThread
	if err != nil && (err == utils.ErrNotATicketThread || err != utils.ErrCantFetchChannel) {
		// An error occurred that was not "not a ticket thread" or "no such thread"
		return
	}

	if !noPing {
		if err == nil {
			message = "Hi <@" + userID.String() + ">!\n" + message
			messageParams.AllowedMentions = &tempest.AllowedMentions{Users: []tempest.Snowflake{userID}}
		} else {
			log.Println("Error fetching user for thread:", err)
			responseMsg = "I couldn't find a user associated with this thread in my database, so I can't ping them." +
				"However, I've sent the message to the thread."
		}
	}

	messageParams.Content = message

	_, err = utils.SendDiscordMessage(itx.Client, itx.ChannelID, messageParams, nil, true)
	if err != nil {
		itx.SendLinearReply("Error sending message to thread: "+err.Error(), true)
		return
	}

	itx.SendLinearReply(responseMsg, true)

}
