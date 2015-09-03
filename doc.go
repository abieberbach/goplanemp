//Multiplayer-Bibliothek für X-Plane.
//Die Umsetzung basiert auf dem Code der libxplanemp Bibliothek.
//Weitere Informationen unter https://github.com/wadesworld2112/libxplanemp.
//
//Für die Darstellung können vorhandene CSL-Pakete genutzt werden. Zusätzlich unterstützt diese Bibliothek auch die Darstellung von OBJ8-Objekten.
//Bei den OBJ8-Objekte werden auch Animationen unterstützt.
//
//Die Bibliothek kann in einem X-Plane Plugin wie folgt genutzt werden:
//
//Erstellung des MultiplayerManagers
//
//Schritt 1 (während der Start-Phase des Plugins) MultiplayerManager angelegen
//
//	config := NewDefaultConfiguration()
//	mpManager, err := New(cslPath, relatedTxt, lightsPng, defaultIcao, config)
//
//Schritt 2 (während der Enable-Phase des Plugins) MultiplayerManager aktivieren
//
//	if mpManager != nil {
//		mpManager.Enable()
//	}
//
//Neue Flugzeuge und Fahrzeuge hinzufügen
//
//Flugzeuge und Fahrzeuge können über die CreatePlane-Funktion hinzugefügt werden. Dazu muss das ICAO-Kürzel des Flugzeugs und der Airline und ggf. ein Livery angegeben werden.
//Welches Objekt der MultiplayerManager auswählt, erfolgt nach dieser Reihenfolge:
//	1. ICAO, Airline und Livery stimmen überein
//	2. ICAO und Airline stimmen überein
//	3. verwandter ICAO, Airline und Livery stimmen überein
//	4. verwandter ICAO und Airline stimmen überein
//	5. ICAO stimmt überein
//	6. verwandter ICAO stimmt überein
//	7. Standard-ICAO stimmt überein
//	8. verwandter Standard-ICAO stimmt überein
//	9. beliebiges Flugzeug nehmen
//Beispiel zum Hinzufügen eines neuen Flugzeugs
//	planeId := mpManager.CreatePlane("B738", "BER", "WL", dataFunc, ref)
package goplanemp
