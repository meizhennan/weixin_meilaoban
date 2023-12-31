package helpers

import (
	"github.com/vicanso/go-charts/v2"
	"io/ioutil"
	//"io/ioutil"
)

func DrawDoubleYaxis(stock string, title string,
	xAxisOption []string,
	legendOption []string,
	seriesList []charts.Series,
	fileName string) (string, error) {

	//加载中文字体文件
	//buf, err := ioutil.ReadFile("./src/static/NotoSansCJKsc-VF.ttf")
	//if err != nil {
	//	panic(err)
	//}
	//err = charts.InstallFont("noto", buf)
	//if err != nil {
	//	panic(err)
	//}
	//font, _ := charts.GetFont("noto")
	//charts.SetDefaultFont(font)

	chartOption := charts.ChartOption{
		Title: charts.TitleOption{
			Text:     title,
			FontSize: 16,
			Left:     "200",
		},
		XAxis:  charts.NewXAxisOption(xAxisOption),
		Legend: charts.NewLegendOption(legendOption, "500"),
		YAxisOptions: []charts.YAxisOption{
			{
				Formatter: "{value}亿元",
				Color: charts.Color{
					R: 84,
					G: 112,
					B: 198,
					A: 255,
				},
			},
			{
				Formatter: "{value}%",
				Color: charts.Color{
					R: 84,
					G: 112,
					B: 198,
					A: 255,
				},
			},
		},
		SeriesList: seriesList,
	}
	p, err := charts.Render(chartOption)
	if err != nil {
		panic(err)
	}

	buf, err := p.Bytes()
	if err != nil {
		panic(err)
	}

	filePath, err := WriteFile(buf, stock, fileName)
	if err != nil {
		panic(err)
	}
	return filePath, err
}

func MultipleLine(stock, title string,
	legendOption []string,
	xAxisOption []string,
	seriesList []charts.Series,
	fileName string) {
	//加载中文字体文件
	buf, err := ioutil.ReadFile("./src/static/NotoSansCJKsc-VF.ttf")
	if err != nil {
		panic(err)
	}
	err = charts.InstallFont("noto", buf)
	if err != nil {
		panic(err)
	}
	font, _ := charts.GetFont("noto")
	charts.SetDefaultFont(font)

	chartOption := charts.ChartOption{
		Title: charts.TitleOption{
			Text:     title,
			FontSize: 16,
			Left:     "100",
		},
		Legend: charts.NewLegendOption(legendOption, "500"),
		XAxis:  charts.NewXAxisOption(xAxisOption),
		//YAxisOptions: []charts.YAxisOption{
		//	{
		//		Formatter: "{value}%",
		//		Color: charts.Color{
		//			R: 84,
		//			G: 112,
		//			B: 198,
		//			A: 255,
		//		},
		//	},
		//},
		SeriesList: seriesList,
	}
	p, err := charts.Render(chartOption)
	if err != nil {
		panic(err)
	}

	buf, err = p.Bytes()
	if err != nil {
		panic(err)
	}

	_, err = WriteFile(buf, stock, fileName)
	if err != nil {
		panic(err)
	}

}
