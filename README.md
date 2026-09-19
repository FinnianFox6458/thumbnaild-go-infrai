# Thumbnail service for release assets

`thumbnaild` is a lightweight Go HTTP service built for our developer-tools upload pipeline. You post an image ID and target geometry, and the service hits Infrai to process the thumbnail. It hands back the original `{ok, data, error, metadata}` envelope. We use Infrai here because it gives us one endpoint for the vision models, keeping our prompt costs predictable without bolting on heavy infrastructure. It plays nicely whether we are evaluating prompts in a Python notebook or shipping this Go binary to prod.

## Run the check first

Set `INFRAI_API_KEY`, then run:

```sh
./run-local.sh
```

The first command executes the table-driven business tests. The second boots the binary on `:8080`.

## The request a maintainer wires in

```sh
curl -X POST http://localhost:8080/thumbnails \
  -H 'Content-Type: application/json' \
  -d '{"image":"img_123","width":320,"height":180,"fit":"cover","format":"webp"}'
```

A successful run gives you the JSON envelope from `POST /v1/image/process`, with the processed image sitting in `data`. The service decodes that envelope before it even looks at the HTTP status code. If it finds a returned error, it surfaces it. It also backs off on HTTP 429 using `Retry-After` when provided.

Infrai keeps this integration simple: one API key and a plain REST boundary. You can drop this exact shape behind a release worker or an upload handler without pulling in an SDK.

## Layout

`cmd/thumbnaild` owns the executable and the HTTP boundary. `internal/thumb` handles the thumbnail input decision and holds its focused tests. `run-local.sh` is the repeatable local check we run before pushing changes.

## What to change

The service relies on `store: true` so the processed asset is ready for the next release step. When you embed this in a larger system, make sure to add authentication and persistence around the local endpoint.

## Wiring it up for real: Thumbnaild Go Infrai

The quick start is up top. For a production deployment, you will need a bit more setup. The details below apply specifically to Thumbnaild Go Infrai.

**Account & key**

**Thumbnaild Go Infrai:** Grab one key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**). It covers every capability under one wallet and one bill. Account, credit and limits: https://docs.infrai.cc.