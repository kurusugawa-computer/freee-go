package freee

import (
	"net/http"
)

// https://developer.freee.co.jp/reference/hr/reference#operations-tag-ログインユーザー

type LoginUser struct {
	ID        int `json:"id"`
	Companies []struct {
		ID          int     `json:"id"`
		Name        string  `json:"name"`
		Role        string  `json:"role"` // company_admin, self_only, clerk
		ExternalCID string  `json:"external_cid"`
		EmployeeID  *int    `json:"employee_id"`
		DisplayName *string `json:"display_name"`
	} `json:"companies"`
}

// GetLoginUser はこのリクエストの認可セッションにおけるログインユーザーの情報を返します。
// freee人事労務では一人のログインユーザーを複数の事業所に関連付けられるため、このユーザーと関連のあるすべての事業所の情報をリストで返します。
// 注意点
// - 他のAPIのパラメータとしてcompany_idが求められる場合は、このAPIで取得したcompany_idを使用します。
// - 給与計算対象外の従業員のemployee_idとdisplay_nameは取得できません。
func (c *Client) GetLoginUser() (LoginUser, error) {
	u := "https://api.freee.co.jp/hr/api/v1/users/me"
	resp, err := c.do(http.MethodGet, u, nil, nil)
	if err != nil {
		return LoginUser{}, err
	}

	loginUser := LoginUser{}
	if err := resp.Parse(&loginUser); err != nil {
		return LoginUser{}, err
	}

	return loginUser, nil
}
