package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/riman/products", app.RimanApiListProductsHandler)

	router.HandlerFunc(http.MethodGet, "/v1/orders", app.listShopifyOrdersHandler)

	router.HandlerFunc(http.MethodPost, "/v1/attributes/metafield", app.updateOrderMetaField)

	router.HandlerFunc(http.MethodGet, "/v1/process/orders", app.processShopifyOrders)

	router.HandlerFunc(http.MethodGet, "/v1/process/order/cart", app.getCartHandler)

	router.HandlerFunc(http.MethodPost, "/v1/process/order", app.processShopifyOrder)

	router.HandlerFunc(http.MethodGet, "/v1/process/order/billing", app.inputBillingHandler)

	router.HandlerFunc(http.MethodGet, "/v1/shopify/products", app.ShopifyApiListProductsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/shopify/orders/note", app.orderUpdateHandler)

	router.HandlerFunc(http.MethodGet, "/v1/riman/test/event", app.getEventHandler)

	router.HandlerFunc(http.MethodGet, "/v1/riman/orders", app.listRimanOrders)
	router.HandlerFunc(http.MethodGet, "/v1/riman/shipment", app.getShipmentHandler)
	router.HandlerFunc(http.MethodGet, "/v1/riman/tracking", app.trackingHandler)

	router.HandlerFunc(http.MethodPost, "/v1/riman/login", app.chromeLoginHandler)
	router.HandlerFunc(http.MethodGet, "/v1/riman/home", app.chromeHomePageHandler)

	router.HandlerFunc(http.MethodPost, "/v1/riman/logout", app.apiRimanLogoutHandler)

	//	router.HandlerFunc(http.MethodPost, "/v1/riman/login", app.createRimanTokenHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users/register", app.registerUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users/login", app.createAuthenticationTokenHandler)

	router.HandlerFunc(http.MethodGet, "/v1/orders/all", app.listAllShopifyOrders)
	router.HandlerFunc(http.MethodGet, "/v1/orders/list", app.listShopifyOrders)

	router.HandlerFunc(http.MethodGet, "/v1/clients", app.listClientsHandler)

	/*

		router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
		router.HandlerFunc(http.MethodPut, "/v1/users/password", app.updateUserPasswordHandler)


		router.HandlerFunc(http.MethodPost, "/v1/tokens/activation", app.createActivationTokenHandler)
		router.HandlerFunc(http.MethodPost, "/v1/tokens/password-reset", app.createPasswordResetTokenHandler)
	*/

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(app.authenticateClient(router))))))
}
