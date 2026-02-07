package viewmodels

import "time"

// EventViewModel represents an event for display purposes
type EventViewModel struct {
	Slug        string
	Name        string
	Date        time.Time
	Location    string
	ImageURL    string
	RaceType    string    // e.g., "Trail Run", "Road Race", "Ultra Marathon"
	Distance    string    // e.g., "10K", "Half Marathon", "50K"
	Description string
	Races       []RaceViewModel
	Photos      []string
	MapURL      string
	Organizer   string
	Price       string
	Capacity    int
	Registered  int
}

// RaceViewModel represents a race within an event
type RaceViewModel struct {
	Name        string
	Distance    string
	Price       string
	StartTime   string
	Capacity    int
	Registered  int
	Description string
}

// FormattedDate returns the date in a human-readable format
func (e EventViewModel) FormattedDate() string {
	return e.Date.Format("2 January 2006")
}

// FormattedDay returns just the day number
func (e EventViewModel) FormattedDay() string {
	return e.Date.Format("02")
}

// FormattedMonth returns the abbreviated month
func (e EventViewModel) FormattedMonth() string {
	return e.Date.Format("Jan")
}

// FormattedYear returns the year
func (e EventViewModel) FormattedYear() string {
	return e.Date.Format("2006")
}

// SpotsRemaining returns the number of spots left
func (e EventViewModel) SpotsRemaining() int {
	return e.Capacity - e.Registered
}

// RegistrationPercentage returns how full the event is
func (e EventViewModel) RegistrationPercentage() int {
	if e.Capacity == 0 {
		return 0
	}
	return (e.Registered * 100) / e.Capacity
}

// GetMockEvents returns sample events for UI mockup
func GetMockEvents() []EventViewModel {
	return []EventViewModel{
		{
			Slug:        "peak-district-ultra-2026",
			Name:        "Peak District Ultra",
			Date:        time.Date(2026, 4, 18, 8, 0, 0, 0, time.UTC),
			Location:    "Castleton, Peak District",
			ImageURL:    "https://images.unsplash.com/photo-1551632811-561732d1e306?w=600&h=400&fit=crop",
			RaceType:    "Ultra Marathon",
			Distance:    "50K",
			Description: "Experience the breathtaking beauty of the Peak District on this challenging 50K ultra marathon. Wind through limestone valleys, climb iconic peaks, and test your limits on some of the finest trails in England.",
			Organizer:   "Peak Running Co",
			Price:       "£65",
			Capacity:    500,
			Registered:  342,
			MapURL:      "https://images.unsplash.com/photo-1524661135-423995f22d0b?w=800&h=400&fit=crop",
			Photos: []string{
				"https://images.unsplash.com/photo-1551632811-561732d1e306?w=800&h=600&fit=crop",
				"https://images.unsplash.com/photo-1469395446868-fb6a048d5ca3?w=800&h=600&fit=crop",
				"https://images.unsplash.com/photo-1483728642387-6c3bdd6c93e5?w=800&h=600&fit=crop",
			},
			Races: []RaceViewModel{
				{Name: "Ultra 50K", Distance: "50K", Price: "£65", StartTime: "06:00", Capacity: 300, Registered: 245, Description: "The main event - a challenging 50K route through the heart of the Peak District."},
				{Name: "Marathon", Distance: "42K", Price: "£55", StartTime: "07:00", Capacity: 200, Registered: 97, Description: "A full marathon distance covering the most scenic sections of the course."},
			},
		},
		{
			Slug:        "lake-district-trail-run",
			Name:        "Lake District Trail Run",
			Date:        time.Date(2026, 5, 9, 9, 0, 0, 0, time.UTC),
			Location:    "Ambleside, Lake District",
			ImageURL:    "https://images.unsplash.com/photo-1571104508999-893933ded431?w=600&h=400&fit=crop",
			RaceType:    "Trail Run",
			Distance:    "Half Marathon",
			Description: "A stunning half marathon through the Lake District National Park. Run alongside crystal-clear lakes, through ancient woodlands, and past iconic fells.",
			Organizer:   "Lakes Events",
			Price:       "£45",
			Capacity:    750,
			Registered:  512,
			MapURL:      "https://images.unsplash.com/photo-1524661135-423995f22d0b?w=800&h=400&fit=crop",
			Photos: []string{
				"https://images.unsplash.com/photo-1571104508999-893933ded431?w=800&h=600&fit=crop",
				"https://images.unsplash.com/photo-1501785888041-af3ef285b470?w=800&h=600&fit=crop",
			},
			Races: []RaceViewModel{
				{Name: "Half Marathon", Distance: "21K", Price: "£45", StartTime: "09:00", Capacity: 500, Registered: 389, Description: "The flagship half marathon with challenging ascents and incredible views."},
				{Name: "10K Fun Run", Distance: "10K", Price: "£25", StartTime: "10:30", Capacity: 250, Registered: 123, Description: "A scenic 10K perfect for beginners and families."},
			},
		},
		{
			Slug:        "yorkshire-three-peaks",
			Name:        "Yorkshire Three Peaks Challenge",
			Date:        time.Date(2026, 6, 14, 7, 0, 0, 0, time.UTC),
			Location:    "Horton-in-Ribblesdale, Yorkshire",
			ImageURL:    "https://images.unsplash.com/photo-1464822759023-fed622ff2c3b?w=600&h=400&fit=crop",
			RaceType:    "Fell Race",
			Distance:    "24 miles",
			Description: "Conquer the legendary Yorkshire Three Peaks in this iconic fell race. Summit Pen-y-ghent, Whernside, and Ingleborough in under 12 hours.",
			Organizer:   "Yorkshire Trails",
			Price:       "£50",
			Capacity:    600,
			Registered:  598,
			MapURL:      "https://images.unsplash.com/photo-1524661135-423995f22d0b?w=800&h=400&fit=crop",
			Photos: []string{
				"https://images.unsplash.com/photo-1464822759023-fed622ff2c3b?w=800&h=600&fit=crop",
				"https://images.unsplash.com/photo-1500534623283-312aade485b7?w=800&h=600&fit=crop",
			},
			Races: []RaceViewModel{
				{Name: "Three Peaks Challenge", Distance: "24 miles", Price: "£50", StartTime: "07:00", Capacity: 600, Registered: 598, Description: "The classic Three Peaks route with a 12-hour cutoff."},
			},
		},
		{
			Slug:        "cotswolds-spring-10k",
			Name:        "Cotswolds Spring 10K",
			Date:        time.Date(2026, 3, 22, 10, 0, 0, 0, time.UTC),
			Location:    "Bourton-on-the-Water, Cotswolds",
			ImageURL:    "https://images.unsplash.com/photo-1508739773434-c26b3d09e071?w=600&h=400&fit=crop",
			RaceType:    "Road Race",
			Distance:    "10K",
			Description: "A beautiful spring road race through picturesque Cotswold villages. Rolling hills, honey-stone cottages, and country lanes await.",
			Organizer:   "Cotswold Running Club",
			Price:       "£28",
			Capacity:    400,
			Registered:  156,
			MapURL:      "https://images.unsplash.com/photo-1524661135-423995f22d0b?w=800&h=400&fit=crop",
			Photos: []string{
				"https://images.unsplash.com/photo-1508739773434-c26b3d09e071?w=800&h=600&fit=crop",
			},
			Races: []RaceViewModel{
				{Name: "10K Race", Distance: "10K", Price: "£28", StartTime: "10:00", Capacity: 300, Registered: 112, Description: "A fast and scenic 10K through the Cotswolds countryside."},
				{Name: "5K Fun Run", Distance: "5K", Price: "£15", StartTime: "11:30", Capacity: 100, Registered: 44, Description: "A family-friendly 5K suitable for all abilities."},
			},
		},
		{
			Slug:        "snowdonia-marathon",
			Name:        "Snowdonia Marathon",
			Date:        time.Date(2026, 10, 24, 8, 0, 0, 0, time.UTC),
			Location:    "Llanberis, Snowdonia",
			ImageURL:    "https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=600&h=400&fit=crop",
			RaceType:    "Mountain Marathon",
			Distance:    "Marathon",
			Description: "One of the most scenic and challenging marathons in the UK. Run beneath the shadow of Snowdon through spectacular Welsh mountain scenery.",
			Organizer:   "Welsh Mountain Events",
			Price:       "£58",
			Capacity:    2000,
			Registered:  1456,
			MapURL:      "https://images.unsplash.com/photo-1524661135-423995f22d0b?w=800&h=400&fit=crop",
			Photos: []string{
				"https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=800&h=600&fit=crop",
				"https://images.unsplash.com/photo-1519904981063-b0cf448d479e?w=800&h=600&fit=crop",
			},
			Races: []RaceViewModel{
				{Name: "Full Marathon", Distance: "42.2K", Price: "£58", StartTime: "08:00", Capacity: 1500, Registered: 1123, Description: "The flagship Snowdonia Marathon with stunning mountain views."},
				{Name: "Half Marathon", Distance: "21.1K", Price: "£38", StartTime: "09:30", Capacity: 500, Registered: 333, Description: "A challenging half marathon through the Snowdonia foothills."},
			},
		},
		{
			Slug:        "south-downs-way-50",
			Name:        "South Downs Way 50",
			Date:        time.Date(2026, 7, 11, 6, 0, 0, 0, time.UTC),
			Location:    "Worthing, West Sussex",
			ImageURL:    "https://images.unsplash.com/photo-1551632811-561732d1e306?w=600&h=400&fit=crop",
			RaceType:    "Ultra Marathon",
			Distance:    "50 miles",
			Description: "Follow the ancient South Downs Way on this epic 50-mile ultra. Chalk grassland, coastal views, and Iron Age hill forts make this a truly memorable race.",
			Organizer:   "Centurion Running",
			Price:       "£95",
			Capacity:    350,
			Registered:  298,
			MapURL:      "https://images.unsplash.com/photo-1524661135-423995f22d0b?w=800&h=400&fit=crop",
			Photos: []string{
				"https://images.unsplash.com/photo-1551632811-561732d1e306?w=800&h=600&fit=crop",
			},
			Races: []RaceViewModel{
				{Name: "50 Mile Ultra", Distance: "50 miles", Price: "£95", StartTime: "06:00", Capacity: 350, Registered: 298, Description: "The full 50-mile route along the South Downs Way."},
			},
		},
	}
}

// GetMockEvent returns a single mock event by slug
func GetMockEvent(slug string) *EventViewModel {
	events := GetMockEvents()
	for _, e := range events {
		if e.Slug == slug {
			return &e
		}
	}
	return nil
}
