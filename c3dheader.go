package goc3d

import "fmt"

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
	ScaleFactor      float32
	FrameRate        float32
	EventLabels      []string
	UsesInteger      bool
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

	if h.UsesInteger == true {
		str = fmt.Sprintf("%sData is given in ints  = true\n", str)
	} else {
		str = fmt.Sprintf("%sData is given in ints  = false\n", str)
	}

	return str
}
