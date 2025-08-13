package riman

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/pistolricks/kbeauty-api/internal/data"
)

type ClientModel struct {
	DB *sql.DB
}

var AnonymousClient = &Client{}

func (c *Client) IsAnonymous() bool {
	return c == AnonymousClient
}

func (m ClientModel) GetByClientUsername(username string) (*Client, error) {
	query := `
        SELECT id, created_at, first_name, middle_name, last_name, suffix, email, mobile, username, riman_user_id, status, organization_type, signup_date, anniversary_date, account_type, sponsor_username, member_id, rank, enrollment_date, personal_orders_volume, personal_clients_volume, total_personal_volume, current_month_sp, current_month_bp, last_order_date, last_order_id, last_order_sp, last_order_bp, lifetime_spend, most_recent_12_month_spend, data
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

func (m ClientModel) GetByClientEmail(email string) (*Client, error) {
	query := `
        SELECT id, created_at, first_name, middle_name, last_name, suffix, email, mobile, username, riman_user_id, status, organization_type, signup_date, anniversary_date, account_type, sponsor_username, member_id, rank, enrollment_date, personal_orders_volume, personal_clients_volume, total_personal_volume, current_month_sp, current_month_bp, last_order_date, last_order_id, last_order_sp, last_order_bp, lifetime_spend, most_recent_12_month_spend, data
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
	SELECT count(*) OVER(),id, created_at, first_name, middle_name, last_name, suffix, email, mobile, username, riman_user_id, status, organization_type, signup_date, anniversary_date, account_type, sponsor_username, member_id, rank, enrollment_date, personal_orders_volume, personal_clients_volume, total_personal_volume, current_month_sp, current_month_bp, last_order_date, last_order_id, last_order_sp, last_order_bp, lifetime_spend, most_recent_12_month_spend, data
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
        SELECT clients.id, clients.created_at, clients.first_name, clients.middle_name, clients.last_name, clients.suffix, clients.email, clients.mobile, clients.username, clients.riman_user_id, clients.status, clients.organization_type, clients.signup_date, clients.anniversary_date, clients.account_type, clients.sponsor_username, clients.member_id, clients.rank, clients.enrollment_date, clients.personal_orders_volume, clients.personal_clients_volume, clients.total_personal_volume, clients.current_month_sp, clients.current_month_bp, clients.last_order_date, clients.last_order_id, clients.last_order_sp, clients.last_order_bp, clients.lifetime_spend, clients.most_recent_12_month_spend, clients.data
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
