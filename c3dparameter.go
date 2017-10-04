package goc3d

import "fmt"

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
