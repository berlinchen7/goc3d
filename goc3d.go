package goc3d

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math"
	"os"
)

////////////////////////////////////////////////////////////////////////////////
// C3D Data
////////////////////////////////////////////////////////////////////////////////

type C3DData struct {
	Analog []C3DAnalog
	Points [][]C3DPoint
}

type C3DPoint struct {
	X        float32
	Y        float32
	Z        float32
	C        byte
	Residual byte
	Valid    bool
}

type C3DAnalog struct {
	x float32
	y float32
	z float32
}

////////////////////////////////////////////////////////////////////////////////
// END C3D Data
////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////
// C3D Group
////////////////////////////////////////////////////////////////////////////////

type C3DGroup struct {
	Name        string
	ID          int
	Description string
}

func (g C3DGroup) String() string {

	str := fmt.Sprintf("\nGroup name = %s\n", g.Name)
	str = fmt.Sprintf("%sGroup ID       = %d\n", str, g.ID)
	str = fmt.Sprintf("%sDescription    = %s\n", str, g.Description)

	return str
}

////////////////////////////////////////////////////////////////////////////////
// END C3D Group
////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////
// C3D Parameter
////////////////////////////////////////////////////////////////////////////////

type C3DParameter struct {
	Name           string
	GroupID        int
	DataType       int
	NrOfDimensions int
	Dimensions     []int
	Description    string
	DataLength     int
	ByteData       []byte
	StringData     []string
	RealData       []float32
	IntegerData    []int16
	Locked         bool
}

func (p C3DParameter) String() string {

	str := fmt.Sprintf("\nParameter name = %s\n", p.Name)
	str = fmt.Sprintf("%sGroup ID       = %d\n", str, p.GroupID)
	switch p.DataType {
	case CHAR:
		str = fmt.Sprintf("%sData Type      = CHAR\n", str)
	case INTEGER:
		str = fmt.Sprintf("%sData Type      = INTEGER\n", str)
	case REAL:
		str = fmt.Sprintf("%sData Type      = REAL\n", str)
	case BYTE:
		str = fmt.Sprintf("%sData Type      = BYTE\n", str)
	}

	dimStr := ""
	if p.NrOfDimensions > 0 {
		dimStr = fmt.Sprintf("%s[%d", dimStr, p.Dimensions[0])
		for i := 1; i < p.NrOfDimensions; i++ {
			dimStr = fmt.Sprintf("%s,%d", dimStr, p.Dimensions[i])
		}
		dimStr = fmt.Sprintf("%s]", dimStr)
	}

	str = fmt.Sprintf("%sDimensions     = %d %s\n", str, p.NrOfDimensions, dimStr)
	str = fmt.Sprintf("%sDescription    = %s\n", str, p.Description)

	str = fmt.Sprintf("%sData Length    = %d\n", str, p.DataLength)
	// str = fmt.Sprintf("%sData           = %s\n", str, p.ByteData)
	switch p.DataType {
	case CHAR:
		for i, v := range p.StringData {
			str = fmt.Sprintf("%s  Data[%d] = %s\n", str, i, v)
		}
	case INTEGER:
		for i, v := range p.IntegerData {
			str = fmt.Sprintf("%s  Data[%d] = %d\n", str, i, v)
		}
	case REAL:
		for i, v := range p.RealData {
			str = fmt.Sprintf("%s  Data[%d] = %f\n", str, i, v)
		}
	}

	if p.Locked {
		str = fmt.Sprintf("%sLocked         = true\n", str)
	} else {
		str = fmt.Sprintf("%sLocked         = false\n", str)
	}

	return str
}

////////////////////////////////////////////////////////////////////////////////
// END C3D Parameter
////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////
// C3D Info
////////////////////////////////////////////////////////////////////////////////

type C3DInfo struct {
	Parameters []C3DParameter
	Groups     []C3DGroup
}

const (
	CHAR int = iota
	BYTE
	INTEGER
	REAL
)

////////////////////////////////////////////////////////////////////////////////
// END C3D Info
////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////
// C3D Header
////////////////////////////////////////////////////////////////////////////////

type C3DHeader struct {
	Valid                  bool
	HasLabels              bool
	Uses4CharLabels        bool
	ParameterSection       int
	NrOfTrajectories       int
	NrOfAnalogMeasurements int
	FirstFrame             int
	LastFrame              int
	InterpolationGap       int
	DataStart              int
	NrOfAnalogSamples      int
	ScaleFactor            float32
	FrameRate              float32
	EventLabels            []string
	UsesInteger            bool
}

func (h C3DHeader) String() string {
	str := ""

	str = fmt.Sprintf("%sParameter section starts at       %d\n", str, h.ParameterSection)
	if h.Valid == true {
		str = fmt.Sprintf("Valid C3D                     = true\n")
	} else {
		str = fmt.Sprintf("Valid C3D                     = false\n")
	}
	str = fmt.Sprintf("%sNumber of trajectories        = %d\n", str, h.NrOfTrajectories)
	str = fmt.Sprintf("%sNumber of analog measurements = %d\n", str, h.NrOfAnalogMeasurements)
	str = fmt.Sprintf("%sFirst frame                   = %d\n", str, h.FirstFrame)
	str = fmt.Sprintf("%sLast frame                    = %d\n", str, h.LastFrame)
	str = fmt.Sprintf("%sMax gap                       = %d\n", str, h.InterpolationGap)
	str = fmt.Sprintf("%sData start                    = %d\n", str, h.DataStart)
	str = fmt.Sprintf("%sNumber of analog samples      = %d\n", str, h.NrOfAnalogSamples)
	str = fmt.Sprintf("%sScale factor                  = %f\n", str, h.ScaleFactor)
	str = fmt.Sprintf("%sFrame rate                    = %f\n", str, h.FrameRate)

	if h.HasLabels == true {
		str = fmt.Sprintf("%sHas labels                    = true\n", str)
	} else {
		str = fmt.Sprintf("%sHas labels                    = false\n", str)
	}

	if h.Uses4CharLabels == true {
		str = fmt.Sprintf("%sUses 4 Char Labels            = true\n", str)
	} else {
		str = fmt.Sprintf("%sUses 4 Char Labels            = false\n", str)
	}

	if h.UsesInteger == true {
		str = fmt.Sprintf("%sData is given in ints         = true\n", str)
	} else {
		str = fmt.Sprintf("%sData is given in ints         = false\n", str)
	}

	return str
}

////////////////////////////////////////////////////////////////////////////////
// END C3D Header
////////////////////////////////////////////////////////////////////////////////

func imin(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func imax(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func bytesToInt(bytes []byte) (r int) {
	if len(bytes) != 2 {
		panic("bytesToInt: only takes two bytes")
	}
	bits := binary.LittleEndian.Uint16(bytes)
	r = int(int16(bits))
	return
}

// this function is only used for the header section
func readWordAsInt(bytes []byte, index int) (r int) {
	wordIndex := index * 2
	b := make([]byte, 2, 2)
	b[0] = bytes[wordIndex]
	b[1] = bytes[wordIndex+1]
	bits := binary.LittleEndian.Uint16(b)
	r = int(int16(bits))
	return
}

// this function is only used for the header section
func read2WordsAsFloat(bytes []byte, index int) (r float32) {
	wordIndex := index * 2
	b := make([]byte, 4, 4)
	b[0] = bytes[wordIndex]
	b[1] = bytes[wordIndex+1]
	b[2] = bytes[wordIndex+2]
	b[3] = bytes[wordIndex+3]
	bits := binary.LittleEndian.Uint32(b)
	r = math.Float32frombits(bits)
	return r
}

func readEventLabels(bytes []byte, start, end int) (r []string) {
	for i := start; i < end; i += 4 {
		r = append(r, string(bytes[i:i+4]))
	}
	return
}

func readHeader(io *bufio.Reader) (r C3DHeader) {
	bytes := make([]byte, 512, 512)
	_, aerr := io.Read(bytes)
	check(aerr)

	if bytes[0] != 2 {
		panic(fmt.Sprintf("Not a C3D file. First byte is not 2 but %d", bytes[0]))
	}

	if bytes[1] != 80 {
		panic(fmt.Sprintf("Not a C3D file. Second byte is not 80 but %d", bytes[1]))
	}

	r.Valid = true

	r.ParameterSection = int(bytes[0])                       // first byte
	r.NrOfTrajectories = readWordAsInt(bytes, 1)             // word 2
	r.NrOfAnalogMeasurements = readWordAsInt(bytes, 2)       // word 3
	r.FirstFrame = readWordAsInt(bytes, 3)                   // word 4
	r.LastFrame = readWordAsInt(bytes, 4)                    // word 5
	r.InterpolationGap = readWordAsInt(bytes, 5)             // word 6
	r.ScaleFactor = read2WordsAsFloat(bytes, 6)              // word 7 - 8
	r.DataStart = readWordAsInt(bytes, 8)                    // word 9
	r.NrOfAnalogSamples = readWordAsInt(bytes, 9)            // word 10
	r.FrameRate = read2WordsAsFloat(bytes, 10)               // word 11 - 12
	r.HasLabels = (readWordAsInt(bytes, 147) == 12345)       // word 148
	r.Uses4CharLabels = (readWordAsInt(bytes, 149) == 12345) // word 150
	r.EventLabels = readEventLabels(bytes, 198, 233)         // word 199 - 234

	if r.ScaleFactor < 0 {
		r.UsesInteger = false
	} else {
		r.UsesInteger = true
	}

	return
}

func parseStringData(data []byte, dimensions []int) (r []string) {
	nr := 1
	strLen := dimensions[len(dimensions)-1]
	for i := 0; i < len(dimensions)-1; i++ { // last dimension is the string length
		nr *= dimensions[i]
	}

	if nr == 0 {
		return
	}

	r = make([]string, nr, nr)

	for i := 0; i < nr; i++ {
		a := i * strLen
		b := (i + 1) * strLen
		r[i] = string(data[a:b])
	}

	return
}

func parseIntData(data []byte, dimensions []int) (r []int16) {
	b := make([]byte, 2, 2)
	nr := 1
	for i := 0; i < len(dimensions); i++ { // last dimension is the string length
		nr *= dimensions[i]
	}

	if nr == 0 {
		return
	}

	r = make([]int16, nr, nr)

	for i := 0; i < nr; i++ {
		b[0] = data[i*2]
		b[1] = data[i*2+1]
		r[i] = int16(binary.LittleEndian.Uint16(b))
	}

	return
}

func parseRealData(data []byte, dimensions []int) (r []float32) {
	b := make([]byte, 4, 4)
	nr := 1
	for i := 0; i < len(dimensions); i++ { // last dimension is the string length
		nr *= dimensions[i]
	}

	if nr == 0 {
		return
	}

	r = make([]float32, nr, nr)

	for i := 0; i < nr; i++ {
		b[0] = data[i*4]
		b[1] = data[i*4+1]
		b[2] = data[i*4+2]
		b[3] = data[i*4+3]
		bits := binary.LittleEndian.Uint32(b)
		r[i] = float32(math.Float32frombits(bits))
	}

	return
}

func parseParameterBlock(bytes []byte, info *C3DInfo) {
	nrOfBytes := len(bytes)
	byteIndex := 0
	b := make([]byte, 2, 2)

	var groups []C3DGroup
	var parameters []C3DParameter

	for byteIndex < nrOfBytes {
		nrOfCharactersInName := int(uint8(bytes[byteIndex]))
		byteIndex++

		// in case there are filling zero-bytes, e.g. at the end of the parameter block
		if nrOfCharactersInName == 0 {
			break
			continue
		}

		groupID := int(int8(bytes[byteIndex])) // byte 2
		byteIndex++
		locked := false
		if nrOfCharactersInName < 0 {
			locked = true
			nrOfCharactersInName = -nrOfCharactersInName
		}

		name := string(bytes[byteIndex : byteIndex+nrOfCharactersInName])
		byteIndex += nrOfCharactersInName

		// reading byte offset to next block
		b[0] = bytes[byteIndex]
		b[1] = bytes[byteIndex+1]
		offset := int16(binary.LittleEndian.Uint16(b))
		nextBlockStartsAt := int(offset) + byteIndex
		byteIndex += 2

		if groupID > 0 { // we have a parameter
			nrOfBytes := int(int8(bytes[byteIndex]))
			byteIndex++

			dataType := CHAR
			switch nrOfBytes {
			case -1:
				dataType = CHAR
			case 1:
				dataType = BYTE
			case 2:
				dataType = INTEGER
			case 4:
				dataType = REAL
			default:
				panic(fmt.Sprintf("Wrong number of bytes detected in parsing of parameters: %d", nrOfBytes))
				// dataType = UNKNOWN
			}

			if nrOfBytes < 0 {
				nrOfBytes = -nrOfBytes
			}

			nrOfDimensions := int(uint8(bytes[byteIndex]))
			byteIndex++

			if nrOfDimensions < 0 || nrOfDimensions > 7 {
				panic(fmt.Sprintf("Number of dimensions must be in [0,7] but is %d", nrOfDimensions))
			}
			dimensions := make([]int, nrOfDimensions, nrOfDimensions)
			dataSize := nrOfBytes
			if dataSize < 0 {
				dataSize = -dataSize
			}
			for i := 0; i < int(nrOfDimensions); i++ {
				n := int(uint8(bytes[byteIndex]))
				dimensions[nrOfDimensions-i-1] = n
				dataSize *= n
				byteIndex++
			}
			data := bytes[byteIndex : byteIndex+dataSize]
			byteIndex += dataSize

			var stringData []string
			var intData []int16
			var realData []float32

			switch dataType {
			case CHAR:
				stringData = parseStringData(data, dimensions)
			case INTEGER:
				intData = parseIntData(data, dimensions)
			case REAL:
				realData = parseRealData(data, dimensions)
			}

			descriptionLength := int(int8(bytes[byteIndex]))
			parameterDescription := ""
			if descriptionLength > 0 {
				parameterDescription = string(bytes[byteIndex : byteIndex+descriptionLength])
			}
			p := C3DParameter{GroupID: groupID,
				Name:           name,
				Description:    parameterDescription,
				DataType:       dataType,
				NrOfDimensions: nrOfDimensions,
				Dimensions:     dimensions,
				DataLength:     dataSize,
				ByteData:       data,
				StringData:     stringData,
				RealData:       realData,
				IntegerData:    intData,
				Locked:         locked}
			parameters = append(parameters, p)

		} else { // we have a group
			groupID = -groupID // groups have negative group ids

			descriptionLength := int(uint8(bytes[byteIndex]))
			groupDescription := ""
			if descriptionLength > 0 {
				groupDescription = string(bytes[byteIndex : byteIndex+descriptionLength])
			}
			g := C3DGroup{ID: groupID, Name: name, Description: groupDescription}
			groups = append(groups, g)

		}
		if offset > 0 {
			byteIndex = nextBlockStartsAt
		} else {
			break
		}
	}
	info.Groups = groups
	info.Parameters = parameters
}

func readParameters(io *bufio.Reader) C3DInfo {
	info := C3DInfo{Parameters: nil, Groups: nil}

	bytes := make([]byte, 4, 4)
	_, err := io.Read(bytes)
	check(err)

	nrOfParameterBlocks := int(uint8(bytes[3]))

	var block []byte
	for i := 0; i < nrOfParameterBlocks*512-4; i++ {
		b, berr := io.ReadByte()
		check(berr)
		block = append(block, b)
	}

	parseParameterBlock(block, &info)

	return info
}

func parseIntegerPointData(bytes []byte, scaleFactor float32) (float32, float32, float32, byte, byte, bool) {

	x := float32(int16(binary.LittleEndian.Uint16(bytes[0:2]))) * scaleFactor
	y := float32(int16(binary.LittleEndian.Uint16(bytes[2:4]))) * scaleFactor
	z := float32(int16(binary.LittleEndian.Uint16(bytes[4:6]))) * scaleFactor

	ok := int16(binary.LittleEndian.Uint16(bytes[6:8])) > 0

	return x, y, z, bytes[6], bytes[7], ok
}

func parseFloatingPointData(bytes []byte) (float32, float32, float32, byte, byte, bool) {

	x := math.Float32frombits(binary.LittleEndian.Uint32(bytes[0:4]))
	y := math.Float32frombits(binary.LittleEndian.Uint32(bytes[4:8]))
	z := math.Float32frombits(binary.LittleEndian.Uint32(bytes[8:12]))

	ok := int32(binary.LittleEndian.Uint16(bytes[12:16])) > 0

	return x, y, z, bytes[12], bytes[13], ok
}

func read3DIntDataOnly(io *bufio.Reader, header C3DHeader) C3DData {
	data := C3DData{Analog: nil, Points: nil}
	bytes := make([]byte, 8, 8)
	nrOfFrames := header.LastFrame - header.FirstFrame
	nrOfTrajectories := header.NrOfTrajectories

	p := make([][]C3DPoint, nrOfTrajectories, nrOfTrajectories)
	for i := 0; i < nrOfTrajectories; i++ {
		p[i] = make([]C3DPoint, nrOfFrames, nrOfFrames)
	}

	for frame := 0; frame < nrOfFrames; frame++ {
		for trajectory := 0; trajectory < nrOfTrajectories; trajectory++ {
			_, err := io.Read(bytes)
			check(err)
			x, y, z, cam, res, ok := parseIntegerPointData(bytes, header.ScaleFactor)
			p[trajectory][frame].X = x
			p[trajectory][frame].Y = y
			p[trajectory][frame].Z = z
			p[trajectory][frame].C = cam
			p[trajectory][frame].Residual = res
			p[trajectory][frame].Valid = ok
		}
	}
	data.Points = p
	return data
}

func read3DFloatDataOnly(io *bufio.Reader, header C3DHeader) C3DData {
	data := C3DData{Analog: nil, Points: nil}
	bytes := make([]byte, 16, 16)
	nrOfFrames := header.LastFrame - header.FirstFrame
	nrOfTrajectories := header.NrOfTrajectories

	fmt.Println("Number of Frames:", nrOfFrames)

	p := make([][]C3DPoint, nrOfTrajectories, nrOfTrajectories)
	for i := 0; i < nrOfTrajectories; i++ {
		p[i] = make([]C3DPoint, nrOfFrames, nrOfFrames)
	}

	for frame := 0; frame < nrOfFrames; frame++ {
		for trajectory := 0; trajectory < nrOfTrajectories; trajectory++ {
			_, err := io.Read(bytes)
			check(err)
			x, y, z, cam, res, ok := parseFloatingPointData(bytes)
			p[trajectory][frame].X = x
			p[trajectory][frame].Y = y
			p[trajectory][frame].Z = z
			p[trajectory][frame].C = cam
			p[trajectory][frame].Residual = res
			p[trajectory][frame].Valid = ok
		}
	}
	data.Points = p
	return data
}

func read3DFloatAndAnalogData(io *bufio.Reader, header C3DHeader) C3DData {
	data := C3DData{Analog: nil, Points: nil}
	bytes := make([]byte, 16, 16)
	b := make([]byte, 4, 4)
	nrOfFrames := header.LastFrame - header.FirstFrame
	nrOfTrajectories := header.NrOfTrajectories
	nfOfSamplesPerFrame := header.NrOfAnalogMeasurements / header.NrOfAnalogSamples

	fmt.Println("Number of Trajectories:", nrOfTrajectories)

	p := make([][]C3DPoint, nrOfTrajectories, nrOfTrajectories)
	for i := 0; i < nrOfTrajectories; i++ {
		p[i] = make([]C3DPoint, nrOfFrames, nrOfFrames)
	}

	for frame := 0; frame < nrOfFrames; frame++ {
		for trajectory := 0; trajectory < nrOfTrajectories; trajectory++ {
			_, err := io.Read(bytes)
			check(err)
			x, y, z, cam, res, ok := parseFloatingPointData(bytes)
			p[trajectory][frame].X = x
			p[trajectory][frame].Y = y
			p[trajectory][frame].Z = z
			p[trajectory][frame].C = cam
			p[trajectory][frame].Residual = res
			p[trajectory][frame].Valid = ok
			fmt.Println(frame, trajectory)

			for a := 0; a < nfOfSamplesPerFrame; a++ {
				_, err = io.Read(b)
				check(err)
			}
		}
	}
	data.Points = p
	return data
}

func read3DIntAndAnalogData(io *bufio.Reader, header C3DHeader) C3DData {
	data := C3DData{Analog: nil, Points: nil}
	bytes := make([]byte, 8, 8)
	nrOfFrames := header.LastFrame - header.FirstFrame
	nrOfTrajectories := header.NrOfTrajectories

	fmt.Println("Number of Trajectories:", nrOfTrajectories)

	p := make([][]C3DPoint, nrOfTrajectories, nrOfTrajectories)
	for i := 0; i < nrOfTrajectories; i++ {
		p[i] = make([]C3DPoint, nrOfFrames, nrOfFrames)
	}

	for frame := 0; frame < nrOfFrames; frame++ {
		for trajectory := 0; trajectory < nrOfTrajectories; trajectory++ {
			_, err := io.Read(bytes)
			check(err)
			x, y, z, cam, res, ok := parseIntegerPointData(bytes, header.ScaleFactor)
			if frame < 2 {
				fmt.Println(frame, trajectory, x, y, z)
			}
			p[trajectory][frame].X = x
			p[trajectory][frame].Y = y
			p[trajectory][frame].Z = z
			p[trajectory][frame].C = cam
			p[trajectory][frame].Residual = res
			p[trajectory][frame].Valid = ok
			for a := 0; a < header.NrOfAnalogMeasurements; a++ {
				_, err = io.Read(bytes)
				check(err)
			}
		}
	}
	data.Points = p
	return data
}

func readData(io *bufio.Reader, header C3DHeader) C3DData {
	var data C3DData

	if header.NrOfAnalogSamples == 0 {
		if header.UsesInteger == true {
			data = read3DIntDataOnly(io, header)
		} else {
			data = read3DFloatDataOnly(io, header)
		}
	} else {
		fmt.Println("hier 3")
		if header.UsesInteger == true {
			fmt.Println("hier 3a")
			data = read3DIntAndAnalogData(io, header)
		} else {
			fmt.Println("hier 3b")
			data = read3DFloatAndAnalogData(io, header)
		}

	}

	return data
}

func ReadC3D(filename string) (C3DHeader, C3DInfo, C3DData) {
	f, err := os.Open(filename)
	defer f.Close()
	check(err)

	bufr := bufio.NewReader(f)
	header := readHeader(bufr)
	info := readParameters(bufr)
	f.Seek(int64(header.DataStart*512), 0)
	data := readData(bufr, header)

	return header, info, data
}
