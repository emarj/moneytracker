package sheet

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/shopspring/decimal"
	"golang.org/x/oauth2/google"
	sheets "google.golang.org/api/sheets/v4"
	"ronche.se/moneytracker/model"
)

type SheetService struct {
	srv     *sheets.Service
	sheetID string
}

func New(authFile string, sheetID string) (*SheetService, error) {
	ctx := context.Background()

	b, err := ioutil.ReadFile(authFile)
	if err != nil {
		return nil, fmt.Errorf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/sheets.googleapis.com-go-quickstart.json
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, fmt.Errorf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)

	srv, err := sheets.New(client)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve Sheets Client %v", err)
	}

	return &SheetService{srv, sheetID}, nil

}

var usrmap = map[string]int{"M": 1, "A": 2}
var usrmaprev = map[string]int{"A": 1, "M": 2}

var empty = "'-"

func (s *SheetService) Insert(e model.Expense) error {
	ctx := context.Background()

	valueInputOption := "USER_ENTERED"
	insertDataOption := "INSERT_ROWS"

	rangeShared := "Comuni!A4"
	rangeUser := "Spese%s!A4"
	rangeu1 := fmt.Sprintf(rangeUser, e.Who)

	//Really bad code. Solution: UserGetAll ?
	rangeu2 := fmt.Sprintf(rangeUser, "A")
	if e.Who == "A" {
		rangeu2 = fmt.Sprintf(rangeUser, "M")
	}

	dateStr := "'" + e.Date.Format("2006-01-02")
	dateCreatedStr := "'" + e.Date.Format("2006-01-02T15:04:05")

	var err error

	if e.Type == 0 {
		u1rows := [][]interface{}{
			{e.UUID, dateCreatedStr, dateStr, e.Who, e.Description, e.Method.Name, e.Amount.Neg().StringFixed(2), e.Shared, e.ShareQuota, e.Category.Name},
		}
		if e.Shared {

			quota := decimal.New(int64(e.ShareQuota), -2)
			amount2 := e.Amount.Mul(quota)
			//amount1 := e.Amount.Sub(amount2)

			u1rows = append(u1rows, []interface{}{e.UUID, dateCreatedStr, dateStr, e.Who, "Storno: " + e.Description, "", amount2.StringFixed(3), e.Shared, e.ShareQuota, e.Category.Name})
			u2rows := [][]interface{}{{e.UUID, dateCreatedStr, dateStr, e.Who, e.Description, "", amount2.Neg().StringFixed(3), e.Shared, e.ShareQuota, e.Category.Name}}

			shrrow := []interface{}{e.UUID, dateCreatedStr, dateStr, e.Who, e.Description, e.Method.Name, e.Amount.Neg().StringFixed(2), e.Shared, e.ShareQuota, e.Category.Name}

			cols := []interface{}{0, 0}

			cols[usrmap[e.Who]-1] = amount2.StringFixed(3)

			shrrow = append(shrrow, cols...)

			_, err = s.srv.Spreadsheets.Values.Append(s.sheetID, rangeShared, &sheets.ValueRange{Values: [][]interface{}{shrrow}}).ValueInputOption(valueInputOption).InsertDataOption(insertDataOption).Context(ctx).Do()
			if err != nil {
				return err
			}

			_, err = s.srv.Spreadsheets.Values.Append(s.sheetID, rangeu2, &sheets.ValueRange{Values: u2rows}).ValueInputOption(valueInputOption).InsertDataOption(insertDataOption).Context(ctx).Do()
			if err != nil {
				return err
			}

		}

		_, err = s.srv.Spreadsheets.Values.Append(s.sheetID, rangeu1, &sheets.ValueRange{Values: u1rows}).ValueInputOption(valueInputOption).InsertDataOption(insertDataOption).Context(ctx).Do()
		if err != nil {
			return err
		}
	}

	if e.Type == 1 {
		shrrow := []interface{}{e.UUID, dateCreatedStr, dateStr, e.Who, e.Description, e.Method.Name, empty, empty, empty, e.Category.Name}
		cols := []interface{}{0, 0}

		cols[usrmaprev[e.Who]-1] = e.Amount.Neg().StringFixed(2)

		shrrow = append(shrrow, cols...)

		_, err = s.srv.Spreadsheets.Values.Append(s.sheetID, rangeShared, &sheets.ValueRange{Values: [][]interface{}{shrrow}}).ValueInputOption(valueInputOption).InsertDataOption(insertDataOption).Context(ctx).Do()
		if err != nil {
			return err
		}

		u1rows := [][]interface{}{
			{e.UUID, dateCreatedStr, dateStr, e.Who, e.Description, e.Method.Name, e.Amount.Neg().StringFixed(2), empty, empty, e.Category.Name},
			{e.UUID, dateCreatedStr, dateStr, e.Who, "Storno: " + e.Description, "", e.Amount.StringFixed(2), empty, empty, e.Category.Name},
		}
		u2rows := [][]interface{}{
			{e.UUID, dateCreatedStr, dateStr, e.Who, e.Description, e.Method.Name, e.Amount.StringFixed(2), empty, empty, e.Category.Name},
			{e.UUID, dateCreatedStr, dateStr, e.Who, "Storno: " + e.Description, "", e.Amount.Neg().StringFixed(2), empty, empty, e.Category.Name},
		}

		_, err = s.srv.Spreadsheets.Values.Append(s.sheetID, rangeu1, &sheets.ValueRange{Values: u1rows}).ValueInputOption(valueInputOption).InsertDataOption(insertDataOption).Context(ctx).Do()
		if err != nil {
			return err
		}

		_, err = s.srv.Spreadsheets.Values.Append(s.sheetID, rangeu2, &sheets.ValueRange{Values: u2rows}).ValueInputOption(valueInputOption).InsertDataOption(insertDataOption).Context(ctx).Do()
		if err != nil {
			return err
		}

	}

	return nil
}
