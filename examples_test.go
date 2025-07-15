package jsoon

import "fmt"

type StripeChargeResponse struct {
	ID                 string         `json:"id"`
	Object             string         `json:"object"`
	Amount             int            `json:"amount"`
	AmountRefunded     int            `json:"amount_refunded"`
	BalanceTransaction string         `json:"balance_transaction"`
	Captured           bool           `json:"captured"`
	Created            int            `json:"created"`
	Currency           string         `json:"currency"`
	Customer           string         `json:"customer"`
	Livemode           bool           `json:"livemode"`
	Outcome            *StripeOutcome `json:"outcome"`
	Paid               bool           `json:"paid"`
	Refunded           bool           `json:"refunded"`
	Refunds            *StripeRefunds `json:"refunds"`
	Source             *StripeSource  `json:"source"`
	Status             string         `json:"status"`
}

func (s *StripeChargeResponse) MarshalJsoon(e *Encoder) (err error) {
	e.String("id", s.ID)
	e.String("object", s.Object)
	e.Number("amount", float64(s.Amount))
	e.Number("amount_refunded", float64(s.AmountRefunded))
	e.String("balance_transaction", s.BalanceTransaction)
	e.Bool("captured", s.Captured)
	e.Number("created", float64(s.Created))
	e.String("currency", s.Currency)
	e.String("customer", s.Customer)
	e.Bool("livemode", s.Livemode)
	e.Object("outcome", s.Outcome)
	e.Bool("paid", s.Paid)
	e.Bool("refunded", s.Refunded)
	e.Object("refunds", s.Refunds)
	e.Object("source", s.Source)
	e.String("status", s.Status)
	return
}

func (s *StripeChargeResponse) UnmarshalJsoon(key string, val *Value) (err error) {
	switch key {
	case "id":
		s.ID, err = val.String()
	case "object":
		s.Object, err = val.String()
	case "amount":
		var a float64
		a, err = val.Number()
		s.Amount = int(a)
	case "amount_refunded":
		var ar float64
		ar, err = val.Number()
		s.AmountRefunded = int(ar)
	case "balance_transaction":
		s.BalanceTransaction, err = val.String()
	case "captured":
		s.Captured, err = val.Bool()
	case "created":
		var cr float64
		cr, err = val.Number()
		s.Created = int(cr)
	case "currency":
		s.Currency, err = val.String()
	case "customer":
		s.Customer, err = val.String()
	case "livemode":
		s.Livemode, err = val.Bool()
	case "outcome":
		s.Outcome = &StripeOutcome{}
		err = val.Object(s.Outcome)
	case "paid":
		s.Paid, err = val.Bool()
	case "refunded":
		s.Refunded, err = val.Bool()
	case "refunds":
		s.Refunds = &StripeRefunds{}
		err = val.Object(s.Refunds)
	case "source":
		s.Source = &StripeSource{}
		err = val.Object(s.Source)
	case "status":
		s.Status, err = val.String()
	}

	return
}

func (s *StripeChargeResponse) Equals(os *StripeChargeResponse) (err error) {
	return
}

type StripeOutcome struct {
	NetworkStatus string `json:"network_status"`
	RiskLevel     string `json:"risk_level"`
	SellerMessage string `json:"seller_message"`
	Type          string `json:"type"`
}

func (s *StripeOutcome) MarshalJsoon(e *Encoder) (err error) {
	e.String("network_status", s.NetworkStatus)
	e.String("risk_level", s.RiskLevel)
	e.String("seller_message", s.SellerMessage)
	e.String("type", s.Type)
	return
}

func (s *StripeOutcome) UnmarshalJsoon(key string, val *Value) (err error) {
	switch key {
	case "network_status":
		s.NetworkStatus, err = val.String()
	case "risk_level":
		s.RiskLevel, err = val.String()
	case "seller_message":
		s.SellerMessage, err = val.String()
	case "type":
		s.Type, err = val.String()
	}

	return
}

func (s *StripeOutcome) Equals(os *StripeOutcome) (err error) {
	return
}

type StripeRefunds struct {
	Object     string `json:"object"`
	HasMore    bool   `json:"has_more"`
	TotalCount int    `json:"total_count"`
	URL        string `json:"url"`
}

func (s *StripeRefunds) MarshalJsoon(e *Encoder) (err error) {
	e.String("object", s.Object)
	e.Bool("has_more", s.HasMore)
	e.Number("total_count", float64(s.TotalCount))
	e.String("url", s.URL)
	return
}

func (s *StripeRefunds) UnmarshalJsoon(key string, val *Value) (err error) {
	switch key {
	case "object":
		s.Object, err = val.String()
	case "has_more":
		s.HasMore, err = val.Bool()
	case "total_count":
		var tc float64
		tc, err = val.Number()
		s.TotalCount = int(tc)
	case "url":
		s.URL, err = val.String()
	}

	return
}

func (s *StripeRefunds) Equals(os *StripeRefunds) (err error) {
	return
}

type StripeSource struct {
	ID       string `json:"id"`
	Object   string `json:"object"`
	Brand    string `json:"brand"`
	Country  string `json:"country"`
	Customer string `json:"customer"`
	CVCCheck string `json:"cvc_check"`
	ExpMonth int    `json:"exp_month"`
	ExpYear  int    `json:"exp_year"`
	Funding  string `json:"funding"`
	Last4    string `json:"last4"`
}

func (s *StripeSource) MarshalJsoon(e *Encoder) (err error) {
	e.String("id", s.ID)
	e.String("object", s.Object)
	e.String("brand", s.Brand)
	e.String("country", s.Country)
	e.String("customer", s.Customer)
	e.String("cvc_check", s.CVCCheck)
	e.Number("exp_month", float64(s.ExpMonth))
	e.Number("exp_year", float64(s.ExpYear))
	e.String("funding", s.Funding)
	e.String("last4", s.Last4)
	return
}

func (s *StripeSource) UnmarshalJsoon(key string, val *Value) (err error) {
	switch key {
	case "id":
		s.ID, err = val.String()
	case "object":
		s.Object, err = val.String()
	case "brand":
		s.Brand, err = val.String()
	case "country":
		s.Country, err = val.String()
	case "customer":
		s.Customer, err = val.String()
	case "cvc_check":
		s.CVCCheck, err = val.String()
	case "exp_month":
		var em float64
		em, err = val.Number()
		s.ExpMonth = int(em)
	case "exp_year":
		var ey float64
		ey, err = val.Number()
		s.ExpYear = int(ey)
	case "funding":
		s.Funding, err = val.String()
	case "last4":
		s.Last4, err = val.String()
	}

	return
}

func (s *StripeSource) Equals(os *StripeSource) (err error) {
	if s.ID != os.ID {
		return fmt.Errorf("id's don't match: %s | %s", s.ID, os.ID)
	}

	if s.Object != os.Object {
		///		// return fmt.Errorf("objects's don't match: %s | %s")
	}

	if s.Brand != os.Brand {
		//	// return fmt.Errorf("brands's don't match: %s | %s")
	}

	if s.Country != os.Country {
		//	// return fmt.Errorf("countries's don't match: %s | %s")
	}

	if s.Customer != os.Customer {
		// return fmt.Errorf("id's don't match: %s | %s")
	}

	if s.CVCCheck != os.CVCCheck {
		// return fmt.Errorf("id's don't match: %s | %s")
	}

	if s.ExpMonth != os.ExpMonth {
		// return fmt.Errorf("id's don't match: %s | %s")
	}

	if s.ExpYear != os.ExpYear {
		// return fmt.Errorf("id's don't match: %s | %s")
	}

	if s.Last4 != os.Last4 {
		// return fmt.Errorf("id's don't match: %s | %s")
	}

	return
}
