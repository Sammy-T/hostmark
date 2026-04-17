# hostmark

A knowledge-base and notes app supporting markdown text.

## Features

- Role based access
- Member-only markdown files
  - Local filesystem access to the app's markdown files (depending on setup)
- Personal markdown notes
  - Configurable privacy
  - Custom tags

<img width="1920" height="909" alt="hm-screen-01" src="https://github.com/user-attachments/assets/337dcfc4-72fc-4b36-be9b-5a9ad162d99b" /><br>
<img width="1920" height="909" alt="hm-screen-02" src="https://github.com/user-attachments/assets/818fa282-81bc-4764-be34-be7b23103938" /><br>
<img width="1920" height="909" alt="hm-screen-03" src="https://github.com/user-attachments/assets/75c3d3e2-2bd9-46b2-b4a4-9c97d3ef454f" /><br>

## Getting Started

### Running with Docker

> [!IMPORTANT]
> If you want local filesystem access to the app's markdown files,<br>
> Create a `.files` subdirectory inside your current working directory **before** creating and running the app container.

Create and run the container:

```bash
# With local filesystem access
docker run --restart=unless-stopped -d -p 4000:3000 --name hostmark -v .files:/app/.files:rw -v hm-data:/app/.data sammytd/hostmark
```

```bash
# Without local filesystem access
docker run --restart=unless-stopped -d -p 4000:3000 --name hostmark -v hm-files:/app/.files -v hm-data:/app/.data sammytd/hostmark
```

Then visit `http://localhost:4000`.

#### Docker run options

| Option | Description |
| --- | --- |
| `--restart` | The container restart policy. |
| `-d, --detach` | Run the container in the background. |
| `-p, --publish` | Publish the container's port to the host. The host port can be whichever available port you want.<br>`<host-port>:<container-port>` |
| `--name` | The container name (optional) |
| `-v, --volume` | A container volume for persistent data. |

See <https://docs.docker.com/reference/cli/docker/container/run/> for full reference.

## Development

### Run the dev server

Requirements:

- [Node.js](https://nodejs.org/)
- [pnpm](https://pnpm.io/)
- [Go](https://go.dev/)

#### Get go dependencies

```bash
go get ./...
```

#### Run the server

```bash
go run -tags dev .
```

Then visit `http://localhost:3000`.

### Build and run the Docker container

Requires [docker](https://www.docker.com/)

#### Build the image

```bash
docker build -t my/hostmark .
```

#### Create and run a container from the image

```bash
docker run -d -p 4000:3000 --name my-hostmark -v .files:/app/.files:rw -v my-hm-data:/app/.data my/hostmark
```

Then visit your host machine's `http://localhost:4000`.
