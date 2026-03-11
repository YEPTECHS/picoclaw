package tools

import (
	"context"
	"errors"
	"testing"
)

func TestReactTool_Name(t *testing.T) {
	rt := NewReactTool()
	if rt.Name() != "react" {
		t.Errorf("expected name 'react', got %q", rt.Name())
	}
}

func TestReactTool_Parameters(t *testing.T) {
	rt := NewReactTool()
	params := rt.Parameters()
	props, ok := params["properties"].(map[string]any)
	if !ok {
		t.Fatal("expected properties map")
	}
	if _, ok := props["emoji"]; !ok {
		t.Error("expected 'emoji' in properties")
	}
	if _, ok := props["remove"]; !ok {
		t.Error("expected 'remove' in properties")
	}
}

func TestReactTool_Execute_MissingEmoji(t *testing.T) {
	rt := NewReactTool()
	result := rt.Execute(context.Background(), map[string]any{})
	if !result.IsError {
		t.Error("expected error for missing emoji")
	}
}

func TestReactTool_Execute_NoCallback(t *testing.T) {
	rt := NewReactTool()
	ctx := WithToolContext(context.Background(), "slack", "chat-1", "msg-1")
	result := rt.Execute(ctx, map[string]any{"emoji": "eyes"})
	if !result.IsError {
		t.Error("expected error when callback not set")
	}
}

func TestReactTool_Execute_MissingContext(t *testing.T) {
	rt := NewReactTool()
	rt.SetReactCallback(func(_, _, _, _ string, _ bool) error { return nil })

	// No messageID in context
	ctx := WithToolContext(context.Background(), "slack", "chat-1", "")
	result := rt.Execute(ctx, map[string]any{"emoji": "eyes"})
	if !result.IsError {
		t.Error("expected error when messageID missing")
	}
}

func TestReactTool_Execute_AddReaction(t *testing.T) {
	var gotChannel, gotChatID, gotMsgID, gotEmoji string
	var gotRemove bool

	rt := NewReactTool()
	rt.SetReactCallback(func(channel, chatID, messageID, emoji string, remove bool) error {
		gotChannel = channel
		gotChatID = chatID
		gotMsgID = messageID
		gotEmoji = emoji
		gotRemove = remove
		return nil
	})

	ctx := WithToolContext(context.Background(), "slack", "C123/T456", "1234567890.123456")
	result := rt.Execute(ctx, map[string]any{"emoji": "mag"})

	if result.IsError {
		t.Errorf("unexpected error: %s", result.ForLLM)
	}
	if !result.Silent {
		t.Error("expected silent result")
	}
	if gotChannel != "slack" {
		t.Errorf("expected channel 'slack', got %q", gotChannel)
	}
	if gotChatID != "C123/T456" {
		t.Errorf("expected chatID 'C123/T456', got %q", gotChatID)
	}
	if gotMsgID != "1234567890.123456" {
		t.Errorf("expected messageID '1234567890.123456', got %q", gotMsgID)
	}
	if gotEmoji != "mag" {
		t.Errorf("expected emoji 'mag', got %q", gotEmoji)
	}
	if gotRemove {
		t.Error("expected remove=false")
	}
}

func TestReactTool_Execute_RemoveReaction(t *testing.T) {
	var gotRemove bool
	rt := NewReactTool()
	rt.SetReactCallback(func(_, _, _, _ string, remove bool) error {
		gotRemove = remove
		return nil
	})

	ctx := WithToolContext(context.Background(), "slack", "C123", "msg-1")
	result := rt.Execute(ctx, map[string]any{"emoji": "eyes", "remove": true})

	if result.IsError {
		t.Errorf("unexpected error: %s", result.ForLLM)
	}
	if !gotRemove {
		t.Error("expected remove=true")
	}
}

func TestReactTool_Execute_CallbackError(t *testing.T) {
	rt := NewReactTool()
	rt.SetReactCallback(func(_, _, _, _ string, _ bool) error {
		return errors.New("slack API error")
	})

	ctx := WithToolContext(context.Background(), "slack", "C123", "msg-1")
	result := rt.Execute(ctx, map[string]any{"emoji": "eyes"})

	if !result.IsError {
		t.Error("expected error when callback fails")
	}
}
