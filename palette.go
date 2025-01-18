package main

type Variant struct {
	Colors map[string]*Color
}

var (
	MainVariant = Variant{
		Colors: map[string]*Color{
			"base": {
				Hex: "191724",
				RGB: [3]int{25, 23, 36},
				HSL: [3]float64{249, 22, 12},
			},
			"surface": {
				Hex: "1f1d2e",
				RGB: [3]int{31, 29, 46},
				HSL: [3]float64{247, 23, 15},
			},
			"overlay": {
				Hex: "26233a",
				RGB: [3]int{38, 35, 58},
				HSL: [3]float64{245, 25, 18},
			},
			"muted": {
				Hex: "6e6a86",
				RGB: [3]int{110, 106, 134},
				HSL: [3]float64{249, 12, 47},
			},
			"subtle": {
				Hex: "908caa",
				RGB: [3]int{144, 140, 170},
				HSL: [3]float64{248, 15, 61},
			},
			"text": {
				Hex: "e0def4",
				RGB: [3]int{224, 222, 244},
				HSL: [3]float64{245, 50, 91},
			},
			"love": {
				Hex: "eb6f92",
				RGB: [3]int{235, 111, 146},
				HSL: [3]float64{343, 76, 68},
			},
			"gold": {
				Hex: "f6c177",
				RGB: [3]int{246, 193, 119},
				HSL: [3]float64{35, 88, 72},
			},
			"rose": {
				Hex: "ebbcba",
				RGB: [3]int{235, 188, 186},
				HSL: [3]float64{2, 55, 83},
			},
			"pine": {
				Hex: "31748f",
				RGB: [3]int{49, 116, 143},
				HSL: [3]float64{197, 49, 38},
			},
			"foam": {
				Hex: "9ccfd8",
				RGB: [3]int{156, 207, 216},
				HSL: [3]float64{189, 43, 73},
			},
			"iris": {
				Hex: "c4a7e7",
				RGB: [3]int{196, 167, 231},
				HSL: [3]float64{267, 57, 78},
			},
			"highlightLow": {
				Hex: "21202e",
				RGB: [3]int{33, 32, 46},
				HSL: [3]float64{244, 18, 15},
			},
			"highlightMed": {
				Hex: "403d52",
				RGB: [3]int{64, 61, 82},
				HSL: [3]float64{247, 15, 28},
			},
			"highlightHigh": {
				Hex: "524f67",
				RGB: [3]int{82, 79, 103},
				HSL: [3]float64{245, 13, 36},
			},
		},
	}

	MoonVariant = Variant{
		Colors: map[string]*Color{
			"base": {
				Hex: "232136",
				RGB: [3]int{35, 33, 54},
				HSL: [3]float64{246, 24, 17},
			},
			"surface": {
				Hex: "2a273f",
				RGB: [3]int{42, 39, 63},
				HSL: [3]float64{248, 24, 20},
			},
			"overlay": {
				Hex: "393552",
				RGB: [3]int{57, 53, 82},
				HSL: [3]float64{248, 21, 26},
			},
			"muted": {
				Hex: "6e6a86",
				RGB: [3]int{110, 106, 134},
				HSL: [3]float64{249, 12, 47},
			},
			"subtle": {
				Hex: "908caa",
				RGB: [3]int{144, 140, 170},
				HSL: [3]float64{248, 15, 61},
			},
			"text": {
				Hex: "e0def4",
				RGB: [3]int{224, 222, 244},
				HSL: [3]float64{245, 50, 91},
			},
			"love": {
				Hex: "eb6f92",
				RGB: [3]int{235, 111, 146},
				HSL: [3]float64{343, 76, 68},
			},
			"gold": {
				Hex: "f6c177",
				RGB: [3]int{246, 193, 119},
				HSL: [3]float64{35, 88, 72},
			},
			"rose": {
				Hex: "ea9a97",
				RGB: [3]int{234, 154, 151},
				HSL: [3]float64{2, 66, 75},
			},
			"pine": {
				Hex: "3e8fb0",
				RGB: [3]int{62, 143, 176},
				HSL: [3]float64{197, 48, 47},
			},
			"foam": {
				Hex: "9ccfd8",
				RGB: [3]int{156, 207, 216},
				HSL: [3]float64{189, 43, 73},
			},
			"iris": {
				Hex: "c4a7e7",
				RGB: [3]int{196, 167, 231},
				HSL: [3]float64{267, 57, 78},
			},
			"highlightLow": {
				Hex: "2a283e",
				RGB: [3]int{42, 40, 62},
				HSL: [3]float64{245, 22, 20},
			},
			"highlightMed": {
				Hex: "44415a",
				RGB: [3]int{68, 65, 90},
				HSL: [3]float64{247, 16, 30},
			},
			"highlightHigh": {
				Hex: "56526e",
				RGB: [3]int{86, 82, 110},
				HSL: [3]float64{249, 15, 38},
			},
		},
	}

	DawnVariant = Variant{
		Colors: map[string]*Color{
			"base": {
				Hex: "faf4ed",
				RGB: [3]int{250, 244, 237},
				HSL: [3]float64{32, 57, 95},
			},
			"surface": {
				Hex: "fffaf3",
				RGB: [3]int{255, 250, 243},
				HSL: [3]float64{35, 100, 98},
			},
			"overlay": {
				Hex: "f2e9e1",
				RGB: [3]int{242, 233, 225},
				HSL: [3]float64{25, 36, 92},
			},
			"muted": {
				Hex: "9893a5",
				RGB: [3]int{152, 147, 165},
				HSL: [3]float64{254, 9, 61},
			},
			"subtle": {
				Hex: "797593",
				RGB: [3]int{121, 117, 147},
				HSL: [3]float64{249, 13, 52},
			},
			"text": {
				Hex: "575279",
				RGB: [3]int{87, 82, 121},
				HSL: [3]float64{248, 19, 40},
			},
			"love": {
				Hex: "b4637a",
				RGB: [3]int{180, 99, 122},
				HSL: [3]float64{343, 35, 55},
			},
			"gold": {
				Hex: "ea9d34",
				RGB: [3]int{234, 157, 52},
				HSL: [3]float64{35, 81, 56},
			},
			"rose": {
				Hex: "d7827e",
				RGB: [3]int{215, 130, 126},
				HSL: [3]float64{2, 55, 67},
			},
			"pine": {
				Hex: "286983",
				RGB: [3]int{40, 105, 131},
				HSL: [3]float64{197, 53, 34},
			},
			"foam": {
				Hex: "56949f",
				RGB: [3]int{86, 148, 159},
				HSL: [3]float64{189, 30, 48},
			},
			"iris": {
				Hex: "907aa9",
				RGB: [3]int{144, 122, 169},
				HSL: [3]float64{267, 22, 57},
			},
			"highlightLow": {
				Hex: "f4ede8",
				RGB: [3]int{244, 237, 232},
				HSL: [3]float64{25, 35, 93},
			},
			"highlightMed": {
				Hex: "dfdad9",
				RGB: [3]int{223, 218, 217},
				HSL: [3]float64{10, 9, 86},
			},
			"highlightHigh": {
				Hex: "cecacd",
				RGB: [3]int{206, 202, 205},
				HSL: [3]float64{315, 4, 80},
			},
		},
	}
)
