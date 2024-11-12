package utils

import (
	"reflect"
	"testing"
)

func TestParseIRCMessage(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		args    args
		want    *ircMessage
		wantErr bool
	}{
		{
			name: "USERNOTICE message with submysterygift tags",
			args: args{line: "@badge-info=subscriber/25;badges=subscriber/3012,sub-gifter/550;color=#FF7F50;display-name=Bent_Knee;emotes=;flags=;id=954b764f-0b73-440f-9528-818fba5eb210;login=bent_knee;mod=0;msg-id=submysterygift;msg-param-community-gift-id=5596808143494217874;msg-param-mass-gift-count=20;msg-param-origin-id=5596808143494217874;msg-param-sender-count=570;msg-param-sub-plan=1000;room-id=204730616;subscriber=1;system-msg=Bent_Knee\\sis\\sgifting\\s20\\sTier\\s1\\sSubs\\sto\\sCinna's\\scommunity!\\sThey've\\sgifted\\sa\\stotal\\sof\\s570\\sin\\sthe\\schannel!;tmi-sent-ts=1730308163912;user-id=183568249;user-type=;vip=0 :tmi.twitch.tv USERNOTICE #cinna"},
			want: &ircMessage{
				Raw: "@badge-info=subscriber/5;badges=subscriber/3,premium/1;color=#FF4700;display-name=Grimm_Star333;emotes=;flags=;id=9bdbf18e-60e7-41c6-bdf8-70d764d5e93e;login=grimm_star333;mod=0;msg-id=submysterygift;msg-param-community-gift-id=8679973651393933659;msg-param-mass-gift-count=5;msg-param-origin-id=8679973651393933659;msg-param-sender-count=5;msg-param-sub-plan=1000;room-id=13240194;subscriber=1;system-msg=Grimm_Star333\\sis\\sgifting\\s5\\sTier\\s1\\sSubs\\sto\\sscump's\\scommunity!\\sThey've\\sgifted\\sa\\stotal\\sof\\s5\\sin\\sthe\\schannel!;tmi-sent-ts=1730317076163;user-id=44603698;user-type=;vip=0 :tmi.twitch.tv USERNOTICE #scump",
				Tags: map[string]string{
					"badge-info":                  "subscriber/5",
					"badges":                      "subscriber/3,premium/1",
					"color":                       "#FF4700",
					"display-name":                "Grimm_Star333",
					"emotes":                      "",
					"flags":                       "",
					"id":                          "9bdbf18e-60e7-41c6-bdf8-70d764d5e93e",
					"login":                       "grimm_star333",
					"mod":                         "0",
					"msg-id":                      "submysterygift",
					"msg-param-community-gift-id": "8679973651393933659",
					"msg-param-mass-gift-count":   "5",
					"msg-param-origin-id":         "8679973651393933659",
					"msg-param-sender-count":      "5",
					"msg-param-sub-plan":          "1000",
					"room-id":                     "13240194",
					"subscriber":                  "1",
					"system-msg":                  "Grimm_Star333 is gifting 5 Tier 1 Subs to scump's community! They've gifted a total of 5 in the channel!",
					"tmi-sent-ts":                 "1730317076163",
					"user-id":                     "44603698",
					"user-type":                   "",
					"vip":                         "0",
				},
				Source: ircMessageSource{
					Nickname: "", Username: "", Host: "tmi.twitch.tv"},
				Command: "USERNOTICE",
				Params:  []string{"#scump"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseIRCMessage(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseIRCMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseIRCMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
