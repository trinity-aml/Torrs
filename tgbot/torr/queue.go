package torr

import (
	"fmt"
	"github.com/dustin/go-humanize"
	tele "gopkg.in/telebot.v4"
	"strconv"
	"time"
)

type DLQueue struct {
	id        int
	c         tele.Context
	hash      string
	fileID    string
	fileName  string
	updateMsg *tele.Message
}

var (
	manager = &Manager{}
)

func Start() {
	manager.Start()
}

func Show(c tele.Context) error {
	//msg := ""
	//mu.Lock()
	//for i, dlQueue := range queue {
	//	s := "#" + strconv.Itoa(i+1) + ":\n<b>Хэш:</b> <code>" + dlQueue.hash + "</code>\n<i>" + filepath.Base(dlQueue.fileName) + "</i>\n"
	//	if len(msg+s) > 1024 {
	//		err := c.Send(msg)
	//		if err != nil {
	//			return err
	//		}
	//		msg = ""
	//	}
	//	msg += s
	//}
	//mu.Unlock()
	//if msg != "" {
	//	return c.Send("Очередь:\n" + msg)
	//} else {
	//	return c.Send("Очередь пуста")
	//}
	return nil
}

func AddAll(c tele.Context, hash string) {
	manager.Add(c, hash, "all")
}

func Add(c tele.Context, hash, fileID string) {
	manager.Add(c, hash, "file:"+fileID)
}

func Cancel(id int) {
	manager.Cancel(id)
}

func updateLoadStatus(wrk *Worker, file *TorrFile, fi, fc int) {
	ti, err := GetTorrentInfo(wrk.torrentHash)
	if err != nil {
		wrk.c.Bot().Edit(wrk.msg, "Ошибка при получении данных о торренте")
	} else if wrk.isCancelled {
		wrk.c.Bot().Edit(wrk.msg, "Остановка...")
	} else {
		wrk.c.Send(tele.UploadingVideo)
		if ti.DownloadSpeed == 0 {
			ti.DownloadSpeed = 1.0
		}
		wait := time.Duration(float64(file.Loaded())/ti.DownloadSpeed) * time.Second
		speed := humanize.Bytes(uint64(ti.DownloadSpeed)) + "/sec"
		peers := fmt.Sprintf("%v · %v/%v", ti.ConnectedSeeders, ti.ActivePeers, ti.TotalPeers)
		prc := fmt.Sprintf("%.2f%% %v / %v", float64(file.offset)*100.0/float64(file.size), humanize.Bytes(uint64(file.offset)), humanize.Bytes(uint64(file.size)))

		name := file.name
		if name == ti.Title {
			name = ""
		}

		msg := "Загрузка торрента:\n" +
			"<b>" + ti.Title + "</b>\n"
		if name != "" {
			msg += "<i>" + name + "</i>\n"
		}
		msg += "<b>Хэш:</b> <code>" + file.hash + "</code>\n"
		if file.offset < file.size {
			msg += "<b>Скорость: </b>" + speed + "\n" +
				"<b>Осталось: </b>" + wait.String() + "\n" +
				"<b>Пиры: </b>" + peers + "\n" +
				"<b>Загружено: </b>" + prc
		}
		if fc > 1 {
			msg += "\n<b>Файлов: </b>" + strconv.Itoa(fi) + "/" + strconv.Itoa(fc)
		}
		if file.offset >= file.size {
			msg += "\n<b>Завершение загрузки, это займет некоторое время</b>"
		}

		torrKbd := &tele.ReplyMarkup{}
		torrKbd.Inline([]tele.Row{torrKbd.Row(torrKbd.Data("Отмена", "cancel", strconv.Itoa(wrk.id)))}...)
		wrk.c.Bot().Edit(wrk.msg, msg, torrKbd)
	}
}
