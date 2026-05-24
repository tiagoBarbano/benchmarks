package user

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=3"`
	Email string `json:"email" validate:"required,email"`
	CEP   string `json:"cep" validate:"required,len=8,numeric"`
}

type AddressResponse struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
}

type UserResponse struct {
	ID      string           `json:"id"`
	Name    string           `json:"name"`
	Email   string           `json:"email"`
	CEP     string           `json:"cep"`
	Address *AddressResponse `json:"address,omitempty"`
}
