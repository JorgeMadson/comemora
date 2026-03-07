package domain

import (
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

const EventTypeDefault EventType = EventTypeOther

func (e *Event) GetContent() string {
	if e.CustomMessage != "" {
		return e.CustomMessage
	}

	tmpl, ok := defaultMessages[e.Type]
	if !ok {
		tmpl = defaultMessages[EventTypeDefault]
	}

	return replacePlaceholder(tmpl, "{name}", e.Name)
}

func replacePlaceholder(tmpl, placeholder, value string) string {
	return strings.ReplaceAll(tmpl, placeholder, value)
}
