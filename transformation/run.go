package transformation

import "context"

func RunDealerTransformation(ctx context.Context) {
	InitDealerTransformation(ctx)
	DealerTransformation(ctx)
}

func RunCarTransformation(ctx context.Context) {
	// Before running section dealer run this Make sure you run dealer script before
	// Brand()
	// BrandToModel("Lexus", " ", "")
	// BrandToModel("ANIVERSARIO", "aniversario", "series")
	// Model()
	// Trim()
	// Color()
	DealerIntoCarTransformation(ctx)
}
