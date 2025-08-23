package services

// MapService provides map-related utilities for League of Legends
type MapService struct {
	// Map dimensions and zone definitions
}

// NewMapService creates a new map service
func NewMapService() *MapService {
	return &MapService{}
}

// GetMapBounds returns the boundaries of Summoner's Rift
func (ms *MapService) GetMapBounds() (width, height int) {
	return 14870, 14870 // Summoner's Rift dimensions
}

// GetStrategicZones returns all strategic zones on the map
func (ms *MapService) GetStrategicZones() []MapZone {
	return []MapZone{
		{
			Name:        "Dragon Pit",
			Coordinates: [][]int{{9800, 4200}, {10200, 4200}, {10200, 4600}, {9800, 4600}},
			Strategic:   true,
			GamePhase:   []string{"mid", "late"},
		},
		{
			Name:        "Baron Pit",
			Coordinates: [][]int{{4800, 10200}, {5200, 10200}, {5200, 10600}, {4800, 10600}},
			Strategic:   true,
			GamePhase:   []string{"late"},
		},
		{
			Name:        "Blue Side Blue Buff",
			Coordinates: [][]int{{3800, 8000}, {4200, 8000}, {4200, 8400}, {3800, 8400}},
			Strategic:   true,
			GamePhase:   []string{"early", "mid"},
		},
		{
			Name:        "Red Side Red Buff",
			Coordinates: [][]int{{10800, 6600}, {11200, 6600}, {11200, 7000}, {10800, 7000}},
			Strategic:   true,
			GamePhase:   []string{"early", "mid"},
		},
		{
			Name:        "River Bushes",
			Coordinates: [][]int{{6000, 6000}, {9000, 6000}, {9000, 9000}, {6000, 9000}},
			Strategic:   true,
			GamePhase:   []string{"early", "mid", "late"},
		},
		{
			Name:        "Top Lane Tri-Bush",
			Coordinates: [][]int{{2300, 10800}, {2700, 10800}, {2700, 11200}, {2300, 11200}},
			Strategic:   false,
			GamePhase:   []string{"early", "mid"},
		},
		{
			Name:        "Bot Lane Tri-Bush",
			Coordinates: [][]int{{12000, 3800}, {12400, 3800}, {12400, 4200}, {12000, 4200}},
			Strategic:   false,
			GamePhase:   []string{"early", "mid"},
		},
	}
}

// MapZone represents a zone on the League of Legends map
type MapZone struct {
	Name        string    `json:"name"`
	Coordinates [][]int   `json:"coordinates"` // Polygon coordinates
	Strategic   bool      `json:"strategic"`   // High-value zone
	GamePhase   []string  `json:"game_phase"`  // When this zone is most important
}

// IsPointInZone checks if a coordinate is within a specific zone
func (ms *MapService) IsPointInZone(x, y int, zone MapZone) bool {
	return ms.isPointInPolygon(x, y, zone.Coordinates)
}

// GetZoneForPoint returns the zone name for a given coordinate
func (ms *MapService) GetZoneForPoint(x, y int) string {
	zones := ms.GetStrategicZones()
	for _, zone := range zones {
		if ms.IsPointInZone(x, y, zone) {
			return zone.Name
		}
	}
	return "unknown"
}

// GetStrategicZonesForPhase returns strategic zones relevant for a game phase
func (ms *MapService) GetStrategicZonesForPhase(phase string) []MapZone {
	zones := ms.GetStrategicZones()
	var filteredZones []MapZone
	
	for _, zone := range zones {
		if zone.Strategic {
			for _, phaseStr := range zone.GamePhase {
				if phaseStr == phase {
					filteredZones = append(filteredZones, zone)
					break
				}
			}
		}
	}
	
	return filteredZones
}

// Ray casting algorithm to determine if point is in polygon
func (ms *MapService) isPointInPolygon(x, y int, polygon [][]int) bool {
	n := len(polygon)
	inside := false
	
	j := n - 1
	for i := 0; i < n; i++ {
		if ((polygon[i][1] > y) != (polygon[j][1] > y)) &&
			(x < (polygon[j][0]-polygon[i][0])*(y-polygon[i][1])/(polygon[j][1]-polygon[i][1])+polygon[i][0]) {
			inside = !inside
		}
		j = i
	}
	
	return inside
}

// GetLaneForPosition determines which lane a coordinate belongs to
func (ms *MapService) GetLaneForPosition(x, y int) string {
	// Simplified lane detection based on map quadrants
	centerX, centerY := 7435, 7435 // Map center
	
	if x < centerX-1000 && y > centerY+1000 {
		return "top"
	} else if x > centerX+1000 && y < centerY-1000 {
		return "bottom"
	} else if abs(x-centerX) < 2000 && abs(y-centerY) < 2000 {
		return "mid"
	} else if (x < centerX && y < centerY) || (x > centerX && y > centerY) {
		return "jungle"
	}
	
	return "unknown"
}

// GetMapSide determines if a position is on blue or red side
func (ms *MapService) GetMapSide(x, y int) string {
	centerX, centerY := 7435, 7435 // Map center
	
	// Blue side is generally bottom-left, red side is top-right
	if x+y < centerX+centerY {
		return "blue"
	}
	return "red"
}

// CalculateDistance calculates distance between two points
func (ms *MapService) CalculateDistance(x1, y1, x2, y2 int) float64 {
	dx := float64(x1 - x2)
	dy := float64(y1 - y2)
	return (dx*dx + dy*dy) // Return squared distance for performance
}

// GetNearbyZones returns zones within a certain distance of a point
func (ms *MapService) GetNearbyZones(x, y int, maxDistance float64) []string {
	zones := ms.GetStrategicZones()
	var nearbyZones []string
	
	for _, zone := range zones {
		// Calculate center of zone
		centerX, centerY := ms.getZoneCenter(zone)
		distance := ms.CalculateDistance(x, y, centerX, centerY)
		
		if distance <= maxDistance*maxDistance { // Compare squared distances
			nearbyZones = append(nearbyZones, zone.Name)
		}
	}
	
	return nearbyZones
}

// getZoneCenter calculates the center point of a zone
func (ms *MapService) getZoneCenter(zone MapZone) (int, int) {
	if len(zone.Coordinates) == 0 {
		return 0, 0
	}
	
	sumX, sumY := 0, 0
	for _, coord := range zone.Coordinates {
		sumX += coord[0]
		sumY += coord[1]
	}
	
	return sumX / len(zone.Coordinates), sumY / len(zone.Coordinates)
}

// Helper function for absolute value
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}