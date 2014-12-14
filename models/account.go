package models

import (
	"github.com/gyepisam/mcf"
	_ "github.com/gyepisam/mcf/scrypt" // Required to have mcf hash the password into scrypt
)

// Account stores essential data for a Lavaboom user, and is thus not encrypted.
type Account struct {
	Resource

	// Billing is a struct containing billing information.
	// TODO Work in progress
	Billing BillingData `json:"billing" gorethink:"billing"`

	// Password is the password used to login to the account.
	// It's hashed and salted using scrypt.
	Password string `json:"-"  gorethink:"password"`

	// PGPKey is the fingerprint of account's default key
	PGPKey string `json:"pgp_key" gorethink:"pgp_key"`

	// Settings contains data needed to customize the user experience.
	// TODO Work in progress
	Settings SettingsData `json:"settings" gorethink:"settings"`

	// Type is the account type.
	// Examples (work in progress):
	//		* beta: while in beta these are full accounts; after beta, these are normal accounts with special privileges
	//		* std: standard, free account
	//		* premium: premium account
	//		* superuser: Lavaboom staff
	Type string `json:"type" gorethink:"type"`

	AltEmail string `json:"alt_email" gorethink:"alt_email"`

	Status string `json:"status" gorethink:"status"`
}

// SetPassword changes the account's password
func (a *Account) SetPassword(password string) error {
	encrypted, err := mcf.Create(password)
	if err != nil {
		return err
	}

	a.Password = encrypted
	return nil
}

// VerifyPassword checks if password is valid and upgrades it if its encrypting scheme was outdated
// Returns isValid, wasUpdated, error
func (a *Account) VerifyPassword(password string) (bool, bool, error) {
	isValid, err := mcf.Verify(password, a.Password)
	if err != nil {
		return false, false, err
	}

	if !isValid {
		return false, false, nil
	}

	isCurrent, err := mcf.IsCurrent(a.Password)
	if err != nil {
		return false, false, err
	}

	if !isCurrent {
		err := a.SetPassword(password)
		if err != nil {
			return true, false, err
		}

		a.Touch()
		return true, true, nil
	}

	return true, false, nil
}

// SettingsData TODO
type SettingsData struct {
}

// BillingData TODO
type BillingData struct {
}
