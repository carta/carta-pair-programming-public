FROM node:22-trixie-slim

WORKDIR /app

RUN corepack enable && corepack prepare pnpm@11.27.1 --activate

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile

COPY . .

EXPOSE 5173

CMD ["pnpm", "run", "dev", "--host", "0.0.0.0"]
