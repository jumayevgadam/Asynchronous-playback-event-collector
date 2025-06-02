package model

type PlaybackEvent struct {
	ID          int
	UserID      string
	VideoID     string
	StartAt     int64
	StopAt      int64
	BitrateKbps int
	DeviceType  string
	ErrorCode   *string
	Region      string
	Model       string
}
