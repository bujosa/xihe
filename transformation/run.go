package transformation

func RunDealerTransformation() {
	InitDealerTransformation()
	DealerTransformation()
}

func RunCarTransformation() {
	// Before running section dealer run this Make sure you run dealer script before
	Brand()
	BrandToModel("Lexus", " ", "")
	BrandToModel("ANIVERSARIO", "aniversario", "series")
	Model()
	Trim()
	Color()
	DealerIntoCarTransformation()
}
