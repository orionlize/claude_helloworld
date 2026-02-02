-- 更新现有的 endpoints 记录,添加示例的响应参数
-- 这样可以测试前端展示功能

UPDATE endpoints
SET
  request_params = '[
    {
      "name": "id",
      "type": "string",
      "param_type": "path",
      "required": true,
      "description": "资源ID",
      "default_value": null,
      "children": null
    }
  ]'::jsonb,
  response_params = '[
    {
      "name": "code",
      "type": "integer",
      "param_type": "body",
      "required": true,
      "description": "响应状态码",
      "default_value": 200,
      "children": null
    },
    {
      "name": "message",
      "type": "string",
      "param_type": "body",
      "required": true,
      "description": "响应消息",
      "default_value": null,
      "children": null
    },
    {
      "name": "data",
      "type": "object",
      "param_type": "body",
      "required": true,
      "description": "响应数据",
      "default_value": null,
      "children": [
        {
          "name": "id",
          "type": "string",
          "param_type": "body",
          "required": true,
          "description": "ID",
          "default_value": null,
          "children": null
        },
        {
          "name": "name",
          "type": "string",
          "param_type": "body",
          "required": true,
          "description": "名称",
          "default_value": null,
          "children": null
        },
        {
          "name": "created_at",
          "type": "string",
          "param_type": "body",
          "required": true,
          "description": "创建时间",
          "default_value": null,
          "children": null
        }
      ]
    }
  ]'::jsonb,
  response_body = '{
    "type": "json",
    "data_type": "object",
    "schema": [
      {
        "name": "code",
        "type": "integer",
        "param_type": "body",
        "required": true,
        "description": "状态码",
        "default_value": null,
        "children": null
      },
      {
        "name": "data",
        "type": "object",
        "param_type": "body",
        "required": true,
        "description": "数据",
        "default_value": null,
        "children": null
      }
    ],
    "example": {
      "code": 200,
      "message": "success",
      "data": {
        "id": "123",
        "name": "示例项目",
        "created_at": "2024-01-01T00:00:00Z"
      }
    },
    "json_schema": ""
  }'::jsonb
WHERE response_params IS NULL OR response_params = '[]'::jsonb;

-- 显示更新结果
SELECT
  name,
  method,
  url,
  jsonb_array_length(response_params) as response_params_count,
  response_body IS NOT NULL as has_response_body
FROM endpoints;
