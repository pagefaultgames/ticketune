package constants

import (
	"log"

	"github.com/amatsagu/tempest"
	"github.com/joho/godotenv"
)

var (
	HELPER_ROLE_ID                 tempest.Snowflake
	TICKET_CHANNEL_ID              tempest.Snowflake
	SUPPORT_CATEGORY_ID            tempest.Snowflake
	BOT_TROUBLESHOOTING_CHANNEL_ID tempest.Snowflake
	DISCORD_GUILD_ID               tempest.Snowflake
)

// "I couldn't find a user associated with this thread in my database, so I can't ping them...."
var COULD_NOT_FIND_USER_TO_PING string = "I couldn't find a user associated with this thread in my database, so I can't ping them." +
	"However, I've sent the requested message to the thread."

// Initialize the constants from environment variables
func init() {
	log.Println("Loading environment variables...")
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln("failed to load env variables", err)
	}
	// This function can be used to initialize constants if needed in the future
	var err error
	HELPER_ROLE_ID, err = tempest.EnvToSnowflake("DISCORD_HELPER_ROLE_ID")
	if err != nil {
		log.Fatalln("failed to parse HELPER_ROLE_ID variable to snowflake", err)
	}
	TICKET_CHANNEL_ID, err = tempest.EnvToSnowflake("TICKET_CHANNEL_ID")
	if err != nil {
		log.Fatalln("failed to parse TICKET_CHANNEL_ID variable to snowflake", err)
	}
	SUPPORT_CATEGORY_ID, err = tempest.EnvToSnowflake("SUPPORT_TICKET_CATEGORY_ID")
	if err != nil {
		log.Fatalln("failed to parse SUPPORT_CATEGORY_ID variable to snowflake", err)
	}
	BOT_TROUBLESHOOTING_CHANNEL_ID, err = tempest.EnvToSnowflake("BOT_TROUBLESHOOTING_CHANNEL_ID")
	if err != nil {
		log.Fatalln("failed to parse BOT_TROUBLESHOOTING_CHANNEL_ID variable to snowflake", err)
	}
	DISCORD_GUILD_ID, err = tempest.EnvToSnowflake("DISCORD_GUILD_ID")
	if err != nil {
		log.Fatalln("failed to parse DISCORD_GUILD_ID variable to snowflake", err)
	}
}
