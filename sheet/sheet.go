package sheet

import (
	"context"
	"fmt"
	"io/ioutil"

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

func (s *SheetService) Insert(e model.Expense) error {
	ctx := context.Background()

	valueInputOption := "RAW"
	insertDataOption := "INSERT_ROWS"

	rangeShared := "Comuni!A2"
	rangeUser := "Spese%s!A2"
	rangeu1 := fmt.Sprintf(rangeUser, e.Who.Name)

	//Really bad code. Solution: UserGetAll ?
	rangeu2 := fmt.Sprintf(rangeUser, "A")
	if e.Who.Name == "A" {
		rangeu2 = fmt.Sprintf(rangeUser, "M")
	}

	u1rows := [][]interface{}{
		{e.UUID, e.Date.Format("02/01/2006"), e.Who.Name, e.Description, e.Method.Name, e.Amount, e.Shared, e.ShareQuota, e.Category.Name},
	}

	if e.Shared {

		amount2 := float64(e.Amount*e.ShareQuota) / float64(100)
		amount1 := amount2 - float64(e.Amount)

		u1rows = append(u1rows, []interface{}{e.UUID, e.Date.Format("02/01/2006"), e.Who.Name, "Storno: " + e.Description, "", amount1, e.Shared, e.ShareQuota, e.Category.Name})
		u2rows := [][]interface{}{{e.UUID, e.Date.Format("02/01/2006"), e.Who.Name, e.Description, "", amount2, e.Shared, e.ShareQuota, e.Category.Name}}

		shrrows := [][]interface{}{
			{e.UUID, e.Date.Format("02/01/2006"), e.Who.Name, e.Description, e.Method.Name, e.Amount, e.Shared, e.ShareQuota, e.Category.Name},
		}

		_, err := s.srv.Spreadsheets.Values.Append(s.sheetID, rangeShared, &sheets.ValueRange{Values: shrrows}).ValueInputOption(valueInputOption).InsertDataOption(insertDataOption).Context(ctx).Do()
		if err != nil {
			return err
		}

		_, err = s.srv.Spreadsheets.Values.Append(s.sheetID, rangeu2, &sheets.ValueRange{Values: u2rows}).ValueInputOption(valueInputOption).InsertDataOption(insertDataOption).Context(ctx).Do()
		if err != nil {
			return err
		}

	}

	_, err := s.srv.Spreadsheets.Values.Append(s.sheetID, rangeu1, &sheets.ValueRange{Values: u1rows}).ValueInputOption(valueInputOption).InsertDataOption(insertDataOption).Context(ctx).Do()
	if err != nil {
		return err
	}

	return nil
}
