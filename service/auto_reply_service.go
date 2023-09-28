package service

import (
	"encoding/json"
	"github.com/vicanso/go-charts/v2"
	"log"
	"net/http"
	"os"
	"time"
	"wxcloudrun-golang/helpers"
	"wxcloudrun-golang/material"
)

type AutoReplyRequest struct {
	ToUserName   string `json:"ToUserName"`
	FromUserName string `json:"FromUserName"`
	CreateTime   int64  `json:"CreateTime"`
	MsgType      string `json:"MsgType"`
	Content      string `json:"Content"`
	MsgId        int64  `json:"MsgId"`
}

type Media struct {
	MediaId string `json:"MediaId"`
}
type AutoReplyResponse struct {
	ToUserName   string `json:"ToUserName"`
	FromUserName string `json:"FromUserName"`
	CreateTime   int64  `json:"CreateTime"`
	MsgType      string `json:"MsgType"`
	Content      string `json:"Content"`
	Image        Media  `json:"Image"`
	Voice        Media  `json:"Voice"`
}

func AutoReplyHandler(w http.ResponseWriter, r *http.Request) {
	var autoReplyRequest AutoReplyRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&autoReplyRequest)
	if err != nil {
		return
	}

	//根据消息内容返回不同的类型
	receiveMsg := autoReplyRequest.Content
	log.Printf("receiveMsg:%s", receiveMsg)
	//receiveMsgType := autoReplyRequest.MsgType

	//根据内容查询对应的数据
	impPath, err := searchDataAndCreateImg(receiveMsg)
	log.Printf("filePath [%s]", impPath)

	materialObj := material.NewMaterial()
	mediaId, mediaUrl, err := materialObj.AddMaterial(material.PermanentMaterialTypeImage, impPath)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("meidaId: %s , mediaUrl: %s ", mediaId, mediaUrl)

	media := &Media{
		MediaId: mediaId,
	}
	res := &AutoReplyResponse{
		ToUserName:   autoReplyRequest.FromUserName,
		FromUserName: autoReplyRequest.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "image",
		//Content:      "测试回复",
		Image: *media,
	}

	msg, err := json.Marshal(res)
	if err != nil {
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}

func searchDataAndCreateImg(msg string) (string, error) {
	// 根据用户发送消息查找是哪个股票
	stock := "000519" //股票代码
	path := helpers.ImageBasePath + stock
	//清空 /tmp/{stock}/*.png 所有图片
	err := deleteFiles(path)
	if err != nil {
		return "", err
	}

	xAxisOption := []string{
		"2016",
		"2017",
		"2018",
		"2019",
		"2020",
		"2021",
		"2022",
	}
	legendOption := []string{
		"营业收入",
		"费用率",
		"毛利率",
		"营业利润率",
	}
	seriesList := []charts.Series{
		//营业收入
		{
			Type: charts.ChartTypeBar,
			Data: charts.NewSeriesDataFromValues([]float64{
				401.55,
				610.63,
				771.99,
				888.54,
				979.93,
				1094.64,
				1275.54,
			}),
		},
		//费用率:销售、管理、财务、研发 综合占营业总收入的比例 （财务费用为正,则相加，为负不进行扣减，计算的比较保守）
		{
			Type: charts.ChartTypeLine,
			Data: charts.NewSeriesDataFromValues([]float64{
				14.61,
				12.62,
				10.23,
				10.63,
				9.53,
				10.22,
				9.65,
			}),
			AxisIndex: 1,
		},
		//毛利率   (营业收入-营业成本)/营业收入
		{
			Type: charts.ChartTypeLine,
			Data: charts.NewSeriesDataFromValues([]float64{
				91.51,
				90.27,
				91.55,
				91.64,
				91.68,
				91.79,
				92.09,
			}),
			AxisIndex: 1,
		},
		//营业利润率   营业利润/营业收入
		{
			Type: charts.ChartTypeLine,
			Data: charts.NewSeriesDataFromValues([]float64{
				41.63,
				44.35,
				45.60,
				46.37,
				47.65,
				47.92,
				49.17,
			}),
			AxisIndex: 1,
		},
	}
	fileName := "income_from_operation.png"
	filePath, err := helpers.DrawDoubleYaxis(stock, "营业利润趋势", xAxisOption, legendOption, seriesList, fileName)

	imagePaths := []string{filePath, filePath, filePath, filePath}

	//合并图片
	outputPath, err := helpers.MergeImage(path, imagePaths)
	if err != nil {
		return "", err
	}
	return outputPath, err
}

func deleteFiles(dir string) error {
	_, err := os.Stat(dir)
	if err != nil { //不存在dir 直接返回nil
		return nil
	}

	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(dir + "/" + name)
		if err != nil {
			return err
		}
	}
	return nil
}
