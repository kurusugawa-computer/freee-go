package freee

import (
	"net/http"
	"net/url"
	"strconv"
)

// https://developer.freee.co.jp/reference/hr/reference#operations-tag-従業員

// DeleteEmployee は指定したIDの従業員を削除します。
// 注意点
// - 管理者権限を持ったユーザーのみ実行可能です。
func (c *Client) DeleteEmployee(companyID int, employeeID int) error {
	u := "https://api.freee.co.jp/hr/api/v1/employees/" + url.PathEscape(strconv.Itoa(employeeID))
	q := url.Values{
		"company_id ": {strconv.Itoa(companyID)},
	}
	resp, err := c.do(http.MethodDelete, u, q, nil)
	if err != nil {
		return err
	}
	resp.Close()
	return nil
}

type CompaniesEmployee struct {
	ID                 int     `json:"id"`
	Num                *string `json:"num"`
	DisplayName        string  `json:"display_name"`
	EntryDate          string  `json:"entry_date"`
	RetireDate         *string `json:"retire_date"`
	UserID             int     `json:"user_id"`
	Email              *string `json:"email"`
	PayrollCalculation bool    `json:"payroll_calculation"`
	ClosingDay         *int    `json:"closing_day"`
	PayDay             *int    `json:"pay_day"`
	MonthOfPayDay      *string `json:"month_of_pay_day"`
}

type ListAllEmployeesOpts struct {
	Limit                    int  // 取得レコードの件数 (デフォルト: 50, 最小: 1, 最大: 100)
	Offset                   int  // 取得レコードのオフセット (デフォルト: 0)
	WithNoPayrollCalculation bool // trueを指定すると給与計算対象外の従業員情報をレスポンスに含めます。
}

// ListEmployees は指定した事業所に所属する従業員をリストで返します。
// 注意点
// - 管理者権限を持ったユーザーのみ実行可能です。
// - 退職ユーザーも含めて取得可能です。
func (c *Client) ListCompaniesEmployees(companyID int, opts *ListAllEmployeesOpts) ([]CompaniesEmployee, error) {
	u := "https://api.freee.co.jp/hr/api/v1/companies/" + url.PathEscape(strconv.Itoa(companyID)) + "/employees"
	q := url.Values{}
	if opts != nil {
		if opts.Limit != 50 {
			q.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Offset > 0 {
			q.Set("offset", strconv.Itoa(opts.Limit))
		}
		if opts.WithNoPayrollCalculation {
			q.Set("with_no_payroll_calculation", strconv.FormatBool(opts.WithNoPayrollCalculation))
		}
	}
	resp, err := c.do(http.MethodGet, u, q, nil)
	if err != nil {
		return nil, err
	}

	employees := []CompaniesEmployee{}
	if err := resp.Parse(&employees); err != nil {
		return nil, err
	}

	return employees, nil
}

type Employee struct {
	ID                                 int     `json:"id"`
	CompanyID                          int     `json:"company_id"`
	Num                                *string `json:"num"`
	DisplayName                        string  `json:"display_name"`
	BasePensionNum                     *string `json:"base_pension_num"`
	EmploymentInsuranceReferenceNumber string  `json:"employment_insurance_reference_number"`
	BirthDate                          string  `json:"birth_date"`
	EntryDate                          string  `json:"entry_date"`
	RetireDate                         *string `json:"retire_date"`
	UserID                             *int    `json:"user_id"`
	ProfileRule                        struct {
		ID                        int     `json:"id"`
		CompanyID                 int     `json:"company_id"`
		EmployeeID                int     `json:"employee_id"`
		LastName                  string  `json:"last_name"`
		FirstName                 string  `json:"first_name"`
		LastNameKana              string  `json:"last_name_kana"`
		FirstNameKana             string  `json:"first_name_kana"`
		Zipcode1                  *string `json:"zipcode1"`
		Zipcode2                  *string `json:"zipcode2"`
		PrefectureCode            *int    `json:"prefecture_code"`
		Address                   *string `json:"address"`
		AddressKana               *string `json:"address_kana"`
		Phone1                    *string `json:"phone1"`
		Phone2                    *string `json:"phone2"`
		Phone3                    *string `json:"phone3"`
		ResidentialZipcode1       *string `json:"residential_zipcode1"`
		ResidentialZipcode2       *string `json:"residential_zipcode2"`
		ResidentialPrefectureCode *int    `json:"residential_prefecture_code"`
		ResidentialAddress        *string `json:"residential_address"`
		ResidentialAddressKana    *string `json:"residential_address_kana"`
		EmploymentType            *string `json:"employment_type"`
		Title                     *string `json:"title"`
		Gender                    string  `json:"gender"`
		Married                   bool    `json:"married"`
		IsWorkingStudent          bool    `json:"is_working_student"`
		WidowType                 string  `json:"widow_type"`
		DisabilityType            string  `json:"disability_type"`
		Email                     *string `json:"email"`
		HouseholderName           string  `json:"householder_name"`
		Householder               *string `json:"householder"`
	} `json:"profile_rule"`
	HealthInsuranceRule struct {
		ID                                          int      `json:"id"`
		CompanyID                                   int      `json:"company_id"`
		EmployeeID                                  int      `json:"employee_id"`
		Entried                                     bool     `json:"entried"`
		HealthInsuranceSalaryCalcType               string   `json:"health_insurance_salary_calc_type"`
		HealthInsuranceBonusCalcType                string   `json:"health_insurance_bonus_calc_type"`
		ManualHealthInsuranceAmountOfEmployeeSalary *int     `json:"manual_health_insurance_amount_of_employee_salary"`
		ManualHealthInsuranceAmountOfEmployeeBonus  *int     `json:"manual_health_insurance_amount_of_employee_bonus"`
		ManualHealthInsuranceAmountOfCompanySalary  *float64 `json:"manual_health_insurance_amount_of_company_salary"`
		ManualHealthInsuranceAmountOfCompanyBonus   *float64 `json:"manual_health_insurance_amount_of_company_bonus"`
		CareInsuranceSalaryCalcType                 string   `json:"care_insurance_salary_calc_type"`
		CareInsuranceBonusCalcType                  string   `json:"care_insurance_bonus_calc_type"`
		ManualCareInsuranceAmountOfEmployeeSalary   *int     `json:"manual_care_insurance_amount_of_employee_salary"`
		ManualCareInsuranceAmountOfEmployeeBonus    *int     `json:"manual_care_insurance_amount_of_employee_bonus"`
		ManualCareInsuranceAmountOfCompanySalary    *float64 `json:"manual_care_insurance_amount_of_company_salary"`
		ManualCareInsuranceAmountOfCompanyBonus     *float64 `json:"manual_care_insurance_amount_of_company_bonus"`
		ReferenceNum                                *string  `json:"reference_num"`
		StandardMonthlyRemuneration                 int      `json:"standard_monthly_remuneration"`
	} `json:"health_insurance_rule"`
	WelfarePensionInsuranceRule struct {
		ID                                                  int      `json:"id"`
		ChildAllowanceContributionBonusCalcType             string   `json:"child_allowance_contribution_bonus_calc_type"`
		ChildAllowanceContributionSalaryCalcType            string   `json:"child_allowance_contribution_salary_calc_type"`
		CompanyID                                           int      `json:"company_id"`
		EmployeeID                                          int      `json:"employee_id"`
		Entried                                             bool     `json:"entried"`
		ManualChildAllowanceContributionAmountBonus         *float64 `json:"manual_child_allowance_contribution_amount_bonus"`
		ManualChildAllowanceContributionAmountSalary        *float64 `json:"manual_child_allowance_contribution_amount_salary"`
		ManualWelfarePensionInsuranceAmountOfCompanyBonus   *float64 `json:"manual_welfare_pension_insurance_amount_of_company_bonus"`
		ManualWelfarePensionInsuranceAmountOfCompanySalary  *float64 `json:"manual_welfare_pension_insurance_amount_of_company_salary"`
		ManualWelfarePensionInsuranceAmountOfEmployeeBonus  *int     `json:"manual_welfare_pension_insurance_amount_of_employee_bonus"`
		ManualWelfarePensionInsuranceAmountOfEmployeeSalary *int     `json:"manual_welfare_pension_insurance_amount_of_employee_salary"`
		ReferenceNum                                        *string  `json:"reference_num"`
		StandardMonthlyRemuneration                         int      `json:"standard_monthly_remuneration"`
		WelfarePensionInsuranceBonusCalcType                string   `json:"welfare_pension_insurance_bonus_calc_type"`
		WelfarePensionInsuranceSalaryCalcType               string   `json:"welfare_pension_insurance_salary_calc_type"`
	} `json:"welfare_pension_insurance_rule"`
	DependentRules []struct {
		ID                                                  int     `json:"id"`
		CompanyID                                           int     `json:"company_id"`
		EmployeeID                                          int     `json:"employee_id"`
		LastName                                            string  `json:"last_name"`
		FirstName                                           string  `json:"first_name"`
		LastNameKana                                        *string `json:"last_name_kana"`
		FirstNameKana                                       *string `json:"first_name_kana"`
		Gender                                              string  `json:"gender"`
		Relationship                                        string  `json:"relationship"`
		BirthDate                                           string  `json:"birth_date"`
		ResidenceType                                       string  `json:"residence_type"`
		Zipcode1                                            *string `json:"zipcode1"`
		Zipcode2                                            *string `json:"zipcode2"`
		PrefectureCode                                      *int    `json:"prefecture_code"`
		Address                                             *string `json:"address"`
		AddressKana                                         *string `json:"address_kana"`
		BasePensionNum                                      *string `json:"base_pension_num"`
		Income                                              int     `json:"income"`
		AnnualRevenue                                       int     `json:"annual_revenue"`
		DisabilityType                                      string  `json:"disability_type"`
		Occupation                                          *string `json:"occupation"`
		AnnualRemittanceAmount                              int     `json:"annual_remittance_amount"`
		EmploymentInsuranceReceiveStatus                    *string `json:"employment_insurance_receive_status"`
		EmploymentInsuranceReceivesFrom                     *string `json:"employment_insurance_receives_from"`
		PhoneType                                           *string `json:"phone_type"`
		Phone1                                              *string `json:"phone1"`
		Phone2                                              *string `json:"phone2"`
		Phone3                                              *string `json:"phone3"`
		SocialInsuranceAndTaxDependent                      string  `json:"social_insurance_and_tax_dependent"`
		SocialInsuranceDependentAcquisitionDate             *string `json:"social_insurance_dependent_acquisition_date"`
		SocialInsuranceDependentAcquisitionReason           string  `json:"social_insurance_dependent_acquisition_reason"`
		SocialInsuranceOtherDependentAcquisitionReason      *string `json:"social_insurance_other_dependent_acquisition_reason"`
		SocialInsuranceDependentDisqualificationDate        *string `json:"social_insurance_dependent_disqualification_date"`
		SocialInsuranceDependentDisqualificationReason      string  `json:"social_insurance_dependent_disqualification_reason"`
		SocialInsuranceOtherDependentDisqualificationReason *string `json:"social_insurance_other_dependent_disqualification_reason"`
		TaxDependentAcquisitionDate                         *string `json:"tax_dependent_acquisition_date"`
		TaxDependentAcquisitionReason                       string  `json:"tax_dependent_acquisition_reason"`
		TaxOtherDependentAcquisitionReason                  *string `json:"tax_other_dependent_acquisition_reason"`
		TaxDependentDisqualificationDate                    *string `json:"tax_dependent_disqualification_date"`
		TaxDependentDisqualificationReason                  string  `json:"tax_dependent_disqualification_reason"`
		TaxOtherDependentDisqualificationReason             *string `json:"tax_other_dependent_disqualification_reason"`
		NonResidentDependentsReason                         string  `json:"non_resident_dependents_reason"`
	} `json:"dependent_rules"`
	BankAccountRule struct {
		ID             int     `json:"id"`
		CompanyID      int     `json:"company_id"`
		EmployeeID     int     `json:"employee_id"`
		BankName       *string `json:"bank_name"`
		BankNameKana   *string `json:"bank_name_kana"`
		BankCode       *string `json:"bank_code"`
		BranchName     *string `json:"branch_name"`
		BranchNameKana *string `json:"branch_name_kana"`
		BranchCode     *string `json:"branch_code"`
		AccountNumber  *string `json:"account_number"`
		AccountName    *string `json:"account_name"`
		AccountType    *string `json:"account_type"`
	} `json:"bank_account_rule"`
	BasicPayRule struct {
		ID          int    `json:"id"`
		CompanyID   int    `json:"company_id"`
		EmployeeID  int    `json:"employee_id"`
		PayCalcType string `json:"pay_calc_type"`
		PayAmount   int    `json:"pay_amount"`
	} `json:"basic_pay_rule"`
	PayrollCalculation           bool    `json:"payroll_calculation"`
	CompanyReferenceDateRuleName *string `json:"company_reference_date_rule_name"`
}

type ListEmployeesOpts struct {
	Limit                    int  // 取得レコードの件数 (デフォルト: 50, 最小: 1, 最大: 100)
	Offset                   int  // 取得レコードのオフセット (デフォルト: 0)
	WithNoPayrollCalculation bool // trueを指定すると給与計算対象外の従業員情報をレスポンスに含めます。
}

// ListEmployees は指定した対象年月に事業所に所属する従業員をリストで返します。
// 注意点
// - 管理者権限を持ったユーザーのみ実行可能です。
// - 指定した年月に退職済みユーザーは取得できません。
// - 保険料計算方法が自動計算の場合、対応する保険料の直接指定金額は無視されnullが返されます。(例: 給与計算時の健康保険料の計算方法が自動計算の場合、給与計算時の健康保険料の直接指定金額はnullが返されます)
// - 事業所が定額制の健康保険組合に加入している場合、保険料の直接指定金額は無視されnullが返されます。
func (c *Client) ListEmployees(companyID int, year int, month int, opts *ListEmployeesOpts) ([]Employee, error) {
	u := "https://api.freee.co.jp/hr/api/v1/employees"
	q := url.Values{
		"company_id": {strconv.Itoa(companyID)},
		"year":       {strconv.Itoa(year)},
		"month":      {strconv.Itoa(month)},
	}
	if opts != nil {
		if opts.Limit != 50 {
			q.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Offset > 0 {
			q.Set("offset", strconv.Itoa(opts.Limit))
		}
		if opts.WithNoPayrollCalculation {
			q.Set("with_no_payroll_calculation", strconv.FormatBool(opts.WithNoPayrollCalculation))
		}
	}
	resp, err := c.do(http.MethodGet, u, q, nil)
	if err != nil {
		return nil, err
	}

	result := struct {
		Employees  []Employee `json:"employees"`
		TotalCount int        `json:"total_count"`
	}{}
	if err := resp.Parse(&result); err != nil {
		return nil, err
	}

	return result.Employees, nil
}

// GetEmployee は指定したIDの従業員を返します。
// 注意点
// - 管理者権限を持ったユーザーのみ実行可能です。
// - 指定した年月に退職済みユーザーは取得できません。
// - 保険料計算方法が自動計算の場合、対応する保険料の直接指定金額は無視されnullが返されます。(例: 給与計算時の健康保険料の計算方法が自動計算の場合、給与計算時の健康保険料の直接指定金額はnullが返されます)
// - 事業所が定額制の健康保険組合に加入している場合、保険料の直接指定金額は無視されnullが返されます。
func (c *Client) GetEmployee(companyID int, employeeID int, year int, month int) (Employee, error) {
	u := "https://api.freee.co.jp/hr/api/v1/employees/" + url.PathEscape(strconv.Itoa(employeeID))
	q := url.Values{
		"company_id": {strconv.Itoa(companyID)},
		"year":       {strconv.Itoa(year)},
		"month":      {strconv.Itoa(month)},
	}
	resp, err := c.do(http.MethodGet, u, q, nil)
	if err != nil {
		return Employee{}, err
	}

	result := struct {
		Employee Employee `json:"employee"`
	}{}
	if err := resp.Parse(&result); err != nil {
		return Employee{}, err
	}

	return result.Employee, nil
}

type CreateEmployeeRequest struct {
	CompanyID int `json:"company_id"`
	Employee  struct {
		Num                          *string `json:"num,omitempty"`
		WorkingHoursSystemName       *string `json:"working_hours_system_name,omitempty"`
		CompanyReferenceDateRuleName *string `json:"company_reference_date_rule_name,omitempty"`
		LastName                     string  `json:"last_name"`
		FirstName                    string  `json:"first_name"`
		LastNameKana                 string  `json:"last_name_kana"`
		FirstNameKana                string  `json:"first_name_kana"`
		BirthDate                    string  `json:"birth_date"`
		EntryDate                    string  `json:"entry_date"`
		PayCalcType                  *string `json:"pay_calc_type,omitempty"`
		PayAmount                    *int    `json:"pay_amount,omitempty"`
		Gender                       *string `json:"gender,omitempty"`
		Married                      *bool   `json:"married,omitempty"`
		NoPayrollCalculation         *bool   `json:"no_payroll_calculation,omitempty"`
	} `json:"employee"`
}

// CreateEmployee は従業員を新規作成します。
// 注意点
// - 管理者権限を持ったユーザーのみ実行可能です。
func (c *Client) CreateEmployee(req *CreateEmployeeRequest) (Employee, error) {
	u := "https://api.freee.co.jp/hr/api/v1/employees"
	resp, err := c.do(http.MethodPost, u, nil, req)
	if err != nil {
		return Employee{}, err
	}

	result := struct {
		Employee Employee `json:"employee"`
	}{}
	if err := resp.Parse(&result); err != nil {
		return Employee{}, err
	}

	return result.Employee, nil
}

type UpdateEmployeeRequest struct {
	CompanyID int  `json:"company_id"`
	Year      *int `json:"year,omitempty"`
	Month     *int `json:"month,omitempty"`
	Employee  struct {
		Num                                *string `json:"num,omitempty"`
		DisplayName                        *string `json:"display_name,omitempty"`
		BasePensionNum                     *string `json:"base_pension_num,omitempty"`
		EmploymentInsuranceReferenceNumber *string `json:"employment_insurance_reference_number,omitempty"`
		BirthDate                          string  `json:"birth_date"`
		EntryDate                          string  `json:"entry_date"`
		RetireDate                         *string `json:"retire_date,omitempty"`
		CompanyReferenceDateRuleName       *string `json:"company_reference_date_rule_name,omitempty"`
	} `json:"employee"`
}

// CreateEmployee は従業員を新規作成します。
// 注意点
// - 管理者権限を持ったユーザーのみ実行可能です。
func (c *Client) UpdateEmployee(employeeID int, request *UpdateEmployeeRequest) (Employee, error) {
	u := "https://api.freee.co.jp/hr/api/v1/employees/" + url.PathEscape(strconv.Itoa(employeeID))
	resp, err := c.do(http.MethodPut, u, nil, request)
	if err != nil {
		return Employee{}, err
	}

	result := struct {
		Employee Employee `json:"employee"`
	}{}
	if err := resp.Parse(&result); err != nil {
		return Employee{}, err
	}

	return result.Employee, nil
}
