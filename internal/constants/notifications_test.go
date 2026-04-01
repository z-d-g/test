package constants

import "testing"

func TestNotificationMessage(t *testing.T) {
	tests := []struct {
		nt   NotificationType
		want string
	}{
		{CodeBlockCopied, "Copied code block"},
		{TextCopied, "Copied text"},
		{CutToClipboard, "Cut to clipboard"},
		{PastedFromClipboard, "Pasted"},
		{ExitConfirmation, "Unsaved changes. Quit? (y/n)"},
		{NotificationType("unknown"), "unknown"},
	}
	for _, tt := range tests {
		if got := tt.nt.Message(); got != tt.want {
			t.Errorf("NotificationType(%q).Message() = %q, want %q", tt.nt, got, tt.want)
		}
	}
}
