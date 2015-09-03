package csl
import (
	"bufio"
	"os"
	"strings"
	"github.com/abieberbach/goplane/extra/logging"
	"errors"
	"fmt"
	"github.com/abieberbach/goplane/xplm/utilities"
	"strconv"
)

var simVersion int

func init() {
	simVersion, _, _ = utilities.GetVersions()
}

func parseCslPackage(aircraftFile *os.File, baseDir string) (*CslPackage, error) {
	scanner := bufio.NewScanner(aircraftFile)
	result := &CslPackage{}
	result.Valid = true
	result.BaseDirectory = baseDir
	result.Dependencies = make([]string, 0)
	var currentAircraft *CslAircraft
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if len(text) == 0 || strings.HasPrefix(text, "#") {
			//Zeile enthält keinen Text bzw. nur Kommentare --> ignorieren
			continue
		}
		tokens := strings.Split(text, " ")
		tokens = trimTokens(tokens)
		var err error
		switch tokens[0] {
		case "EXPORT_NAME":
			err = handleExportName(tokens, result)
		case "DEPENDENCY":
			err = handleDependency(tokens, result)
		case "AIRCRAFT":
			currentAircraft, err = handleAircraft(tokens, currentAircraft, result, aircraftFile)
		case "OBJECT":
			currentAircraft, err = handleObject7(tokens, currentAircraft, result)
		case "OBJ8_AIRCRAFT":
			currentAircraft, err = handleObject8Aircraft(tokens, currentAircraft, result)
		case "OBJ8":
			err = handleObj8(tokens, currentAircraft, aircraftFile)
		case "TEXTURE":
			err = handleTexture(tokens, currentAircraft)
		case "HASGEAR":
			err = handleHasGear(tokens, currentAircraft, aircraftFile)
		case "ICAO":
			err = handleIcao(tokens, currentAircraft)
		case "AIRLINE":
			err = handleAirline(tokens, currentAircraft)
		case "LIVERY":
			err = handleLivery(tokens, currentAircraft)
		default:
			logging.Warningf("unexpected command \"%v\" in file %v", tokens[0], aircraftFile.Name())
		}
		if err != nil {
			return nil, err
		}
	}
	if currentAircraft != nil {
		result.Aircrafts = append(result.Aircrafts, currentAircraft)
	}
	return result, nil
}

func handleExportName(tokens []string, cslPackage *CslPackage) error {
	if len(tokens) != 2 {
		return errors.New("Command <EXPORT_NAME> missing package name")
	}
	cslPackage.Name = tokens[1]
	return nil
}

func handleDependency(tokens []string, cslPackage *CslPackage) error {
	if len(tokens) != 2 {
		return errors.New("Command <DEPENDENCY> missing package name")
	}
	cslPackage.Dependencies = append(cslPackage.Dependencies, tokens[1])
	return nil
}

func handleAircraft(tokens []string, currentAircraft *CslAircraft, cslPackage *CslPackage, aircraftFile *os.File) (*CslAircraft, error) {
	if len(tokens) != 4 {
		return nil, errors.New(fmt.Sprintf("Command <AIRCRAFT> invalid token count, expected 4, current %v", len(tokens)))
	}
	logging.Warningf("ignore austin aircraft (not supported): %v (file: %v)", tokens[3], aircraftFile.Name())
	return nil, nil

	/*
	minVersion := tokens[1]
	maxVersion := tokens[2]
	if !checkVersion(minVersion, maxVersion) {
	}
	newAircraft := appendAircraft(currentAircraft, cslPackage)
	newAircraft.AcfFile = tokens[3]
	newAircraft.CslType=CSL_TYPE_AUSTIN
	return newAircraft, nil
	*/
}

func handleObject7(tokens []string, currentAircraft *CslAircraft, cslPackage *CslPackage) (*CslAircraft, error) {
	if len(tokens) != 2 {
		return nil, errors.New("Command <OBJECT> missing object file name")
	}
	newAircraft := appendAircraft(currentAircraft, cslPackage)
	newAircraft.ObjFile = tokens[1]
	newAircraft.CslType = CSL_TYPE_OBJ7
	return newAircraft, nil
}

func handleObject8Aircraft(tokens []string, currentAircraft *CslAircraft, cslPackage *CslPackage) (*CslAircraft, error) {
	if len(tokens) != 1 {
		return nil, errors.New(fmt.Sprintf("Command <OBJ8_AIRCRAFT> invalid token count, expected 1, current %v", len(tokens)))
	}
	newAircraft := appendAircraft(currentAircraft, cslPackage)
	newAircraft.CslType = CSL_TYPE_OBJ8
	return newAircraft, nil
}


func handleObj8(tokens []string, currentAircraft *CslAircraft, aircraftFile *os.File) (error) {
	if len(tokens) < 4 || len(tokens) > 5 {
		return errors.New(fmt.Sprintf("Command <OBJ8> invalid token count, expected 4 (5), current %v", len(tokens)))
	}
	if currentAircraft == nil || currentAircraft.CslType != CSL_TYPE_OBJ8 {
		//aktuelles Flugzeug ist kein OBJ8-Flugzeug --> Eintrag ignorieren
		return nil
	}
	attachment := &CslAircraftAttachment{}
	attachment.ObjectLoadRequested = false
	if (tokens[1] == "GLASS") {
		attachment.DrawType = CSL_DRAW_GLASS;
	} else if (tokens[1] == "LIGHTS") {
		attachment.DrawType = CSL_DRAW_LIGHTS;
	} else if (tokens[1] == "LOW_LOD") {
		attachment.DrawType = CSL_DRAW_LOW_LOD
	}else if (tokens[1] == "SOLID") {
		attachment.DrawType = CSL_DRAW_SOLID
	} else {
		return errors.New(fmt.Sprintf("invalid group value: %v (file: %v)", tokens[1], aircraftFile.Name()))
	}

	if (tokens[2] == "YES") {
		attachment.NeedsAnimation = true
	}else if (tokens[2] == "NO") {
		attachment.NeedsAnimation = false
	} else {
		return errors.New(fmt.Sprintf("invalid animate value: %v (file: %v)", tokens[2], aircraftFile.Name()))
	}
	attachment.ObjFile = tokens[3]
	if len(tokens) > 4 {
		attachment.ReplaceTexture = tokens[4]
	}
	currentAircraft.Attachments = append(currentAircraft.Attachments, attachment)
	return nil
}

func handleTexture(tokens []string, currentAircraft *CslAircraft) error {
	if len(tokens) != 2 {
		return errors.New("Command <TEXTURE> missing texture file name")
	}
	if currentAircraft != nil {
		currentAircraft.TextureFile = tokens[1]
	}
	return nil
}

func handleHasGear(tokens []string, currentAircraft *CslAircraft, aircraftFile *os.File) error {
	if len(tokens) != 2 {
		return errors.New("Command <HASGEAR> missing value")
	}
	if currentAircraft != nil {
		switch tokens[1] {
		case "YES":
			currentAircraft.HasGear = true
		case "NO":
			currentAircraft.HasGear = false
		default:
			logging.Warningf("invalid HASGEAR-Value: %v (file: %v)", tokens[1], aircraftFile.Name())
		}
	}
	return nil
}

func handleIcao(tokens []string, currentAircraft *CslAircraft) error {
	if len(tokens) < 2 {
		return errors.New(fmt.Sprintf("Command <ICAO> invalid token count, expected 2, current %v", len(tokens)))
	}
	if currentAircraft != nil {
		currentAircraft.Icao = tokens[1]
	}
	return nil
}


func handleAirline(tokens []string, currentAircraft *CslAircraft) error {
	if len(tokens) < 3 {
		return errors.New(fmt.Sprintf("Command <AIRLINE> invalid token count, expected 3, current %v", len(tokens)))
	}
	if currentAircraft != nil {
		currentAircraft.Icao = tokens[1]
		currentAircraft.Airline = tokens[2]
	}
	return nil
}

func handleLivery(tokens []string, currentAircraft *CslAircraft) error {
	if len(tokens) < 4 {
		return errors.New(fmt.Sprintf("Command <LIVERY> invalid token count, expected 4, current %v", len(tokens)))
	}
	if currentAircraft != nil {
		currentAircraft.Icao = tokens[1]
		currentAircraft.Airline = tokens[2]
		currentAircraft.Livery = tokens[3]
	}
	return nil
}


func appendAircraft(aircraft *CslAircraft, cslPackage *CslPackage) *CslAircraft {
	if aircraft != nil {
		cslPackage.Aircrafts = append(cslPackage.Aircrafts, aircraft)
	}
	result := &CslAircraft{}
	result.HasGear = true
	result.Attachments = make([]*CslAircraftAttachment, 0)
	return result
}

func checkVersion(sMinVersion, sMaxVersion string) bool {
	minVersion, err := strconv.Atoi(sMinVersion)
	if err != nil {
		return false
	}
	maxVersion, err := strconv.Atoi(sMaxVersion)
	if err != nil {
		return false
	}
	return minVersion <= simVersion&&simVersion <= maxVersion
}

