package converterandformatter_test

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	uuid "github.com/kevinburke/go.uuid"
	converterandformatter "github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/serverutils"
	"github.com/stretchr/testify/assert"
)

// CoverageThreshold sets the test coverage threshold below which the tests will fail
const CoverageThreshold = 0.75

func TestMain(m *testing.M) {
	os.Setenv("MESSAGE_KEY", "this-is-a-test-key$$$")
	os.Setenv("ENVIRONMENT", "staging")
	err := os.Setenv("ROOT_COLLECTION_SUFFIX", "staging")
	if err != nil {
		if serverutils.IsDebug() {
			log.Printf("can't set root collection suffix in env: %s", err)
		}
		os.Exit(-1)
	}
	existingDebug, err := serverutils.GetEnvVar("DEBUG")
	if err != nil {
		existingDebug = "false"
	}

	os.Setenv("DEBUG", "true")

	rc := m.Run()
	// Restore DEBUG envar to original value after running test
	os.Setenv("DEBUG", existingDebug)

	// rc 0 means we've passed,
	// and CoverMode will be non empty if run with -cover
	if rc == 0 && testing.CoverMode() != "" {
		c := testing.Coverage()
		if c < CoverageThreshold {
			fmt.Println("Tests passed but coverage failed at", c)
			rc = -1
		}
	}

	os.Exit(rc)
}

func TestIsMSISDNValid(t *testing.T) {

	tests := []struct {
		name   string
		msisdn string
		want   bool
	}{
		{
			name:   "valid : kenyan with code",
			msisdn: "+254722000000",
			want:   true,
		},
		{
			name:   "valid : kenyan without code",
			msisdn: "0722000000",
			want:   true,
		},
		{
			name:   "valid : kenyan without code and spaces",
			msisdn: "0722 000 000",
			want:   true,
		},
		{
			name:   "valid : kenyan without plus sign",
			msisdn: "+254722000000",
			want:   true,
		},
		{
			name:   "valid : kenyan without plus sign and spaces",
			msisdn: "+254 722 000 000",
			want:   true,
		},
		{
			name:   "invalid : kenyan with alphanumeric1",
			msisdn: "+25472abc0000",
			want:   false,
		},
		{
			name:   "invalid : kenyan with alphanumeric2",
			msisdn: "072abc0000",
			want:   false,
		},
		{
			name:   "invalid : kenyan short length",
			msisdn: "0720000",
			want:   false,
		},
		{
			name:   "invalid : kenyan with unwanted characters : asterisk",
			msisdn: "072*120000",
			want:   false,
		},
		{
			name:   "invalid : kenyan without code with plus sign as prefix",
			msisdn: "+0722000000",
			want:   false,
		},
		{
			name:   "ivalid : international with alphanumeric",
			msisdn: "90191919qwe",
			want:   false,
		},
		{
			name:   "invalid : international with unwanted characters : asterisk",
			msisdn: "(+351) 282 *3 50 50",
			want:   false,
		},
		{
			name:   "invalid : international with unwanted characters : assorted",
			msisdn: "(+351) $82 *3 50 50",
			want:   false,
		},
		{
			name:   "valid : usa number",
			msisdn: "+12028569601",
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := converterandformatter.IsMSISDNValid(tt.msisdn); got != tt.want {
				t.Errorf("IsMSISDNValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalizeMSISDN(t *testing.T) {
	type args struct {
		msisdn string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "good Kenyan number, full E164 format",
			args: args{
				"+254723002959",
			},
			want:    "+254723002959",
			wantErr: false,
		},
		{
			name: "good Kenyan number, no + prefix",
			args: args{
				"254723002959",
			},
			want:    "+254723002959",
			wantErr: false,
		},
		{
			name: "good Kenyan number, no international dialling code",
			args: args{
				"0723002959",
			},
			want:    "+254723002959",
			wantErr: false,
		},
		{
			name: "good US number, full E164 format",
			args: args{
				"+16125409037",
			},
			want:    "+16125409037",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converterandformatter.NormalizeMSISDN(tt.args.msisdn)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeMSISDN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if *got != tt.want {
				t.Errorf("NormalizeMSISDN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateMSISDN(t *testing.T) {
	fc := &firebasetools.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	assert.Nil(t, err)

	ctx := firebasetools.GetAuthenticatedContext(t)
	firestoreClient, err := firebaseApp.Firestore(ctx)
	assert.Nil(t, err)

	otpMsisdn := "+254722000000"
	normalized, err := converterandformatter.NormalizeMSISDN(otpMsisdn)
	assert.Nil(t, err)

	validOtpCode := rand.Int()
	validOtpData := map[string]interface{}{
		"authorizationCode": strconv.Itoa(validOtpCode),
		"isValid":           true,
		"message":           "testing OTP message",
		"msisdn":            normalized,
		"timestamp":         time.Now(),
	}
	_, err = firebasetools.SaveDataToFirestore(firestoreClient, firebasetools.SuffixCollection(converterandformatter.OTPCollectionName), validOtpData)
	assert.Nil(t, err)

	invalidOtpCode := rand.Int()
	invalidOtpData := map[string]interface{}{
		"authorizationCode": strconv.Itoa(invalidOtpCode),
		"isValid":           false,
		"message":           "testing OTP message",
		"msisdn":            normalized,
		"timestamp":         time.Now(),
	}
	_, err = firebasetools.SaveDataToFirestore(firestoreClient, firebasetools.SuffixCollection(converterandformatter.OTPCollectionName), invalidOtpData)
	assert.Nil(t, err)

	type args struct {
		msisdn           string
		verificationCode string
		isUSSD           bool
		firestoreClient  *firestore.Client
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "invalid phone format",
			args: args{
				msisdn: "not a valid phone format",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "ussd session validation",
			args: args{
				msisdn:           "0722000000",
				verificationCode: uuid.NewV1().String(),
				isUSSD:           true,
				firestoreClient:  firestoreClient,
			},
			want:    "+254722000000",
			wantErr: false,
		},
		{
			name: "non existent verification code for non USSD",
			args: args{
				msisdn:           "0722000000",
				verificationCode: uuid.NewV1().String(),
				isUSSD:           false,
				firestoreClient:  firestoreClient,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "valid verification code for non USSD",
			args: args{
				msisdn:           "0722000000",
				verificationCode: strconv.Itoa(validOtpCode),
				isUSSD:           false,
				firestoreClient:  firestoreClient,
			},
			want:    "+254722000000",
			wantErr: false,
		},
		{
			name: "used (invalid) verification code for non USSD",
			args: args{
				msisdn:           "0722000000",
				verificationCode: strconv.Itoa(invalidOtpCode),
				isUSSD:           false,
				firestoreClient:  firestoreClient,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converterandformatter.ValidateMSISDN(tt.args.msisdn, tt.args.verificationCode, tt.args.isUSSD, tt.args.firestoreClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMSISDN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateMSISDN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateAndSaveMSISDN(t *testing.T) {
	fc, _ := firebasetools.GetFirestoreClient(context.Background())

	type args struct {
		msisdn           string
		verificationCode string
		isUSSD           bool
		optIn            bool
		firestoreClient  *firestore.Client
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "invalid phone number/OTP, non USSD",
			args: args{
				msisdn:           "0722000000",
				verificationCode: "not a real one",
				isUSSD:           false,
				optIn:            true,
				firestoreClient:  fc,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "valid phone number, USSD, opt in true",
			args: args{
				msisdn:           "0722000000",
				verificationCode: "this is a ussd session ID from the telco",
				isUSSD:           true,
				optIn:            true,
				firestoreClient:  fc,
			},
			want:    "+254722000000",
			wantErr: false,
		},
		{
			name: "valid phone number, USSD, opt in false",
			args: args{
				msisdn:           "0722000000",
				verificationCode: "this is a ussd session ID from the telco",
				isUSSD:           true,
				optIn:            false,
				firestoreClient:  fc,
			},
			want:    "+254722000000",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converterandformatter.ValidateAndSaveMSISDN(tt.args.msisdn, tt.args.verificationCode, tt.args.isUSSD, tt.args.optIn, tt.args.firestoreClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAndSaveMSISDN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateAndSaveMSISDN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringSliceContains(t *testing.T) {
	type args struct {
		s []string
		e string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "string found in slice",
			args: args{
				s: []string{"a", "b", "c", "d", "e"},
				e: "a",
			},
			want: true,
		},
		{
			name: "string not found in slice",
			args: args{
				s: []string{"a", "b", "c", "d", "e"},
				e: "z",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := converterandformatter.StringSliceContains(tt.args.s, tt.args.e); got != tt.want {
				t.Errorf("StringSliceContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSliceContains(t *testing.T) {
	type args struct {
		s []int
		e int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "slice which contains the int",
			args: args{
				s: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
				e: 7,
			},
			want: true,
		},
		{
			name: "slice which does NOT contain the int",
			args: args{
				s: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
				e: 79,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := converterandformatter.IntSliceContains(tt.args.s, tt.args.e); got != tt.want {
				t.Errorf("IntSliceContains() = %v, want %v", got, tt.want)
			}
		})
	}
}
