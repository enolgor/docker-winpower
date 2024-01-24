package service

type UpsResponse struct {
	Status UpsStatus `json:"status"`
}

type UpsStatus string

const StatusNormal UpsStatus = "Normal"
const StatusACFail UpsStatus = "AC Fail"
const StatusUnknown UpsStatus = "Unknown"
