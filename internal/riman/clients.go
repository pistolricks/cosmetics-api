package riman

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/pistolricks/cosmetics-api/internal/data"
	"resty.dev/v3"
)

type ClientModel struct {
	DB *sql.DB
}

type ReissueTokenResponse = map[string]any

func (m ClientModel) ReissueToken(token string) (*ReissueTokenResponse, error) {

	logoutUrl := fmt.Sprintf("https://security-api.riman.com/api/v2/token/reissue")

	fmt.Println(logoutUrl)

	client := resty.New()
	defer client.Close()

	res, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(token).
		SetResult(&ReissueTokenResponse{}).
		SetError(&Errors{}).
		Post(logoutUrl)

	if err != nil {
		return nil, err
	}

	fmt.Println(res.String())
	fmt.Println("string | cart")
	fmt.Println(res.Result().(*ReissueTokenResponse))

	return res.Result().(*ReissueTokenResponse), err

}

type LogoutResponse = map[string]any

func (m ClientModel) Logout(token string) (*LogoutResponse, error) {

	logoutUrl := fmt.Sprintf("https://security-api.riman.com/api/v2/token/logout")

	fmt.Println(logoutUrl)

	client := resty.New()
	defer client.Close()

	res, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(token).
		SetResult(&LogoutResponse{}).
		SetError(&Errors{}).
		Post(logoutUrl)

	if err != nil {
		return nil, err
	}

	fmt.Println(res.String())
	fmt.Println("string | cart")
	fmt.Println(res.Result().(*LogoutResponse))

	return res.Result().(*LogoutResponse), err
}

func (m ClientModel) GetByClientUserName(username string) (*Client, error) {
	query := `
        SELECT id, created_at, first_name, middle_name, last_name, suffix, email, mobile, username, riman_user_id, status, organization_type, signup_date, anniversary_date, account_type, sponsor_username, member_id, rank, enrollment_date, personal_orders_volume, personal_clients_volume, total_personal_volume, current_month_sp, current_month_bp, last_order_date, last_order_id, last_order_sp, last_order_bp, lifetime_spend, most_recent_12_month_spend, data, token
        FROM clients
        WHERE username = $1`

	var client Client

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, username).Scan(
		&client.ID,
		&client.CreatedAt,
		&client.FirstName,
		&client.MiddleName,
		&client.LastName,
		&client.Suffix,
		&client.Email,
		&client.Mobile,
		&client.Username,
		&client.RimanUserId,
		&client.Status,
		&client.OrganizationType,
		&client.SignupDate,
		&client.AnniversaryDate,
		&client.AccountType,
		&client.SponsorUsername,
		&client.MemberId,
		&client.Rank,
		&client.EnrollmentDate,
		&client.PersonalOrdersVolume,
		&client.PersonalClientsVolume,
		&client.TotalPersonalVolume,
		&client.CurrentMonthSp,
		&client.CurrentMonthBp,
		&client.LastOrderDate,
		&client.LastOrderId,
		&client.LastOrderSp,
		&client.LastOrderBp,
		&client.LifetimeSpend,
		&client.MostRecent12MonthSpend,
		&client.Data,
		&client.Token,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &client, nil
}

func (m ClientModel) GetByClientEmail(email string) (*Client, error) {
	query := `
        SELECT id, created_at, first_name, middle_name, last_name, suffix, email, mobile, username, riman_user_id, status, organization_type, signup_date, anniversary_date, account_type, sponsor_username, member_id, rank, enrollment_date, personal_orders_volume, personal_clients_volume, total_personal_volume, current_month_sp, current_month_bp, last_order_date, last_order_id, last_order_sp, last_order_bp, lifetime_spend, most_recent_12_month_spend, data, token
        FROM clients
        WHERE email = $1`

	var client Client

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&client.ID,
		&client.CreatedAt,
		&client.FirstName,
		&client.MiddleName,
		&client.LastName,
		&client.Suffix,
		&client.Email,
		&client.Mobile,
		&client.Username,
		&client.RimanUserId,
		&client.Status,
		&client.OrganizationType,
		&client.SignupDate,
		&client.AnniversaryDate,
		&client.AccountType,
		&client.SponsorUsername,
		&client.MemberId,
		&client.Rank,
		&client.EnrollmentDate,
		&client.PersonalOrdersVolume,
		&client.PersonalClientsVolume,
		&client.TotalPersonalVolume,
		&client.CurrentMonthSp,
		&client.CurrentMonthBp,
		&client.LastOrderDate,
		&client.LastOrderId,
		&client.LastOrderSp,
		&client.LastOrderBp,
		&client.LifetimeSpend,
		&client.MostRecent12MonthSpend,
		&client.Data,
		&client.Token,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, data.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &client, nil
}

func (m ClientModel) GetAll() ([]*Client, data.Metadata, error) {

	query := fmt.Sprintf(`
	SELECT count(*) OVER(),id, created_at, first_name, middle_name, last_name, suffix, email, mobile, username, riman_user_id, status, organization_type, signup_date, anniversary_date, account_type, sponsor_username, member_id, rank, enrollment_date, personal_orders_volume, personal_clients_volume, total_personal_volume, current_month_sp, current_month_bp, last_order_date, last_order_id, last_order_sp, last_order_bp, lifetime_spend, most_recent_12_month_spend, data, token
	FROM clients
	`)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, data.Metadata{}, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	totalRecords := 0
	clients := []*Client{}

	for rows.Next() {
		var client Client

		err := rows.Scan(
			&totalRecords,
			&client.ID,
			&client.CreatedAt,
			&client.FirstName,
			&client.MiddleName,
			&client.LastName,
			&client.Suffix,
			&client.Email,
			&client.Mobile,
			&client.Username,
			&client.RimanUserId,
			&client.Status,
			&client.OrganizationType,
			&client.SignupDate,
			&client.AnniversaryDate,
			&client.AccountType,
			&client.SponsorUsername,
			&client.MemberId,
			&client.Rank,
			&client.EnrollmentDate,
			&client.PersonalOrdersVolume,
			&client.PersonalClientsVolume,
			&client.TotalPersonalVolume,
			&client.CurrentMonthSp,
			&client.CurrentMonthBp,
			&client.LastOrderDate,
			&client.LastOrderId,
			&client.LastOrderSp,
			&client.LastOrderBp,
			&client.LifetimeSpend,
			&client.MostRecent12MonthSpend,
			&client.Data,
			&client.Token,
		)
		if err != nil {
			return nil, data.Metadata{}, err

		}

		clients = append(clients, &client)
	}

	if err = rows.Err(); err != nil {
		return nil, data.Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, 1, 20)

	return clients, metadata, nil
}

func (m ClientModel) GetForRimanToken(tokenScope, tokenPlaintext string) (*Client, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `
        SELECT clients.id, clients.created_at, clients.first_name, clients.middle_name, clients.last_name, clients.suffix, clients.email, clients.mobile, clients.username, clients.riman_user_id, clients.status, clients.organization_type, clients.signup_date, clients.anniversary_date, clients.account_type, clients.sponsor_username, clients.member_id, clients.rank, clients.enrollment_date, clients.personal_orders_volume, clients.personal_clients_volume, clients.total_personal_volume, clients.current_month_sp, clients.current_month_bp, clients.last_order_date, clients.last_order_id, clients.last_order_sp, clients.last_order_bp, clients.lifetime_spend, clients.most_recent_12_month_spend, clients.data, clients.token
        FROM clients
        INNER JOIN sessions
        ON clients.id = sessions.client_id
        WHERE sessions.hash = $1
        AND sessions.scope = $2 
        AND sessions.expiry > $3
		`

	args := []any{tokenHash[:], tokenScope, time.Now()}

	var client Client

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&client.ID,
		&client.CreatedAt,
		&client.FirstName,
		&client.MiddleName,
		&client.LastName,
		&client.Suffix,
		&client.Email,
		&client.Mobile,
		&client.Username,
		&client.RimanUserId,
		&client.Status,
		&client.OrganizationType,
		&client.SignupDate,
		&client.AnniversaryDate,
		&client.AccountType,
		&client.SponsorUsername,
		&client.MemberId,
		&client.Rank,
		&client.EnrollmentDate,
		&client.PersonalOrdersVolume,
		&client.PersonalClientsVolume,
		&client.TotalPersonalVolume,
		&client.CurrentMonthSp,
		&client.CurrentMonthBp,
		&client.LastOrderDate,
		&client.LastOrderId,
		&client.LastOrderSp,
		&client.LastOrderBp,
		&client.LifetimeSpend,
		&client.MostRecent12MonthSpend,
		&client.Data,
		&client.Token,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, data.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &client, nil
}

func (m ClientModel) Update(client *Client) error {
	query := `
        UPDATE clients
               SET token = $1, data = $2, first_name = $3, middle_name = $4, last_name = $5, suffix = $6, email = $7, mobile = $8, username = $9, riman_user_id = $10, status = $11, organization_type = $12, signup_date = $13, anniversary_date = $14, account_type = $15, sponsor_username = $16, member_id = $17, rank = $18, enrollment_date = $19, personal_orders_volume = $20, personal_clients_volume = $21, total_personal_volume = $22, current_month_sp = $23, current_month_bp = $24, last_order_date = $25, last_order_id = $26, last_order_sp = $27, last_order_bp = $28, lifetime_spend = $29, most_recent_12_month_spend = $30
	WHERE id = $31
     RETURNING token`

	args := []any{
		client.Token,
		client.Data,
		client.FirstName,
		client.MiddleName,
		client.LastName,
		client.Suffix,
		client.Email,
		client.Mobile,
		client.Username,
		client.RimanUserId,
		client.Status,
		client.OrganizationType,
		client.SignupDate,
		client.AnniversaryDate,
		client.AccountType,
		client.SponsorUsername,
		client.MemberId,
		client.Rank,
		client.EnrollmentDate,
		client.PersonalOrdersVolume,
		client.PersonalClientsVolume,
		client.TotalPersonalVolume,
		client.CurrentMonthSp,
		client.CurrentMonthBp,
		client.LastOrderDate,
		client.LastOrderId,
		client.LastOrderSp,
		client.LastOrderBp,
		client.LifetimeSpend,
		client.MostRecent12MonthSpend,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&client.Token)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrEditConflict
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}
