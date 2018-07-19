package main

type Target struct {
	Circle
	Enabled bool
}

func NewTarget(x, y float64) *Target {
	t := &Target{
		Circle:  NewCircle(x, y, 5),
		Enabled: true,
	}

	t.Circle.Color = colorGreen

	return t
}
