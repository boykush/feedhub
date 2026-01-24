package domain

import (
	"errors"
	"time"
)

// Category represents the category of a transaction
type Category int

const (
	CategoryUnspecified Category = iota
	CategoryIncome
	CategoryFood
	CategoryDailyNecessities
	CategoryHobby
	CategorySocial
	CategoryTransportation
	CategoryClothing
	CategoryHealth
	CategoryAutomobile
	CategoryEducation
	CategorySpecial
	CategoryCashCard
	CategoryUtilities
	CategoryCommunication
	CategoryHousing
	CategoryTaxSocialSecurity
	CategoryInsurance
	CategoryOther
)

// Money represents an amount of money with high precision
type Money struct {
	units int64
	nanos int32
}

// NewMoney creates a new Money value with validation
func NewMoney(units int64, nanos int32) (Money, error) {
	if nanos < -999999999 || nanos > 999999999 {
		return Money{}, errors.New("nanos must be between -999999999 and 999999999")
	}
	// Sign must match
	if units > 0 && nanos < 0 {
		return Money{}, errors.New("sign of units and nanos must match")
	}
	if units < 0 && nanos > 0 {
		return Money{}, errors.New("sign of units and nanos must match")
	}
	return Money{units: units, nanos: nanos}, nil
}

// Units returns the whole units of the amount
func (m Money) Units() int64 {
	return m.units
}

// Nanos returns the nano units of the amount
func (m Money) Nanos() int32 {
	return m.nanos
}

// Transaction represents a single financial transaction
type Transaction struct {
	id       string
	date     time.Time
	title    string
	amount   Money
	category Category
}

// NewTransaction creates a new Transaction with validation
func NewTransaction(id string, date time.Time, title string, amount Money, category Category) (Transaction, error) {
	if id == "" {
		return Transaction{}, errors.New("id is required")
	}
	if title == "" {
		return Transaction{}, errors.New("title is required")
	}
	return Transaction{
		id:       id,
		date:     date,
		title:    title,
		amount:   amount,
		category: category,
	}, nil
}

// ID returns the transaction ID
func (t Transaction) ID() string {
	return t.id
}

// Date returns the transaction date
func (t Transaction) Date() time.Time {
	return t.date
}

// Title returns the transaction title
func (t Transaction) Title() string {
	return t.title
}

// Amount returns the transaction amount
func (t Transaction) Amount() Money {
	return t.amount
}

// Category returns the transaction category
func (t Transaction) Category() Category {
	return t.category
}
