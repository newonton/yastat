package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/newonton/yastat/internal/timezone"
)

type Client struct {
	APIKey string
	AppID  int
}

type response struct {
	Data struct {
		Points []struct {
			Measures []struct {
				Shows  int     `json:"shows"`
				Reward float64 `json:"partner_wo_nds"`
			} `json:"measures"`
		} `json:"points"`
	} `json:"data"`
}

func (c *Client) Fetch() (int, float64, error) {
	var lastErr error

	for i := range 5 {
		shows, reward, err := c.fetchOnce()
		if err == nil {
			return shows, reward, nil
		}
		lastErr = err
		time.Sleep(time.Second * time.Duration(1<<i))
	}
	return 0, 0, lastErr
}

func (c *Client) fetchOnce() (int, float64, error) {
	req, _ := http.NewRequest("GET",
		"https://partner.yandex.ru/api/statistics2/get.json", nil)

	q := req.URL.Query()
	q.Set("lang", "ru")
	q.Set("dimension_field", "date|day")
	q.Set("period", time.Now().In(timezone.MoscowZone).Format(time.DateOnly))
	q.Set("entity_field", "page_id")
	q.Add("field", "shows")
	q.Add("field", "partner_wo_nds")
	q.Set("filter", fmt.Sprintf("[\"page_id\",\"=\",\"%d\"]", c.AppID))
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "OAuth "+c.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return 0, 0, fmt.Errorf("server error %d", resp.StatusCode)
	}

	var r response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return 0, 0, err
	}

	m := r.Data.Points[0].Measures[0]
	return m.Shows, m.Reward, nil
}
