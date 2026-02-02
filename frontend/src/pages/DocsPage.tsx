import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { ArrowLeft, Download, FileText } from 'lucide-react'
import api from '@/lib/api'

export default function DocsPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [content, setContent] = useState('')
  const [format, setFormat] = useState<'markdown' | 'html' | 'openapi'>('markdown')
  const [isLoading, setIsLoading] = useState(true)
  const [projectName, setProjectName] = useState('')

  useEffect(() => {
    loadDocs()
  }, [id, format])

  const loadDocs = async () => {
    setIsLoading(true)
    try {
      const response = await api.get(`/projects/${id}`)
      if (response.data.success) {
        setProjectName(response.data.data.name)
      }

      const queryParams = format === 'markdown' ? '' : `?format=${format}`
      const docsResponse = await api.get(`/projects/${id}/docs${queryParams}`, {
        headers: {
          Accept: format === 'openapi' ? 'application/json' : 'text/*',
        },
      })

      if (format === 'openapi') {
        // Format JSON for display
        setContent(JSON.stringify(docsResponse.data, null, 2))
      } else {
        setContent(docsResponse.data)
      }
    } catch (error) {
      console.error('Failed to load docs:', error)
    } finally {
      setIsLoading(false)
    }
  }

  const handleExport = async () => {
    try {
      const response = await api.get(`/projects/${id}/docs/export`, {
        responseType: 'blob',
      })

      // Create download link
      const url = window.URL.createObjectURL(new Blob([response.data]))
      const link = document.createElement('a')
      link.href = url
      link.download = `${projectName}.postman_collection.json`
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
    } catch (error) {
      console.error('Failed to export:', error)
      alert('Failed to export collection')
    }
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white border-b">
        <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
          <div className="flex items-center gap-4">
            <Button variant="ghost" size="icon" onClick={() => navigate(`/project/${id}`)}>
              <ArrowLeft className="w-5 h-5" />
            </Button>
            <div className="flex items-center gap-2">
              <FileText className="w-5 h-5 text-primary" />
              <h1 className="text-xl font-bold">API Documentation</h1>
            </div>
          </div>
          <Button onClick={handleExport} variant="outline">
            <Download className="w-4 h-4 mr-2" />
            Export Postman
          </Button>
        </div>
      </header>

      <main className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Card>
          <CardHeader>
            <div className="flex justify-between items-center">
              <CardTitle className="text-lg">Documentation</CardTitle>
              <div className="flex gap-2">
                <select
                  value={format}
                  onChange={(e) => setFormat(e.target.value as any)}
                  className="flex h-10 rounded-md border border-input bg-background px-3 py-2"
                >
                  <option value="markdown">Markdown</option>
                  <option value="html">HTML</option>
                  <option value="openapi">OpenAPI 3.0</option>
                </select>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="text-center py-8">Loading documentation...</div>
            ) : (
              <div className="prose max-w-none">
                {format === 'html' ? (
                  <iframe
                    srcDoc={content}
                    className="w-full border-0 rounded-md"
                    style={{ minHeight: '400px' }}
                    sandbox="allow-same-origin"
                  />
                ) : (
                  <pre className="bg-gray-900 text-gray-100 p-4 rounded-md overflow-auto max-h-[600px] text-sm whitespace-pre-wrap">
                    {content}
                  </pre>
                )}
              </div>
            )}
          </CardContent>
        </Card>
      </main>
    </div>
  )
}
