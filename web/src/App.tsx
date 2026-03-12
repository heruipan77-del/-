import { useMemo, useState } from 'react'

type Credentials = {
  ip: string
  port: number
  deviceId: string
  username: string
  password: string
  insecure: boolean
}

type GenericRow = Record<string, string>

const defaultCreds: Credentials = {
  ip: '',
  port: 8088,
  deviceId: '',
  username: 'admin',
  password: '',
  insecure: true,
}

export function App() {
  const [creds, setCreds] = useState<Credentials>(defaultCreds)
  const [lunId, setLunId] = useState('')
  const [rows, setRows] = useState<GenericRow[]>([])
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const cols = useMemo(() => (rows[0] ? Object.keys(rows[0]) : []), [rows])

  async function callApi(path: '/api/storagepools' | '/api/lun') {
    setLoading(true)
    setError('')
    setRows([])
    try {
      const payload = path === '/api/lun' ? { ...creds, lunId } : creds
      const resp = await fetch(path, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      })
      const data = await resp.json()
      if (!resp.ok) {
        setError(data.error ?? 'request failed')
        return
      }
      const content = data.data
      setRows(Array.isArray(content) ? content : [content])
    } catch (e) {
      setError(e instanceof Error ? e.message : 'unknown error')
    } finally {
      setLoading(false)
    }
  }

  return (
    <main>
      <h1>SmartOps 存储运维控制台</h1>
      <section className="panel">
        <h2>设备连接</h2>
        <div className="grid">
          <label>IP<input value={creds.ip} onChange={(e) => setCreds({ ...creds, ip: e.target.value })} /></label>
          <label>Port<input type="number" value={creds.port} onChange={(e) => setCreds({ ...creds, port: Number(e.target.value) })} /></label>
          <label>Device ID<input value={creds.deviceId} onChange={(e) => setCreds({ ...creds, deviceId: e.target.value })} /></label>
          <label>User<input value={creds.username} onChange={(e) => setCreds({ ...creds, username: e.target.value })} /></label>
          <label>Password<input type="password" value={creds.password} onChange={(e) => setCreds({ ...creds, password: e.target.value })} /></label>
          <label className="checkbox"><input type="checkbox" checked={creds.insecure} onChange={(e) => setCreds({ ...creds, insecure: e.target.checked })} /> 跳过 TLS 校验</label>
        </div>
      </section>

      <section className="panel">
        <h2>查询</h2>
        <div className="actions">
          <button onClick={() => callApi('/api/storagepools')} disabled={loading}>查询 StoragePool</button>
          <input placeholder="LUN ID" value={lunId} onChange={(e) => setLunId(e.target.value)} />
          <button onClick={() => callApi('/api/lun')} disabled={loading || !lunId}>查询 LUN</button>
        </div>
      </section>

      {error && <p className="error">错误: {error}</p>}

      {rows.length > 0 && (
        <section className="panel">
          <h2>结果</h2>
          <table>
            <thead>
              <tr>{cols.map((c) => <th key={c}>{c}</th>)}</tr>
            </thead>
            <tbody>
              {rows.map((r, idx) => (
                <tr key={idx}>{cols.map((c) => <td key={c}>{String(r[c] ?? '')}</td>)}</tr>
              ))}
            </tbody>
          </table>
        </section>
      )}
    </main>
  )
}
