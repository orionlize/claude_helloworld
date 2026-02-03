package main

import (
  "encoding/json"
  "log"
  "math"
  "net/http"
  "strings"
  "time"
)

type Metric struct {
  ID          string `json:"id"`
  Title       string `json:"title"`
  Value       string `json:"value"`
  Description string `json:"description"`
  Trend       string `json:"trend"`
  TrendType   string `json:"trendType"`
}

type SummaryStatus struct {
  Project     string `json:"project"`
  Environment string `json:"environment"`
  SyncState   string `json:"syncState"`
  LastSync    string `json:"lastSync"`
}

type Summary struct {
  Metrics []Metric      `json:"metrics"`
  Status  SummaryStatus `json:"status"`
}

type TrafficPoint struct {
  Label     string `json:"label"`
  Primary   int    `json:"primary"`
  Secondary int    `json:"secondary"`
}

type InterfaceItem struct {
  ID        string `json:"id"`
  Name      string `json:"name"`
  Method    string `json:"method"`
  Path      string `json:"path"`
  Category  string `json:"category"`
  Status    string `json:"status"`
  UpdatedAt string `json:"updatedAt"`
}

type FieldNode struct {
  Name        string      `json:"name"`
  Type        string      `json:"type"`
  Required    bool        `json:"required,omitempty"`
  Description string      `json:"description,omitempty"`
  Children    []FieldNode `json:"children,omitempty"`
}

type InterfaceDetail struct {
  InterfaceItem
  Description    string      `json:"description"`
  RequestFields  []FieldNode `json:"requestFields"`
  ResponseFields []FieldNode `json:"responseFields"`
}

type SyncLog struct {
  ID     string `json:"id"`
  Time   string `json:"time"`
  Action string `json:"action"`
  Status string `json:"status"`
  Detail string `json:"detail"`
  Author string `json:"author"`
}

type Store struct {
  Summary    Summary
  Traffic    map[string][]TrafficPoint
  Interfaces []InterfaceDetail
  SyncLogs   []SyncLog
}

func main() {
  store := newStore()

  mux := http.NewServeMux()
  mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
    writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
  })
  mux.HandleFunc("/api/summary", func(w http.ResponseWriter, r *http.Request) {
    writeJSON(w, http.StatusOK, store.Summary)
  })
  mux.HandleFunc("/api/traffic", func(w http.ResponseWriter, r *http.Request) {
    rangeParam := r.URL.Query().Get("range")
    if rangeParam == "" {
      rangeParam = "90d"
    }
    if data, ok := store.Traffic[rangeParam]; ok {
      writeJSON(w, http.StatusOK, data)
      return
    }
    writeJSON(w, http.StatusOK, store.Traffic["90d"])
  })
  mux.HandleFunc("/api/interfaces", func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/api/interfaces" {
      http.NotFound(w, r)
      return
    }
    list := make([]InterfaceItem, 0, len(store.Interfaces))
    for _, item := range store.Interfaces {
      list = append(list, item.InterfaceItem)
    }
    writeJSON(w, http.StatusOK, list)
  })
  mux.HandleFunc("/api/interfaces/", func(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/api/interfaces/")
    if id == "" {
      http.NotFound(w, r)
      return
    }
    for _, item := range store.Interfaces {
      if item.ID == id {
        writeJSON(w, http.StatusOK, item)
        return
      }
    }
    http.NotFound(w, r)
  })
  mux.HandleFunc("/api/sync-logs", func(w http.ResponseWriter, r *http.Request) {
    writeJSON(w, http.StatusOK, store.SyncLogs)
  })

  handler := withCORS(mux)

  srv := &http.Server{
    Addr:         ":8080",
    Handler:      handler,
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 5 * time.Second,
  }

  log.Println("API server listening on :8080")
  if err := srv.ListenAndServe(); err != nil {
    log.Fatal(err)
  }
}

func withCORS(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    if r.Method == http.MethodOptions {
      w.WriteHeader(http.StatusNoContent)
      return
    }
    next.ServeHTTP(w, r)
  })
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(status)
  if err := json.NewEncoder(w).Encode(payload); err != nil {
    log.Println("encode error:", err)
  }
}

func newStore() Store {
  now := time.Now()
  summary := Summary{
    Metrics: []Metric{
      {
        ID:          "interfaces",
        Title:       "同步接口数",
        Value:       "1,284",
        Description: "本周新增 38",
        Trend:       "+4.6%",
        TrendType:   "up",
      },
      {
        ID:          "sync-rate",
        Title:       "同步成功率",
        Value:       "99.92%",
        Description: "过去 7 天",
        Trend:       "+0.12%",
        TrendType:   "up",
      },
      {
        ID:          "test-pass",
        Title:       "测试通过率",
        Value:       "92.4%",
        Description: "最近 120 组",
        Trend:       "-1.4%",
        TrendType:   "down",
      },
      {
        ID:          "doc-visits",
        Title:       "文档访问量",
        Value:       "45,678",
        Description: "30 天内",
        Trend:       "+12.5%",
        TrendType:   "up",
      },
    },
    Status: SummaryStatus{
      Project:     "API 平台核心服务",
      Environment: "staging",
      SyncState:   "healthy",
      LastSync:    now.Format("2006-01-02 15:04"),
    },
  }

  traffic := map[string][]TrafficPoint{
    "90d": generateTraffic(12, now),
    "30d": generateTraffic(10, now),
    "7d":  generateTraffic(7, now),
  }

  interfaces := []InterfaceDetail{
    {
      InterfaceItem: InterfaceItem{
        ID:        "intf_001",
        Name:      "同步接口列表",
        Method:    "GET",
        Path:      "/api/interfaces",
        Category:  "同步层",
        Status:    "active",
        UpdatedAt: now.Format("2006-01-02 15:04"),
      },
      Description:   "按条件分页查询当前项目同步的接口列表。",
      RequestFields: requestFields(),
      ResponseFields: []FieldNode{
        {
          Name: "success",
          Type: "boolean",
        },
        {
          Name: "data",
          Type: "object",
          Children: []FieldNode{
            {
              Name: "items",
              Type: "array",
              Children: []FieldNode{
                {Name: "interfaceId", Type: "string"},
                {Name: "name", Type: "string"},
                {Name: "method", Type: "string"},
                {Name: "status", Type: "string"},
              },
            },
            {Name: "total", Type: "number"},
          },
        },
        {Name: "traceId", Type: "string"},
      },
    },
    {
      InterfaceItem: InterfaceItem{
        ID:        "intf_002",
        Name:      "创建测试套件",
        Method:    "POST",
        Path:      "/api/tests/suites",
        Category:  "测试层",
        Status:    "active",
        UpdatedAt: now.Add(-45 * time.Minute).Format("2006-01-02 15:04"),
      },
      Description:   "基于接口定义创建测试套件。",
      RequestFields: requestFields(),
      ResponseFields: []FieldNode{
        {Name: "suiteId", Type: "string"},
        {Name: "createdAt", Type: "string"},
      },
    },
    {
      InterfaceItem: InterfaceItem{
        ID:        "intf_003",
        Name:      "Mock 规则列表",
        Method:    "GET",
        Path:      "/api/mock/rules",
        Category:  "Mock 层",
        Status:    "active",
        UpdatedAt: now.Add(-2 * time.Hour).Format("2006-01-02 15:04"),
      },
      Description:   "查询项目的 Mock 规则列表。",
      RequestFields: requestFields(),
      ResponseFields: []FieldNode{
        {Name: "items", Type: "array"},
        {Name: "total", Type: "number"},
      },
    },
    {
      InterfaceItem: InterfaceItem{
        ID:        "intf_004",
        Name:      "文档分享设置",
        Method:    "PATCH",
        Path:      "/api/docs/share",
        Category:  "文档层",
        Status:    "deprecated",
        UpdatedAt: now.Add(-12 * time.Hour).Format("2006-01-02 15:04"),
      },
      Description:   "更新文档分享权限配置。",
      RequestFields: requestFields(),
      ResponseFields: []FieldNode{
        {Name: "shareId", Type: "string"},
        {Name: "status", Type: "string"},
      },
    },
  }

  logs := []SyncLog{
    {
      ID:     "log_01",
      Time:   now.Add(-8 * time.Minute).Format("15:04"),
      Action: "接口更新",
      Status: "success",
      Detail: "同步 12 条接口，2 条分类",
      Author: "easy-yapi",
    },
    {
      ID:     "log_02",
      Time:   now.Add(-30 * time.Minute).Format("15:04"),
      Action: "字段变更",
      Status: "warning",
      Detail: "检测到响应字段缺少描述",
      Author: "easy-yapi",
    },
    {
      ID:     "log_03",
      Time:   now.Add(-58 * time.Minute).Format("15:04"),
      Action: "接口删除",
      Status: "success",
      Detail: "已移除废弃接口 1 条",
      Author: "easy-yapi",
    },
    {
      ID:     "log_04",
      Time:   now.Add(-90 * time.Minute).Format("15:04"),
      Action: "同步失败",
      Status: "failed",
      Detail: "鉴权失败，已重试",
      Author: "easy-yapi",
    },
  }

  return Store{
    Summary:    summary,
    Traffic:    traffic,
    Interfaces: interfaces,
    SyncLogs:   logs,
  }
}

func requestFields() []FieldNode {
  return []FieldNode{
    {
      Name:        "tenantId",
      Type:        "string",
      Required:    true,
      Description: "租户标识",
    },
    {
      Name: "filters",
      Type: "object",
      Children: []FieldNode{
        {Name: "env", Type: "string", Description: "环境"},
        {Name: "status", Type: "string", Description: "状态"},
        {
          Name: "timeRange",
          Type: "object",
          Children: []FieldNode{
            {Name: "from", Type: "string"},
            {Name: "to", Type: "string"},
          },
        },
      },
    },
    {
      Name: "pagination",
      Type: "object",
      Children: []FieldNode{{Name: "page", Type: "number"}, {Name: "pageSize", Type: "number"}},
    },
  }
}

func generateTraffic(points int, now time.Time) []TrafficPoint {
  if points < 2 {
    points = 2
  }
  start := now.AddDate(0, 0, -points*2)
  series := make([]TrafficPoint, 0, points)
  for i := 0; i < points; i++ {
    day := start.AddDate(0, 0, i*2)
    base := 120 + 18*math.Sin(float64(i)/2.0)
    wave := 40 * math.Sin(float64(i)/1.4)
    primary := int(base + wave + float64(i*12))
    secondary := int(base + float64(i*8))
    series = append(series, TrafficPoint{
      Label:     day.Format("Jan 2"),
      Primary:   primary,
      Secondary: secondary,
    })
  }
  return series
}
