package cdn

import (
	"crypto/md5"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

func GetCDNAuthUrl(baseUrl, key, path string) string {
	if len(path) < 10 {
		return ""
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	ts := time.Now().Unix()
	rn := strings.ReplaceAll(uuid.New().String(), "-", "")
	uid := 0
	md5Str := fmt.Sprintf("%s-%d-%s-%d-%s", path, ts, rn, uid, key)
	md5Val := fmt.Sprintf("%x", md5.Sum([]byte(md5Str)))
	return fmt.Sprintf("%s%s?auth_key=%d-%s-%d-%s", baseUrl, path, ts, rn, uid, md5Val)
}
