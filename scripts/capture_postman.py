#!/usr/bin/env python3
"""Captura request/response REALES de todos los endpoints y genera artefactos.

Corre contra el server en ENV=test (auth mockeada: user 1 / org de ceros) con los
seeds de prueba aplicados (alumnos 9001-9003, aula 9001, ramp 9001, devices
9001-9002, contenido pedagógico 1-3). NO requiere token.

Genera:
  - docs/api/postman_capture.json          (captura cruda: req + resp + expected)
  - alizia-inclusion.postman_collection.json (colección Postman v2.1, con pm.test)
  - docs/api/ENDPOINTS.md                   (referencia: endpoint/body/expected/response)

Uso:  python scripts/capture_postman.py [base_url]
"""
import json
import os
import sys
import urllib.request
import urllib.error

BASE = sys.argv[1] if len(sys.argv) > 1 else "http://127.0.0.1:8080"
API = BASE + "/api/v1"
HERE = os.path.dirname(os.path.abspath(__file__))
ROOT = os.path.dirname(HERE)
DOCS = os.path.join(ROOT, "docs", "api")
os.makedirs(DOCS, exist_ok=True)

records = []  # cada item documenta un request capturado
ctx = {}      # ids creados durante el flujo


def call(method, path, body=None, raw_path=False):
    url = (BASE + path) if raw_path else (API + path)
    data = json.dumps(body).encode() if body is not None else None
    req = urllib.request.Request(url, data=data, method=method,
                                 headers={"Content-Type": "application/json"})
    try:
        with urllib.request.urlopen(req, timeout=30) as r:
            txt = r.read().decode()
            ct = r.headers.get("Content-Type", "")
            return r.status, txt, ct
    except urllib.error.HTTPError as e:
        return e.code, e.read().decode(), e.headers.get("Content-Type", "")


def unwrap(txt):
    try:
        p = json.loads(txt) if txt else {}
    except json.JSONDecodeError:
        return None, txt
    if isinstance(p, dict) and "description" in p:
        return p, p["description"]
    return p, p


def rec(group, name, method, path, expect, body=None, raw_path=False,
        desc="", binary=False):
    status, txt, ct = call(method, path, body, raw_path)
    parsed = None
    if not binary:
        try:
            parsed = json.loads(txt) if txt else None
        except json.JSONDecodeError:
            parsed = None
    ok = status == expect
    full_path = (BASE + path) if raw_path else (API + path)
    item = {
        "group": group, "name": name, "method": method,
        "url": full_path, "api_path": path, "raw_path": raw_path,
        "body": body, "expect_status": expect, "status": status,
        "ok": ok, "content_type": ct, "desc": desc,
        "response_text": (txt[:4000] if binary else txt),
        "response_json": parsed, "binary": binary,
    }
    records.append(item)
    flag = "PASS" if ok else f"FAIL(got {status})"
    print(f"{flag:14} | {method:6} {path}")
    return status, txt, parsed


def d(parsed):
    """Devuelve el payload util desenvolviendo {description:...} del toolkit."""
    if isinstance(parsed, dict) and "description" in parsed:
        return parsed["description"]
    return parsed


# ============================================================
# 01 — Health & Auth
# ============================================================
print("\n##### 01 - Health & Auth #####")
rec("01 - Health & Auth", "Health Check", "GET", "/health", 200, raw_path=True,
    desc="Liveness + ping a la DB. Sin auth ni prefijo /api/v1.")
rec("01 - Health & Auth", "Get Me", "GET", "/auth/me", 200,
    desc="Identidad autenticada. En ENV=test devuelve el usuario mock (test@educabot.com).")

# ============================================================
# 02 — Classrooms
# ============================================================
print("\n##### 02 - Classrooms #####")
rec("02 - Classrooms", "List Classrooms", "GET", "/classrooms", 200,
    desc="Lista las aulas de la organización.")
st, txt, p = rec("02 - Classrooms", "Create Classroom", "POST", "/classrooms", 201,
                 body={"name": "Aula Newman", "grade": "6", "section": "B"},
                 desc="Crea un aula. Guarda el id para los siguientes requests.")
ctx["classroom_id"] = (d(p) or {}).get("id")
cid = ctx["classroom_id"]
rec("02 - Classrooms", "Get Classroom (created)", "GET", f"/classrooms/{cid}", 200,
    desc="Detalle del aula recién creada.")
rec("02 - Classrooms", "Update Classroom (created)", "PUT", f"/classrooms/{cid}", 200,
    body={"name": "Aula Newman (editada)"},
    desc="Actualiza el aula creada.")
rec("02 - Classrooms", "List Classroom Students (seed 9001)", "GET",
    "/classrooms/9001/students", 200,
    desc="Alumnos del aula de prueba (9001): Tomás, Lucía, Mateo.")

# ============================================================
# 03 — Teachers
# ============================================================
print("\n##### 03 - Teachers #####")
rec("03 - Teachers", "List Teachers", "GET", "/teachers", 200,
    desc="Lista los docentes de la organización.")

# ============================================================
# 04 — Students
# ============================================================
print("\n##### 04 - Students #####")
rec("04 - Students", "List Students", "GET", "/students", 200,
    desc="Lista los alumnos de la organización.")
rec("04 - Students", "List Students (by classroom)", "GET",
    "/students?classroom_id=9001", 200,
    desc="Filtra alumnos por aula.")
rec("04 - Students", "Get Student (seed 9001)", "GET", "/students/9001", 200,
    desc="Detalle de Tomás Prueba (perfil rico de prueba).")
st, txt, p = rec("04 - Students", "Create Student", "POST", "/students", 201,
                 body={"name": "Alumno Newman", "classroom_id": 9001},
                 desc="Crea un alumno en el aula de prueba.")
ctx["student_id"] = (d(p) or {}).get("id")
sid = ctx["student_id"]
rec("04 - Students", "Update Student (created)", "PUT", f"/students/{sid}", 200,
    body={"name": "Alumno Newman (editado)"},
    desc="Actualiza el alumno creado.")

# ============================================================
# 05 — Student Profiles
# ============================================================
print("\n##### 05 - Student Profiles #####")
rec("05 - Student Profiles", "Get Student Profile (seed 9001)", "GET",
    "/students/9001/profile", 200,
    desc="Perfil del alumno (dificultades, descripción libre).")
rec("05 - Student Profiles", "Upsert Student Profile (created)", "PUT",
    f"/students/{sid}/profile", 200,
    body={"is_transitory": False,
          "difficulties": ["motricidad_fina", "atencion"],
          "free_description": "Prueba: dificultades de motricidad fina y atención."},
    desc="Crea/actualiza el perfil del alumno.")

# ============================================================
# 06 — Catalog (Ramps & Devices)
# ============================================================
print("\n##### 06 - Catalog #####")
rec("06 - Catalog", "List Ramps", "GET", "/ramps", 200,
    desc="Lista las rampas (agrupadores de dispositivos de la valija).")
rec("06 - Catalog", "Get Ramp (9001)", "GET", "/ramps/9001", 200,
    desc="Detalle de una rampa.")
rec("06 - Catalog", "List Devices", "GET", "/devices", 200,
    desc="Lista los dispositivos de la valija adaptativa.")
rec("06 - Catalog", "List Devices (by ramp)", "GET", "/devices?ramp_id=9001", 200,
    desc="Filtra dispositivos por rampa.")
rec("06 - Catalog", "Get Device (9001)", "GET", "/devices/9001", 200,
    desc="Detalle de un dispositivo.")

# ============================================================
# 07 — Adaptations
# ============================================================
print("\n##### 07 - Adaptations #####")
rec("07 - Adaptations", "List Adaptations", "GET", "/adaptations", 200,
    desc="Lista las adaptaciones (recursos) del docente.")
st, txt, p = rec("07 - Adaptations", "Create Adaptation", "POST", "/adaptations", 201,
                 body={"student_id": 9001, "device_id": 9001, "subject": "Matemática",
                       "title": "Fracciones con material concreto",
                       "activity_description": "Fracciones con tablero de anticipación",
                       "adaptation_strategy": "Anticipar pasos con apoyo visual",
                       "adaptation_type": "actividad_adaptada",
                       "notes": "Responde bien a estímulos visuales"},
                 desc="Crea una adaptación vinculada a un alumno y un dispositivo.")
ctx["adaptation_id"] = (d(p) or {}).get("id")
aid = ctx["adaptation_id"]
rec("07 - Adaptations", "List Adaptations (by student)", "GET",
    "/adaptations?student_id=9001", 200,
    desc="Filtra adaptaciones por alumno.")
rec("07 - Adaptations", "Get Adaptation (created)", "GET", f"/adaptations/{aid}", 200,
    desc="Detalle de la adaptación creada.")
rec("07 - Adaptations", "Update Adaptation (created)", "PUT", f"/adaptations/{aid}", 200,
    body={"status": "probado", "outcome": "Completó la actividad con más autonomía",
          "notes": "Funcionó, repetir la próxima clase"},
    desc="Actualiza estado/resultado de la adaptación.")
rec("07 - Adaptations", "List Adaptation Resources (created)", "GET",
    f"/adaptations/{aid}/resources", 200,
    desc="Recursos (materiales) ligados a la adaptación.")
rec("07 - Adaptations", "Export Adaptation (md)", "GET",
    f"/adaptations/{aid}/export?format=md", 200, binary=True,
    desc="Exporta la adaptación a Markdown (descarga). Content-Type text/markdown.")

# ============================================================
# 08 — Chat History
# ============================================================
print("\n##### 08 - Chat History #####")
rec("08 - Chat History", "Chat History (assist)", "GET", "/chat/history/assist", 200,
    desc="Historial de conversaciones del docente en el modo indicado (contextId).")

# ============================================================
# 09 — Dashboard
# ============================================================
print("\n##### 09 - Dashboard #####")
rec("09 - Dashboard", "Get Metrics", "GET", "/dashboard/metrics", 200,
    desc="Métricas agregadas (alumnos, adaptaciones, aulas).")
rec("09 - Dashboard", "Get AI Usage", "GET", "/dashboard/ai-usage", 200,
    desc="Consumo de tokens de IA por la organización.")

# ============================================================
# 10 — Context Engine · HU-1 apertura (/inclusion/open)
# ============================================================
print("\n##### 10 - HU-1 open #####")
g1 = "10 - Context Engine · HU-1 Apertura"
rec(g1, "Open · sin dimensión (saludo)", "POST", "/inclusion/open", 200, body={},
    desc="Sin dimensión: saluda y pregunta de qué hablar (needs_dimension=true).")
rec(g1, "Open · dimensión ambigua", "POST", "/inclusion/open", 200,
    body={"dimension": "loquesea"},
    desc="Dimensión no reconocida: repregunta.")
rec(g1, "Open · alumno sin id", "POST", "/inclusion/open", 200,
    body={"dimension": "alumno"},
    desc="Falta student_id: repregunta por el alumno.")
rec(g1, "Open · alumno 9001", "POST", "/inclusion/open", 200,
    body={"dimension": "alumno", "student_id": 9001},
    desc="Abre sobre Tomás (9001) y recupera resúmenes previos del alumno.")
rec(g1, "Open · alumno inexistente", "POST", "/inclusion/open", 404,
    body={"dimension": "alumno", "student_id": 999999},
    desc="Alumno inexistente: error.")
rec(g1, "Open · tema sin topic", "POST", "/inclusion/open", 200,
    body={"dimension": "tema"},
    desc="Falta el tema: repregunta.")
rec(g1, "Open · tema TEA", "POST", "/inclusion/open", 200,
    body={"dimension": "tema", "topic": "TEA"},
    desc="Abre sobre el tema TEA y recupera resúmenes por keyword.")
rec(g1, "Open · valija", "POST", "/inclusion/open", 200,
    body={"dimension": "valija"},
    desc="Abre sobre la valija (sin alumno).")

# ============================================================
# 11 — Context Engine · HU-2 contexto (/inclusion/context)
# ============================================================
print("\n##### 11 - HU-2 context #####")
g2 = "11 - Context Engine · HU-2 Contexto"
rec(g2, "Context · alumno 9001 (perfil rico)", "POST", "/inclusion/context", 200,
    body={"dimension": "alumno", "student_id": 9001},
    desc="Arma el contexto del alumno: perfil, PPI, diagnósticos, situaciones.")
rec(g2, "Context · alumno 9003 (degradación)", "POST", "/inclusion/context", 200,
    body={"dimension": "alumno", "student_id": 9003},
    desc="Alumno sin PPI: degrada con elegancia (missing_data incluye 'ppi').")
rec(g2, "Context · valija (lazy)", "POST", "/inclusion/context", 200,
    body={"dimension": "valija"},
    desc="Contexto de valija: catálogo + situaciones, sin target_student.")
rec(g2, "Context · alumno inexistente", "POST", "/inclusion/context", 404,
    body={"dimension": "alumno", "student_id": 999999},
    desc="Alumno inexistente: error.")

# ============================================================
# 12 — Context Engine · HU-3 RAG (/inclusion/search-content)
# ============================================================
print("\n##### 12 - HU-3 search-content #####")
g3 = "12 - Context Engine · HU-3 Contenido"
rec(g3, "Search · 'autismo desregula' (match TEA)", "POST", "/inclusion/search-content",
    200, body={"query": "autismo desregula"},
    desc="Busca material pedagógico por keyword/full-text. Devuelve el doc TEA.")
rec(g3, "Search · 'lectura dislexia'", "POST", "/inclusion/search-content", 200,
    body={"query": "lectura dislexia"},
    desc="Devuelve el documento de dislexia.")
rec(g3, "Search · 'quimica organica' (sin match)", "POST", "/inclusion/search-content",
    200, body={"query": "quimica organica"},
    desc="Sin coincidencias: no inventa, devuelve results vacío.")
rec(g3, "Search · query vacía", "POST", "/inclusion/search-content", 200,
    body={"query": ""},
    desc="Query vacía: results vacío.")

# ============================================================
# 13 — Context Engine · HU-4/IA (/inclusion/recommend, /assist)
# ============================================================
print("\n##### 13 - IA recommend/assist #####")
g4 = "13 - Context Engine · IA (recommend / assist)"
rec(g4, "Recommend Device", "POST", "/inclusion/recommend", 200,
    body={"student_id": 9001, "subject": "Lengua",
          "objective": "Producción de textos narrativos", "duration": "40",
          "dynamic": "Trabajo individual con apoyo docente",
          "materials": "Cuaderno, lápiz", "history": []},
    desc="Recomienda dispositivos de la valija para una actividad (LLM).")
st, txt, p = rec(g4, "Assist Classroom", "POST", "/inclusion/assist", 200,
                 body={"classroom_id": 9001, "student_id": 9001,
                       "message": "Tomás no logra sostener la atención al escribir. ¿Qué hago?",
                       "history": []},
                 desc="Chat asistente (loop agéntico con tools). Persiste la conversación.")
ctx["conversation_id"] = (d(p) or {}).get("conversation_id")

# ============================================================
# 14 — Context Engine · HU-5 memoria (/inclusion/close)
# ============================================================
print("\n##### 14 - HU-5 close #####")
g5 = "14 - Context Engine · HU-5 Memoria (close)"
conv = ctx.get("conversation_id") or 1
rec(g5, "Close Session (compacta conversación)", "POST", "/inclusion/close", 200,
    body={"conversation_id": conv},
    desc="Compacta la conversación al cerrar: genera el resumen + tags (alumno/tema/valija).")
rec(g5, "Close Session · conversación inexistente", "POST", "/inclusion/close", 404,
    body={"conversation_id": 999999},
    desc="Conversación inexistente: error.")
rec(g5, "Close Session · sin conversation_id", "POST", "/inclusion/close", 400,
    body={},
    desc="Falta conversation_id: validación (400).")

# ============================================================
# Cleanup (borra lo creado en este flujo)
# ============================================================
print("\n##### 15 - Cleanup #####")
gc = "15 - Cleanup"
rec(gc, "Delete Adaptation (created)", "DELETE", f"/adaptations/{aid}", 204,
    desc="Borra la adaptación creada.")
rec(gc, "Delete Student (created)", "DELETE", f"/students/{sid}", 204,
    desc="Borra el alumno creado.")
rec(gc, "Delete Classroom (created)", "DELETE", f"/classrooms/{cid}", 204,
    desc="Borra el aula creada.")

# ============================================================
# Escribir captura cruda
# ============================================================
capture_path = os.path.join(DOCS, "postman_capture.json")
with open(capture_path, "w", encoding="utf-8") as f:
    json.dump(records, f, ensure_ascii=False, indent=2)

passed = sum(1 for r in records if r["ok"])
print(f"\n================================")
print(f"CAPTURA: {passed}/{len(records)} con status esperado")
print(f"  -> {capture_path}")

# ============================================================
# Generar colección Postman v2.1
# ============================================================


def pretty(obj):
    return json.dumps(obj, ensure_ascii=False, indent=4)


def url_obj(api_path, raw_path):
    if raw_path:
        # /health (sin /api/v1)
        clean = api_path.lstrip("/")
        return {"raw": "{{host}}/" + clean, "host": ["{{host}}"],
                "path": clean.split("/")}
    # separa querystring
    path_part, _, qs = api_path.lstrip("/").partition("?")
    segs = path_part.split("/")
    u = {"raw": "{{base_url}}/" + api_path.lstrip("/"),
         "host": ["{{base_url}}"], "path": segs}
    if qs:
        query = []
        for kv in qs.split("&"):
            k, _, v = kv.partition("=")
            query.append({"key": k, "value": v})
        u["query"] = query
    return u


def test_script(r):
    lines = [f"pm.test('status {r['expect_status']}', function () {{ "
             f"pm.response.to.have.status({r['expect_status']}); }});"]
    # guardar ids creados
    nm = r["name"]
    if nm == "Create Classroom":
        lines += ["var o = pm.response.json(); o = o.description || o;",
                  "pm.collectionVariables.set('created_classroom_id', o.id);"]
    elif nm == "Create Student":
        lines += ["var o = pm.response.json(); o = o.description || o;",
                  "pm.collectionVariables.set('created_student_id', o.id);"]
    elif nm == "Create Adaptation":
        lines += ["var o = pm.response.json(); o = o.description || o;",
                  "pm.collectionVariables.set('created_adaptation_id', o.id);"]
    elif nm == "Assist Classroom":
        lines += ["var o = pm.response.json(); o = o.description || o;",
                  "pm.collectionVariables.set('conversation_id', o.conversation_id);",
                  "pm.test('has response', function () { pm.expect(o.response).to.be.a('string'); });"]
    return lines


def req_url_for_collection(r):
    # Reemplaza ids dinámicos por variables de colección en la URL.
    path = r["api_path"]
    repl = {
        f"/classrooms/{ctx.get('classroom_id')}": "/classrooms/{{created_classroom_id}}",
        f"/students/{ctx.get('student_id')}/profile": "/students/{{created_student_id}}/profile",
        f"/students/{ctx.get('student_id')}": "/students/{{created_student_id}}",
        f"/adaptations/{ctx.get('adaptation_id')}/resources": "/adaptations/{{created_adaptation_id}}/resources",
        f"/adaptations/{ctx.get('adaptation_id')}/export": "/adaptations/{{created_adaptation_id}}/export",
        f"/adaptations/{ctx.get('adaptation_id')}": "/adaptations/{{created_adaptation_id}}",
    }
    for real, var in repl.items():
        if ctx.get('classroom_id') and path.startswith(real):
            path = path.replace(real, var, 1)
            break
    body = r["body"]
    if r["name"] == "Close Session (compacta conversación)":
        body = {"conversation_id": "{{conversation_id}}"} if False else body
    return path, body


def build_request(r):
    path, body = req_url_for_collection(r)
    request = {"method": r["method"], "header": [], "url": url_obj(path, r["raw_path"])}
    if r["raw_path"]:
        request["auth"] = {"type": "noauth"}
    if body is not None:
        request["header"].append({"key": "Content-Type", "value": "application/json"})
        request["body"] = {"mode": "raw", "raw": pretty(body),
                           "options": {"raw": {"language": "json"}}}
    # respuesta de ejemplo capturada
    resp_body = r["response_text"]
    example = [{
        "name": f"Ejemplo ({r['status']})",
        "originalRequest": {"method": r["method"], "header": request["header"],
                            "url": request["url"],
                            **({"body": request["body"]} if "body" in request else {})},
        "status": "OK" if r["status"] < 300 else "Error",
        "code": r["status"],
        "_postman_previewlanguage": "json",
        "header": [{"key": "Content-Type", "value": r["content_type"] or "application/json"}],
        "body": resp_body,
    }]
    return {
        "name": r["name"],
        "event": [{"listen": "test", "script": {"type": "text/javascript",
                                                 "exec": test_script(r)}}],
        "request": {**request, "description": r["desc"]},
        "response": example,
    }


# agrupar por carpeta
groups = {}
order = []
for r in records:
    if r["group"] not in groups:
        groups[r["group"]] = []
        order.append(r["group"])
    groups[r["group"]].append(build_request(r))

collection = {
    "info": {
        "name": "alizia-inclusion (ENV=test)",
        "description": (
            "Colección completa de alizia-inclusion-be — TODOS los endpoints "
            f"({len(records)} requests).\n\n"
            "Pensada para correr contra el server en **ENV=test** (auth mockeada: "
            "user 1 / org 00000000-...-0001, NO requiere token) con los seeds de "
            "prueba aplicados (alumnos 9001-9003, aula 9001, ramp 9001, devices "
            "9001-9002, contenido pedagógico 1-3).\n\n"
            "Flujo secuencial (Runner/Newman): los Create guardan ids en variables, "
            "Assist guarda conversation_id, Cleanup borra lo creado. Cada request "
            "trae un Ejemplo de response REAL capturado.\n\n"
            "Para un deploy real detrás del auth-service, completá la variable "
            "{{token}} (auth bearer a nivel colección)."
        ),
        "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
    },
    "variable": [
        {"key": "host", "value": "http://localhost:8080"},
        {"key": "base_url", "value": "http://localhost:8080/api/v1"},
        {"key": "token", "value": ""},
        {"key": "created_classroom_id", "value": ""},
        {"key": "created_student_id", "value": ""},
        {"key": "created_adaptation_id", "value": ""},
        {"key": "conversation_id", "value": "1"},
    ],
    "auth": {"type": "bearer", "bearer": [{"key": "token", "value": "{{token}}", "type": "string"}]},
    "item": [{"name": g, "item": groups[g]} for g in order],
}

coll_path = os.path.join(ROOT, "alizia-inclusion.postman_collection.json")
with open(coll_path, "w", encoding="utf-8") as f:
    json.dump(collection, f, ensure_ascii=False, indent=2)
print(f"  -> {coll_path}  ({len(records)} requests)")

# ============================================================
# Generar ENDPOINTS.md
# ============================================================


def fmt_json_block(text, limit=1500):
    if text is None or text == "":
        return "_(vacío)_"
    try:
        obj = json.loads(text)
        pretty_txt = json.dumps(obj, ensure_ascii=False, indent=2)
    except (json.JSONDecodeError, TypeError):
        pretty_txt = text
    if len(pretty_txt) > limit:
        pretty_txt = pretty_txt[:limit] + "\n... (truncado)"
    return "```json\n" + pretty_txt + "\n```"


md = []
md.append("# API Reference — alizia-inclusion-be\n")
md.append("> Generado por `scripts/capture_postman.py` con respuestas **reales** "
          "capturadas contra el server en `ENV=test` (auth mockeada, org de ceros, "
          "seeds de prueba).\n")
md.append(f"> Total: **{len(records)} endpoints/casos**. "
          f"`{passed}/{len(records)}` devolvieron el status esperado.\n")
md.append("\n## Notas\n")
md.append("- **Auth**: en `ENV=test` no se requiere token (middleware mock). "
          "En prod, bearer token del auth-service.\n")
md.append("- **Envoltura**: las respuestas OK del toolkit vienen como "
          "`{\"description\": <payload>}`.\n")
md.append("- **IA**: con el stub client (sin Azure key) `recommend`/`assist`/`close` "
          "devuelven contenido `[stub]` determinista.\n")

cur = None
for r in records:
    if r["group"] != cur:
        cur = r["group"]
        md.append(f"\n---\n\n## {cur}\n")
    md.append(f"\n### {r['name']}\n")
    md.append(f"{r['desc']}\n")
    md.append(f"\n**`{r['method']} {r['url']}`** — expected `{r['expect_status']}`, "
              f"got `{r['status']}` {'✅' if r['ok'] else '❌'}\n")
    if r["body"] is not None:
        md.append("\n_Request body:_\n")
        md.append(fmt_json_block(json.dumps(r["body"], ensure_ascii=False)))
    md.append("\n_Response:_\n")
    if r["binary"]:
        md.append(f"```\nContent-Type: {r['content_type']}\n"
                  f"{r['response_text'][:600]}\n```")
    else:
        md.append(fmt_json_block(r["response_text"]))
    md.append("")

md_path = os.path.join(DOCS, "ENDPOINTS.md")
with open(md_path, "w", encoding="utf-8") as f:
    f.write("\n".join(md))
print(f"  -> {md_path}")

sys.exit(0 if passed == len(records) else 1)
