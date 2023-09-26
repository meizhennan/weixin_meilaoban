package helpers

import (
	"github.com/vicanso/go-charts/v2"
	"io/ioutil"
)

func DrawDoubleYaxis(title string,
	xAxisOption []string,
	legendOption []string,
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

	buf, err = p.Bytes()
	if err != nil {
		panic(err)
	}

	err = WriteFile(buf, fileName)
	if err != nil {
		panic(err)
	}
}
