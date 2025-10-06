package records

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//GetRecordNames関数
//正常系：mp4のみ抽出してファイル名の配列を返す
func TestGetRecordNames_Success(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.mp4"), []byte(""), 0644)
	os.WriteFile(filepath.Join(dir, "b.txt"), []byte(""), 0644)
	os.WriteFile(filepath.Join(dir, "c.mp4"), []byte(""), 0644)

	names, err := GetRecordNames(dir)
	assert.NoError(t, err)
	assert.Equal(t, []string{"a.mp4", "c.mp4"}, names)
}

//異常系：存在しないディレクトリ（カメラ名）を指定したとき、適切なエラーが返る
func TestGetRecordNames_NotExist(t *testing.T) {
	// 存在しないディレクトリ
	names, err := GetRecordNames("nonexistent_dir")
	assert.Nil(t, names)
	assert.ErrorIs(t, err, ErrCameraNotFound)
}

//異常系：ディレクトリの読み取りエラー時に適切なエラーが返る
func TestGetRecordNames_OtherError(t *testing.T) {
	file := filepath.Join(t.TempDir(), "not_a_dir")
	os.WriteFile(file, []byte(""), 0644)

	names, err := GetRecordNames(file)
	assert.Nil(t, names)
	assert.ErrorIs(t, err, ErrGetRecordsFailed)
}

//GetRecordInfo関数
//正常系：フィルターなしで正常に動作する
func TestGetRecordInfo_NoFilter(t *testing.T) {
	dir := t.TempDir()
	name := "rec.mp4"
	path := filepath.Join(dir, name)
	data := []byte("hello")
	err := os.WriteFile(path, data, 0644)
	require.NoError(t, err)
	modTime := time.Date(2025, 7, 2, 15, 4, 5, 0, time.UTC)
	err = os.Chtimes(path, modTime, modTime)
	require.NoError(t, err)


	rec, err := GetRecordInfoWithDateFilter(
		name,
		dir,
		false,
		time.Time{},
		"camera1",
		"http://base",
	)
	assert.NoError(t, err)
	assert.Equal(t, name, rec.Name)
	assert.Equal(t, "http://base/camera1/rec.mp4", rec.Url)
	assert.True(t, rec.Timestamp.Equal(modTime.UTC()))
	assert.Equal(t, int64(len(data)), rec.Size)
}

//正常系：フィルターを使用したとき、正常に動作する
func TestGetRecordInfo_FilterMatch(t *testing.T) {
	dir := t.TempDir()
	name := "rec2.mp4"
	path := filepath.Join(dir, name)
	err := os.WriteFile(path, []byte("xx"), 0644)
	require.NoError(t, err)
	modTime := time.Date(2025, 7, 2, 10, 0, 0, 0, time.UTC)
	err = os.Chtimes(path, modTime, modTime)
	require.NoError(t, err)

	filterDate := modTime
	rec, err := GetRecordInfoWithDateFilter(
		name,
		dir,
		true,
		filterDate,
		"cam2",
		"http://base2/",
	)
	assert.NoError(t, err)
	assert.Equal(t, name, rec.Name)
	assert.Equal(t, "http://base2/cam2/rec2.mp4", rec.Url)
}

//異常系：フィルタと作成日時が一致しない時、ErrDateNotMatchエラーが返る
func TestGetRecordInfo_FilterMismatch(t *testing.T) {
	dir := t.TempDir()
	name := "rec3.mp4"
	path := filepath.Join(dir, name)
	err := os.WriteFile(path, []byte(""), 0644)
	require.NoError(t, err)
	modTime := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
	err = os.Chtimes(path, modTime, modTime)
	require.NoError(t, err)

	filterDate := time.Date(2025, 7, 2, 0, 0, 0, 0, time.UTC)
	_, err = GetRecordInfoWithDateFilter(
		name,
		dir,
		true,
		filterDate,
		"cam3",
		"http://b",
	)
	assert.ErrorIs(t, err, ErrDateNotMatch)
}

//異常系；ファイルが存在しないとき、ErrGetFileInfoが返る
func TestGetRecordInfo_FileNotFound(t *testing.T) {
	_, err := GetRecordInfoWithDateFilter(
		"nofile.mp4",
		"/nonexistent_path",
		false,
		time.Time{},
		"cam",
		"http://b",
	)
	assert.ErrorIs(t, err, ErrGetFileInfo)
}
