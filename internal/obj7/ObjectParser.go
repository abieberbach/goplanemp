//Copyright (c) 2015. The goplanemp AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

package obj7
import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"os"
	"strings"
	"bytes"
	"io"
	"path/filepath"
)

type ObjectParser struct {
	scanner     *bufio.Scanner
	baseDir     string
	texture     string
	commands    []*ObjectCommand
	err         error
	currentLine int
}

func (self *ObjectParser) GetTexture() string {
	if self.err != nil {
		return ""
	}
	return self.texture
}

func (self *ObjectParser) GetCommands() []*ObjectCommand {
	if self.err != nil {
		return make([]*ObjectCommand, 0)
	}
	return self.commands
}


func (self *ObjectParser) Err() error {
	return self.err
}

func (self *ObjectParser) CurrentLine() int {
	return self.currentLine
}

func (self *ObjectParser) Parse() {
	//LineFeed wird ignoriert
	self.parseLineFeed()
	version := self.parseVersion()
	if version != -1&&version != 7 {
		self.err = errors.New(fmt.Sprintf("only obj7 files version 7 are supported, current version: %v", version))
	}
	//Zeile mit "OBJ" parsen --> diese wird ignoriert
	self.parseObjLine()
	//Textur einlesen
	self.texture = filepath.Join(self.baseDir,self.parseTexture())
	for {
		//Geometrie einlesen
		command := self.parseCommand()
		if command == nil {
			//Dateiende wurde erreicht oder es gab ein anderen Fehler --> Schleife verlassen
			break
		}
		self.commands = append(self.commands, command)

	}
}

func (self *ObjectParser) parseLineFeed() string {
	tokens := self.readLine(true)
	if self.err != nil {
		return ""
	}
	return tokens[0]
}

func (self *ObjectParser) parseVersion() int {
	tokens := self.readLine(true)
	if self.err != nil {
		return -1
	}
	switch tokens[0] {
	case "700":
		return 7
	case "2":
		return 2
	default:
		return 1
	}
}

func (self *ObjectParser) parseObjLine() {
	tokens := self.readLine(true)
	if self.err != nil {
		return
	}
	if len(tokens) != 1 || tokens[0] != "OBJ" {
		self.err = errors.New("line with content \"OBJ\" expected")
	}
}

func (self *ObjectParser) parseTexture() string {
	tokens := self.readLine(true)
	if self.err != nil {
		return ""
	}
	return tokens[0]
}

func (self *ObjectParser) parseCommand() *ObjectCommand {
	if self.err != nil {
		return nil
	}
	tokens := self.readLine(false)
	if len(tokens) == 0 {
		//es gibt nichts mehr zu lesen
		return nil
	}
	command, elementCount := FindObjectCommand(tokens[0])
	if command == nil {
		//Kommando konnte nicht gefunden werden
		self.err = errors.New(fmt.Sprintf("command not found \"%v\"", tokens[0]))
		return nil
	}
	//Kommando existiert --> entsprechend des Kommandos weiter einlesen
	switch command.CmdType {
	case CommandType_None:
		if command.CmdId != CommandId_ObjEnd {
			self.err = errors.New("id \"obj_end\" expected")
		}
		//Ende des Objektes wurde erreicht  --> es gibt nichts mehr zu lesen
		return nil
	case CommandType_PtLine:
		if elementCount == 0 && len(tokens) > 1 {
			//Anzahl wird dynamisch ermittelt
			elementCount, self.err = strconv.Atoi(tokens[1])
		}
		command = self.parsePtLine(command, elementCount)
	case CommandType_Poly:
		if elementCount == 0 && len(tokens) > 1 {
			//Anzahl wird dynamisch ermittelt
			elementCount, self.err = strconv.Atoi(tokens[1])
		}
		command = self.parsePoly(command, elementCount)
	case CommandType_Attr:
		command = self.parseAttr(tokens, command, elementCount)
	}
	return command
}

func (self *ObjectParser) parsePtLine(command *ObjectCommand, count int) *ObjectCommand {
	for i := 0; i < count; i++ {
		if self.err != nil {
			return nil
		}
		tokens := self.readLine(true)
		if (len(tokens) > 5) {
			vrgb := VecRGB{}
			vrgb.V[0], self.err = parseFloat32(tokens[0])
			vrgb.V[1], self.err = parseFloat32(tokens[1])
			vrgb.V[2], self.err = parseFloat32(tokens[2])
			vrgb.RGB[0], self.err = parseFloat32(tokens[3])
			vrgb.RGB[1], self.err = parseFloat32(tokens[4])
			vrgb.RGB[2], self.err = parseFloat32(tokens[5])
			command.RGB = append(command.RGB, vrgb)
		} else {
			self.err = errors.New(fmt.Sprintf("unexpected token count (%v <= 5)", len(tokens)))
			return nil
		}

	}
	return command
}

func (self *ObjectParser) parsePoly(command *ObjectCommand, count int) *ObjectCommand {
	for i := 0; i < count; i++ {
		if self.err != nil {
			return nil
		}
		tokens := self.readLine(true)
		if len(tokens) > 4 {
			vst := VecTex{}
			vst.V[0], self.err = parseFloat32(tokens[0])
			vst.V[1], self.err = parseFloat32(tokens[1])
			vst.V[2], self.err = parseFloat32(tokens[2])
			vst.St[0], self.err = parseFloat32(tokens[3])
			vst.St[1], self.err = parseFloat32(tokens[4])
			command.St = append(command.St, vst)

			if len(tokens) > 9 {
				i++
				vst := VecTex{}
				vst.V[0], self.err = parseFloat32(tokens[5])
				vst.V[1], self.err = parseFloat32(tokens[6])
				vst.V[2], self.err = parseFloat32(tokens[7])
				vst.St[0], self.err = parseFloat32(tokens[8])
				vst.St[1], self.err = parseFloat32(tokens[9])
				command.St = append(command.St, vst)
			}
		} else {
			self.err = errors.New(fmt.Sprintf("unexpected token count (%v <= 4)", len(tokens)))
			return nil
		}
	}
	return command
}

func (self *ObjectParser) parseAttr(tokens []string, command *ObjectCommand, count int) *ObjectCommand {
	if len(tokens) > count {
		for n := 0; n < count; n++ {
			var value float32
			value, self.err = parseFloat32(tokens[n + 1])
			command.Attributes = append(command.Attributes, value)
		}
	} else {
		self.err = errors.New(fmt.Sprintf("unexpected token count (%v <= 4)", len(tokens)))
		return nil
	}
	return command
}

func (self *ObjectParser) readLine(mandatoryLine bool) []string {
	if self.err != nil {
		return []string{}
	}
	for self.scanner.Scan() {
		self.currentLine++
		text := self.scanner.Text()
		if self.scanner.Err() != nil {
			self.err = self.scanner.Err()
			return []string{}
		}
		//Kommentare abschneiden
		pos := strings.Index(text, "//")
		if pos >= 0 {
			text = text[:pos]
		}
		text = strings.TrimSpace(text)
		if len(text) == 0 {
			//Zeile enthält keinen Text ignorieren
			continue
		}
		text = strings.Replace(text, "\t", " ", -1)
		tokens := make([]string, 0)
		for _, token := range strings.Split(text, " ") {
			token = strings.TrimSpace(token)
			if len(token) > 0 {
				tokens = append(tokens, token)
			}
		}
		return tokens
	}
	if self.scanner.Err() != nil {
		self.err = self.scanner.Err()
		return []string{}
	}
	//es gibt nichts zu lesen --> prüfen ob eine Zeile gewünscht war
	if mandatoryLine {
		//ja es sollte eine Zeile gelesen werden --> Fehler
		self.err = io.ErrUnexpectedEOF
	}
	return []string{}
}

func NewObjectParser(objFilePath string) (*ObjectParser, error) {
	objFile, err := os.Open(objFilePath)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(objFile)
	scanner.Split(scanObjLines)
	return &ObjectParser{scanner, filepath.Dir(objFilePath), "", make([]*ObjectCommand, 0), nil, 0}, nil
}


func parseFloat32(token string) (float32, error) {
	value, err := strconv.ParseFloat(token, 32)
	return float32(value), err
}

func scanObjLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0:i], nil
	}
	if i := bytes.IndexByte(data, '\r'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}