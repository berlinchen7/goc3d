package goc3d

import "fmt"

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
