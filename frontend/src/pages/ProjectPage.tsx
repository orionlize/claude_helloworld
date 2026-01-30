import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { projectsApi, collectionsApi, endpointsApi } from '@/lib/api'
import { useProjectStore } from '@/store/project'
import type { Collection, Endpoint } from '@/types'
import { ArrowLeft, Plus, Folder, File, Trash2, ChevronRight } from 'lucide-react'

export default function ProjectPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const {
    currentProject,
    collections,
    endpoints,
    selectedCollection,
    setCurrentProject,
    setCollections,
    setEndpoints,
    setSelectedCollection,
  } = useProjectStore()

  const [isLoading, setIsLoading] = useState(true)
  const [showEndpointForm, setShowEndpointForm] = useState(false)
  const [endpointForm, setEndpointForm] = useState({
    name: '',
    method: 'GET',
    url: '',
    body: '',
    description: '',
  })

  useEffect(() => {
    loadProject()
    loadCollections()
  }, [id])

  const loadProject = async () => {
    try {
      const response = await projectsApi.get(id!)
      if (response.data.success) {
        setCurrentProject(response.data.data!)
      }
    } catch (error) {
      console.error('Failed to load project:', error)
    } finally {
      setIsLoading(false)
    }
  }

  const loadCollections = async () => {
    try {
      const response = await collectionsApi.list(id!)
      if (response.data.success) {
        setCollections(response.data.data || [])
      }
    } catch (error) {
      console.error('Failed to load collections:', error)
    }
  }

  const loadEndpoints = async (collectionId: string) => {
    try {
      const response = await endpointsApi.list(collectionId)
      if (response.data.success) {
        setEndpoints(response.data.data || [])
      }
    } catch (error) {
      console.error('Failed to load endpoints:', error)
    }
  }

  const handleCreateCollection = async () => {
    const name = prompt('Enter collection name:')
    if (!name) return

    const description = prompt('Enter collection description (optional):') || ''

    try {
      await collectionsApi.create(id!, { name, description })
      loadCollections()
    } catch (error: any) {
      alert(error.response?.data?.error || 'Failed to create collection')
    }
  }

  const handleDeleteCollection = async (collectionId: string) => {
    if (!confirm('Delete this collection and all its endpoints?')) return
    try {
      await collectionsApi.delete(collectionId)
      loadCollections()
      if (selectedCollection === collectionId) {
        setSelectedCollection(null)
        setEndpoints([])
      }
    } catch (error: any) {
      alert(error.response?.data?.error || 'Failed to delete collection')
    }
  }

  const handleCreateEndpoint = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!selectedCollection) return

    try {
      await endpointsApi.create(selectedCollection, endpointForm)
      loadEndpoints(selectedCollection)
      setShowEndpointForm(false)
      setEndpointForm({ name: '', method: 'GET', url: '', body: '', description: '' })
    } catch (error: any) {
      alert(error.response?.data?.error || 'Failed to create endpoint')
    }
  }

  const handleDeleteEndpoint = async (endpointId: string) => {
    if (!confirm('Delete this endpoint?')) return
    try {
      await endpointsApi.delete(endpointId)
      if (selectedCollection) loadEndpoints(selectedCollection)
    } catch (error: any) {
      alert(error.response?.data?.error || 'Failed to delete endpoint')
    }
  }

  const getMethodColor = (method: string) => {
    const colors: Record<string, string> = {
      GET: 'bg-blue-100 text-blue-700',
      POST: 'bg-green-100 text-green-700',
      PUT: 'bg-yellow-100 text-yellow-700',
      PATCH: 'bg-orange-100 text-orange-700',
      DELETE: 'bg-red-100 text-red-700',
    }
    return colors[method] || 'bg-gray-100 text-gray-700'
  }

  if (isLoading) {
    return <div className="flex items-center justify-center min-h-screen">Loading...</div>
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
          <div className="flex items-center gap-4">
            <Button variant="ghost" size="icon" onClick={() => navigate('/')}>
              <ArrowLeft className="w-5 h-5" />
            </Button>
            <div>
              <h1 className="text-xl font-bold">{currentProject?.name}</h1>
              <p className="text-sm text-gray-600">{currentProject?.description}</p>
            </div>
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="grid grid-cols-12 gap-6">
          {/* Sidebar - Collections */}
          <div className="col-span-4">
            <Card>
              <CardHeader>
                <div className="flex justify-between items-center">
                  <CardTitle className="text-lg">Collections</CardTitle>
                  <Button size="sm" onClick={handleCreateCollection}>
                    <Plus className="w-4 h-4" />
                  </Button>
                </div>
              </CardHeader>
              <CardContent className="p-0">
                {collections.length === 0 ? (
                  <div className="p-4 text-center text-gray-500 text-sm">
                    No collections yet
                  </div>
                ) : (
                  <div className="divide-y">
                    {collections.map((collection) => (
                      <div
                        key={collection.id}
                        className={`flex items-center justify-between p-3 cursor-pointer hover:bg-gray-50 ${
                          selectedCollection === collection.id ? 'bg-blue-50' : ''
                        }`}
                        onClick={() => {
                          setSelectedCollection(collection.id)
                          loadEndpoints(collection.id)
                        }}
                      >
                        <div className="flex items-center gap-2 flex-1">
                          <Folder className="w-4 h-4 text-gray-400" />
                          <span className="text-sm font-medium">{collection.name}</span>
                        </div>
                        <div className="flex items-center gap-1">
                          <ChevronRight className="w-4 h-4 text-gray-400" />
                          <Button
                            variant="ghost"
                            size="icon"
                            className="h-6 w-6"
                            onClick={(e) => {
                              e.stopPropagation()
                              handleDeleteCollection(collection.id)
                            }}
                          >
                            <Trash2 className="w-3 h-3 text-destructive" />
                          </Button>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>
          </div>

          {/* Main Content - Endpoints */}
          <div className="col-span-8">
            {!selectedCollection ? (
              <Card>
                <CardContent className="py-12 text-center text-gray-500">
                  Select a collection to view endpoints
                </CardContent>
              </Card>
            ) : (
              <div className="space-y-4">
                <div className="flex justify-between items-center">
                  <h2 className="text-xl font-semibold">Endpoints</h2>
                  <Button size="sm" onClick={() => setShowEndpointForm(!showEndpointForm)}>
                    <Plus className="w-4 h-4 mr-2" />
                    New Endpoint
                  </Button>
                </div>

                {showEndpointForm && (
                  <Card>
                    <CardHeader>
                      <CardTitle className="text-lg">Create Endpoint</CardTitle>
                    </CardHeader>
                    <CardContent>
                      <form onSubmit={handleCreateEndpoint} className="space-y-4">
                        <div>
                          <Label htmlFor="name">Name</Label>
                          <Input
                            id="name"
                            value={endpointForm.name}
                            onChange={(e) =>
                              setEndpointForm({ ...endpointForm, name: e.target.value })
                            }
                            required
                          />
                        </div>
                        <div>
                          <Label htmlFor="method">Method</Label>
                          <select
                            id="method"
                            value={endpointForm.method}
                            onChange={(e) =>
                              setEndpointForm({ ...endpointForm, method: e.target.value })
                            }
                            className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2"
                          >
                            {['GET', 'POST', 'PUT', 'PATCH', 'DELETE'].map((method) => (
                              <option key={method} value={method}>
                                {method}
                              </option>
                            ))}
                          </select>
                        </div>
                        <div>
                          <Label htmlFor="url">URL</Label>
                          <Input
                            id="url"
                            placeholder="https://api.example.com/endpoint"
                            value={endpointForm.url}
                            onChange={(e) =>
                              setEndpointForm({ ...endpointForm, url: e.target.value })
                            }
                            required
                          />
                        </div>
                        <div>
                          <Label htmlFor="body">Body (JSON)</Label>
                          <textarea
                            id="body"
                            rows={4}
                            className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                            value={endpointForm.body}
                            onChange={(e) =>
                              setEndpointForm({ ...endpointForm, body: e.target.value })
                            }
                          />
                        </div>
                        <div>
                          <Label htmlFor="description">Description</Label>
                          <Input
                            id="description"
                            value={endpointForm.description}
                            onChange={(e) =>
                              setEndpointForm({ ...endpointForm, description: e.target.value })
                            }
                          />
                        </div>
                        <div className="flex gap-2">
                          <Button type="submit">Create</Button>
                          <Button
                            type="button"
                            variant="outline"
                            onClick={() => setShowEndpointForm(false)}
                          >
                            Cancel
                          </Button>
                        </div>
                      </form>
                    </CardContent>
                  </Card>
                )}

                <div className="space-y-2">
                  {endpoints.length === 0 ? (
                    <Card>
                      <CardContent className="py-8 text-center text-gray-500">
                        No endpoints in this collection
                      </CardContent>
                    </Card>
                  ) : (
                    endpoints.map((endpoint) => (
                      <Card key={endpoint.id}>
                        <CardContent className="p-4">
                          <div className="flex items-center justify-between">
                            <div className="flex items-center gap-3 flex-1">
                              <span
                                className={`px-2 py-1 text-xs font-bold rounded ${getMethodColor(
                                  endpoint.method
                                )}`}
                              >
                                {endpoint.method}
                              </span>
                              <div className="flex-1">
                                <div className="font-medium">{endpoint.name}</div>
                                <div className="text-sm text-gray-500">{endpoint.url}</div>
                              </div>
                            </div>
                            <div className="flex items-center gap-2">
                              <File className="w-4 h-4 text-gray-400" />
                              <Button
                                variant="ghost"
                                size="icon"
                                className="h-8 w-8"
                                onClick={() => handleDeleteEndpoint(endpoint.id)}
                              >
                                <Trash2 className="w-4 h-4 text-destructive" />
                              </Button>
                            </div>
                          </div>
                        </CardContent>
                      </Card>
                    ))
                  )}
                </div>
              </div>
            )}
          </div>
        </div>
      </main>
    </div>
  )
}
