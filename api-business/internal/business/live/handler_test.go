package live

import "testing"

func TestRoomResponseMasksPlaybackBeforeLiving(t *testing.T) {
	room := roomRow{ID: 7, Title: "七禧直播", Status: "scheduled", PlayURL: "https://play.example/live.m3u8"}
	result := roomResponse(room, nil)
	if result["play_url"] != "" || result["live_status"] != 102 {
		t.Fatalf("scheduled room must not expose playback: %#v", result)
	}
}

func TestRoomResponseExposesPublicPlaybackWhenLiving(t *testing.T) {
	room := roomRow{ID: 8, Title: "七禧直播", Status: "living", PlayURL: "https://play.example/live.m3u8"}
	result := roomResponse(room, nil)
	if result["play_url"] != room.PlayURL || result["live_status"] != 101 {
		t.Fatalf("living room response mismatch: %#v", result)
	}
}
