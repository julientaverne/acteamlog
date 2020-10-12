package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"	
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	//"errors"
	"log"
	"bytes"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jz222/logowl/internal/models"
	"github.com/jz222/logowl/internal/services"
	"github.com/jz222/logowl/internal/store"
	"github.com/jz222/logowl/internal/utils"
)

type LoggingControllers struct {
	LoggingService services.InterfaceLogging
}

func decrypt(key []byte, securemess string) (decodedmess string, err error) {
	iv := "TestingIV1234567"
	ciphertext, err := base64.StdEncoding.DecodeString(securemess)
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, ciphertext)

	decodedmess = string(ciphertext)
	return
}

func (l *LoggingControllers) RegisterError(c *gin.Context) {
	
	errorEvent := models.Error{
		Badges:    map[string]string{},
		UserAgent: c.Request.UserAgent(),
		Count:     1,
		Timestamp: time.Now().Unix(),
	}
	

	CIPHER_KEY := []byte("NFd6N3v1nbL47FK0xpZjxZ7NY4fYpNYd")

	buf := new(bytes.Buffer)
    buf.ReadFrom(c.Request.Body)
    newStr := buf.String()

	log.Printf("ENCRYPTED: %s\n", newStr)

	if decrypted, err := decrypt(CIPHER_KEY, newStr); err != nil {
		log.Println(">>>>>>>>>", err)
	} else {
		log.Printf("+++++++++++ DECRYPTED: %s\n", decrypted)

		
		err := json.NewDecoder(strings.NewReader(decrypted)).Decode(&errorEvent)
		if err != nil {
			utils.RespondWithError(c, http.StatusBadRequest, err.Error())
			return
		}
	
		if !errorEvent.AnonymizeData {
			errorEvent.ClientIP = c.ClientIP()
		}
	
		if !errorEvent.IsValid() {
			utils.RespondWithError(c, http.StatusBadRequest, "the provided data is too large")
			return
		}
	
		go l.LoggingService.SaveError(errorEvent)
	
		utils.RespondWithSuccess(c)
		

	}


}

func (l *LoggingControllers) RegisterAnalyticEvent(c *gin.Context) {
	var analyticEvent models.AnalyticEvent

	err := json.NewDecoder(c.Request.Body).Decode(&analyticEvent)
	if err != nil {
		fmt.Println(err.Error())
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	if analyticEvent.Ticket == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "the ticket was not provided")
		return
	}

	analyticEvent.UserAgent = c.Request.UserAgent()

	go l.LoggingService.SaveAnalyticEvent(analyticEvent)

	utils.RespondWithSuccess(c)
}

func GetLoggingController(store store.InterfaceStore) LoggingControllers {
	loggingService := services.GetLoggingService(store)

	return LoggingControllers{
		LoggingService: &loggingService,
	}
}
