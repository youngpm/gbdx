package gbdx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Order holds the ID and all the associated Acquisitions of a submitted order.
type Order struct {
	ID           string        `json:"order_id"`
	Acquisitions []Acquisition `json:"acquisitions"`
}

// Acquisition holds the ID, ordering state, and location (an s3 path) of an order.
type Acquisition struct {
	ID       string `json:"acquisition_id"`
	State    string `json:"state"`
	Location string `json:"location"`
}

// OrderStatus returns the status of the order given the string identifying it.
func (a *Api) OrderStatus(orderID string) (*Order, error) {

	url := fmt.Sprintf("%s%s", endpoints.orders, orderID)
	resp, err := a.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Order status check returned a bad status code: %s", resp.Status)
	}

	var order *Order
	err = json.NewDecoder(resp.Body).Decode(&order)
	if err != nil {
		return nil, fmt.Errorf("Order status failed to decode properly: %v", err)
	}
	return order, err
}

// NewOrder submits a new order.
func (a *Api) NewOrder(IDs ...string) (*Order, error) {

	jsonBody, err := json.Marshal(IDs)
	if err != nil {
		return nil, fmt.Errorf("NewOrder failed to marshall input IDs: %v", err)
	}

	resp, err := a.client.Post(endpoints.orders, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("NewOrder failed to post: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NewOrder post returned a bad status code: %s %v", resp.Status, err)
	}

	var order *Order
	err = json.NewDecoder(resp.Body).Decode(&order)
	if err != nil {
		return nil, fmt.Errorf("NewOrder response failed to decode properly: %v", err)
	}
	return order, err
}

// OrderLocation returns the location of orders in AWS.
func (a *Api) OrderLocation(IDs ...string) (*Order, error) {

	u, err := url.Parse(endpoints.ordersLocation)
	if err != nil {
		return nil, fmt.Errorf("OrderLocation failed parsing the url: %v", err)
	}

	idBytes, err := json.Marshal(IDs)
	if err != nil {
		return nil, fmt.Errorf("NewOrder failed to marshall input IDs: %v", err)
	}

	q := u.Query()
	q.Set("acquisitionIds", string(idBytes))
	u.RawQuery = q.Encode()

	resp, err := a.client.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("OrderLocation failed to get: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OrderLocation get returned a bad status code: %s", resp.Status)
	}

	var order *Order
	err = json.NewDecoder(resp.Body).Decode(&order)
	if err != nil {
		return nil, fmt.Errorf("OrderLocation response failed to decode properly: %v", err)
	}
	return order, err
}

// OrderingHeartbeat checks if the ordering endpoint is alive and well.
func OrderingHeartbeat() error {

	url := endpoints.ordersHeartbeat

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("ordering heartbeat failed %s: %v", resp.Status, err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ordering heartbeat returned a bad status code: %s", resp.Status)
	}
	defer resp.Body.Close()

	return err
}
