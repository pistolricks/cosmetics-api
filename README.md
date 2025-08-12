# KBeauty API

SCOPES https://shopify.dev/docs/api/admin-rest/usage/access-scopes

Using headless chromium to automate orders.

## read_assigned_fulfillment_orders

# Riman API

GetRimanCustomer=https://profile-api.riman.com/api/v1/users/51311

GetRimanReferrer=https://profile-api.riman.com/api/v1/member/referrer
PatchCartReferrer=https://cart-api.riman.com/api/v1/shopping/64586b7d-a56d-409a-999c-a934e9b8fa78
--{"mainFk": 51311,"referrerSiteUrl": "WeKBeauty","countryCode": "US","culture": "en-US"}

GetMemberShippingAddress=https://profile-api.riman.com/api/v1/member/address/shipping?isOTG=true
GetBillingAddress=https://profile-api.riman.com/api/v1/member/address/billing
GetShippingAddress=https://profile-api.riman.com/api/v2/shipping-address
AddressRequest=https://profile-api.riman.com/api/v1/address/auto-complete?addressString=666%20Capela%20Way&country=US&latitude=undefined&longitude=undefined

PatchCartShippingTypeFK=https://cart-api.riman.com/api/v1/shopping/64586b7d-a56d-409a-999c-a934e9b8fa78
--{"shippingTypeFk": 1163 }

PatchCartMainFK=https://cart-api.riman.com/api/v1/shopping/64586b7d-a56d-409a-999c-a934e9b8fa78
--{"mainFk": 51311 }

ShippingCostRequest=https://cart-api.riman.com/api/v2/shipping/shipping/options?cartKey=64586b7d-a56d-409a-999c-a934e9b8fa78&country=US&state=CA&zip=95831&city=Sacramento&address1=666%20Capela%20Way&address2=&address3=undefined&mainOrderPK=0

PutShippingAddress=https://cart-api.riman.com/api/v2/order/shipping/address/2306672
--{"mainOrdersFK": 2306672,"firstName": "Xuan","lastName": "Nguyen","address1": "666 Capela Way","address2": "","
city": "Sacramento","state": "CA","postalCode": "95831","country": "US","phone": "9164029611","fax": "","email": "
embreday9@gmail.com","verifiedShippingAddress": true }

PutSetShipping=https://cart-api.riman.com/api/v2/order/set-shipping
--{"orderPk": 2306672,"chartTypeFk": 1163,"signatureRequired": false,"cartKey": "64586b7d-a56d-409a-999c-a934e9b8fa78"}

GetCart=https://cart-api.riman.com/api/v1/shopping/72357ede-62c0-4e03-9904-f47be1a7f900

GetRemainingProduct=https://cart-api.riman.com/api/v2/products/remaining-products?mainPK=51311&productIds=50466,51343
GetHasCompleted=https://cart-api.riman.com/api/v1/order/51311/has-completed

GetRepSite=https://profile-api.riman.com/api/v1/repsite

PutCart=https://cart-api.riman.com/api/v2/order/items/add-from-cart/2306672/64586b7d-a56d-409a-999c-a934e9b8fa78

RimanConfirmOrderDetails=https://cart-api.riman.com/api/v1/order/2306672/confirmation-details

PostCurrency=https://cms-api.riman.com/api/v1/format/currency-number

GetRimanDollarsBalance=https://cart-api.riman.com/api/v1/order/getRimanDollarsBalance/51311/2306672

PostOrder=https://cart-api.riman.com/api/v2/order
payload [
cartKey: "72357ede-62c0-4e03-9904-f47be1a7f900"
countryCode: "US"
mainId: 51311
mainOrderType: 4
salesCampaignFK: null
]

GetLocaleInfo=https://cms-api.riman.com/api/v1/application/locale-info
GetProductsNotLoggedIn=https://cart-api.riman.com/api/v2/products?cartType=R&countryCode=US&culture=en-US&isCart=true&repSiteUrl=rmnsocial
GetProducts=https://cart-api.riman.com/api/v2/products?cartType=R&countryCode=US&culture=en-US&isCart=true&repSiteUrl=WeKBeauty

GetOrders=https://cart-api.riman.com/api/v1/orders?mainSiteUrl=2043124962&getEnrollerOrders=&getCustomerOrders=&orderNumber=&shipmentNumber=&trackingNumber=&isRefunded=&paidStatus=true&orderType=&orderLevel=&weChatOrderNumber=&startDate=&endDate=&offset=0&limit=20&orderBy=-mainOrdersPK
GetShipping=https://cart-api.riman.com/api/v1/orders/{OrderNumber}/shipment-products

## Don't forget the Processes in Excalidraw

