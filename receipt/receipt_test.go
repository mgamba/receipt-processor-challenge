package receipt

import (
  "testing"
  "strings"
  "fmt"
)

func jsonToPoints(json string) int {
  receipt, _ := Parse(strings.NewReader(json))
  return receipt.Points()
}

func TestRoundTotal(t *testing.T) {
  receiptFor := func(total string) string {
    return `{
      "retailer": "Target",
      "purchaseDate": "2022-01-01",
      "purchaseTime": "13:01",
      "items": [],
      "total": "` + total + `"
    }`
  }
  roundTotal := "35.00"
  nonRoundTotal := "35.35"

  pointDiff := jsonToPoints(receiptFor(roundTotal)) - jsonToPoints(receiptFor(nonRoundTotal))

  expectedPoints := 50 + 25 // extra 25 for multiple of 25
  if (pointDiff != expectedPoints) {
    t.Errorf("Expected total to be %v more than non-round, but got %v", expectedPoints, pointDiff)
  }
}

func TestMultiple25(t *testing.T) {
  receiptFor := func(totalCents string) string {
    return `{
      "retailer": "Target",
      "purchaseDate": "2022-01-01",
      "purchaseTime": "13:01",
      "items": [],
      "total": "35.` + totalCents + `"
    }`
  }
  testCases := map[string]int{
    "00": 25 + 50, // extra 50 for round dollar amount
    "25": 25,
    "50": 25,
    "75": 25,
  }
  pointsWithout25Bonus := jsonToPoints(receiptFor("13"))

  for centString, expectedPoints := range testCases {
    actualPoints := jsonToPoints(receiptFor(centString)) - pointsWithout25Bonus

    if (actualPoints != expectedPoints) {
      t.Errorf("Expected extra points for \"%s\" cents to be %v, but got %v", centString, expectedPoints, actualPoints)
    }
  }
}

func TestItemPairs(t *testing.T) {
  receiptFor := func(items string) string {
    return `{
      "retailer": "Target",
      "purchaseDate": "2022-01-01",
      "purchaseTime": "13:01",
      "items": [` + items + `],
      "total": "35.00"
    }`
  }
  testCases := map[int]int{
    0: 0,
    1: 0,
    2: 5,
    3: 5,
    4: 10,
  }
  zeroItemPoints := jsonToPoints(receiptFor(""))

  for itemCount, expectedPoints := range testCases {
    items := make([]string, itemCount)
    for i := 0; i < itemCount; i++ {
      items[i] = `{"shortDescription":"d","price":"0.00"}`
    }

    actualPoints := jsonToPoints(receiptFor(strings.Join(items, ","))) - zeroItemPoints

    if (actualPoints != expectedPoints) {
      t.Errorf("Expected extra points for \"%v\" items to be %v, but got %v", itemCount, expectedPoints, actualPoints)
    }
  }
}

func TestDescriptionLength(t *testing.T) {
  receiptFor := func(items string) string {
    return `{
      "retailer": "Target",
      "purchaseDate": "2022-01-01",
      "purchaseTime": "13:01",
      "items": [` + items + `],
      "total": "35.00"
    }`
  }
  zeroItemPoints := jsonToPoints(receiptFor(""))


  testCases := map[string]int{
    // for description with length % 3
    `{"shortDescription":"ddd","price":"21.01"}`: 5,

    // for description with length % 3 and whitespace
    `{"shortDescription":"   ddd\n\n\n","price":"21.01"}`: 5,

    // for description with length not mod 3
    `{"shortDescription":"dddd","price":"21.01"}`: 0,

    // for multiple descriptions
    `{"shortDescription":"dddddd","price":"21.01"},
     {"shortDescription":"  ddd\n","price":"21.01"},
     {"shortDescription":"dddd","price":"21.01"}`: 10 + 5, // extra 5 for item pair bonus
  }

  for items, expectedBonus := range testCases {
    actualPoints := jsonToPoints(receiptFor(items))
    diff := actualPoints - zeroItemPoints
    if (diff != expectedBonus) {
      t.Errorf("Expected %v points for Description but got %v", expectedBonus, diff)
    }
  }
}

func TestOddDays(t *testing.T) {
  receiptFor := func(purchaseDay string) string {
    return `{
      "retailer": "Target",
      "purchaseDate": "2022-01-` + purchaseDay + `",
      "purchaseTime": "13:01",
      "items": [],
      "total": "35.00"
    }`
  }
  testCases := map[int]int{
    0: 0,
    1: 6,
    2: 0,
    3: 6,
    4: 0,
    5: 6,
  }
  zeroItemPoints := jsonToPoints(receiptFor("00"))

  for purchaseDay, expectedPoints := range testCases {
    paddedDay := fmt.Sprintf("%02d", purchaseDay)

    actualPoints := jsonToPoints(receiptFor(paddedDay)) - zeroItemPoints

    if (actualPoints != expectedPoints) {
      t.Errorf("Expected extra points for day \"%v\" to be %v, but got %v", paddedDay, expectedPoints, actualPoints)
    }
  }
}

func TestTimeOfPurchase(t *testing.T) {
  receiptFor := func(purchaseTime string) string {
    return `{
      "retailer": "Target",
      "purchaseDate": "2022-01-01",
      "purchaseTime": "` + purchaseTime + `",
      "items": [],
      "total": "35.00"
    }`
  }
  testCases := map[string]int{
    "03:00": 0,
    "14:00": 10,
    "15:59": 10,
    "16:00": 0,
  }
  zeroItemPoints := jsonToPoints(receiptFor("00:00"))

  for purchaseTime, expectedPoints := range testCases {
    actualPoints := jsonToPoints(receiptFor(purchaseTime)) - zeroItemPoints

    if (actualPoints != expectedPoints) {
      t.Errorf("Expected extra points for time \"%v\" to be %v, but got %v", purchaseTime, expectedPoints, actualPoints)
    }
  }
}
