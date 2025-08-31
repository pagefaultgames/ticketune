package utils

import (
	"errors"
	ticketune_db "ticketune-bot/ticketune-db"

	"github.com/amatsagu/tempest"
)

// Get the channel and user ID associated with a command interaction
// Errors if the
func GetUserFromThread(itx *tempest.CommandInteraction) (userID tempest.Snowflake, err error) {
	// If this is not a thread in the ticket channel, do nothing
	channel, err := GetChannelFromID(itx.Client, itx.ChannelID)
	if err != nil {
		itx.SendLinearReply("Error fetching channel information", true)
		return
	}

	if !CheckIfPasswordTicketChannel(channel) {
		itx.SendLinearReply("This command can only be used on a password ticket thread", true)
		err = errors.New("not a password ticket thread")
	}

	threadID := itx.ChannelID
	userID, err = ticketune_db.GetDB().GetThreadUser(threadID)

	return
}
