# Possible improvements

This was a small, scoped task, so I kept the implementation matched to that. This list is meant to show how I'd think about it if it had to grow into something bigger, real users, real scale, real abuse.

**Moderation and abuse prevention.** It's an anonymous image feed with no rate limiting and no content moderation. Right now nothing stops someone from scripting a few thousand uploads in a row, or uploading something inappropriate. A real version would need at minimum a rate limit per IP on `/api/upload`, and realistically some kind of image moderation (even a basic NSFW classifier) before anything goes live.

**Aspect-ratio preserving crop.** Covered in the architecture notes, uploads currently get squashed into a 1080x1080 square instead of being center-cropped. Cheap fix, just didn't get to it.

**Object storage instead of local disk.** Uploaded images live in a Docker volume on the backend host, no redundancy, and it blocks horizontal scaling. In production I'd use S3 with pre-signed upload URLs so files go straight from the browser to S3, and serve the feed straight from S3 (behind a CDN like CloudFront) instead of proxying reads through the backend too.

**WebSocket fan-out across instances.** The live feed hub is a single in-memory map of connected clients. It works great with one backend instance and falls apart with more than one, clients on instance A never hear about uploads that land on instance B. Would need Redis pub/sub or similar to broadcast across instances.

**Cache invalidation is approximate.** Explained in the architecture notes too, uploading only busts the cache for the default (untagged, first-page) view. Filtered/paginated views just age out of the 30 second TTL on their own. Fine for now, but a pattern-based Redis invalidation (or a version counter baked into the cache key) would make it exact instead of eventually consistent.

**No thumbnails.** The feed grid shows the same full 1080x1080 JPEG that gets served in the lightbox view. A dedicated smaller thumbnail (say 300x300) for the grid would cut bandwidth noticeably, especially on the "load more" pagination where a dozen images load at once.

**Duplicate detection.** Nothing stops the same image being uploaded five times under different titles. A perceptual hash or even a simple content hash check on upload would catch exact duplicates at least.

**Observability.** Logging is just Chi's request logger middleware. No metrics (upload rate, cache hit ratio, WS connection count) and no tracing. Would want at least basic Prometheus metrics before running this anywhere real.

**No way to report or remove content.** Once something's uploaded, there's no delete endpoint and no reporting flow. For an anonymous feed that's a fairly important piece missing. As the task's description only mentioned handling upload for now so I did that, but if I wanted to improve this product there would be a delete endpoint available to moderators and a reporting workflow available to users to flag any disallowed content.
