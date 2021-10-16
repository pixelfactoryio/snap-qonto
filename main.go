package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/pixelfactoryio/goqonto/v2"
)

// AuthTransport structs holds company Slug and  Secret key
type AuthTransport struct {
	*http.Transport
	Slug   string
	Secret string
}

// RoundTrip set "Authorization" header
func (t AuthTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Authorization", fmt.Sprintf("%s:%s", t.Slug, t.Secret))
	return t.Transport.RoundTrip(r)
}

func main() {
	godotenv.Load()
	orgID := os.Getenv("QONTO_ORG_ID")
	IBAN := os.Getenv("QONTO_ORG_IBAN")
	userLogin := os.Getenv("QONTO_USER_LOGIN")
	userSecretKey := os.Getenv("QONTO_SECRET_KEY")
	outputDataPath := os.Getenv("SNAP_OUTPUT_DATA_PATH")

	client := http.Client{
		Transport: AuthTransport{
			&http.Transport{},
			userLogin,
			userSecretKey,
		},
	}

	qonto := goqonto.NewClient(&client)
	ctx := context.Background()

	// Get Organization
	orga, resp, err := qonto.Organizations.Get(ctx, orgID)
	if err != nil && resp.StatusCode != http.StatusOK {
		panic(err.Error())
	}

	// List Transactions
	params := &goqonto.TransactionsOptions{
		Slug: orga.Slug,
		IBAN: IBAN,
	}

	transactions, resp, err := qonto.Transactions.List(ctx, params)
	if err != nil && resp.StatusCode != http.StatusOK {
		panic(err.Error())
	}

	file, _ := json.MarshalIndent(transactions, "", " ")
	err = ioutil.WriteFile(outputDataPath, file, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func prettyPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}
