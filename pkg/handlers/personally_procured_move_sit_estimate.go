package handlers

import (
	"time"

	"github.com/go-openapi/runtime/middleware"

	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/unit"
)

// ShowPPMSitEstimateHandler returns PPM SIT estimate for a weight, move date,
type ShowPPMSitEstimateHandler HandlerContext

// Handle calculates SIT charge and retrieves SIT discount rate.
// It returns the discount rate applied to relevant SIT charge.
func (h ShowPPMSitEstimateHandler) Handle(params ppmop.ShowPPMSitEstimateParams) middleware.Responder {
	engine := rateengine.NewRateEngine(h.db, h.logger, h.planner)
	sitZip3 := rateengine.Zip5ToZip3(params.DestinationZip)
	cwtWeight := unit.Pound(params.WeightEstimate).ToCWT()
	plannedMoveDateTime := time.Time(params.PlannedMoveDate)

	_, sitDiscount, err := PPMDiscountFetch(h.db,
		h.logger,
		params.OriginZip,
		params.DestinationZip,
		time.Time(params.PlannedMoveDate),
	)
	if err != nil {
		return responseForError(h.logger, err)
	}
	sitTotal, err := engine.SitCharge(cwtWeight, int(params.DaysInStorage), sitZip3, plannedMoveDateTime, true)

	if err != nil {
		return responseForError(h.logger, err)
	}

	// Swagger returns int64 when using the integer type
	sitCharge := int64(sitDiscount.Apply(sitTotal))

	ppmSitEstimate := internalmessages.PPMSitEstimate{
		Estimate: &sitCharge,
	}
	return ppmop.NewShowPPMSitEstimateOK().WithPayload(&ppmSitEstimate)
}
