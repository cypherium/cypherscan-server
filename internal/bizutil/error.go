package bizutil

import (
	"fmt"

	"github.com/cypherium/cypherscan-server/internal/util"
	log "github.com/sirupsen/logrus"
)

// HandleError is to print error with detiled info
func HandleError(err error, message string) {
	key := 1
	logError(err, func(err error) { log.Errorf("%#+v", err) })
	fmt.Printf("[%v] %v", key, message)
	logError(err, func(err error) { log.Errorf("%s", err.Error()) })
	log.WithFields(log.Fields{
		"id":  key,
		"err": err,
	}).Error(message)
}

// HandleFatal is to log stack and quit app
func HandleFatal(err error, message string) {
	key := 1
	logError(err, func(err error) { log.Errorf("%#+v", err) })
	fmt.Printf("[%v] %v", key, message)
	logError(err, func(err error) { log.Errorf("%s", err.Error()) })
	log.WithFields(log.Fields{
		"id":  key,
		"err": err,
	}).Fatal(message)
}

type print func(err error)

func logError(err error, p print) {
	if err == nil {
		return
	}
	switch err.(type) {
	case *util.MyError:
		logError((err.(*util.MyError)).Inner, p)
		p(err)
	default:
		p(err)
	}
}
