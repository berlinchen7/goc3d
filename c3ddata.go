package goc3d

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
