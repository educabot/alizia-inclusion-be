# Flujos de producto — Director y Profesor

> Visión **funcional / de usuario**. Qué puede hacer cada rol, qué ve y cómo se mueve por la app.
> Sin detalle técnico ni de base de datos.

---

## 🧭 En una frase

- **Director** → mira y organiza. Crea las aulas, ve a los docentes y consulta cómo va todo (tablero con números).
- **Profesor** → trabaja el día a día. Carga a sus alumnos, les arma el perfil, planifica con ayuda de la IA (Alizia) y crea las adaptaciones para cada chico.

---

## 👔 DIRECTOR — Qué puede hacer

```mermaid
flowchart TD
    D([Director entra a la app])

    D --> A[🏫 Gestionar aulas]
    D --> T[👩‍🏫 Ver docentes]
    D --> M[📊 Ver el tablero]

    A --> A1[Crear un aula nueva]
    A --> A2[Editar datos de un aula]
    A --> A3[Eliminar un aula]
    A --> A4[Ver lista de aulas]

    T --> T1[Ver qué docentes hay en la escuela]

    M --> M1[Cuántos alumnos / adaptaciones hay]
    M --> M2[Cómo vienen las adaptaciones<br/>en curso, funcionaron, etc.]
    M --> M3[Qué recursos se usan más]
    M --> M4[Cuánto se usó la IA Alizia]

    style D fill:#1565c0,color:#fff
    style A fill:#e3f2fd
    style T fill:#e3f2fd
    style M fill:#e3f2fd
```

### Su día típico
1. Entra y ve el **tablero general** de la escuela.
2. Si arranca un ciclo nuevo, **crea las aulas**.
3. Revisa **qué docentes** están cargados.
4. Cada tanto mira los **números**: cuántas adaptaciones se hicieron, cómo van, cuánto se apoyó el equipo en la IA.

**El director NO crea alumnos ni adaptaciones.** Su rol es organizar y supervisar.

---

## 👩‍🏫 PROFESOR — Qué puede hacer

```mermaid
flowchart TD
    P([Profesor entra a la app])

    P --> AL[🧑‍🎓 Gestionar alumnos]
    P --> PE[📋 Armar perfil del alumno]
    P --> IA[🤖 Pedir ayuda a Alizia IA]
    P --> AD[📝 Crear adaptaciones]
    P --> EX[📄 Exportar / compartir]

    AL --> AL1[Crear un alumno]
    AL --> AL2[Editar / eliminar alumno]
    AL --> AL3[Ver alumnos de su aula]

    PE --> PE1[Cargar necesidades,<br/>diagnóstico, apoyos que necesita]

    IA --> IA1[Recomendación rápida:<br/>qué recurso le sirve a este chico]
    IA --> IA2[Asistente conversacional:<br/>charlar y planificar el aula]

    AD --> AD1[Crear adaptación para un alumno]
    AD --> AD2[Marcar cómo le fue<br/>en curso → funcionó / ajustar]
    AD --> AD3[Editar / eliminar]

    EX --> EX1[Descargar la adaptación en PDF]

    style P fill:#2e7d32,color:#fff
    style AL fill:#e8f5e9
    style PE fill:#e8f5e9
    style IA fill:#fff3e0
    style AD fill:#e8f5e9
    style EX fill:#e8f5e9
```

### Su día típico (recorrido completo)

```mermaid
flowchart LR
    S1[1️⃣ Carga<br/>al alumno] --> S2[2️⃣ Arma su perfil<br/>necesidades/apoyos]
    S2 --> S3[3️⃣ Le pregunta a Alizia<br/>qué hacer]
    S3 --> S4[4️⃣ Crea la<br/>adaptación]
    S4 --> S5[5️⃣ La aplica en el aula<br/>y marca el resultado]
    S5 --> S6[6️⃣ Exporta el PDF<br/>para compartir]

    style S1 fill:#e8f5e9
    style S2 fill:#e8f5e9
    style S3 fill:#fff3e0
    style S4 fill:#e8f5e9
    style S5 fill:#e8f5e9
    style S6 fill:#e8f5e9
```

---

## 🤖 Cómo funciona la ayuda de Alizia (IA) — para el Profesor

Alizia tiene **dos modos**, según lo que necesite el docente:

```mermaid
flowchart TD
    Q{¿Qué necesita<br/>el profesor?}

    Q -->|Algo puntual y rápido| R[💡 Recomendación]
    Q -->|Charlar y planificar| C[💬 Asistente]

    R --> R1[Le da una sugerencia<br/>de recurso para ese alumno]

    C --> C1[Conversa con el docente]
    C1 --> C2[Alizia consulta sola<br/>los datos del aula y alumnos]
    C2 --> C3[Propone una adaptación<br/>lista para guardar]
    C3 --> C4[El profesor la acepta<br/>y queda creada]

    style Q fill:#ef6c00,color:#fff
    style R fill:#fff3e0
    style C fill:#fff3e0
```

- **Recomendación** → respuesta rápida tipo "para este chico te sirve tal recurso".
- **Asistente** → es una conversación. Alizia entiende el contexto del aula, propone una adaptación concreta y el docente la puede guardar de una.

---

## ⚖️ Director vs Profesor — de un vistazo

| | 👔 Director | 👩‍🏫 Profesor |
|---|---|---|
| **Foco** | Organizar y supervisar | Trabajar con cada alumno |
| **Crea aulas** | ✅ Sí | ❌ No |
| **Ve docentes** | ✅ Sí | — |
| **Crea alumnos** | ❌ No | ✅ Sí |
| **Arma perfiles** | ❌ No | ✅ Sí |
| **Usa la IA Alizia** | ❌ No (solo ve cuánto se usó) | ✅ Sí, es su herramienta clave |
| **Crea adaptaciones** | ❌ No | ✅ Sí |
| **Exporta PDF** | — | ✅ Sí |
| **Ve el tablero con números** | ✅ Sí | — |

---

## 📌 Nota importante sobre los roles

Hoy la app **no bloquea por rol técnicamente**: cualquier usuario logueado podría, en teoría, entrar a las pantallas del otro. La separación "Director hace esto / Profesor hace esto otro" es **de producto y de cómo se diseñó la experiencia**, no una restricción forzada por el sistema.

> Si el negocio necesita que un profesor **no pueda** tocar lo del director (o viceversa), eso hay que pedirlo como una funcionalidad a agregar.
