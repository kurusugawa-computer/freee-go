package freee

import (
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// https://developer.freee.co.jp/reference/hr/reference#operations-tag-勤怠

type WorkRecord struct {
	BreakRecords []struct {
		ClockInAt  time.Time `json:"clock_in_at"`
		ClockOutAt time.Time `json:"clock_out_at"`
	} `json:"break_records"`
	ClockInAt                                 *time.Time `json:"clock_in_at"`
	ClockOutAt                                *time.Time `json:"clock_out_at"`
	Date                                      time.Time  `json:"date"`
	DayPattern                                string     `json:"day_pattern"`
	SchedulePattern                           string     `json:"schedule_pattern"`
	EarlyLeavingMins                          int        `json:"early_leaving_mins"`
	HalfPaidHolidayMins                       int        `json:"half_paid_holiday_mins"`
	HalfSpecialHolidayMins                    int        `json:"half_special_holiday_mins"`
	HourlyPaidHolidayMins                     int        `json:"hourly_paid_holiday_mins"`
	HourlySpecialHolidayMins                  int        `json:"hourly_special_holiday_mins"`
	IsAbsence                                 bool       `json:"is_absence"`
	IsEditable                                bool       `json:"is_editable"`
	LatenessMins                              int        `json:"lateness_mins"`
	NormalWorkClockInAt                       *time.Time `json:"normal_work_clock_in_at"`
	NormalWorkClockOutAt                      *time.Time `json:"normal_work_clock_out_at"`
	NormalWorkMins                            int        `json:"normal_work_mins"`
	Note                                      string     `json:"note"`
	PaidHoliday                               int        `json:"paid_holiday"`
	SpecialHoliday                            int        `json:"special_holiday"`
	SpecialHolidaySettingID                   *int       `json:"special_holiday_setting_id"`
	UseAttendanceDeduction                    bool       `json:"use_attendance_deduction"`
	UseDefaultWorkPattern                     bool       `json:"use_default_work_pattern"`
	UseHalfCompensatoryHoliday                bool       `json:"use_half_compensatory_holiday"`
	TotalOvertimeWorkMins                     int        `json:"total_overtime_work_mins"`
	TotalHolidayWorkMins                      int        `json:"total_holiday_work_mins"`
	TotalLatenightWorkMins                    int        `json:"total_latenight_work_mins"`
	NotAutoCalcWorkTime                       bool       `json:"not_auto_calc_work_time"`
	TotalExcessStatutoryWorkMins              int        `json:"total_excess_statutory_work_mins"`
	TotalLatenightExcessStatutoryWorkMins     int        `json:"total_latenight_excess_statutory_work_mins"`
	TotalOvertimeExceptNormalWorkMins         int        `json:"total_overtime_except_normal_work_mins"`
	TotalLatenightOvertimeExceptNormalWorkMin int        `json:"total_latenight_overtime_except_normal_work_min"`
}

// DeleteWorkRecord は指定した従業員の勤怠情報を削除します。
func (c *Client) DeleteWorkRecord(employeeID int, date string) error {
	u := "https://api.freee.co.jp/hr/api/v1/employees/" + url.PathEscape(strconv.Itoa(employeeID)) + "/work_records/" + url.PathEscape(date)
	resp, err := c.do(http.MethodDelete, u, nil, nil)
	if err != nil {
		return err
	}
	resp.Close()
	return nil
}

// GetWorkRecord は指定した従業員・日付の勤怠情報を返します。
func (c *Client) GetWorkRecord(employeeID int, date string) (WorkRecord, error) {
	u := "https://api.freee.co.jp/hr/api/v1/employees/" + url.PathEscape(strconv.Itoa(employeeID)) + "/work_records/" + url.PathEscape(date)
	resp, err := c.do(http.MethodGet, u, nil, nil)
	if err != nil {
		return WorkRecord{}, err
	}

	workRecord := WorkRecord{}
	if err := resp.Parse(&workRecord); err != nil {
		return WorkRecord{}, err
	}

	return workRecord, nil
}

type PutWorkRecordRequest struct {
	CompanyID    int `json:"company_id"`
	BreakRecords []struct {
		ClockInAt  time.Time `json:"clock_in_at"`
		ClockOutAt time.Time `json:"clock_out_at"`
	} `json:"break_records,omitempty"`
	ClockInAt                *time.Time `json:"clock_in_at,omitempty"`
	ClockOutAt               *time.Time `json:"clock_out_at,omitempty"`
	DayPattern               *string    `json:"day_pattern,omitempty"`
	EarlyLeavingMins         *int       `json:"early_leaving_mins,omitempty"`
	IsAbsence                *bool      `json:"is_absence,omitempty"`
	LatenessMins             *int       `json:"lateness_mins,omitempty"`
	NormalWorkClockInAt      *time.Time `json:"normal_work_clock_in_at,omitempty"`
	NormalWorkClockOutAt     *time.Time `json:"normal_work_clock_out_at,omitempty"`
	NormalWorkMins           *int       `json:"normal_work_mins,omitempty"`
	Note                     *string    `json:"note,omitempty"`
	PaidHoliday              *string    `json:"paid_holiday,omitempty"`
	HalfPaidHolidayMins      *int       `json:"half_paid_holiday_mins,omitempty"`
	HourlyPaidHolidayMins    *int       `json:"hourly_paid_holiday_mins,omitempty"`
	SpecialHoliday           *int       `json:"special_holiday,omitempty"`
	SpecialHolidaySettingID  *int       `json:"special_holiday_setting_id,omitempty"`
	HalfSpecialHolidayMins   *int       `json:"half_special_holiday_mins,omitempty"`
	HourlySpecialHolidayMins *int       `json:"hourly_special_holiday_mins,omitempty"`
	UseAttendanceDeduction   *bool      `json:"use_attendance_deduction,omitempty"`
	UseDefaultWorkPattern    *bool      `json:"use_default_work_pattern,omitempty"`
}

// PutWorkRecord は指定した従業員の勤怠情報を更新します。
// 注意点
// - 振替出勤・振替休日・代休出勤・代休の登録はAPIでは行うことができません。
func (c *Client) PutWorkRecord(employeeID int, date string, request *PutWorkRecordRequest) (WorkRecord, error) {
	u := "https://api.freee.co.jp/hr/api/v1/employees/" + url.PathEscape(strconv.Itoa(employeeID)) + "/work_records/" + url.PathEscape(date)
	resp, err := c.do(http.MethodPut, u, nil, request)
	if err != nil {
		return WorkRecord{}, err
	}

	workRecord := WorkRecord{}
	if err := resp.Parse(&workRecord); err != nil {
		return WorkRecord{}, err
	}

	return workRecord, nil
}
