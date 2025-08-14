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

	router.HandlerFunc(http.MethodGet, "/v1/platform/products", app.ShopifyApiListProductsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/platform/notes", app.orderUpdateHandler)
	router.HandlerFunc(http.MethodPost, "/v1/platform/fields", app.updateOrderFields)
	router.HandlerFunc(http.MethodPost, "/v1/platform/tracking", app.updateFulfillmentHandler)
	router.HandlerFunc(http.MethodGet, "/v1/platform/orders", app.listShopifyOrdersHandler)
	router.HandlerFunc(http.MethodGet, "/v1/platform/orders/list", app.listShopifyOrders)
	router.HandlerFunc(http.MethodGet, "/v1/platform/orders/all", app.listAllShopifyOrders)

	router.HandlerFunc(http.MethodPost, "/v1/vendors/login", app.createRimanSessionHandler)
	router.HandlerFunc(http.MethodGet, "/v1/vendors/products", app.RimanApiListProductsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/vendors/carts", app.getCartHandler)
	router.HandlerFunc(http.MethodGet, "/v1/vendors/orders", app.listRimanOrders)
	router.HandlerFunc(http.MethodGet, "/v1/vendors/tracking", app.trackingHandler)
	router.HandlerFunc(http.MethodPost, "/v1/vendors/shipment", app.getShipmentHandler)
	router.HandlerFunc(http.MethodGet, "/v1/vendors/clients", app.listClientsHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users/login", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users/register", app.registerUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/navigation/login", app.chromeLoginHandler)
	router.HandlerFunc(http.MethodPost, "/v1/navigation/logout", app.apiRimanLogoutHandler)
	router.HandlerFunc(http.MethodGet, "/v1/navigation/home", app.chromeHomePageHandler)

	router.HandlerFunc(http.MethodPost, "/v1/navigation/orders", app.processShopifyOrder)
	router.HandlerFunc(http.MethodGet, "/v1/navigation/orders", app.processShopifyOrders)

	router.HandlerFunc(http.MethodGet, "/v1/navigation/test/event", app.getEventHandler)

	router.HandlerFunc(http.MethodGet, "/v1/navigation/shipping", app.inputShippingHandler)
	router.HandlerFunc(http.MethodGet, "/v1/navigation/billing", app.inputBillingHandler)
	router.HandlerFunc(http.MethodGet, "/v1/navigation/forms/submit", app.submitFormHandler)

	/*

		router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
		router.HandlerFunc(http.MethodPut, "/v1/users/password", app.updateUserPasswordHandler)


		router.HandlerFunc(http.MethodPost, "/v1/tokens/activation", app.createActivationTokenHandler)
		router.HandlerFunc(http.MethodPost, "/v1/tokens/password-reset", app.createPasswordResetTokenHandler)
	*/

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(app.authenticateClient(router))))))
}
