package tsparser

import (
	"encoding/hex"
	"testing"
)

func TestParseProgramMapSection(t *testing.T) {
	raw, _ := hex.DecodeString("02b03105a0c50000e150f007c10188de02efff1be151f0035201810fe152f00352018306e154f008520187fd030012ad")
	table := Table(raw)

	sec := ParseProgramMapSection(table)
	if sec == nil {
		t.Errorf("couldnt parse")
	}
	if sec.PCR() != 0x150 {
		t.Errorf("incorrect PCR: 0x%x", sec.PCR())
	}

	streamEntries := sec.StreamEntries()
	if len(streamEntries) != 3 {
		t.Errorf("3 stream entry expected, but %d", len(streamEntries))
	}
	for _, entry := range streamEntries {
		switch {
		case entry.PID() == 0x151 && entry.StreamType() == 0x1b:
		case entry.PID() == 0x152 && entry.StreamType() == 0x0f:
		case entry.PID() == 0x154 && entry.StreamType() == 0x06:
		default:
			t.Errorf("unexpected entry: %v", entry)
		}
	}
}

func TestParseServiceDescriptionSection(t *testing.T) {
	raw, _ := hex.DecodeString("42f0787fe4d300007fe4ff0420f30018480a0100071b7cd5b8c6ecd3c10184cf0701fe01f002ffff0421f30014480a0100071b7cd5b8c6ecd3c10184cf0302fe010422f30014480a0100071b7cd5b8c6ecd3c10184cf0302fe0105a0e50018480ac000071b7cd5b8c6ecd3c10188cf07030e8946554a49")
	table := Table(raw)

	sec := ParseServiceDescriptionSection(table)
	if sec == nil {
		t.Errorf("couldnt parse")
	}

	serviceEntries := sec.ServiceEntries()
	if len(serviceEntries) != 4 {
		t.Errorf("4 services expected, but %d", len(serviceEntries))
	}
	for _, entry := range serviceEntries {
		if len(entry.Descriptors()) == 0 {
			t.Errorf("invalid service descriptors")
		}
	}
}
