package handle

import (
	"bytes"
	"fileServer/app/pkg/minio"
	"fileServer/app/utils/fontToImg"
	"fmt"
	"image/png"
	"net/http"
)

func NickName(w http.ResponseWriter, r *http.Request) {
	// 拿到参数
	query := r.URL.Query()
	nickName := query.Get("name")

	// 生成图片
	rgba := fontToImg.ToImg(nickName)
	proverbs := new(bytes.Buffer)
	err := png.Encode(proverbs, rgba)
	if err != nil {
		fmt.Println(err)
		return
	}
	readSeeker := bytes.NewReader(proverbs.Bytes())

	// 传到服务器中
	minio.Upload("nickName"+".png", readSeeker)

	// 响应
	fmt.Fprintf(w, "{\"code\":0}")
}
