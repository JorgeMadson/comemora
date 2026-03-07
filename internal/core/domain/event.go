package domain

import (
	"errors"
	"strings"
	"time"
)

type EventType string

const (
	EventTypeBirthday    EventType = "aniversario"
	EventTypeWedding     EventType = "casamento"
	EventTypeDating      EventType = "namoro"
	EventTypePet         EventType = "pet"
	EventTypeWork        EventType = "trabalho"
	EventTypeBereavement EventType = "luto"
	EventTypeOther       EventType = "outro"
)

type NotificationChannel string

const (
	ChannelEmail    NotificationChannel = "Email"
	ChannelTeams    NotificationChannel = "Teams"
	ChannelWhatsApp NotificationChannel = "WhatsApp"
	ChannelSMS      NotificationChannel = "SMS"
	ChannelTelegram NotificationChannel = "Telegram"
	ChannelDiscord  NotificationChannel = "Discord"
)

type Event struct {
	ID                 uint                `json:"id" gorm:"primaryKey"`
	Name               string              `json:"name"`
	Day                int                 `json:"day"`
	Month              int                 `json:"month"`
	Year               int                 `json:"year,omitempty"` // 0 if unknown/irrelevant
	Type               EventType           `json:"type"`
	Tags               []string            `json:"tags" gorm:"serializer:json"` // Requires GORM JSON serializer
	PreferredChannel   NotificationChannel `json:"preferred_channel"`
	ContactDestination string              `json:"contact_destination"`
	CustomMessage      string              `json:"custom_message"`
	IsImportant        bool                `json:"is_important"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
}

var defaultMessages = map[EventType]string{
	EventTypeBirthday:    "Feliz Aniversário! {name} 🎈🍰🎁 Que seja um dia inesquecível e o início de um novo ano na sua vida cheio de felicidade e muitas realizações. 😘❤️",
	EventTypeWedding:     "Feliz Aniversário de Casamento! {name} 💍💑 Que o amor de vocês continue crescendo a cada dia!",
	EventTypeWork:        "Parabéns pelo tempo de casa! {name} 🚀 Obrigado por fazer parte da nossa jornada.",
	EventTypeDating:      "Feliz dia do nosso amor! {name} ❤️",
	EventTypePet:         "Parabéns para o nosso pet querido! {name} 🐾🦴",
	EventTypeBereavement: "Hoje lembramos com carinho de {name}. 🖤",
	EventTypeOther:       "Olá {name}, hoje é um dia especial! 🎉",
}

func (e *Event) GetContent() string {
	if e.CustomMessage != "" {
		return e.CustomMessage
	}

	tmpl, ok := defaultMessages[e.Type]
	if !ok {
		tmpl = defaultMessages[EventTypeOther]
	}

	return strings.ReplaceAll(tmpl, "{name}", e.Name)
}

func (e EventType) IsValid() bool {
	_, ok := defaultMessages[e]
	return ok
}

func (c NotificationChannel) IsValid() bool {
	switch c {
	case ChannelEmail, ChannelTeams, ChannelWhatsApp, ChannelSMS, ChannelTelegram, ChannelDiscord:
		return true
	}
	return false
}

func (e *Event) Validate() error {
	if strings.TrimSpace(e.Name) == "" {
		return errors.New("name is required")
	}
	if e.Day < 1 || e.Day > 31 {
		return errors.New("day must be between 1 and 31")
	}
	if e.Month < 1 || e.Month > 12 {
		return errors.New("month must be between 1 and 12")
	}
	if !e.Type.IsValid() {
		return errors.New("invalid event type: must be one of aniversario, casamento, namoro, pet, trabalho, luto, outro")
	}
	if e.PreferredChannel != "" && !e.PreferredChannel.IsValid() {
		return errors.New("invalid preferred_channel: must be one of Email, Teams, WhatsApp, SMS, Telegram, Discord")
	}
	if strings.TrimSpace(e.ContactDestination) == "" {
		return errors.New("contact_destination is required")
	}
	return nil
}

