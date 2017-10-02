package goc3d

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math"
	"os"

	pb "gopkg.in/cheggaaa/pb.v1"
)

const (
	CHAR int = iota
	BYTE
	INTEGER
	REAL
	UNKNOWN
)

type C3DHeader struct {
	Valid            bool
	HasLabels        bool
	Uses4CharLabels  bool
	ParameterSection int
	NrOfTrajectories int
	NrOfMeasurements int
	FirstFrame       int
	LastFrame        int
	InterpolationGap int
	DataStart        int
	NrOfSamples      int
	ScaleFactor      float64
	FrameRate        float64
	EventLabels      []string
}

type C3DParameter struct {
	Name        string
	GroupID     int
	DataType    int
	Dimensions  int
	Description string
	Data        []byte
	Locked      bool
}

type C3DGroup struct {
	Name        string
	ID          int
	Description string
}

type C3DInfo struct {
	Parameters []C3DParameter
	Groups     []C3DGroup
}

func (h C3DHeader) String() string {
	str := ""

	if h.Valid == true {
		str = fmt.Sprintf("Valid C3D              = true\n")
	} else {
		str = fmt.Sprintf("Valid C3D              = false\n")
	}

	str = fmt.Sprintf("%sParameter section starts at %d\n", str, h.ParameterSection)
	str = fmt.Sprintf("%sNumber of trajectories = %d\n", str, h.NrOfTrajectories)
	str = fmt.Sprintf("%sNumber of measurements = %d\n", str, h.NrOfMeasurements)
	str = fmt.Sprintf("%sFirst frame            = %d\n", str, h.FirstFrame)
	str = fmt.Sprintf("%sLast frame             = %d\n", str, h.LastFrame)
	str = fmt.Sprintf("%sMax gap                = %d\n", str, h.InterpolationGap)
	str = fmt.Sprintf("%sData start             = %d\n", str, h.DataStart)
	str = fmt.Sprintf("%sNumber of samples      = %d\n", str, h.NrOfSamples)
	str = fmt.Sprintf("%sScale factor           = %f\n", str, h.ScaleFactor)
	str = fmt.Sprintf("%sFrame rate             = %f\n", str, h.FrameRate)

	if h.HasLabels == true {
		str = fmt.Sprintf("%sHas labels             = true\n", str)
	} else {
		str = fmt.Sprintf("%sHas labels             = false\n", str)
	}

	if h.Uses4CharLabels == true {
		str = fmt.Sprintf("%sUses 4 Char Labels     = true\n", str)
	} else {
		str = fmt.Sprintf("%sUses 4 Char Labels     = false\n", str)
	}

	return str
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func readWordAsInt(bytes []byte, index int) (r int) {
	wordIndex := index * 2
	b := make([]byte, 2, 2)
	b[0] = bytes[wordIndex]
	b[1] = bytes[wordIndex+1]
	bits := binary.LittleEndian.Uint16(b)
	r = int(int16(bits))
	return
}

func read2WordsAsFloat(bytes []byte, index int) (r float64) {
	wordIndex := index * 2
	b := make([]byte, 4, 4)
	b[0] = bytes[wordIndex]
	b[1] = bytes[wordIndex+1]
	b[2] = bytes[wordIndex+2]
	b[3] = bytes[wordIndex+3]
	bits := binary.LittleEndian.Uint32(b)
	r = float64(math.Float32frombits(bits))
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
		fmt.Println("Not a C3D file")
		fmt.Println(fmt.Sprintf("First byte is not 2 but %d", bytes[0]))
		os.Exit(-1)
	}

	if bytes[1] != 80 {
		fmt.Println("Not a C3D file")
		fmt.Println(fmt.Sprintf("Second byte is not 80 but %d", bytes[1]))
		os.Exit(-1)
	}

	r.Valid = true

	r.ParameterSection = int(bytes[0])                       // first byte
	r.NrOfTrajectories = readWordAsInt(bytes, 1)             // word 2
	r.NrOfMeasurements = readWordAsInt(bytes, 2)             // word 3
	r.FirstFrame = readWordAsInt(bytes, 3)                   // word 4
	r.LastFrame = readWordAsInt(bytes, 4)                    // word 5
	r.InterpolationGap = readWordAsInt(bytes, 5)             // word 6
	r.ScaleFactor = read2WordsAsFloat(bytes, 6)              // word 7 - 8
	r.DataStart = readWordAsInt(bytes, 8)                    // word 9
	r.NrOfSamples = readWordAsInt(bytes, 9)                  // word 10
	r.FrameRate = read2WordsAsFloat(bytes, 10)               // word 11 - 12
	r.HasLabels = (readWordAsInt(bytes, 147) == 12345)       // word 148
	r.Uses4CharLabels = (readWordAsInt(bytes, 149) == 12345) // word 150
	r.EventLabels = readEventLabels(bytes, 198, 233)         // word 199 - 234

	return
}

func parseParameterBlock(bytes []byte, info *C3DInfo) {
	nrOfBytes := len(bytes)
	byteIndex := 0
	b := make([]byte, 2, 2)

	var groups []C3DGroup
	var parameters []C3DParameter

	fmt.Println("First 20 bytes:")
	fmt.Println(bytes[0:20])

	for byteIndex < nrOfBytes {
		nrOfCharactersInName := int(uint8(bytes[byteIndex]))
		byteIndex++

		// in case there are filling zero-bytes, e.g. at the end of the parameter block
		if nrOfCharactersInName == 0 {
			continue
		}

		groupID := int(int8(bytes[byteIndex]))
		byteIndex++
		fmt.Println("Found group id ", groupID)
		locked := false
		if nrOfCharactersInName < 0 {
			locked = true
			nrOfCharactersInName = -nrOfCharactersInName
		}

		name := string(bytes[byteIndex : byteIndex+nrOfCharactersInName])
		byteIndex += nrOfCharactersInName
		fmt.Println(fmt.Sprintf("Found name \"%s\"", name))

		// reading byte offset to next block
		b[0] = bytes[byteIndex]
		byteIndex++
		b[1] = bytes[byteIndex]
		byteIndex++
		bits := binary.LittleEndian.Uint16(b)
		nextBlockStartsAt := int(int16(bits)) + byteIndex - 2
		fmt.Println("Next block starts at", nextBlockStartsAt)

		if groupID > 0 { // we have a parameter
			fmt.Println("Working on Parameter")

			fmt.Println("byte:", bytes[byteIndex])
			nrOfBytes := int(int8(bytes[byteIndex]))
			byteIndex++
			fmt.Println("Nr of Bytes", nrOfBytes)

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
				// panic(fmt.Sprintf("Wrong number of bytes detected in parsing of parameters: %d", nrOfBytes))
				dataType = UNKNOWN
			}

			fmt.Println("Data Type", dataType)

			nrOfDimensions := int(int8(bytes[byteIndex]))
			byteIndex++
			fmt.Println("Nr of Dimensions", nrOfDimensions)

			if nrOfDimensions < 0 || nrOfDimensions > 7 {
				panic(fmt.Sprintf("Number of dimensions must be in [0,7] but is %d", nrOfDimensions))
			}
			dimensions := make([]int, nrOfDimensions, nrOfDimensions)
			dataSize := 1
			for i := 0; i < int(nrOfDimensions); i++ {
				dimensions[i] = int(int8(bytes[byteIndex]))
				dataSize *= int(dimensions[i])
				byteIndex++
			}
			data := bytes[byteIndex : byteIndex+dataSize]
			byteIndex += dataSize + 1

			descriptionLength := int(uint8(bytes[byteIndex]))
			parameterDescription := ""
			if descriptionLength > 0 {
				parameterDescription = string(bytes[byteIndex : byteIndex+descriptionLength])
			}
			p := C3DParameter{GroupID: groupID,
				Name:        name,
				Description: parameterDescription,
				DataType:    dataType,
				Data:        data,
				Locked:      locked}
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
		byteIndex = nextBlockStartsAt
	}
	info.Groups = groups
	info.Parameters = parameters
}

func readParameters(io *bufio.Reader, info *C3DInfo) {
	bytes := make([]byte, 4, 4)
	_, err := io.Read(bytes)
	check(err)

	nrOfParameterBlocks := int(uint8(bytes[3]))

	fmt.Println(fmt.Sprintf("Nr of Parameter Blocks %d", nrOfParameterBlocks))

	bytes = make([]byte, nrOfParameterBlocks*512, nrOfParameterBlocks*512)
	_, err = io.Read(bytes)
	check(err)

	parseParameterBlock(bytes, info)
}

func ReadC3D(filename string, eta bool) {
	f, err := os.Open(filename)
	defer f.Close()
	check(err)

	stats, statsErr := f.Stat()
	check(statsErr)

	var size int = int(stats.Size())
	var bar *pb.ProgressBar

	if eta == true {
		bar = pb.StartNew(size)
	}

	bufr := bufio.NewReader(f)
	header := readHeader(bufr)
	fmt.Println(header)

	info := C3DInfo{Parameters: nil, Groups: nil}

	readParameters(bufr, &info)
	fmt.Println(info)

	// for i := 0; i < size; i++ {
	// v := readWord(bufr)

	// fmt.Println(v)

	// if eta == true {
	// bar.Increment()
	// }
	// }

	if eta == true {
		bar.Finish()
	}
}
