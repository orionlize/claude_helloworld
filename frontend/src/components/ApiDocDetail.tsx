import { Card, CardContent } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { ChevronRight, ChevronDown, Copy, Check, FileText, Braces } from 'lucide-react'
import { useState } from 'react'
import type { APIParam, APIBody, Endpoint } from '@/types'

interface ApiDocDetailProps {
  endpoint: Endpoint
}

export default function ApiDocDetail({ endpoint }: ApiDocDetailProps) {
  const [expandedKeys, setExpandedKeys] = useState<Set<string>>(new Set())
  const [copiedField, setCopiedField] = useState<string | null>(null)

  const copyToClipboard = async (text: string, fieldId: string) => {
    try {
      await navigator.clipboard.writeText(text)
      setCopiedField(fieldId)
      setTimeout(() => setCopiedField(null), 2000)
    } catch (err) {
      console.error('Failed to copy:', err)
    }
  }

  const toggleExpand = (key: string) => {
    const newExpanded = new Set(expandedKeys)
    if (newExpanded.has(key)) {
      newExpanded.delete(key)
    } else {
      newExpanded.add(key)
    }
    setExpandedKeys(newExpanded)
  }

  const getMethodColor = (method: string) => {
    const colors: Record<string, string> = {
      GET: 'bg-blue-100 text-blue-700 border-blue-200',
      POST: 'bg-green-100 text-green-700 border-green-200',
      PUT: 'bg-yellow-100 text-yellow-700 border-yellow-200',
      PATCH: 'bg-orange-100 text-orange-700 border-orange-200',
      DELETE: 'bg-red-100 text-red-700 border-red-200',
    }
    return colors[method] || 'bg-gray-100 text-gray-700 border-gray-200'
  }

  // 递归渲染参数树
  const renderParamTree = (params: APIParam[], level = 0, parentKey = ''): React.ReactNode => {
    if (!params || params.length === 0) return null

    return params.map((param, index) => {
      const key = `${parentKey}-${param.name}-${index}`
      const hasChildren = param.children && param.children.length > 0
      const isExpanded = expandedKeys.has(key)

      return (
        <div key={key} className="select-none">
          {/* 参数行 */}
          <div
            className="flex items-center py-2 hover:bg-gray-50 transition-colors group"
            style={{ paddingLeft: `${level * 24 + 12}px` }}
          >
            {/* 展开/折叠图标 */}
            <div
              className={`w-5 h-5 flex items-center justify-center cursor-pointer ${
                hasChildren ? 'opacity-100' : 'opacity-0'
              }`}
              onClick={() => hasChildren && toggleExpand(key)}
            >
              {hasChildren ? (
                isExpanded ? (
                  <ChevronDown className="w-4 h-4 text-gray-400" />
                ) : (
                  <ChevronRight className="w-4 h-4 text-gray-400" />
                )
              ) : null}
            </div>

            {/* 参数名 */}
            <div className="flex-1 min-w-0 flex items-center gap-2">
              <span className="font-mono text-sm font-medium text-gray-800">{param.name}</span>

              {/* 类型 Badge */}
              <Badge
                variant="outline"
                className="text-xs font-normal shrink-0"
                title="字段类型"
              >
                {param.type}
              </Badge>

              {/* 必填标识 */}
              {param.required && (
                <Badge
                  variant="destructive"
                  className="text-xs px-1.5 py-0 shrink-0"
                  title="必填字段"
                >
                  必填
                </Badge>
              )}

              {/* 参数类型 Badge */}
              <Badge
                variant="secondary"
                className="text-xs font-normal shrink-0"
                title="参数位置"
              >
                {param.param_type}
              </Badge>
            </div>

            {/* 描述和默认值 */}
            <div className="flex-1 min-w-0 px-4">
              <div className="text-sm text-gray-600 truncate" title={param.description}>
                {param.description || '-'}
              </div>
              {param.default_value !== undefined && (
                <div className="text-xs text-gray-400">
                  默认: {String(param.default_value)}
                </div>
              )}
            </div>
          </div>

          {/* 子参数 */}
          {hasChildren && isExpanded && (
            <div className="border-l border-gray-200 ml-6">
              {renderParamTree(param.children!, level + 1, key)}
            </div>
          )}
        </div>
      )
    })
  }

  // 渲染请求体或响应体的 Schema
  const renderSchemaTree = (schema?: APIBody, title = 'Body') => {
    if (!schema) return null

    return (
      <div className="space-y-3">
        {/* Schema 信息头 */}
        <div className="flex items-center gap-3 px-3 py-2 bg-gray-50 rounded">
          <Braces className="w-4 h-4 text-gray-500" />
          <span className="text-sm font-medium text-gray-700">{title}</span>
          <Badge variant="outline" className="text-xs">
            {schema.type}
          </Badge>
          <Badge variant="secondary" className="text-xs">
            {schema.data_type}
          </Badge>
        </div>

        {/* Schema 树 */}
        {schema.schema && schema.schema.length > 0 ? (
          <Card>
            <CardContent className="p-0">
              <div className="py-2">{renderParamTree(schema.schema)}</div>
            </CardContent>
          </Card>
        ) : null}

        {/* 示例数据 */}
        {schema.example && (
          <div className="space-y-2">
            <div className="flex items-center justify-between px-1">
              <span className="text-xs font-medium text-gray-700">示例数据</span>
              <Button
                variant="ghost"
                size="sm"
                className="h-6 px-2 text-xs"
                onClick={() =>
                  copyToClipboard(JSON.stringify(schema.example, null, 2), 'example')
                }
              >
                {copiedField === 'example' ? (
                  <>
                    <Check className="w-3 h-3 mr-1" />
                    已复制
                  </>
                ) : (
                  <>
                    <Copy className="w-3 h-3 mr-1" />
                    复制
                  </>
                )}
              </Button>
            </div>
            <pre className="bg-gray-900 text-gray-100 p-4 rounded text-xs overflow-x-auto">
              {JSON.stringify(schema.example, null, 2)}
            </pre>
          </div>
        )}

        {/* JSON Schema */}
        {schema.json_schema && (
          <div className="space-y-2">
            <span className="text-xs font-medium text-gray-700 px-1">JSON Schema</span>
            <pre className="bg-gray-50 p-3 rounded text-xs overflow-x-auto text-gray-700">
              {schema.json_schema}
            </pre>
          </div>
        )}
      </div>
    )
  }

  const hasRequestParams =
    endpoint.request_params && endpoint.request_params.length > 0
  const hasRequestBody = endpoint.request_body
  const hasRequestHeaders = endpoint.headers && Object.keys(endpoint.headers).length > 0

  const hasResponseParams =
    endpoint.response_params && endpoint.response_params.length > 0
  const hasResponseBody = endpoint.response_body

  return (
    <div className="space-y-4">
      {/* 接口基本信息卡片 */}
      <Card className="border-2">
        <CardContent className="p-6">
          {/* 方法和名称 */}
          <div className="flex items-center gap-3 mb-4">
            <Badge className={`text-sm px-3 py-1.5 font-bold ${getMethodColor(endpoint.method)}`}>
              {endpoint.method}
            </Badge>
            <h2 className="text-xl font-bold text-gray-900">{endpoint.name}</h2>
          </div>

          {/* URL */}
          <div className="mb-4">
            <label className="text-xs font-medium text-gray-500 mb-1 block">接口地址</label>
            <div className="flex items-center gap-2 p-3 bg-gray-50 rounded-lg border">
              <code className="flex-1 text-sm text-gray-800 font-mono break-all">
                {endpoint.url}
              </code>
              <Button
                variant="ghost"
                size="sm"
                className="h-8 w-8 p-0 shrink-0"
                onClick={() => copyToClipboard(endpoint.url, 'url')}
              >
                {copiedField === 'url' ? (
                  <Check className="w-4 h-4 text-green-600" />
                ) : (
                  <Copy className="w-4 h-4 text-gray-500" />
                )}
              </Button>
            </div>
          </div>

          {/* 描述 */}
          {endpoint.description && (
            <div>
              <label className="text-xs font-medium text-gray-500 mb-1 block">接口描述</label>
              <p className="text-sm text-gray-700 bg-gray-50 p-3 rounded-lg">
                {endpoint.description}
              </p>
            </div>
          )}
        </CardContent>
      </Card>

      {/* 参数详情 Tabs */}
      <Tabs defaultValue="request" className="w-full">
        <TabsList className="grid w-full grid-cols-2">
          <TabsTrigger value="request">请求参数</TabsTrigger>
          <TabsTrigger value="response">响应参数</TabsTrigger>
        </TabsList>

        {/* 请求参数 Tab */}
        <TabsContent value="request" className="space-y-4 mt-4">
          {!hasRequestParams && !hasRequestBody && !hasRequestHeaders ? (
            <Card>
              <CardContent className="py-12 text-center text-gray-500">
                <FileText className="w-12 h-12 mx-auto mb-3 text-gray-300" />
                <p className="text-sm">该接口没有请求参数</p>
              </CardContent>
            </Card>
          ) : (
            <>
              {/* 路径参数、查询参数、请求头 */}
              {hasRequestParams && (
                <div className="space-y-2">
                  <h4 className="text-sm font-semibold text-gray-700 px-1">参数列表</h4>
                  <Card>
                    <CardContent className="p-0">
                      <div className="py-2">{renderParamTree(endpoint.request_params || [])}</div>
                    </CardContent>
                  </Card>
                </div>
              )}

              {/* 请求体 */}
              {hasRequestBody && renderSchemaTree(endpoint.request_body, '请求体 (Request Body)')}

              {/* 自定义请求头 */}
              {hasRequestHeaders && (
                <div className="space-y-2">
                  <h4 className="text-sm font-semibold text-gray-700 px-1">请求头 (Headers)</h4>
                  <Card>
                    <CardContent className="p-0">
                      <div className="divide-y">
                        {Object.entries(endpoint.headers).map(([key, value]) => (
                          <div
                            key={key}
                            className="flex items-center justify-between py-2 px-4 hover:bg-gray-50"
                          >
                            <span className="text-sm font-medium font-mono text-gray-700">{key}</span>
                            <span className="text-sm text-gray-600">{value}</span>
                          </div>
                        ))}
                      </div>
                    </CardContent>
                  </Card>
                </div>
              )}
            </>
          )}
        </TabsContent>

        {/* 响应参数 Tab */}
        <TabsContent value="response" className="space-y-4 mt-4">
          {!hasResponseParams && !hasResponseBody ? (
            <Card>
              <CardContent className="py-12 text-center text-gray-500">
                <FileText className="w-12 h-12 mx-auto mb-3 text-gray-300" />
                <p className="text-sm">该接口没有响应参数说明</p>
              </CardContent>
            </Card>
          ) : (
            <>
              {/* 响应参数 */}
              {hasResponseParams && (
                <div className="space-y-2">
                  <h4 className="text-sm font-semibold text-gray-700 px-1">响应字段</h4>
                  <Card>
                    <CardContent className="p-0">
                      <div className="py-2">{renderParamTree(endpoint.response_params || [])}</div>
                    </CardContent>
                  </Card>
                </div>
              )}

              {/* 响应体 */}
              {hasResponseBody && renderSchemaTree(endpoint.response_body, '响应体 (Response Body)')}
            </>
          )}
        </TabsContent>
      </Tabs>

      {/* 原始请求体(向后兼容) */}
      {endpoint.body && (
        <Card>
          <CardContent className="p-4">
            <h4 className="text-sm font-semibold text-gray-700 mb-3">原始请求体</h4>
            <pre className="bg-gray-50 p-3 rounded text-xs overflow-x-auto">
              {endpoint.body}
            </pre>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
