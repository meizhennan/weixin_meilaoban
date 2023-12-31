package material

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/silenceper/wechat/v2/util"
)

const (
	addNewsURL          = "https://api.weixin.qq.com/cgi-bin/material/add_news"
	updateNewsURL       = "https://api.weixin.qq.com/cgi-bin/material/update_news"
	addMaterialURL      = "https://api.weixin.qq.com/cgi-bin/material/add_material"
	delMaterialURL      = "https://api.weixin.qq.com/cgi-bin/material/del_material"
	getMaterialURL      = "https://api.weixin.qq.com/cgi-bin/material/get_material"
	getMaterialCountURL = "https://api.weixin.qq.com/cgi-bin/material/get_materialcount"
	batchGetMaterialURL = "https://api.weixin.qq.com/cgi-bin/material/batchget_material"
)

//PermanentMaterialType 永久素材类型
type PermanentMaterialType string

const (
	//PermanentMaterialTypeImage 永久素材图片类型（image）
	PermanentMaterialTypeImage PermanentMaterialType = "image"
	//PermanentMaterialTypeVideo 永久素材视频类型（video）
	PermanentMaterialTypeVideo PermanentMaterialType = "video"
	//PermanentMaterialTypeVoice 永久素材语音类型 （voice）
	PermanentMaterialTypeVoice PermanentMaterialType = "voice"
	//PermanentMaterialTypeNews 永久素材图文类型（news）
	PermanentMaterialTypeNews PermanentMaterialType = "news"
)

//Material 素材管理
type Material struct {
}

//NewMaterial init
func NewMaterial() *Material {
	material := new(Material)
	return material
}

//Article 永久图文素材
type Article struct {
	Title            string `json:"title"`
	ThumbMediaID     string `json:"thumb_media_id"`
	ThumbURL         string `json:"thumb_url"`
	Author           string `json:"author"`
	Digest           string `json:"digest"`
	ShowCoverPic     int    `json:"show_cover_pic"`
	Content          string `json:"content"`
	ContentSourceURL string `json:"content_source_url"`
	URL              string `json:"url"`
	DownURL          string `json:"down_url"`
}

//// GetNews 获取/下载永久素材
//func (material *Material) GetNews(id string) ([]*Article, error) {
//
//	uri := fmt.Sprintf("%s", getMaterialURL)
//
//	var req struct {
//		MediaID string `json:"media_id"`
//	}
//	req.MediaID = id
//	responseBytes, err := util.PostJSON(uri, req)
//	if err != nil {
//		return nil, err
//	}
//
//	var res struct {
//		NewsItem []*Article `json:"news_item"`
//	}
//
//	err = util.DecodeWithCustomerStruct(responseBytes, &res, "GetNews")
//	if err != nil {
//		return nil, err
//	}
//
//	return res.NewsItem, nil
//}

//reqArticles 永久性图文素材请求信息
type reqArticles struct {
	Articles []*Article `json:"articles"`
}

//resArticles 永久性图文素材返回结果
type resArticles struct {
	util.CommonError

	MediaID string `json:"media_id"`
}

//AddNews 新增永久图文素材
func (material *Material) AddNews(articles []*Article) (mediaID string, err error) {
	req := &reqArticles{articles}

	uri := fmt.Sprintf("%s", addNewsURL)
	responseBytes, err := util.PostJSON(uri, req)
	if err != nil {
		return
	}
	var res resArticles
	err = json.Unmarshal(responseBytes, &res)
	if err != nil {
		return
	}
	mediaID = res.MediaID
	return
}

//reqUpdateArticle 更新永久性图文素材请求信息
type reqUpdateArticle struct {
	MediaID  string   `json:"media_id"`
	Index    int64    `json:"index"`
	Articles *Article `json:"articles"`
}

// UpdateNews 更新永久图文素材
func (material *Material) UpdateNews(article *Article, mediaID string, index int64) (err error) {
	req := &reqUpdateArticle{mediaID, index, article}

	uri := fmt.Sprintf("%s", updateNewsURL)
	var response []byte
	response, err = util.PostJSON(uri, req)
	if err != nil {
		return
	}
	return util.DecodeWithCommonError(response, "UpdateNews")
}

//resAddMaterial 永久性素材上传返回的结果
type resAddMaterial struct {
	util.CommonError

	MediaID string `json:"media_id"`
	URL     string `json:"url"`
}

//AddMaterial 上传永久性素材（处理视频需要单独上传）
func (material *Material) AddMaterial(mediaType PermanentMaterialType, filename string) (mediaID string, url string, err error) {
	if mediaType == PermanentMaterialTypeVideo {
		err = errors.New("永久视频素材上传使用 AddVideo 方法")
		return
	}

	uri := fmt.Sprintf("%s?type=%s", addMaterialURL, mediaType)
	var response []byte
	response, err = util.PostFile("media", filename, uri)
	if err != nil {
		return
	}
	var resMaterial resAddMaterial
	err = util.DecodeWithError(response, &resMaterial, "AddMaterial")
	if err != nil {
		return
	}

	mediaID = resMaterial.MediaID
	url = resMaterial.URL
	return
}

type reqVideo struct {
	Title        string `json:"title"`
	Introduction string `json:"introduction"`
}

//AddVideo 永久视频素材文件上传
func (material *Material) AddVideo(filename, title, introduction string) (mediaID string, url string, err error) {

	uri := fmt.Sprintf("%s?type=video", addMaterialURL)

	videoDesc := &reqVideo{
		Title:        title,
		Introduction: introduction,
	}
	var fieldValue []byte
	fieldValue, err = json.Marshal(videoDesc)
	if err != nil {
		return
	}

	fields := []util.MultipartFormField{
		{
			IsFile:    true,
			Fieldname: "media",
			Filename:  filename,
		},
		{
			IsFile:    false,
			Fieldname: "description",
			Value:     fieldValue,
		},
	}

	var response []byte
	response, err = util.PostMultipartForm(fields, uri)
	if err != nil {
		return
	}

	var resMaterial resAddMaterial

	err = util.DecodeWithError(response, &resMaterial, "AddVideo")
	if err != nil {
		return
	}

	mediaID = resMaterial.MediaID
	url = resMaterial.URL
	return
}

type reqDeleteMaterial struct {
	MediaID string `json:"media_id"`
}

//DeleteMaterial 删除永久素材
func (material *Material) DeleteMaterial(mediaID string) error {

	uri := fmt.Sprintf("%s", delMaterialURL)
	response, err := util.PostJSON(uri, reqDeleteMaterial{mediaID})
	if err != nil {
		return err
	}

	return util.DecodeWithCommonError(response, "DeleteMaterial")
}

//ArticleList 永久素材列表
type ArticleList struct {
	util.CommonError
	TotalCount int64             `json:"total_count"`
	ItemCount  int64             `json:"item_count"`
	Item       []ArticleListItem `json:"item"`
}

//ArticleListItem 用于ArticleList的item节点
type ArticleListItem struct {
	MediaID    string             `json:"media_id"`
	Content    ArticleListContent `json:"content"`
	Name       string             `json:"name"`
	URL        string             `json:"url"`
	UpdateTime int64              `json:"update_time"`
}

//ArticleListContent 用于ArticleListItem的content节点
type ArticleListContent struct {
	NewsItem   []Article `json:"news_item"`
	UpdateTime int64     `json:"update_time"`
	CreateTime int64     `json:"create_time"`
}

//reqBatchGetMaterial BatchGetMaterial请求参数
type reqBatchGetMaterial struct {
	Type   PermanentMaterialType `json:"type"`
	Count  int64                 `json:"count"`
	Offset int64                 `json:"offset"`
}

// BatchGetMaterial 批量获取永久素材
//reference:https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Get_materials_list.html
func (material *Material) BatchGetMaterial(permanentMaterialType PermanentMaterialType, offset, count int64) (list ArticleList, err error) {

	uri := fmt.Sprintf("%s", batchGetMaterialURL)

	req := reqBatchGetMaterial{
		Type:   permanentMaterialType,
		Offset: offset,
		Count:  count,
	}

	var response []byte
	response, err = util.PostJSON(uri, req)
	if err != nil {
		return
	}

	err = util.DecodeWithError(response, &list, "BatchGetMaterial")
	return
}

// ResMaterialCount 素材总数
type ResMaterialCount struct {
	util.CommonError
	VoiceCount int64 `json:"voice_count"` // 语音总数量
	VideoCount int64 `json:"video_count"` // 视频总数量
	ImageCount int64 `json:"image_count"` // 图片总数量
	NewsCount  int64 `json:"news_count"`  // 图文总数量
}

// GetMaterialCount 获取素材总数.
func (material *Material) GetMaterialCount() (res ResMaterialCount, err error) {

	uri := fmt.Sprintf("%s", getMaterialCountURL)
	var response []byte
	response, err = util.HTTPGet(uri)
	if err != nil {
		return
	}
	err = util.DecodeWithError(response, &res, "GetMaterialCount")
	return
}
