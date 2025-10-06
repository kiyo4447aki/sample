package handlers

import (
	"api/config"
	"api/records"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//ListRecordHandler
//正常系；フィルター不使用
func TestListRecordsHandler_Success_NoDate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	base := t.TempDir()
	cam := filepath.Join(base, "cam")
	err := os.Mkdir(cam, 0755)
	require.NoError(t, err)
	name := "r.mp4"
	path := filepath.Join(cam, name)
	err = os.WriteFile(path, []byte("data"), 0644)
	require.NoError(t, err)
	mod := time.Date(2025, 7, 2, 12, 0, 0, 0, time.UTC)
	err = os.Chtimes(path, mod, mod)
	require.NoError(t, err)

	cfg := &config.Config{BasePath: base, BaseUrl: "http://base"}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/?camera=cam", nil)

	handler := ListRecordsHandler(cfg)
	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp records.RecordsResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Len(t, resp.Recordings, 1)
	assert.Equal(t, 1, resp.Count)
	rec := resp.Recordings[0]
	assert.Equal(t, name, rec.Name)
	assert.Equal(t, "http://base/cam/r.mp4", rec.Url)
	assert.True(t, rec.Timestamp.Equal(mod.UTC()))
}

// 正常系：フィルタあり、マッチするレコードあり
func TestListRecordsHandler_Success_WithDateFilterMatch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	base := t.TempDir()
	cam := filepath.Join(base, "cam")
	err := os.Mkdir(cam, 0755)
	require.NoError(t, err)
	name := "match.mp4"
	path := filepath.Join(cam, name)
	err = os.WriteFile(path, []byte("x"), 0644)
	require.NoError(t, err)
	filter := time.Date(2025, 7, 2, 0, 0, 0, 0, time.UTC)
	err = os.Chtimes(path, filter, filter)
	require.NoError(t, err)


	cfg := &config.Config{BasePath: base, BaseUrl: "http://base"}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/?camera=cam&date=2025-07-02", nil)

	handler := ListRecordsHandler(cfg)
	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp records.RecordsResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Len(t, resp.Recordings, 1)
	assert.Equal(t, 1, resp.Count)
	rec := resp.Recordings[0]
	assert.Equal(t, name, rec.Name)
	assert.Equal(t, "http://base/cam/match.mp4", rec.Url)
	assert.True(t, rec.Timestamp.Equal(filter.UTC()))
}

//正常系：フィルタあり、マッチするレコードなし
func TestListRecordsHandler_Success_WithDateFilterMismatch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	base := t.TempDir()
	cam := filepath.Join(base, "cam")
	err := os.Mkdir(cam, 0755)
	require.NoError(t, err)
	name := "r2.mp4"
	path := filepath.Join(cam, name)
	err = os.WriteFile(path, []byte(""), 0644)
	require.NoError(t, err)
	mod := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
	err = os.Chtimes(path, mod, mod)
	require.NoError(t, err)

	cfg := &config.Config{BasePath: base, BaseUrl: "http://base"}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/?camera=cam&date=2025-07-02", nil)

	handler := ListRecordsHandler(cfg)
	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp records.RecordsResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Len(t, resp.Recordings, 0)
	assert.Equal(t, 0, resp.Count)
}

// 異常系：camera パラメータが空文字の場合
func TestListRecordsHandler_EmptyCameraParam(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{BasePath: t.TempDir(), BaseUrl: "http://base"}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/?camera=", nil)

	handler := ListRecordsHandler(cfg)
	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var body map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "failed", body["status"])
	assert.Equal(t, "invalid camera name", body["error"])
}

//異常系：カメラ名に上位ディレクトリのパスが含まれているとき、エラーを返す
func TestListRecordsHandler_InvalidCameraName(t *testing.T){
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{BasePath: t.TempDir(), BaseUrl: "http://base"}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/?camera=../secret", nil)

	handler := ListRecordsHandler(cfg)
	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var body map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "failed", body["status"])
	assert.Equal(t, "invalid camera name", body["error"])
}

//異常系：フィルターの日付が不正な値のとき、エラーが返る
func TestListRecordsHandler_InvalidDateFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{BasePath: t.TempDir(), BaseUrl: "http://base"}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/?camera=cam&date=2025-13-01", nil)

	handler := ListRecordsHandler(cfg)
	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var body map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "failed", body["status"])
	assert.Equal(t, "invalid date format", body["error"])
}

//異常系；カメラが存在しない時、エラーが返る
func TestListRecordsHandler_CameraNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{BasePath: t.TempDir(), BaseUrl: "http://base"}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/?camera=noExistCamera", nil)

	handler := ListRecordsHandler(cfg)
	handler(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var body map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "failed", body["status"])
	assert.Equal(t, "camera not found", body["error"])
}

//異常系；レコード情報の取得に失敗した時、エラーが返る
func TestListRecordsHandler_GetRecordsFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// ファイルのパスをディレクトリパスとして渡す
	base := t.TempDir()
	file := filepath.Join(base, "camfile")
	err := os.WriteFile(file, []byte(""), 0644)
	require.NoError(t, err)
	cfg := &config.Config{BasePath: base, BaseUrl: "http://base"}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/?camera=camfile", nil)

	handler := ListRecordsHandler(cfg)
	handler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var body map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, "failed", body["status"])
	assert.Equal(t, "failed to get records", body["error"])
}


