import React, { useEffect, useRef, useState } from 'react';
import { getApiUrl } from '../utils/api-config';
import Papa from 'papaparse';
import { Row } from './App';
import { parseNumber, inferBoolean } from './util';

interface Props { onLoadRows: (rows: Row[]) => void }
interface JobFile { name: string; size: number; mtime: number }

const numericHints = new Set([
  'kills','deaths','assists','kda','kp','cs','cs_per_min','gold','gpm','dmg_to_champs','dpm','vision_score','duration_s','gold_share','dmg_share','vision_share','cs10','gold10','xp10','kills_per_min','wards_placed','wards_killed'
]);

export const Exporter: React.FC<Props> = ({ onLoadRows }) => {
  const [platform, setPlatform] = useState(()=> localStorage.getItem('export_platform') || 'euw1');
  const [riotId, setRiotId] = useState(()=> localStorage.getItem('export_riotId') || '');
  const COUNT_MAX = 1000;
  const [queues, setQueues] = useState<string>('420,440');
  const [lang, setLang] = useState('fr_FR');
  const [timeline, setTimeline] = useState(true);
  const [all, setAll] = useState(true);
  const [light, setLight] = useState(true);
  const [jobId, setJobId] = useState<string>('');
  const [lines, setLines] = useState<string[]>([]);
  const [status, setStatus] = useState<string>('idle');
  const [exitCode, setExitCode] = useState<number | null>(null);
  const [files, setFiles] = useState<JobFile[]>([]);
  const [progress, setProgress] = useState<number>(0);
  const [useSSE, setUseSSE] = useState<boolean>(false);
  const [jobs, setJobs] = useState<any[]>([]);
  const [showJobs, setShowJobs] = useState(false);

  const fetchJobs = async (prune=false) => {
    try {
      const url = `${API}/jobs${prune? '?prune=true':''}`;
      const r = await fetch(url);
      if (!r.ok) return;
      const j = await r.json();
      setJobs(j.jobs||[]);
    } catch {}
  };
  useEffect(()=> { if (showJobs) fetchJobs(false); }, [showJobs]);
  const lineCountRef = useRef(0);
  const API = getApiUrl('');

  useEffect(()=> {
    let t: any; let es: EventSource | null = null;
    let fallbackTimer: any;
    if (jobId && status === 'running') {
      // Try SSE first
      try {
        const esUrl = `${API}/api/jobs/${jobId}/logs`;
        es = new EventSource(esUrl);
        setUseSSE(true);
        es.onopen = () => { /* connection established */ };
        es.onmessage = (e) => {
          if (e.type === 'log') {
            try {
              const data = JSON.parse(e.data);
              const line = data.log;
              if (!line) return;
              setLines(prev => [...prev, line]);
              const m = line.match(/\b(\d+)[\/](\d+)\b/);
              if (m) { 
                const cur = parseInt(m[1],10); 
                const tot = parseInt(m[2],10); 
                if (tot>0) setProgress(Math.min(100, Math.round(cur*100/tot))); 
              }
            } catch {}
          }
        };
        es.addEventListener('heartbeat', (e: any) => {
          // Keep connection alive
        });
        es.onerror = () => {
          if (es) es.close();
          setUseSSE(false);
        };
      } catch {
        setUseSSE(false);
      }
      
      // Fallback quick check: if après 500ms aucune ligne via SSE, lancer polling
      fallbackTimer = setTimeout(()=> {
        if (lines.length === 0) {
          setUseSSE(false); // déclenchera polling
        }
      }, 500);
      
      // Fallback polling if SSE not active
      const poll = async () => {
        if (useSSE) return; // skip if SSE running
        try {
          const r = await fetch(`${API}/api/jobs/${jobId}`);
          if (!r.ok) return;
          const j = await r.json();
          
          setLines(j.log_lines || []);
          setProgress(j.progress || 0);
          setStatus(j.status);
          
          if (j.status === 'completed') {
            setProgress(100);
            setExitCode(0);
            if (j.zip_path) {
              setFiles([{name: 'export_bundle.zip', size: 0, mtime: Date.now()/1000}]);
            }
          } else if (j.status === 'failed') {
            setExitCode(1);
            if (j.error) {
              setLines(prev => [...prev, `Error: ${j.error}`]);
            }
          }
          
          if (j.status === 'running') t = setTimeout(poll, 1000);
        } catch {}
      };
      if (!useSSE) poll();
    }
    return ()=> { t && clearTimeout(t); fallbackTimer && clearTimeout(fallbackTimer); if (es) es.close(); };
  }, [jobId, status, useSSE, lines.length]);

  const launch = async () => {
    setLines([]); setExitCode(null); setStatus('running'); lineCountRef.current = 0;
    setProgress(0);
    try { localStorage.setItem('export_platform', platform); localStorage.setItem('export_riotId', riotId); } catch {}
    
    // Parse username and tagline from riotId
    const parts = riotId.split('#');
    if (parts.length !== 2) {
      setStatus('error'); 
      setLines(['Format incorrect: utilisez Username#Tagline']);
      return;
    }
    
    try {
      const body = { 
        username: parts[0].trim(), 
        tagline: parts[1].trim(), 
        gameCount: COUNT_MAX 
      };
      const r = await fetch(`${API}/api/export`, {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(body)});
      if (!r.ok) throw new Error('Export failed');
      const j = await r.json();
      setJobId(j.job_id);
      
      // immediate initial status fetch
      setTimeout(async () => {
        try {
          const rs = await fetch(`${API}/api/jobs/${j.job_id}`);
          if (rs.ok) {
            const js = await rs.json();
            setLines(js.log_lines||[]);
            setProgress(js.progress || 0);
            setStatus(js.status);
          }
        } catch {}
      }, 100);
    } catch (e:any) {
      setStatus('error'); setLines([e.message]);
    }
  };

  const loadCsv = async (fname: string) => {
    if (!jobId) return;
    const r = await fetch(`${API}/api/jobs/${jobId}/download`);
    if (!r.ok) return;
    const text = await r.text();
    const parsed = Papa.parse(text, {header:true, skipEmptyLines:true});
    const rows: Row[] = (parsed.data as any[]).map(r => {
      const out: Row = {} as Row;
      for (const [k,v] of Object.entries(r)) {
        if (v === '') { (out as any)[k]=null; continue; }
        if (k === 'win') { (out as any)[k]= inferBoolean(v); continue; }
        if (k === 'date') { (out as any)[k]= String(v); continue; }
        if (numericHints.has(k)) { (out as any)[k]= parseNumber(v as any); continue; }
        (out as any)[k]= v;
      }
      return out;
    });
    onLoadRows(rows);
  };

  const cancel = async () => {
    if (!jobId) return;
    try { 
      await fetch(`${API}/api/jobs/${jobId}`, {method:'DELETE'}); 
      setStatus('cancelled');
    } catch {}
  };

  return (
    <div style={{fontSize:'.65rem'}}>
  <h2 style={{marginTop:'1.2rem'}}>Export Direct (max {COUNT_MAX} matchs)</h2>
      <div style={{display:'flex', flexDirection:'column', gap:4}}>
  <input placeholder='Riot ID (Game#Tag)' value={riotId} onChange={e=>setRiotId(e.target.value)} />
        <div style={{display:'flex', gap:4}}>
          <select value={platform} onChange={e=>setPlatform(e.target.value)}>
            {['euw1','eun1','na1','kr','br1','la1','la2','oc1','jp1'].map(p=> <option key={p}>{p}</option>)}
          </select>
          <input disabled value={COUNT_MAX} style={{width:80, background:'#222', color:'#aaa'}} />
          <input value={queues} onChange={e=>setQueues(e.target.value)} style={{flex:1}} />
        </div>
        <div style={{display:'flex', gap:4}}>
          <select value={lang} onChange={e=>setLang(e.target.value)}>
            {['fr_FR','en_US','es_ES','de_DE'].map(l=> <option key={l}>{l}</option>)}
          </select>
          <label style={{display:'flex',alignItems:'center', gap:4}}><input type='checkbox' checked={timeline} onChange={e=>setTimeline(e.target.checked)} /> timeline</label>
          <label style={{display:'flex',alignItems:'center', gap:4}}><input type='checkbox' checked={all} onChange={e=>setAll(e.target.checked)} /> all</label>
          <label style={{display:'flex',alignItems:'center', gap:4}} title='Génère un CSV réduit (_light)'><input type='checkbox' checked={light} onChange={e=>setLight(e.target.checked)} /> light</label>
        </div>
        <div style={{display:'flex', gap:4}}>
          <button disabled={!riotId || status==='running'} onClick={launch}>Lancer</button>
          <button disabled={status!=='running'} onClick={cancel}>Stop</button>
        </div>
        <div>
          État: <span style={{color: status==='error'? '#f85149': status==='finished' && exitCode===0? '#3fb950': '#d29922'}}>{status}{exitCode!=null? ` (${exitCode})`:''}</span>{status==='running' && <span> – {progress}%</span>}
          {status==='running' && (
            <div className="progress-bar-wrap">
              <div className="progress-bar" style={{width: progress+'%'}} />
            </div>
          )}
        </div>
      </div>
      {lines.length>0 && (
        <div style={{marginTop:8, maxHeight:140, overflow:'auto', background:'#161b22', border:'1px solid #30363d', padding:6, fontFamily:'monospace', fontSize:10}}>
          {lines.map((l,i)=> <div key={i}>{l}</div>)}
        </div>
      )}
      {status==='finished' && exitCode===0 && files.length>0 && (
        <div style={{marginTop:8}}>
          <strong>Fichiers:</strong>
          <ul style={{listStyle:'none', padding:0, margin:0}}>
            {files.map(f=> (
              <li key={f.name} style={{display:'flex', gap:6, alignItems:'center'}}>
                <span>{f.name}</span>
                <button onClick={()=> loadCsv(f.name)}>Charger</button>
                <a style={{color:'#58a6ff'}} href={`${API}/api/jobs/${jobId}/download`} target='_blank' rel='noreferrer'>DL</a>
              </li>
            ))}
          </ul>
      {files.some(f=> f.name==='export_bundle.zip') && <div style={{marginTop:4, fontStyle:'italic'}}>Archive zip disponible (export_bundle.zip)</div>}
        </div>
      )}
      <div style={{marginTop:12}}>
        <button onClick={()=> setShowJobs(s=> !s)}>{showJobs? 'Fermer Jobs':'Jobs'}</button>
        {showJobs && (
          <div style={{marginTop:6, border:'1px solid #30363d', padding:6}}>
            <div style={{display:'flex', gap:6}}>
              <button onClick={()=> fetchJobs(false)}>Rafraîchir</button>
              <button onClick={()=> fetchJobs(true)}>Prune</button>
            </div>
            <table style={{width:'100%', fontSize:10, marginTop:6}}>
              <thead><tr><th>ID</th><th>Status</th><th>Exit</th><th>Lignes</th><th>Âge (min)</th></tr></thead>
              <tbody>
                {jobs.map(j=> {
                  const ageMin = ((Date.now()/1000)-j.created)/60;
                  return <tr key={j.id} style={{opacity: j.status==='running'?1:.7}}>
                    <td>{j.id}</td>
                    <td>{j.status}</td>
                    <td>{j.exitCode==null? '': j.exitCode}</td>
                    <td>{j.lines}</td>
                    <td>{ageMin.toFixed(1)}</td>
                  </tr>;
                })}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
};
