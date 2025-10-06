package records

import "fmt"

var ErrCameraNotFound error = fmt.Errorf("camera not found")
var ErrGetRecordsFailed error = fmt.Errorf("failed to get records")
var ErrGetFileInfo error = fmt.Errorf("failed to get file info")
var ErrDateNotMatch error = fmt.Errorf("filter date not match")

