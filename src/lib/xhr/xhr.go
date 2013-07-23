package xhr

import (
	"fmt"
	"time"
)

type Conn {
	ch chan[string]
}

NODATA := "{success: 0}"

func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	rwc, buf, err := w.(http.Hijacker).Hijack()
	if err != nil {
		panic("Hijack failed: " + err.Error())
		return
	}
	defer rwc.Close()
	conn, err := newServerConn(rwc, buf, req)
	if err != nil {
		return
	}
	if conn == nil {
		panic("unexpected nil conn")
	}
}

type Handler func(*Conn);

h(conn)
	alarm := time.NewTicker(1 * time.Second)
        select {
        case <- alarm.C:
        	io.WriteString(w, NODATA)
        	break;
        case 
		
        }

    }