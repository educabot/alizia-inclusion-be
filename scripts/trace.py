#!/usr/bin/env python3
"""
trace.py — reconstruye conversaciones del agente de Alizia desde los logs de Railway.

El backend (alizia-inclusion-be) emite logs JSON estructurado (slog) con CHAT_TRACE_VERBOSE=true,
así que cada turno deja en producción el prompt, las tool calls (con argumentos), los resultados de
cada tool y la respuesta. Esta herramienta baja esos logs y los reconstruye turno por turno para
auditar e iterar el agente.

Vía de acceso: `railway logs -d --json` devuelve cada evento como objeto JSON con TODOS sus atributos
y los correlation ids (request_id, org_id, user_id). OJO: el flag `--filter` de Railway sólo busca en
el nombre del evento (`message`), no en los atributos — por eso el filtrado por alumno/usuario se hace
LOCALMENTE acá.

Requiere la CLI de Railway instalada y logueada (`railway whoami`). Sólo lectura.

Ejemplos:
    python3 scripts/trace.py last                      # la última conversación
    python3 scripts/trace.py last --user 5 --n 2       # las 2 últimas del docente user_id=5
    python3 scripts/trace.py student Francisco         # conversaciones que mencionan a Francisco
    python3 scripts/trace.py conversation 238          # todos los turnos de la conversación 238
    python3 scripts/trace.py request <request_id>      # un turno puntual
    python3 scripts/trace.py last --format md > out.md # export markdown
    cat dump.ndjson | python3 scripts/trace.py last --stdin   # re-analizar un dump guardado

Flags comunes: --since/--until/--lines (ventana de logs), --user, --format {pretty,md,json},
--full (no recortar textos), --rag (incluir eventos rag.*/embed.*), --stdin (leer NDJSON de stdin).
"""

import argparse
import json
import os
import subprocess
import sys
from collections import defaultdict, OrderedDict

# --- Defaults de Railway (override por flag o env) -------------------------------------------------
DEF_PROJECT = os.environ.get("RAILWAY_PROJECT_ID", "54b0088b-c001-4991-8ddd-8aaded0c270f")
DEF_SERVICE = os.environ.get("RAILWAY_SERVICE", "Alizia-Inclusion")
DEF_ENV = os.environ.get("RAILWAY_ENVIRONMENT", "production")

# Eventos del pipeline de chat que sabemos renderizar. Cualquier otro se ignora con gracia.
KNOWN = {
    "chat.context_loaded", "chat.context_build_failed", "chat.prompt_built",
    "chat.agentic_iteration", "chat.tool_call", "chat.tool_result", "chat.tool_error",
    "ai.request", "ai.response", "chat.sources_used", "chat.guardrail_tripped",
    "chat.turn_done", "chat.raw_response",
    "rag.request", "rag.embed_ok", "rag.search", "embed.request",
}

# Códigos ANSI (se desactivan si la salida no es una TTY o --format != pretty).
class C:
    DIM = "\033[2m"; BOLD = "\033[1m"; RESET = "\033[0m"
    CYAN = "\033[36m"; YELLOW = "\033[33m"; GREEN = "\033[32m"; RED = "\033[31m"; MAGENTA = "\033[35m"


def color_enabled(fmt):
    return fmt == "pretty" and sys.stdout.isatty()


# --- Carga de logs ---------------------------------------------------------------------------------

def fetch_logs(args):
    """Devuelve una lista de eventos (dicts). Lee de stdin o corre la CLI de Railway."""
    if args.stdin:
        raw = sys.stdin.read().splitlines()
    else:
        cmd = ["railway", "logs", "-d", "--json",
               "-p", args.project, "-s", args.service, "-e", args.environment]
        # Railway hace streaming si no se acota; siempre pasamos --lines (acota y desactiva el stream).
        if args.since:
            cmd += ["--since", args.since]
        if args.until:
            cmd += ["--until", args.until]
        cmd += ["--lines", str(args.lines or 1000)]
        try:
            out = subprocess.run(cmd, capture_output=True, text=True, timeout=180)
        except FileNotFoundError:
            sys.exit("ERROR: no se encontró la CLI `railway`. Instalá con `brew install railway` y `railway login`.")
        except subprocess.TimeoutExpired:
            sys.exit("ERROR: `railway logs` tardó demasiado. Acotá con --lines o --since/--until.")
        if out.returncode != 0:
            sys.exit(f"ERROR al correr railway logs:\n{out.stderr.strip()}\n"
                     "¿Estás logueada? Probá `railway whoami`.")
        raw = out.stdout.splitlines()

    events = []
    for line in raw:
        line = line.strip()
        if not line.startswith("{"):
            continue
        try:
            events.append(json.loads(line))
        except json.JSONDecodeError:
            continue
    return events


# --- Estructuración --------------------------------------------------------------------------------

def build_turns(events):
    """Agrupa eventos por request_id (= turno), ordenados por tiempo. Devuelve OrderedDict rid->turn."""
    by_rid = defaultdict(list)
    for e in events:
        rid = e.get("request_id") or "(sin-request-id)"
        by_rid[rid].append(e)

    turns = {}
    for rid, evs in by_rid.items():
        evs.sort(key=lambda d: d.get("time", d.get("timestamp", "")))
        td = next((e for e in evs if e.get("message") == "chat.turn_done"), None)
        pb = next((e for e in evs if e.get("message") == "chat.prompt_built"), None)
        turns[rid] = {
            "rid": rid,
            "events": evs,
            "time": evs[0].get("time", evs[0].get("timestamp", "")),
            "end_time": evs[-1].get("time", evs[-1].get("timestamp", "")),
            "user_id": next((e.get("user_id") for e in evs if e.get("user_id")), None),
            "conversation_id": td.get("conversation_id") if td else None,
            "has_chat": any(e.get("message") in KNOWN for e in evs),
            "turn_done": td,
            "prompt_built": pb,
        }
    # ordenar turnos por tiempo
    return OrderedDict(sorted(turns.items(), key=lambda kv: kv[1]["time"]))


def turn_text_blob(turn):
    """Concatena los textos de un turno para búsqueda local por alumno. Excluye system_prompt y
    payload a propósito: el system prompt lista a TODOS los alumnos de la org, así que incluirlo
    haría que cualquier conversación matchee cualquier nombre. Buscamos solo en lo que el docente
    escribió y en lo que las tools pidieron/devolvieron."""
    parts = []
    for e in turn["events"]:
        for k in ("user_message", "args", "result", "response"):
            v = e.get(k)
            if isinstance(v, str):
                parts.append(v)
    return "\n".join(parts).lower()


def group_conversations(turns):
    """conv_id -> lista de turnos (ordenados). Los turnos sin conversation_id quedan sueltos por rid."""
    convs = defaultdict(list)
    for rid, t in turns.items():
        key = t["conversation_id"] if t["conversation_id"] is not None else f"rid:{rid}"
        convs[key].append(t)
    for k in convs:
        convs[k].sort(key=lambda t: t["time"])
    return convs


# --- Render ----------------------------------------------------------------------------------------

def trunc(s, n, full):
    if s is None:
        return ""
    s = str(s)
    if full or len(s) <= n:
        return s
    return s[:n] + f" …(+{len(s) - n} chars)"


def unwrap_json(s):
    """Los campos args/result llegan como string JSON con un escape de más que Railway agrega
    (ej. '{\\"student_id\\": 2}'). Lo limpiamos y devolvemos el JSON compacto y legible; si no
    parsea, devolvemos el string tal cual."""
    if not isinstance(s, str):
        return s
    candidate = s.replace('\\"', '"') if '\\"' in s else s
    try:
        return json.dumps(json.loads(candidate), ensure_ascii=False)
    except Exception:
        return s


def fmt_payload(s, n, full):
    return trunc(unwrap_json(s), n, full)


def clean_text(s):
    """Des-escapa el nivel extra que Railway agrega a los textos libres (\\n, \\", \\t)."""
    if not isinstance(s, str):
        return s
    return s.replace('\\n', '\n').replace('\\t', '\t').replace('\\"', '"').replace('\\\\', '\\')


def fmt_text(s, n, full, oneline=False):
    s = clean_text(s)
    if oneline and isinstance(s, str):
        s = s.replace('\n', ' · ')
    return trunc(s, n, full)


def detect_verbose_off(events):
    """Verbose off ⇒ aparece <x>_len SIN su contraparte de texto <x>. (result_len siempre coexiste
    con result en chat.tool_result, así que no lo usamos como señal para evitar falsos positivos.)"""
    for e in events:
        for base in ("system_prompt", "user_message", "args", "response"):
            if (base + "_len") in e and base not in e:
                return True
    return False


def render_turn_pretty(turn, args, col):
    out = []
    c = C if col else type("N", (), {k: "" for k in vars(C) if not k.startswith("_")})
    pb = turn["prompt_built"] or {}
    td = turn["turn_done"] or {}
    head = f"{c.BOLD}── turno {turn['rid'][:8]} ─ {turn['time'][:19]} ─ user {turn['user_id']}{c.RESET}"
    out.append(head)
    if pb.get("user_message") is not None:
        out.append(f"  {c.CYAN}USER{c.RESET}  {fmt_text(pb.get('user_message'), 300, args.full, oneline=True)}")
    if args.full and pb.get("system_prompt"):
        out.append(f"  {c.DIM}SYSTEM PROMPT:{c.RESET}\n{clean_text(pb.get('system_prompt'))}")

    for e in turn["events"]:
        m = e.get("message")
        if m == "chat.agentic_iteration":
            out.append(f"  {c.DIM}· iteración agéntica: {e.get('tool_calls')} tool call(s){c.RESET}")
        elif m == "chat.tool_call":
            out.append(f"  {c.YELLOW}→ TOOL{c.RESET} {c.BOLD}{e.get('tool')}{c.RESET}  args={fmt_payload(e.get('args'), 200, args.full)}")
        elif m == "chat.tool_result":
            out.append(f"    {c.GREEN}← result{c.RESET} ({e.get('result_len')} chars) {fmt_payload(e.get('result'), 160, args.full)}")
        elif m == "chat.tool_error":
            out.append(f"    {c.RED}← ERROR{c.RESET} {e.get('tool')}: {e.get('error')}")
        elif m == "ai.response":
            out.append(f"  {c.DIM}· modelo: {e.get('total_tokens')} tokens, {e.get('tool_calls')} tool calls pedidas, {e.get('duration_ms')}ms{c.RESET}")
        elif m == "chat.sources_used":
            tools = e.get("tools") or []
            flags = []
            if e.get("used_valija"): flags.append("valija")
            if e.get("used_student"): flags.append("alumno")
            if e.get("used_rag"): flags.append(f"RAG({e.get('rag_hits')} hits)")
            out.append(f"  {c.MAGENTA}FUENTES{c.RESET} {', '.join(flags) or 'ninguna'}  tools={tools or '—'}")
        elif m == "chat.guardrail_tripped":
            out.append(f"  {c.RED}⚠ GUARDRAIL{c.RESET} {e.get('reason')}")
        elif m in ("rag.request", "rag.embed_ok", "rag.search", "embed.request") and args.rag:
            extra = {k: v for k, v in e.items() if k not in ("message", "time", "timestamp", "level", "request_id", "org_id", "user_id")}
            out.append(f"  {c.DIM}· {m} {trunc(json.dumps(extra, ensure_ascii=False), 160, args.full)}{c.RESET}")

    resp = td.get("response")
    if resp is not None:
        out.append(f"  {c.CYAN}ALIZIA{c.RESET} {fmt_text(resp, 400, args.full, oneline=True)}")
    meta = []
    if td.get("identified_student"): meta.append(f"alumno={td['identified_student']}")
    if td.get("recommended_device"): meta.append(f"device={td['recommended_device']}")
    if td.get("has_adaptation"): meta.append("guardó adaptación")
    if td.get("questions_count"): meta.append(f"{td['questions_count']} preguntas")
    if td.get("referenced_count"): meta.append(f"{td['referenced_count']} refs")
    if td.get("total_tokens"): meta.append(f"{td['total_tokens']} tokens")
    if meta:
        out.append(f"  {c.DIM}{' · '.join(meta)}{c.RESET}")
    return "\n".join(out)


def render_turn_md(turn, args):
    out = []
    pb = turn["prompt_built"] or {}
    td = turn["turn_done"] or {}
    out.append(f"### Turno `{turn['rid'][:8]}` — {turn['time'][:19]} — user {turn['user_id']}")
    if pb.get("user_message") is not None:
        out.append(f"**Usuario:** {fmt_text(pb.get('user_message'), 500, args.full)}")
    for e in turn["events"]:
        m = e.get("message")
        if m == "chat.tool_call":
            out.append(f"- 🔧 `{e.get('tool')}` — args: `{fmt_payload(e.get('args'), 200, args.full)}`")
        elif m == "chat.tool_result":
            out.append(f"  - ← {e.get('result_len')} chars: `{fmt_payload(e.get('result'), 160, args.full)}`")
        elif m == "chat.tool_error":
            out.append(f"  - ❌ {e.get('tool')}: {e.get('error')}")
        elif m == "chat.sources_used":
            out.append(f"- **Fuentes:** valija={e.get('used_valija')} alumno={e.get('used_student')} rag={e.get('used_rag')}({e.get('rag_hits')}) tools={e.get('tools')}")
        elif m == "chat.guardrail_tripped":
            out.append(f"- ⚠️ **Guardrail:** {e.get('reason')}")
    if td.get("response") is not None:
        out.append(f"**Alizia:** {fmt_text(td.get('response'), 600, args.full)}")
        out.append(f"<sub>conv={td.get('conversation_id')} alumno={td.get('identified_student')} device={td.get('recommended_device')} adaptación={td.get('has_adaptation')} tokens={td.get('total_tokens')}</sub>")
    return "\n".join(out)


def conv_summary(turnlist, col):
    c = C if col else type("N", (), {k: "" for k in vars(C) if not k.startswith("_")})
    total_tokens = sum((t["turn_done"] or {}).get("total_tokens", 0) or 0 for t in turnlist)
    used_rag = any(any(e.get("message") == "chat.sources_used" and e.get("used_rag") for e in t["events"]) for t in turnlist)
    all_tools = []
    for t in turnlist:
        for e in t["events"]:
            if e.get("message") == "chat.tool_call":
                all_tools.append(e.get("tool"))
    line = (f"{c.BOLD}resumen:{c.RESET} {len(turnlist)} turno(s) · {total_tokens} tokens · "
            f"RAG {'sí' if used_rag else 'NO'} · tools usadas: {', '.join(sorted(set(all_tools))) or 'ninguna'}")
    return line


# --- Selección por subcomando ----------------------------------------------------------------------

def filter_user(turns, user):
    if not user:
        return turns
    return OrderedDict((rid, t) for rid, t in turns.items() if str(t["user_id"]) == str(user))


def cmd_last(turns, args):
    turns = filter_user(turns, args.user)
    convs = group_conversations(turns)
    # ordenar conversaciones por su último turno
    ordered = sorted(convs.items(), key=lambda kv: kv[1][-1]["time"], reverse=True)
    return [tl for _, tl in ordered[: args.n]]


def cmd_student(turns, args):
    name = args.name.lower()
    turns = filter_user(turns, args.user)
    convs = group_conversations(turns)
    # Una conversación entra si CUALQUIER turno suyo menciona al alumno; se muestra completa.
    matched = [(k, tl) for k, tl in convs.items() if any(name in turn_text_blob(t) for t in tl)]
    matched.sort(key=lambda kv: kv[1][-1]["time"], reverse=True)
    return [tl for _, tl in matched[: args.n]]


def cmd_conversation(turns, args):
    sel = [t for t in turns.values() if str(t["conversation_id"]) == str(args.id)]
    sel.sort(key=lambda t: t["time"])
    return [sel] if sel else []


def cmd_request(turns, args):
    t = turns.get(args.id) or next((t for rid, t in turns.items() if rid.startswith(args.id)), None)
    return [[t]] if t else []


# --- Main ------------------------------------------------------------------------------------------

def main():
    # Flags comunes en un parent parser → válidos antes o después del subcomando.
    common = argparse.ArgumentParser(add_help=False)
    common.add_argument("--project", default=DEF_PROJECT)
    common.add_argument("--service", default=DEF_SERVICE)
    common.add_argument("--environment", "-e", default=DEF_ENV)
    common.add_argument("--since")
    common.add_argument("--until")
    common.add_argument("--lines", type=int)
    common.add_argument("--user", help="filtrar por user_id")
    common.add_argument("--format", choices=["pretty", "md", "json"], default="pretty")
    common.add_argument("--full", action="store_true", help="no recortar textos")
    common.add_argument("--rag", action="store_true", help="incluir eventos rag.*/embed.*")
    common.add_argument("--stdin", action="store_true", help="leer NDJSON de stdin en vez de Railway")
    common.add_argument("-n", "--n", type=int, default=1, help="cantidad de conversaciones (last/student)")

    p = argparse.ArgumentParser(
        parents=[common],
        description="Reconstruye conversaciones del agente de Alizia desde los logs de Railway.")
    sub = p.add_subparsers(dest="cmd", required=True)
    sub.add_parser("last", parents=[common], help="última(s) conversación(es)")
    sp = sub.add_parser("student", parents=[common], help="conversaciones que mencionan a un alumno"); sp.add_argument("name")
    sp = sub.add_parser("conversation", parents=[common], help="todos los turnos de una conversación"); sp.add_argument("id")
    sp = sub.add_parser("request", parents=[common], help="un turno puntual"); sp.add_argument("id")

    args = p.parse_args()

    events = fetch_logs(args)
    if not events:
        sys.exit("No se obtuvieron eventos. Ampliá la ventana con --lines/--since o revisá el acceso a Railway.")
    if detect_verbose_off(events):
        print("AVISO: hay claves <x>_len → CHAT_TRACE_VERBOSE parece estar apagado; faltan los textos completos.\n",
              file=sys.stderr)

    turns = build_turns(events)
    turns = OrderedDict((rid, t) for rid, t in turns.items() if t["has_chat"])

    dispatch = {"last": cmd_last, "student": cmd_student, "conversation": cmd_conversation, "request": cmd_request}
    convlist = dispatch[args.cmd](turns, args)

    if not convlist:
        sys.exit("Sin resultados para ese criterio en la ventana de logs traída. Probá ampliar --lines/--since.")

    col = color_enabled(args.format)

    if args.format == "json":
        norm = []
        for tl in convlist:
            norm.append({
                "conversation_id": tl[0]["conversation_id"],
                "turns": [{
                    "request_id": t["rid"], "time": t["time"], "user_id": t["user_id"],
                    "user_message": (t["prompt_built"] or {}).get("user_message"),
                    "tool_calls": [{"tool": e.get("tool"), "args": unwrap_json(e.get("args"))}
                                   for e in t["events"] if e.get("message") == "chat.tool_call"],
                    "sources_used": next((e for e in t["events"] if e.get("message") == "chat.sources_used"), None),
                    "response": (t["turn_done"] or {}).get("response"),
                    "turn_done": t["turn_done"],
                } for t in tl],
            })
        print(json.dumps(norm, ensure_ascii=False, indent=2))
        return

    for tl in convlist:
        conv_id = tl[0]["conversation_id"]
        if args.format == "md":
            print(f"## Conversación {conv_id if conv_id is not None else '(suelta)'}\n")
            for t in tl:
                print(render_turn_md(t, args)); print()
        else:
            header = f"{C.BOLD}══ Conversación {conv_id if conv_id is not None else '(suelta)'} ══{C.RESET}" if col \
                else f"══ Conversación {conv_id if conv_id is not None else '(suelta)'} ══"
            print(header)
            for t in tl:
                print(render_turn_pretty(t, args, col)); print()
            print(conv_summary(tl, col)); print()


if __name__ == "__main__":
    main()
