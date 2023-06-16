package rest

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) DebugImport(c *gin.Context) {
	keystoreName := "Passwords"
	b := []byte("")

	csvReader := csv.NewReader(bytes.NewReader(b))

	records, err := csvReader.ReadAll()
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	k, err := s.app.CreateKeystore(keystoreName)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	for _, row := range records {
		if len(row) != 4 {
			continue
		}

		website := row[0]
		// url := row[1]
		username := row[2]
		password := row[3]

		_, err := s.app.CreateEntry(
			k.Id,
			website,
			username,
			password,
			"",
		)
		if err != nil {
			log.Error(err)
			Error(c, http.StatusInternalServerError)
			return
		}
		fmt.Println("Imported: ", website, username, "***")

	}

	c.Status(http.StatusOK)
}
