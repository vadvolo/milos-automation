package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var (
	metricurl = "https://monitoring.api.cloud.yandex.net/monitoring/v2/data/write?folderId=b1guh823ahdh649bbdv2&service=custom"
	iamurl    = "https://iam.api.cloud.yandex.net/iam/v1/tokens"
	token     = os.Getenv("YC_TOKEN")
)

type MetricPost struct {
	Timestamp time.Time `json:"ts"`
	Labels    any       `json:"labels"`
	Metrics   []Metric  `json:"metrics"`
}

func NewMetricPost() (*MetricPost, error) {
	ts := time.Now()
	fmt.Println(ts.Format(time.RFC3339))
	labels := map[string]string{
		"object": "milos",
	}
	return &MetricPost{
		Timestamp: ts,
		Labels:    labels,
		Metrics: []Metric{
			Metric{
				Name:      "TesName",
				Labels:    labels,
				Type:      "DGAUGE",
				Timestamp: ts,
				Value:     "123.4",
				Timeseries: []Tseries{
					Tseries{
						Timestamp: ts,
						Value:     "234.5",
					},
				},
			},
		},
	}, nil
}

func (m *MetricPost) Post() error {
	s, _ := json.MarshalIndent(m, "", "  ")
	fmt.Println(string(s))
	body, err := json.Marshal(m)
	if err != nil {
		return err
	}
	r, err := http.NewRequest("POST", metricurl, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	r.Header.Add("Content-Type", "application/json")
	iam_token, err := GetYandexCloudIAMToken()
	if err != nil {
		return err
	}
	r.Header.Add("Authorization", "Bearer " + iam_token)
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return err
	}
	fmt.Println(res)
	defer res.Body.Close()

	post := &MetricResponse{}
	derr := json.NewDecoder(res.Body).Decode(post)
	if derr != nil {
		return err
	}

	if res.StatusCode != http.StatusCreated {
		return err
	}

	fmt.Println(post.WrittenMetricsCount, post.ErrorMessage)

	return nil
}

type Metric struct {
	Name       string    `json:"name"`
	Labels     any       `json:"labels"`
	Type       string    `json:"type"`
	Timestamp  time.Time `json:"ts"`
	Value      string    `json:"value"`
	Timeseries []Tseries `json:"timeseries"`
}

type Tseries struct {
	Timestamp time.Time `json:"ts"`
	Value     string    `json:"value"`
}

type MetricResponse struct {
	WrittenMetricsCount string `json:"writtenMetricsCount"`
	ErrorMessage        string `json:"errorMessage"`
}

func GetYandexCloudIAMToken() (string, error) {
	body, err := json.Marshal(IAMRequest{
		Token: token,
	})
	if err != nil {
		return "", err
	}
	r, err := http.NewRequest("POST", iamurl, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	r.Header.Add("Content-Type", "application/json")
	r.ParseForm()	

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))

	post := &IAMResponse{}
	json.Unmarshal(b, post)

	fmt.Println(post.IamToken, post.ExpiresAt)

	

	return post.IamToken, nil

	// return post.IamToken, nil
}

type IAMRequest struct {
	Token string `json:"yandexPassportOauthToken"`
}

type IAMResponse struct {
	IamToken  string `json:"iamToken"`
	ExpiresAt string `json:"expiresAt"`
}
