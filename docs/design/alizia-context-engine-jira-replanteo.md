# Replanteo del backlog Jira contra el criterio de Juan/Francisco

> **Estado:** ANÁLISIS + propuesta. No ejecutado en Jira. Para revisar con Juan/Francisco.
> **Fecha:** 2026-06-01 · Proyecto ALZ.

## 1. El criterio (Slack 2026-06-01)

- **Historia** = un caso de uso **testeable de forma completa**, con condiciones previas y resultado esperado.
- Formato: **"Como {user} quiero {resultado} para {motivo}"**.
- **Lenguaje común, no técnico** (Francisco: para gestionarlo o derivarlo a un agente sin preguntar).
- Subtarea = accionable técnico, delegable a agentes.

## 2. Auditoría del árbol actual (lo que está hoy en Jira)

¿Cada HU actual es un caso de uso de cara a un *user*, o es infra técnica disfrazada de historia?

| HU actual | ¿Caso de uso testeable de un user? | Veredicto |
|---|---|---|
| HU-1.1 Trazar cada turno | No tiene user; es instrumentación | ❌ técnica |
| HU-2.1 Usa el contexto del alumno | Sí — **docente** | ✅ historia |
| HU-2.2 Tools agénticas | Mecanismo interno, no caso de uso | ⚠️ habilitador |
| HU-3.1 Editar/publicar prompt sin deploy | Sí — **producto** (interno) | ✅ historia (interna) |
| HU-4.1 Conversación larga conserva memoria | Sí — **docente** | ✅ historia |
| HU-4.2 Memoria viva del alumno | Sí — **docente** | ✅ historia |
| HU-5.1 Golden alimentan few-shot | Habilitador de calidad | ❌ técnica |
| HU-5.2 Medir si mejora | Sí — **equipo** | ✅ historia (interna) |
| HU-5.3 Promover golden + A/B | Mecanismo interno | ❌ técnica |
| HU-6.1 Contenido pedagógico base | Trabajo de autoría, no caso de uso | ⚠️ habilitador |

**Conclusión:** de 10 "historias", solo ~5 son casos de uso reales. El resto son piezas técnicas. El árbol sigue siendo **demasiado de plataforma** — el mismo problema que señaló Juan, solo que un nivel más arriba. Las 6 épicas (Traza, Prompts versionados, Flywheel…) son **capas técnicas del sistema**, no objetivos de producto.

## 3. La decisión de fondo: dos modelos incompatibles

**Modelo PLATAFORMA (lo que hay hoy):** el backlog refleja la arquitectura del Context Engine (6 capas técnicas). Bueno para devs; malo para el criterio de Juan (historias técnicas).

**Modelo PRODUCTO (lo que pide Juan):** el backlog se organiza por **casos de uso del usuario**. Pocas historias `Como… quiero… para…`, y **todo el Context Engine baja a subtareas** que las habilitan. El "Context Engine" deja de verse como tal; se disuelve en habilitadores.

Seguir a Juan al pie de la letra = pasar al Modelo Producto.

## 4. Propuesta — árbol Modelo Producto

4 épicas (objetivos de producto) · 8 historias (casos de uso testeables) · las **mismas 30 tareas** técnicas, redistribuidas como subtareas.

### EP · Alizia entiende a cada alumno
- **HU · Como docente quiero que Alizia tenga en cuenta el contexto de mi alumno (perfil, situaciones, PPI) al responder, para que las sugerencias sean pertinentes y no genéricas.**
  Subtareas: teacher_profiles · students/student_profiles · situations_catalog · diagnósticos · ppi · rol integradora · Context Assembler *(T-2.1…2.7)*
- **HU · Como docente quiero que Alizia traiga más datos del alumno cuando hace falta, para que profundice sin que yo le repita todo.**
  Subtareas: tools get_student_* *(T-2.8)*

### EP · Alizia con memoria
- **HU · Como docente quiero que Alizia no pierda el hilo en conversaciones largas, para seguir sin recontextualizar.**
  Subtareas: conversation_summaries + job · usar resumen en capMessages *(T-4.1, T-4.2)*
- **HU · Como docente quiero que Alizia recuerde entre sesiones qué funcionó con mi alumno, para no explicárselo cada vez.**
  Subtareas: student_insights + job *(T-4.3)*

### EP · Alizia mejora con el uso
- **HU · Como equipo de Alizia quiero medir si las recomendaciones mejoran versión a versión, para decidir con datos.**
  Subtareas: ai_usage con traza · señal de aceptación implícita · win-rate + métricas *(T-1.1, T-1.2, T-5.3)*
- **HU · Como equipo de producto quiero ajustar y publicar el comportamiento de Alizia sin deploy (con validación + fallback), para iterar rápido y sin riesgo.**
  Subtareas: modelo versionado · renderer · validación · fallback · migrar capa 1 · seed de los 3 body *(T-3.1…3.6)*
- **HU · Como equipo de Alizia quiero que Alizia mejore sola promoviendo buenos ejemplos, para subir la calidad sin trabajo manual.**
  Subtareas: response_examples + seed · selección few-shot · A/B + eval · job batch flywheel *(T-5.1, T-5.2, T-5.4, T-5.5)*

### EP · Alizia con criterio pedagógico
- **HU · Como equipo pedagógico quiero que las respuestas sigan nuestros lineamientos de inclusión, para confiar en lo que Alizia sugiere.**
  Subtareas: identidad/persona · marco pedagógico · límites · formato · ~15 situaciones · few-shot golden *(T-6.1…6.6)*

## 5. Trade-offs honestos

- ✅ Las 8 historias son casos de uso testeables, en formato y lenguaje común. Cumple Juan/Francisco.
- ✅ Las 30 tareas técnicas no se pierden — pasan a subtareas.
- ⚠️ Se **disuelve** el Context Engine como estructura técnica visible (las "capas" ya no son épicas).
- ⚠️ Hay **historias grandes** (la de "iterar sin deploy" tiene 6 subtareas técnicas pesadas). Testeable como caso de uso, pero gorda.
- ⚠️ Implica **rehacer Jira otra vez**: re-crear 4 épicas + 8 historias, reparentar las 30 tareas, borrar la estructura EP-1..6 / HU actuales.

## 6. Decisión pendiente

1. ¿Vamos al **Modelo Producto** (este árbol) o mantenemos el **Modelo Plataforma** actual con las historias reescritas en formato Como/quiero/para?
2. Si Modelo Producto: ¿los nombres de épicas/historias así, o ajustamos?
3. ¿Validamos esto con Juan/Francisco antes de re-ejecutar en Jira? (recomendado: ya rehicimos una vez).
