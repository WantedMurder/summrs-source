package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/WantedMurder/summrs-source/utils"
)

func (cmd *Commands) Ban(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	member, err := s.GuildMember(m.GuildID, m.Mentions[0].ID)
	if err != nil {
		return
	}
	var (
		userRole   = utils.HighestRole(s, m.GuildID, m.Member)
		targetRole = utils.HighestRole(s, m.GuildID, member)
	)

	if !utils.IsAbove(userRole, targetRole) && m.Author.ID != utils.GetGuildOwner(s, m.GuildID) {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<:X_Mark:857046727530250271> | You do not have the proper permissions to ban <@%s>", m.Mentions[0].ID))
		return
	}

	err = s.GuildBanCreateWithReason(m.GuildID, m.Mentions[0].ID, fmt.Sprintf("%s | Command ban", s.State.User.Username), 0)

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<:X_Mark:857046727530250271> | Could not ban <@%s>", m.Mentions[0].ID))
		return
	}

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Description: fmt.Sprintf("<:V_Mark:857046695577518090> | Successfully banned <@%s> from the server.", m.Mentions[0].ID),
		Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("banned by: %s", m.Author.Username)},
		Color:       0x36393F,
	})
}

func (cmd *Commands) Kick(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	member, err := s.GuildMember(m.GuildID, m.Mentions[0].ID)
	if err != nil {
		return
	}

	var (
		userRole   = utils.HighestRole(s, m.GuildID, m.Member)
		targetRole = utils.HighestRole(s, m.GuildID, member)
	)

	if !utils.IsAbove(userRole, targetRole) && m.Author.ID != utils.GetGuildOwner(s, m.GuildID) {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<:X_Mark:857046727530250271> | You do not have the proper permissions to kick <@%s>", m.Mentions[0].ID))
		return
	}

	err = s.GuildMemberDeleteWithReason(m.GuildID, m.Mentions[0].ID, fmt.Sprintf("%s | Command kick", s.State.User.Username))

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<:X_Mark:857046727530250271> | Could not kick <@%s>", m.Mentions[0].ID))
		return
	}
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Description: fmt.Sprintf("<:V_Mark:857046695577518090> | Successfully kicked <@%s> from the server.", m.Mentions[0].ID),
		Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Kicked by: %s", m.Author.Username)},
		Color:       0x36393F,
	})
}

func (cmd *Commands) Lockdown(s *discordgo.Session, m *discordgo.Message, ctx *Context) {

	err := s.ChannelPermissionSet(m.ChannelID, m.GuildID, "role", 0, discordgo.PermissionSendMessages)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "<:X_Mark:857046727530250271> | Couldn't lock this channel. Maybe check perms?")
		return
	}
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Description: fmt.Sprintf("<:V_Mark:857046695577518090> | Successfully locked <#%s>.", m.ChannelID),
		Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Locked by: %s", m.Author.Username)},
		Color:       0x36393F,
	})
}

func (cmd *Commands) Purge(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	amount, err := strconv.Atoi(ctx.Fields[0])
	if err != nil || amount > 100 {
		s.ChannelMessageSend(m.ChannelID, "<:X_Mark:857046727530250271> | You have to specify a number and it has to be inbetween 1-100")
		return
	}
	var tempMessages []string

	messages, err := s.ChannelMessages(m.ChannelID, amount, m.ID, "", "")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "<:X_Mark:857046727530250271> | Could not fetch channel messages")
		return
	}

	for _, message := range messages {
		tempMessages = append(tempMessages, message.ID)
	}
	tempMessages = append(tempMessages, m.ID)

	err = s.ChannelMessagesBulkDelete(m.ChannelID, tempMessages)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Could not purge: %s", err.Error()))
	}
}

func (cmd *Commands) SlowMode(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	seconds, err := strconv.Atoi(ctx.Fields[0])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "You have to specify a number.")
		return
	}

	channel, err := s.ChannelEditComplex(m.ChannelID, &discordgo.ChannelEdit{
		RateLimitPerUser: seconds,
	})

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "<:X_Mark:857046727530250271> | Couldn't set slowmode on this channel. Maybe check perms?")
		return
	}
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Description: fmt.Sprintf("<:V_Mark:857046695577518090> | Set <#%s> slowmode to %d seconds", channel.ID, seconds),
		Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Slowmode set by: %s", m.Author.Username)},
		Color:       0x36393F,
	})
}

func (cmd *Commands) Unban(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	var (
		nukebot    int
		unbancount int
	)

	bans, err := s.GuildBans(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "<:X_Mark:857046727530250271> | Could not fetch the guild bans")
		return
	}

	for _, ban := range bans {
		if !strings.Contains(ban.Reason, s.State.User.Username) {
			s.GuildBanDelete(m.GuildID, ban.User.ID)
			unbancount++
			continue
		}
		nukebot++
	}
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Description: fmt.Sprintf("<:V_Mark:857046695577518090> | Unbanned %d members | <:X_Mark:857046727530250271> | Didn't unban %d account(s) banned by %s.", unbancount, nukebot, s.State.User.Username),
		Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Mass-Unban started by: %s", m.Author.Username)},
		Color:       0x36393F,
	})
}

func (cmd *Commands) UnLockdown(s *discordgo.Session, m *discordgo.Message, ctx *Context) {

	err := s.ChannelPermissionSet(m.ChannelID, m.GuildID, "role", discordgo.PermissionSendMessages, discordgo.PermissionMentionEveryone|discordgo.PermissionSendTTSMessages)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "<:X_Mark:857046727530250271> | Couldn't unlock this channel. Maybe check perms?")
		return
	}
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Description: fmt.Sprintf("<:V_Mark:857046695577518090> | Successfully unlocked <#%s>.", m.ChannelID),
		Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Unlocked by: %s", m.Author.Username)},
		Color:       0x36393F,
	})
}

func (cmd *Commands) UnSlowMode(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	channel, err := s.ChannelEditComplex(m.ChannelID, &discordgo.ChannelEdit{
		RateLimitPerUser: 0,
	})

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "<:X_Mark:857046727530250271> | Couldn't turn off slowmode on this channel. Maybe check perms?")
		return
	}
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Description: fmt.Sprintf("<:V_Mark:857046695577518090> | Successfully disabled slowmode for <#%s>.", channel.ID),
		Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Slowmode turned off by: %s", m.Author.Username)},
		Color:       0x36393F,
	})
}
