package models_test

import (
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestFindDutyStations() {
	t := suite.T()

	address := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "12345",
	}
	suite.mustSave(&address)

	station1 := models.DutyStation{
		Name:      "First Station",
		Branch:    internalmessages.MilitaryBranchARMY,
		AddressID: address.ID,
	}
	suite.mustSave(&station1)

	station2 := models.DutyStation{
		Name:      "Second Station",
		Branch:    internalmessages.MilitaryBranchARMY,
		AddressID: address.ID,
	}
	suite.mustSave(&station2)

	stations, err := models.FindDutyStations(suite.db, "first", string(internalmessages.MilitaryBranchARMY))
	if err != nil {
		t.Errorf("Find duty stations error: %v", err)
	}

	if len(stations) != 1 {
		t.Errorf("Should have only got 1 response, got %v", len(stations))
	}

	if stations[0].Name != "First Station" {
		t.Errorf("Station name should have been \"First Station \", got %v", stations[0].Name)
	}

	if stations[0].Address.City != "city" {
		t.Error("Address should have been loaded")
	}
}
