package cmd

import (
	"net/http"
	"strconv"
	"time"

	"github.com/amirhnajafiz/Blue-sky/internal/pion/media"
	"github.com/amirhnajafiz/Blue-sky/internal/pion/signal"
	"github.com/amirhnajafiz/Blue-sky/internal/pion/track"
	"github.com/gin-gonic/gin"
	"github.com/pion/webrtc/v2"
)

const (
	rtcpPLIInterval = time.Second * 3
)

type Sdp struct {
	Sdp string
}

func Exec() {
	router := gin.Default()

	peerConnectionMap := make(map[string]chan *webrtc.Track)
	api := media.GetMediaAPI()
	peerConnectionConfig := media.GetPeerConfig()

	router.POST("/webrtc/sdp/m/:meetingId/c/:userID/p/:peerId/s/:isSender", func(c *gin.Context) {
		isSender, _ := strconv.ParseBool(c.Param("isSender"))
		userID := c.Param("userID")
		peerID := c.Param("peerId")

		var session Sdp
		if err := c.ShouldBindJSON(&session); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		offer := webrtc.SessionDescription{}
		signal.Decode(session.Sdp, &offer)

		peerConnection, _ := api.NewPeerConnection(peerConnectionConfig)
		_ = peerConnection.SetRemoteDescription(offer)

		if !isSender {
			track.ReceiveTrack(peerConnection, peerConnectionMap, peerID)
		}
	})

	_ = router.Run(":8080")
}
