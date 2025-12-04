package writerastochnoy

import (
    "bytes"
    "context"
    "encoding/json"
    "log"
    "net/http"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"

    rass "rastochnoy/genproto/rastochnoy"
)

var (
    clientsDB9      = make(map[*websocket.Conn]bool)
    clientsMuDB9    sync.Mutex

    lastFullDB9     []byte
    lastFullMuDB9   sync.Mutex
    initializedDB9  bool

    upgraderDB9 = websocket.Upgrader{
        CheckOrigin: func(r *http.Request) bool { return true },
    }

    repoReaderDB9 func(context.Context, *rass.ReadWriteRastochnoyReq) (*rass.ReadWriteRastochnoyRes, error)
)

// Init
func InitRastochnoyWSDB9(reader func(context.Context, *rass.ReadWriteRastochnoyReq) (*rass.ReadWriteRastochnoyRes, error)) {
    repoReaderDB9 = reader
    go pollingDB9()
}

// WebSocket handler
func WebSocketHandlerDB9(c *gin.Context) {
    conn, err := upgraderDB9.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        return
    }

    clientsMuDB9.Lock()
    clientsDB9[conn] = true
    clientsMuDB9.Unlock()

    log.Println("New WS client connected (DB9)")

    // 1) Dastlab FULL yuborish
    if repoReaderDB9 != nil {
        data, err := repoReaderDB9(context.Background(), &rass.ReadWriteRastochnoyReq{})
        if err == nil && data != nil {
            msg, _ := json.Marshal(data)
            conn.WriteMessage(websocket.TextMessage, msg)

            lastFullMuDB9.Lock()
            if !initializedDB9 {
                lastFullDB9 = msg
                initializedDB9 = true
            }
            lastFullMuDB9.Unlock()
        }
    }

    go func() {
        defer func() {
            clientsMuDB9.Lock()
            delete(clientsDB9, conn)
            clientsMuDB9.Unlock()
            conn.Close()
        }()
        for {
            if _, _, err := conn.ReadMessage(); err != nil {
                break
            }
        }
    }()
}

// Polling: send only diff
func pollingDB9() {
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()

    for range ticker.C {

        if repoReaderDB9 == nil {
            continue
        }

        newData, err := repoReaderDB9(context.Background(), &rass.ReadWriteRastochnoyReq{})
        if err != nil || newData == nil {
            continue
        }

        // ---- Yangi map ID->item ----
        newMap := make(map[string]*rass.ReadWriteRastoshnoyItem)
        for _, v := range newData.Data {
            newMap[*v.Id] = v
        }

        lastFullMuDB9.Lock()
        oldBytes := lastFullDB9
        isInit := initializedDB9
        lastFullMuDB9.Unlock()

        // DASTLABKI holat – FULL saqlaymiz lekin diff yubormaymiz
        if !isInit || len(oldBytes) == 0 {
            b, _ := json.Marshal(newData)

            lastFullMuDB9.Lock()
            lastFullDB9 = b
            initializedDB9 = true
            lastFullMuDB9.Unlock()

            continue
        }

        // ---- Eski map ID->item ----
        var oldData rass.ReadWriteRastochnoyRes
        json.Unmarshal(oldBytes, &oldData)

        oldMap := make(map[string]*rass.ReadWriteRastoshnoyItem)
        for _, v := range oldData.Data {
            oldMap[*v.Id] = v
        }

        // ---- Faqat o'zgarganlarni topamiz ----
        diffArr := make([]*rass.ReadWriteRastoshnoyItem, 0)

        for id, newItem := range newMap {
            oldItem, exists := oldMap[id]
            if !exists {
                diffArr = append(diffArr, newItem)
                continue
            }

            if *oldItem.Value != *newItem.Value || *oldItem.Offsett != *newItem.Offsett {
                diffArr = append(diffArr, newItem)
            }
        }

        // Hech narsa o‘zgarmagan
        if len(diffArr) == 0 {
            continue
        }

        // ---- Diff yuboramiz ----
        diffBody := map[string]interface{}{
            "data": diffArr,
        }

        diffBytes, _ := json.Marshal(diffBody)
        broadcastDB9(diffBytes)

        // ---- To‘liq data ni yangilaymiz ----
        b, _ := json.Marshal(newData)
        lastFullMuDB9.Lock()
        lastFullDB9 = b
        lastFullMuDB9.Unlock()
    }
}

func broadcastDB9(msg []byte) {
    clientsMuDB9.Lock()
    defer clientsMuDB9.Unlock()

    for conn := range clientsDB9 {
        if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
            conn.Close()
            delete(clientsDB9, conn)
        }
    }
}

///////////////////////////////////////////////////////
// DIFF UTILITIES (o'zgarmaydi)
///////////////////////////////////////////////////////

func diffMaps(oldMap, newMap map[string]interface{}) map[string]interface{} {
    diff := make(map[string]interface{})

    for k, newVal := range newMap {
        oldVal, exists := oldMap[k]

        if !exists {
            diff[k] = newVal
            continue
        }

        oldSub, okOld := oldVal.(map[string]interface{})
        newSub, okNew := newVal.(map[string]interface{})
        if okOld && okNew {
            sub := diffMaps(oldSub, newSub)
            if len(sub) > 0 {
                diff[k] = sub
            }
            continue
        }

        oldSlice, okOldArr := oldVal.([]interface{})
        newSlice, okNewArr := newVal.([]interface{})
        if okOldArr && okNewArr {
            if !equalSlices(oldSlice, newSlice) {
                diff[k] = newVal
            }
            continue
        }

        if !valuesEqual(oldVal, newVal) {
            diff[k] = newVal
        }
    }

    return diff
}

func equalSlices(a, b []interface{}) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if !valuesEqual(a[i], b[i]) {
            return false
        }
    }
    return true
}

func valuesEqual(a, b interface{}) bool {
    aj, _ := json.Marshal(a)
    bj, _ := json.Marshal(b)
    return bytes.Equal(aj, bj)
}
