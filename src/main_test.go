package main

import "testing"

func TestMainApp(t *testing.T) {
	result := "SunGo"
	if result != "SunGo" {
		t.Errorf("Oczekiwano SunGo, otrzymano %s", result)
	}
}

func TestNewTappableSelect(t *testing.T) {
	// TODO: wywołaj newTappableSelect z realnymi danymi wejściowymi i sprawdź wynik
	t.Skip("TODO: implement test for newTappableSelect")
}

func TestNewCustomEntry(t *testing.T) {
	// TODO: wywołaj NewCustomEntry z realnymi danymi wejściowymi i sprawdź wynik
	t.Skip("TODO: implement test for NewCustomEntry")
}

func TestSaveSettingsToFile(t *testing.T) {
	// TODO: wywołaj saveSettingsToFile z realnymi danymi wejściowymi i sprawdź wynik
	t.Skip("TODO: implement test for saveSettingsToFile")
}

func TestCreateBlinker(t *testing.T) {
	// TODO: wywołaj createBlinker z realnymi danymi wejściowymi i sprawdź wynik
	t.Skip("TODO: implement test for createBlinker")
}
