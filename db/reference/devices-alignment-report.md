# Alineación de devices con la Valija Chubut

**Fuente de la verdad:** pestaña **"Listado de productos - Chubut"** del sheet "Listado de
productos - Valija chubut" (owner anabela@educabot.com,
[sheet](https://docs.google.com/spreadsheets/d/1P62Omhma1oYrHvSJkl3zYPlqq-pisZAYoXDWTZ0IjOU/edit?gid=296161033)).
Solo esa pestaña es canónica. Datos crudos en `db/reference/valija-chubut.csv`
(24 productos + 2 kits: `ETE-I10822-EB`, `ETE-I10857-EB`).

## Resultado: 24 devices activos = exactamente el sheet

Granularidad **C**: 1 device por producto, con `product_code` (ETE), `product_family` y `stage`
para agrupar/filtrar etapas. Categoría = `ramp` (las 5 ya existentes).

| Acción | Cantidad | Detalle |
|---|---|---|
| Backfill `product_code` (siguen activos) | 6 | Auriculares, Auriculares cancelación, Mouse trackball, Pulsador USB, Pen reader, Elástico→Bouncyband (ETE-I10818) |
| Devices nuevos (upsert) | 18 | Resto del sheet, con código + ramp + family/stage |
| Desactivados (`is_active=FALSE`, NO borrados) | 9 | 6 sin respaldo en el sheet + 3 genéricos superados por sus etapas |

**Desactivados:** Time Timer, Tablet educativa 10", Finger focus, Teclado CLEVY, Soporte
flexible, Pelota antiestrés *(sin respaldo en el sheet)* · Regla de lectura con ventana, Pinzas
de escritura, Tijeras adaptadas *(genéricos reemplazados por sus versiones por etapa)*.

## Familias / etapas
- `soporte_lapiz`: ETE-I10824 (etapa 4), 10825 (3), 10826 (2), 10827 (1)
- `ayuda_lectura`: ETE-I10829 (1), 10830 (3), 10831 (4)
- `tijera_adaptada`: ETE-I10836 (1), 10837 (2) · `tijera`: ETE-I10835 (zurdos)

## Cambios en el schema (migración 000023)
`devices` += `product_code VARCHAR(50)`, `product_family VARCHAR(80)`, `stage SMALLINT`,
`is_active BOOLEAN NOT NULL DEFAULT TRUE` + índice único parcial
`uq_devices_org_product_code (organization_id, product_code) WHERE product_code IS NOT NULL`.

## Notas de ejecución (prod)
- **Org de producción real:** `Alizia Inclusión` = `b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22`
  (la org demo `Escuela Demo Inclusión` del seed NO existe en prod). El seed commiteado apunta a
  la org demo (consistente con `db/seeds/seed.sql`); para correr en prod se usa una variante que
  resuelve la org por `name = 'Alizia Inclusión'`. Prod ya tiene los 5 ramps y los 15 devices con
  nombres idénticos → el backfill/desactivación por nombre aplica sin cambios.
- Correr **después** del seed base (los INSERT resuelven `ramp_id` por nombre).

## Pendiente (follow-up, fuera de este cambio)
- Exponer los nuevos campos (`product_code`, `is_active`, `product_family`, `stage`) en la entidad
  Go / API y filtrar por `is_active` en los listados del catálogo.
