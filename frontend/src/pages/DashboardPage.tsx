import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { projectsApi } from '@/lib/api'
import { useAuthStore } from '@/store/auth'
import type { Project } from '@/types'
import { Plus, FolderOpen, Trash2 } from 'lucide-react'

export default function DashboardPage() {
  const navigate = useNavigate()
  const { user, clearAuth } = useAuthStore()
  const [projects, setProjects] = useState<Project[]>([])
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    loadProjects()
  }, [])

  const loadProjects = async () => {
    try {
      const response = await projectsApi.list()
      if (response.data.success) {
        setProjects(response.data.data || [])
      }
    } catch (error) {
      console.error('Failed to load projects:', error)
    } finally {
      setIsLoading(false)
    }
  }

  const handleCreateProject = async () => {
    const name = prompt('Enter project name:')
    if (!name) return

    const description = prompt('Enter project description (optional):') || ''

    try {
      const response = await projectsApi.create({ name, description })
      if (response.data.success) {
        loadProjects()
      }
    } catch (error: any) {
      alert(error.response?.data?.error || 'Failed to create project')
    }
  }

  const handleDeleteProject = async (id: string) => {
    if (!confirm('Are you sure you want to delete this project?')) return

    try {
      await projectsApi.delete(id)
      loadProjects()
    } catch (error: any) {
      alert(error.response?.data?.error || 'Failed to delete project')
    }
  }

  const handleLogout = () => {
    clearAuth()
    navigate('/login')
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
          <h1 className="text-2xl font-bold text-primary">APIHub</h1>
          <div className="flex items-center gap-4">
            <span className="text-sm text-gray-600">{user?.username}</span>
            <Button variant="outline" size="sm" onClick={handleLogout}>
              Logout
            </Button>
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="flex justify-between items-center mb-6">
          <div>
            <h2 className="text-3xl font-bold text-gray-900">Your Projects</h2>
            <p className="text-gray-600 mt-1">Manage your API collections and endpoints</p>
          </div>
          <Button onClick={handleCreateProject}>
            <Plus className="w-4 h-4 mr-2" />
            New Project
          </Button>
        </div>

        {isLoading ? (
          <div className="text-center py-12">
            <div className="text-gray-500">Loading projects...</div>
          </div>
        ) : projects.length === 0 ? (
          <Card className="text-center py-12">
            <CardContent>
              <FolderOpen className="w-16 h-16 mx-auto text-gray-400 mb-4" />
              <h3 className="text-lg font-semibold mb-2">No projects yet</h3>
              <p className="text-gray-600 mb-4">Create your first project to get started</p>
              <Button onClick={handleCreateProject}>
                <Plus className="w-4 h-4 mr-2" />
                Create Project
              </Button>
            </CardContent>
          </Card>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {projects.map((project) => (
              <Card
                key={project.id}
                className="cursor-pointer hover:shadow-lg transition-shadow"
                onClick={() => navigate(`/project/${project.id}`)}
              >
                <CardHeader>
                  <CardTitle className="flex justify-between items-start">
                    <span className="truncate">{project.name}</span>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-8 w-8"
                      onClick={(e) => {
                        e.stopPropagation()
                        handleDeleteProject(project.id)
                      }}
                    >
                      <Trash2 className="w-4 h-4 text-destructive" />
                    </Button>
                  </CardTitle>
                  <CardDescription className="line-clamp-2">
                    {project.description || 'No description'}
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="text-sm text-gray-500">
                    Created {new Date(project.created_at).toLocaleDateString()}
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </main>
    </div>
  )
}
