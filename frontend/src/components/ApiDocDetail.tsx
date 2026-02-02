import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { ChevronRight, ChevronDown, FileText, Copy, Check } from 'lucide-react'
import { useState } from 'react'
import type { APIParam, APIBody, Endpoint } from '@/types'

interface ApiDocDetailProps {
  endpoint: Endpoint
}

export default function ApiDocDetail({ endpoint }: ApiDocDetailProps) {
  const [expandedSections, setExpandedSections] = useState<Set<string>>(new Set())
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

  const toggleSection = (sectionId: string) => {
    const newExpanded = new Set(expandedSections)
    if (newExpanded.has(sectionId)) {
      newExpanded.delete(sectionId)
    } else {
      newExpanded.add(sectionId)
    }
    setExpandedSections(newExpanded)
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

  const renderParamItem = (param: APIParam, level = 0) => {
    const hasChildren = param.children && param.children.length > 0
    const sectionId = `${param.name}-${level}`

    return (
      <div key={param.name} className={`${level > 0 ? 'ml-4' : ''}`}>
        <div
          className={`flex items-start gap-2 py-2 px-3 hover:bg-gray-50 rounded transition-colors ${
            hasChildren ? 'cursor-pointer' : ''
          }`}
          onClick={() => hasChildren && toggleSection(sectionId)}
        >
          {hasChildren && (
            <span className="mt-0.5">
              {expandedSections.has(sectionId) ? (
                <ChevronDown className="w-4 h-4 text-gray-400" />
              ) : (
                <ChevronRight className="w-4 h-4 text-gray-400" />
              )}
            </span>
          )}
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 flex-wrap">
              <span className="font-medium text-sm">{param.name}</span>
              <Badge variant="outline" className="text-xs">
                {param.type}
              </Badge>
              {param.required && (
                <Badge variant="destructive" className="text-xs">
                  必填
                </Badge>
              )}
              <Badge variant="secondary" className="text-xs">
                {param.param_type}
              </Badge>
            </div>
            {param.description && (
              <p className="text-xs text-gray-600 mt-1">{param.description}</p>
            )}
            {param.default_value !== undefined && (
              <p className="text-xs text-gray-500 mt-1">
                默认值: {JSON.stringify(param.default_value)}
              </p>
            )}
          </div>
        </div>
        {hasChildren && expandedSections.has(sectionId) && (
          <div className="border-l-2 border-gray-200 ml-2">
            {param.children?.map((child) => renderParamItem(child, level + 1))}
          </div>
        )}
      </div>
    )
  }

  const renderParamsSection = (title: string, params?: APIParam[]) => {
    if (!params || params.length === 0) return null

    return (
      <div className="space-y-2">
        <h4 className="text-sm font-semibold text-gray-700">{title}</h4>
        <div className="border rounded-lg divide-y">{params.map((param) => renderParamItem(param))}</div>
      </div>
    )
  }

  const renderBodySection = (title: string, body?: APIBody) => {
    if (!body) return null

    return (
      <div className="space-y-3">
        <h4 className="text-sm font-semibold text-gray-700">{title}</h4>
        <div className="border rounded-lg p-4">
          <div className="flex items-center gap-2 mb-3">
            <Badge variant="outline">{body.type}</Badge>
            <Badge variant="secondary">{body.data_type}</Badge>
          </div>
          {body.schema && body.schema.length > 0 && (
            <div className="space-y-1">
              {body.schema.map((param) => renderParamItem(param))}
            </div>
          )}
          {body.example && (
            <div className="mt-4">
              <div className="flex items-center justify-between mb-2">
                <h5 className="text-xs font-medium text-gray-700">示例</h5>
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-6 px-2"
                  onClick={() => copyToClipboard(JSON.stringify(body.example, null, 2), 'example')}
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
              <pre className="bg-gray-50 p-3 rounded text-xs overflow-x-auto">
                {JSON.stringify(body.example, null, 2)}
              </pre>
            </div>
          )}
          {body.json_schema && (
            <div className="mt-4">
              <h5 className="text-xs font-medium text-gray-700 mb-2">JSON Schema</h5>
              <pre className="bg-gray-50 p-3 rounded text-xs overflow-x-auto">
                {body.json_schema}
              </pre>
            </div>
          )}
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {/* Header */}
      <Card>
        <CardHeader>
          <div className="flex items-start justify-between">
            <div className="flex-1">
              <div className="flex items-center gap-3 mb-2">
                <Badge className={`text-sm px-3 py-1 ${getMethodColor(endpoint.method)}`}>
                  {endpoint.method}
                </Badge>
                <CardTitle className="text-xl">{endpoint.name}</CardTitle>
              </div>
              <div className="mt-2 p-3 bg-gray-50 rounded flex items-center justify-between gap-2">
                <code className="text-sm text-gray-800 break-all flex-1">{endpoint.url}</code>
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-8 px-2 shrink-0"
                  onClick={() => copyToClipboard(endpoint.url, 'url')}
                >
                  {copiedField === 'url' ? (
                    <Check className="w-4 h-4 text-green-600" />
                  ) : (
                    <Copy className="w-4 h-4 text-gray-500" />
                  )}
                </Button>
              </div>
              {endpoint.description && (
                <p className="text-sm text-gray-600 mt-3">{endpoint.description}</p>
              )}
            </div>
          </div>
        </CardHeader>
      </Card>

      {/* Documentation Tabs */}
      <Tabs defaultValue="request" className="w-full">
        <TabsList className="grid w-full grid-cols-2">
          <TabsTrigger value="request">请求参数</TabsTrigger>
          <TabsTrigger value="response">响应参数</TabsTrigger>
        </TabsList>

        {/* Request Tab */}
        <TabsContent value="request" className="space-y-4 mt-4">
          {/* Request Params */}
          {renderParamsSection('路径参数 & 查询参数 & 请求头', endpoint.request_params)}

          {/* Request Body */}
          {renderBodySection('请求体 (Request Body)', endpoint.request_body)}

          {/* Basic Auth / Headers */}
          {endpoint.headers && Object.keys(endpoint.headers).length > 0 && (
            <div className="space-y-2">
              <h4 className="text-sm font-semibold text-gray-700">请求头</h4>
              <div className="border rounded-lg divide-y">
                {Object.entries(endpoint.headers).map(([key, value]) => (
                  <div key={key} className="flex items-center justify-between py-2 px-3">
                    <span className="text-sm font-medium">{key}</span>
                    <span className="text-sm text-gray-600">{value}</span>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* No Request Info */}
          {!endpoint.request_params?.length &&
            !endpoint.request_body &&
            (!endpoint.headers || Object.keys(endpoint.headers).length === 0) && (
              <Card>
                <CardContent className="py-8 text-center text-gray-500">
                  <FileText className="w-12 h-12 mx-auto mb-3 text-gray-300" />
                  <p>该接口没有请求参数</p>
                </CardContent>
              </Card>
            )}
        </TabsContent>

        {/* Response Tab */}
        <TabsContent value="response" className="space-y-4 mt-4">
          {/* Response Params */}
          {renderParamsSection('响应参数', endpoint.response_params)}

          {/* Response Body */}
          {renderBodySection('响应体 (Response Body)', endpoint.response_body)}

          {/* No Response Info */}
          {!endpoint.response_params?.length && !endpoint.response_body && (
            <Card>
              <CardContent className="py-8 text-center text-gray-500">
                <FileText className="w-12 h-12 mx-auto mb-3 text-gray-300" />
                <p>该接口没有响应参数说明</p>
              </CardContent>
            </Card>
          )}
        </TabsContent>
      </Tabs>

      {/* Raw Body (for backward compatibility) */}
      {endpoint.body && (
        <Card>
          <CardHeader>
            <CardTitle className="text-base">原始请求体</CardTitle>
          </CardHeader>
          <CardContent>
            <pre className="bg-gray-50 p-4 rounded text-xs overflow-x-auto">
              {endpoint.body}
            </pre>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
