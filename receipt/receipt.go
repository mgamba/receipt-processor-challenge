package receipt

import (
  "encoding/json"
	"io"
  "strings"
  "math"
  "regexp"
  "strconv"
)

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price string `json:"price"`
}

type Receipt struct {
	Retailer string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items []Item `json:"items"`
	Total string `json:"total"`
}

func Parse(jsonString io.Reader) (Receipt, error) {
	var receipt Receipt
	dec := json.NewDecoder(jsonString)
	if err := dec.Decode(&receipt); err == nil {
    return receipt, nil
	} else {
    return receipt, err
  }
}

func (r Receipt) Points() (int) {
  return r.retailer_name_points() +
  r.cents_points() +
  r.total_mult_points() +
  r.item_pair_points() +
  r.item_description_points() +
  r.odd_day_points() +
  r.afternoon_points()
}

// One point for every alphanumeric character in the retailer name.
func(r Receipt) retailer_name_points() (int) {
  alphanum := regexp.MustCompile("[a-zA-Z0-9]")
  matches := alphanum.FindAllStringIndex(r.Retailer, -1)
  return len(matches)
}

// 50 points if the total is a round dollar amount with no cents.
func(r Receipt) cents_points() (int) {
  re, _ := regexp.Compile("^\\d+\\.(\\d{2})$")
  cents := re.FindStringSubmatch(r.Total)[0]
  if (cents == "00") {
    return 50
  } else {
    return 0
  }
}

// 25 points if the total is a multiple of 0.25.
func(r Receipt) total_mult_points() (int) {
  re, _ := regexp.Compile("^\\d+\\.(\\d{2})$")
  cents := re.FindStringSubmatch(r.Total)[0]
  switch cents {
  case "00", "25", "50", "75":
    return 25
  default:
    return 0
  }
}

// 5 points for every two items on the receipt.
func(r Receipt) item_pair_points() (int) {
  return 5 * (len(r.Items) / 2)
}

// If the trimmed length of the item description is a multiple of 3,
//  - multiply the price by 0.2
//  - and round up to the nearest integer.
func(r Receipt) item_description_points() (int) {
  result := 0
  multiplier := 0.2
  for _, item := range r.Items {
    if (len(strings.TrimSpace(item.ShortDescription)) % 3 == 0) {
      price, _ := strconv.ParseFloat(item.Price, 64)
      result += int(math.Ceil(price * multiplier))
    }
  }
  return result
}

// 6 points if the day in the purchase date is odd.
func(r Receipt) odd_day_points() (int) {
  day, _ := strconv.Atoi(r.PurchaseDate[len(r.PurchaseDate)-1:])
  return (day % 2) * 6
}

// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
func(r Receipt) afternoon_points() (int) {
  hour := r.PurchaseTime[0:1]
  if (hour == "14" || hour == "15") { // assume 2pm means some time after 2
    return 10
  }
  return 0
}
