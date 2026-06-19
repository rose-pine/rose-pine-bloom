package color

type Palette map[string]*Color

type VariantMeta struct {
	Id          string
	Name        string
	Appearance  string
	Description string
	Colors      Palette
}

var Accents = []string{
	"love", "gold", "rose", "pine", "foam", "iris",
}

var MainPalette = Palette{
	"base":          {HSL: HSL{249, 22, 12}, RGB: RGB{25, 23, 36}},
	"surface":       {HSL: HSL{247, 23, 15}, RGB: RGB{31, 29, 46}},
	"overlay":       {HSL: HSL{245, 25, 18}, RGB: RGB{38, 35, 58}},
	"muted":         {HSL: HSL{249, 12, 47}, RGB: RGB{110, 106, 134}},
	"subtle":        {HSL: HSL{248, 15, 61}, RGB: RGB{144, 140, 170}},
	"text":          {HSL: HSL{245, 50, 91}, RGB: RGB{224, 222, 244}},
	"love":          {HSL: HSL{343, 76, 68}, RGB: RGB{235, 111, 146}, On: "text"},
	"gold":          {HSL: HSL{35, 88, 72}, RGB: RGB{246, 193, 119}, On: "surface"},
	"rose":          {HSL: HSL{2, 55, 83}, RGB: RGB{235, 188, 186}, On: "surface"},
	"pine":          {HSL: HSL{197, 49, 38}, RGB: RGB{49, 116, 143}, On: "text"},
	"foam":          {HSL: HSL{189, 43, 73}, RGB: RGB{156, 207, 216}, On: "surface"},
	"iris":          {HSL: HSL{267, 57, 78}, RGB: RGB{196, 167, 231}, On: "surface"},
	"highlightLow":  {HSL: HSL{244, 18, 15}, RGB: RGB{33, 32, 46}},
	"highlightMed":  {HSL: HSL{247, 15, 28}, RGB: RGB{64, 61, 82}},
	"highlightHigh": {HSL: HSL{245, 13, 36}, RGB: RGB{82, 79, 103}},
}

var MoonPalette = Palette{
	"base":          {HSL: HSL{246, 24, 17}, RGB: RGB{35, 33, 54}},
	"surface":       {HSL: HSL{248, 24, 20}, RGB: RGB{42, 39, 63}},
	"overlay":       {HSL: HSL{248, 21, 26}, RGB: RGB{57, 53, 82}},
	"muted":         {HSL: HSL{249, 12, 47}, RGB: RGB{110, 106, 134}},
	"subtle":        {HSL: HSL{248, 15, 61}, RGB: RGB{144, 140, 170}},
	"text":          {HSL: HSL{245, 50, 91}, RGB: RGB{224, 222, 244}},
	"love":          {HSL: HSL{343, 76, 68}, RGB: RGB{235, 111, 146}, On: "text"},
	"gold":          {HSL: HSL{35, 88, 72}, RGB: RGB{246, 193, 119}, On: "surface"},
	"rose":          {HSL: HSL{2, 66, 75}, RGB: RGB{234, 154, 151}, On: "surface"},
	"pine":          {HSL: HSL{197, 48, 47}, RGB: RGB{62, 143, 176}, On: "text"},
	"foam":          {HSL: HSL{189, 43, 73}, RGB: RGB{156, 207, 216}, On: "surface"},
	"iris":          {HSL: HSL{267, 57, 78}, RGB: RGB{196, 167, 231}, On: "surface"},
	"highlightLow":  {HSL: HSL{245, 22, 20}, RGB: RGB{42, 40, 62}},
	"highlightMed":  {HSL: HSL{247, 16, 30}, RGB: RGB{68, 65, 90}},
	"highlightHigh": {HSL: HSL{249, 15, 38}, RGB: RGB{86, 82, 110}},
}

var DawnPalette = Palette{
	"base":          {HSL: HSL{32, 57, 95}, RGB: RGB{250, 244, 237}},
	"surface":       {HSL: HSL{35, 100, 98}, RGB: RGB{255, 250, 243}},
	"overlay":       {HSL: HSL{25, 36, 92}, RGB: RGB{242, 233, 225}},
	"muted":         {HSL: HSL{254, 9, 61}, RGB: RGB{152, 147, 165}},
	"subtle":        {HSL: HSL{249, 13, 52}, RGB: RGB{121, 117, 147}},
	"text":          {HSL: HSL{248, 19, 40}, RGB: RGB{87, 82, 121}},
	"love":          {HSL: HSL{343, 35, 55}, RGB: RGB{180, 99, 122}, On: "surface"},
	"gold":          {HSL: HSL{35, 81, 56}, RGB: RGB{234, 157, 52}, On: "surface"},
	"rose":          {HSL: HSL{2, 55, 67}, RGB: RGB{215, 130, 126}, On: "surface"},
	"pine":          {HSL: HSL{197, 53, 34}, RGB: RGB{40, 105, 131}, On: "surface"},
	"foam":          {HSL: HSL{189, 30, 48}, RGB: RGB{86, 148, 159}, On: "surface"},
	"iris":          {HSL: HSL{267, 22, 57}, RGB: RGB{144, 122, 169}, On: "surface"},
	"highlightLow":  {HSL: HSL{25, 35, 93}, RGB: RGB{244, 237, 232}},
	"highlightMed":  {HSL: HSL{10, 9, 86}, RGB: RGB{223, 218, 217}},
	"highlightHigh": {HSL: HSL{315, 4, 80}, RGB: RGB{206, 202, 205}},
}

const description = "All natural pine, faux fur and a bit of soho vibes for the classy minimalist"

var (
	MainVariantMeta = VariantMeta{
		Id:          "rose-pine",
		Name:        "Rosé Pine",
		Appearance:  "dark",
		Description: description,
		Colors:      MainPalette,
	}

	MoonVariantMeta = VariantMeta{
		Id:          "rose-pine-moon",
		Name:        "Rosé Pine Moon",
		Appearance:  "dark",
		Description: description,
		Colors:      MoonPalette,
	}

	DawnVariantMeta = VariantMeta{
		Id:          "rose-pine-dawn",
		Name:        "Rosé Pine Dawn",
		Appearance:  "light",
		Description: description,
		Colors:      DawnPalette,
	}
)

var Variants = []VariantMeta{
	MainVariantMeta,
	MoonVariantMeta,
	DawnVariantMeta,
}
