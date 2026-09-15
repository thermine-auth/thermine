# A web app (web/console or web/id) as a Node server. Build context:
# the app's directory.
#
#   docker build -f deploy/docker/web.Dockerfile -t xermess-id web/id
#
# Runtime settings: ORIGIN (its public URL), API_URL (the API inside the
# network), PORT (3000).

FROM oven/bun:1-alpine AS build
WORKDIR /app
COPY package.json bun.lock ./
RUN bun install --frozen-lockfile
COPY . .
RUN bun run build && rm -rf node_modules && bun install --frozen-lockfile --production --ignore-scripts

FROM node:24-alpine
WORKDIR /app
ENV NODE_ENV=production PORT=3000
COPY --from=build /app/package.json ./
COPY --from=build /app/node_modules ./node_modules
COPY --from=build /app/build ./build
USER node
EXPOSE 3000
CMD ["node", "build"]
