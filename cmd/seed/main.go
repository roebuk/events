package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"firecrest/db"
)

type raceSpec struct {
	Name        string
	Slug        string
	MaxCapacity int32
	PriceUnits  int32 // pence
}

type eventSpec struct {
	Name  string
	Slug  string
	Year  int32
	Races []raceSpec
}

type orgSpec struct {
	Name   string
	Events []eventSpec
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	_ = godotenv.Load()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_NAME", "firecrest"),
		getEnv("DB_SSLMODE", "disable"),
	)

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()

	queries := db.New(pool)

	orgs := seedData()

	for _, o := range orgs {
		org, err := queries.CreateOrganisation(ctx, o.Name)
		if err != nil {
			return fmt.Errorf("create organisation %q: %w", o.Name, err)
		}
		fmt.Printf("Created organisation: %s (id=%d)\n", org.Name, org.ID)

		for _, e := range o.Events {
			event, err := queries.CreateEvent(ctx, db.CreateEventParams{
				OrganisationID: org.ID,
				Name:           e.Name,
				Slug:           e.Slug,
				Year:           e.Year,
			})
			if err != nil {
				return fmt.Errorf("create event %q: %w", e.Name, err)
			}
			fmt.Printf("  Event: %s (id=%d, year=%d)\n", event.Name, event.ID, event.Year)

			regOpen := time.Date(int(e.Year), 1, 15, 9, 0, 0, 0, time.UTC)
			regClose := time.Date(int(e.Year), 6, 1, 23, 59, 0, 0, time.UTC)

			for _, r := range e.Races {
				race, err := queries.CreateRace(ctx, db.CreateRaceParams{
					EventID: event.ID,
					Name:    r.Name,
					Slug:    r.Slug,
					RegistrationOpenDate: pgtype.Timestamptz{
						Time:  regOpen,
						Valid: true,
					},
					RegistrationCloseDate: pgtype.Timestamptz{
						Time:  regClose,
						Valid: true,
					},
					MaxCapacity: r.MaxCapacity,
					PriceUnits: pgtype.Int4{
						Int32: r.PriceUnits,
						Valid: true,
					},
					Currency: pgtype.Text{
						String: "GBP",
						Valid:  true,
					},
				})
				if err != nil {
					return fmt.Errorf("create race %q for event %q: %w", r.Name, e.Name, err)
				}
				fmt.Printf("    Race: %s (id=%d, capacity=%d, price=£%.2f)\n",
					race.Name, race.ID, race.MaxCapacity, float64(r.PriceUnits)/100)
			}
		}
	}

	fmt.Println("\nSeeding complete.")
	return nil
}

func seedData() []orgSpec {
	return []orgSpec{
		{
			Name: "Northern Trail Events",
			Events: []eventSpec{
				{
					Name: "Pennine Way Ultra", Slug: "pennine-way-ultra", Year: 2026,
					Races: []raceSpec{
						{Name: "50 Mile Ultra", Slug: "50-mile", MaxCapacity: 200, PriceUnits: 8500},
						{Name: "30 Mile Challenge", Slug: "30-mile", MaxCapacity: 300, PriceUnits: 6500},
					},
				},
				{
					Name: "Yorkshire Three Peaks Fell Race", Slug: "yorkshire-three-peaks", Year: 2026,
					Races: []raceSpec{
						{Name: "Three Peaks", Slug: "three-peaks", MaxCapacity: 500, PriceUnits: 3500},
					},
				},
				{
					Name: "Lakeland 100", Slug: "lakeland-100", Year: 2026,
					Races: []raceSpec{
						{Name: "100 Mile", Slug: "100-mile", MaxCapacity: 250, PriceUnits: 15000},
						{Name: "50 Mile", Slug: "50-mile", MaxCapacity: 300, PriceUnits: 9500},
					},
				},
				{
					Name: "Kielder Dark Skies Trail Marathon", Slug: "kielder-dark-skies", Year: 2026,
					Races: []raceSpec{
						{Name: "Marathon", Slug: "marathon", MaxCapacity: 400, PriceUnits: 4500},
						{Name: "Half Marathon", Slug: "half-marathon", MaxCapacity: 600, PriceUnits: 3200},
					},
				},
				{
					Name: "Hadrian's Wall Ultra", Slug: "hadrians-wall-ultra", Year: 2026,
					Races: []raceSpec{
						{Name: "69 Mile End-to-End", Slug: "69-mile", MaxCapacity: 150, PriceUnits: 11000},
					},
				},
			},
		},
		{
			Name: "South West Running Co",
			Events: []eventSpec{
				{
					Name: "Bath Half Marathon", Slug: "bath-half-marathon", Year: 2026,
					Races: []raceSpec{
						{Name: "Half Marathon", Slug: "half-marathon", MaxCapacity: 12000, PriceUnits: 3900},
					},
				},
				{
					Name: "Cotswolds Classic 10K", Slug: "cotswolds-classic-10k", Year: 2026,
					Races: []raceSpec{
						{Name: "10K", Slug: "10k", MaxCapacity: 1500, PriceUnits: 2800},
					},
				},
				{
					Name: "Exmoor Coastal Trail", Slug: "exmoor-coastal-trail", Year: 2026,
					Races: []raceSpec{
						{Name: "Ultra 35 Mile", Slug: "ultra-35", MaxCapacity: 200, PriceUnits: 5500},
						{Name: "Half Marathon", Slug: "half-marathon", MaxCapacity: 400, PriceUnits: 3500},
						{Name: "10K", Slug: "10k", MaxCapacity: 600, PriceUnits: 2500},
					},
				},
				{
					Name: "Bristol 10K", Slug: "bristol-10k", Year: 2026,
					Races: []raceSpec{
						{Name: "10K", Slug: "10k", MaxCapacity: 8000, PriceUnits: 2600},
					},
				},
				{
					Name: "Jurassic Coast Marathon", Slug: "jurassic-coast-marathon", Year: 2026,
					Races: []raceSpec{
						{Name: "Marathon", Slug: "marathon", MaxCapacity: 800, PriceUnits: 4800},
						{Name: "Half Marathon", Slug: "half-marathon", MaxCapacity: 1200, PriceUnits: 3400},
					},
				},
			},
		},
		{
			Name: "Peak District Multisport",
			Events: []eventSpec{
				{
					Name: "Dark Peak Fell Race", Slug: "dark-peak-fell-race", Year: 2026,
					Races: []raceSpec{
						{Name: "Long Route (18 miles)", Slug: "long", MaxCapacity: 300, PriceUnits: 2800},
						{Name: "Short Route (9 miles)", Slug: "short", MaxCapacity: 400, PriceUnits: 2000},
					},
				},
				{
					Name: "Tour of the Peaks Sportive", Slug: "tour-of-the-peaks", Year: 2026,
					Races: []raceSpec{
						{Name: "Epic (100 miles)", Slug: "epic", MaxCapacity: 500, PriceUnits: 4500},
						{Name: "Standard (70 miles)", Slug: "standard", MaxCapacity: 800, PriceUnits: 3800},
						{Name: "Short (40 miles)", Slug: "short", MaxCapacity: 1000, PriceUnits: 3000},
					},
				},
				{
					Name: "Chatsworth Sportive", Slug: "chatsworth-sportive", Year: 2026,
					Races: []raceSpec{
						{Name: "Gran Fondo (85 miles)", Slug: "gran-fondo", MaxCapacity: 600, PriceUnits: 4200},
						{Name: "Medio Fondo (55 miles)", Slug: "medio-fondo", MaxCapacity: 800, PriceUnits: 3500},
					},
				},
				{
					Name: "Peak District Gravel 100", Slug: "peak-gravel-100", Year: 2026,
					Races: []raceSpec{
						{Name: "100K Gravel", Slug: "100k", MaxCapacity: 300, PriceUnits: 5000},
						{Name: "50K Gravel", Slug: "50k", MaxCapacity: 400, PriceUnits: 3500},
					},
				},
				{
					Name: "Derwent Valley Triathlon", Slug: "derwent-valley-tri", Year: 2026,
					Races: []raceSpec{
						{Name: "Olympic Distance", Slug: "olympic", MaxCapacity: 400, PriceUnits: 6500},
						{Name: "Sprint Distance", Slug: "sprint", MaxCapacity: 600, PriceUnits: 4500},
					},
				},
			},
		},
		{
			Name: "Essex Endurance Events",
			Events: []eventSpec{
				{
					Name: "Southend Half Marathon", Slug: "southend-half-marathon", Year: 2026,
					Races: []raceSpec{
						{Name: "Half Marathon", Slug: "half-marathon", MaxCapacity: 3000, PriceUnits: 3200},
					},
				},
				{
					Name: "Essex Gravel Century", Slug: "essex-gravel-century", Year: 2026,
					Races: []raceSpec{
						{Name: "100 Mile Gravel", Slug: "100-mile", MaxCapacity: 250, PriceUnits: 5500},
						{Name: "60 Mile Gravel", Slug: "60-mile", MaxCapacity: 400, PriceUnits: 4000},
						{Name: "30 Mile Intro", Slug: "30-mile", MaxCapacity: 500, PriceUnits: 2800},
					},
				},
				{
					Name: "Chelmsford Marathon", Slug: "chelmsford-marathon", Year: 2026,
					Races: []raceSpec{
						{Name: "Marathon", Slug: "marathon", MaxCapacity: 2000, PriceUnits: 4200},
						{Name: "Half Marathon", Slug: "half-marathon", MaxCapacity: 3000, PriceUnits: 3000},
						{Name: "10K Fun Run", Slug: "10k", MaxCapacity: 5000, PriceUnits: 2200},
					},
				},
				{
					Name: "Lee Valley Triathlon", Slug: "lee-valley-tri", Year: 2026,
					Races: []raceSpec{
						{Name: "Middle Distance", Slug: "middle", MaxCapacity: 300, PriceUnits: 9500},
						{Name: "Olympic Distance", Slug: "olympic", MaxCapacity: 500, PriceUnits: 6500},
						{Name: "Sprint Distance", Slug: "sprint", MaxCapacity: 700, PriceUnits: 4500},
					},
				},
				{
					Name: "Epping Forest Trail 10K", Slug: "epping-forest-10k", Year: 2026,
					Races: []raceSpec{
						{Name: "10K Trail", Slug: "10k", MaxCapacity: 800, PriceUnits: 2400},
					},
				},
			},
		},
		{
			Name: "Welsh Mountain Events",
			Events: []eventSpec{
				{
					Name: "Snowdonia Marathon Eryri", Slug: "snowdonia-marathon-eryri", Year: 2026,
					Races: []raceSpec{
						{Name: "Marathon", Slug: "marathon", MaxCapacity: 2500, PriceUnits: 5800},
					},
				},
				{
					Name: "Brecon Beacons Ultra", Slug: "brecon-beacons-ultra", Year: 2026,
					Races: []raceSpec{
						{Name: "Ultra 46 Mile", Slug: "ultra-46", MaxCapacity: 250, PriceUnits: 7500},
						{Name: "Marathon", Slug: "marathon", MaxCapacity: 400, PriceUnits: 4800},
					},
				},
				{
					Name: "Dragon Ride Sportive", Slug: "dragon-ride-sportive", Year: 2026,
					Races: []raceSpec{
						{Name: "Gran Fondo (140 miles)", Slug: "gran-fondo", MaxCapacity: 400, PriceUnits: 6500},
						{Name: "Medio Fondo (100 miles)", Slug: "medio-fondo", MaxCapacity: 600, PriceUnits: 5000},
						{Name: "Piccolo Fondo (68 miles)", Slug: "piccolo-fondo", MaxCapacity: 800, PriceUnits: 3800},
					},
				},
				{
					Name: "Pembrokeshire Coastal Triathlon", Slug: "pembrokeshire-tri", Year: 2026,
					Races: []raceSpec{
						{Name: "Olympic Distance", Slug: "olympic", MaxCapacity: 350, PriceUnits: 7000},
						{Name: "Sprint Distance", Slug: "sprint", MaxCapacity: 500, PriceUnits: 4800},
					},
				},
				{
					Name: "Cader Idris Fell Race", Slug: "cader-idris-fell-race", Year: 2026,
					Races: []raceSpec{
						{Name: "Fell Race (10 miles)", Slug: "fell-race", MaxCapacity: 200, PriceUnits: 2200},
					},
				},
			},
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
