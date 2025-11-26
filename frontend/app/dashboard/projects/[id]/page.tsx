'use client'

import { useEffect, useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import { ArrowLeft, Rocket, ExternalLink } from 'lucide-react'
import Link from 'next/link'
import { projects, deployments } from '@/lib/api'
import { formatDate, getStatusColor } from '@/lib/utils'

interface Project {
  id: string
  name: string
  repo_url: string
  build_command: string
  output_dir: string
  created_at: string
}

interface Deployment {
  id: string
  status: string
  subdomain: string
  commit_hash: string
  created_at: string
}

export default function ProjectDetailPage() {
  const params = useParams()
  const router = useRouter()
  const [project, setProject] = useState<Project | null>(null)
  const [deploymentsList, setDeploymentsList] = useState<Deployment[]>([])
  const [loading, setLoading] = useState(true)
  const [deploying, setDeploying] = useState(false)

  useEffect(() => {
    loadProject()
  }, [params.id])

  const loadProject = async () => {
    try {
      const response = await projects.get(params.id as string)
      setProject(response.data)
      // TODO: Load deployments for this project
      setDeploymentsList([])
    } catch (error) {
      console.error('Failed to load project:', error)
      router.push('/dashboard')
    } finally {
      setLoading(false)
    }
  }

  const handleDeploy = async () => {
    if (!project) return

    setDeploying(true)
    try {
      const response = await deployments.trigger(project.id)
      router.push(`/dashboard/deployments/${response.data.id}`)
    } catch (error) {
      console.error('Failed to trigger deployment:', error)
      alert('Deploy failed')
    } finally {
      setDeploying(false)
    }
  }

  if (loading) {
    return (
      <div className="text-center py-12">
        <div className="animate-spin h-8 w-8 border-4 border-blue-600 border-t-transparent rounded-full mx-auto" />
      </div>
    )
  }

  if (!project) return null

  return (
    <div>
      <Link
        href="/dashboard"
        className="inline-flex items-center text-gray-600 hover:text-gray-900 mb-6"
      >
        <ArrowLeft className="h-4 w-4 mr-2" />
        Back to Projects
      </Link>

      <div className="flex items-start justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">{project.name}</h1>
          <a
            href={project.repo_url}
            target="_blank"
            rel="noopener noreferrer"
            className="text-blue-600 hover:underline mt-2 inline-flex items-center"
          >
            {project.repo_url}
            <ExternalLink className="h-4 w-4 ml-1" />
          </a>
        </div>
        <button
          onClick={handleDeploy}
          disabled={deploying}
          className="flex items-center space-x-2 px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
        >
          <Rocket className="h-5 w-5" />
          <span>{deploying ? 'Deploying...' : 'Deploy Now'}</span>
        </button>
      </div>

      <div className="grid md:grid-cols-2 gap-6 mb-8">
        <div className="bg-white rounded-xl p-6 shadow-sm">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">
            Configuration
          </h3>
          <dl className="space-y-3">
            <div>
              <dt className="text-sm text-gray-600">Build Command</dt>
              <dd className="text-gray-900 font-mono text-sm">
                {project.build_command}
              </dd>
            </div>
            <div>
              <dt className="text-sm text-gray-600">Output Directory</dt>
              <dd className="text-gray-900 font-mono text-sm">
                {project.output_dir}
              </dd>
            </div>
            <div>
              <dt className="text-sm text-gray-600">Created</dt>
              <dd className="text-gray-900 text-sm">
                {formatDate(project.created_at)}
              </dd>
            </div>
          </dl>
        </div>

        <div className="bg-white rounded-xl p-6 shadow-sm">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">
            Quick Stats
          </h3>
          <dl className="space-y-3">
            <div>
              <dt className="text-sm text-gray-600">Total Deployments</dt>
              <dd className="text-2xl font-bold text-gray-900">
                {deploymentsList.length}
              </dd>
            </div>
            <div>
              <dt className="text-sm text-gray-600">Last Deployed</dt>
              <dd className="text-gray-900 text-sm">
                {deploymentsList.length > 0
                  ? formatDate(deploymentsList[0].created_at)
                  : 'Never'}
              </dd>
            </div>
          </dl>
        </div>
      </div>

      <div className="bg-white rounded-xl shadow-sm">
        <div className="p-6 border-b border-gray-200">
          <h3 className="text-lg font-semibold text-gray-900">
            Recent Deployments
          </h3>
        </div>
        <div className="p-6">
          {deploymentsList.length === 0 ? (
            <p className="text-center text-gray-600 py-8">
              Belum ada deployment
            </p>
          ) : (
            <div className="space-y-4">
              {deploymentsList.map((deployment) => (
                <div
                  key={deployment.id}
                  className="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:border-blue-300 transition-colors"
                >
                  <div>
                    <div className="flex items-center space-x-3 mb-2">
                      <span
                        className={`px-3 py-1 rounded-full text-xs font-medium ${getStatusColor(
                          deployment.status
                        )}`}
                      >
                        {deployment.status}
                      </span>
                      <span className="text-gray-600 text-sm">
                        {deployment.commit_hash?.slice(0, 7) || 'N/A'}
                      </span>
                    </div>
                    <p className="text-gray-600 text-sm">
                      {formatDate(deployment.created_at)}
                    </p>
                  </div>
                  <Link
                    href={`/dashboard/deployments/${deployment.id}`}
                    className="px-4 py-2 text-blue-600 hover:bg-blue-50 rounded-lg"
                  >
                    View Details
                  </Link>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

