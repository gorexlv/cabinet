import axios from 'axios'
import { useEffect, useState } from 'react'

interface Article {
  id: number
  title: string
  content: string
  summary: string
  tags: string
  source_url: string
  created_at: string
  updated_at: string
}

function App() {
  const [articles, setArticles] = useState<Article[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetchArticles()
  }, [])

  const fetchArticles = async () => {
    try {
      const response = await axios.get('/api/v1/articles')
      setArticles(response.data)
    } catch (err) {
      setError('Failed to fetch articles')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  if (loading) return <div>Loading...</div>
  if (error) return <div>Error: {error}</div>

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">文章列表</h1>
      <div className="grid gap-6">
        {articles.map((article) => (
          <div key={article.id} className="bg-white rounded-lg shadow-md p-6">
            <h2 className="text-xl font-semibold mb-2">{article.title}</h2>
            <p className="text-gray-600 mb-4">{article.summary}</p>
            <div className="flex flex-wrap gap-2 mb-4">
              {article.tags.split(',').map((tag, index) => (
                <span
                  key={index}
                  className="bg-blue-100 text-blue-800 text-sm px-2 py-1 rounded"
                >
                  {tag.trim()}
                </span>
              ))}
            </div>
            <div className="text-sm text-gray-500">
              <a
                href={article.source_url}
                target="_blank"
                rel="noopener noreferrer"
                className="text-blue-600 hover:underline"
              >
                原文链接
              </a>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}

export default App 