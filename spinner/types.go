package spinner

//Sequence enum for choosing spinner
type Sequence []string

//Circle1 spinner sequence that is made of ◐◓◑◒
func Circle1() Sequence {
	return []string{"◐", "◓", "◑", "◒"}
}

//Pole1 spinner sequence that is made of ▁▃▄▅▆▇█▇▆▅▄▃
func Pole1() Sequence {
	return []string{"▁", "▃", "▄", "▅", "▆", "▇", "█", "▇", "▆", "▅", "▄", "▃"}
}

//Pole2 spinner sequence that is made of ▉▊▋▌▍▎▏▎▍▌▋▊▉
func Pole2() Sequence {
	return []string{"▉", "▊", "▌", "▍", "▎", "▏", "▎", "▍", "▌", "▋", "▊", "▉"}
}
