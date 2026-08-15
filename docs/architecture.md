# Architecture & design notes

This is a write-up of how the app is put together and why I made the choices I did. It's not a spec, just notes for the reasoning behind the code.

## The stack, and why

Backend is Go with Chi. I didn't want a full framework like Gin, Chi is basically a router on top of the standard `net/http`, so handlers are just regular `http.HandlerFunc`s and there's nothing framework-specific to learn to read the code. For a project this size that felt like the right amount of structure, enough for routing and middleware (request IDs, logging, recover, CORS), not much more.

Frontend is React + TypeScript with Vite, no state management library. The whole app is small enough that `App.tsx` holding a bit of state (which tags are selected, a refresh counter, whether the upload modal is open) and passing it down as props is genuinely simpler than wiring up Redux or Zustand for three pieces of state.

Postgres for the actual data, Redis for caching. Images are stored on disk in a Docker volume.

## Data model

Three tables, `images`, `tags`, and a join table `image_tags`. Tags are `citext` so `"Cats"` and `"cats"` collide into the same row instead of creating duplicates, that one bit me early on when I was testing and ended up with three different casings of the same tag cluttering the filter bar.

`images` has a `status` column that's always `'ready'` right now, uploads are processed synchronously in the request handler. I added the column anyway because if uploads ever needed to move to a background job (say, for slower processing or some scanning), I didn't want to touch the schema again. It's unused optionality, but it shows that its future proof. It can also be useful in the queue mechanism requested as bonus point on the task description, but it seemed too big of an architecture for such a small task.

There are two indexes beyond the primary keys. `idx_images_created_at` on `images(created_at DESC, id DESC)` matches the feed query's `ORDER BY` exactly, so pagination doesn't need a full table sort. `idx_image_tags_tag_id` on `image_tags(tag_id)` backs the tag filter, without it, resolving a tag name to its matching images would scan the whole join table, since its primary key `(image_id, tag_id)` only helps once you already know the `image_id`.

## Pagination

The image feed uses cursor pagination, `(created_at, id)` as a tuple, instead of `OFFSET`/`LIMIT` page numbers. With offset pagination, deep pages get slower as the offset grows because Postgres still has to scan and discard all the earlier rows. Keyset pagination avoids that, and it also plays nicer with a live feed, if a new image lands while you're scrolled down, offset-based paging would shift and could show you a duplicate or skip a row. The cursor is just `created_at + id` base64-encoded, `id` is the tiebreaker for rows with an identical timestamp.

The same cursor is what powers infinite scroll on the frontend too, for a continuous feed closer to Instagram. An `IntersectionObserver` watches a sentinel element at the bottom of the grid and fetches the next page automatically once it scrolls into view, no backend changes needed, it just reuses the existing cursor pagination.

## Image normalization

Every upload gets decoded, resized to a fixed 1080x1080, and re-encoded as JPEG at quality 85, regardless of the original format (jpeg/png/webp are all accepted in). This keeps the feed grid uniform and keeps storage and bandwidth predictable, no matter what someone uploads. Uploads are also capped at 5MB at the handler level, before normalization even runs, so nothing oversized gets that far in the first place.

The honest caveat is it resizes to fill the square without preserving aspect ratio, so a wide landscape photo gets visibly squashed. A center-crop would look better and it's on the list below, I just wanted normalization working end to end before polishing how it crops.

## Caching

`GET /api/images` and `GET /api/tags` are cached in Redis for 30 seconds, keyed by the tag filter and cursor. On upload, the handler deletes the cache entry for the untagged first page and the tag list, so a fresh upload shows up immediately for someone browsing the default feed. Cached entries for other filter/cursor combinations aren't explicitly busted, they just fall out of the cache naturally within 30 seconds. That's a deliberate shortcut, with a 30 second TTL, the staleness window is short enough that I didn't think pattern-based invalidation (e.g. scanning and deleting all `images:*` keys) was worth the added complexity for this app's scale, but it's a real trade-off, not an oversight.

## Live updates

New uploads get broadcast over a WebSocket (`/ws`) to every connected client. The hub itself is a small in-memory struct, a map of connected clients guarded by a mutex. When a client gets an `image.created` message it just triggers a full refetch of the feed rather than trying to splice the new image into local state. That's simpler to reason about and avoids the client and server ever disagreeing about ordering or duplicate entries, at the cost of an extra round trip per upload. Given uploads aren't a high-frequency event, that trade felt fine.

The one real limitation here is that the hub is in-process. If the backend ever ran as more than one instance, clients connected to instance A wouldn't hear about an upload that hit instance B. Fixing that would mean fanning broadcasts out through Redis pub/sub instead of holding client connections in a local map.

## API docs

`/docs` serves a Swagger UI page backed by a hand-written OpenAPI spec (`backend/internal/api/docs/openapi.yaml`). I didn't reach for a codegen tool like `swaggo/swag`, which builds the spec from comment annotations sprinkled across every handler, because the API surface here is tiny (three REST endpoints), and hand-writing one YAML file was less overhead than adding annotations everywhere and a build step to regenerate them. Both the spec and the UI page are embedded into the Go binary with `go:embed` so they ship inside the Docker image without any extra file copying.

## Local dev vs Docker

Everything runs through `docker compose up`, Postgres, Redis, a one-shot `migrate` container that runs migrations and exits, then the backend and frontend. The `migrate` service uses `depends_on: condition: service_completed_successfully` so the backend doesn't start against a database that hasn't been migrated yet. Running things locally without Docker is also supported (see the README) for faster iteration on the Go code specifically, since rebuilding a Docker image on every change is slower than `go run`.
