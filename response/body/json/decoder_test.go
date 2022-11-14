package json_test

import (
	"bytes"
	stdJSON "encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/deelawn/go-http/response/body/json"
)

type SoundOfMusic struct {
	Doe string `json:"doe,omitempty"`
	Re  string `json:"re,omitempty"`
	Mi  string `json:"mi,omitempty"`
	Fa  string `json:"fa,omitempty"`
	So  string `json:"so,omitempty"`
	La  string `json:"la,omitempty"`
	Ti  string `json:"ti,omitempty"`
	Do  string `json:"do,omitempty"`
}

func (s SoundOfMusic) String() string {
	return fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s", s.Doe, s.Re, s.Mi, s.Fa, s.So, s.La, s.Ti, s.Do)
}

func newSoundOfMusicReader() io.Reader {

	som := SoundOfMusic{
		Doe: "a deer a female deer",
		Re:  "a drop of golden sun",
		Mi:  "a name i call myself",
		Fa:  "a long long way to run",
		So:  "a needle pulling thread",
		La:  "a note to follow so",
		Ti:  "a drink with jam and bread",
		Do:  "which brings us back to do. o. o. o.",
	}

	bSOM, _ := stdJSON.Marshal(som)
	return bytes.NewBuffer(bSOM)
}

func TestDecoder_Decode(t *testing.T) {

	tests := []struct {
		name           string
		source         io.Reader
		target         any
		expTargetValue SoundOfMusic
		expErrText     string
	}{
		{
			name:       "nil source",
			expErrText: json.ErrNilDecodeSource.Error(),
		},
		{
			name:       "nil target",
			source:     newSoundOfMusicReader(),
			expErrText: "json: Unmarshal(nil)",
		},
		{
			name:       "non writable target",
			source:     newSoundOfMusicReader(),
			target:     SoundOfMusic{},
			expErrText: "json: Unmarshal(non-pointer json_test.SoundOfMusic)",
		},
		{
			name:   "okay",
			source: newSoundOfMusicReader(),
			target: new(SoundOfMusic),
			expTargetValue: SoundOfMusic{
				Doe: "a deer a female deer",
				Re:  "a drop of golden sun",
				Mi:  "a name i call myself",
				Fa:  "a long long way to run",
				So:  "a needle pulling thread",
				La:  "a note to follow so",
				Ti:  "a drink with jam and bread",
				Do:  "which brings us back to do. o. o. o.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var errText string

			decoder := new(json.Decoder)
			err := decoder.Decode(tt.source, tt.target)
			if err != nil {
				errText = err.Error()
			}

			if errText != tt.expErrText {
				t.Errorf("wanted error %s, got %s", tt.expErrText, errText)
				return
			} else if tt.expErrText != "" {
				return
			}

			if tt.target.(fmt.Stringer).String() != tt.expTargetValue.String() {
				t.Errorf("wanted target %s, got %s", tt.expTargetValue, tt.target)
			}
		})
	}
}
