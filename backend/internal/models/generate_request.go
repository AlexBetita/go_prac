package models

import "mime/multipart"

type GenerateRequest struct {
	Message       string                  `json:"message"   validate:"required"`
	InteractionID *string                 `json:"interaction_id,omitempty"`
	Model         string                  `json:"model,omitempty"`   // overrides interaction/default
	Plugins       []string                `json:"plugins,omitempty"` // optional function names
	Stream        bool                    `json:"stream,omitempty"`
	SystemPrompt  *string                 `json:"system_prompt,omitempty"` // overrides interaction/default
	Files         []*multipart.FileHeader `json:"-"`                       // filled by middleware
}
