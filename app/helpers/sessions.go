package helpers

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("secret-key"))

const (
	SessionKey  = "session-key"
	SessionName = "session-name"
)

func AppendFlash(w http.ResponseWriter, r *http.Request, key string, value string) {
	session, _ := store.Get(r, SessionName)
	session.AddFlash(value, key)
	session.Save(r, w)
}

func getFlashMessages(w http.ResponseWriter, r *http.Request, session *sessions.Session, key string) []string {
	flashes := session.Flashes(key)
	session.Save(r, w)
	messages := []string{}
	for _, v := range flashes {
		messages = append(messages, v.(string))
	}
	return messages
}

func GetFlashes(w http.ResponseWriter, r *http.Request) struct {
	FlashSuccess []string
	FlashError   []string
	FlashNotice  []string
} {
	session, _ := store.Get(r, SessionName)

	flash := struct {
		FlashSuccess []string
		FlashError   []string
		FlashNotice  []string
	}{
		FlashSuccess: getFlashMessages(w, r, session, FlashSuccess),
		FlashError:   getFlashMessages(w, r, session, FlashError),
		FlashNotice:  getFlashMessages(w, r, session, FlashNotice),
	}

	log.Printf("flash: %v", flash)
	return flash
}

func ClearFlashes(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, SessionName)
	session.Flashes(FlashError)
	session.Flashes(FlashSuccess)
	session.Flashes(FlashNotice)
	session.Save(r, w)
}

func CreateSession(w http.ResponseWriter, r *http.Request, userUUID string) {
	session, _ := store.Get(r, SessionName)
	session.Values["userUUID"] = userUUID
	session.Save(r, w)
}

func GetSession(r *http.Request) string {
	session, _ := store.Get(r, SessionName)
	userUUID := session.Values["userUUID"]
	if userUUID == nil {
		return ""
	}
	return userUUID.(string)
}
