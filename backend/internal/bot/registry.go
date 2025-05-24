package bot

import (
	"context"
	"encoding/json"

	openai "github.com/openai/openai-go"
)

type Handler func(ctx context.Context, args json.RawMessage) (any, error)

type Spec struct {
	Definition openai.FunctionDefinition
	Handle     Handler
}

var Registry = map[string]Spec{}

func Register(spec Spec) { Registry[spec.Definition.Name] = spec }
