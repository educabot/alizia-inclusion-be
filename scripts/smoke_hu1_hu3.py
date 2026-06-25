#!/usr/bin/env python3
"""Smoke / Postman-equivalente para HU-1 a HU-3 (endpoints sin LLM).

Requiere el server corriendo con ENV=test (auth mockeada: user 1 / org de ceros)
y los seeds de prueba aplicados (alumnos 9001-9003 + corpus pedagógico).

Uso:  python scripts/smoke_hu1_hu3.py [base_url]
"""
import json
import sys
import urllib.request

BASE = sys.argv[1] if len(sys.argv) > 1 else "http://127.0.0.1:8080"
API = BASE + "/api/v1"

passed = 0
failed = 0


def call(method, path, body=None):
    url = BASE + path if path.startswith("/health") else API + path
    data = json.dumps(body).encode() if body is not None else None
    req = urllib.request.Request(url, data=data, method=method,
                                 headers={"Content-Type": "application/json"})
    try:
        with urllib.request.urlopen(req, timeout=15) as r:
            return r.status, json.loads(r.read().decode() or "{}")
    except urllib.error.HTTPError as e:
        return e.code, json.loads(e.read().decode() or "{}")


def unwrap(payload):
    # El toolkit envuelve la respuesta OK en {"description": {...}}.
    return payload.get("description", payload) if isinstance(payload, dict) else payload


def check(name, cond):
    global passed, failed
    if cond:
        passed += 1
        print(f"PASS | {name}")
    else:
        failed += 1
        print(f"FAIL | {name}")


def section(t):
    print(f"\n===== {t} =====")


# ---------------- HEALTH ----------------
section("HEALTH")
st, body = call("GET", "/health")
check("health 200 + db ok", st == 200 and body.get("db") == "ok")

# ---------------- HU-1 · /inclusion/open ----------------
section("HU-1 · apertura guiada (router)")

st, b = call("POST", "/inclusion/open", {})
d = unwrap(b)
check("1. sin dimension -> pregunta + needs_dimension",
      st == 200 and d.get("needs_dimension") is True and "alumno" in d.get("greeting", ""))

st, b = call("POST", "/inclusion/open", {"dimension": "loquesea"})
d = unwrap(b)
check("2. dimension ambigua -> reclarifica",
      st == 200 and d.get("needs_dimension") is True and "No me qued" in d.get("greeting", ""))

st, b = call("POST", "/inclusion/open", {"dimension": "alumno"})
d = unwrap(b)
check("3. alumno sin id -> repregunta", st == 200 and d.get("needs_dimension") is True)

st, b = call("POST", "/inclusion/open", {"dimension": "alumno", "student_id": 9001})
d = unwrap(b)
check("4. alumno 9001 -> Tomás + dimension alumno",
      st == 200 and d.get("dimension") == "alumno" and not d.get("needs_dimension")
      and d.get("student", {}).get("name") == "Tomás Prueba")

st, b = call("POST", "/inclusion/open", {"dimension": "alumno", "student_id": 9002})
d = unwrap(b)
check("5. alumno 9002 -> Lucía", st == 200 and d.get("student", {}).get("name") == "Lucía Demo")

st, b = call("POST", "/inclusion/open", {"dimension": "alumno", "student_id": 999999})
check("6. alumno inexistente -> error 4xx/5xx", st >= 400)

st, b = call("POST", "/inclusion/open", {"dimension": "tema"})
d = unwrap(b)
check("7. tema sin topic -> repregunta", st == 200 and d.get("needs_dimension") is True)

st, b = call("POST", "/inclusion/open", {"dimension": "tema", "topic": "TEA"})
d = unwrap(b)
check("8. tema TEA -> dimension tema", st == 200 and d.get("dimension") == "tema" and "TEA" in d.get("greeting", ""))

st, b = call("POST", "/inclusion/open", {"dimension": "valija"})
d = unwrap(b)
check("9. valija -> dimension valija sin alumno",
      st == 200 and d.get("dimension") == "valija" and not d.get("needs_dimension") and d.get("student") is None)

# ---------------- HU-2 · /inclusion/context ----------------
section("HU-2 · Context Assembler")

st, b = call("POST", "/inclusion/context", {"dimension": "alumno", "student_id": 9001})
d = unwrap(b)
ts = d.get("target_student") or {}
check("10. ctx 9001 -> target Tomás + perfil rico",
      st == 200 and ts.get("name") == "Tomás Prueba" and ts.get("profile", {}).get("support_level") == "medio")
check("11. ctx 9001 -> PPI presente", bool(d.get("ppi")))
check("12. ctx 9001 -> diagnósticos presentes", len(d.get("diagnoses") or []) >= 1)
check("13. ctx 9001 -> situaciones (estático) cargadas", len(d.get("situations") or []) >= 10)
check("14. ctx 9001 -> missing_data incluye perfil_docente", "perfil_docente" in (d.get("missing_data") or []))

st, b = call("POST", "/inclusion/context", {"dimension": "alumno", "student_id": 9003})
d = unwrap(b)
check("15. ctx 9003 -> sin PPI (degradación)", d.get("ppi") is None and "ppi" in (d.get("missing_data") or []))

st, b = call("POST", "/inclusion/context", {"dimension": "valija"})
d = unwrap(b)
check("16. ctx valija -> sin target_student (lazy)", st == 200 and d.get("target_student") is None and len(d.get("situations") or []) >= 10)

st, b = call("POST", "/inclusion/context", {"dimension": "alumno", "student_id": 999999})
check("17. ctx alumno inexistente -> error", st >= 400)

# ---------------- HU-3 · /inclusion/search-content ----------------
section("HU-3 · RAG contenido pedagógico")

st, b = call("POST", "/inclusion/search-content", {"query": "autismo desregula"})
d = unwrap(b)
res = d.get("results") or []
check("18. 'autismo desregula' -> match TEA",
      st == 200 and len(res) >= 1 and "TEA" in res[0].get("title", ""))

st, b = call("POST", "/inclusion/search-content", {"query": "estrategias TEA atencion"})
d = unwrap(b)
res = d.get("results") or []
check("19. 'estrategias TEA atencion' -> >=2 y TEA primero",
      st == 200 and len(res) >= 2 and res[0].get("content_id") == 2
      and res[0].get("score", 0) >= res[1].get("score", 0))

st, b = call("POST", "/inclusion/search-content", {"query": "lectura dislexia"})
d = unwrap(b)
res = d.get("results") or []
check("20. 'lectura dislexia' -> doc dislexia (id 3)", st == 200 and len(res) >= 1 and res[0].get("content_id") == 3)

st, b = call("POST", "/inclusion/search-content", {"query": "quimica organica"})
d = unwrap(b)
check("21. 'quimica organica' -> sin match (no inventa)", st == 200 and (d.get("results") or []) == [])

st, b = call("POST", "/inclusion/search-content", {"query": ""})
d = unwrap(b)
check("22. query vacía -> resultados vacíos", st == 200 and (d.get("results") or []) == [])

# ---------------- TOTAL ----------------
print(f"\n================================")
print(f"RESULTADO: {passed} PASS / {failed} FAIL  (total {passed + failed})")
sys.exit(1 if failed else 0)
