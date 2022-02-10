package handler

import (
	"context"
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
)

func checkDuplicateHub(s *Server, name, id string) validation.RuleFunc {
	return func(value interface{}) error {
		resp, _ := s.st.GetHubBy(context.Background(), trim(name))
		if resp == nil || resp.ID == id {
			return nil
		}
		return errors.New(name + AlrEx)
	}
}

func checkDuplicateBranch(s *Server, name, id string) validation.RuleFunc {
	return func(value interface{}) error {
		resp, _ := s.st.GetBranchBy(context.Background(), trim(name))
		if resp == nil || resp.ID == id {
			return nil
		}
		return errors.New(name + AlrEx)
	}
}

func checkHubPosition(s *Server, position int32, id string) validation.RuleFunc {
	return func(value interface{}) error {
		resp, _ := s.st.GetHubByPosition(context.Background(), position)
		if resp == nil || resp.ID == id {
			return nil
		}
		return fmt.Errorf(PosEx, position)
	}
}

func checkDuplicateDistrict(s *Server, name, id string) validation.RuleFunc {
	return func(value interface{}) error {
		resp, _ := s.st.GetDistrictBy(context.Background(), trim(name))
		if resp == nil || resp.ID == id {
			return nil
		}
		return errors.New(name + AlrEx)
	}
}

func checkDistrictPosition(s *Server, position int32, id string) validation.RuleFunc {
	return func(value interface{}) error {
		resp, _ := s.st.GetDistrictByPosition(context.Background(), position)
		if resp == nil || resp.ID == id {
			return nil
		}
		return fmt.Errorf(PosEx, position)
	}
}

func checkBranchPosition(s *Server, position int32, id string) validation.RuleFunc {
	return func(value interface{}) error {
		resp, _ := s.st.GetBranchByPosition(context.Background(), position)
		if resp == nil || resp.ID == id {
			return nil
		}
		return fmt.Errorf(PosEx, position)
	}
}

func checkDeliveryCompanyPosition(s *Server, position int32, id string) validation.RuleFunc {
	return func(value interface{}) error {
		resp, _ := s.st.GetDeliveryCompanyByPosition(context.Background(), position)
		if resp == nil || resp.ID == id {
			return nil
		}
		return fmt.Errorf(PosEx, position)
	}
}

func checkCountryExists(s *Server, cid string) validation.RuleFunc {
	return func(value interface{}) error {
		str := fmt.Errorf(" Country not exists")
		res, _ := s.st.GetCountryBy(context.Background(), cid)
		if res == nil {
			return str
		}
		if res != nil || res.ID == cid {
			return nil
		}
		return str

	}
}

func checkDuplicateHubPhone(s *Server, phone string, id string) validation.RuleFunc {
	return func(value interface{}) error {
		resp, _ := s.st.GetHubBy(context.Background(), phone)
		if resp == nil || resp.ID == id {
			return nil
		}
		return fmt.Errorf(PhnEx, phone)
	}
}

func checkDuplicateDeliveryCompanyPhone(s *Server, phone string, id string) validation.RuleFunc {
	return func(value interface{}) error {
		resp, _ := s.st.GetDeliveryCompanyBy(context.Background(), phone)
		if resp == nil || resp.ID == id {
			return nil
		}
		return fmt.Errorf(PhnEx, phone)
	}
}

func checkDuplicateBranchPhone(s *Server, phone string, id string) validation.RuleFunc {
	return func(value interface{}) error {
		resp, _ := s.st.GetBranchBy(context.Background(), phone)
		if resp == nil || resp.ID == id {
			return nil
		}
		return fmt.Errorf(PhnEx, phone)
	}
}

func checkDuplicateBranchEmail(s *Server, email string, id string) validation.RuleFunc {
	return func(value interface{}) error {
		resp, _ := s.st.GetBranchBy(context.Background(), email)
		if resp == nil || resp.ID == id {
			return nil
		}
		return fmt.Errorf(PhnEx, email)
	}
}

func checkDuplicateDeliveryCompanyEmail(s *Server, email string, id string) validation.RuleFunc {
	return func(value interface{}) error {
		resp, _ := s.st.GetDeliveryCompanyBy(context.Background(), email)
		if resp == nil || resp.ID == id {
			return nil
		}
		return fmt.Errorf(PhnEx, email)
	}
}

func checkDuplicateDeliveryCompany(s *Server, email string, id string) validation.RuleFunc {
	return func(value interface{}) error {
		resp, _ := s.st.GetHubBy(context.Background(), email)
		if resp == nil || resp.ID == id {
			return nil
		}
		return fmt.Errorf(PhnEx, email)
	}
}

func checkDuplicateHubEmail(s *Server, email string, id string) validation.RuleFunc {
	return func(value interface{}) error {
		resp, _ := s.st.GetHubBy(context.Background(), email)
		if resp == nil || resp.ID == id {
			return nil
		}
		return fmt.Errorf(PhnEx, email)
	}
}

func checkDuplicateStation(s *Server, name, id string) validation.RuleFunc {
	return func(value interface{}) error {
		resp, _ := s.st.GetStationBy(context.Background(), trim(name))
		if resp == nil || resp.ID == id {
			return nil
		}
		return errors.New(name + AlrEx)
	}
}

func checkStationPosition(s *Server, position int32, id string) validation.RuleFunc {
	return func(value interface{}) error {
		resp, _ := s.st.GetStationByPosition(context.Background(), position)
		if resp == nil || resp.ID == id {
			return nil
		}
		return fmt.Errorf(PosEx, position)
	}
}

func checkDistrictExists(s *Server, id string) validation.RuleFunc {
	return func(value interface{}) error {
		str := fmt.Errorf(" district not exists")
		res, _ := s.st.GetDistrictBy(context.Background(), id)
		if res == nil {
			return str
		}
		if res != nil || res.ID == id {
			return nil
		}
		return str
	}
}

func checkDuplicateDesignation(s *Server, name, id string) validation.RuleFunc {
	return func(value interface{}) error {
		resp, _ := s.st.GetDesignationBy(context.Background(), trim(name))
		if resp == nil || resp.ID == id {
			return nil
		}
		return errors.New(name + AlrEx)
	}
}

func checkDesignationPosition(s *Server, position int32, id string) validation.RuleFunc {
	return func(value interface{}) error {
		resp, _ := s.st.GetDesignationByPosition(context.Background(), position)
		if resp == nil || resp.ID == id {
			return nil
		}
		return fmt.Errorf(PosEx, position)
	}
}
func checkDuplicateCountry(s *Server, name, id string) validation.RuleFunc {
	return func(value interface{}) error {
		resp, _ := s.st.GetCountryBy(context.Background(), trim(name))
		if resp == nil || resp.ID == id {
			return nil
		}
		return errors.New(name + AlrEx)
	}
}

func checkCountryPosition(s *Server, position int32, id string) validation.RuleFunc {
	return func(value interface{}) error {
		resp, _ := s.st.GetCountryByPosition(context.Background(), position)
		if resp == nil || resp.ID == id {
			return nil
		}
		return fmt.Errorf(PosEx, position)
	}
}
