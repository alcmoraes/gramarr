package handler

import (
	"fmt"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (e *Env) HandleAuth(m *tb.Message) {
	var msg []string
	pass := m.Payload
	user, exists := e.Users.User(m.Sender.ID)

	// Empty Password?
	if pass == "" {
		Send(e.Bot, m.Sender, "Usage: `/auth [password]`")
		return
	}

	// Is User Already Admin?
	if exists && user.IsAdmin() {
		// Notify User
		msg = append(msg, "You're already authorized.")
		msg = append(msg, "Type /start to begin.")
		Send(e.Bot, m.Sender, strings.Join(msg, "\n"))
		return
	}

	// Check if pass is Admin Password
	if pass == e.Config.Bot.AdminPassword {
		if exists {
			user.Access = UAAdmin
			e.Users.Update(user)
		} else {
			newUser := User{
				ID:        m.Sender.ID,
				FirstName: m.Sender.FirstName,
				LastName:  m.Sender.LastName,
				Username:  m.Sender.Username,
				Access:    UAAdmin,
			}
			e.Users.Create(newUser)
		}

		// Notify User
		msg = append(msg, "You have been authorized as an *admin*.")
		msg = append(msg, "Type /start to begin.")
		Send(e.Bot, m.Sender, strings.Join(msg, "\n"))

		// Notify Admin
		adminMsg := fmt.Sprintf("%s has been granted admin access.", DisplayName(m.Sender))
		SendAdmin(e.Bot, e.Users.Admins(), adminMsg)

		return
	}

	// Check if pass is User Password
	if pass == e.Config.Bot.Password {
		if exists {
			// Notify User
			msg = append(msg, "You're already authorized.")
			msg = append(msg, "Type /start to begin.")
			Send(e.Bot, m.Sender, strings.Join(msg, "\n"))
			return
		}
		newUser := User{
			ID:        m.Sender.ID,
			Username:  m.Sender.Username,
			FirstName: m.Sender.FirstName,
			LastName:  m.Sender.LastName,
			Access:    UAMember,
		}
		e.Users.Create(newUser)

		// Notify User
		msg = append(msg, "You have been authorized.")
		msg = append(msg, "Type /start to begin.")
		Send(e.Bot, m.Sender, strings.Join(msg, "\n"))

		// Notify Admin
		adminMsg := fmt.Sprintf("%s has been granted acccess.", DisplayName(m.Sender))
		SendAdmin(e.Bot, e.Users.Admins(), adminMsg)
		return
	}

	// Notify User
	SendError(e.Bot, m.Sender, "Your password is invalid.")

	// Notify Admin
	adminMsg := "%s made an invalid auth request with password: %s"
	adminMsg = fmt.Sprintf(adminMsg, DisplayName(m.Sender), EscapeMarkdown(m.Payload))
	SendAdmin(e.Bot, e.Users.Admins(), adminMsg)
}
