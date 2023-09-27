package model

type DBModel struct {
	Order  *Order  `json:"order,omitempty" bson:"order,omitempty"`
	Wallet *Wallet `json:"wallet,omitempty" bson:"wallet,omitempty"`
	Coin   *Coin   `json:"coin,omitempty" bson:"coin,omitempty"`
}
