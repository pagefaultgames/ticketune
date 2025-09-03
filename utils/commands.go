package utils

import (
	"database/sql"
	"log"

	"github.com/amatsagu/tempest"
)

// Base say command functionality reusable by multiple command implementations
// Parameters:
// `itx“: The command interaction to respond to
// `content“: The message content to send to the thread
// `invokerResponse`: The message to send back to the command invoker on success. On error, a relevant error message will be sent instead.
func SayCommandTemplate(itx *tempest.CommandInteraction,
	content string,
	invokerResponse string,
) {
	// Get the user associated with this thread (this handles responding to the interaction on error)
	userID, err := GetUserFromThread(itx)
	if err != sql.ErrNoRows && err != nil {
		return
	}

	// The message to send publicly to the thread
	msg := "Hi <@" + userID.String() + ">!\n" + content

	if err != nil {
		log.Println("Error fetching user for thread:", err)
		msg = content
		invokerResponse = "I couldn't find a user associated with this thread in my database, so I can't ping them." +
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

	itx.SendLinearReply(invokerResponse, true)
}
