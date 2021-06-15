package converters_and_formatters

const (
	defaultRegion = "KE"

	// OTPCollectionName is the name of the collection used to persist single
	// use verification codes on Firebase
	OTPCollectionName         = "otps"
	PhoneOptInCollectionName  = "phone_opt_ins"
	USSDSessionCollectionName = "ussd_signup_sessions"

	// AuthTokenContextKey is used to add/retrieve the Firebase UID on the context
	AuthTokenContextKey = ContextKey("UID")

	// TestUserEmail is used by integration tests
	TestUserEmail = "be.well@bewell.co.ke"
)
