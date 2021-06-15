package converters_and_formatters

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/savannahghi/firebasetools"
	"github.com/ttacon/libphonenumber"
)

// IsMSISDNValid uses regular expression to validate the a phone number
func IsMSISDNValid(msisdn string) bool {
	if len(msisdn) < 10 {
		return false
	}
	reKen := regexp.MustCompile(`^(?:254|\+254|0)?((7|1)(?:(?:[129][0-9])|(?:0[0-8])|(4[0-1]))[0-9]{6})$`)
	re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	if !reKen.MatchString(msisdn) {
		return re.MatchString(msisdn)
	}
	return reKen.MatchString(msisdn)
}

// NormalizeMSISDN validates the input phone number.
// For valid phone numbers, it normalizes them to international format
// e.g +2547........
func NormalizeMSISDN(msisdn string) (*string, error) {
	if !IsMSISDNValid(msisdn) {
		return nil, fmt.Errorf("invalid phone number: %s", msisdn)
	}
	num, err := libphonenumber.Parse(msisdn, defaultRegion)
	if err != nil {
		return nil, err
	}
	formatted := libphonenumber.Format(num, libphonenumber.INTERNATIONAL)
	cleaned := strings.ReplaceAll(formatted, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	return &cleaned, nil
}

// ValidateMSISDN returns an error if the MSISDN format is wrong or the
// supplied verification code is not valid

// Deprecated: Should implement `VerifyOTP` instead. This helps to confirm if a phonenumber
// is valid by verifying the code sent to it.
func ValidateMSISDN(
	msisdn, verificationCode string,
	isUSSD bool, firestoreClient *firestore.Client) (string, error) {

	// check the format
	normalized, err := NormalizeMSISDN(msisdn)
	if err != nil {
		return "", fmt.Errorf("invalid phone format: %v", err)
	}

	// save a USSD log for USSD registrations
	if isUSSD {
		log := USSDSessionLog{
			MSISDN:    msisdn,
			SessionID: verificationCode,
		}
		_, err = firebasetools.SaveDataToFirestore(
			firestoreClient, firebasetools.SuffixCollection(USSDSessionCollectionName), log)
		if err != nil {
			return "", fmt.Errorf("unable to save USSD session: %v", err)
		}
		return *normalized, nil
	}

	// check if the OTP is on file / known
	query := firestoreClient.Collection(firebasetools.SuffixCollection(OTPCollectionName)).Where(
		"isValid", "==", true,
	).Where(
		"msisdn", "==", normalized,
	).Where(
		"authorizationCode", "==", verificationCode,
	)
	ctx := context.Background()
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve verification codes: %v", err)
	}
	if len(docs) == 0 {
		return "", fmt.Errorf("no matching verification codes found")
	}

	for _, doc := range docs {
		otpData := doc.Data()
		otpData["isValid"] = false
		err = firebasetools.UpdateRecordOnFirestore(
			firestoreClient, firebasetools.SuffixCollection(OTPCollectionName), doc.Ref.ID, otpData)
		if err != nil {
			return "", fmt.Errorf("unable to save updated OTP document: %v", err)
		}
	}

	return *normalized, nil
}

// ValidateAndSaveMSISDN returns an error if the MSISDN format is wrong or the
// supplied verification code is not valid
func ValidateAndSaveMSISDN(
	msisdn, verificationCode string, isUSSD bool, optIn bool,
	firestoreClient *firestore.Client) (string, error) {
	validated, err := ValidateMSISDN(
		msisdn, verificationCode, isUSSD, firestoreClient)
	if err != nil {
		return "", fmt.Errorf("invalid MSISDN: %s", err)
	}
	if optIn {
		data := PhoneOptIn{
			MSISDN:  validated,
			OptedIn: optIn,
		}
		_, err = firebasetools.SaveDataToFirestore(
			firestoreClient, PhoneOptInCollectionName, data)
		if err != nil {
			return "", fmt.Errorf("unable to save email opt in: %v", err)
		}
	}
	return validated, nil
}

// StringSliceContains tests if a string is contained in a slice of strings
func StringSliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
