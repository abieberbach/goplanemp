package csl

type matchingResult struct {
	exactIcaoAirlineLiveryMatch   *CslAircraft
	exactIcaoAirlineMatch         *CslAircraft
	exactIcaoMatch                *CslAircraft
	relatedIcaoAirlineLiveryMatch *CslAircraft
	relatedIcaoAirlineMatch       *CslAircraft
	relatedIcaoMatch              *CslAircraft
}

func matchPlane(aircrafts []*CslAircraft, icao, airline, livery string) matchingResult {
	result := matchingResult{}
	for _, currentAircraft := range aircrafts {
		icaoMatched := currentAircraft.Icao == icao
		airlineMatched := currentAircraft.Airline == airline
		liveryMatched := currentAircraft.Livery == livery
		if icaoMatched {
			//ICAO stimmt überein
			if airlineMatched {
				//Airline stimmt überein
				if liveryMatched {
					//Livery stimmt auch überein
					result.exactIcaoAirlineLiveryMatch = currentAircraft
				}
				if result.exactIcaoAirlineMatch == nil || currentAircraft.Livery == "" {
					//Flugzeug nur übernehmen, wenn es noch kein Ergebnis gibt oder das Livery leer ist
					result.exactIcaoAirlineMatch = currentAircraft
				}
			}
			if result.exactIcaoMatch == nil || currentAircraft.Airline == "" {
				//Flugzeug nur übernehmen, wenn es noch kein Ergebnis gibt oder die Airline leer ist
				result.exactIcaoMatch = currentAircraft
			}
		} else {
			//es ist ein verwandter ICAO
			if airlineMatched {
				//Airline stimmt überein
				if liveryMatched {
					//Livery stimmt auch überein
					result.relatedIcaoAirlineLiveryMatch = currentAircraft
				}
				if result.relatedIcaoAirlineMatch == nil || currentAircraft.Livery == "" {
					//Flugzeug nur übernehmen, wenn es noch kein Ergebnis gibt oder das Livery leer ist
					result.relatedIcaoAirlineMatch = currentAircraft
				}
			}
			if result.relatedIcaoMatch == nil || currentAircraft.Airline == "" {
				//Flugzeug nur übernehmen, wenn es noch kein Ergebnis gibt oder die Airline leer ist
				result.relatedIcaoMatch = currentAircraft
			}
		}
	}
	return result
}