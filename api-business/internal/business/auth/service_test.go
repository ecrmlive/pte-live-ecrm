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

func TestSMSValidationAndHash(t *testing.T) {
	if !validMobile("13800000000") || validMobile("12800000000") || validMobile("1380000000a") {
		t.Fatal("mobile validation mismatch")
	}
	if !validSMSPurpose("login") || !validSMSPurpose("binding") || !validSMSPurpose("reset_password") || !validSMSPurpose("change_mobile") || validSMSPurpose("reset") {
		t.Fatal("sms purpose whitelist mismatch")
	}
	if smsHash("13800000000", "login", "123456") == smsHash("13800000000", "binding", "123456") {
		t.Fatal("purpose must be bound into sms hash")
	}
}
