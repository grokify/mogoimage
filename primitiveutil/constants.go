package primitiveutil

/*
https://github.com/fogleman/primitive
mode: 0=combo, 1=triangle, 2=rect, 3=ellipse, 4=circle, 5=rotatedrect, 6=beziers, 7=rotatedellipse, 8=polygon
*/

type Mode int

const (
	ModeCombo Mode = iota
	ModeTriangle
	ModeRectangle
	ModeEllipose
	ModeCircle
	ModeRotatedRectangle
	ModeBeziers
	ModeRotatedEllipse
	ModePolygon
)
