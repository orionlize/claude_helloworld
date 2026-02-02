import { useState, useEffect } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { environmentsApi } from '@/lib/api'
import { Plus, Trash2, Check } from 'lucide-react'
import type { Environment } from '@/types'

interface EnvironmentSelectorProps {
  projectId: string
  selectedEnvironment?: Environment
  onEnvironmentChange?: (env: Environment) => void
}

export default function EnvironmentSelector({
  projectId,
  selectedEnvironment,
  onEnvironmentChange,
}: EnvironmentSelectorProps) {
  const [environments, setEnvironments] = useState<Environment[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [showNewEnvForm, setShowNewEnvForm] = useState(false)
  const [newEnvName, setNewEnvName] = useState('')
  const [editingEnv, setEditingEnv] = useState<Environment | null>(null)

  useEffect(() => {
    loadEnvironments()
  }, [projectId])

  const loadEnvironments = async () => {
    try {
      const response = await environmentsApi.list(projectId)
      if (response.data.success) {
        setEnvironments(response.data.data || [])
        // Auto-select default environment
        const defaultEnv = response.data.data?.find((e: Environment) => e.is_default)
        if (defaultEnv && !selectedEnvironment) {
          onEnvironmentChange?.(defaultEnv)
        }
      }
    } catch (error) {
      console.error('Failed to load environments:', error)
    } finally {
      setIsLoading(false)
    }
  }

  const handleCreateEnvironment = async () => {
    if (!newEnvName.trim()) return

    try {
      const response = await environmentsApi.create(projectId, {
        name: newEnvName,
        variables: {},
        is_default: environments.length === 0,
      })
      if (response.data.success) {
        setNewEnvName('')
        setShowNewEnvForm(false)
        loadEnvironments()
      }
    } catch (error: any) {
      alert(error.response?.data?.error || 'Failed to create environment')
    }
  }

  const handleUpdateEnvironment = async () => {
    if (!editingEnv) return

    try {
      await environmentsApi.update(editingEnv.id, {
        name: editingEnv.name,
        variables: editingEnv.variables,
        is_default: editingEnv.is_default,
      })
      setEditingEnv(null)
      loadEnvironments()
    } catch (error: any) {
      alert(error.response?.data?.error || 'Failed to update environment')
    }
  }

  const handleSetAsDefault = async (env: Environment) => {
    // Unset all others first
    for (const e of environments) {
      if (e.is_default && e.id !== env.id) {
        await environmentsApi.update(e.id, {
          name: e.name,
          variables: e.variables,
          is_default: false,
        })
      }
    }

    // Set new default
    await environmentsApi.update(env.id, {
      name: env.name,
      variables: env.variables,
      is_default: true,
    })
    loadEnvironments()
  }

  const handleDeleteEnvironment = async (envId: string) => {
    if (!confirm('Delete this environment?')) return

    try {
      await environmentsApi.delete(envId)
      loadEnvironments()
    } catch (error: any) {
      alert(error.response?.data?.error || 'Failed to delete environment')
    }
  }

  if (isLoading) {
    return <div>Loading...</div>
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex justify-between items-center">
          <CardTitle className="text-lg">Environments</CardTitle>
          <Button size="sm" onClick={() => setShowNewEnvForm(!showNewEnvForm)}>
            <Plus className="w-4 h-4" />
          </Button>
        </div>
      </CardHeader>
      <CardContent className="space-y-3">
        {/* New Environment Form */}
        {showNewEnvForm && (
          <div className="flex gap-2 items-center p-2 bg-gray-50 rounded-md">
            <Input
              placeholder="Environment name"
              value={newEnvName}
              onChange={(e) => setNewEnvName(e.target.value)}
              className="flex-1"
            />
            <Button size="sm" onClick={handleCreateEnvironment}>
              <Check className="w-4 h-4" />
            </Button>
            <Button size="sm" variant="ghost" onClick={() => setShowNewEnvForm(false)}>
              ×
            </Button>
          </div>
        )}

        {/* Environment List */}
        <div className="space-y-2">
          {environments.map((env) => (
            <div
              key={env.id}
              className={`p-2 rounded-md border cursor-pointer transition-colors ${
                selectedEnvironment?.id === env.id
                  ? 'bg-primary/10 border-primary'
                  : 'hover:bg-gray-50'
              }`}
              onClick={() => {
                onEnvironmentChange?.(env)
                setEditingEnv(null)
              }}
            >
              {editingEnv?.id === env.id ? (
                <div className="space-y-2" onClick={(e) => e.stopPropagation()}>
                  <Input
                    value={editingEnv.name}
                    onChange={(e) =>
                      setEditingEnv({ ...editingEnv, name: e.target.value })
                    }
                  />
                  <div className="text-xs text-gray-500">Variables (JSON):</div>
                  <textarea
                    rows={4}
                    value={JSON.stringify(editingEnv.variables, null, 2)}
                    onChange={(e) => {
                      try {
                        const variables = JSON.parse(e.target.value)
                        setEditingEnv({ ...editingEnv, variables })
                      } catch {
                        // Invalid JSON, ignore
                      }
                    }}
                    className="w-full rounded-md border border-input bg-background px-3 py-2 text-xs font-mono"
                  />
                  <div className="flex gap-2">
                    <Button size="sm" onClick={handleUpdateEnvironment}>
                      <Check className="w-3 h-3" />
                    </Button>
                    <Button
                      size="sm"
                      variant="ghost"
                      onClick={() => setEditingEnv(null)}
                    >
                      Cancel
                    </Button>
                  </div>
                </div>
              ) : (
                <div className="flex justify-between items-center">
                  <div className="flex items-center gap-2">
                    <span className="font-medium">{env.name}</span>
                    {env.is_default && (
                      <span className="text-xs bg-primary/10 text-primary px-2 py-0.5 rounded">
                        Default
                      </span>
                    )}
                  </div>
                  <div className="flex items-center gap-1">
                    {!env.is_default && (
                      <Button
                        size="sm"
                        variant="ghost"
                        onClick={(e) => {
                          e.stopPropagation()
                          handleSetAsDefault(env)
                        }}
                      >
                        Set Default
                      </Button>
                    )}
                    <Button
                      size="sm"
                      variant="ghost"
                      onClick={(e) => {
                        e.stopPropagation()
                        setEditingEnv(env)
                      }}
                    >
                      Edit
                    </Button>
                    <Button
                      size="sm"
                      variant="ghost"
                      className="text-destructive"
                      onClick={(e) => {
                        e.stopPropagation()
                        handleDeleteEnvironment(env.id)
                      }}
                    >
                      <Trash2 className="w-3 h-3" />
                    </Button>
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>

        {environments.length === 0 && !showNewEnvForm && (
          <p className="text-sm text-gray-500 text-center py-2">
            No environments. Create one to manage variables.
          </p>
        )}

        {/* Selected Environment Variables */}
        {selectedEnvironment && !editingEnv && (
          <div className="mt-4 pt-4 border-t">
            <h4 className="text-sm font-semibold mb-2">Variables: {selectedEnvironment.name}</h4>
            <div className="bg-gray-50 rounded-md p-3 text-sm">
              {Object.keys(selectedEnvironment.variables).length === 0 ? (
                <p className="text-gray-500">No variables defined</p>
              ) : (
                <div className="space-y-1">
                  {Object.entries(selectedEnvironment.variables).map(([key, value]) => (
                    <div key={key} className="flex gap-2">
                      <span className="font-medium text-gray-700">{key}:</span>
                      <span className="text-gray-600 font-mono">{value}</span>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
