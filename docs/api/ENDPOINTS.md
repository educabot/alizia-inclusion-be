# API Reference — alizia-inclusion-be

> Generado por `scripts/capture_postman.py` con respuestas **reales** capturadas contra el server en `ENV=test` (auth mockeada, org de ceros, seeds de prueba).

> Total: **54 endpoints/casos**. `54/54` devolvieron el status esperado.


## Notas

- **Auth**: en `ENV=test` no se requiere token (middleware mock). En prod, bearer token del auth-service.

- **Envoltura**: las respuestas OK del toolkit vienen como `{"description": <payload>}`.

- **IA**: con el stub client (sin Azure key) `recommend`/`assist`/`close` devuelven contenido `[stub]` determinista.


---

## 01 - Health & Auth


### Health Check

Liveness + ping a la DB. Sin auth ni prefijo /api/v1.


**`GET http://127.0.0.1:8080/health`** — expected `200`, got `200` ✅


_Response:_

```json
{
  "db": "ok",
  "status": "ok"
}
```


### Get Me

Identidad autenticada. En ENV=test devuelve el usuario mock (test@educabot.com).


**`GET http://127.0.0.1:8080/api/v1/auth/me`** — expected `200`, got `200` ✅


_Response:_

```json
{
  "description": {
    "id": 1,
    "name": "Test User",
    "email": "test@educabot.com",
    "role": "teacher"
  }
}
```


---

## 02 - Classrooms


### List Classrooms

Lista las aulas de la organización.


**`GET http://127.0.0.1:8080/api/v1/classrooms`** — expected `200`, got `200` ✅


_Response:_

```json
{
  "description": [
    {
      "id": 9001,
      "name": "Aula de prueba",
      "grade": "4to",
      "section": "A",
      "student_count": 3
    }
  ]
}
```


### Create Classroom

Crea un aula. Guarda el id para los siguientes requests.


**`POST http://127.0.0.1:8080/api/v1/classrooms`** — expected `201`, got `201` ✅


_Request body:_

```json
{
  "name": "Aula Newman",
  "grade": "6",
  "section": "B"
}
```

_Response:_

```json
{
  "id": 1,
  "name": "Aula Newman",
  "grade": "6",
  "section": "B",
  "student_count": 0
}
```


### Get Classroom (created)

Detalle del aula recién creada.


**`GET http://127.0.0.1:8080/api/v1/classrooms/1`** — expected `200`, got `200` ✅


_Response:_

```json
{
  "description": {
    "id": 1,
    "name": "Aula Newman",
    "grade": "6",
    "section": "B",
    "student_count": 0
  }
}
```


### Update Classroom (created)

Actualiza el aula creada.


**`PUT http://127.0.0.1:8080/api/v1/classrooms/1`** — expected `200`, got `200` ✅


_Request body:_

```json
{
  "name": "Aula Newman (editada)"
}
```

_Response:_

```json
{
  "description": {
    "id": 1,
    "name": "Aula Newman (editada)",
    "grade": "6",
    "section": "B",
    "student_count": 0
  }
}
```


### List Classroom Students (seed 9001)

Alumnos del aula de prueba (9001): Tomás, Lucía, Mateo.


**`GET http://127.0.0.1:8080/api/v1/classrooms/9001/students`** — expected `200`, got `200` ✅


_Response:_

```json
{
  "description": [
    {
      "id": 9002,
      "name": "Lucía Demo",
      "classroom_id": 9001,
      "profile": {
        "id": 2,
        "student_id": 9002,
        "student_name": "Lucía Demo",
        "is_transitory": true,
        "difficulties": [
          "dificultad_lectura"
        ],
        "free_description": "Dificultad transitoria de lectura tras cambio de escuela; en proceso de adaptación."
      },
      "created_at": "2026-06-10T21:22:47Z"
    },
    {
      "id": 9003,
      "name": "Mateo Test",
      "classroom_id": 9001,
      "profile": {
        "id": 3,
        "student_id": 9003,
        "student_name": "Mateo Test",
        "is_transitory": false,
        "difficulties": [
          "hipersensibilidad_sensorial",
          "rutinas_rigidas"
        ],
        "free_description": "Necesita anticipación de cambios y entornos con baja estimulación sensorial."
      },
      "created_at": "2026-06-10T21:22:47Z"
    },
    {
      "id": 9001,
      "name": "Tomás Prueba",
      "classroom_id": 9001,
      "profile": {
        "id": 1,
        "student_id": 9001,
        "student_name": "Tomás Prueba",
        "is_transitory": false,
        "difficulties": [
          "se_distrae_facilmente",
          "impulsividad"
        ],
        "free_description": "Le cuesta sostener la atención en consignas largas; responde bien a pausas activas."
      },
      "created_at": "2026-06-10T21:22:47Z"
    }
  ]
}
```


---

## 03 - Teachers


### List Teachers

Lista los docentes de la organización.


**`GET http://127.0.0.1:8080/api/v1/teachers`** — expected `200`, got `200` ✅


_Response:_

```json
{
  "description": [
    {
      "id": 1,
      "name": "Test User",
      "email": "test@educabot.com",
      "role": "teacher"
    }
  ]
}
```


---

## 04 - Students


### List Students

Lista los alumnos de la organización.


**`GET http://127.0.0.1:8080/api/v1/students`** — expected `200`, got `200` ✅


_Response:_

```json
{
  "description": [
    {
      "id": 9002,
      "name": "Lucía Demo",
      "classroom_id": 9001,
      "profile": {
        "id": 2,
        "student_id": 9002,
        "student_name": "Lucía Demo",
        "is_transitory": true,
        "difficulties": [
          "dificultad_lectura"
        ],
        "free_description": "Dificultad transitoria de lectura tras cambio de escuela; en proceso de adaptación."
      },
      "created_at": "2026-06-10T21:22:47Z"
    },
    {
      "id": 9003,
      "name": "Mateo Test",
      "classroom_id": 9001,
      "profile": {
        "id": 3,
        "student_id": 9003,
        "student_name": "Mateo Test",
        "is_transitory": false,
        "difficulties": [
          "hipersensibilidad_sensorial",
          "rutinas_rigidas"
        ],
        "free_description": "Necesita anticipación de cambios y entornos con baja estimulación sensorial."
      },
      "created_at": "2026-06-10T21:22:47Z"
    },
    {
      "id": 9001,
      "name": "Tomás Prueba",
      "classroom_id": 9001,
      "profile": {
        "id": 1,
        "student_id": 9001,
        "student_name": "Tomás Prueba",
        "is_transitory": false,
        "difficulties": [
          "se_distrae_facilmente",
          "impulsividad"
        ],
        "free_description": "Le cuesta sostener la atención en consignas largas; responde bien a pausas activas."
      },
      "created_at": "2026-06-10T21:22:47Z"
    }
  ]
}
```


### List Students (by classroom)

Filtra alumnos por aula.


**`GET http://127.0.0.1:8080/api/v1/students?classroom_id=9001`** — expected `200`, got `200` ✅


_Response:_

```json
{
  "description": [
    {
      "id": 9002,
      "name": "Lucía Demo",
      "classroom_id": 9001,
      "profile": {
        "id": 2,
        "student_id": 9002,
        "student_name": "Lucía Demo",
        "is_transitory": true,
        "difficulties": [
          "dificultad_lectura"
        ],
        "free_description": "Dificultad transitoria de lectura tras cambio de escuela; en proceso de adaptación."
      },
      "created_at": "2026-06-10T21:22:47Z"
    },
    {
      "id": 9003,
      "name": "Mateo Test",
      "classroom_id": 9001,
      "profile": {
        "id": 3,
        "student_id": 9003,
        "student_name": "Mateo Test",
        "is_transitory": false,
        "difficulties": [
          "hipersensibilidad_sensorial",
          "rutinas_rigidas"
        ],
        "free_description": "Necesita anticipación de cambios y entornos con baja estimulación sensorial."
      },
      "created_at": "2026-06-10T21:22:47Z"
    },
    {
      "id": 9001,
      "name": "Tomás Prueba",
      "classroom_id": 9001,
      "profile": {
        "id": 1,
        "student_id": 9001,
        "student_name": "Tomás Prueba",
        "is_transitory": false,
        "difficulties": [
          "se_distrae_facilmente",
          "impulsividad"
        ],
        "free_description": "Le cuesta sostener la atención en consignas largas; responde bien a pausas activas."
      },
      "created_at": "2026-06-10T21:22:47Z"
    }
  ]
}
```


### Get Student (seed 9001)

Detalle de Tomás Prueba (perfil rico de prueba).


**`GET http://127.0.0.1:8080/api/v1/students/9001`** — expected `200`, got `200` ✅


_Response:_

```json
{
  "description": {
    "id": 9001,
    "name": "Tomás Prueba",
    "classroom_id": 9001,
    "profile": {
      "id": 1,
      "student_id": 9001,
      "student_name": "Tomás Prueba",
      "is_transitory": false,
      "difficulties": [
        "se_distrae_facilmente",
        "impulsividad"
      ],
      "free_description": "Le cuesta sostener la atención en consignas largas; responde bien a pausas activas."
    },
    "created_at": "2026-06-10T21:22:47Z"
  }
}
```


### Create Student

Crea un alumno en el aula de prueba.


**`POST http://127.0.0.1:8080/api/v1/students`** — expected `201`, got `201` ✅


_Request body:_

```json
{
  "name": "Alumno Newman",
  "classroom_id": 9001
}
```

_Response:_

```json
{
  "id": 1,
  "name": "Alumno Newman",
  "classroom_id": 9001,
  "created_at": "2026-06-10T19:24:42-03:00"
}
```


### Update Student (created)

Actualiza el alumno creado.


**`PUT http://127.0.0.1:8080/api/v1/students/1`** — expected `200`, got `200` ✅


_Request body:_

```json
{
  "name": "Alumno Newman (editado)"
}
```

_Response:_

```json
{
  "description": {
    "id": 1,
    "name": "Alumno Newman (editado)",
    "classroom_id": 9001,
    "created_at": "2026-06-10T19:24:42Z"
  }
}
```


---

## 05 - Student Profiles


### Get Student Profile (seed 9001)

Perfil del alumno (dificultades, descripción libre).


**`GET http://127.0.0.1:8080/api/v1/students/9001/profile`** — expected `200`, got `200` ✅


_Response:_

```json
{
  "description": {
    "id": 9001,
    "name": "Tomás Prueba",
    "classroom_id": 9001,
    "profile": {
      "id": 1,
      "student_id": 9001,
      "student_name": "Tomás Prueba",
      "is_transitory": false,
      "difficulties": [
        "se_distrae_facilmente",
        "impulsividad"
      ],
      "free_description": "Le cuesta sostener la atención en consignas largas; responde bien a pausas activas."
    },
    "created_at": "2026-06-10T21:22:47Z"
  }
}
```


### Upsert Student Profile (created)

Crea/actualiza el perfil del alumno.


**`PUT http://127.0.0.1:8080/api/v1/students/1/profile`** — expected `200`, got `200` ✅


_Request body:_

```json
{
  "is_transitory": false,
  "difficulties": [
    "motricidad_fina",
    "atencion"
  ],
  "free_description": "Prueba: dificultades de motricidad fina y atención."
}
```

_Response:_

```json
{
  "description": {
    "id": 1,
    "name": "Alumno Newman (editado)",
    "classroom_id": 9001,
    "profile": {
      "id": 4,
      "student_id": 1,
      "student_name": "Alumno Newman (editado)",
      "is_transitory": false,
      "difficulties": [
        "motricidad_fina",
        "atencion"
      ],
      "free_description": "Prueba: dificultades de motricidad fina y atención."
    },
    "created_at": "2026-06-10T19:24:42Z"
  }
}
```


---

## 06 - Catalog


### List Ramps

Lista las rampas (agrupadores de dispositivos de la valija).


**`GET http://127.0.0.1:8080/api/v1/ramps`** — expected `200`, got `200` ✅


_Response:_

```json
{
  "description": [
    {
      "id": 9001,
      "name": "Rampa de Atención",
      "description": "Apoyos para atención y autorregulación",
      "short_description": "Atención",
      "sort_order": 1,
      "devices": [
        {
          "id": 9001,
          "ramp_id": 9001,
          "name": "Tablero de anticipación",
          "description": "Secuencia visual de la rutina de clase",
          "qr_code": "QR-TEST-9001",
          "useful_when": "Cuando el alumno no inicia la tarea o necesita anticipar pasos",
          "quantity": 1,
          "sort_order": 1
        },
        {
          "id": 9002,
          "ramp_id": 9001,
          "name": "Auriculares de aislamiento",
          "description": "Reducen estímulos sonoros para sostener la atención",
          "qr_code": "QR-TEST-9002",
          "useful_when": "Cuando el alumno se distrae o se desregula por ruido",
          "quantity": 1,
          "sort_order": 2
        }
      ]
    }
  ]
}
```


### Get Ramp (9001)

Detalle de una rampa.


**`GET http://127.0.0.1:8080/api/v1/ramps/9001`** — expected `200`, got `200` ✅


_Response:_

```json
{
  "description": {
    "id": 9001,
    "name": "Rampa de Atención",
    "description": "Apoyos para atención y autorregulación",
    "short_description": "Atención",
    "sort_order": 1,
    "devices": [
      {
        "id": 9001,
        "ramp_id": 9001,
        "name": "Tablero de anticipación",
        "description": "Secuencia visual de la rutina de clase",
        "qr_code": "QR-TEST-9001",
        "useful_when": "Cuando el alumno no inicia la tarea o necesita anticipar pasos",
        "quantity": 1,
        "sort_order": 1
      },
      {
        "id": 9002,
        "ramp_id": 9001,
        "name": "Auriculares de aislamiento",
        "description": "Reducen estímulos sonoros para sostener la atención",
        "qr_code": "QR-TEST-9002",
        "useful_when": "Cuando el alumno se distrae o se desregula por ruido",
        "quantity": 1,
        "sort_order": 2
      }
    ]
  }
}
```


### List Devices

Lista los dispositivos de la valija adaptativa.


**`GET http://127.0.0.1:8080/api/v1/devices`** — expected `200`, got `200` ✅


_Response:_

```json
{
  "description": [
    {
      "id": 9001,
      "ramp_id": 9001,
      "name": "Tablero de anticipación",
      "description": "Secuencia visual de la rutina de clase",
      "qr_code": "QR-TEST-9001",
      "useful_when": "Cuando el alumno no inicia la tarea o necesita anticipar pasos",
      "quantity": 1,
      "sort_order": 1,
      "ramp_name": "Rampa de Atención"
    },
    {
      "id": 9002,
      "ramp_id": 9001,
      "name": "Auriculares de aislamiento",
      "description": "Reducen estímulos sonoros para sostener la atención",
      "qr_code": "QR-TEST-9002",
      "useful_when": "Cuando el alumno se distrae o se desregula por ruido",
      "quantity": 1,
      "sort_order": 2,
      "ramp_name": "Rampa de Atención"
    }
  ]
}
```


### List Devices (by ramp)

Filtra dispositivos por rampa.


**`GET http://127.0.0.1:8080/api/v1/devices?ramp_id=9001`** — expected `200`, got `200` ✅


_Response:_

```json
{
  "description": [
    {
      "id": 9001,
      "ramp_id": 9001,
      "name": "Tablero de anticipación",
      "description": "Secuencia visual de la rutina de clase",
      "qr_code": "QR-TEST-9001",
      "useful_when": "Cuando el alumno no inicia la tarea o necesita anticipar pasos",
      "quantity": 1,
      "sort_order": 1,
      "ramp_name": "Rampa de Atención"
    },
    {
      "id": 9002,
      "ramp_id": 9001,
      "name": "Auriculares de aislamiento",
      "description": "Reducen estímulos sonoros para sostener la atención",
      "qr_code": "QR-TEST-9002",
      "useful_when": "Cuando el alumno se distrae o se desregula por ruido",
      "quantity": 1,
      "sort_order": 2,
      "ramp_name": "Rampa de Atención"
    }
  ]
}
```


### Get Device (9001)

Detalle de un dispositivo.


**`GET http://127.0.0.1:8080/api/v1/devices/9001`** — expected `200`, got `200` ✅


_Response:_

```json
{
  "description": {
    "id": 9001,
    "ramp_id": 9001,
    "name": "Tablero de anticipación",
    "description": "Secuencia visual de la rutina de clase",
    "qr_code": "QR-TEST-9001",
    "useful_when": "Cuando el alumno no inicia la tarea o necesita anticipar pasos",
    "quantity": 1,
    "sort_order": 1,
    "ramp_name": "Rampa de Atención"
  }
}
```


---

## 07 - Adaptations


### List Adaptations

Lista las adaptaciones (recursos) del docente.


**`GET http://127.0.0.1:8080/api/v1/adaptations`** — expected `200`, got `200` ✅


_Response:_

```json
{
  "description": []
}
```


### Create Adaptation

Crea una adaptación vinculada a un alumno y un dispositivo.


**`POST http://127.0.0.1:8080/api/v1/adaptations`** — expected `201`, got `201` ✅


_Request body:_

```json
{
  "student_id": 9001,
  "device_id": 9001,
  "subject": "Matemática",
  "title": "Fracciones con material concreto",
  "activity_description": "Fracciones con tablero de anticipación",
  "adaptation_strategy": "Anticipar pasos con apoyo visual",
  "adaptation_type": "actividad_adaptada",
  "notes": "Responde bien a estímulos visuales"
}
```

_Response:_

```json
{
  "id": 1,
  "student_id": 9001,
  "student_name": "Tomás Prueba",
  "teacher_id": 1,
  "teacher_name": "Test User",
  "device_id": 9001,
  "device_name": "Tablero de anticipación",
  "device_ids": [],
  "device_names": [],
  "title": "Fracciones con material concreto",
  "subject": "Matemática",
  "activity_description": "Fracciones con tablero de anticipación",
  "adaptation_strategy": "Anticipar pasos con apoyo visual",
  "adaptation_type": "actividad_adaptada",
  "notes": "Responde bien a estímulos visuales",
  "status": "en_curso",
  "was_edited": false,
  "created_at": "2026-06-10T19:24:42Z",
  "updated_at": "2026-06-10T19:24:42Z"
}
```


### List Adaptations (by student)

Filtra adaptaciones por alumno.


**`GET http://127.0.0.1:8080/api/v1/adaptations?student_id=9001`** — expected `200`, got `200` ✅


_Response:_

```json
{
  "description": [
    {
      "id": 1,
      "student_id": 9001,
      "student_name": "Tomás Prueba",
      "teacher_id": 1,
      "teacher_name": "Test User",
      "device_id": 9001,
      "device_name": "Tablero de anticipación",
      "device_ids": [],
      "device_names": [],
      "title": "Fracciones con material concreto",
      "subject": "Matemática",
      "activity_description": "Fracciones con tablero de anticipación",
      "adaptation_strategy": "Anticipar pasos con apoyo visual",
      "adaptation_type": "actividad_adaptada",
      "notes": "Responde bien a estímulos visuales",
      "status": "en_curso",
      "was_edited": false,
      "created_at": "2026-06-10T19:24:42Z",
      "updated_at": "2026-06-10T19:24:42Z"
    }
  ]
}
```


### Get Adaptation (created)

Detalle de la adaptación creada.


**`GET http://127.0.0.1:8080/api/v1/adaptations/1`** — expected `200`, got `200` ✅


_Response:_

```json
{
  "description": {
    "id": 1,
    "student_id": 9001,
    "student_name": "Tomás Prueba",
    "teacher_id": 1,
    "teacher_name": "Test User",
    "device_id": 9001,
    "device_name": "Tablero de anticipación",
    "device_ids": [],
    "device_names": [],
    "title": "Fracciones con material concreto",
    "subject": "Matemática",
    "activity_description": "Fracciones con tablero de anticipación",
    "adaptation_strategy": "Anticipar pasos con apoyo visual",
    "adaptation_type": "actividad_adaptada",
    "notes": "Responde bien a estímulos visuales",
    "status": "en_curso",
    "was_edited": false,
    "created_at": "2026-06-10T19:24:42Z",
    "updated_at": "2026-06-10T19:24:42Z"
  }
}
```


### Update Adaptation (created)

Actualiza estado/resultado de la adaptación.


**`PUT http://127.0.0.1:8080/api/v1/adaptations/1`** — expected `200`, got `200` ✅


_Request body:_

```json
{
  "status": "probado",
  "outcome": "Completó la actividad con más autonomía",
  "notes": "Funcionó, repetir la próxima clase"
}
```

_Response:_

```json
{
  "description": {
    "id": 1,
    "student_id": 9001,
    "student_name": "Tomás Prueba",
    "teacher_id": 1,
    "teacher_name": "Test User",
    "device_id": 9001,
    "device_name": "Tablero de anticipación",
    "device_ids": [],
    "device_names": [],
    "title": "Fracciones con material concreto",
    "subject": "Matemática",
    "activity_description": "Fracciones con tablero de anticipación",
    "adaptation_strategy": "Anticipar pasos con apoyo visual",
    "adaptation_type": "actividad_adaptada",
    "outcome": "Completó la actividad con más autonomía",
    "notes": "Funcionó, repetir la próxima clase",
    "status": "probado",
    "was_edited": false,
    "created_at": "2026-06-10T19:24:42Z",
    "updated_at": "2026-06-10T19:24:42Z"
  }
}
```


### List Adaptation Resources (created)

Recursos (materiales) ligados a la adaptación.


**`GET http://127.0.0.1:8080/api/v1/adaptations/1/resources`** — expected `200`, got `200` ✅


_Response:_

```json
{
  "description": []
}
```


### Export Adaptation (md)

Exporta la adaptación a Markdown (descarga). Content-Type text/markdown.


**`GET http://127.0.0.1:8080/api/v1/adaptations/1/export?format=md`** — expected `200`, got `200` ✅


_Response:_

```
Content-Type: text/markdown; charset=utf-8
# Matemática

**Alumno:** Tomás Prueba  
**Docente:** Test User  
**Tipo:** actividad_adaptada  
**Estado:** Probado

## Actividad

Fracciones con tablero de anticipación

## Estrategia

Anticipar pasos con apoyo visual

## Dispositivos sugeridos

### Tablero de anticipación


## Notas para el docente

Funcionó, repetir la próxima clase

## Resultado

Completó la actividad con más autonomía

---

_Generado por Alizia · Educabot · 10/06/2026 · Adaptación #1_

```


---

## 08 - Chat History


### Chat History (assist)

Historial de conversaciones del docente en el modo indicado (contextId).


**`GET http://127.0.0.1:8080/api/v1/chat/history/assist`** — expected `200`, got `200` ✅


_Response:_

```json
{
  "description": [
    {
      "id": 1,
      "mode": "assist",
      "messages": [
        {
          "role": "user",
          "content": "Hola",
          "created_at": "2026-06-10T18:39:08Z"
        },
        {
          "role": "assistant",
          "content": "[stub] chat response",
          "created_at": "2026-06-10T18:39:08Z"
        }
      ],
      "created_at": "2026-06-10T18:39:08Z"
    }
  ]
}
```


---

## 09 - Dashboard


### Get Metrics

Métricas agregadas (alumnos, adaptaciones, aulas).


**`GET http://127.0.0.1:8080/api/v1/dashboard/metrics`** — expected `200`, got `200` ✅


_Response:_

```json
{
  "description": {
    "total_students": 4,
    "students_with_profiles": 4,
    "total_adaptations": 1,
    "adaptations_by_status": {
      "probado": 1
    },
    "adaptations_by_type": {
      "actividad_adaptada": 1
    },
    "top_used_devices": [],
    "adaptations_this_week": 1,
    "classroom_count": 2
  }
}
```


### Get AI Usage

Consumo de tokens de IA por la organización.


**`GET http://127.0.0.1:8080/api/v1/dashboard/ai-usage`** — expected `200`, got `200` ✅


_Response:_

```json
{
  "description": {
    "window_days": 30,
    "total_requests": 0,
    "prompt_tokens": 0,
    "completion_tokens": 0,
    "total_tokens": 0,
    "by_mode": []
  }
}
```


---

## 10 - Context Engine · HU-1 Apertura


### Open · sin dimensión (saludo)

Sin dimensión: saluda y pregunta de qué hablar (needs_dimension=true).


**`POST http://127.0.0.1:8080/api/v1/inclusion/open`** — expected `200`, got `200` ✅


_Request body:_

```json
{}
```

_Response:_

```json
{
  "description": {
    "greeting": "¡Hola! Soy Alizia, tu asistente de inclusión. ¿De qué querés hablar: de un alumno, de la valija o de un tema?",
    "needs_dimension": true
  }
}
```


### Open · dimensión ambigua

Dimensión no reconocida: repregunta.


**`POST http://127.0.0.1:8080/api/v1/inclusion/open`** — expected `200`, got `200` ✅


_Request body:_

```json
{
  "dimension": "loquesea"
}
```

_Response:_

```json
{
  "description": {
    "greeting": "No me quedó claro. ¿Querés hablar de un alumno, de la valija o de un tema?",
    "needs_dimension": true
  }
}
```


### Open · alumno sin id

Falta student_id: repregunta por el alumno.


**`POST http://127.0.0.1:8080/api/v1/inclusion/open`** — expected `200`, got `200` ✅


_Request body:_

```json
{
  "dimension": "alumno"
}
```

_Response:_

```json
{
  "description": {
    "greeting": "Dale, ¿de qué alumno querés hablar?",
    "needs_dimension": true
  }
}
```


### Open · alumno 9001

Abre sobre Tomás (9001) y recupera resúmenes previos del alumno.


**`POST http://127.0.0.1:8080/api/v1/inclusion/open`** — expected `200`, got `200` ✅


_Request body:_

```json
{
  "dimension": "alumno",
  "student_id": 9001
}
```

_Response:_

```json
{
  "description": {
    "greeting": "Listo, hablemos de Tomás Prueba.",
    "needs_dimension": false,
    "dimension": "alumno",
    "student": {
      "id": 9001,
      "organization_id": "00000000-0000-0000-0000-000000000001",
      "classroom_id": 9001,
      "name": "Tomás Prueba",
      "age_range": "8-9",
      "grade_level": "4to",
      "preferred_name": "Tomi",
      "profile": {
        "id": 1,
        "student_id": 9001,
        "is_transitory": false,
        "difficulties": [
          "se_distrae_facilmente",
          "impulsividad"
        ],
        "free_description": "Le cuesta sostener la atención en consignas largas; responde bien a pausas activas.",
        "support_level": "medio",
        "strengths": [
          "memoria visual",
          "creatividad"
        ],
        "interests": [
          "dinosaurios",
          "dibujo"
        ],
        "triggers": [
          "ruidos fuertes",
          "consignas extensas"
        ],
        "effective_strategies": [
          "pausas activas",
          "consignas cortas paso a paso"
        ],
        "ineffective_strategies": [
          "retos en público"
        ],
        "situation_codes": [
          "no_inicia_tarea",
          "se_distrae"
        ],
        "has_therapeutic_companion": true,
        "environment_notes": "Acompañante terapéutico 3 veces por semana; familia muy presente.",
        "created_at": "2026-06-10T21:22:47.30357Z",
        "updated_at": "2026-06-10T21:22:47.30357Z"
  
... (truncado)
```


### Open · alumno inexistente

Alumno inexistente: error.


**`POST http://127.0.0.1:8080/api/v1/inclusion/open`** — expected `404`, got `404` ✅


_Request body:_

```json
{
  "dimension": "alumno",
  "student_id": 999999
}
```

_Response:_

```json
{
  "code": "not_found",
  "description": "not found"
}
```


### Open · tema sin topic

Falta el tema: repregunta.


**`POST http://127.0.0.1:8080/api/v1/inclusion/open`** — expected `200`, got `200` ✅


_Request body:_

```json
{
  "dimension": "tema"
}
```

_Response:_

```json
{
  "description": {
    "greeting": "Perfecto, ¿sobre qué tema querés que busquemos?",
    "needs_dimension": true
  }
}
```


### Open · tema TEA

Abre sobre el tema TEA y recupera resúmenes por keyword.


**`POST http://127.0.0.1:8080/api/v1/inclusion/open`** — expected `200`, got `200` ✅


_Request body:_

```json
{
  "dimension": "tema",
  "topic": "TEA"
}
```

_Response:_

```json
{
  "description": {
    "greeting": "Buenísimo, trabajemos sobre TEA.",
    "needs_dimension": false,
    "dimension": "tema"
  }
}
```


### Open · valija

Abre sobre la valija (sin alumno).


**`POST http://127.0.0.1:8080/api/v1/inclusion/open`** — expected `200`, got `200` ✅


_Request body:_

```json
{
  "dimension": "valija"
}
```

_Response:_

```json
{
  "description": {
    "greeting": "Genial, miremos la valija. ¿Qué necesitás resolver?",
    "needs_dimension": false,
    "dimension": "valija"
  }
}
```


---

## 11 - Context Engine · HU-2 Contexto


### Context · alumno 9001 (perfil rico)

Arma el contexto del alumno: perfil, PPI, diagnósticos, situaciones.


**`POST http://127.0.0.1:8080/api/v1/inclusion/context`** — expected `200`, got `200` ✅


_Request body:_

```json
{
  "dimension": "alumno",
  "student_id": 9001
}
```

_Response:_

```json
{
  "description": {
    "device_catalog": [
      {
        "id": 9001,
        "organization_id": "00000000-0000-0000-0000-000000000001",
        "ramp_id": 9001,
        "name": "Tablero de anticipación",
        "description": "Secuencia visual de la rutina de clase",
        "qr_code": "QR-TEST-9001",
        "useful_when": "Cuando el alumno no inicia la tarea o necesita anticipar pasos",
        "quantity": 1,
        "sort_order": 1,
        "ramp": {
          "id": 9001,
          "organization_id": "00000000-0000-0000-0000-000000000001",
          "name": "Rampa de Atención",
          "description": "Apoyos para atención y autorregulación",
          "short_description": "Atención",
          "sort_order": 1,
          "created_at": "2026-06-10T21:40:42.968007Z",
          "updated_at": "2026-06-10T21:40:42.968007Z"
        },
        "created_at": "2026-06-10T21:40:42.968007Z",
        "updated_at": "2026-06-10T21:40:42.968007Z"
      },
      {
        "id": 9002,
        "organization_id": "00000000-0000-0000-0000-000000000001",
        "ramp_id": 9001,
        "name": "Auriculares de aislamiento",
        "description": "Reducen estímulos sonoros para sostener la atención",
        "qr_code": "QR-TEST-9002",
        "useful_when": "Cuando el alumno se distrae o se desregula por ruido",
        "quantity": 1,
        "sort_order": 2,
        "ramp": {
          "id": 9001,
          "organization_id": "00000000-0000-0000-0000-000000000001",
          "name": "Ra
... (truncado)
```


### Context · alumno 9003 (degradación)

Alumno sin PPI: degrada con elegancia (missing_data incluye 'ppi').


**`POST http://127.0.0.1:8080/api/v1/inclusion/context`** — expected `200`, got `200` ✅


_Request body:_

```json
{
  "dimension": "alumno",
  "student_id": 9003
}
```

_Response:_

```json
{
  "description": {
    "device_catalog": [
      {
        "id": 9001,
        "organization_id": "00000000-0000-0000-0000-000000000001",
        "ramp_id": 9001,
        "name": "Tablero de anticipación",
        "description": "Secuencia visual de la rutina de clase",
        "qr_code": "QR-TEST-9001",
        "useful_when": "Cuando el alumno no inicia la tarea o necesita anticipar pasos",
        "quantity": 1,
        "sort_order": 1,
        "ramp": {
          "id": 9001,
          "organization_id": "00000000-0000-0000-0000-000000000001",
          "name": "Rampa de Atención",
          "description": "Apoyos para atención y autorregulación",
          "short_description": "Atención",
          "sort_order": 1,
          "created_at": "2026-06-10T21:40:42.968007Z",
          "updated_at": "2026-06-10T21:40:42.968007Z"
        },
        "created_at": "2026-06-10T21:40:42.968007Z",
        "updated_at": "2026-06-10T21:40:42.968007Z"
      },
      {
        "id": 9002,
        "organization_id": "00000000-0000-0000-0000-000000000001",
        "ramp_id": 9001,
        "name": "Auriculares de aislamiento",
        "description": "Reducen estímulos sonoros para sostener la atención",
        "qr_code": "QR-TEST-9002",
        "useful_when": "Cuando el alumno se distrae o se desregula por ruido",
        "quantity": 1,
        "sort_order": 2,
        "ramp": {
          "id": 9001,
          "organization_id": "00000000-0000-0000-0000-000000000001",
          "name": "Ra
... (truncado)
```


### Context · valija (lazy)

Contexto de valija: catálogo + situaciones, sin target_student.


**`POST http://127.0.0.1:8080/api/v1/inclusion/context`** — expected `200`, got `200` ✅


_Request body:_

```json
{
  "dimension": "valija"
}
```

_Response:_

```json
{
  "description": {
    "device_catalog": [
      {
        "id": 9001,
        "organization_id": "00000000-0000-0000-0000-000000000001",
        "ramp_id": 9001,
        "name": "Tablero de anticipación",
        "description": "Secuencia visual de la rutina de clase",
        "qr_code": "QR-TEST-9001",
        "useful_when": "Cuando el alumno no inicia la tarea o necesita anticipar pasos",
        "quantity": 1,
        "sort_order": 1,
        "ramp": {
          "id": 9001,
          "organization_id": "00000000-0000-0000-0000-000000000001",
          "name": "Rampa de Atención",
          "description": "Apoyos para atención y autorregulación",
          "short_description": "Atención",
          "sort_order": 1,
          "created_at": "2026-06-10T21:40:42.968007Z",
          "updated_at": "2026-06-10T21:40:42.968007Z"
        },
        "created_at": "2026-06-10T21:40:42.968007Z",
        "updated_at": "2026-06-10T21:40:42.968007Z"
      },
      {
        "id": 9002,
        "organization_id": "00000000-0000-0000-0000-000000000001",
        "ramp_id": 9001,
        "name": "Auriculares de aislamiento",
        "description": "Reducen estímulos sonoros para sostener la atención",
        "qr_code": "QR-TEST-9002",
        "useful_when": "Cuando el alumno se distrae o se desregula por ruido",
        "quantity": 1,
        "sort_order": 2,
        "ramp": {
          "id": 9001,
          "organization_id": "00000000-0000-0000-0000-000000000001",
          "name": "Ra
... (truncado)
```


### Context · alumno inexistente

Alumno inexistente: error.


**`POST http://127.0.0.1:8080/api/v1/inclusion/context`** — expected `404`, got `404` ✅


_Request body:_

```json
{
  "dimension": "alumno",
  "student_id": 999999
}
```

_Response:_

```json
{
  "code": "not_found",
  "description": "not found"
}
```


---

## 12 - Context Engine · HU-3 Contenido


### Search · 'autismo desregula' (match TEA)

Busca material pedagógico por keyword/full-text. Devuelve el doc TEA.


**`POST http://127.0.0.1:8080/api/v1/inclusion/search-content`** — expected `200`, got `200` ✅


_Request body:_

```json
{
  "query": "autismo desregula"
}
```

_Response:_

```json
{
  "description": {
    "query": "autismo desregula",
    "results": [
      {
        "content_id": 2,
        "chunk_id": 2,
        "title": "Estrategias para acompañar a estudiantes con TEA en el aula",
        "type": "material",
        "keywords": [
          "TEA",
          "autismo",
          "autorregulacion",
          "anticipacion",
          "transiciones"
        ],
        "preview": "Anticipar los cambios con apoyos visuales, ofrecer un rincón de calma y materiales para regular estímulos (auriculares, time timer). Dividir las consignas en pasos cortos y mantener rutinas predecibles favorece la autorregulación.",
        "score": 0.33435988426208496
      }
    ]
  }
}
```


### Search · 'lectura dislexia'

Devuelve el documento de dislexia.


**`POST http://127.0.0.1:8080/api/v1/inclusion/search-content`** — expected `200`, got `200` ✅


_Request body:_

```json
{
  "query": "lectura dislexia"
}
```

_Response:_

```json
{
  "description": {
    "query": "lectura dislexia",
    "results": [
      {
        "content_id": 3,
        "chunk_id": 3,
        "title": "Técnicas de lectura para estudiantes con dislexia",
        "type": "material",
        "keywords": [
          "dislexia",
          "lectura",
          "comprension",
          "accesibilidad"
        ],
        "preview": "Presentar textos en fragmentos cortos, con tipografía accesible e interlineado amplio. Combinar lectura con audio y evaluar la comprensión de forma oral, separada de la decodificación.",
        "score": 0.790727436542511
      }
    ]
  }
}
```


### Search · 'quimica organica' (sin match)

Sin coincidencias: no inventa, devuelve results vacío.


**`POST http://127.0.0.1:8080/api/v1/inclusion/search-content`** — expected `200`, got `200` ✅


_Request body:_

```json
{
  "query": "quimica organica"
}
```

_Response:_

```json
{
  "description": {
    "query": "quimica organica",
    "results": []
  }
}
```


### Search · query vacía

Query vacía: results vacío.


**`POST http://127.0.0.1:8080/api/v1/inclusion/search-content`** — expected `200`, got `200` ✅


_Request body:_

```json
{
  "query": ""
}
```

_Response:_

```json
{
  "description": {
    "query": "",
    "results": []
  }
}
```


---

## 13 - Context Engine · IA (recommend / assist)


### Recommend Device

Recomienda dispositivos de la valija para una actividad (LLM).


**`POST http://127.0.0.1:8080/api/v1/inclusion/recommend`** — expected `200`, got `200` ✅


_Request body:_

```json
{
  "student_id": 9001,
  "subject": "Lengua",
  "objective": "Producción de textos narrativos",
  "duration": "40",
  "dynamic": "Trabajo individual con apoyo docente",
  "materials": "Cuaderno, lápiz",
  "history": []
}
```

_Response:_

```json
{
  "description": {
    "response": "[stub] chat response",
    "conversation_id": 2
  }
}
```


### Assist Classroom

Chat asistente (loop agéntico con tools). Persiste la conversación.


**`POST http://127.0.0.1:8080/api/v1/inclusion/assist`** — expected `200`, got `200` ✅


_Request body:_

```json
{
  "classroom_id": 9001,
  "student_id": 9001,
  "message": "Tomás no logra sostener la atención al escribir. ¿Qué hago?",
  "history": []
}
```

_Response:_

```json
{
  "description": {
    "response": "[stub] chat response",
    "conversation_id": 3
  }
}
```


---

## 14 - Context Engine · HU-5 Memoria (close)


### Close Session (compacta conversación)

Compacta la conversación al cerrar: genera el resumen + tags (alumno/tema/valija).


**`POST http://127.0.0.1:8080/api/v1/inclusion/close`** — expected `200`, got `200` ✅


_Request body:_

```json
{
  "conversation_id": 3
}
```

_Response:_

```json
{
  "description": {
    "conversation_id": 3,
    "summary": "[stub] chat response",
    "topic_keywords": null,
    "student_ids": [
      9001
    ]
  }
}
```


### Close Session · conversación inexistente

Conversación inexistente: error.


**`POST http://127.0.0.1:8080/api/v1/inclusion/close`** — expected `404`, got `404` ✅


_Request body:_

```json
{
  "conversation_id": 999999
}
```

_Response:_

```json
{
  "code": "not_found",
  "description": "not found"
}
```


### Close Session · sin conversation_id

Falta conversation_id: validación (400).


**`POST http://127.0.0.1:8080/api/v1/inclusion/close`** — expected `400`, got `400` ✅


_Request body:_

```json
{}
```

_Response:_

```json
{
  "code": "validation_error",
  "description": "validation error: conversation_id is required"
}
```


---

## 15 - Cleanup


### Delete Adaptation (created)

Borra la adaptación creada.


**`DELETE http://127.0.0.1:8080/api/v1/adaptations/1`** — expected `204`, got `204` ✅


_Response:_

_(vacío)_


### Delete Student (created)

Borra el alumno creado.


**`DELETE http://127.0.0.1:8080/api/v1/students/1`** — expected `204`, got `204` ✅


_Response:_

_(vacío)_


### Delete Classroom (created)

Borra el aula creada.


**`DELETE http://127.0.0.1:8080/api/v1/classrooms/1`** — expected `204`, got `204` ✅


_Response:_

_(vacío)_
