package user

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string             `json:"name" bson:"name"`
	Email   string             `json:"email" bson:"email"`
	Address *Address           `json:"address,omitempty" bson:"address,omitempty"`
	CEP     string             `json:"cep,omitempty" bson:"cep,omitempty"`
}

type Address struct {
	CEP         string `json:"cep" bson:"cep"`
	Logradouro  string `json:"logradouro" bson:"logradouro"`
	Complemento string `json:"complemento" bson:"complemento"`
	Bairro      string `json:"bairro" bson:"bairro"`
	Localidade  string `json:"localidade" bson:"localidade"`
	UF          string `json:"uf" bson:"uf"`
}

type DtoUserResponse struct {
	ID            string  `json:"id" bson:"_id"`
	Empresa       string  `json:"empresa" bson:"empresa"`
	Cotacao_final float64 `json:"cotacao_final" bson:"cotacao_final"`
	Created_at    string  `json:"created_at" bson:"created_at"`
	Update_at     string  `json:"updated_at" bson:"updated_at"`
	Deleted       bool    `json:"deleted" bson:"deleted"`
}
