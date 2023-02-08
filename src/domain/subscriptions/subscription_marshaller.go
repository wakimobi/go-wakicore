package subscriptions

import (
	"encoding/json"
)

type PublicSubscription struct {
	ID     int64  `json:"id" xml:"id"`
	Status string `json:"status" xml:"status"`
}

type PrivateSubscription struct {
	ID     int    `json:"id" xml:"id"`
	Status string `json:"status" xml:"status"`
	Msisdn string `json:"msisdn" xml:"msisdn"`
}

func (sub *Subscription) Marshall(isPublic bool) interface{} {
	if isPublic {
		return PublicSubscription{
			ID:     sub.ID,
			Status: "Success",
		}
	}
	subJson, _ := json.Marshal(sub)
	var privateSubscription PrivateSubscription
	json.Unmarshal(subJson, &privateSubscription)
	return privateSubscription
}
