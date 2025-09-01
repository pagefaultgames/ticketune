package commands

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/pagefaultgames/ticketune/constants"
	"github.com/pagefaultgames/ticketune/db"
	"github.com/pagefaultgames/ticketune/utils"

	"github.com/amatsagu/tempest"
)

var closeCommandDescription = "Close the current password ticket thread and remove the associated user's permission overrides"

var CloseCommand = tempest.Command{
	Name:                "close",
	Description:         closeCommandDescription,
	RequiredPermissions: tempest.ADMINISTRATOR_PERMISSION_FLAG,
	SlashCommandHandler: closeTicketCommandImpl,
	Contexts:            []tempest.InteractionContextType{tempest.GUILD_CONTEXT_TYPE},
}

func closeTicketCommandImpl(itx *tempest.CommandInteraction) {
	// If this is not a thread in the ticket channel, do nothing
	channel, err := utils.GetChannelFromID(itx.Client, itx.ChannelID)
	if err != nil {
		log.Println("Error fetching channel info:", err)
		//return // will be caught by next check
	}

	// ParentID is the ID of the parent channel for threads, or the category ID for channels
	// If ParentID is not the ticket channel ID, this is not a valid ticket thread
	if channel.ParentID != constants.TICKET_CHANNEL_ID || channel.Type != tempest.GUILD_PRIVATE_THREAD_CHANNEL_TYPE {
		itx.SendLinearReply("This command can only be used on a password ticket thread", true)
		return
	}

	user, err := db.Get().CloseThread(itx.ChannelID)
	if err == sql.ErrNoRows {
		// If no rows were returned, tell the initiator of the commands.
		itx.SendLinearReply("Error: I couldn't find a user associated with this thread in my database. You'll have to close the thread manually.", true)
		return
	}

	// Delete the channel permissions for the user
	err = deleteChannelPermissionForUser(itx.Client, user)
	if err != nil {
		log.Println("Error deleting channel permission for user:", err)
		itx.SendLinearReply("Error: I couldn't remove the user's permissions to access this thread. You'll have to close the thread manually.", true)
		return
	}

	// Delete the thread
	_, err = itx.Client.Rest.Request(
		http.MethodDelete,
		fmt.Sprintf("/channels/%d", itx.ChannelID),
		nil,
	)
	if err != nil {
		itx.SendLinearReply(
			fmt.Sprintf("I removed the user's access to the ticket, but ran into an error deleting the thread: %s.",
				err.Error()),
			// Not ephemeral so that a Helper can show a dev what went wrong
			false,
		)
	}
}

// Remove permission overrides for the user in the ticket channel
func deleteChannelPermissionForUser(client *tempest.Client, userID tempest.Snowflake) error {
	_, err := client.Rest.Request(
		http.MethodDelete,
		fmt.Sprintf("/channels/%d/permissions/%d", constants.TICKET_CHANNEL_ID, userID),
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}
