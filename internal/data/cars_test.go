package data

import (
	"gearboxd/internal/assert"
	"gearboxd/internal/validator"
	"testing"

	"github.com/shopspring/decimal"
)

func TestValidateCar(t *testing.T) {
	tests := []struct {
		name           string
		Make           string
		Model          string
		Year           int
		Description    string
		ImageURL       string
		Gearbox        string
		Drivetrain     string
		Horsepower     int
		Fuel           string
		PriceNew       float64
		expectedErrors map[string]string
	}{
		{
			name: "valid car",
			Make: "BMW", Model: "M3", Year: 2020,
			Description: "Sport sedan",
			ImageURL:    "https://img.com/m3.jpg",
			Gearbox:     "automatic", Drivetrain: "RWD",
			Horsepower: 480, Fuel: "gas", PriceNew: 50000,
			expectedErrors: map[string]string{},
		},
		{
			name: "missing make",
			Make: "", Model: "M3", Year: 2020,
			Description: "Sport sedan",
			ImageURL:    "https://img.com/m3.jpg",
			Gearbox:     "automatic", Drivetrain: "RWD",
			Horsepower: 480, Fuel: "gas", PriceNew: 50000,
			expectedErrors: map[string]string{
				"make": "make cannot be empty",
			},
		},
		{
			name: "year below minimum",
			Make: "BMW", Model: "M3", Year: 1800,
			Description: "Sport sedan",
			ImageURL:    "https://img.com/m3.jpg",
			Gearbox:     "automatic", Drivetrain: "RWD",
			Horsepower: 480, Fuel: "gas", PriceNew: 50000,
			expectedErrors: map[string]string{
				"year": "the minimum year is 1886, when the first car was patented",
			},
		},
		{
			name: "invalid gearbox",
			Make: "BMW", Model: "M3", Year: 2020,
			Description: "Sport sedan",
			ImageURL:    "https://img.com/m3.jpg",
			Gearbox:     "semi-auto", Drivetrain: "RWD",
			Horsepower: 480, Fuel: "gas", PriceNew: 50000,
			expectedErrors: map[string]string{
				"gearbox": "gearbox must be one of: [manual automatic DCT CVT]",
			},
		},
		{
			name: "invalid drivetrain",
			Make: "BMW", Model: "M3", Year: 2020,
			Description: "Sport sedan",
			ImageURL:    "https://img.com/m3.jpg",
			Gearbox:     "automatic", Drivetrain: "awd",
			Horsepower: 480, Fuel: "gas", PriceNew: 50000,
			expectedErrors: map[string]string{
				"drivetrain": "drivetrain must be one of: [FWD RWD AWD 4WD]",
			},
		},
		{
			name: "invalid fuel",
			Make: "BMW", Model: "M3", Year: 2020,
			Description: "Sport sedan",
			ImageURL:    "https://img.com/m3.jpg",
			Gearbox:     "automatic", Drivetrain: "RWD",
			Horsepower: 480, Fuel: "petrol", PriceNew: 50000,
			expectedErrors: map[string]string{
				"fuel": "fuel must be one of: [diesel gas electric hybrid plug-in-hybrid hydrogen lpg cng]",
			},
		},
		{
			name: "price zero",
			Make: "BMW", Model: "M3", Year: 2020,
			Description: "Sport sedan",
			ImageURL:    "https://img.com/m3.jpg",
			Gearbox:     "automatic", Drivetrain: "RWD",
			Horsepower: 480, Fuel: "gas", PriceNew: 0,
			expectedErrors: map[string]string{
				"price_new": "price_new must be a positive number higher than 0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			car := &Car{
				Make:        tt.Make,
				Model:       tt.Model,
				Year:        tt.Year,
				Description: tt.Description,
				ImageURL:    tt.ImageURL,
				Gearbox:     tt.Gearbox,
				Drivetrain:  tt.Drivetrain,
				Horsepower:  tt.Horsepower,
				Fuel:        tt.Fuel,
				PriceNew:    decimal.NewFromFloat(tt.PriceNew),
			}

			v := validator.New()

			ValidateCar(v, car)

			assert.Equal(t, v.Errors, tt.expectedErrors)
		})
	}
}
