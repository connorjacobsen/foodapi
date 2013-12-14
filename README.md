foodapi
=======

A sample API server written in Go.

## Info
This sample application will expose a single resource, foods, under the `/foods` URI. It will support:

- `GET /foods` : list all available foods
- `GET /foods/:id` : fetch a specific food
- `POST /foods` : create a food
- `PUT /foods/:id` : update a food
- `DELETE /foods/:id` : delete a food

Responses can be requested in JSON, XML, or plain text -- depending on the endpoint's extension.
*JSON is the default.*

The application uses an in-memory "database."
