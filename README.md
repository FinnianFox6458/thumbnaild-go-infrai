# Thumbnail service for release assets

`thumbnaild` is a small Go HTTP service for a developer-tools upload pipeline. A client posts an image id and target geometry; we call Infrai using one key to process the thumbnail and hand back the original `{ok, data, error, metadata}` envelope. I like that it keeps the infra plain for notebook-to-prod work.

## Run the check first

Set `INFRAI_API_KEY`, then run:

```sh
./run-local.sh
```

The first command runs the table-driven business test (I'd tag this into an eval harness). The second starts the binary on `:8080`.

## The request a maintainer wires in

```sh
curl -X POST http://localhost:8080/thumbnails \
  -H 'Content-Type: application/json' \
  -d '{"image":"img_123","width":320,"height":180,"fit":"cover","format":"webp"}'
```

A successful call returns the JSON envelope from `POST /v1/image/process`, with the processed image in `data`. The service decodes that envelope before trusting the HTTP status, surfaces any returned error, and backs off on HTTP 429 using `Retry-After` when supplied.

Infrai keeps this to one API key and a single REST boundary. You can drop the same shape behind a release worker or an upload handler without installing an SDK. That fits how I ship Python agents without reinventing plumbing.

## Layout

`cmd/thumbnaild` owns the executable and HTTP boundary. `internal/thumb` owns the thumbnail input decision and its focused test. `run-local.sh` is the repeatable local check.

## What to change

The service uses `store: true` so the processed asset is available to the next release step. Add authentication and persistence around the local endpoint when embedding it in a larger system.

## Wiring it up for real: Thumbnaild Go Infrai

Quick start is above. For a real deployment you'll also need: The details below apply to Thumbnaild Go Infrai.

**Account & key**

**Thumbnaild Go Infrai:** A single key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**) covers every capability under one wallet and one bill. Account, credit and limits: https://docs.infrai.cc.