package price

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type PriceService interface {
	UserDiscountSet(uid uint32, id string, startTime, endTime string) error
	UserRebateSet(uid uint32, id string, startTime, endTime string) error
	UserPackageSet(uid uint32, id string, startTime, endTime string) error
}

type priceService struct {
	host   string
	client *http.Client
}

func NewPriceService(host string, tr http.RoundTripper) PriceService {
	return &priceService{
		host:   host,
		client: &http.Client{Transport: tr},
	}
}

func (s *priceService) UserDiscountSet(uid uint32, id string, startTime, endTime string) error {
	resp, err := s.client.PostForm(s.host+"/v3/user/discount/set", url.Values{
		"uid":         {strconv.FormatUint(uint64(uid), 10)},
		"id":          {id},
		"effect_time": {startTime},
		"dead_time":   {endTime},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("response code: %d", resp.StatusCode)
	}
	return nil
}

func (s *priceService) UserRebateSet(uid uint32, id string, startTime, endTime string) error {
	resp, err := s.client.PostForm(s.host+"/v3/user/rebate/set", url.Values{
		"uid":         {strconv.FormatUint(uint64(uid), 10)},
		"id":          {id},
		"effect_time": {startTime},
		"dead_time":   {endTime},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("response code: %d", resp.StatusCode)
	}
	return nil
}

func (s *priceService) UserPackageSet(uid uint32, id string, startTime, endTime string) error {
	resp, err := s.client.PostForm(s.host+"/v3/user/package/set", url.Values{
		"uid":         {strconv.FormatUint(uint64(uid), 10)},
		"id":          {id},
		"effect_time": {startTime},
		"dead_time":   {endTime},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("response code: %d", resp.StatusCode)
	}
	return nil
}
