package transformation

import (
	"context"

	"github.com/bujosa/xihe/scripts"
)

func RunDealerTransformation(ctx context.Context) {
	InitDealerTransformation(ctx)
	DealerTransformation(ctx)
}

func RunCarTransformation(ctx context.Context) {
	// Before running section dealer run this Make sure you run dealer script before

	// First step is about lookup brand and add some fields necesary for the transformation
	Brand(ctx)

	// Model transformation
	BrandToModel(ctx, "ANIVERSARIO", "aniversario", "series")
	BrandToModel(ctx, "Runner", "4-runner", "4runner")
	Model(ctx, 1)

	// Color transformation
	Color(ctx)

	// Fueltype transformation
	Fueltype(ctx)

	// Transmission transformation
	Transmission(ctx)

	// BodyStyle transformation
	BodyStyle(ctx)

	// DriveTrain transformation
	DriveTrain(ctx)

	// Model UnMatched Strategy
	UnMatchedModelLayerTwo(ctx)
	scripts.FixTrimNameForModelMatchLayer(ctx, 2)

	// Model UnMatched Strategy Layer 3
	UnMatchedModelLayerThree(ctx)
	scripts.FixTrimNameForModelMatchLayer(ctx, 3)

	// Trim transformation
	Trim(ctx)

	// Dealer into car transformation
	DealerIntoCarTransformation(ctx)
}
