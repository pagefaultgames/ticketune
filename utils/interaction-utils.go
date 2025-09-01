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

var MissingOptionErr = errors.New("Option is missing")
var WrongTypeErr = errors.New("Option is of the wrong type")

// A constraint that matches all numeric types
type Numeric interface {
	~float32 |
		~float64 |
		~int |
		~int8 |
		~int16 |
		~int32 |
		~int64 |
		~uint |
		~uint8 |
		~uint16 |
		~uint32 |
		~uint64
}

// Convenience function to cast a number of type float64 to a numeric type T
func ConvertFloat64ToNumeric[T Numeric](val float64) (res T) {
	return T(val)
}

// Helper function for a slash command interaction to get an option from a numeric type.
// Returns an error if the option is missing or not an integer
// If `sendReply` is true, sends a reply to the interaction on error
func GetNumericOption[T Numeric](itx *tempest.CommandInteraction, name string, sendReply bool) (res T, err error) {
	val, present := itx.GetOptionValue(name)
	if !present {
		if sendReply {
			itx.SendLinearReply("Error: "+name+" is missing", true)
		}
		return res, MissingOptionErr
	}

	// Try direct type assertion first
	res, ok := val.(T)
	if ok {
		return res, nil
	}
	// If float64, try conversion to numeric type
	if t, ok := any(val).(float64); ok {
		return ConvertFloat64ToNumeric[T](t), nil
	} else if t, ok := any(val).(T); ok {
		return t, nil
	}

	if sendReply {
		itx.SendLinearReply("Error: "+name+" is invalid.", true)
	}
	return res, WrongTypeErr

}

type DiscordOptionResponse interface {
	~float64 | ~string | ~bool
}

// Generic helper function for a slash command interaction to get an option of any type.
func GetOption[T DiscordOptionResponse](itx *tempest.CommandInteraction, name string, sendReply bool) (res T, err error) {
	val, present := itx.GetOptionValue(name)
	if !present {
		if sendReply {
			itx.SendLinearReply("Error: "+name+" is missing", true)
		}
		return res, MissingOptionErr
	}

	// Try direct type assertion first
	if res, ok := val.(T); ok {
		return res, nil
	}
	if sendReply {
		itx.SendLinearReply("Error: "+name+" is invalid.", true)
	}
	return res, WrongTypeErr

}
