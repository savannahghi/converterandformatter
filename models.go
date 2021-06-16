package converterandformatter

// USSDSessionLog is used to persist a log of USSD sessions
type USSDSessionLog struct {
	MSISDN    string `json:"msisdn" firestore:"msisdn"`
	SessionID string `json:"sessionID" firestore:"sessionID"`
}

//IsEntity ...
func (p PhoneOptIn) IsEntity() {}

// PhoneOptIn is used to persist and manage phone communication whitelists
type PhoneOptIn struct {
	MSISDN  string `json:"msisdn" firestore:"msisdn"`
	OptedIn bool   `json:"optedIn" firestore:"optedIn"`
}

//IsEntity ...
func (p USSDSessionLog) IsEntity() {}
