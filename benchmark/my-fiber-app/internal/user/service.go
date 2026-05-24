package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"my-fiber-app/pkg/redis"
)

type Service interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetAll(ctx context.Context) ([]User, error)
	GetByID(ctx context.Context, id string) (*DtoUserResponse, error)
	Update(ctx context.Context, id string, user *User) (*User, error)
	Delete(ctx context.Context, id string) error
}

type userService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &userService{repo: repo}
}

const redisUserKeyPrefix = "user:"

func (s *userService) Create(ctx context.Context, user *User) (*User, error) {
	if user.CEP != "" {
		addr, err := fetchAddressFromViaCEP(user.CEP)
		if err != nil {
			// Pode optar por ignorar erro do cep e continuar, ou retornar erro
			return nil, err
		}
		user.Address = addr
	}
	createdUser, err := s.repo.Create(ctx, user)
	return createdUser, err
}

func (s *userService) GetAll(ctx context.Context) ([]User, error) {
	return s.repo.GetAll(ctx)
}

func (s *userService) GetByID(ctx context.Context, id string) (*DtoUserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil || user == nil {
		return user, err
	}

	return user, nil
}

// func (s *userService) GetByID(ctx context.Context, id string) (*User, error) {
// 	cacheKey := redisUserKeyPrefix + id

// 	// 1. Tenta buscar no Redis
// 	cached, err := redis.Client.Get(ctx, cacheKey).Result()
// 	if err == nil {
// 		var user User
// 		if err := json.Unmarshal([]byte(cached), &user); err == nil {
// 			return &user, nil
// 		}
// 	}

// 	// 2. Se não achar no cache, busca no Mongo
// 	user, err := s.repo.GetByID(ctx, id)
// 	if err != nil || user == nil {
// 		return user, err
// 	}

// 	// 3. Salva no cache Redis (expire 60 segundos)
// 	data, err := json.Marshal(user)
// 	if err == nil {
// 		redis.Client.Set(ctx, cacheKey, data, 60*time.Second)
// 	}

// 	return user, nil
// }

func (s *userService) Update(ctx context.Context, id string, user *User) (*User, error) {
	updatedUser, err := s.repo.Update(ctx, id, user)
	if err != nil {
		return nil, err
	}
	// Atualiza cache
	cacheKey := redisUserKeyPrefix + id
	data, err := json.Marshal(updatedUser)
	if err == nil {
		redis.Client.Set(ctx, cacheKey, data, 60*time.Second)
	}
	return updatedUser, nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	// Remove do cache
	cacheKey := redisUserKeyPrefix + id
	redis.Client.Del(ctx, cacheKey)
	return nil
}

func fetchAddressFromViaCEP(cep string) (*Address, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("ViaCEP returned status %d", resp.StatusCode)
	}

	var address Address
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		return nil, err
	}

	if address.CEP == "" {
		return nil, fmt.Errorf("CEP inválido ou não encontrado")
	}

	return &address, nil
}
