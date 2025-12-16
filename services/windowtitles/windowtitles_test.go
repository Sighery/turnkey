package windowtitles

import (
	"fmt"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	validTitle := "L:A_N:application_ID:com.lab126.booklet.reader_M:false_PC:N_RC:true_" +
		"WT:true_ASR:true_O:URL_WTNB:true_WTPB:true_DM:N_S:-2"
	testCases := []struct {
		title      KindleTitle
		key        string
		want_value string
		want_found bool
	}{
		{title: KindleTitle(validTitle), key: "O", want_value: "URL", want_found: true},
		{title: KindleTitle(validTitle), key: "nil", want_value: "", want_found: false},
		{title: KindleTitle(""), key: "O", want_value: "", want_found: false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%+v", tc), func(t *testing.T) {
			got_value, got_found := tc.title.Get(tc.key)
			if got_value != tc.want_value || got_found != tc.want_found {
				t.Errorf(
					"Invalid. Want val %s found %t, got val %s found %t",
					tc.want_value, tc.want_found, got_value, got_found,
				)
			}
		})
	}
}

func TestSet(t *testing.T) {
	validTitle := "L:A_N:application_ID:com.lab126.booklet.reader_M:false_PC:N_RC:true_" +
		"WT:true_ASR:true_O:URL_WTNB:true_WTPB:true_DM:N_S:-2"
	testCases := []struct {
		title KindleTitle
		key   string
		value string
		want  KindleTitle
	}{
		{
			title: KindleTitle(validTitle), key: "O", value: "URLD",
			want: KindleTitle(strings.Replace(validTitle, "O:URL", "O:URLD", -1)),
		},
		{
			title: KindleTitle(validTitle), key: "O", value: "",
			want: KindleTitle(strings.Replace(validTitle, "O:URL", "O:", -1)),
		},
		{
			title: KindleTitle(validTitle), key: "S", value: "1",
			want: KindleTitle(strings.Replace(validTitle, "S:-2", "S:1", -1)),
		},
		{
			title: KindleTitle(validTitle), key: "L", value: "B24",
			want: KindleTitle(strings.Replace(validTitle, "L:A", "L:B24", -1)),
		},
		{
			title: KindleTitle(validTitle), key: "nil", value: "", want: KindleTitle(validTitle),
		},
		{
			title: KindleTitle(""), key: "O", value: "URLD", want: KindleTitle(""),
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%+v", tc), func(t *testing.T) {
			got := tc.title.Set(tc.key, tc.value)
			if got != tc.want {
				t.Errorf("Invalid. Want %s, got %s", tc.want, got)
			}
		})
	}
}

func TestIsApplication(t *testing.T) {
	readerTitle := "L:A_N:application_ID:com.lab126.booklet.reader_M:false_PC:N_RC:true_" +
		"WT:true_ASR:true_O:URL_WTNB:true_WTPB:true_DM:N_S:-2"
	footerTitle := "L:D_N:footerBar_AKB:true_HIDET1:200_ID:com.lab126.kppAaMenu_" +
		"A:com.lab126.kppAaMenu_owner:com.lab126.booklet.reader_M:dismissible_" +
		"PAIRID:AaMenuWindow_S:-1_SHOWT1:200"

	testCases := []struct {
		title KindleTitle
		name  string
		want  bool
	}{
		{title: KindleTitle(readerTitle), name: ApplicationReader, want: true},
		{title: KindleTitle(readerTitle), name: "unknown.app", want: false},
		{title: KindleTitle(footerTitle), name: ApplicationReader, want: false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%+v", tc), func(t *testing.T) {
			got := tc.title.IsApplication(tc.name)
			if got != tc.want {
				t.Errorf("Invalid. Want %t, got %t", tc.want, got)
			}
		})
	}
}

// func TestNewKindleTitle(t *testing.T) {
// 	tc := "L:A_N:application_ID:com.lab126.booklet.reader_M:false_PC:N_RC:true_" +
// 		"WT:true_ASR:true_O:URL_WTNB:true_WTPB:true_DM:N_S:-2"
// 	got := NewKindleTitle(tc)

// 	want := KindleTitle{
// 		"L":    "A",
// 		"N":    "application",
// 		"ID":   "com.lab126.booklet.reader",
// 		"M":    "false",
// 		"PC":   "N",
// 		"RC":   "true",
// 		"WT":   "true",
// 		"ASR":  "true",
// 		"O":    "URL",
// 		"WTNB": "true",
// 		"WTPB": "true",
// 		"DM":   "N",
// 		"S":    "-2",
// 	}

// 	if !maps.Equal(got, want) {
// 		t.Errorf("Invalid value, got: %+v, want: %+v", got, want)
// 	}
// }

// func TestNewKindleTitleDropsInvalidComponent(t *testing.T) {
// 	tc := "L:A:C_N:application"
// 	got := NewKindleTitle(tc)

// 	want := KindleTitle{"N": "application"}

// 	if !maps.Equal(got, want) {
// 		t.Errorf("Invalid value, got: %+v, want: %+v", got, want)
// 	}
// }

// func TestToWindowName(t *testing.T) {
// 	tc := KindleTitle{
// 		"L":    "A",
// 		"N":    "application",
// 		"ID":   "com.lab126.booklet.reader",
// 		"M":    "false",
// 		"PC":   "N",
// 		"RC":   "true",
// 		"WT":   "true",
// 		"ASR":  "true",
// 		"O":    "URL",
// 		"WTNB": "true",
// 		"WTPB": "true",
// 		"DM":   "N",
// 		"S":    "-2",
// 	}
// 	got := tc.ToWindowName()

// 	want := "L:A_N:application_ID:com.lab126.booklet.reader_M:false_PC:N_RC:true_" +
// 		"WT:true_ASR:true_O:URL_WTNB:true_WTPB:true_DM:N_S:-2"

// 	if got != want {
// 		t.Errorf("Invalid value, got: %s, want: %s", got, want)
// 	}
// }
