package service

import (
	"receipt-processor/domain"

	"log"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ReceiptService struct {
	ReceiptsById map[string]domain.Receipt
}

func (s ReceiptService) ProcessReceipt(receipt domain.Receipt) string {
	uuid := uuid.New().String()
	s.ReceiptsById[uuid] = receipt
	return uuid
}

func (s ReceiptService) GetPoints(id string) (int, bool) {
	receipt, ok := s.ReceiptsById[id]
	if ok {
		return getPoints(receipt), true
	} else {
		return 0, false
	}
}

func getPoints(receipt domain.Receipt) int {
	points := 0
	points += scoreRetailerName(receipt.Retailer)
	points += scoreItems(receipt.Items)
	points += scoreTotal(receipt.Total)
	points += scoreTime(receipt.PurchaseDate + " " + receipt.PurchaseTime)
	return points
}

// One point for every alphanumeric character in the retailer name.
func scoreRetailerName(retailer string) int {
	points := 0
	for _, ch := range retailer {
		if isAlphanumeric(ch) {
			points += 1
		}
	}
	return points
}

func isAlphanumeric(ch rune) bool {
	return ('0' <= ch && ch <= '9') || ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

// 5 points for every two items on the receipt.
// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
func scoreItems(items []domain.Item) int {
	points := 0
	points += 5 * (len(items) / 2)
	for _, item := range items {
		price, err := parseFloat(item.Price)
		if err != nil {
			log.Println("Warning: not able to parse price for item.")
			continue
		}
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			a := new(big.Float).Mul(price, big.NewFloat(0.2))
			if !a.IsInt() {
				a = new(big.Float).Add(a, big.NewFloat(1))
			}
			b, _ := a.Int64()
			points += int(b)
		}
	}
	return points
}

// 50 points if the total is a round dollar amount with no cents.
// 25 points if the total is a multiple of 0.25
func scoreTotal(totalStr string) int {
	total, err := parseFloat(totalStr)
	if err != nil {
		log.Println("Warning: not able to parse total amount for receipt.")
	}

	t, _ := new(big.Float).Mul(big.NewFloat(100.0), total).Int64()
	points := 0
	if t%100 == 0 {
		points += 50
	}
	if t%25 == 0 {
		points += 25
	}
	return points
}

func parseFloat(s string) (*big.Float, error) {
	base := 10
	precision := uint(15)
	decimal, _, err := big.ParseFloat(s, base, precision, big.AwayFromZero)
	if err != nil {
		return nil, err
	}
	return decimal, nil
}

// 6 points if the day in the purchase date is odd.
// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
func scoreTime(datetimeString string) int {
	points := 0
	datetime, err := parseDate(datetimeString)
	if err != nil {
		log.Println("Warning: error while parsing date.")
	}
	dateIsOdd := datetime.Day()%2 != 0
	if dateIsOdd {
		points += 6
	}
	startHour := 14
	endHour := 16
	if isInTimeBand(datetime, startHour, endHour) {
		points += 10
	}
	return points
}

func isInTimeBand(datetime time.Time, startHour int, endHour int) bool {
	afterStart := (startHour < datetime.Hour() || (startHour == datetime.Hour() && datetime.Minute() > 0))
	beforeEnd := datetime.Hour() < endHour
	return afterStart && beforeEnd
}

func parseDate(s string) (time.Time, error) {
	layout := "2006-01-02 15:04"
	return time.Parse(layout, s)
}
