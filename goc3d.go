package goc3d

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math"
	"os"

	pb "gopkg.in/cheggaaa/pb.v1"
)

type C3DHeader struct {
	Valid            bool
	HasLabels        bool
	Uses4CharLabels  bool
	NrOfTrajectories int64
	NrOfMeasurements int64
	FirstFrame       int64
	LastFrame        int64
	InterpolationGap int64
	DataStart        int64
	NrOfSamples      int64
	ScaleFactor      float64
	FrameRate        float64
}

func (h C3DHeader) String() string {
	str := ""

	if h.Valid == true {
		str = fmt.Sprintf("Valid C3D                = true\n")
	} else {
		str = fmt.Sprintf("Valid C3D                = false\n")
	}

	str = fmt.Sprintf("Valid C3D                = false\n")
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

func read2BytesAsInt(io *bufio.Reader) (r int64) {
	bytes := make([]byte, 2, 2)
	_, aerr := io.Read(bytes)
	check(aerr)
	bits := binary.LittleEndian.Uint16(bytes)
	r = int64(int16(bits))
	return
}

func read4BytesAsFloat(io *bufio.Reader) (r float64) {
	bytes := make([]byte, 4, 4)
	_, aerr := io.Read(bytes)
	check(aerr)
	bits := binary.LittleEndian.Uint32(bytes)
	r = float64(math.Float32frombits(bits))
	return r
}

func readHeader(io *bufio.Reader) (r C3DHeader) {
	a, aerr := io.ReadByte()
	check(aerr)

	if a != 2 {
		fmt.Println("Not a C3D file")
		fmt.Println(fmt.Sprintf("First byte is not 2 but %d", a))
		os.Exit(-1)
	}

	b, berr := io.ReadByte()
	check(berr)

	if b != 80 {
		fmt.Println("Not a C3D file")
		fmt.Println(fmt.Sprintf("Second byte is not 80 but %d", b))
		os.Exit(-1)
	}

	r.Valid = true

	r.NrOfTrajectories = read2BytesAsInt(io) // word 2
	r.NrOfMeasurements = read2BytesAsInt(io) // word 3
	r.FirstFrame = read2BytesAsInt(io)       // word 4
	r.LastFrame = read2BytesAsInt(io)        // word 5
	r.InterpolationGap = read2BytesAsInt(io) // word 6
	r.ScaleFactor = read4BytesAsFloat(io)    // word 7 - 8
	r.DataStart = read2BytesAsInt(io)        // word 9
	r.NrOfSamples = read2BytesAsInt(io)      // word 10
	r.FrameRate = read4BytesAsFloat(io)      // word 11 - 12

	bytes := make([]byte, 135, 135)
	_, err := io.Read(bytes)
	check(err)

	r.HasLabels = (read2BytesAsInt(io) == 12345)

	bytes = make([]byte, 2, 2)
	_, err = io.Read(bytes)
	check(err)

	r.Uses4CharLabels = (read2BytesAsInt(io) == 12345)

	return
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
