#!/usr/bin/env python3
"""Runner único de smoke HU-1 a HU-4 (endpoints sin LLM).

Corre las suites smoke_hu1_hu3.py y smoke_hu4.py y agrega el resultado.
Requiere el server con ENV=test y los seeds de prueba aplicados.

Uso:  python scripts/smoke_all.py [base_url]
"""
import os
import subprocess
import sys

BASE = sys.argv[1] if len(sys.argv) > 1 else "http://127.0.0.1:8080"
HERE = os.path.dirname(os.path.abspath(__file__))
SUITES = ["smoke_hu1_hu3.py", "smoke_hu4.py"]

rc_total = 0
for suite in SUITES:
    print(f"\n########## {suite} ##########", flush=True)
    rc = subprocess.call([sys.executable, os.path.join(HERE, suite), BASE])
    rc_total = rc_total or rc

print("\n================================")
print("RESULTADO GLOBAL HU-1..HU-4:", "TODO OK" if rc_total == 0 else "HAY FALLOS")
sys.exit(rc_total)
