import { useEffect, useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Clock } from 'lucide-react'

interface ResponseData {
  status_code?: number
  status?: string
  headers?: Record<string, string>
  body?: string
  duration?: number
  error?: string
  timestamp?: string
}

export default function ResponseViewer() {
  const [response, setResponse] = useState<ResponseData | null>(null)

  useEffect(() => {
    // Listen for response updates via storage event (from other tabs) or custom event
    const handleResponseUpdate = () => {
      const stored = sessionStorage.getItem('lastResponse')
      if (stored) {
        try {
          setResponse(JSON.parse(stored))
        } catch {
          setResponse(null)
        }
      }
    }

    // Check for initial response
    handleResponseUpdate()

    // Listen for custom event from same tab
    window.addEventListener('response-updated', handleResponseUpdate)

    // Poll for updates (same tab scenario)
    const interval = setInterval(handleResponseUpdate, 500)

    return () => {
      window.removeEventListener('response-updated', handleResponseUpdate)
      clearInterval(interval)
    }
  }, [])

  const getStatusCodeColor = (code?: number) => {
    if (!code) return 'bg-gray-100 text-gray-700'
    if (code >= 200 && code < 300) return 'bg-green-100 text-green-700'
    if (code >= 300 && code < 400) return 'bg-yellow-100 text-yellow-700'
    if (code >= 400 && code < 500) return 'bg-orange-100 text-orange-700'
    if (code >= 500) return 'bg-red-100 text-red-700'
    return 'bg-gray-100 text-gray-700'
  }

  const formatJson = (str: string) => {
    try {
      const parsed = JSON.parse(str)
      return JSON.stringify(parsed, null, 2)
    } catch {
      return str
    }
  }

  if (!response) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Response</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-gray-500 text-center py-8">Send a request to see the response here</p>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex justify-between items-center">
          <CardTitle className="text-lg">Response</CardTitle>
          {response.duration !== undefined && (
            <div className="flex items-center gap-1 text-sm text-gray-600">
              <Clock className="w-4 h-4" />
              {response.duration}ms
            </div>
          )}
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Status Code */}
        <div className="flex items-center gap-3">
          {response.error ? (
            <Badge variant="destructive" className="px-3 py-1">Error</Badge>
          ) : (
            response.status_code && (
              <Badge className={`px-3 py-1 ${getStatusCodeColor(response.status_code)}`}>
                {response.status_code}
              </Badge>
            )
          )}
          {response.status && <span className="text-sm text-gray-600">{response.status}</span>}
        </div>

        {/* Error Message */}
        {response.error && (
          <div className="p-3 bg-red-50 border border-red-200 rounded-md">
            <p className="text-red-700 text-sm">{response.error}</p>
          </div>
        )}

        {/* Headers */}
        {response.headers && Object.keys(response.headers).length > 0 && (
          <div>
            <h3 className="text-sm font-semibold mb-2">Response Headers</h3>
            <div className="bg-gray-50 rounded-md p-3 text-sm">
              {Object.entries(response.headers).map(([key, value]) => (
                <div key={key} className="flex gap-2 py-1">
                  <span className="font-medium text-gray-700 w-1/3 truncate">{key}:</span>
                  <span className="text-gray-600 break-all">{value}</span>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Body */}
        {response.body && (
          <div>
            <h3 className="text-sm font-semibold mb-2">Response Body</h3>
            <div className="bg-gray-900 rounded-md p-3 overflow-auto max-h-96">
              <pre className="text-sm text-gray-100 font-mono whitespace-pre-wrap">
                {formatJson(response.body)}
              </pre>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
