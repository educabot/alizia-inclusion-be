#!/usr/bin/env python3
"""Smoke / Postman-equivalente para HU-4 (recomendar valija + guardar adaptaciones).

Valida la persistencia con origen IA + was_edited, el scope privado por docente
y los filtros del listado (alumno / material / tema), sin LLM. El guardado
confirmado de las tools de acción (create_recurso/create_student con flag
confirmed) está cubierto por los unit tests del dispatcher (agentic_hu4_test.go),
ya que viven en el loop agéntico que requiere Azure.

Requiere ENV=test y los seeds de prueba (incluye user 1 + alumno 9001).

Uso:  python scripts/smoke_hu4.py [base_url]
"""
import json
import sys
import urllib.request
import urllib.parse

BASE = sys.argv[1] if len(sys.argv) > 1 else "http://127.0.0.1:8080"
API = BASE + "/api/v1"

passed = 0
failed = 0


def call(method, path, body=None):
    data = json.dumps(body).encode() if body is not None else None
    req = urllib.request.Request(API + path, data=data, method=method,
                                 headers={"Content-Type": "application/json"})
    try:
        with urllib.request.urlopen(req, timeout=15) as r:
            return r.status, json.loads(r.read().decode() or "{}")
    except urllib.error.HTTPError as e:
        return e.code, json.loads(e.read().decode() or "{}")


def unwrap(payload):
    return payload.get("description", payload) if isinstance(payload, dict) else payload


def check(name, cond):
    global passed, failed
    if cond:
        passed += 1
        print(f"PASS | {name}")
    else:
        failed += 1
        print(f"FAIL | {name}")


print("===== HU-4 · recomendar valija + guardar adaptaciones =====")

# 1. Crear recurso con origen IA + was_edited (persistencia).
marker = "HU4 smoke narrativa"
st, b = call("POST", "/adaptations", {
    "student_id": 9001, "subject": marker, "title": "Texto narrativo",
    "activity_description": "Producción guiada paso a paso", "was_edited": True,
})
created = unwrap(b)
created_id = created.get("id")
check("1. crear recurso -> 201 + teacher_id del docente mock + was_edited",
      st == 201 and created.get("teacher_id") == 1 and created.get("was_edited") is True)

# 2. Listado privado del docente: todos los resultados son del teacher 1.
st, b = call("GET", "/adaptations?student_id=9001")
lst = unwrap(b) or []
teachers = {a.get("teacher_id") for a in lst}
check("2. listado privado del docente (solo teacher_id=1)",
      st == 200 and teachers.issubset({1}) and any(a.get("id") == created_id for a in lst))

# 3. Filtro por tema (q) encuentra el recurso recién creado.
st, b = call("GET", "/adaptations?q=" + urllib.parse.quote("narrativo"))
lst = unwrap(b) or []
check("3. filtro por tema q=narrativo encuentra el recurso",
      st == 200 and any(a.get("id") == created_id for a in lst))

# 4. Filtro por tema inexistente -> vacío.
st, b = call("GET", "/adaptations?q=zzznoexistexyz")
check("4. filtro por tema inexistente -> vacío", st == 200 and (unwrap(b) or []) == [])

# 5. Filtro por alumno inexistente -> vacío.
st, b = call("GET", "/adaptations?student_id=999999")
check("5. filtro por alumno inexistente -> vacío", st == 200 and (unwrap(b) or []) == [])

print("\n================================")
print(f"RESULTADO HU-4: {passed} PASS / {failed} FAIL  (total {passed + failed})")
sys.exit(1 if failed else 0)
