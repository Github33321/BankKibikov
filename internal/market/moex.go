package market

import (
	"BankKibikov/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var MoexTickers = []string{
	"SBER", "GAZP", "LKOH", "ROSN", "TATN",
	"VTBR", "GMKN", "ALRS", "CHMF", "POLY",
	"NVTK", "SNGS", "AFLT", "YNDX", "MAGN",
}

func GetMoexPrices(ctx context.Context) ([]models.PriceQuote, error) {
	tickers := strings.Join(MoexTickers, ",")

	url := fmt.Sprintf(
		"https://iss.moex.com/iss/engines/stock/markets/shares/securities.json?securities=%s",
		tickers,
	)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("failed to fetch moex data")
	}

	var raw struct {
		Marketdata struct {
			Columns []string        `json:"columns"`
			Data    [][]interface{} `json:"data"`
		} `json:"marketdata"`
		Securities struct {
			Columns []string        `json:"columns"`
			Data    [][]interface{} `json:"data"`
		} `json:"securities"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	var idxSECID, idxLAST int = -1, -1
	for i, col := range raw.Marketdata.Columns {
		if col == "SECID" {
			idxSECID = i
		}
		if col == "LAST" {
			idxLAST = i
		}
	}

	var idxSECID2, idxNAME int = -1, -1
	for i, col := range raw.Securities.Columns {
		if col == "SECID" {
			idxSECID2 = i
		}
		if col == "SECNAME" {
			idxNAME = i
		}
	}

	if idxSECID == -1 || idxLAST == -1 || idxSECID2 == -1 || idxNAME == -1 {
		return nil, errors.New("unexpected moex response format")
	}

	secNames := map[string]string{}
	for _, row := range raw.Securities.Data {
		secid := row[idxSECID2].(string)
		name := row[idxNAME].(string)
		secNames[secid] = name
	}

	secPrices := map[string]float64{}
	for _, row := range raw.Marketdata.Data {
		secid := row[idxSECID].(string)
		if secPrices[secid] != 0 {
			continue
		}
		if row[idxLAST] != nil {
			if price, ok := row[idxLAST].(float64); ok && price > 0 {
				secPrices[secid] = price
			}
		}
	}

	var quotes []models.PriceQuote
	for _, t := range MoexTickers {
		quotes = append(quotes, models.PriceQuote{
			Name:  fmt.Sprintf("%s (%s)", t, secNames[t]),
			Price: secPrices[t],
		})
	}

	return quotes, nil
}
