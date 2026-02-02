import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Send, Loader2 } from 'lucide-react'
import { testApi } from '@/lib/api'

interface RequestPanelProps {
  environment?: Record<string, string>
  onRequestSent?: () => void
}

export default function RequestPanel({ environment = {}, onRequestSent }: RequestPanelProps) {
  const [method, setMethod] = useState('GET')
  const [url, setUrl] = useState('')
  const [headers, setHeaders] = useState<Record<string, string>>({})
  const [body, setBody] = useState('')
  const [isLoading, setIsLoading] = useState(false)

  const methods = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'OPTIONS', 'HEAD']

  // Replace environment variables in URL
  const replaceEnvVars = (text: string): string => {
    return text.replace(/\{\{(\w+)\}\}/g, (match, key) => {
      return environment[key] || match
    })
  }

  const handleSendRequest = async () => {
    if (!url) {
      alert('Please enter a URL')
      return
    }

    setIsLoading(true)
    try {
      // Replace environment variables
      const processedUrl = replaceEnvVars(url)

      // Replace environment variables in headers
      const processedHeaders: Record<string, string> = {}
      Object.entries(headers).forEach(([key, value]) => {
        processedHeaders[key] = replaceEnvVars(value)
      })

      const response = await testApi.send({
        method,
        url: processedUrl,
        headers: processedHeaders,
        body: body || undefined,
      })

      if (response.data.success) {
        // Store response in sessionStorage for ResponseViewer to pick up
        sessionStorage.setItem('lastResponse', JSON.stringify({
          ...response.data.data,
          timestamp: new Date().toISOString(),
        }))
        onRequestSent?.()
      }
    } catch (error: any) {
      console.error('Request failed:', error)
      sessionStorage.setItem('lastResponse', JSON.stringify({
        error: error.response?.data?.error || error.message || 'Request failed',
        timestamp: new Date().toISOString(),
      }))
      onRequestSent?.()
    } finally {
      setIsLoading(false)
    }
  }

  const handleAddHeader = () => {
    setHeaders({ ...headers, '': '' })
  }

  const handleRemoveHeader = (key: string) => {
    const newHeaders = { ...headers }
    delete newHeaders[key]
    setHeaders(newHeaders)
  }

  const handleHeaderKeyChange = (oldKey: string, newKey: string) => {
    if (oldKey === newKey) return
    const value = headers[oldKey]
    const newHeaders = { ...headers }
    delete newHeaders[oldKey]
    newHeaders[newKey] = value
    setHeaders(newHeaders)
  }

  const handleHeaderValueChange = (key: string, value: string) => {
    setHeaders({ ...headers, [key]: value })
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-lg">API Request</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Method and URL */}
        <div className="flex gap-2">
          <select
            value={method}
            onChange={(e) => setMethod(e.target.value)}
            className="flex h-10 w-32 rounded-md border border-input bg-background px-3 py-2 text-sm"
          >
            {methods.map((m) => (
              <option key={m} value={m}>
                {m}
              </option>
            ))}
          </select>
          <Input
            placeholder="https://api.example.com/endpoint"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            className="flex-1"
          />
        </div>

        {/* Environment hint */}
        {Object.keys(environment).length > 0 && (
          <p className="text-xs text-gray-500">
            Use <span className="font-mono bg-gray-100 px-1 rounded">{`{{variable}}`}</span> to use environment variables. Available: {Object.keys(environment).join(', ')}
          </p>
        )}

        {/* Headers */}
        <div>
          <div className="flex justify-between items-center mb-2">
            <Label>Headers</Label>
            <Button size="sm" variant="outline" onClick={handleAddHeader}>
              + Add Header
            </Button>
          </div>
          <div className="space-y-2">
            {Object.entries(headers).map(([key, value]) => (
              <div key={key} className="flex gap-2 items-center">
                <Input
                  placeholder="Header name"
                  value={key}
                  onChange={(e) => handleHeaderKeyChange(key, e.target.value)}
                  className="flex-1"
                />
                <span className="text-gray-400">:</span>
                <Input
                  placeholder="Header value"
                  value={value}
                  onChange={(e) => handleHeaderValueChange(key, e.target.value)}
                  className="flex-1"
                />
                <Button
                  size="sm"
                  variant="ghost"
                  onClick={() => handleRemoveHeader(key)}
                  className="text-destructive"
                >
                  ×
                </Button>
              </div>
            ))}
            {Object.keys(headers).length === 0 && (
              <p className="text-sm text-gray-500 text-center py-2">No headers. Click "+ Add Header" to add one.</p>
            )}
          </div>
        </div>

        {/* Body */}
        {['POST', 'PUT', 'PATCH'].includes(method) && (
          <div>
            <Label>Body (JSON)</Label>
            <textarea
              placeholder='{"key": "value"}'
              value={body}
              onChange={(e) => setBody(e.target.value)}
              rows={6}
              className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm font-mono"
            />
          </div>
        )}

        {/* Send Button */}
        <Button
          onClick={handleSendRequest}
          disabled={isLoading || !url}
          className="w-full"
        >
          {isLoading ? (
            <>
              <Loader2 className="w-4 h-4 mr-2 animate-spin" />
              Sending...
            </>
          ) : (
            <>
              <Send className="w-4 h-4 mr-2" />
              Send Request
            </>
          )}
        </Button>
      </CardContent>
    </Card>
  )
}
