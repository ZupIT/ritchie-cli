package template

type Template []string

var (
	Default = Template{"⣷", "⣯", "⣟", "⡿", "⢿", "⣻", "⣽", "⣾"}
	Grow    = Template{"▏", "▎", "▍", "▌", "▋", "▊", "▉", "▊", "▋", "▌", "▍", "▎"}
	Arrow   = Template{"⬆️ ", "↗️ ", "➡️ ", "↘️ ", "⬇️ ", "↙️ ", "⬅️ ", "↖️ "}
	Clock   = Template{"🕛 ", "🕐 ", "🕑 ", "🕒 ", "🕓 ", "🕔 ", "🕕 ", "🕖 ", "🕗 ", "🕘 ", "🕙 ", "🕚 "}
	Point   = Template{"∙∙∙", "●∙∙", "∙●∙", "∙∙●", "∙∙∙"}
)
