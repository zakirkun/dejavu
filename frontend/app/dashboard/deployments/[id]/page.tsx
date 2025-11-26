'use client'

import { useEffect, useState, useRef } from 'react'
import { useParams } from 'next/navigation'
import { ArrowLeft, ExternalLink } from 'lucide-react'
import Link from 'next/link'
import { deployments } from '@/lib/api'
import { formatDate, getStatusColor } from '@/lib/utils'

interface Deployment {
  id: string
  project_id: string
  status: string
  subdomain: string
  image_url: string
  commit_hash: string
  build_logs: string
  created_at: string
}

export default function DeploymentDetailPage() {
  const params = useParams()
  const [deployment, setDeployment] = useState<Deployment | null>(null)
  const [logs, setLogs] = useState<string[]>([])
  const [loading, setLoading] = useState(true)
  const wsRef = useRef<WebSocket | null>(null)
  const logsEndRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    loadDeployment()
    connectWebSocket()

    return () => {
      if (wsRef.current) {
        wsRef.current.close()
      }
    }
  }, [params.id])

  useEffect(() => {
    logsEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [logs])

  const loadDeployment = async () => {
    try {
      const response = await deployments.getStatus(params.id as string)
      setDeployment(response.data)
      if (response.data.build_logs) {
        setLogs(response.data.build_logs.split('\n'))
      }
    } catch (error) {
      console.error('Failed to load deployment:', error)
    } finally {
      setLoading(false)
    }
  }

  const connectWebSocket = () => {
    try {
      const ws = deployments.connectLogs(params.id as string)
      wsRef.current = ws

      ws.onmessage = (event) => {
        const newLog = event.data
        setLogs((prev) => [...prev, newLog])
      }

      ws.onerror = (error) => {
        console.error('WebSocket error:', error)
      }

      ws.onclose = () => {
        console.log('WebSocket closed')
      }
    } catch (error) {
      console.error('Failed to connect WebSocket:', error)
    }
  }

  if (loading) {
    return (
      <div className="text-center py-12">
        <div className="animate-spin h-8 w-8 border-4 border-blue-600 border-t-transparent rounded-full mx-auto" />
      </div>
    )
  }

  if (!deployment) return null

  const deploymentUrl = `https://${deployment.subdomain}.dejavu.id`

  return (
    <div>
      <Link
        href="/dashboard"
        className="inline-flex items-center text-gray-600 hover:text-gray-900 mb-6"
      >
        <ArrowLeft className="h-4 w-4 mr-2" />
        Back to Dashboard
      </Link>

      <div className="mb-8">
        <div className="flex items-center space-x-3 mb-4">
          <h1 className="text-3xl font-bold text-gray-900">Deployment</h1>
          <span
            className={`px-4 py-2 rounded-full text-sm font-medium ${getStatusColor(
              deployment.status
            )}`}
          >
            {deployment.status}
          </span>
        </div>
        
        {deployment.status === 'ready' && (
          <a
            href={deploymentUrl}
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center text-blue-600 hover:underline"
          >
            {deploymentUrl}
            <ExternalLink className="h-4 w-4 ml-1" />
          </a>
        )}
      </div>

      <div className="grid md:grid-cols-3 gap-6 mb-8">
        <div className="bg-white rounded-xl p-6 shadow-sm">
          <h3 className="text-sm text-gray-600 mb-2">Subdomain</h3>
          <p className="text-lg font-semibold text-gray-900">
            {deployment.subdomain}
          </p>
        </div>

        <div className="bg-white rounded-xl p-6 shadow-sm">
          <h3 className="text-sm text-gray-600 mb-2">Commit</h3>
          <p className="text-lg font-mono text-gray-900">
            {deployment.commit_hash?.slice(0, 7) || 'N/A'}
          </p>
        </div>

        <div className="bg-white rounded-xl p-6 shadow-sm">
          <h3 className="text-sm text-gray-600 mb-2">Created</h3>
          <p className="text-lg text-gray-900">
            {formatDate(deployment.created_at)}
          </p>
        </div>
      </div>

      <div className="bg-white rounded-xl shadow-sm">
        <div className="p-6 border-b border-gray-200">
          <h3 className="text-lg font-semibold text-gray-900">Build Logs</h3>
        </div>
        <div className="p-6">
          <div className="bg-gray-900 rounded-lg p-4 h-96 overflow-y-auto font-mono text-sm">
            {logs.length === 0 ? (
              <p className="text-gray-500">Waiting for logs...</p>
            ) : (
              logs.map((log, i) => (
                <div key={i} className="text-green-400">
                  {log}
                </div>
              ))
            )}
            <div ref={logsEndRef} />
          </div>
        </div>
      </div>
    </div>
  )
}

