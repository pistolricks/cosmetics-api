package v2

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"github.com/vinhluan/go-graphql-client"
	graphify "github.com/vinhluan/go-shopify-graphql"
)

type ServiceV2 struct {
	DB      *sql.DB
	Client  *graphify.Client
	Graphql *graphql.GraphQL
}

func (v2 ServiceV2) Select() {

}

func (v ServiceV2) GraphQLClient() graphql.GraphQL {
	return *v.Graphql
}

func (v ServiceV2) Mutate(ctx context.Context, m interface{}, variables map[string]interface{}) error {
	var retries = 0
	for {
		r, err := c.gql.Mutate(ctx, m, variables)
		if err != nil {
			if r != nil {
				wait := CalculateWaitTime(r.Extensions)
				if wait > 0 {
					retries++
					time.Sleep(wait)
					continue
				}
			}
			if IsConnectionError(err) {
				retries++
				if retries > c.retries {
					return fmt.Errorf("after %v tries: %w", retries, err)
				}
				time.Sleep(time.Duration(retries) * time.Second)
				continue
			}
			return err
		}
		break
	}

	return nil
}

func (v ServiceV2) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	var retries = 0
	for {
		r, err := c.gql.Query(ctx, q, variables)
		if err != nil {
			if r != nil {
				wait := CalculateWaitTime(r.Extensions)
				if wait > 0 {
					retries++
					time.Sleep(wait)
					continue
				}
			}
			if uerr, isURLErr := err.(*url.Error); isURLErr && (uerr.Timeout() || uerr.Temporary()) || IsConnectionError(err) {
				retries++
				if retries > c.retries {
					return fmt.Errorf("after %v tries: %w", retries, err)
				}
				time.Sleep(time.Duration(retries) * time.Second)
				continue
			}
			return err
		}
		break
	}

	return nil
}

func (v ServiceV2) QueryString(ctx context.Context, q string, variables map[string]interface{}, out interface{}) error {
	var retries = 0
	for {
		r, err := c.gql.QueryString(ctx, q, variables, out)
		if err != nil {
			if r != nil {
				wait := CalculateWaitTime(r.Extensions)
				if wait > 0 {
					retries++
					time.Sleep(wait)
					continue
				}
			}
			if uerr, isURLErr := err.(*url.Error); isURLErr && (uerr.Timeout() || uerr.Temporary()) || IsConnectionError(err) {
				retries++
				if retries > c.retries {
					return fmt.Errorf("after %v tries: %w", retries, err)
				}
				time.Sleep(time.Duration(retries) * time.Second)
				continue
			}
			return err
		}
		break
	}

	return nil
}
