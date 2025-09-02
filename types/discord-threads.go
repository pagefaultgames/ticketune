/*
 * SPDX-FileCopyrightText: 2025 Pagefault Games
 * SPDX-FileContributor: SirzBenjie
 *
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package types

import (
	"time"

	"github.com/amatsagu/tempest"
)

type ArchiveDuration int // Values for AutoArchiveDuration

const (
	AUTO_ARCHIVE_HOUR ArchiveDuration = 60
	AUTO_ARCHIVE_DAY  ArchiveDuration = 1440
	AUTO_ARCHIVE_3DAY ArchiveDuration = 4320
	AUTO_ARCHIVE_WEEK ArchiveDuration = 10080
)

type ThreadType int // Values for Channel.Type when creating threads

const (
	THREAD_TYPE_PUBLIC  ThreadType = 11
	THREAD_TYPE_PRIVATE ThreadType = 12
)

type CreateThreadWithoutMessageParams struct {
	Name                string          `json:"name"`                            // 1-100 character thread name
	AutoArchiveDuration ArchiveDuration `json:"auto_archive_duration,omitempty"` // Duration in minutes to automatically archive the thread after recent activity, can be set to: 60, 1440, 4320, 10080
	Type                ThreadType      `json:"type,omitempty"`                  // The type of thread to create
	RateLimitPerUser    int             `json:"rate_limit_per_user,omitempty"`   // Amount of seconds a user has to wait before sending another message (0-21600)
	Invitable           bool            `json:"invitable"`                       // Whether non-moderators can add other non-moderators to a thread; only available when creating a private thread, and defaults to true if omitted
}

// Channel represents a Discord channel or thread object (partial, for threads).
// https://discord.com/developers/docs/resources/channel#channel-object-channel-structure
type Channel struct {
	ID                         tempest.Snowflake   `json:"id"`                                      // the id of this channel
	Type                       tempest.ChannelType `json:"type"`                                    // the type of channel
	GuildID                    tempest.Snowflake   `json:"guild_id,omitempty"`                      // the id of the guild (may be missing for some channel objects received over gateway guild dispatches)
	Name                       string              `json:"name"`                                    // The name of the channel (1-100 characters).
	ParentID                   tempest.Snowflake   `json:"parent_id,omitempty"`                     // for guild channels: id of the parent category for a channel (each parent category can contain up to 50 channels), for threads: id of the text channel this thread was created
	ThreadMetadata             *ThreadMetadata     `json:"thread_metadata,omitempty"`               // thread-specific fields not needed by other channels
	DefaultAutoArchiveDuration int                 `json:"default_auto_archive_duration,omitempty"` // default duration, copied onto newly created threads, in minutes, threads will stop showing in the channel list after the specified period of inactivity, can be set to: 60, 1440, 4320, 10080
	Flags                      int                 `json:"flags,omitempty"`                         // https://discord.com/developers/docs/resources/channel#channel-object-channel-flags
	// MemberCount                int                 `json:"member_count,omitempty"`
	// TotalMessageSent           int                 `json:"total_message_sent,omitempty"`
	// LastPinTimestamp           string              `json:"last_pin_timestamp,omitempty"`
	// OwnerID                    string              `json:"owner_id,omitempty"`
	// Position                   int                 `json:"position,omitempty"`
	// PermissionOverwrites []struct {
	//  ID    string `json:"id"`
	//  Type  int    `json:"type"`
	//  Allow string `json:"allow"`
	//  Deny  string `json:"deny"`
	// } `json:"permission_overwrites,omitempty"`
	// RTCRegion                  string          `json:"rtc_region,omitempty"`
	// MessageCount               int             `json:"message_count,omitempty"`
	// Permissions                string          `json:"permissions,omitempty"`
	// Bitrate                    int             `json:"bitrate,omitempty"` // for voice channels only
	// NSFW                       bool            `json:"nsfw,omitempty"`
	// Unused fields, omitted
	// Member                     *ThreadMember   `json:"member,omitempty"`
	// VideoQualityMode           int             `json:"video_quality_mode,omitempty"`
}

func (c *Channel) Mention() string {
	return "<#" + c.ID.String() + ">"
}

// ThreadMetadata represents thread-specific metadata.
type ThreadMetadata struct {
	Archived            bool            `json:"archived"`              // whether the thread is archived
	AutoArchiveDuration ArchiveDuration `json:"auto_archive_duration"` // the thread will stop showing in the channel list after `auto_archive_duration` minutes of inactivity, can be set to: 60, 1440, 4320, 10080
	ArchiveTimestamp    time.Time       `json:"archive_timestamp"`     // timestamp when the thread's archive status was last changed, used for calculating recent activity
	Locked              bool            `json:"locked"`                // whether the thread is locked; when a thread is locked, only users with MANAGE_THREADS can unarchive it
	Invitable           bool            `json:"invitable,omitempty"`   // whether non-moderators can add other non-moderators to a thread; only available on private threads
	CreateTimestamp     time.Time       `json:"create_timestamp"`      // timestamp when the thread was created; only populated for threads created after 2022-01-09
}

// ThreadMember represents a thread member object.
type ThreadMember struct {
	ID            tempest.Snowflake `json:"id,omitempty"`      // ID of the thread
	UserID        tempest.Snowflake `json:"user_id,omitempty"` // ID of the thread
	JoinTimestamp time.Time         `json:"join_timestamp"`    // Time the user last joined the thread
	// Flags         int               `json:"flags"`             // Any user-thread settings, currently only used for notifications
}

type ROLE_OR_MEMBER uint8 // Values for the `type` field in EditChannelPermissionsParams

const (
	ROLE_TYPE   ROLE_OR_MEMBER = 0
	MEMBER_TYPE ROLE_OR_MEMBER = 1
)

type EditChannelPermissionsParams struct {
	Allow tempest.PermissionFlags `json:"allow,string,omitempty"` // bitwise value of allowed permissions
	Deny  tempest.PermissionFlags `json:"deny,string,omitempty"`  // bitwise value of denied permissions
	Type  ROLE_OR_MEMBER          `json:"type"`                   // 1 for role or 2 for member
}
