package transforms

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Minh2009/pv_soa/internal/models"
	"github.com/Minh2009/pv_soa/pkgs/utils"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"net"
	"net/http"
	"time"
)

type ProductsReq struct {
	Reference  []string   `json:"references,omitempty"`
	Names      []string   `json:"names,omitempty"`
	AddFrom    *time.Time `json:"add_from,omitempty"`
	AddTo      *time.Time `json:"add_to,omitempty"`
	Status     []int      `json:"status,omitempty"`
	Categories []string   `json:"categories,omitempty"`
	Cities     []string   `json:"cities,omitempty"`

	Offset int `json:"offset,omitempty"`
	Limit  int `json:"limit,omitempty"`
}

func (r ProductsReq) GetOffset() int {
	return r.Offset
}

func (r ProductsReq) GetLimit() int {
	if r.Limit == 0 {
		return 20
	}
	return r.Limit
}

func DecodeProductsReq(_ context.Context, r *http.Request) (interface{}, error) {
	var req ProductsReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

type ProductReq struct {
	Id string `json:"id" validate:"required"`
}

func DecodeProductReq(_ context.Context, r *http.Request) (interface{}, error) {
	uid := mux.Vars(r)["uid"]
	if uid == "" {
		return nil, errors.New("invalid uid")
	}
	return ProductReq{Id: uid}, nil
}

type ProductCreateReq struct {
	Name     string           `json:"name" validate:"required"`
	Price    *decimal.Decimal `json:"price" validate:"required"`
	Quantity *int64           `json:"quantity"`
	Status   int              `json:"status,omitempty"`

	Categories []string `json:"categories"`
	CityId     string   `json:"city_id" validate:"required"`
	SupplierId string   `json:"supplier_id" validate:"required"`
}

func DecodeProductCreateReq(_ context.Context, r *http.Request) (interface{}, error) {
	var req ProductCreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	err := utils.Validate(req)
	if err != nil {
		return nil, err
	}
	if len(req.Categories) == 0 {
		return nil, errors.New("product must be belong to a category")
	}
	valid := []models.ProductStatus{models.OnOrder, models.Available, models.OutOfStock}
	if !utils.Contains(valid, models.ProductStatus(req.Status)) {
		req.Status = int(models.Available)
	}
	return req, nil
}

type ProductUpdateReq struct {
	ProductCreateReq
	ProductReq
}

func DecodeProductUpdateReq(_ context.Context, r *http.Request) (interface{}, error) {
	var req ProductUpdateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	err := utils.Validate(req.ProductReq)
	if err != nil {
		return nil, err
	}
	if len(req.Categories) == 0 {
		return nil, errors.New("product must be belong to a category")
	}
	valid := []models.ProductStatus{models.OnOrder, models.Available, models.OutOfStock}
	if !utils.Contains(valid, models.ProductStatus(req.Status)) {
		req.Status = 0
	}
	return req, nil
}

type ProductDistanceReq struct {
	ProductReq
	IP string `json:"ip" validate:"required"`
}

func DecodeProductDistanceReq(_ context.Context, r *http.Request) (interface{}, error) {
	var req ProductDistanceReq
	req.Id = mux.Vars(r)["uid"]
	// Get client IP
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	if testIP := r.URL.Query().Get("ip"); testIP != "" {
		ip = testIP
	}
	ipAddr, _, err := net.SplitHostPort(ip)
	if err != nil {
		// If error, the IP might not have port information
		ipAddr = ip
	}
	req.IP = ipAddr
	err = utils.Validate(req.ProductReq)
	if err != nil {
		return nil, err
	}
	return req, nil
}
