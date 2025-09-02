package image

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

// 图片链接
type ImageLinks struct {
	URL              string `json:"url"`
	HTML             string `json:"html"`
	BBCode           string `json:"bbcode"`
	Markdown         string `json:"markdown"`
	MarkdownWithLink string `json:"markdown_with_link"`
	ThumbnailURL     string `json:"thumbnail_url"`
	DeleteURL        string `json:"delete_url"`
}

// 图片数据
type ImageData struct {
	Key        string     `json:"key"`
	Name       string     `json:"name"`
	Pathname   string     `json:"pathname"`
	OriginName string     `json:"origin_name"`
	Size       string     `json:"size"`
	Mimetype   string     `json:"mimetype"`
	Extension  string     `json:"extension"`
	MD5        string     `json:"md5"`
	SHA1       string     `json:"sha1"`
	Links      ImageLinks `json:"links"`
}

type APIResponse struct {
	Status  bool      `json:"status"`
	Message string    `json:"message"`
	Data    ImageData `json:"data"`
}

type ImageHelper struct {
	token   string
	baseUrl string
	albumId int
	client  *http.Client
}

// 创建一个新的 ImageHelper 实例
func NewImageHelper() *ImageHelper {
	return &ImageHelper{
		token:   "1550|K0tdQUigggjMnUxru5jZb6lLMsipUFa6fVwAHu0q",
		baseUrl: "https://picui.cn/api/v1",
		albumId: 1591,
		client:  &http.Client{},
	}
}

// 上传文件
func (h *ImageHelper) UploadImage(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(fileHeader.Filename))
	if err != nil {
		return "", fmt.Errorf("创建文件表单字段失败: %w", err)
	}

	if _, err = io.Copy(part, file); err != nil {
		return "", fmt.Errorf("拷贝文件内容失败: %w", err)
	}

	if h.albumId != 0 {
		if err = writer.WriteField("album_id", fmt.Sprintf("%d", h.albumId)); err != nil {
			return "", fmt.Errorf("写入 album_id 字段失败: %w", err)
		}
	}

	if err = writer.Close(); err != nil {
		return "", fmt.Errorf("关闭 multipart writer 失败: %w", err)
	}

	req, err := http.NewRequest("POST", h.baseUrl+"/upload", body)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+h.token)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := h.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应体失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		// 尝试解析错误响应
		var errorResp struct {
			Message string `json:"message"`
		}
		if json.Unmarshal(respBody, &errorResp) == nil && errorResp.Message != "" {
			return "", fmt.Errorf("API 请求失败，状态码: %d, 错误: %s",
				resp.StatusCode, errorResp.Message)
		}

		return "", fmt.Errorf("API 请求失败，状态码: %d, 响应: %s",
			resp.StatusCode, string(respBody))
	}

	var apiResponse APIResponse
	if err := json.Unmarshal(respBody, &apiResponse); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	return apiResponse.Data.Links.URL, nil
}
