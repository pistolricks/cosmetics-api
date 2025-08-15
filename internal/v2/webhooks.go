package v2

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pistolricks/cosmetics-api/graph/model"
	graphify "github.com/vinhluan/go-shopify-graphql"
)

type WebhookService interface {
	CreateWebhookSubscription(ctx context.Context, topic model.WebhookSubscriptionTopic, input model.WebhookSubscriptionInput) (*model.WebhookSubscription, error)
	CreateEventBridgeWebhookSubscription(ctx context.Context, topic model.WebhookSubscriptionTopic, input model.EventBridgeWebhookSubscriptionInput) (*model.WebhookSubscription, error)
	ListWebhookSubscriptions(ctx context.Context, topics []model.WebhookSubscriptionTopic) ([]*model.WebhookSubscription, error)
	DeleteWebhook(ctx context.Context, webhookID string) (deletedID *string, err error)
	UpdateWebhookSubscription(ctx context.Context, webhookID string, input model.WebhookSubscriptionInput) (*model.WebhookSubscription, error)
}

type WebhookServiceOp struct {
	client *Client
}

var _ WebhookService = &WebhookServiceOp{}

type mutationWebhookCreate struct {
	WebhookCreateResult *model.WebhookSubscriptionCreatePayload `graphql:"webhookSubscriptionCreate(topic: $topic, webhookSubscription: $webhookSubscription)" json:"webhookSubscriptionCreate"`
}

type mutationWebhookUpdate struct {
	WebhookUpdateResult *model.WebhookSubscriptionUpdatePayload `graphql:"webhookSubscriptionUpdate(id: $id, webhookSubscription: $webhookSubscription)" json:"webhookSubscriptionUpdate"`
}

type mutationWebhookDelete struct {
	WebhookDeleteResult *model.WebhookSubscriptionDeletePayload `graphql:"webhookSubscriptionDelete(id: $id)" json:"webhookSubscriptionDelete"`
}

type mutationEventBridgeWebhookCreate struct {
	EventBridgeWebhookCreateResult *model.EventBridgeWebhookSubscriptionCreatePayload `graphql:"eventBridgeWebhookSubscriptionCreate(topic: $topic, webhookSubscription: $webhookSubscription)" json:"eventBridgeWebhookSubscriptionCreate"`
}

// NOTE: Have to use this because writeQuery function will not write structs that implements UnmarshalJSON function
const webhookSubscriptionMutationSelects = `
userErrors {
	field
	message
}
webhookSubscription {
	apiVersion {
		displayName
		handle
		supported
	}
	callbackUrl
	createdAt
	format
	id
	includeFields
	legacyResourceId
	metafieldNamespaces
	privateMetafieldNamespaces
	topic
	updatedAt
	endpoint {
		__typename
		...on WebhookEventBridgeEndpoint {
			arn
		}
		...on WebhookHttpEndpoint {
			callbackUrl
		}
	}
}`

type WebhookV2 struct {
	DB     *sql.DB
	Client *graphify.Client
}

func (s WebhookServiceOp) CreateWebhookSubscription(ctx context.Context, topic model.WebhookSubscriptionTopic, input model.WebhookSubscriptionInput) (*model.WebhookSubscription, error) {
	m := mutationWebhookCreate{}
	vars := map[string]interface{}{
		"topic":               topic,
		"webhookSubscription": input,
	}
	err := s.client.Mutate(ctx, &m, vars)
	if err != nil {
		return nil, fmt.Errorf("mutation: %w", err)
	}

	if len(m.WebhookCreateResult.UserErrors) > 0 {
		return nil, fmt.Errorf("%+v", m.WebhookCreateResult.UserErrors)
	}

	return m.WebhookCreateResult.WebhookSubscription, nil
}

func (s WebhookServiceOp) CreateEventBridgeWebhookSubscription(ctx context.Context, topic model.WebhookSubscriptionTopic, input model.EventBridgeWebhookSubscriptionInput) (*model.WebhookSubscription, error) {
	m := mutationEventBridgeWebhookCreate{}
	vars := map[string]interface{}{
		"topic":               topic,
		"webhookSubscription": input,
	}

	err := s.client.Mutate(ctx, &m, vars)
	if err != nil {
		return nil, fmt.Errorf("mutation: %w", err)
	}

	if len(m.EventBridgeWebhookCreateResult.UserErrors) > 0 {
		return nil, fmt.Errorf("%+v", m.EventBridgeWebhookCreateResult.UserErrors)
	}

	return m.EventBridgeWebhookCreateResult.WebhookSubscription, nil
}

func (s WebhookServiceOp) DeleteWebhook(ctx context.Context, webhookID string) (*string, error) {
	m := mutationWebhookDelete{}
	vars := map[string]interface{}{
		"id": webhookID,
	}
	err := s.client.Mutate(ctx, &m, vars)
	if err != nil {
		return nil, fmt.Errorf("mutation: %w", err)
	}

	if len(m.WebhookDeleteResult.UserErrors) > 0 {
		return nil, fmt.Errorf("%+v", m.WebhookDeleteResult.UserErrors)
	}
	return m.WebhookDeleteResult.DeletedWebhookSubscriptionID, nil
}

func (s WebhookServiceOp) ListWebhookSubscriptions(ctx context.Context, topics []model.WebhookSubscriptionTopic) ([]*model.WebhookSubscription, error) {
	queryFormat := `query webhookSubscriptions($first: Int!, $topics: [WebhookSubscriptionTopic!]%s) {
		webhookSubscriptions(first: $first, topics: $topics%s) {
			edges {
				cursor
				node {
					id
					topic
					apiVersion {
						displayName
						handle
						supported
					}
					endpoint {
						__typename
						... on WebhookHttpEndpoint {
							callbackUrl
						}
						... on WebhookEventBridgeEndpoint{
							arn
						}
					}
					callbackUrl
					format
					topic
					includeFields
					createdAt
					updatedAt
				}
			}
			pageInfo {
				hasNextPage
			}
		}
	}`

	var (
		cursor string
		vars   = map[string]interface{}{
			"first":  200,
			"topics": topics,
		}
		output = make([]*model.WebhookSubscription, 0)
	)
	for {
		var (
			query string
			out   model.QueryRoot
		)
		if cursor != "" {
			vars["after"] = cursor
			query = fmt.Sprintf(queryFormat, ", $after: String", ", after: $after")
		} else {
			query = fmt.Sprintf(queryFormat, "", "")
		}
		err := s.client.QueryString(ctx, query, vars, &out)
		if err != nil {
			return nil, fmt.Errorf("query: %w", err)
		}
		for _, wh := range out.WebhookSubscriptions.Edges {
			output = append(output, wh.Node)
		}
		if out.WebhookSubscriptions.PageInfo.HasNextPage {
			cursor = out.WebhookSubscriptions.Edges[len(out.WebhookSubscriptions.Edges)-1].Cursor
		} else {
			break
		}
	}
	return output, nil
}

func (s WebhookServiceOp) UpdateWebhookSubscription(ctx context.Context, webhookID string, input model.WebhookSubscriptionInput) (*model.WebhookSubscription, error) {
	m := mutationWebhookUpdate{}
	vars := map[string]interface{}{
		"id":                  webhookID,
		"webhookSubscription": input,
	}
	err := s.client.Mutate(ctx, &m, vars)
	if err != nil {
		return nil, fmt.Errorf("mutation: %w", err)
	}

	if len(m.WebhookUpdateResult.UserErrors) > 0 {
		return nil, fmt.Errorf("%+v", m.WebhookUpdateResult.UserErrors)
	}

	return m.WebhookUpdateResult.WebhookSubscription, nil
}
