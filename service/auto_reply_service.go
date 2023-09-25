package service

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
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
	//impPath,err := searchDataAndCreateImg(receiveMsg)

	materialObj := material.NewMaterial()
	mediaId, mediaUrl, err := materialObj.AddMaterial(material.PermanentMaterialTypeImage, "./static/cost_ratio.png")
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
		Image:        *media,
	}

	msg, err := json.Marshal(res)
	if err != nil {
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}

func searchDataAndCreateImg(msg string) (string, error) {
	return "", nil
}
