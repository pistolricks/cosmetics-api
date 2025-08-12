package shopify

import (
	"context"
	"fmt"
	"time"

	goshopify "github.com/bold-commerce/go-shopify/v4"
)

type ShopifyProduct struct {
	Id             int64     `json:"id"`
	Title          string    `json:"title"`
	BodyHtml       string    `json:"body_html"`
	Vendor         string    `json:"vendor"`
	ProductType    string    `json:"product_type"`
	Handle         string    `json:"handle"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	PublishedAt    time.Time `json:"published_at"`
	PublishedScope string    `json:"published_scope"`
	Tags           string    `json:"tags"`
	Status         string    `json:"status"`
	Options        []struct {
		Id        int64    `json:"id"`
		ProductId int64    `json:"product_id"`
		Name      string   `json:"name"`
		Position  int      `json:"position"`
		Values    []string `json:"values"`
	} `json:"options"`
	Variants []struct {
		Id                   int64     `json:"id"`
		ProductId            int64     `json:"product_id"`
		Title                string    `json:"title"`
		Sku                  string    `json:"sku"`
		Position             int       `json:"position"`
		Grams                int       `json:"grams"`
		InventoryPolicy      string    `json:"inventory_policy"`
		Price                string    `json:"price"`
		FulfillmentService   string    `json:"fulfillment_service"`
		InventoryManagement  string    `json:"inventory_management"`
		InventoryItemId      int64     `json:"inventory_item_id"`
		Option1              string    `json:"option1"`
		CreatedAt            time.Time `json:"created_at"`
		UpdatedAt            time.Time `json:"updated_at"`
		Taxable              bool      `json:"taxable"`
		Barcode              string    `json:"barcode"`
		InventoryQuantity    int       `json:"inventory_quantity"`
		Weight               string    `json:"weight"`
		WeightUnit           string    `json:"weight_unit"`
		OldInventoryQuantity int       `json:"old_inventory_quantity"`
		RequiresShipping     bool      `json:"requires_shipping"`
		AdminGraphqlApiId    string    `json:"admin_graphql_api_id"`
	} `json:"variants"`
	Image struct {
		Id                int64     `json:"id"`
		ProductId         int64     `json:"product_id"`
		Position          int       `json:"position"`
		CreatedAt         time.Time `json:"created_at"`
		UpdatedAt         time.Time `json:"updated_at"`
		Width             int       `json:"width"`
		Height            int       `json:"height"`
		Src               string    `json:"src"`
		Alt               string    `json:"alt"`
		AdminGraphqlApiId string    `json:"admin_graphql_api_id"`
	} `json:"image"`
	Images []struct {
		Id                int64     `json:"id"`
		ProductId         int64     `json:"product_id"`
		Position          int       `json:"position"`
		CreatedAt         time.Time `json:"created_at"`
		UpdatedAt         time.Time `json:"updated_at"`
		Width             int       `json:"width"`
		Height            int       `json:"height"`
		Src               string    `json:"src"`
		Alt               string    `json:"alt"`
		AdminGraphqlApiId string    `json:"admin_graphql_api_id"`
	} `json:"images"`
	TemplateSuffix    string `json:"template_suffix"`
	AdminGraphqlApiId string `json:"admin_graphql_api_id"`
}

func GetProducts(client *goshopify.Client) ([]goshopify.Product, int, error) {

	// Get all collections
	count, err := client.Product.Count(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
	}

	products, err := client.Product.List(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
	}

	return products, count, err
}
