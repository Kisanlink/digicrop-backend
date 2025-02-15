package models

type VoiceRequest struct {
	AudioData string `json:"audio_data"`
}

type VoiceResponse struct {
	Response string `json:"response"`
}
