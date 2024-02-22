package freee

import (
	"net/http"
	"net/url"
	"strconv"
)

// https://developer.freee.co.jp/reference/hr/reference#operations-tag-勤怠

type WorkRecord struct {
	BreakRecords []struct {
		ClockInAt  string `json:"clock_in_at"`
		ClockOutAt string `json:"clock_out_at"`
	} `json:"break_records"`
	ClockInAt                                 *string `json:"clock_in_at"`
	ClockOutAt                                *string `json:"clock_out_at"`
	Date                                      string  `json:"date"`
	DayPattern                                string  `json:"day_pattern"`
	SchedulePattern                           string  `json:"schedule_pattern"`
	EarlyLeavingMins                          int     `json:"early_leaving_mins"`
	HalfPaidHolidayMins                       int     `json:"half_paid_holiday_mins"`
	HalfSpecialHolidayMins                    int     `json:"half_special_holiday_mins"`
	HourlyPaidHolidayMins                     int     `json:"hourly_paid_holiday_mins"`
	HourlySpecialHolidayMins                  int     `json:"hourly_special_holiday_mins"`
	IsAbsence                                 bool    `json:"is_absence"`
	IsEditable                                bool    `json:"is_editable"`
	LatenessMins                              int     `json:"lateness_mins"`
	NormalWorkClockInAt                       *string `json:"normal_work_clock_in_at"`
	NormalWorkClockOutAt                      *string `json:"normal_work_clock_out_at"`
	NormalWorkMins                            int     `json:"normal_work_mins"`
	Note                                      string  `json:"note"`
	PaidHoliday                               float32 `json:"paid_holiday"`
	SpecialHoliday                            float32 `json:"special_holiday"`
	SpecialHolidaySettingID                   *int    `json:"special_holiday_setting_id"`
	UseAttendanceDeduction                    bool    `json:"use_attendance_deduction"`
	UseDefaultWorkPattern                     bool    `json:"use_default_work_pattern"`
	UseHalfCompensatoryHoliday                bool    `json:"use_half_compensatory_holiday"`
	TotalOvertimeWorkMins                     int     `json:"total_overtime_work_mins"`
	TotalHolidayWorkMins                      int     `json:"total_holiday_work_mins"`
	TotalLatenightWorkMins                    int     `json:"total_latenight_work_mins"`
	NotAutoCalcWorkTime                       bool    `json:"not_auto_calc_work_time"`
	TotalExcessStatutoryWorkMins              int     `json:"total_excess_statutory_work_mins"`
	TotalLatenightExcessStatutoryWorkMins     int     `json:"total_latenight_excess_statutory_work_mins"`
	TotalOvertimeExceptNormalWorkMins         int     `json:"total_overtime_except_normal_work_mins"`
	TotalLatenightOvertimeExceptNormalWorkMin int     `json:"total_latenight_overtime_except_normal_work_min"`
}

// DeleteWorkRecord は指定した従業員の勤怠情報を削除します。
func (c *Client) DeleteWorkRecord(companyID int, employeeID int, date Date) error {
	u := "https://api.freee.co.jp/hr/api/v1/employees/" + url.PathEscape(strconv.Itoa(employeeID)) + "/work_records/" + url.PathEscape(date.String())
	q := url.Values{
		"company_id": {strconv.Itoa(companyID)},
	}
	resp, err := c.do(http.MethodDelete, u, q, nil)
	if err != nil {
		return err
	}
	resp.Close()
	return nil
}

// GetWorkRecord は指定した従業員・日付の勤怠情報を返します。
func (c *Client) GetWorkRecord(employeeID int, companyID int, date Date) (WorkRecord, error) {
	u := "https://api.freee.co.jp/hr/api/v1/employees/" + url.PathEscape(strconv.Itoa(employeeID)) + "/work_records/" + url.PathEscape(date.String())
	pu, err := url.Parse(u)
	if err != nil {
		return WorkRecord{}, err
	}
	q := pu.Query()
	q.Set("company_id", strconv.Itoa(companyID))
	pu.RawQuery = q.Encode()
	resp, err := c.do(http.MethodGet, pu.String(), nil, nil)
	if err != nil {
		return WorkRecord{}, err
	}

	workRecord := WorkRecord{}
	if err := resp.Parse(&workRecord); err != nil {
		return WorkRecord{}, err
	}

	return workRecord, nil
}

type DayPattern string

const (
	NormalDay         DayPattern = "normal_day"
	PrescribedHoliday DayPattern = "prescribed_holiday"
	LegalHoliday      DayPattern = "legal_holiday"
)

type PutWorkRecordRequest struct {
	CompanyID                int                        `json:"company_id"`
	BreakRecords             []PutWorkRecordBreakRecord `json:"break_records,omitempty"`
	ClockInAt                *DateTime                  `json:"clock_in_at,omitempty"`
	ClockOutAt               *DateTime                  `json:"clock_out_at,omitempty"`
	DayPattern               *DayPattern                `json:"day_pattern,omitempty"`
	EarlyLeavingMins         *int                       `json:"early_leaving_mins,omitempty"`
	IsAbsence                *bool                      `json:"is_absence,omitempty"`
	LatenessMins             *int                       `json:"lateness_mins,omitempty"`
	NormalWorkClockInAt      *DateTime                  `json:"normal_work_clock_in_at,omitempty"`
	NormalWorkClockOutAt     *DateTime                  `json:"normal_work_clock_out_at,omitempty"`
	NormalWorkMins           *int                       `json:"normal_work_mins,omitempty"`
	Note                     *string                    `json:"note,omitempty"`
	PaidHoliday              *int                       `json:"paid_holiday,omitempty"`
	HalfPaidHolidayMins      *int                       `json:"half_paid_holiday_mins,omitempty"`
	HourlyPaidHolidayMins    *int                       `json:"hourly_paid_holiday_mins,omitempty"`
	SpecialHoliday           *int                       `json:"special_holiday,omitempty"`
	SpecialHolidaySettingID  *int                       `json:"special_holiday_setting_id,omitempty"`
	HalfSpecialHolidayMins   *int                       `json:"half_special_holiday_mins,omitempty"`
	HourlySpecialHolidayMins *int                       `json:"hourly_special_holiday_mins,omitempty"`
	UseAttendanceDeduction   *bool                      `json:"use_attendance_deduction,omitempty"`
	UseDefaultWorkPattern    *bool                      `json:"use_default_work_pattern,omitempty"`
}

type PutWorkRecordBreakRecord struct {
	ClockInAt  DateTime `json:"clock_in_at,omitempty"`
	ClockOutAt DateTime `json:"clock_out_at,omitempty"`
}

// PutWorkRecord は指定した従業員の勤怠情報を更新します。
// 注意点
// - 振替出勤・振替休日・代休出勤・代休の登録はAPIでは行うことができません。
func (c *Client) PutWorkRecord(employeeID int, date Date, request *PutWorkRecordRequest) (WorkRecord, error) {
	u := "https://api.freee.co.jp/hr/api/v1/employees/" + url.PathEscape(strconv.Itoa(employeeID)) + "/work_records/" + url.PathEscape(date.String())
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

type DaysAndHours struct {
	Days  float32 `json:"days"`
	Hours int     `json:"hours"`
}

type WorkRecordSummariesWage struct {
	Name                string `json:"name"`
	TotalNormalTimeMins int    `json:"total_normal_time_mins"`
}

type WorkRecordSummaries struct {
	Year                                        int                       `json:"year"`
	Month                                       int                       `json:"month"`
	StartDate                                   string                    `json:"start_date"`
	EndDate                                     string                    `json:"end_date"`
	WorkDays                                    float32                   `json:"work_days"`
	TotalWorkMins                               int                       `json:"total_work_mins"`
	TotalNormalWorkMins                         int                       `json:"total_normal_work_mins"`
	TotalExcessStatutoryWorkMins                int                       `json:"total_excess_statutory_work_mins"`
	TotalOvertimeExceptNormalWorkMins           int                       `json:"total_overtime_except_normal_work_mins"`
	TotalOvertimeWithinNormalWorkMins           int                       `json:"total_overtime_within_normal_work_mins"`
	TotalHolidayWorkMins                        int                       `json:"total_holiday_work_mins"`
	TotalLatenightWorkMins                      int                       `json:"total_latenight_work_mins"`
	NumAbsences                                 float32                   `json:"num_absences"`
	NumPaidHolidays                             float32                   `json:"num_paid_holidays"`
	NumPaidHolidaysAndHours                     DaysAndHours              `json:"num_paid_holidays_and_hours"`
	NumPaidHolidaysLeft                         float32                   `json:"num_paid_holidays_left"`
	NumPaidHolidaysAndHoursLeft                 DaysAndHours              `json:"num_paid_holidays_and_hours_left"`
	NumSubstituteHolidaysUsed                   float32                   `json:"num_substitute_holidays_used"`
	NumCompensatoryHolidaysUsed                 float32                   `json:"num_compensatory_holidays_used"`
	NumSpecialHolidaysUsed                      float32                   `json:"num_special_holidays_used"`
	NumSpecialHolidaysAndHoursUsed              DaysAndHours              `json:"num_special_holidays_and_hours_used"`
	TotalLatenessAndEarlyLeavingMins            int                       `json:"total_lateness_and_early_leaving_mins"`
	MultiHourlyWages                            []WorkRecordSummariesWage `json:"multi_hourly_wages"`
	WorkRecords                                 []WorkRecord              `json:"work_records"`
	TotalShortageWorkMins                       *int                      `json:"total_shortage_work_mins"`
	TotalDeemedPaidExcessStatutoryWorkMins      *int                      `json:"total_deemed_paid_excess_statutory_work_mins"`
	TotalDeemedPaidOvertimeExceptNormalWorkMins *int                      `json:"total_deemed_paid_overtime_except_normal_work_mins"`
}

type GetWorkRecordOpts struct {
	WorkRecords bool // サマリ情報に日次の勤怠情報を含める(true/false)(デフォルト: false)
}

// GetWorkRecordSummariesは、指定した従業員、月の勤怠情報のサマリを返します。
// 注意点
// - work_recordsオプションにtrueを指定することで、明細となる日次の勤怠情報もあわせて返却します。
func (c *Client) GetWorkRecordSummaries(employeeID int, companyID int, year int, month int, opts *GetWorkRecordOpts) (WorkRecordSummaries, error) {
	u := "https://api.freee.co.jp/hr/api/v1/employees/" + url.PathEscape(strconv.Itoa(employeeID)) + "/work_record_summaries/" + url.PathEscape(strconv.Itoa(year)) + "/" + url.PathEscape(strconv.Itoa(month))
	q := url.Values{
		"company_id": {strconv.Itoa(companyID)},
	}

	if opts != nil && opts.WorkRecords {
		q.Set("work_records", "true")
	}

	resp, err := c.do(http.MethodGet, u, q, nil)
	if err != nil {
		return WorkRecordSummaries{}, err
	}

	summaries := WorkRecordSummaries{}
	if err := resp.Parse(&summaries); err != nil {
		return WorkRecordSummaries{}, nil
	}

	return summaries, nil
}

// 値が設定されなかった場合は自動的に0が設定されます
type PutWorkRecordSummariesRequest struct {
	CompanyID                                      int     `json:"company_id"`                                                   // 事業所ID（必須）
	WorkDays                                       float32 `json:"work_days,omitempty"`                                          // 総勤務日数
	WorkDaysOnWeekdays                             float32 `json:"work_days_on_weekdays,omitempty"`                              // 所定労働日の勤務日数
	WorkDaysOnPrescribedHolidays                   float32 `json:"work_days_on_prescribed_holidays,omitempty"`                   // 所定休日の勤務日数
	WorkDaysOnLegalHolidays                        float32 `json:"work_days_on_legal_holidays,omitempty"`                        // 法定休日の勤務日数
	TotalWorkMins                                  int     `json:"total_work_mins,omitempty"`                                    // 労働時間（分）
	TotalNormalWorkMins                            int     `json:"total_normal_work_mins,omitempty"`                             // 所定労働時間（分）
	TotalExcessStatutoryWorkMins                   int     `json:"total_excess_statutory_work_mins,omitempty"`                   // 給与計算に用いられる法定内残業時間（分）
	TotalHolidayWorkMins                           int     `json:"total_holiday_work_mins,omitempty"`                            // 法定休日労働時間（分）
	TotalLatenightWorkMins                         int     `json:"total_latenight_work_mins,omitempty"`                          // 深夜労働時間（分）
	TotalActualExcessStatutoryWorkMins             int     `json:"total_actual_excess_statutory_work_mins,omitempty"`            // 実労働時間ベースの法定内残業時間（分）
	TotalOvertimeWorkMins                          int     `json:"total_overtime_work_mins,omitempty"`                           // 時間外労働時間（分）
	NumAbsences                                    float32 `json:"num_absences,omitempty"`                                       // 欠勤日数
	NumAbsencesForDeduction                        float32 `json:"num_absences_for_deduction,omitempty"`                         // 控除対象の欠勤日数
	TotalLatenessMins                              int     `json:"total_lateness_mins,omitempty"`                                // 遅刻時間（分）
	TotalLatenessMinsForDeduction                  int     `json:"total_lateness_mins_for_deduction,omitempty"`                  // 控除対象の遅刻時間（分）
	TotalEarlyLeavingMins                          int     `json:"total_early_leaving_mins,omitempty"`                           // 早退時間（分）
	TotalEarlyLeavingMinsForDeduction              int     `json:"total_early_leaving_mins_for_deduction,omitempty"`             // 控除対象の早退時間（分）
	NumPaidHolidays                                float32 `json:"num_paid_holidays,omitempty"`                                  // 有給取得日数
	TotalShortageWorkMins                          int     `json:"total_shortage_work_mins,omitempty"`                           // 不足時間（分）（フレックスタイム制でのみ使用）
	TotalDeemedPaidOvertimeExcessStatutoryWorkMins int     `json:"total_deemed_paid_excess_statutory_work_mins,omitempty"`       // 支給対象の法定内残業時間（分）（裁量労働制でのみ使用）
	TOtalDeemedPaidOvertimeExceptNormalWorkMins    int     `json:"total_deemed_paid_overtime_except_normal_work_mins,omitempty"` // 支給対象の時間外労働時間（分）（裁量労働制でのみ使用）
}

// PutWorkRecordSummariesは、指定した従業員、月の勤怠情報のサマリを更新します。
// 注意点
// - 管理者権限を持ったユーザーのみ実行可能です。
// - 日毎の勤怠の更新はこのAPIではできません。日毎の勤怠の操作には勤怠APIを使用して下さい。
// - 勤怠データが存在しない場合は新規作成、既に存在する場合は上書き更新されます。
// - 値が設定された項目のみ更新されます。値が設定されなかった場合は自動的に0が設定されます。
func (c *Client) PutWorkRecordSummaries(employeeID int, year int, month int, request *PutWorkRecordSummariesRequest) (WorkRecordSummaries, error) {
	u := "https://api.freee.co.jp/hr/api/v1/employees/" + url.PathEscape(strconv.Itoa(employeeID)) + "/work_record_summaries/" + url.PathEscape(strconv.Itoa(year)) + "/" + url.PathEscape(strconv.Itoa(month))
	resp, err := c.do(http.MethodPut, u, nil, request)
	if err != nil {
		return WorkRecordSummaries{}, err
	}
	summaries := WorkRecordSummaries{}
	if err := resp.Parse(&summaries); err != nil {
		return WorkRecordSummaries{}, err
	}
	return summaries, nil
}
