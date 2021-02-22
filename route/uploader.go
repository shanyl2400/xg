package route

import (
	"errors"
	"net/http"
	"strings"
	"xg/conf"
	"xg/log"
	"xg/service"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

var (
	ErrInvalidPartition = errors.New("invalid partition")
)

type Partition string

func (p Partition) PartitionPath() string {
	uploads := conf.Get().UploadPath
	if p == "avatar" {
		return uploads + "/avatar"
	}
	return uploads + "/others"
}

func NewPartition(p string) (Partition, error) {
	switch p {
	case "avatar":
		return "avatar", nil
	case "org_attach":
		return "org_attach", nil
	default:
		return "", ErrInvalidPartition
	}
}

func FileName(fileName string) string {
	id := bson.NewObjectId()
	parts := strings.Split(fileName, ".")
	if len(parts) < 2 {
		return id.Hex()
	}
	ext := parts[len(parts)-1]
	return id.Hex() + "." + ext
}

// @Summary uploadFile
// @Description upload a file
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param partition path string true "upload file partition"
// @Tags upload
// @Success 200 {object} FileNameResponse
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/upload/{partition} [post]
func (s *Server) uploadFile(c *gin.Context) {
	partition, err := NewPartition(c.Param("partition"))
	if err != nil {
		log.Error.Println(err)
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	user := s.getJWTUser(c)

	f, err := c.FormFile("file")
	if err != nil {
		log.Error.Println(err)
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	name := FileName(f.Filename)
	err = c.SaveUploadedFile(f, partition.PartitionPath()+"/"+name)
	if err != nil {
		log.Error.Println(err)
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	//上传头像，更新用户信息
	if partition == "avatar" {
		err = service.GetUserService().UpdateUserAvatar(c.Request.Context(), name, user)
		if err != nil {
			log.Error.Println(err)
			s.responseErr(c, http.StatusInternalServerError, err)
			return
		}
	}

	c.JSON(http.StatusOK, FileNameResponse{
		Name:   name,
		ErrMsg: "success",
	})
}
