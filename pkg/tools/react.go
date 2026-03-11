package tools

import "context"

// ReactCallback is invoked by the react tool to add/remove emoji reactions.
// channel, chatID, and messageID identify the target message.
// emoji is the reaction name (e.g. "eyes", "mag", "white_check_mark").
// If remove is true, the reaction is removed; otherwise it is added.
type ReactCallback func(channel, chatID, messageID, emoji string, remove bool) error

// ReactTool allows the agent to add or remove emoji reactions on the user's
// message, providing visible progress feedback (e.g. 🔍→📝→✅).
type ReactTool struct {
	reactCallback ReactCallback
}

func NewReactTool() *ReactTool { return &ReactTool{} }

func (t *ReactTool) SetReactCallback(cb ReactCallback) { t.reactCallback = cb }

func (t *ReactTool) Name() string { return "react" }

func (t *ReactTool) Description() string {
	return "Add or remove an emoji reaction on the user's message. Use to signal progress (e.g. eyes, mag, pencil, white_check_mark)."
}

func (t *ReactTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"emoji": map[string]any{
				"type":        "string",
				"description": "Emoji name without colons (e.g. 'mag', 'white_check_mark', 'eyes', 'pencil').",
			},
			"remove": map[string]any{
				"type":        "boolean",
				"description": "If true, remove the reaction instead of adding it. Defaults to false.",
			},
		},
		"required": []string{"emoji"},
	}
}

func (t *ReactTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	emoji, _ := args["emoji"].(string)
	if emoji == "" {
		return ErrorResult("missing required parameter: emoji")
	}

	remove, _ := args["remove"].(bool)

	if t.reactCallback == nil {
		return ErrorResult("react tool not configured: no callback set")
	}

	channel := ToolChannel(ctx)
	chatID := ToolChatID(ctx)
	messageID := ToolMessageID(ctx)

	if channel == "" || chatID == "" || messageID == "" {
		return ErrorResult("react requires channel, chatID, and messageID in context")
	}

	if err := t.reactCallback(channel, chatID, messageID, emoji, remove); err != nil {
		return ErrorResult("react failed: " + err.Error())
	}

	action := "added"
	if remove {
		action = "removed"
	}
	return SilentResult("Reaction :" + emoji + ": " + action)
}
