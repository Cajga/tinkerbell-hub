package archive

// This package handles the pulling and management of images

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Write will pull an image and write it to local storage device
// image must be a tar file or tar.gz file as set by archiveType
func Write(archiveURL, archiveType, path string, checksum string) error {

	req, err := http.NewRequest("GET", archiveURL, nil)
	if err != nil {
		return err
	}

	//timeout added due to clients timing out when provisioning
	//multiple systems at once.  May want to add this as an
	//input parameter
	client := &http.Client{Timeout: time.Second * 300}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 300 {
		// Customise response for the 404 to make degugging simpler
		if resp.StatusCode == 404 {
			return fmt.Errorf("%s not found", archiveURL)
		}
		return fmt.Errorf("%s", resp.Status)
	}

	absPath, err := filepath.Abs(path)

	switch strings.ToLower(archiveType) {
	case "gz":
		// With compression run data through gzip writer
		// zipOUT, err := gzip.NewReader(resp.Body)
		// if err != nil {
		// 	fmt.Println("[ERROR] New gzip reader:", err)
		// }
		// defer zipOUT.Close()
		// out = zipOUT
	case "tar":
		err := extractTarDirectory(absPath, checksum, resp.Body)
		if err != nil {
		 	log.Fatalf("[ERROR] New gzip reader:", err)
		}
	case "targz":
		err := extractTarGzip(absPath, checksum, resp.Body)		
		if err != nil {
			log.Fatalf("[ERROR] New gzip reader:", err)
		}
	default:
		log.Fatalf("[ERROR] Unknown archiveType supplied:", archiveType)
	}
	return nil
}
