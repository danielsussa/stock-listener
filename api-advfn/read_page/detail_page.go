package read_page

import (
	"strconv"
	"strings"

	"github.com/go-resty/resty"
)

func ReadDetailPage(opt *option) {
	res, err := resty.R().Get("https:" + opt.Url)
	if err != nil {
		panic(err)
	}
	getOptionDetail(string(res.Body()), opt)

}

func getOptionDetail(s string, opt *option) {

	s = strings.Replace(s, "<b>", "", -1)
	s = strings.Replace(s, "</b>", "", -1)
	s = strings.Replace(s, "</span>", "", -1)
	s = strings.Replace(s, ",", ".", -1)
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\t", "", -1)

	tables := strings.Split(s, "<table>")
	tables = tables[1 : len(tables)-1]

	for _, table := range tables {
		table = strings.Split(table, "</table>")[0]

		mapHeader := make(map[int]string)
		headers := strings.Split(table, "<th")
		for i, th := range headers[1 : len(headers)-1] {
			if strings.Contains(th, "Código da Opção") {
				mapHeader[i] = "code_opcao"
			}
			if strings.Contains(th, "Preço de Exercício") {
				mapHeader[i] = "strike"
			}
			if strings.Contains(th, "Último Preço") {
				mapHeader[i] = "last_price"
			}
			if strings.Contains(th, "Dias até Vencimento") {
				mapHeader[i] = "exp_days"
			}
			if strings.Contains(th, "Tipo de Negócio") {
				mapHeader[i] = "kind"
			}
			if strings.Contains(th, "Estilo de Opção") {
				mapHeader[i] = "style"
			}
			if strings.Contains(th, "Número de Negócios") {
				mapHeader[i] = "qtd_neg"
			}
			if strings.Contains(th, "Volume de Ações Negociadas") {
				mapHeader[i] = "vol_neg"
			}
		}

		content := strings.Split(table, "<td")
		for i, td := range content[1 : len(content)-1] {
			valSpl := strings.Split(td, ">")
			val := valSpl[len(valSpl)-2]
			val = strings.Split(val, "<")[0]
			if mapHeader[i] == "code_opcao" {
				opt.Name = val
			}
			if mapHeader[i] == "strike" {
				price, err := strconv.ParseFloat(val, 32)
				if err != nil {
					panic(err)
				}
				opt.Strike = price
			}
			if mapHeader[i] == "last_price" {
				price, err := strconv.ParseFloat(val, 32)
				if err != nil {
					continue
				}
				opt.Price = price
			}
			if mapHeader[i] == "exp_days" {
				val = strings.Replace(val, " dias", "", -1)
				days, err := strconv.ParseFloat(val, 32)
				if err != nil {
					panic(err)
				}
				opt.Expiration = days
			}
			if mapHeader[i] == "kind" {
				if strings.Contains(val, "Call") {
					opt.Kind = "C"
				}
				if strings.Contains(val, "Put") {
					opt.Kind = "P"
				}
			}
			if mapHeader[i] == "style" {
				if strings.Contains(val, "Americano") {
					opt.Style = "A"
				}
				if strings.Contains(val, "Europeu") {
					opt.Style = "E"
				}
			}
			if mapHeader[i] == "qtd_neg" {
				v, err := strconv.ParseFloat(val, 32)
				if err != nil {
					panic(err)
				}
				opt.QtdNegs = v
			}
			if mapHeader[i] == "vol_neg" {
				val = strings.Replace(val, ".", "", -1)
				v, err := strconv.ParseFloat(val, 32)
				if err != nil {
					panic(err)
				}
				opt.VolNegs = v
			}
		}

	}

	// {
	// 	spl := strings.Split(s, "quoteElementPiece")
	// 	for _, v := range spl {
	// 		if len(v) == 0 {
	// 			continue
	// 		}

	// 		if string(v[0:2]) == "10" {
	// 			pStr := v[strings.Index(v, ">")+1 : strings.Index(v, "<")]
	// 			pStr = strings.Replace(pStr, ",", ".", -1)
	// 			price, _ := strconv.ParseFloat(pStr, 32)
	// 			opt.Price = price
	// 		}
	// 	}
	// }

	// // Get expiration date
	// {
	// 	spl := strings.Split(s, "dias</td>")
	// 	txt := spl[0][len(spl[0])-5:]
	// 	txt = strings.TrimSpace(strings.Split(txt, ">")[1])
	// 	expDate, _ := strconv.ParseFloat(txt, 32)
	// 	opt.Expiration = expDate
	// }
}
