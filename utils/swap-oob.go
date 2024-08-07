package utils

type OobSwap struct {
	Extension string
	Mode      string
	Selector  string
}

func (o *OobSwap) Build() string {
	result := ""
	// atm it's hard to use extendsions in oob-swaps
	// as oob-swap attr is trating extension params as
	// a css selector
	// if o.Extension != "" {
	// 	result += o.Extension + ":"
	// }
	if o.Mode != "" {
		result += o.Mode
	} else {
		result += "innerHtml"
	}
	if o.Selector != "" {
		result += ":" + o.Selector
	}
	return result
}
