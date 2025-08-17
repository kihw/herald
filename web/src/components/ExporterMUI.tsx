import React, { useEffect, useRef, useState } from 'react';
import {
  Card,
  CardContent,
  CardHeader,
  TextField,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Button,
  Box,
  Typography,
  LinearProgress,
  Alert,
  Checkbox,
  FormControlLabel,
  Grid,
  Chip,
  List,
  ListItem,
  ListItemText,
  ListItemSecondaryAction,
  IconButton,
  Collapse,
  Divider,
} from '@mui/material';
import {
  PlayArrow,
  Stop,
  Download,
  Refresh,
  ExpandMore,
  ExpandLess,
} from '@mui/icons-material';
import Papa from 'papaparse';
import { Row } from '../types';
import { parseNumber, inferBoolean } from '../utils';

interface ExporterProps {
  onLoadRows: (rows: Row[], puuid?: string) => void;
  onLoadingChange?: (loading: boolean) => void;
  onErrorChange?: (error: string | null) => void;
}

interface JobFile { 
  name: string; 
  size: number; 
  mtime: number; 
}

const numericHints = new Set([
  'kills','deaths','assists','kda','kp','cs','cs_per_min','gold','gpm','dmg_to_champs','dpm','vision_score','duration_s','gold_share','dmg_share','vision_share','cs10','gold10','xp10','kills_per_min','wards_placed','wards_killed'
]);

const queueTypes = [
  { id: '420', label: 'Solo/Duo Classé', description: 'Ranked Solo/Duo Queue' },
  { id: '440', label: 'Flex Classé', description: 'Ranked Flex Queue' },
  { id: '450', label: 'ARAM', description: 'Howling Abyss' },
  { id: '400', label: 'Normal 5v5', description: 'Normal Draft Pick' },
  { id: '430', label: 'Normal Aveugle', description: 'Normal Blind Pick' },
  { id: '900', label: 'URF', description: 'Ultra Rapid Fire' },
  { id: '1020', label: 'One for All', description: 'One for All' },
  { id: '1300', label: 'Nexus Blitz', description: 'Nexus Blitz' },
];

export const ExporterMUI: React.FC<ExporterProps> = ({
  onLoadRows,
  onLoadingChange,
  onErrorChange,
}) => {
  const [platform, setPlatform] = useState(() => localStorage.getItem('export_platform') || 'euw1');
  const [riotId, setRiotId] = useState(() => localStorage.getItem('export_riotId') || '');
  const COUNT_MAX = 1000;
  const [queues, setQueues] = useState<string>('420,440');
  const [selectedQueues, setSelectedQueues] = useState<Set<string>>(new Set(['420', '440']));
  const [lang, setLang] = useState('fr_FR');
  const [timeline, setTimeline] = useState(true);
  const [all, setAll] = useState(true);
  const [light, setLight] = useState(true);
  const [season, setSeason] = useState<number | null>(null);
  const [useCache, setUseCache] = useState(true);
  const [jobId, setJobId] = useState<string>('');
  const [lines, setLines] = useState<string[]>([]);
  const [status, setStatus] = useState<string>('idle');
  const [exitCode, setExitCode] = useState<number | null>(null);
  const [files, setFiles] = useState<JobFile[]>([]);
  const [progress, setProgress] = useState<number>(0);
  const [useSSE, setUseSSE] = useState<boolean>(false);
  const [jobs, setJobs] = useState<any[]>([]);
  const [showJobs, setShowJobs] = useState(false);
  const [showAdvanced, setShowAdvanced] = useState(false);

  const fetchJobs = async (prune = false) => {
    try {
      const url = `${API}/jobs${prune ? '?prune=true' : ''}`;
      const r = await fetch(url);
      if (!r.ok) return;
      const j = await r.json();
      setJobs(j.jobs || []);
    } catch {}
  };

  useEffect(() => { 
    if (showJobs) fetchJobs(false); 
  }, [showJobs]);

  const lineCountRef = useRef(0);
  const API = ((window as any).VITE_API_BASE || 'http://localhost:8000');

  const handleQueueChange = (queueId: string, checked: boolean) => {
    const newSelected = new Set(selectedQueues);
    if (checked) {
      newSelected.add(queueId);
    } else {
      newSelected.delete(queueId);
    }
    setSelectedQueues(newSelected);
    setQueues(Array.from(newSelected).join(','));
  };

  // Effet pour le polling SSE (logique existante adaptée)
  useEffect(() => {
    let t: any; 
    let es: EventSource | null = null;
    let fallbackTimer: any;
    
    if (jobId && status === 'running') {
      onLoadingChange?.(true);
      
      // Vérifier d'abord le statut du job avant de créer l'EventSource
      fetch(`${API}/export/${jobId}/status?since=0`)
        .then(r => r.json())
        .then(j => {
          if (j.status === 'running') {
            // Le job est vraiment en cours, créer l'EventSource
            try {
              const esUrl = `${API}/export/${jobId}/events`;
              es = new EventSource(esUrl);
              setUseSSE(true);
        
              es.onmessage = (e) => {
                const line = e.data;
                if (!line) return;
                setLines(prev => [...prev, line]);
                const m = line.match(/\\b(\\d+)[/](\\d+)\\b/);
                if (m) { 
                  const cur = parseInt(m[1], 10); 
                  const tot = parseInt(m[2], 10); 
                  if (tot > 0) setProgress(Math.min(100, Math.round(cur * 100 / tot))); 
                }
              };
              
              es.addEventListener('status', (e: any) => {
                try {
                  const data = JSON.parse(e.data);
                  setStatus(data.status);
                  setExitCode(data.exitCode ?? null);
                  setProgress(100);
                  onLoadingChange?.(false);
                  
                  fetch(`${API}/export/${jobId}/status?since=0`)
                    .then(r => r.json())
                    .then(j => {
                      setFiles(j.files || []);
                      const ranked = (j.files || []).find((f: JobFile) => f.name === 'ranked.csv') || 
                                    (j.files || []).find((f: JobFile) => f.name.endsWith('.csv'));
                      if (ranked) loadCsv(ranked.name);
                    });
                } catch {}
                es && es.close();
              });
              
              es.onerror = () => {
                if (es) es.close();
                setUseSSE(false);
                onLoadingChange?.(false);
              };
            } catch {
              setUseSSE(false);
              onLoadingChange?.(false);
            }
          } else {
            // Le job est déjà terminé, récupérer directement les données
            setStatus(j.status);
            setExitCode(j.exitCode ?? null);
            setProgress(100);
            setLines(j.lines || []);
            setFiles(j.files || []);
            onLoadingChange?.(false);
            
            const ranked = (j.files || []).find((f: JobFile) => f.name === 'ranked.csv') || 
                          (j.files || []).find((f: JobFile) => f.name.endsWith('.csv'));
            if (ranked) loadCsv(ranked.name);
          }
        })
        .catch(() => {
          setUseSSE(false);
          onLoadingChange?.(false);
        });
    }
    
    return () => { 
      t && clearTimeout(t); 
      fallbackTimer && clearTimeout(fallbackTimer); 
      if (es) es.close(); 
    };
  }, [jobId, status, useSSE, lines.length]);

  const launch = async () => {
    setLines([]); 
    setExitCode(null); 
    setStatus('running'); 
    lineCountRef.current = 0;
    setProgress(0);
    onErrorChange?.(null);
    
    try { 
      localStorage.setItem('export_platform', platform); 
      localStorage.setItem('export_riotId', riotId); 
    } catch {}
    
    try {
      const body: any = { 
        platform, 
        riotId, 
        count: COUNT_MAX, 
        queues: queues.split(',').map(q => parseInt(q.trim(), 10)).filter(n => !isNaN(n)), 
        lang, 
        timeline, 
        all, 
        pretty: true, 
        out: 'ranked.csv', 
        light,
        use_cache: useCache
      };
      
      if (season) {
        body.season = season;
      }
      
      const r = await fetch(`${API}/export`, {
        method: 'POST', 
        headers: { 'Content-Type': 'application/json' }, 
        body: JSON.stringify(body)
      });
      
      if (!r.ok) throw new Error('Export failed');
      const j = await r.json();
      setJobId(j.job_id);
      
      setTimeout(async () => {
        try {
          const rs = await fetch(`${API}/export/${j.job_id}/status?since=0`);
          if (rs.ok) {
            const js = await rs.json();
            setLines(js.lines || []);
            lineCountRef.current = js.lineCount || (js.lines ? js.lines.length : 0);
          }
        } catch {}
      }, 100);
    } catch (e: any) {
      setStatus('error'); 
      setLines([e.message]);
      onErrorChange?.(e.message);
      onLoadingChange?.(false);
    }
  };

  const loadCsv = async (fname: string) => {
    if (!jobId) return;
    
    try {
      const r = await fetch(`${API}/export/${jobId}/download/${fname}`);
      if (!r.ok) return;
      const text = await r.text();
      const parsed = Papa.parse(text, { header: true, skipEmptyLines: true });
      
      const rows: Row[] = (parsed.data as any[]).map(r => {
        const out: Row = {} as Row;
        for (const [k, v] of Object.entries(r)) {
          if (v === '') { (out as any)[k] = null; continue; }
          if (k === 'win') { (out as any)[k] = inferBoolean(v); continue; }
          if (k === 'date') { (out as any)[k] = String(v); continue; }
          if (numericHints.has(k)) { (out as any)[k] = parseNumber(v as any); continue; }
          (out as any)[k] = v;
        }
        return out;
      });
      
      // For now, use riotId as identifier for analytics (in production, would convert to PUUID)
      onLoadRows(rows, riotId);
    } catch (e: any) {
      onErrorChange?.(e.message);
    }
  };

  const cancel = async () => {
    if (!jobId) return;
    try { 
      await fetch(`${API}/export/${jobId}`, { method: 'DELETE' }); 
      onLoadingChange?.(false);
    } catch {}
  };

  const getStatusColor = () => {
    if (status === 'error') return 'error';
    if (status === 'finished' && exitCode === 0) return 'success';
    if (status === 'running') return 'warning';
    return 'info';
  };

  return (
    <Card>
      <CardHeader
        title="Export de données League of Legends"
        subheader={`Export direct via l'API Riot Games (max ${COUNT_MAX} matchs)`}
        action={
          <Button
            startIcon={showAdvanced ? <ExpandLess /> : <ExpandMore />}
            onClick={() => setShowAdvanced(!showAdvanced)}
            size="small"
          >
            {showAdvanced ? 'Réduire' : 'Options'}
          </Button>
        }
      />
      
      <CardContent>
        <Grid container spacing={2}>
          {/* Configuration principale */}
          <Grid item xs={12}>
            <TextField
              fullWidth
              label="Riot ID (GameName#Tag)"
              placeholder="MonPseudo#EUW"
              value={riotId}
              onChange={(e) => setRiotId(e.target.value)}
              size="small"
            />
          </Grid>
          
          <Grid item xs={12} md={4}>
            <FormControl fullWidth size="small">
              <InputLabel>Région</InputLabel>
              <Select
                value={platform}
                label="Région"
                onChange={(e) => setPlatform(e.target.value)}
              >
                {['euw1', 'eun1', 'na1', 'kr', 'br1', 'la1', 'la2', 'oc1', 'jp1'].map(p => (
                  <MenuItem key={p} value={p}>{p.toUpperCase()}</MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>
          
          <Grid item xs={12} md={8}>
            <Box>
              <Typography variant="subtitle2" sx={{ mb: 1 }}>
                Type de parties
              </Typography>
              <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                {queueTypes.map(queue => (
                  <FormControlLabel
                    key={queue.id}
                    control={
                      <Checkbox
                        size="small"
                        checked={selectedQueues.has(queue.id)}
                        onChange={(e) => handleQueueChange(queue.id, e.target.checked)}
                      />
                    }
                    label={queue.label}
                    sx={{ mr: 1 }}
                  />
                ))}
              </Box>
            </Box>
          </Grid>
        </Grid>

        {/* Options avancées */}
        <Collapse in={showAdvanced}>
          <Box sx={{ mt: 2 }}>
            <Grid container spacing={2}>
              <Grid item xs={12} md={3}>
                <FormControl fullWidth size="small">
                  <InputLabel>Langue</InputLabel>
                  <Select
                    value={lang}
                    label="Langue"
                    onChange={(e) => setLang(e.target.value)}
                  >
                    {['fr_FR', 'en_US', 'es_ES', 'de_DE'].map(l => (
                      <MenuItem key={l} value={l}>{l}</MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Grid>
              
              <Grid item xs={12} md={3}>
                <FormControl fullWidth size="small">
                  <InputLabel>Saison</InputLabel>
                  <Select
                    value={season || ''}
                    label="Saison"
                    onChange={(e) => setSeason(e.target.value ? Number(e.target.value) : null)}
                  >
                    <MenuItem value="">Toutes</MenuItem>
                    <MenuItem value="2024">2024</MenuItem>
                    <MenuItem value="2023">2023</MenuItem>
                    <MenuItem value="2022">2022</MenuItem>
                    <MenuItem value="2021">2021</MenuItem>
                  </Select>
                </FormControl>
              </Grid>
              
              <Grid item xs={12} md={6}>
                <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
                  <FormControlLabel
                    control={
                      <Checkbox 
                        checked={timeline} 
                        onChange={(e) => setTimeline(e.target.checked)} 
                      />
                    }
                    label="Timeline"
                  />
                  <FormControlLabel
                    control={
                      <Checkbox 
                        checked={all} 
                        onChange={(e) => setAll(e.target.checked)} 
                      />
                    }
                    label="Données complètes"
                  />
                  <FormControlLabel
                    control={
                      <Checkbox 
                        checked={light} 
                        onChange={(e) => setLight(e.target.checked)} 
                      />
                    }
                    label="Version allégée"
                  />
                  <FormControlLabel
                    control={
                      <Checkbox 
                        checked={useCache} 
                        onChange={(e) => setUseCache(e.target.checked)} 
                      />
                    }
                    label="Utiliser le cache"
                  />
                </Box>
              </Grid>
            </Grid>
          </Box>
        </Collapse>

        {/* Actions */}
        <Box sx={{ mt: 3, display: 'flex', gap: 1, alignItems: 'center' }}>
          <Button
            variant="contained"
            startIcon={<PlayArrow />}
            disabled={!riotId || status === 'running'}
            onClick={launch}
          >
            Lancer l'export
          </Button>
          
          <Button
            variant="outlined"
            startIcon={<Stop />}
            disabled={status !== 'running'}
            onClick={cancel}
          >
            Arrêter
          </Button>
          
          <Box sx={{ flexGrow: 1 }} />
          
          <Chip
            label={`${status}${exitCode != null ? ` (${exitCode})` : ''}`}
            color={getStatusColor()}
            variant="outlined"
          />
          
          {status === 'running' && (
            <Chip label={`${progress}%`} color="info" />
          )}
        </Box>

        {/* Progress bar */}
        {status === 'running' && (
          <Box sx={{ mt: 2 }}>
            <LinearProgress 
              variant="determinate" 
              value={progress} 
              sx={{ height: 8, borderRadius: 4 }}
            />
          </Box>
        )}

        {/* Console de logs */}
        {lines.length > 0 && (
          <Box sx={{ mt: 2 }}>
            <Typography variant="subtitle2" sx={{ mb: 1 }}>
              Journal d'export:
            </Typography>
            <Box
              sx={{
                maxHeight: 120,
                overflow: 'auto',
                backgroundColor: 'grey.900',
                color: 'grey.100',
                p: 1,
                borderRadius: 1,
                fontFamily: 'monospace',
                fontSize: '0.75rem',
              }}
            >
              {lines.map((line, i) => (
                <div key={i}>{line}</div>
              ))}
            </Box>
          </Box>
        )}

        {/* Export terminé - Message simplifié */}
        {status === 'finished' && exitCode === 0 && (
          <Alert severity="success" sx={{ mt: 2 }}>
            Export terminé avec succès ! Les données ont été traitées et sauvegardées.
          </Alert>
        )}
      </CardContent>
    </Card>
  );
};