package app

import (
	"bytes"
	"encoding/hex"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"
	"word/pkg/unique"
)

type fileTypeMap struct {
	TypeMap sync.Map
	sync.Once
}

// SaveHandler 自定义文件上传之后的保存操作
type SaveHandler interface {
	// 保存文件并返回文件最终路径
	Save(file *multipart.FileHeader, fileName string) string
}

// DefaultSaveHandler 默认文件保存器
type DefaultSaveHandler struct {
	prefix  string
	dst     string
	context *gin.Context
}

// SetDst 设置保存位置
func (defaultSaveHandler *DefaultSaveHandler) SetDst(dst string) *DefaultSaveHandler {
	defaultSaveHandler.dst = dst
	return defaultSaveHandler
}

// SetPrefix 设置前缀
func (defaultSaveHandler *DefaultSaveHandler) SetPrefix(prefix string) *DefaultSaveHandler {
	defaultSaveHandler.prefix = prefix
	return defaultSaveHandler
}

// Save 保存
func (defaultSaveHandler *DefaultSaveHandler) Save(file *multipart.FileHeader, fileName string) string {
	filePath := defaultSaveHandler.dst + defaultSaveHandler.prefix + fileName
	err := defaultSaveHandler.context.SaveUploadedFile(file, filePath)
	if err != nil {
		Logger().WithField("log_type", "pkg.app.upload").Error(err)
		return ""
	}

	return strings.ReplaceAll(filePath, Root(), "")
}

type (
	// OSSConnectInfo 连接配置
	OSSConnectInfo struct {
		EndPoint        string `mapstructure:"end_point"`
		AccessKeyID     string `mapstructure:"access_key_id"`
		AccessKeySecret string `mapstructure:"access_key_secret"`
		BucketName      string `mapstructure:"bucket_name"`
		ResourceDomain  string `mapstructure:"resource_domain"`
	}

	// OSSUploader 上传
	OSSUploader struct {
		Config OSSConnectInfo
		client *oss.Client
	}
)

// NewOSSConnector 实例化上传器
func NewOSSConnector(config OSSConnectInfo) *OSSUploader {
	return &OSSUploader{
		Config: config,
	}
}

// Do 执行上传
func (uploader *OSSUploader) Do(filename string, data io.Reader) (err error) {
	uploader.client, err = oss.New(uploader.Config.EndPoint, uploader.Config.AccessKeyID, uploader.Config.AccessKeySecret)
	if err != nil {
		return
	}

	bucket, err := uploader.client.Bucket(uploader.Config.BucketName)
	if err != nil {
		return
	}

	err = bucket.PutObject(filename, data, oss.ObjectACL(oss.ACLPublicRead))
	return
}

// AsyncDo 异步上传
// 可以传入第三个参数以自定义上传错误处理方法
func (uploader *OSSUploader) AsyncDo(filename string, data io.Reader, errHandler ...func(err error)) {
	go func() {
		err := uploader.Do(filename, data)
		if err != nil && len(errHandler) == 1 {
			errHandler[0](err)
		}
	}()
}

// OSS 阿里云存储
type OSS struct {
	prefix string
	dst    string
	async  bool
}

// SetDst 设置保存位置
func (oss *OSS) setDst(host string) *OSS {
	oss.dst = host
	return oss
}

// SetPrefix 设置前缀, 通常是这样的 a/b/c 最前边一定不要有反斜杠
func (oss *OSS) SetPrefix(prefix string) *OSS {
	oss.prefix = prefix
	return oss
}

// Async 异步上传
func (oss *OSS) Async() *OSS {
	oss.async = true
	return oss
}

// Save 保存
func (oss *OSS) Save(file *multipart.FileHeader, fileName string) string {
	var config OSSConnectInfo
	err := InitConfig("aliyun_oss", &config)
	if err != nil {
		Logger().Warn("aliyun oss connect config error")
	}
	var connector = NewOSSConnector(config)
	var filename = unique.IDStr() + path.Ext(file.Filename)
	oss.setDst(connector.Config.ResourceDomain)

	if oss.prefix == "" {
		oss.prefix = "base/assets/"
	}

	var reader, _ = file.Open()
	if oss.async {
		connector.AsyncDo(oss.prefix+filename, reader, func(err error) {
			Logger().Warn(err)
		})
	} else {
		_ = connector.Do(oss.prefix+filename, reader)
	}

	return oss.dst + oss.prefix + filename
}

// FileType 文件类型
var FileType = new(fileTypeMap)

func (fileTypeMap *fileTypeMap) lazyLoad() {
	fileTypeMap.Do(func() {
		fileTypeMap.TypeMap.Store("ffd8ffe", "jpg")                     // JPEG (jpg)
		fileTypeMap.TypeMap.Store("89504e470d0a1a0a0000", "png")        // PNG (png)
		fileTypeMap.TypeMap.Store("474946383961", "gif")                // GIF (gif)
		fileTypeMap.TypeMap.Store("49492a00227105008037", "tif")        // TIFF (tif)
		fileTypeMap.TypeMap.Store("424d228c010000000000", "bmp")        // 16色位图(bmp)
		fileTypeMap.TypeMap.Store("424d8240090000000000", "bmp")        // 24位位图(bmp)
		fileTypeMap.TypeMap.Store("424d8e1b030000000000", "bmp")        // 256色位图(bmp)
		fileTypeMap.TypeMap.Store("41433130313500000000", "dwg")        // CAD (dwg)
		fileTypeMap.TypeMap.Store("3c21444f435459504520", "html")       // HTML (html)   3c68746d6c3e0  3c68746d6c3e0
		fileTypeMap.TypeMap.Store("3c68746d6c3e0", "html")              // HTML (html)   3c68746d6c3e0  3c68746d6c3e0
		fileTypeMap.TypeMap.Store("3c21646f637479706520", "htm")        // HTM (htm)
		fileTypeMap.TypeMap.Store("48544d4c207b0d0a0942", "css")        // css
		fileTypeMap.TypeMap.Store("696b2e71623d696b2e71", "js")         // js
		fileTypeMap.TypeMap.Store("7b5c727466315c616e73", "rtf")        // Rich Text Format (rtf)
		fileTypeMap.TypeMap.Store("38425053000100000000", "psd")        // Photoshop (psd)
		fileTypeMap.TypeMap.Store("46726f6d3a203d3f6762", "eml")        // Email [Outlook Express 6] (eml)
		fileTypeMap.TypeMap.Store("d0cf11e0a1b11ae10000", "doc")        // MS Excel 注意：word、msi 和 excel的文件头一样
		fileTypeMap.TypeMap.Store("d0cf11e0a1b11ae10000", "vsd")        // Visio 绘图
		fileTypeMap.TypeMap.Store("5374616E64617264204A", "mdb")        // MS Access (mdb)
		fileTypeMap.TypeMap.Store("252150532D41646F6265", "ps")         //
		fileTypeMap.TypeMap.Store("255044462d312", "pdf")               // Adobe Acrobat (pdf)
		fileTypeMap.TypeMap.Store("2e524d46000000120001", "rmvb")       // rmvb/rm相同
		fileTypeMap.TypeMap.Store("464c5601050000000900", "flv")        // flv与f4v相同
		fileTypeMap.TypeMap.Store("00000020667479706d70", "mp4")        //
		fileTypeMap.TypeMap.Store("49443303000000002176", "mp3")        //
		fileTypeMap.TypeMap.Store("000001ba210001000180", "mpg")        //
		fileTypeMap.TypeMap.Store("3026b2758e66cf11a6d9", "wmv")        // wmv与asf相同
		fileTypeMap.TypeMap.Store("52494646e27807005741", "wav")        // Wave (wav)
		fileTypeMap.TypeMap.Store("52494646d07d60074156", "avi")        //
		fileTypeMap.TypeMap.Store("4d546864000000060001", "mid")        // MIDI (mid)
		fileTypeMap.TypeMap.Store("504b0304140000000800", "zip")        //
		fileTypeMap.TypeMap.Store("526172211a0700cf9073", "rar")        //
		fileTypeMap.TypeMap.Store("235468697320636f6e66", "ini")        //
		fileTypeMap.TypeMap.Store("504b03040a0000000000", "jar")        //
		fileTypeMap.TypeMap.Store("4d5a9000030000000400", "exe")        // 可执行文件
		fileTypeMap.TypeMap.Store("3c25402070616765206c", "jsp")        // jsp文件
		fileTypeMap.TypeMap.Store("4d616e69666573742d56", "mf")         // MF文件
		fileTypeMap.TypeMap.Store("3c3f786d6c2076657273", "xml")        // xml文件
		fileTypeMap.TypeMap.Store("494e5345525420494e54", "sql")        // xml文件
		fileTypeMap.TypeMap.Store("7061636b616765207765", "java")       // java文件
		fileTypeMap.TypeMap.Store("406563686f206f66660d", "bat")        // bat文件
		fileTypeMap.TypeMap.Store("1f8b0800000000000000", "gz")         // gz文件
		fileTypeMap.TypeMap.Store("6c6f67346a2e726f6f74", "properties") // bat文件
		fileTypeMap.TypeMap.Store("cafebabe0000002e0041", "class")      // bat文件
		fileTypeMap.TypeMap.Store("49545346030000006000", "chm")        // bat文件
		fileTypeMap.TypeMap.Store("04000000010000001300", "mxp")        // bat文件
		fileTypeMap.TypeMap.Store("504b0304140006000800", "docx")       // docx文件
		fileTypeMap.TypeMap.Store("d0cf11e0a1b11ae10000", "wps")        // WPS文字wps、表格et、演示dps都是一样的
		fileTypeMap.TypeMap.Store("6431303a637265617465", "torrent")    //
		fileTypeMap.TypeMap.Store("6D6F6F76", "mov")                    // Quicktime (mov)
		fileTypeMap.TypeMap.Store("FF575043", "wpd")                    // WordPerfect (wpd)
		fileTypeMap.TypeMap.Store("CFAD12FEC5FD746F", "dbx")            // Outlook Express (dbx)
		fileTypeMap.TypeMap.Store("2142444E", "pst")                    // Outlook (pst)
		fileTypeMap.TypeMap.Store("AC9EBD8F", "qdf")                    // Quicken (qdf)
		fileTypeMap.TypeMap.Store("E3828596", "pwl")                    // Windows Password (pwl)
		fileTypeMap.TypeMap.Store("2E7261FD", "ram")                    // Real Audio (ram)
	})
}

// 获取前面结果字节的二进制
func (fileTypeMap *fileTypeMap) bytesToHexString(src []byte) string {
	res := bytes.Buffer{}
	if src == nil || len(src) <= 0 {
		return ""
	}
	temp := make([]byte, 0)
	for _, v := range src {
		sub := v & 0xFF
		hv := hex.EncodeToString(append(temp, sub))
		if len(hv) < 2 {
			res.WriteString(strconv.FormatInt(int64(0), 10))
		}
		res.WriteString(hv)
	}
	return res.String()
}

// 用文件前面几个字节来判断
// fSrc: 文件字节流（就用前面几个字节）
func (fileTypeMap *fileTypeMap) GetFileType(fSrc []byte) string {
	fileTypeMap.lazyLoad()
	var fileType string
	fileCode := fileTypeMap.bytesToHexString(fSrc)
	fileTypeMap.TypeMap.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)
		if strings.HasPrefix(fileCode, strings.ToLower(k)) ||
			strings.HasPrefix(k, strings.ToLower(fileCode)) {
			fileType = v
			return false
		}
		return true
	})
	return fileType
}

// Upload 文件上传公共方法
//  key 上传文件的表单name, 如果是多文件需要加上中括号[]
//  dst 存放路径 注意:无论这里传什么路径, 最后边都会追加 filename.xxx
func Upload(key string, saveHandler SaveHandler, allowedTyp ...string) gin.HandlerFunc {
	return func(context *gin.Context) {
		form, _ := context.MultipartForm()
		files := form.File[key]
		formKey := context.PostForm("key")

		var response = NewResponse(Success, nil, "success")
		var data = make(map[string][]string, 0)
		for _, file := range files {
			f, err := file.Open()
			if err != nil {
				data[formKey] = make([]string, 0, 0)
			} else {
				byt, err := ioutil.ReadAll(f)
				if err != nil {
					data[formKey] = make([]string, 0, 0)
				} else {
					fileType := FileType.GetFileType(byt[:10])
					var typAllow = false
					for _, typ := range allowedTyp {
						typAllow = typAllow || typ == fileType
					}
					if typAllow {
						fileName := strconv.Itoa(int(unique.ID())) + "." + fileType
						data[formKey] = append(data[formKey], saveHandler.Save(file, fileName))
					} else {
						response.Message = Translator(context).Translate("upload file type not allow", gin.H{"type": fileType})
					}
				}
				_ = f.Close()
			}

		}

		response.Data = data
		context.JSON(http.StatusOK, response)
	}
}
