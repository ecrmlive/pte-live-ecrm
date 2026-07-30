package auth

import "testing"

func TestParseChannel(t *testing.T) {
	channels := []LoginChannel{
		ChannelWechat, ChannelMiniProgram, ChannelH5, ChannelPC,
		ChannelIOS, ChannelAndroid, ChannelHarmony,
	}
	for _, channel := range channels {
		got, err := ParseChannel(string(channel))
		if err != nil || got != channel {
			t.Fatalf("ParseChannel(%q) = %q, %v", channel, got, err)
		}
	}
	if _, err := ParseChannel("unknown"); err == nil {
		t.Fatal("ParseChannel(unknown) error = nil")
	}
}
