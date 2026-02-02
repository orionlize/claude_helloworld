import { useState } from 'react'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Loader2, CheckCircle, AlertCircle } from 'lucide-react'
import { yapiApi } from '@/lib/api'

interface YAPISyncDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  projectId: string
  onSuccess?: () => void
}

export default function YAPISyncDialog({
  open,
  onOpenChange,
  projectId,
  onSuccess,
}: YAPISyncDialogProps) {
  const [yapiUrl, setYapiUrl] = useState('')
  const [token, setToken] = useState('')
  const [yapiProjectId, setYapiProjectId] = useState('')

  // Connection test state
  const [isTesting, setIsTesting] = useState(false)
  const [testResult, setTestResult] = useState<{
    success: boolean
    message: string
    projectInfo?: { name: string; categories: number; interfaces: number }
  } | null>(null)

  // Sync state
  const [isSyncing, setIsSyncing] = useState(false)
  const [syncResult, setSyncResult] = useState<{
    success: boolean
    message: string
    stats?: {
      created_collections: number
      updated_collections: number
      created_endpoints: number
      updated_endpoints: number
    }
  } | null>(null)

  const handleTestConnection = async () => {
    if (!yapiUrl || !token || !yapiProjectId) {
      setTestResult({ success: false, message: '请填写所有必填字段' })
      return
    }

    setIsTesting(true)
    setTestResult(null)

    try {
      const response = await yapiApi.test({
        yapi_url: yapiUrl,
        yapi_token: token,
        yapi_project_id: parseInt(yapiProjectId),
      })

      if (response.data.success && response.data.data) {
        setTestResult({
          success: true,
          message: '连接成功',
          projectInfo: {
            name: response.data.data.yapi_project.name,
            categories: response.data.data.total_categories,
            interfaces: response.data.data.total_interfaces,
          },
        })
      }
    } catch (error: any) {
      setTestResult({
        success: false,
        message: error.response?.data?.error || '连接失败，请检查配置',
      })
    } finally {
      setIsTesting(false)
    }
  }

  const handleSync = async () => {
    setIsSyncing(true)
    setSyncResult(null)

    try {
      const response = await yapiApi.sync(projectId, {
        yapi_url: yapiUrl,
        yapi_token: token,
        yapi_project_id: parseInt(yapiProjectId),
      })

      if (response.data.success && response.data.data) {
        setSyncResult({
          success: true,
          message: response.data.data.message,
          stats: response.data.data.stats,
        })
        onSuccess?.()
      }
    } catch (error: any) {
      setSyncResult({
        success: false,
        message: error.response?.data?.error || '同步失败',
      })
    } finally {
      setIsSyncing(false)
    }
  }

  const handleClose = () => {
    onOpenChange(false)
    // Reset state
    setTimeout(() => {
      setTestResult(null)
      setSyncResult(null)
    }, 200)
  }

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>从 YAPI 同步接口</DialogTitle>
          <DialogDescription>
            配置 YAPI 服务器信息，将接口数据同步到当前项目
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4 py-4">
          {/* YAPI URL */}
          <div>
            <Label htmlFor="yapi-url">
              YAPI 服务器地址 <span className="text-destructive">*</span>
            </Label>
            <Input
              id="yapi-url"
              placeholder="http://yapi.example.com"
              value={yapiUrl}
              onChange={(e) => setYapiUrl(e.target.value)}
              disabled={isSyncing}
            />
            <p className="text-xs text-gray-500 mt-1">
              YAPI 服务器地址，不包含路径
            </p>
          </div>

          {/* Token */}
          <div>
            <Label htmlFor="token">
              项目 Token <span className="text-destructive">*</span>
            </Label>
            <Input
              id="token"
              type="password"
              placeholder="在 YAPI 项目设置中获取"
              value={token}
              onChange={(e) => setToken(e.target.value)}
              disabled={isSyncing}
            />
            <p className="text-xs text-gray-500 mt-1">
              在 YAPI 项目设置 - token 生成中获取
            </p>
          </div>

          {/* YAPI Project ID */}
          <div>
            <Label htmlFor="yapi-project-id">
              YAPI 项目 ID <span className="text-destructive">*</span>
            </Label>
            <Input
              id="yapi-project-id"
              type="number"
              placeholder="123"
              value={yapiProjectId}
              onChange={(e) => setYapiProjectId(e.target.value)}
              disabled={isSyncing}
            />
            <p className="text-xs text-gray-500 mt-1">
              在 YAPI 项目页面 URL 中可以看到项目 ID
            </p>
          </div>

          {/* Test Connection Button */}
          <div className="flex items-center gap-2">
            <Button
              type="button"
              variant="outline"
              onClick={handleTestConnection}
              disabled={isTesting || isSyncing}
              className="flex-1"
            >
              {isTesting && <Loader2 className="w-4 h-4 mr-2 animate-spin" />}
              测试连接
            </Button>
          </div>

          {/* Test Result */}
          {testResult && (
            <div
              className={`p-3 rounded-md flex items-start gap-2 ${
                testResult.success ? 'bg-green-50 border border-green-200' : 'bg-red-50 border border-red-200'
              }`}
            >
              {testResult.success ? (
                <CheckCircle className="w-5 h-5 text-green-600 mt-0.5" />
              ) : (
                <AlertCircle className="w-5 h-5 text-red-600 mt-0.5" />
              )}
              <div className="flex-1">
                <p
                  className={`text-sm font-medium ${
                    testResult.success ? 'text-green-800' : 'text-red-800'
                  }`}
                >
                  {testResult.message}
                </p>
                {testResult.projectInfo && (
                  <div className="text-xs text-green-700 mt-1">
                    <p>项目: {testResult.projectInfo.name}</p>
                    <p>
                      分类: {testResult.projectInfo.categories} | 接口:{' '}
                      {testResult.projectInfo.interfaces}
                    </p>
                  </div>
                )}
              </div>
            </div>
          )}

          {/* Sync Result */}
          {syncResult && (
            <div
              className={`p-3 rounded-md flex items-start gap-2 ${
                syncResult.success ? 'bg-green-50 border border-green-200' : 'bg-red-50 border border-red-200'
              }`}
            >
              {syncResult.success ? (
                <CheckCircle className="w-5 h-5 text-green-600 mt-0.5" />
              ) : (
                <AlertCircle className="w-5 h-5 text-red-600 mt-0.5" />
              )}
              <div className="flex-1">
                <p
                  className={`text-sm font-medium ${
                    syncResult.success ? 'text-green-800' : 'text-red-800'
                  }`}
                >
                  {syncResult.message}
                </p>
                {syncResult.stats && (
                  <div className="text-xs text-green-700 mt-2 space-y-1">
                    <p>新增分类: {syncResult.stats.created_collections}</p>
                    <p>更新分类: {syncResult.stats.updated_collections}</p>
                    <p>新增接口: {syncResult.stats.created_endpoints}</p>
                    <p>更新接口: {syncResult.stats.updated_endpoints}</p>
                  </div>
                )}
              </div>
            </div>
          )}
        </div>

        <DialogFooter>
          <Button type="button" variant="outline" onClick={handleClose} disabled={isSyncing}>
            取消
          </Button>
          <Button
            onClick={handleSync}
            disabled={isSyncing || !testResult?.success}
          >
            {isSyncing && <Loader2 className="w-4 h-4 mr-2 animate-spin" />}
            开始同步
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
