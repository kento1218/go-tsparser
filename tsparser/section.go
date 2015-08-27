// Copyright (c) 2014 Kohei YOSHIDA. All rights reserved.
// This software is licensed under the 3-Clause BSD License
// that can be found in LICENSE file.

package tsparser

type ProgramAssociationSection struct {
	network    PID
	programMap map[uint16]PID
}

func ParseProgramAssociationSection(table Table) *ProgramAssociationSection {
	sec := new(ProgramAssociationSection)
	sec.programMap = make(map[uint16]PID)

	payload := table.Data()
	for i := 0; i < len(payload); i += 4 {
		programNumber := uint16(payload[i])<<8 | uint16(payload[i+1])
		if programNumber == 0 {
			sec.network = PID(payload[i+2]&0x1f)<<8 | PID(payload[i+3])
		} else {
			sec.programMap[programNumber] = PID(payload[i+2]&0x1f)<<8 | PID(payload[i+3])
		}
	}

	return sec
}

func (s *ProgramAssociationSection) NetworkPID() PID {
	return s.network
}

func (s *ProgramAssociationSection) ProgramMap() map[uint16]PID {
	return s.programMap
}

type ProgramMapSection struct {
	pcr           PID
	descriptors   []Descriptor
	streamEntries []*StreamEntry
}

type StreamEntry struct {
	streamType    uint16
	elementaryPID PID
	descriptors   []Descriptor
}

func ParseProgramMapSection(table Table) *ProgramMapSection {
	sec := &ProgramMapSection{}

	data := table.Data()
	sec.pcr = PID(data[0]&0x1f)<<8 | PID(data[1])
	desclen := uint(data[2]&0x0f)<<8 | uint(data[3])
	sec.descriptors = ParseDescriptors(data[4 : 4+desclen])

	entryData := data[4+desclen:]
	for len(entryData) > 0 {
		pid := PID(entryData[1]&0x1f)<<8 | PID(entryData[2])

		entry := &StreamEntry{
			streamType:    uint16(entryData[0]),
			elementaryPID: pid,
		}
		sec.streamEntries = append(sec.streamEntries, entry)

		infolen := uint(entryData[3]&0x0f)<<8 | uint(entryData[4])
		entry.descriptors = ParseDescriptors(entryData[5 : 5+infolen])

		entryData = entryData[5+infolen:]
	}

	return sec
}

func (s *ProgramMapSection) PCR() PID {
	return s.pcr
}

func (s *ProgramMapSection) StreamEntries() []*StreamEntry {
	return s.streamEntries
}

func (s *ProgramMapSection) Descriptors() []Descriptor {
	return s.descriptors
}

func (e *StreamEntry) StreamType() uint16 {
	return e.streamType
}

func (e *StreamEntry) PID() PID {
	return e.elementaryPID
}

func (e *StreamEntry) Descriptors() []Descriptor {
	return e.descriptors
}
