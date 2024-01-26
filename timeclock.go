package freee

import (
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// https://developer.freee.co.jp/reference/hr/reference#operations-tag-タイムレコーダー(打刻)

type TimeClock struct {
	ID               int       `json:"id"`
	Date             string    `json:"date"`
	Type             string    `json:"type"`
	Datetime         time.Time `json:"datetime"`
	OriginalDatetime time.Time `json:"original_datetime"`
	Note             string    `json:"note"`
}

type ListTimeClocksOps struct {
	FromDate string // 取得する打刻期間の開始日(YYYY-MM-DD)(例:2018-08-01)(デフォルト: 当月の打刻開始日)
	ToDate   string // 取得する打刻期間の終了日(YYYY-MM-DD)(例:2018-08-31)(デフォルト: 当日)
	Limit    int    // 取得レコードの件数 (デフォルト: 50, 最小: 1, 最大: 100)
	Offset   int    // 取得レコードのオフセット (デフォルト: 0)
}

// ListTimeClocks は指定した従業員・期間の打刻情報を返します。
func (c *Client) ListTimeClocks(companyID int, employeeID int, opts *ListTimeClocksOps) ([]TimeClock, error) {
	u := "https://api.freee.co.jp/hr/api/v1/employees/" + url.PathEscape(strconv.Itoa(employeeID)) + "/time_clocks"
	q := url.Values{
		"company_id": {strconv.Itoa(companyID)},
	}
	if opts != nil {
		if opts.FromDate != "" {
			q.Set("from_date", opts.FromDate)
		}
		if opts.ToDate != "" {
			q.Set("to_date", opts.ToDate)
		}
		if opts.Limit != 50 {
			q.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Offset > 0 {
			q.Set("offset", strconv.Itoa(opts.Limit))
		}
	}
	resp, err := c.do(http.MethodGet, u, q, nil)
	if err != nil {
		return nil, err
	}

	timeClocks := []TimeClock{}
	if err := resp.Parse(&timeClocks); err != nil {
		return nil, err
	}

	return timeClocks, nil
}

// GetTimeClock は指定した従業員・指定した打刻の詳細情報を返します。
func (c *Client) GetTimeClock(companyID int, employeeID int, timeClockID int) (TimeClock, error) {
	u := "https://api.freee.co.jp/hr/api/v1/employees/" + url.PathEscape(strconv.Itoa(employeeID)) + "/time_clocks/" + url.PathEscape(strconv.Itoa(timeClockID))
	q := url.Values{
		"company_id": {strconv.Itoa(companyID)},
	}
	resp, err := c.do(http.MethodGet, u, q, nil)
	if err != nil {
		return TimeClock{}, err
	}

	result := struct {
		EmployeeTimeClock TimeClock `json:"employee_time_clock"`
	}{}
	if err := resp.Parse(&result); err != nil {
		return TimeClock{}, err
	}

	return result.EmployeeTimeClock, nil
}

type AvailableTypes struct {
	AvailableTypes []string `json:"available_types"`
	BaseDate       string   `json:"base_date"`
}

type GetAvailableTypesOpts struct {
	Date string // 従業員情報を取得したい年月日(YYYY-MM-DD)(例:2018-08-01)(デフォルト：当日)
}

// GetAvailableTypes は指定した従業員・日付の打刻可能種別と打刻基準日を返します。
// 例: すでに出勤した状態だと、休憩開始、退勤が配列で返ります。
func (c *Client) GetAvailableTypes(companyID int, employeeID int, opts *GetAvailableTypesOpts) (AvailableTypes, error) {
	u := "https://api.freee.co.jp/hr/api/v1/employees/" + url.PathEscape(strconv.Itoa(employeeID)) + "/time_clocks/available_types"
	q := url.Values{
		"company_id": {strconv.Itoa(companyID)},
	}
	if opts != nil {
		if opts.Date != "" {
			q.Set("date", opts.Date)
		}
	}
	resp, err := c.do(http.MethodGet, u, q, nil)
	if err != nil {
		return AvailableTypes{}, err
	}

	availableTypes := AvailableTypes{}
	if err := resp.Parse(&availableTypes); err != nil {
		return AvailableTypes{}, err
	}

	return availableTypes, nil
}

type CreateTimeClockRequest struct {
	CompanyID int     `json:"company_id"`
	Type      string  `json:"type"`
	BaseDate  *string `json:"base_date,omitempty"`
	Datetime  *string `json:"datetime,omitempty"`
}

// CreateTimeClock は指定した従業員の打刻情報を登録します。
// 注意点
// - 休憩開始の連続や退勤のみなど、整合性の取れていない打刻は登録できません。 打刻可能種別の取得APIを呼ぶことで、その従業員がその時点で登録可能な打刻種別が取得できます。
// - 出勤の打刻は
//   - 前日の出勤時刻から24時間以内の場合、前日の退勤打刻が必須です。
//   - 前日の出勤時刻から24時間経過している場合は、前日の退勤打刻がなくとも出勤打刻を登録することができます。
//
// - 退勤の打刻は
//   - 『退勤を自動打刻する』の設定を使用している場合は、出勤打刻から24時間経過しても退勤打刻がない場合に、退勤打刻が自動で登録されます。
//   - すでに登録されている退勤打刻よりも後の時刻であれば上書き登録することができます。
//
// - 打刻が日をまたぐ場合は、base_date(打刻日)に前日の日付を指定してください。
// - datetime(打刻日時)を指定できるのは管理者か事務担当者の権限を持ったユーザーのみです。
func (c *Client) CreateTimeClock(companyID int, employeeID int, request *CreateTimeClockRequest) (TimeClock, error) {
	u := "https://api.freee.co.jp/hr/api/v1/employees/" + url.PathEscape(strconv.Itoa(employeeID)) + "/time_clocks"
	resp, err := c.do(http.MethodPost, u, nil, request)
	if err != nil {
		return TimeClock{}, err
	}

	result := struct {
		EmployeeTimeClock TimeClock `json:"employee_time_clock"`
	}{}
	if err := resp.Parse(&result); err != nil {
		return TimeClock{}, err
	}

	return result.EmployeeTimeClock, nil
}
