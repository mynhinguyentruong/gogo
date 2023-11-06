package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72/webhook"
	stripe "github.com/stripe/stripe-go/v72"
	// "github.com/stripe/stripe-go/v76/checkout/session"
)


func handleWebhookRoute (c *gin.Context) {
  const MaxBodyBytes = int64(65536)

  reqBody := http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

  body, err := io.ReadAll(reqBody)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
    c.AbortWithStatus(http.StatusServiceUnavailable)
    return
  }

  // read request body that was sent from Stripe
  // ?client_reference_id=123

  endpointSecret:= os.Getenv("endpointSecret")

  if endpointSecret == "" {
      log.Fatal("set env endpointSecret")
  }

  event, err := webhook.ConstructEvent(body, c.Request.Header.Get("Stripe-Signature"), endpointSecret)

  if err != nil {
    fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
    c.AbortWithStatus(http.StatusBadRequest) // Return a 400 error on a bad signature
    return 

  }

 // Handle the checkout.session.completed event
  if event.Type == "checkout.session.completed" {
    var session stripe.CheckoutSession
    err := json.Unmarshal(event.Data.Raw, &session)
    if err != nil {
      fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
      c.AbortWithStatus(http.StatusBadRequest)
      return
    }

    if customer_id := session.ClientReferenceID; customer_id == "" {
      fmt.Errorf("empty client_reference_id")
      c.AbortWithStatusJSON(http.StatusBadRequest, "empty client_reference_id")
      return
    } else {
      fmt.Println("client_reference_id: ", customer_id)
      fmt.Println("the whole event look like this: ", session)
          FulfillOrder(customer_id)
    } 

    // params := &stripe.CheckoutSessionParams{}
    // params.AddExpand("line_items")

    // Retrieve the session. If you require line items in the response, you may include them by expanding line_items.
    // sessionWithLineItems, _ := session.Get(session.ID, params)
    // lineItems := sessionWithLineItems.LineItems
    // Fulfill the purchase...

  }
  // on checkout.session.completed
  // check for client_reference_id
  // update the credit in database
  // based on product id

  fmt.Println("event: ", event)

  c.IndentedJSON(http.StatusOK, event)
  
}

func FulfillOrder(customer_id string) {
 fmt.Println("increase their credit here: ", customer_id)
}
