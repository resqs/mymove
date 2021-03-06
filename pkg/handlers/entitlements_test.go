package handlers

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"
	entitlementop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/entitlements"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestValidateEntitlementHandlerReturns200() {
	// Given: a set of orders, a move, user, servicemember and a PPM
	ppm, _ := testdatagen.MakePPM(suite.db)
	move := ppm.Move

	// When: rank is E1, the orders have dependents and spouse gear, and
	// the weight estimate stored is under entitlement of 10500
	ppm.WeightEstimate = swag.Int64(10000)
	suite.mustSave(&ppm)

	// And: the context contains the auth values
	request := httptest.NewRequest("GET", "/entitlements/move_id", nil)
	request = suite.authenticateRequest(request, move.Orders.ServiceMember)

	params := entitlementop.ValidateEntitlementParams{
		HTTPRequest: request,
		MoveID:      strfmt.UUID(move.ID.String()),
	}

	// And: validate entitlements endpoint is hit
	handler := ValidateEntitlementHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&entitlementop.ValidateEntitlementOK{}, response)

}

func (suite *HandlerSuite) TestValidateEntitlementHandlerReturns409() {
	// Given: a set of orders, a move, user, servicemember and a PPM
	ppm, _ := testdatagen.MakePPM(suite.db)
	move := ppm.Move

	// When: rank is E1, the orders have dependents and spouse gear, and
	// the weight estimate stored is over entitlement of 10500
	ppm.WeightEstimate = swag.Int64(14000)
	suite.mustSave(&ppm)

	// And: the context contains the auth values
	request := httptest.NewRequest("GET", "/entitlements/move_id", nil)
	request = suite.authenticateRequest(request, move.Orders.ServiceMember)

	params := entitlementop.ValidateEntitlementParams{
		HTTPRequest: request,
		MoveID:      strfmt.UUID(move.ID.String()),
	}

	// And: validate entitlements endpoint is hit
	handler := ValidateEntitlementHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	// Then: expect a 409 status code
	suite.Assertions.IsType(&errResponse{}, response)
	errResponse := response.(*errResponse)

	// Then: expect a 409 status code
	suite.Assertions.Equal(http.StatusConflict, errResponse.code)
}

func (suite *HandlerSuite) TestValidateEntitlementHandlerReturns404IfNoPpm() {
	// Given: a set of orders, a move, user, servicemember but NO ppm
	move, _ := testdatagen.MakeMove(suite.db)

	// When: rank is E1, the orders have dependents and spouse gear
	// And: the context contains the auth values
	request := httptest.NewRequest("GET", "/entitlements/move_id", nil)
	request = suite.authenticateRequest(request, move.Orders.ServiceMember)

	params := entitlementop.ValidateEntitlementParams{
		HTTPRequest: request,
		MoveID:      strfmt.UUID(move.ID.String()),
	}

	// And: validate entitlements endpoint is hit
	handler := ValidateEntitlementHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	// Then: expect a 404 status code
	suite.Assertions.IsType(&entitlementop.ValidateEntitlementNotFound{}, response)
}

func (suite *HandlerSuite) TestValidateEntitlementHandlerReturns404IfNoMoveOrOrders() {
	// Given: a user, servicemember but NO Move
	serviceMember, _ := testdatagen.MakeServiceMember(suite.db)

	// When: rank is E1, the orders have dependents and spouse gear
	// And: the context contains the auth values
	request := httptest.NewRequest("GET", "/entitlements/move_id", nil)
	request = suite.authenticateRequest(request, serviceMember)

	badMoveID := uuid.Must(uuid.NewV4())

	params := entitlementop.ValidateEntitlementParams{
		HTTPRequest: request,
		MoveID:      strfmt.UUID(badMoveID.String()),
	}

	// And: validate entitlements endpoint is hit
	handler := ValidateEntitlementHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	// Then: expect a 404 status code
	suite.Assertions.IsType(&errResponse{}, response)
	errResponse := response.(*errResponse)

	// Then: expect a 404 status code
	suite.Assertions.Equal(http.StatusNotFound, errResponse.code)
}

func (suite *HandlerSuite) TestValidateEntitlementHandlerReturns404IfNoRank() {
	// Given: a set of orders, a move, user, servicemember and a PPM
	ppm, _ := testdatagen.MakePPM(suite.db)
	move := ppm.Move

	// When: rank is E1, the orders have dependents and spouse gear, and
	// the weight estimate stored is under entitlement of 10500
	ppm.WeightEstimate = swag.Int64(10000)
	suite.mustSave(&ppm)

	move.Orders.ServiceMember.Rank = nil
	suite.mustSave(&move.Orders.ServiceMember)

	// And: the context contains the auth values
	request := httptest.NewRequest("GET", "/entitlements/move_id", nil)
	request = suite.authenticateRequest(request, move.Orders.ServiceMember)

	params := entitlementop.ValidateEntitlementParams{
		HTTPRequest: request,
		MoveID:      strfmt.UUID(move.ID.String()),
	}

	// And: validate entitlements endpoint is hit
	handler := ValidateEntitlementHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	// Then: expect a 404 status code
	suite.Assertions.IsType(&entitlementop.ValidateEntitlementNotFound{}, response)
}

func (suite *HandlerSuite) TestGetEntitlementWithValidValues() {
	E1 := internalmessages.ServiceMemberRankE1

	// When: E1 has dependents and spouse gear
	suite.Assertions.Equal(10500, getEntitlement(E1, true, true))
	// When: E1 doesn't have dependents or spouse gear
	suite.Assertions.Equal(7000, getEntitlement(E1, false, false))
	// When: E1 doesn't have dependents but has spouse gear - impossible state
	suite.Assertions.Equal(7000, getEntitlement(E1, false, true))
	// When: E1 has dependents but no spouse gear
	suite.Assertions.Equal(10000, getEntitlement(E1, true, false))

}
