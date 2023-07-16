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
	BrandToModel(ctx, "Lexus", " ", "")
	BrandToModel(ctx, "ANIVERSARIO", "aniversario", "series")
	Model(ctx)

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
	scripts.FixTrimNameForModelMatchLayerTwo(ctx)

	// Second if lookup
	Trim(ctx)
	Color(ctx)
	DealerIntoCarTransformation(ctx)
}
