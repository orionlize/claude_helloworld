package main

import (
	"encoding/json"
	"fmt"

	"apihub/internal/model"
)

// 模拟完整的数据流程:数据库 -> 后端存储层 -> API响应
func main() {
	fmt.Println("=== 模拟完整的数据库到API响应流程 ===\n")

	// 1. 模拟数据库中存储的数据(假设从YAPI同步得到)
	endpointInDB := &model.Endpoint{
		ID:          "test-endpoint-1",
		CollectionID: "collection-1",
		Name:        "获取用户列表",
		Method:      "GET",
		URL:         "/api/users",
		Description: "获取用户列表",
		Headers: map[string]string{
			"Authorization": "Bearer {token}",
		},

		// 请求参数
		RequestParams: []model.APIParam{
			{
				Name:      "page",
				Type:      "integer",
				ParamType: "query",
				Required:  false,
				Description: "页码",
				DefaultValue: 1,
			},
		},

		// 响应参数 - 模拟YAPI同步后数据库中应该有的数据
		ResponseParams: []model.APIParam{
			{
				Name:      "code",
				Type:      "integer",
				ParamType: "body",
				Required:  true,
				Description: "状态码",
			},
			{
				Name:      "message",
				Type:      "string",
				ParamType: "body",
				Required:  true,
				Description: "消息",
			},
			{
				Name:      "data",
				Type:      "object",
				ParamType: "body",
				Required:  true,
				Description: "数据列表",
				Children: []model.APIParam{
					{
						Name:      "users",
						Type:      "array",
						ParamType: "body",
						Required:  true,
						Description: "用户数组",
					},
					{
						Name:      "total",
						Type:      "integer",
						ParamType: "body",
						Required:  true,
						Description: "总数",
					},
				},
			},
		},

		// 响应体
		ResponseBody: &model.APIBody{
			Type:     "json",
			DataType: "object",
			Schema: []model.APIParam{
				{
					Name:      "code",
					Type:      "integer",
					ParamType: "body",
					Required:  true,
					Description: "状态码",
				},
				{
					Name:      "data",
					Type:      "object",
					ParamType: "body",
					Required:  true,
					Description: "数据",
				},
			},
			Example: map[string]interface{}{
				"code": 200,
				"message": "success",
				"data": map[string]interface{}{
					"users": []interface{}{},
					"total": 100,
				},
			},
		},
	}

	// 2. 模拟序列化到数据库(存储为JSONB)
	fmt.Println("=== 步骤1: 序列化到数据库 ===")
	dbJSON, _ := json.MarshalIndent(endpointInDB, "", "  ")
	fmt.Printf("存储到数据库的JSON:\n%s\n\n", string(dbJSON))

	// 3. 模拟从数据库读取并反序列化
	fmt.Println("=== 步骤2: 从数据库读取 ===")
	var endpointFromDB model.Endpoint
	if err := json.Unmarshal(dbJSON, &endpointFromDB); err != nil {
		fmt.Printf("反序列化失败: %v\n", err)
		return
	}

	// 4. 检查数据完整性
	fmt.Println("=== 步骤3: 验证数据完整性 ===")
	fmt.Printf("✓ RequestParams: %d 个参数\n", len(endpointFromDB.RequestParams))
	fmt.Printf("✓ ResponseParams: %d 个参数\n", len(endpointFromDB.ResponseParams))
	fmt.Printf("✓ ResponseBody: %v\n", endpointFromDB.ResponseBody != nil)

	if len(endpointFromDB.ResponseParams) > 0 {
		fmt.Println("\n响应参数详情:")
		for i, p := range endpointFromDB.ResponseParams {
			fmt.Printf("  %d. %s (%s) - %s\n", i+1, p.Name, p.Type, p.Description)
			if len(p.Children) > 0 {
				fmt.Printf("     子字段 (%d个):\n", len(p.Children))
				for _, c := range p.Children {
					fmt.Printf("       - %s (%s)\n", c.Name, c.Type)
				}
			}
		}
	}

	// 5. 模拟API响应(序列化为JSON发送给前端)
	fmt.Println("\n=== 步骤4: API响应(发送给前端) ===")
	apiResponseJSON, _ := json.MarshalIndent(endpointFromDB, "", "  ")
	fmt.Printf("API返回给前端的JSON:\n%s\n", string(apiResponseJSON))

	// 6. 验证关键字段是否存在
	fmt.Println("\n=== 步骤5: 验证关键字段 ===")
	var apiData map[string]interface{}
	json.Unmarshal(apiResponseJSON, &apiData)

	if _, ok := apiData["response_params"]; ok {
		fmt.Println("✅ response_params 字段存在")
		if rpArray, ok := apiData["response_params"].([]interface{}); ok {
			fmt.Printf("   包含 %d 个参数\n", len(rpArray))
		}
	} else {
		fmt.Println("❌ response_params 字段缺失")
	}

	if _, ok := apiData["response_body"]; ok {
		fmt.Println("✅ response_body 字段存在")
	} else {
		fmt.Println("❌ response_body 字段缺失")
	}
}
