#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
import sys
import time
import argparse
import threading
from typing import Dict, Any, List, Optional, Tuple
from datetime import datetime, timezone
import json

import requests
from urllib.parse import quote
try:
    from dotenv import load_dotenv  # type: ignore
except Exception:  # pragma: no cover
    load_dotenv = None
import pandas as pd

# ------------------------------
# Constantes & tables de mapping
# ------------------------------

# Platform routing (host) pour LoL (utilisé par Summoner-V4 by-puuid)
PLATFORMS = {
    "euw1": "https://euw1.api.riotgames.com",
    "eun1": "https://eun1.api.riotgames.com",
    "tr1":  "https://tr1.api.riotgames.com",
    "ru":   "https://ru.api.riotgames.com",
    "na1":  "https://na1.api.riotgames.com",
    "la1":  "https://la1.api.riotgames.com",
    "la2":  "https://la2.api.riotgames.com",
    "br1":  "https://br1.api.riotgames.com",
    "kr":   "https://kr.api.riotgames.com",
    "jp1":  "https://jp1.api.riotgames.com",
    "oc1":  "https://oc1.api.riotgames.com",
    "ph2":  "https://ph2.api.riotgames.com",
    "sg2":  "https://sg2.api.riotgames.com",
    "th2":  "https://th2.api.riotgames.com",
    "tw2":  "https://tw2.api.riotgames.com",
    "vn2":  "https://vn2.api.riotgames.com",
    "me1":  "https://me1.api.riotgames.com",
}

# Regional routing pour MATCH-V5 (historique)
REGIONAL_MATCH_ROUTING = {
    # Europe
    "euw1": "https://europe.api.riotgames.com",
    "eun1": "https://europe.api.riotgames.com",
    "tr1":  "https://europe.api.riotgames.com",
    "ru":   "https://europe.api.riotgames.com",
    "me1":  "https://europe.api.riotgames.com",

    # Amériques
    "na1":  "https://americas.api.riotgames.com",
    "la1":  "https://americas.api.riotgames.com",
    "la2":  "https://americas.api.riotgames.com",
    "br1":  "https://americas.api.riotgames.com",

    # Asie principale
    "kr":   "https://asia.api.riotgames.com",
    "jp1":  "https://asia.api.riotgames.com",

    # Océanie & SEA -> SEA cluster pour MATCH-V5
    "oc1":  "https://sea.api.riotgames.com",
    "ph2":  "https://sea.api.riotgames.com",
    "sg2":  "https://sea.api.riotgames.com",
    "th2":  "https://sea.api.riotgames.com",
    "tw2":  "https://sea.api.riotgames.com",
    "vn2":  "https://sea.api.riotgames.com",
}

# Regional routing pour ACCOUNT-V1 (3 clusters officiels : AMERICAS / EUROPE / ASIA)
# NB: Pour les plateformes OCE & SEA, `account-v1` s'appelle sur **ASIA** (pas SEA).
REGIONAL_ACCOUNT_ROUTING = {
    # Europe
    "euw1": "https://europe.api.riotgames.com",
    "eun1": "https://europe.api.riotgames.com",
    "tr1":  "https://europe.api.riotgames.com",
    "ru":   "https://europe.api.riotgames.com",
    "me1":  "https://europe.api.riotgames.com",

    # Amériques
    "na1":  "https://americas.api.riotgames.com",
    "la1":  "https://americas.api.riotgames.com",
    "la2":  "https://americas.api.riotgames.com",
    "br1":  "https://americas.api.riotgames.com",

    # Asie (inclut OCE/SEA pour account-v1)
    "kr":   "https://asia.api.riotgames.com",
    "jp1":  "https://asia.api.riotgames.com",
    "oc1":  "https://asia.api.riotgames.com",
    "ph2":  "https://asia.api.riotgames.com",
    "sg2":  "https://asia.api.riotgames.com",
    "th2":  "https://asia.api.riotgames.com",
    "tw2":  "https://asia.api.riotgames.com",
    "vn2":  "https://asia.api.riotgames.com",
}

# Queue map minimal (sera étendu dynamiquement si possible)
QUEUE_MAP_BASE = {
    400: "Draft Normal",
    420: "Ranked Solo/Duo",
    430: "Blind Normal",
    440: "Ranked Flex",
    450: "ARAM",
    700: "Clash",
    490: "Quickplay",
}

_DYNAMIC_QUEUE_MAP: Optional[Dict[int, str]] = None

def load_queue_map() -> Dict[int, str]:
    global _DYNAMIC_QUEUE_MAP
    if _DYNAMIC_QUEUE_MAP is not None:
        return _DYNAMIC_QUEUE_MAP
    try:
        r = requests.get("https://static.developer.riotgames.com/docs/lol/queues.json", timeout=20)
        qdata = r.json()
        mapping = {}
        for entry in qdata:
            qid = entry.get("queueId")
            if isinstance(qid, int):
                desc = entry.get("description") or entry.get("map") or str(qid)
                mapping[qid] = desc
        # merge base overrides
        mapping.update(QUEUE_MAP_BASE)
        _DYNAMIC_QUEUE_MAP = mapping
    except Exception:
        _DYNAMIC_QUEUE_MAP = dict(QUEUE_MAP_BASE)
    return _DYNAMIC_QUEUE_MAP

# -------------
# Utilitaires
# -------------

def load_env_once():
    """Charge un fichier .env local (une seule fois) si python-dotenv est installé.
    Recherche dans l'ordre: chemin explicite via RIOT_DOTENV, sinon .env dans cwd.
    """
    if getattr(load_env_once, "_done", False):  # idempotent
        return
    path = os.getenv("RIOT_DOTENV") or ".env"
    if load_dotenv and os.path.exists(path):
        load_dotenv(path)
    load_env_once._done = True  # type: ignore


def get_api_key(arg_key: Optional[str] = None) -> str:
    load_env_once()
    key = arg_key or os.getenv("RIOT_API_KEY")
    if not key:
        print("Erreur: Aucune clé d'API Riot fournie. Utilisez --api-key ou définissez RIOT_API_KEY dans .env.", file=sys.stderr)
        sys.exit(1)
    return key


class RateLimiter:
    """Rate limiter combiné pour respecter 2 fenêtres:
    - 20 requêtes / 1 seconde
    - 100 requêtes / 120 secondes (2 minutes)

    Implémentation: on conserve des timestamps (epoch float) dans deux deque.
    Avant chaque requête, on purge les timestamps expirés et on dort juste assez
    pour que la requête respecte les quotas. Thread-safe pour permettre un futur parallélisme.
    """
    def __init__(self, short_max=20, short_window=1.0, long_max=100, long_window=120.0):
        from collections import deque
        self.short_max = short_max
        self.short_window = short_window
        self.long_max = long_max
        self.long_window = long_window
        self._short = deque()  # type: ignore
        self._long = deque()   # type: ignore
        self._lock = threading.Lock()

    def acquire(self):
        while True:
            with self._lock:
                now = time.time()
                # purge
                while self._short and now - self._short[0] >= self.short_window:
                    self._short.popleft()
                while self._long and now - self._long[0] >= self.long_window:
                    self._long.popleft()

                # si on peut consommer un slot
                can_short = len(self._short) < self.short_max
                can_long = len(self._long) < self.long_max
                if can_short and can_long:
                    self._short.append(now)
                    self._long.append(now)
                    return  # autorisé

                # Calcul sommeil nécessaire
                wait_short = 0.0
                wait_long = 0.0
                if not can_short and self._short:
                    wait_short = self.short_window - (now - self._short[0])
                if not can_long and self._long:
                    wait_long = self.long_window - (now - self._long[0])
                wait_time = max(wait_short, wait_long, 0.001)
            time.sleep(wait_time)


# Rate limiter global simple (créé à l'import)
RATE_LIMITER = RateLimiter()

def platform_base_url(platform: str) -> str:
    base = PLATFORMS.get(platform.lower())
    if not base:
        print(f"Plateforme inconnue: {platform}. Ex: euw1, na1, kr ...", file=sys.stderr)
        sys.exit(1)
    return base

def regional_match_base(platform: str) -> str:
    base = REGIONAL_MATCH_ROUTING.get(platform.lower())
    if not base:
        print(f"Impossible de déduire la route régionale MATCH-V5 pour {platform}.", file=sys.stderr)
        sys.exit(1)
    return base

def regional_account_base(platform: str) -> str:
    base = REGIONAL_ACCOUNT_ROUTING.get(platform.lower())
    if not base:
        print(f"Impossible de déduire la route régionale ACCOUNT-V1 pour {platform}.", file=sys.stderr)
        sys.exit(1)
    return base

def request_with_backoff(session: requests.Session, method: str, url: str, **kwargs) -> requests.Response:
    """Appels API avec rate limiting proactif + backoff sur 429/5xx.

    - Respecte quotas définis dans RATE_LIMITER avant chaque call.
    - Sur 429, lit Retry-After si présent sinon backoff exponentiel.
    - Sur 5xx réessaie avec backoff exponentiel.
    """
    max_attempts = 6
    backoff = 1.0
    resp: Optional[requests.Response] = None
    for attempt in range(1, max_attempts + 1):
        RATE_LIMITER.acquire()
        resp = session.request(method, url, timeout=25, **kwargs)
        if resp.status_code == 429:
            retry_after = resp.headers.get("Retry-After")
            sleep_s = float(retry_after) if retry_after else backoff
            time.sleep(sleep_s)
            backoff = min(backoff * 2, 16)
            continue
        if resp.status_code >= 500:
            time.sleep(backoff)
            backoff = min(backoff * 2, 16)
            continue
        return resp
    if resp is None:
        raise RuntimeError("Échec de la requête: aucune réponse obtenue après retries")
    return resp

def ts_to_local(ts_ms: int) -> datetime:
    # L'API renvoie des timestamps en ms UTC; on convertit vers le fuseau local de la machine.
    return datetime.fromtimestamp(ts_ms / 1000, tz=timezone.utc).astimezone()

def safe_div(a: float, b: float) -> float:
    return a / b if b else 0.0

def join_items(ids: List[int], id_to_name: Dict[int, str]) -> str:
    names = []
    for i in ids:
        if i and i in id_to_name:
            names.append(id_to_name[i])
        elif i:
            names.append(str(i))
    return " | ".join(names)

# ----------------------
# Data Dragon (facultatif)
# ----------------------

def load_ddragon_maps(lang: str = "en_US") -> Tuple[Dict[int, str], Dict[int, str], Dict[int, str], Dict[int, str], Dict[int, str]]:
    """
    Télécharge les mappings (dernière version) :
      - items: id -> nom
      - champions: key numérique -> nom
    - runes: id -> nom (individual runes)
    - rune styles: styleId -> nom (trees)
    - summoner spells: id -> nom
    """
    try:
        vlist = requests.get("https://ddragon.leagueoflegends.com/api/versions.json", timeout=15).json()
        version = vlist[0]
        items = requests.get(f"https://ddragon.leagueoflegends.com/cdn/{version}/data/{lang}/item.json", timeout=20).json()
        champs = requests.get(f"https://ddragon.leagueoflegends.com/cdn/{version}/data/{lang}/champion.json", timeout=20).json()
        runes = requests.get(f"https://ddragon.leagueoflegends.com/cdn/{version}/data/{lang}/runesReforged.json", timeout=20).json()
        spells = requests.get(f"https://ddragon.leagueoflegends.com/cdn/{version}/data/{lang}/summoner.json", timeout=20).json()
    except Exception as e:
        print(f"⚠️ Impossible de charger Data Dragon: {e}. Les IDs bruts seront utilisés.")
        return {}, {}, {}, {}, {}

    item_map = {int(k): v.get("name", str(k)) for k, v in items.get("data", {}).items()}
    champ_key_to_name = {}
    for cname, cdata in champs.get("data", {}).items():
        try:
            champ_key_to_name[int(cdata["key"])] = cdata["name"]
        except Exception:
            pass

    rune_map = {}
    rune_style_map = {}
    try:
        for tree in runes:
            style_id = tree.get("id")
            if isinstance(style_id, int):
                rune_style_map[style_id] = tree.get("name", str(style_id))
            for slot in tree.get("slots", []):
                for rune in slot.get("runes", []):
                    rune_map[int(rune["id"])] = rune["name"]
    except Exception:
        pass
    spell_map: Dict[int, str] = {}
    try:
        for sp_name, sp_data in spells.get("data", {}).items():
            key = sp_data.get("key")
            if key and key.isdigit():
                spell_map[int(key)] = sp_data.get("name", sp_name)
    except Exception:
        pass

    return item_map, champ_key_to_name, rune_map, rune_style_map, spell_map

# ----------------------
# Extraction des matchs
# ----------------------

def parse_riot_id(riot_id: str) -> Tuple[str, str]:
    if "#" not in riot_id:
        raise ValueError("Format Riot ID invalide. Utilisez 'GameName#TagLine'.")
    game_name, tag_line = riot_id.split("#", 1)
    if not game_name or not tag_line:
        raise ValueError("Riot ID incomplet. Exemple: 'MonPseudo#EUW'.")
    return game_name, tag_line

def get_puuid_from_riotid(session: requests.Session, platform: str, riot_id: Optional[str], game_name: Optional[str], tag_line: Optional[str], api_key: str) -> str:
    if riot_id:
        game_name, tag_line = parse_riot_id(riot_id)
    if not (game_name and tag_line):
        raise RuntimeError("Vous devez fournir --riot-id \"GameName#Tag\" OU --game-name et --tag-line, OU bien --puuid.")
    base = regional_account_base(platform)
    url = f"{base}/riot/account/v1/accounts/by-riot-id/{quote(game_name)}/{quote(tag_line)}"
    headers = {"X-Riot-Token": api_key}
    r = request_with_backoff(session, "GET", url, headers=headers)
    if r.status_code != 200:
        raise RuntimeError(f"Riot ID non trouvé ({r.status_code}): {r.text}")
    data = r.json()
    return data["puuid"]

def get_match_ids(session: requests.Session, platform: str, puuid: str, api_key: str, count: int, queues: Optional[List[int]], start_time: Optional[int], end_time: Optional[int]) -> List[str]:
    base = regional_match_base(platform)
    params = {"start": 0, "count": min(count, 100)}
    if start_time:
        params["startTime"] = start_time
    if end_time:
        params["endTime"] = end_time
    headers = {"X-Riot-Token": api_key}
    match_ids = []

    def fetch_batch(extra_params):
        url = f"{base}/lol/match/v5/matches/by-puuid/{puuid}/ids"
        r = request_with_backoff(session, "GET", url, headers=headers, params=extra_params)
        if r.status_code != 200:
            raise RuntimeError(f"Erreur match ids ({r.status_code}): {r.text}")
        return r.json()

    remaining = count
    start_idx = 0
    if queues and len(queues) > 1:
        for q in queues:
            local_params = dict(params)
            local_params["queue"] = q
            local_params["start"] = 0
            while True:
                local_params["count"] = min(remaining, 100)
                if local_params["count"] <= 0:
                    break
                batch = fetch_batch(local_params)
                if not batch:
                    break
                match_ids.extend(batch)
                remaining -= len(batch)
                local_params["start"] += len(batch)
                if len(batch) < local_params["count"]:
                    break
            if remaining <= 0:
                break
    else:
        if queues:
            params["queue"] = queues[0]
        while remaining > 0:
            params["start"] = start_idx
            params["count"] = min(remaining, 100)
            batch = fetch_batch(params)
            if not batch:
                break
            match_ids.extend(batch)
            remaining -= len(batch)
            start_idx += len(batch)
            if len(batch) < params["count"]:
                break

    match_ids = list(dict.fromkeys(match_ids))
    return match_ids[:count]

def get_match(session: requests.Session, platform: str, match_id: str, api_key: str) -> Dict[str, Any]:
    base = regional_match_base(platform)
    headers = {"X-Riot-Token": api_key}
    url = f"{base}/lol/match/v5/matches/{match_id}"
    r = request_with_backoff(session, "GET", url, headers=headers)
    if r.status_code != 200:
        raise RuntimeError(f"Erreur match {match_id} ({r.status_code}): {r.text}")
    return r.json()

def optional_summoner_profile(session: requests.Session, platform: str, puuid: str, api_key: str) -> Dict[str, Any]:
    """Essaye d'obtenir quelques métadonnées du compte via summoner-v4 by-puuid. Non bloquant."""
    try:
        base = platform_base_url(platform)
        headers = {"X-Riot-Token": api_key}
        url = f"{base}/lol/summoner/v4/summoners/by-puuid/{puuid}"
        r = request_with_backoff(session, "GET", url, headers=headers)
        if r.status_code != 200:
            return {}
        return r.json()
    except Exception:
        return {}

def extract_row_for_puuid(match: Dict[str, Any], puuid: str,
                          item_map: Dict[int, str], champ_map: Dict[int, str],
                          rune_map: Dict[int, str], rune_style_map: Dict[int, str],
                          spell_map: Dict[int, str], pretty: bool,
                          timeline_extra: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
    info = match.get("info", {})
    meta = match.get("metadata", {})
    participants = info.get("participants", [])
    me = None
    for p in participants:
        if p.get("puuid") == puuid:
            me = p
            break
    if not me:
        raise ValueError("PUUID non trouvé dans le match (custom? mode non supporté).")

    my_team_id = me.get("teamId")
    team_kills = sum(p.get("kills", 0) for p in participants if p.get("teamId") == my_team_id)

    item_ids = [me.get(f"item{i}", 0) for i in range(0, 7)]
    if pretty and item_map:
        items = join_items(item_ids, item_map)
        item_slots = [item_map.get(i, "") if i else "" for i in item_ids]
    else:
        items = " | ".join(str(i) for i in item_ids if i)
        item_slots = [str(i) if i else "" for i in item_ids]

    primary_style = me.get("perkPrimaryStyle")
    sub_style = me.get("perkSubStyle")
    primary_keystone = None
    perk_ids: List[int] = []
    try:
        styles = me.get("perks", {}).get("styles", [])
        for s in styles:
            sels = s.get("selections", [])
            for sel in sels:
                pid = sel.get("perk")
                if isinstance(pid, int):
                    perk_ids.append(pid)
            if s.get("description") == "primaryStyle" and sels:
                primary_keystone = sels[0].get("perk")
        # Stat perks
        stat_perks = me.get("perks", {}).get("statPerks", {})
    except Exception:
        stat_perks = {}

    if pretty:
        rune_primary = rune_style_map.get(primary_style, str(primary_style)) if primary_style else ""
        rune_sub = rune_style_map.get(sub_style, str(sub_style)) if sub_style else ""
        keystone = rune_map.get(primary_keystone, str(primary_keystone)) if primary_keystone else ""
        perk_names = [rune_map.get(pid, str(pid)) for pid in perk_ids]
    else:
        rune_primary = str(primary_style) if primary_style else ""
        rune_sub = str(sub_style) if sub_style else ""
        keystone = str(primary_keystone) if primary_keystone else ""
        perk_names = [str(pid) for pid in perk_ids]

    champ_id = me.get("championId")
    champ_name = champ_map.get(champ_id, me.get("championName")) if (pretty and champ_map) else me.get("championName")

    # Durée: privilégier gameStartTimestamp & gameEndTimestamp si présents (plus fiable que gameDuration).
    raw_start = info.get("gameStartTimestamp", info.get("gameCreation"))
    raw_end = info.get("gameEndTimestamp")
    if isinstance(raw_start, int) and isinstance(raw_end, int) and raw_end > raw_start:
        duration_s = (raw_end - raw_start) // 1000
    else:
        duration_s = info.get("gameDuration", 0)
    minutes = duration_s / 60 if duration_s else 0.0

    # Récupération brute des positions/roles fournis par l'API
    lane_raw = me.get("lane") or ""
    team_pos = me.get("teamPosition") or ""
    indiv_pos = me.get("individualPosition") or ""
    role_raw = me.get("role") or ""  # ex: DUO_CARRY / DUO_SUPPORT / SOLO / NONE

    # ------------------------------
    # LANE ACCURATE: On dérive un lane_std prioritaire, 100% basé sur teamPosition/lane + role_raw pour ADC/SUPPORT.
    # Règles:
    # 1. teamPosition si dans TOP/JUNGLE/MIDDLE/BOTTOM/UTILITY => base.
    # 2. Différenciation BOT: role_raw DUO_SUPPORT => SUPPORT, sinon ADC.
    # 3. Si teamPosition absent/invalide, on retombe sur lane_raw avec la même logique.
    # 4. Normalisation MIDDLE => MID, UTILITY => SUPPORT.
    # 5. Si rien: UNKNOWN.
    def derive_lane(tp: str, lr: str, rr: str) -> str:
        up_tp = (tp or "").upper()
        up_lr = (lr or "").upper()
        up_rr = (rr or "").upper()

        def classify(base: str) -> Optional[str]:  # base = TOP/JUNGLE/MIDDLE/BOTTOM/UTILITY
            if base == "TOP":
                return "TOP"
            if base == "JUNGLE":
                return "JUNGLE"
            if base in ("MIDDLE", "MID"):
                return "MID"
            if base in ("UTILITY", "SUPPORT"):
                return "SUPPORT"
            if base == "BOTTOM":
                # Distinction ADC / SUPPORT via role_raw
                if up_rr == "DUO_SUPPORT":
                    return "SUPPORT"
                if up_rr == "DUO_CARRY":
                    return "ADC"
                # Si role inconnu: heuristique -> SUPPORT si rr contient SUPPORT sinon ADC
                if "SUPPORT" in up_rr:
                    return "SUPPORT"
                return "ADC"
            return None

        # 1. teamPosition prioritaire
        lane_std = classify(up_tp) if up_tp else None
        # 2. fallback lane_raw si pas déterminé
        if not lane_std:
            lane_std = classify(up_lr) if up_lr else None
        # 3. fallback role_raw pour déduire ADC/SUPPORT seulement si on a un BOTTOM implicite
        if not lane_std:
            if up_rr in ("DUO_CARRY",):
                lane_std = "ADC"
            elif up_rr in ("DUO_SUPPORT",):
                lane_std = "SUPPORT"
        if not lane_std:
            lane_std = "UNKNOWN"
        return lane_std

    lane_std = derive_lane(team_pos, lane_raw, role_raw)
    role_std = lane_std  # on unifie: rôle standard = lane déterminé
    lane = lane_std      # on remplace la valeur lane exportée par lane_std pour assurer exactitude
    role = role_raw      # rôle brut d'API conservé pour référence

    queue_id = info.get("queueId")
    queue_name = load_queue_map().get(queue_id, str(queue_id))
    # Catégorie simplifiée
    if queue_id == 420:
        queue_cat = "solo"
    elif queue_id == 440:
        queue_cat = "flex"
    elif queue_id == 450:
        queue_cat = "aram"
    elif queue_id in (400, 430, 490):
        queue_cat = "normal"
    else:
        queue_cat = "other"

    # Patch simplifié (ex: 14.12) dérivé de gameVersion "14.12.456.1234"
    patch = None
    gv = info.get("gameVersion")
    if isinstance(gv, str):
        parts = gv.split('.')
        if len(parts) >= 2:
            patch = f"{parts[0]}.{parts[1]}"

    game_start = ts_to_local(info.get("gameStartTimestamp", info.get("gameCreation", 0)))
    game_end_dt = None
    if isinstance(raw_end, int) and raw_end > 0:
        try:
            game_end_dt = ts_to_local(raw_end)
        except Exception:
            game_end_dt = None

    kills = me.get("kills", 0)
    deaths = me.get("deaths", 0)
    assists = me.get("assists", 0)
    kda = (kills + assists) / max(1, deaths)
    kp = safe_div(kills + assists, team_kills) if team_kills else 0.0

    cs = (me.get("totalMinionsKilled", 0) + me.get("neutralMinionsKilled", 0))
    cspm = safe_div(cs, minutes)
    dchamps = me.get("totalDamageDealtToChampions", 0)
    dpm = safe_div(dchamps, minutes)
    gold = me.get("goldEarned", 0)
    gpm = safe_div(gold, minutes)
    vis = me.get("visionScore", 0)

    # Team aggregates for shares
    team_participants = [p for p in participants if p.get("teamId") == my_team_id]
    team_total_gold = sum(p.get("goldEarned", 0) for p in team_participants)
    team_total_dmg_champs = sum(p.get("totalDamageDealtToChampions", 0) for p in team_participants)
    team_total_vision = sum(p.get("visionScore", 0) for p in team_participants)

    gold_share = safe_div(gold, team_total_gold) if team_total_gold else 0.0
    dmg_share = safe_div(dchamps, team_total_dmg_champs) if team_total_dmg_champs else 0.0
    vision_share = safe_div(vis, team_total_vision) if team_total_vision else 0.0

    # Spells
    spell1 = me.get("summoner1Id")
    spell2 = me.get("summoner2Id")
    if pretty:
        spell1_name = spell_map.get(spell1, str(spell1)) if spell1 else ""
        spell2_name = spell_map.get(spell2, str(spell2)) if spell2 else ""
    else:
        spell1_name = str(spell1) if spell1 else ""
        spell2_name = str(spell2) if spell2 else ""

    # Detect mythic & boots & trinket (simple heuristics)
    mythic = ""
    boots = ""
    trinket = me.get("item6")  # vrai trinket slot (id ou nom selon pretty)
    pinks = me.get("visionWardsBoughtInGame")
    if item_map and pretty:
        for iid in item_ids:
            if not iid:
                continue
            name = item_map.get(iid, "")
            low = name.lower()
            if not mythic and ("mythic" in low or "mythique" in low):
                mythic = name
            if not boots and ("boot" in low or "bottes" in low):
                boots = name

    # Timeline extras
    cs10 = gold10 = xp10 = None
    if timeline_extra:
        cs10 = timeline_extra.get("cs10")
        gold10 = timeline_extra.get("gold10")
        xp10 = timeline_extra.get("xp10")

    remake = bool(duration_s and duration_s < 300)

    row = {
        "matchId": meta.get("matchId"),
        "gameVersion": info.get("gameVersion"),
        "patch": patch,
        "gameMode": info.get("gameMode"),
        "gameType": info.get("gameType"),
        "mapId": info.get("mapId"),
        "queueId": queue_id,
        "queue": queue_name,
        "queue_category": queue_cat,
        "date": game_start.isoformat(),
        "gameEnd": game_end_dt.isoformat() if game_end_dt else None,
        "duration_s": duration_s,
        "remake": remake,
        "win": bool(me.get("win")),
        "champion": champ_name,
        "championId": champ_id,
        "lane": lane,
        "role": role,
        "role_std": role_std,
        "summonerSpell1": spell1_name,
        "summonerSpell2": spell2_name,
        "kills": kills,
    "kills_per_min": round(safe_div(kills, minutes), 3) if minutes else 0.0,
        "deaths": deaths,
        "assists": assists,
        "kda": round(kda, 3),
        "kp": round(kp, 3),
        "cs": cs,
        "cs_per_min": round(cspm, 3),
        "gold": gold,
        "gpm": round(gpm, 3),
        "dmg_to_champs": dchamps,
        "dpm": round(dpm, 3),
        "vision_score": vis,
        "wards_placed": me.get("wardsPlaced"),
        "wards_killed": me.get("wardsKilled"),
        "control_wards_bought": pinks,
        "detectorWardsPlaced": me.get("detectorWardsPlaced"),
        "champLevel": me.get("champLevel"),
        "summonerName": me.get("summonerName"),
        "teamId": my_team_id,
        "items": items,
        "item0": item_slots[0],
        "item1": item_slots[1],
        "item2": item_slots[2],
        "item3": item_slots[3],
        "item4": item_slots[4],
        "item5": item_slots[5],
        "item6": item_slots[6],
        "primary_rune_style": rune_primary,
        "sub_rune_style": rune_sub,
        "keystone": keystone,
        "perks_list": " | ".join(perk_names),
        "perk0": perk_names[0] if len(perk_names) > 0 else "",
        "perk1": perk_names[1] if len(perk_names) > 1 else "",
        "perk2": perk_names[2] if len(perk_names) > 2 else "",
        "perk3": perk_names[3] if len(perk_names) > 3 else "",
        "perk4": perk_names[4] if len(perk_names) > 4 else "",
        "perk5": perk_names[5] if len(perk_names) > 5 else "",
        "gold_share": round(gold_share, 4),
        "dmg_share": round(dmg_share, 4),
        "vision_share": round(vision_share, 4),
        "mythic": mythic,
        "boots": boots,
        "trinket": trinket,
    "cs10": cs10,
    "gold10": gold10,
    "xp10": xp10,
    "goldDiff10": None,
    "goldDiff15": None,
    "goldDiff20": None,
    "goldDiff25": None,
        # Damage breakdown
        "totalDamageDealt": me.get("totalDamageDealt"),
        "physicalDamageDealtToChampions": me.get("physicalDamageDealtToChampions"),
        "magicDamageDealtToChampions": me.get("magicDamageDealtToChampions"),
        "trueDamageDealtToChampions": me.get("trueDamageDealtToChampions"),
        "physicalDamageTaken": me.get("physicalDamageTaken"),
        "magicDamageTaken": me.get("magicDamageTaken"),
        "trueDamageTaken": me.get("trueDamageTaken"),
        "totalDamageTaken": me.get("totalDamageTaken"),
        "damageSelfMitigated": me.get("damageSelfMitigated"),
        "totalHeal": me.get("totalHeal"),
        "totalHealsOnTeammates": me.get("totalHealsOnTeammates"),
        "totalDamageShieldedOnTeammates": me.get("totalDamageShieldedOnTeammates"),
        "timeCCingOthers": me.get("timeCCingOthers"),
        "totalTimeCCDealt": me.get("totalTimeCCDealt"),
        "turretTakedowns": me.get("turretTakedowns"),
        "inhibitorTakedowns": me.get("inhibitorTakedowns"),
        "inhibitorKills": me.get("inhibitorKills"),
        "nexusKills": me.get("nexusKills"),
        "champExperience": me.get("champExperience"),
    }

    # Ajouter les stat perks si présents
    for k, v in stat_perks.items():
        row[f"statPerk_{k}"] = v

    # Challenges (prefixed ch_)
    challenges = me.get("challenges") or {}
    if isinstance(challenges, dict):
        for ck, cv in challenges.items():
            if isinstance(cv, (int, float, str)):
                row[f"ch_{ck}"] = cv

    # Spell casts
    for sidx in range(1, 5):
        key = f"spell{sidx}Casts"
        if key in me:
            row[key] = me.get(key)

    # Objectifs équipes & bans
    teams_info = info.get("teams", [])
    my_team_data = None
    opp_team_data = None
    for t in teams_info:
        if t.get("teamId") == my_team_id:
            my_team_data = t
        elif t.get("teamId") in (100, 200):
            opp_team_data = t
    def add_obj(prefix: str, tdata: Optional[Dict[str, Any]]):
        if not tdata:
            return
        objectives = tdata.get("objectives") or {}
        for oname, oval in objectives.items():
            if isinstance(oval, dict):
                row[f"{prefix}obj_{oname}_kills"] = oval.get("kills")
                row[f"{prefix}obj_{oname}_first"] = oval.get("first")
        bans = tdata.get("bans") or []
        ban_ids = []
        for b in bans:
            cid = b.get("championId")
            if cid is not None:
                ban_ids.append(cid)
        if ban_ids:
            if pretty and champ_map:
                row[f"{prefix}bans"] = " | ".join(champ_map.get(cid, str(cid)) for cid in ban_ids)
            else:
                row[f"{prefix}bans"] = " | ".join(str(cid) for cid in ban_ids)
    add_obj("team_", my_team_data)
    add_obj("opp_", opp_team_data)

    # Objectives per minute
    if minutes:
        for side in ("team_","opp_"):
            for obj in ("dragon","baron","riftHerald","tower","inhibitor"):
                key = f"{side}obj_{obj}_kills"
                if key in row and row[key] is not None:
                    row[f"{key}_per_min"] = round(row[key] / minutes, 4)
    # Team kills stats
    row["team_kills"] = team_kills
    row["team_kills_per_min"] = round(safe_div(team_kills, minutes), 3) if minutes else 0.0

    # Inject goldDiff snapshots if timeline provided
    if timeline_extra:
        for minute in (10,15,20,25):
            key = f"goldDiff{minute}"
            if key in timeline_extra:
                row[key] = timeline_extra.get(key)

    return row

def get_timeline(session: requests.Session, platform: str, match_id: str, api_key: str) -> Optional[Dict[str, Any]]:
    try:
        base = regional_match_base(platform)
        headers = {"X-Riot-Token": api_key}
        url = f"{base}/lol/match/v5/matches/{match_id}/timeline"
        r = request_with_backoff(session, "GET", url, headers=headers)
        if r.status_code != 200:
            return None
        return r.json()
    except Exception:
        return None

def parse_timeline(timeline: Dict[str, Any], puuid: str) -> Tuple[Dict[str, Any], List[Dict[str, Any]]]:
    metrics: Dict[str, Any] = {}
    progression: List[Dict[str, Any]] = []
    try:
        info = timeline.get("info", {})
        frames = info.get("frames", [])
        meta_participants = timeline.get("metadata", {}).get("participants", [])
        if puuid not in meta_participants:
            return metrics, progression
        participant_id = meta_participants.index(puuid) + 1
        targets_ms = [600000, 900000, 1200000, 1500000]
        captured = {ms: False for ms in targets_ms}
        dragons_100: List[str] = []
        dragons_200: List[str] = []
        for f in frames:
            ts = f.get("timestamp", 0)
            pframes = f.get("participantFrames", {})
            team100 = 0
            team200 = 0
            for pid_str, pdata in pframes.items():
                try:
                    pid = int(pid_str)
                except ValueError:
                    continue
                tg = pdata.get("totalGold", 0)
                if 1 <= pid <= 5:
                    team100 += tg
                else:
                    team200 += tg
            progression.append({
                "timestamp_s": ts / 1000.0,
                "team100_gold": team100,
                "team200_gold": team200,
                "gold_diff": team100 - team200,
            })
            for target in targets_ms:
                if (not captured[target]) and ts >= target:
                    pf = pframes.get(str(participant_id))
                    if pf:
                        minions = pf.get("minionsKilled", 0) + pf.get("jungleMinionsKilled", 0)
                        if target == 600000:
                            metrics["cs10"] = minions
                            metrics["gold10"] = pf.get("totalGold")
                            metrics["xp10"] = pf.get("xp")
                        label = str(target // 60000)
                        metrics[f"goldDiff{label}"] = team100 - team200
                        metrics[f"goldTotal{label}"] = team100 + team200
                    captured[target] = True
            for ev in f.get("events", []) or []:
                if ev.get("type") == "ELITE_MONSTER_KILL" and ev.get("monsterType") == "DRAGON":
                    subtype = ev.get("monsterSubType") or "DRAGON"
                    killer_team = ev.get("killerTeamId")
                    if killer_team == 100:
                        dragons_100.append(subtype)
                    elif killer_team == 200:
                        dragons_200.append(subtype)
        my_team_id = 100 if participant_id <= 5 else 200
        if my_team_id == 100:
            metrics["team_dragons_order"] = "->".join(dragons_100)
            metrics["opp_dragons_order"] = "->".join(dragons_200)
        else:
            metrics["team_dragons_order"] = "->".join(dragons_200)
            metrics["opp_dragons_order"] = "->".join(dragons_100)
        if progression:
            last = progression[-1]
            diff_final = last.get("gold_diff")
            total_final = (last.get("team100_gold", 0) or 0) + (last.get("team200_gold", 0) or 0)
            metrics["goldDiffFinal"] = diff_final
            metrics["goldDiffFinalNorm"] = (diff_final / total_final) if (total_final and diff_final is not None) else None
        return metrics, progression
    except Exception:
        return metrics, progression

# ----------------------
# Programme principal
# ----------------------

def parse_args():
    p = argparse.ArgumentParser(description="Export LoL match history & stats vers CSV/Parquet (API 2025).")
    p.add_argument("--platform", required=True, help="Plateforme ex: euw1, na1, kr ...")
    # Identité joueur (au moins l'un des trois)
    p.add_argument("--riot-id", help='Riot ID "GameName#TagLine" (recommandé).')
    p.add_argument("--game-name", help="Riot GameName (si --tag-line est fourni).")
    p.add_argument("--tag-line", help="Riot TagLine (si --game-name est fourni).")
    p.add_argument("--puuid", help="PUUID direct (skip account-v1).")

    p.add_argument("--api-key", default=None, help="Clé API Riot (sinon utilise RIOT_API_KEY).")
    p.add_argument("--count", type=int, default=100, help="Nombre max de matchs à récupérer (1-1000).")
    p.add_argument("--queues", type=int, nargs='*', default=None, help="Filtrer par queueId (ex: 420, 440, 450, 490).")
    p.add_argument("--start-time", type=int, default=None, help="Epoch seconds min (filtre).")
    p.add_argument("--end-time", type=int, default=None, help="Epoch seconds max (filtre).")
    p.add_argument("--pretty", action="store_true", help="Afficher noms d'objets/champions/runes via Data Dragon.")
    p.add_argument("--lang", default="en_US", help="Langue Data Dragon (ex: fr_FR, en_US).")
    p.add_argument("--out", required=True, help="Chemin de sortie (.csv ou .parquet).")
    p.add_argument("--extended", action="store_true", help="Ajoute colonnes supplémentaires (spells, partages, patch, etc.).")
    p.add_argument("--timeline", action="store_true", help="Récupère la timeline pour métriques early (cs10, gold10, xp10).")
    p.add_argument("--league", action="store_true", help="Ajoute rang Solo/Flex (league-v4) répété sur chaque ligne.")
    p.add_argument("--all", action="store_true", help="Récupère le maximum de données disponibles (implique --extended --timeline --league).")
    return p.parse_args()

def main():
    args = parse_args()
    api_key = get_api_key(args.api_key)

    # --all implique tous les enrichissements
    if getattr(args, "all", False):
        args.extended = True
        args.timeline = True
        args.league = True
        args.pretty = True

    if not (1 <= args.count <= 1000):
        print("⚠️ --count doit être entre 1 et 1000. Ajustement à 100.", file=sys.stderr)
        args.count = 100

    session = requests.Session()

    # 1) Résoudre le PUUID
    if args.puuid:
        puuid = args.puuid
        print(f"🔑 PUUID fourni: {puuid[:12]}...")
    else:
        print("🔎 Résolution du PUUID via ACCOUNT-V1 (Riot ID requis)...")
        puuid = get_puuid_from_riotid(session, args.platform, args.riot_id, args.game_name, args.tag_line, api_key)
        print(f"✅ PUUID trouvé: {puuid[:12]}...")

    # Métadonnées de compte (optionnel)
    summ_profile = optional_summoner_profile(session, args.platform, puuid, api_key)
    if summ_profile:
        print(f"👤 Invocateur: {summ_profile.get('name')} (level {summ_profile.get('summonerLevel')})")
    else:
        print("ℹ️ Impossible d'obtenir le profil Summoner (non bloquant).")

    # Data Dragon
    item_map = {}
    champ_map = {}
    rune_map = {}
    if args.pretty or args.extended:
        print("⬇️ Chargement des données Data Dragon (items/champions/runes/spells)...")
        item_map, champ_map, rune_map, rune_style_map, spell_map = load_ddragon_maps(args.lang)
    else:
        item_map, champ_map, rune_map, rune_style_map, spell_map = {}, {}, {}, {}, {}

    # League ranks (optionnel)
    league_info = {}
    if args.league and summ_profile and summ_profile.get("id"):
        try:
            base = platform_base_url(args.platform)
            url = f"{base}/lol/league/v4/entries/by-summoner/{summ_profile['id']}"
            headers = {"X-Riot-Token": api_key}
            r = request_with_backoff(session, "GET", url, headers=headers)
            if r.status_code == 200:
                entries = r.json()
                for e in entries:
                    qtype = e.get("queueType")
                    if qtype == "RANKED_SOLO_5x5":
                        league_info.update({
                            "rank_solo_tier": e.get("tier"),
                            "rank_solo_div": e.get("rank"),
                            "rank_solo_lp": e.get("leaguePoints"),
                            "rank_solo_wins": e.get("wins"),
                            "rank_solo_losses": e.get("losses"),
                        })
                    elif qtype == "RANKED_FLEX_SR":
                        league_info.update({
                            "rank_flex_tier": e.get("tier"),
                            "rank_flex_div": e.get("rank"),
                            "rank_flex_lp": e.get("leaguePoints"),
                            "rank_flex_wins": e.get("wins"),
                            "rank_flex_losses": e.get("losses"),
                        })
        except Exception as e:
            print(f"ℹ️ League non récupéré: {e}")

    # 2) Match IDs (mode multi-catégories automatique si --all et aucune queue explicite)
    queue_groups = {
        "solo": [420],
        "flex": [440],
        "aram": [450],
        "normal": [400, 430, 490],
    }
    category_match_ids: Dict[str, List[str]] = {}
    if getattr(args, "all", False) and not args.queues:
        total_est = args.count * len(queue_groups)
        print(f"⬇️ Récupération multi-catégories (~{total_est} matchs max) ...")
        for cat, qlist in queue_groups.items():
            try:
                print(f"  ▶ {cat} (queues {qlist}) ...")
                mids = get_match_ids(
                    session,
                    args.platform,
                    puuid,
                    api_key,
                    count=args.count,
                    queues=qlist,
                    start_time=args.start_time,
                    end_time=args.end_time,
                )
                category_match_ids[cat] = mids
                print(f"    → {len(mids)} matchId(s)")
            except Exception as e:
                print(f"⚠️ Impossible de récupérer {cat}: {e}")
        match_ids = [m for lst in category_match_ids.values() for m in lst]
        print(f"📦 Total collecté: {len(match_ids)} matchId(s).")
    else:
        print(f"⬇️ Récupération de {args.count} matchId(s)...")
        match_ids = get_match_ids(
            session,
            args.platform,
            puuid,
            api_key,
            count=args.count,
            queues=args.queues,
            start_time=args.start_time,
            end_time=args.end_time,
        )
        print(f"📦 {len(match_ids)} matchId(s) collectés.")

    rows = []
    raw_timeline_dir = None
    if getattr(args, "all", False):
        base_no_ext, _ext = os.path.splitext(args.out)
        raw_timeline_dir = base_no_ext + "_timelines"
        if args.timeline and not os.path.isdir(raw_timeline_dir):
            try:
                os.makedirs(raw_timeline_dir, exist_ok=True)
            except Exception:
                raw_timeline_dir = None
    team_gold_progression_rows: List[Dict[str, Any]] = []
    for i, mid in enumerate(match_ids, 1):
        try:
            m = get_match(session, args.platform, mid, api_key)
            timeline_metrics = None
            tl = None
            if args.timeline:
                tl = get_timeline(session, args.platform, mid, api_key)
                if tl:
                    timeline_metrics, gold_prog = parse_timeline(tl, puuid)
                    for entry in gold_prog:
                        rec = dict(entry)
                        rec["matchId"] = mid
                        team_gold_progression_rows.append(rec)
                    if raw_timeline_dir:
                        try:
                            with open(os.path.join(raw_timeline_dir, f"{mid}.json"), "w", encoding="utf-8") as f:
                                json.dump(tl, f)
                        except Exception:
                            pass
            row = extract_row_for_puuid(
                m, puuid, item_map, champ_map, rune_map, rune_style_map, spell_map,
                pretty=(args.pretty or args.extended), timeline_extra=timeline_metrics)
            if league_info:
                row.update(league_info)
            rows.append(row)
            if i % 10 == 0 or i == len(match_ids):
                print(f"  - {i}/{len(match_ids)}")
        except Exception as e:
            print(f"⚠️ Erreur sur {mid}: {e}", file=sys.stderr)
            continue
        # Plus besoin de sleep fixe; le RATE_LIMITER régule.

    if not rows:
        print("Aucune ligne à écrire. Arrêt.", file=sys.stderr)
        sys.exit(2)

    df = pd.DataFrame(rows)
    df.sort_values("date", inplace=True)

    out_path = args.out
    ext = os.path.splitext(out_path)[1].lower()
    if ext == ".csv":
        df.to_csv(out_path, index=False)
    elif ext == ".parquet":
        try:
            df.to_parquet(out_path, index=False)
        except Exception as e:
            print(f"⚠️ Écriture Parquet impossible ({e}). On repasse en CSV.", file=sys.stderr)
            out_path = os.path.splitext(out_path)[0] + ".csv"
            df.to_csv(out_path, index=False)
    else:
        print("⚠️ Extension non reconnue. Utilisez .csv ou .parquet. Écriture CSV par défaut.", file=sys.stderr)
        out_path = out_path + ".csv"
        df.to_csv(out_path, index=False)

    # Fichiers supplémentaires: mastery & résumé par catégorie
    base_no_ext, _ext = os.path.splitext(out_path)
    # Champion mastery
    if getattr(args, "all", False) and summ_profile.get("id"):
        try:
            base_plat = platform_base_url(args.platform)
            url = f"{base_plat}/lol/champion-mastery/v4/champion-masteries/by-summoner/{summ_profile['id']}"
            headers = {"X-Riot-Token": api_key}
            r = request_with_backoff(session, "GET", url, headers=headers)
            if r.status_code == 200:
                mastery = r.json()
                # Map champion IDs to names if available
                for entry in mastery:
                    cid = entry.get("championId")
                    if (args.pretty or args.extended) and champ_map:
                        entry["championName"] = champ_map.get(cid, entry.get("championId"))
                pd.DataFrame(mastery).to_csv(base_no_ext + "_champion_mastery.csv", index=False)
                print("📄 Maîtrise champions écrite.")
        except Exception as e:
            print(f"ℹ️ Champion mastery ignorée: {e}")

    # Résumé par queue_category
    if "queue_category" in df.columns:
        try:
            agg = df.groupby("queue_category").agg({
                "matchId": "count",
                "win": lambda s: float(sum(s))/len(s),
                "kda": "mean",
                "kp": "mean",
                "cs_per_min": "mean",
                "gpm": "mean",
                "dpm": "mean",
                "vision_score": "mean",
            }).rename(columns={"matchId": "games", "win": "winrate"})
            agg.reset_index().to_csv(base_no_ext + "_queue_summary.csv", index=False)
            print("📄 Résumé par catégorie écrit.")
        except Exception as e:
            print(f"ℹ️ Résumé par catégorie non généré: {e}")

    # Export progression gold équipe
    if team_gold_progression_rows:
        try:
            pd.DataFrame(team_gold_progression_rows).to_csv(base_no_ext + "_team_gold_progression.csv", index=False)
            print("📄 Progression gold d'équipe écrite.")
        except Exception as e:
            print(f"ℹ️ Progression gold non écrite: {e}")

    print(f"✅ Fichier écrit : {out_path}")
    print("Mode ALL: colonnes étendues + mastery + résumés + timelines (si activé).")
    print("Bon jeu ! 🎮")

if __name__ == "__main__":
    main()
