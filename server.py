"""FastAPI orchestration server for lol_match_exporter.

Features:
 - Launch export jobs (subprocess) and stream logs via SSE
 - Optional light CSV reduction (checkbox in UI)
 - Optional API key protection (env EXPORTER_API_KEY) via header x-api-key OR ?api_key= query (for EventSource)
 - Automatic zip bundling of all job files (export_bundle.zip) on success
"""

from __future__ import annotations

import uuid
import contextlib
import asyncio
import datetime as dt
import json
import os
import subprocess
import sys
import zipfile
from pathlib import Path
from typing import Dict, Optional

from fastapi import FastAPI, HTTPException, Depends, Request
from contextlib import asynccontextmanager
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import StreamingResponse
from fastapi.staticfiles import StaticFiles
from pydantic import BaseModel, Field

BASE_DIR = Path(__file__).parent
JOBS_DIR = BASE_DIR / "jobs"
JOBS_DIR.mkdir(exist_ok=True)


class ExportRequest(BaseModel):
    platform: str
    riotId: str = Field(..., alias="riotId")
    count: int = 100
    queues: Optional[list[int]] = None
    lang: str = "fr_FR"
    pretty: bool = True
    timeline: bool = False
    all: bool = False
    out: str = "ranked.csv"
    light: bool = False  # génère ranked_light.csv avec colonnes réduites


class JobState(BaseModel):
    id: str
    created: float
    status: str
    exit_code: Optional[int] = None
    out_file: Optional[str] = None
    lines: list[str] = []


jobs: Dict[str, JobState] = {}
processes: Dict[str, subprocess.Popen] = {}
log_queues: Dict[str, asyncio.Queue] = {}

API_KEY = os.getenv("EXPORTER_API_KEY") or None


async def require_api_key(request: Request):  # FastAPI dependency
    if not API_KEY:  # not enabled
        return True
    sent = (
        request.headers.get("x-api-key")
        or request.headers.get("X-Api-Key")
        or request.query_params.get("api_key")  # allow query for EventSource
    )
    if sent != API_KEY:
        raise HTTPException(401, "invalid api key")
    return True


@asynccontextmanager
async def lifespan(app: FastAPI):  # manage background prune loop
    task = asyncio.create_task(_prune_loop())
    try:
        yield
    finally:
        task.cancel()
        with contextlib.suppress(Exception):
            await task

app = FastAPI(title="LoL Match Exporter API", lifespan=lifespan)
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # adjust for production
    allow_methods=["*"],
    allow_headers=["*"],
)

# Static files will be mounted after all API routes are defined


def build_command(req: ExportRequest) -> list[str]:
    cmd = [
        sys.executable,
        "-u",            # unbuffered stdout
        "-X",
        "utf8",
        "lol_match_exporter.py",
        "--platform",
        req.platform,
        "--riot-id",
        req.riotId,
        "--count",
        str(req.count),
        "--out",
        req.out,
        "--lang",
        req.lang,
    ]
    if req.queues:
        cmd += ["--queues", *[str(q) for q in req.queues]]
    if req.pretty:
        cmd.append("--pretty")
    if req.timeline:
        cmd.append("--timeline")
    if req.all:
        cmd.append("--all")
    return cmd


@app.post("/export")
async def start_export(req: ExportRequest):
    job_id = uuid.uuid4().hex[:12]
    job_dir = JOBS_DIR / job_id
    job_dir.mkdir(parents=True, exist_ok=True)
    out_path = job_dir / req.out
    cmd = build_command(req)
    # patch output path inside command list
    try:
        cmd[cmd.index("--out") + 1] = str(out_path)
    except ValueError:
        pass
    log_path = job_dir / "stdout.log"
    state = JobState(
        id=job_id,
        created=dt.datetime.utcnow().timestamp(),
        status="running",
        out_file=str(out_path),
    )
    jobs[job_id] = state

    # Force unbuffered via env also (belt + suspenders)
    env = os.environ.copy()
    env["PYTHONUNBUFFERED"] = "1"
    env["PYTHONIOENCODING"] = "utf-8"
    # Early user feedback
    first_line = f"[job {job_id}] lancement export {req.riotId} platform={req.platform} count={req.count}"
    state.lines.append(first_line)
    state.lines.append(f"[job {job_id}] démarrage subprocess Python...")
    print(f"[SERVER] Starting job {job_id} for {req.riotId}")
    print(f"[SERVER] Command: {' '.join(cmd[:5])}...")
    try:
        proc = subprocess.Popen(
            cmd,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            text=True,
            encoding="utf-8",
            errors="replace",
            cwd=BASE_DIR,
            env=env,
        )
        print(f"[SERVER] Subprocess created with PID {proc.pid}")
        processes[job_id] = proc
        q: asyncio.Queue = asyncio.Queue()
        log_queues[job_id] = q
    except Exception as e:
        print(f"[SERVER] Failed to create subprocess: {e}")
        state.status = "error"
        state.lines.append(f"[job {job_id}] ERREUR subprocess: {e}")
        return {"job_id": job_id}
    # push initial lines to SSE backlog immediately
    try:
        q.put_nowait(first_line)
        q.put_nowait(f"[job {job_id}] démarrage subprocess Python...")
    except Exception:
        pass

    async def reader():
        import traceback
        try:
            print(f"[SERVER] Reader task started for job {job_id}")
            assert proc.stdout
            with log_path.open("w", encoding="utf-8") as f:
                line_count = 0
                for line in proc.stdout:
                    line = line.rstrip("\n")
                    line_count += 1
                    if line_count == 1:
                        print(f"[SERVER] First line received for job {job_id}: {line[:100]}")
                    state.lines.append(line)
                    f.write(line + "\n")
                    try:
                        q.put_nowait(line)
                    except Exception:
                        pass
            rc = proc.wait()
            state.exit_code = rc
            state.status = "finished" if rc == 0 else "error"
            # light file generation
            try:
                if rc == 0 and req.light and state.out_file and Path(state.out_file).exists():
                    original = Path(state.out_file)
                    light_path = original.with_name(original.stem + "_light.csv")
                    keep = {
                        "matchId",
                        "date",
                        "queue",
                        "queue_category",
                        "win",
                        "champion",
                        "lane",
                        "kills",
                        "deaths",
                        "assists",
                        "kda",
                        "kp",
                        "cs",
                        "cs_per_min",
                        "gold",
                        "gpm",
                        "dmg_to_champs",
                        "dpm",
                        "vision_score",
                        "cs10",
                        "gold10",
                        "xp10",
                        "gold_share",
                        "dmg_share",
                        "vision_share",
                    }
                    import csv

                    with original.open("r", newline="", encoding="utf-8") as fin, light_path.open(
                        "w", newline="", encoding="utf-8"
                    ) as fout:
                        reader = csv.DictReader(fin)
                        cols = [c for c in (reader.fieldnames or []) if c in keep]
                        writer = csv.DictWriter(fout, fieldnames=cols)
                        writer.writeheader()
                        for row in reader:
                            writer.writerow({k: row.get(k, "") for k in cols})
                    state.lines.append(f"[light] Fichier réduit écrit: {light_path.name}")
            except Exception as e:  # noqa: BLE001
                state.lines.append(f"[light] Erreur génération light: {e}")
                print("[EXCEPTION light]", traceback.format_exc())
            # zip bundle
            try:
                if rc == 0:
                    bundle = job_dir / "export_bundle.zip"
                    with zipfile.ZipFile(bundle, "w", compression=zipfile.ZIP_DEFLATED) as zf:
                        for p in job_dir.iterdir():
                            if p.is_file() and p.name != bundle.name:
                                zf.write(p, arcname=p.name)
                    state.lines.append(f"[zip] Archive créée: {bundle.name}")
            except Exception as e:  # noqa: BLE001
                state.lines.append(f"[zip] Erreur archive: {e}")
                print("[EXCEPTION zip]", traceback.format_exc())
            # final status event for SSE
            try:
                q.put_nowait(
                    json.dumps(
                        {
                            "__status__": state.status,
                            "exitCode": state.exit_code,
                        }
                    )
                )
            except Exception:
                pass
        except Exception as e:
            print("[EXCEPTION reader]", traceback.format_exc())
            state.status = "error"
            state.lines.append(f"[reader] Exception: {e}")

    print(f"[SERVER] Creating reader task for job {job_id}")
    reader_task = asyncio.ensure_future(reader())

    async def heartbeat():
        print(f"[SERVER] Heartbeat task starting for job {job_id}")
        await asyncio.sleep(1)  # Start heartbeat immediately
        last_count = len(state.lines)
        heartbeat_count = 0
        while state.status == "running":
            heartbeat_count += 1
            if state.status != "running":
                break
            if len(state.lines) == last_count:  # no new output
                hb = f"[heartbeat #{heartbeat_count}] {dt.datetime.utcnow().isoformat(timespec='seconds')} subprocess démarrage..."
                state.lines.append(hb)
                try:
                    q.put_nowait(hb)
                except Exception:
                    pass
                print(f"[SERVER] Heartbeat #{heartbeat_count} sent for job {job_id}")
            else:
                last_count = len(state.lines)
                print(f"[SERVER] New lines detected, last_count={last_count}")
            await asyncio.sleep(2)
        print(f"[SERVER] Heartbeat task ending for job {job_id}")
        # done
    print(f"[SERVER] Creating heartbeat task for job {job_id}")
    heartbeat_task = asyncio.ensure_future(heartbeat())
    print(f"[SERVER] Job {job_id} setup complete, returning job_id")
    return {"job_id": job_id}


@app.get("/export/{job_id}/status")
async def get_status(job_id: str, since: int = 0):
    state = jobs.get(job_id)
    if not state:
        raise HTTPException(404, "job not found")
    lines = state.lines[since:]
    files = []
    job_dir = JOBS_DIR / job_id
    if job_dir.exists():
        for p in job_dir.iterdir():
            if p.is_file():
                files.append(
                    {
                        "name": p.name,
                        "size": p.stat().st_size,
                        "mtime": p.stat().st_mtime,
                    }
                )
    return {
        "status": state.status,
        "exitCode": state.exit_code,
        "lines": lines,
        "lineCount": len(state.lines),
        "files": files,
    }


@app.get("/export/{job_id}/download/{filename}")
async def download_file(job_id: str, filename: str):
    from fastapi.responses import FileResponse

    # Check if job directory exists (don't require job to be in memory)
    job_dir = JOBS_DIR / job_id
    if not job_dir.exists():
        raise HTTPException(404, "job directory not found")
    
    file_path = job_dir / filename
    if not file_path.exists():
        raise HTTPException(404, "file not found")
    
    return FileResponse(file_path, filename=filename)


@app.delete("/export/{job_id}")
async def cancel_job(job_id: str):
    proc = processes.get(job_id)
    state = jobs.get(job_id)
    if not state:
        raise HTTPException(404, "job not found")
    if proc and proc.poll() is None:
        proc.terminate()
        try:
            proc.wait(timeout=5)
        except Exception:  # noqa: BLE001
            proc.kill()
        state.status = "cancelled"
    return {"status": state.status}


@app.get("/health")
async def health():
    return {"ok": True, "jobs": len(jobs)}


@app.get("/api/stats")
async def api_stats():
    """Get Riot API usage statistics (if enhanced module is available)."""
    try:
        from riot_api_enhanced import EnhancedRiotAPI
        # This would need to be initialized elsewhere and stored globally
        # For now, return a placeholder
        return {
            "status": "Enhanced API module available",
            "features": [
                "Adaptive rate limiting",
                "In-memory caching with TTL",
                "Seasonal data segmentation",
                "Request queue with priority",
                "Automatic retry with backoff"
            ],
            "note": "Statistics will be available when using enhanced export"
        }
    except ImportError:
        return {
            "status": "Standard API mode",
            "features": ["Basic rate limiting"],
            "note": "Enhanced module not available"
        }


@app.get("/jobs")
async def list_jobs(prune: bool = False):
    """List all jobs. If prune=true, remove finished/error/cancelled older than 6h."""
    now = dt.datetime.utcnow().timestamp()
    removed = 0
    if prune:
        to_delete = [
            jid for jid, st in jobs.items() if st.status in {"finished", "error", "cancelled"} and (now - st.created) > 6 * 3600
        ]
        for jid in to_delete:
            # attempt to remove directory
            job_dir = JOBS_DIR / jid
            try:
                if job_dir.exists():
                    for p in job_dir.iterdir():
                        try:
                            p.unlink()
                        except Exception:
                            pass
                    job_dir.rmdir()
            except Exception:
                pass
            jobs.pop(jid, None)
            processes.pop(jid, None)
            log_queues.pop(jid, None)
            removed += 1
    return {
        "count": len(jobs),
        "removed": removed,
        "jobs": [
            {
                "id": st.id,
                "created": st.created,
                "status": st.status,
                "exitCode": st.exit_code,
                "out": st.out_file,
                "lines": len(st.lines),
            }
            for st in jobs.values()
        ],
    }


def _prune_old(threshold_seconds: int = 6 * 3600) -> int:
    """Internal prune utility; returns number removed."""
    now = dt.datetime.utcnow().timestamp()
    removed = 0
    to_delete = [
        jid
        for jid, st in jobs.items()
        if st.status in {"finished", "error", "cancelled"} and (now - st.created) > threshold_seconds
    ]
    for jid in to_delete:
        job_dir = JOBS_DIR / jid
        try:
            if job_dir.exists():
                for p in job_dir.iterdir():
                    try:
                        p.unlink()
                    except Exception:
                        pass
                job_dir.rmdir()
        except Exception:
            pass
        jobs.pop(jid, None)
        processes.pop(jid, None)
        log_queues.pop(jid, None)
        removed += 1
    return removed


async def _prune_loop():
    while True:
        await asyncio.sleep(3600)
        try:
            _prune_old()
        except Exception:
            pass


@app.get("/export/{job_id}/events")
async def stream_events(job_id: str):
    state = jobs.get(job_id)
    if not state:
        raise HTTPException(404, "job not found")
    q = log_queues.get(job_id)
    if not q:
        raise HTTPException(404, "queue not found")

    async def gen():
        for line in state.lines:
            yield f"data: {line}\n\n"
        while True:
            try:
                item = await q.get()
            except asyncio.CancelledError:  # pragma: no cover - runtime cancel
                break
            if item is None:
                break
            try:
                obj = json.loads(item)
                if isinstance(obj, dict) and "__status__" in obj:
                    payload = json.dumps(
                        {"status": obj["__status__"], "exitCode": obj.get("exitCode")}
                    )
                    yield f"event: status\ndata: {payload}\n\n"
                    break
            except Exception:
                pass
            safe = item.replace("\n", " ")
            yield f"data: {safe}\n\n"

    return StreamingResponse(gen(), media_type="text/event-stream")


# Mount static files for frontend after all API routes
WEB_DIR = BASE_DIR / "web" / "dist"
if WEB_DIR.exists():
    app.mount("/static", StaticFiles(directory=WEB_DIR / "assets"), name="static")
    app.mount("/", StaticFiles(directory=WEB_DIR, html=True), name="frontend")


if __name__ == "__main__":  # pragma: no cover
    import uvicorn
    # Pour éviter que l'écriture des logs dans jobs/ déclenche un reload intempestif
    # on désactive reload par défaut. Activer avec RELOAD=1 si nécessaire.
    reload_enabled = os.getenv("RELOAD") == "1"
    uvicorn.run(
        "server:app",
        host="0.0.0.0",
        port=8000,
        reload=reload_enabled,
        reload_includes=["*.py"] if reload_enabled else None,
        reload_excludes=["jobs/*", "*.log", "*.csv", "*.zip"] if reload_enabled else None,
    )
