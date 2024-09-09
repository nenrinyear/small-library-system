FROM node:18-alpine AS base

# Step 1. Rebuild the source code only when needed
FROM base AS builder

WORKDIR /app

# Install dependencies based on the preferred package manager
COPY package.json yarn.lock* package-lock.json* pnpm-lock.yaml* ./
# Omit --production flag for TypeScript devDependencies
RUN \
  if [ -f yarn.lock ]; then yarn --frozen-lockfile; \
  elif [ -f package-lock.json ]; then npm ci; \
  elif [ -f pnpm-lock.yaml ]; then corepack enable pnpm && pnpm i; \
  # Allow install without lockfile, so example works even without Node.js installed locally
  else echo "Warning: Lockfile not found. It is recommended to commit lockfiles to version control." && yarn install; \
  fi

COPY src ./src
COPY public ./public
COPY next.config.mjs .
COPY tsconfig.json .
COPY tailwind.config.ts .
COPY postcss.config.mjs .

# Environment variables must be present at build time
# https://github.com/vercel/next.js/discussions/14030
ARG NEXT_PUBLIC_SYSTEM_NAME
ARG NEXT_PUBLIC_SYSTEM_NAME_SHORT
ARG NEXT_PUBLIC_SYSTEM_MANAGE_ITEM
ARG NEXT_PUBLIC_SYSTEM_HOST
ARG NEXT_PUBLIC_SYSTEM_DESCRIPTION
ARG NEXT_PUBLIC_MAX_HISTORY_COUNT
ARG MYSQL_HOST
ARG MYSQL_PORT
ARG MYSQL_DATABASE
ARG MYSQL_USER
ARG MYSQL_PASSWORD
ARG MYSQL_DATABASE_URL
ENV NEXT_PUBLIC_SYSTEM_NAME=${NEXT_PUBLIC_SYSTEM_NAME}
ENV NEXT_PUBLIC_SYSTEM_NAME_SHORT=${NEXT_PUBLIC_SYSTEM_NAME_SHORT}
ENV NEXT_PUBLIC_SYSTEM_MANAGE_ITEM=${NEXT_PUBLIC_SYSTEM_MANAGE_ITEM}
ENV NEXT_PUBLIC_SYSTEM_HOST=${NEXT_PUBLIC_SYSTEM_HOST}
ENV NEXT_PUBLIC_SYSTEM_DESCRIPTION=${NEXT_PUBLIC_SYSTEM_DESCRIPTION}
ENV NEXT_PUBLIC_MAX_HISTORY_COUNT=${NEXT_PUBLIC_MAX_HISTORY_COUNT}
ENV MYSQL_HOST=${MYSQL_HOST}
ENV MYSQL_PORT=${MYSQL_PORT}
ENV MYSQL_DATABASE=${MYSQL_DATABASE}
ENV MYSQL_USER=${MYSQL_USER}
ENV MYSQL_PASSWORD=${MYSQL_PASSWORD}
ENV MYSQL_DATABASE_URL=${MYSQL_DATABASE_URL}

# Next.js collects completely anonymous telemetry data about general usage. Learn more here: https://nextjs.org/telemetry
# Uncomment the following line to disable telemetry at build time
# ENV NEXT_TELEMETRY_DISABLED 1

# Build Next.js based on the preferred package manager
RUN \
  if [ -f yarn.lock ]; then yarn build; \
  elif [ -f package-lock.json ]; then npm run build; \
  elif [ -f pnpm-lock.yaml ]; then pnpm build; \
  else npm run build; \
  fi

# Note: It is not necessary to add an intermediate step that does a full copy of `node_modules` here

# Step 2. Production image, copy all the files and run next
FROM base AS runner

WORKDIR /app

# Don't run production as root
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs
USER nextjs

COPY --from=builder /app/public ./public

# Automatically leverage output traces to reduce image size
# https://nextjs.org/docs/advanced-features/output-file-tracing
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

# Environment variables must be redefined at run time
ARG NEXT_PUBLIC_SYSTEM_NAME
ARG NEXT_PUBLIC_SYSTEM_NAME_SHORT
ARG NEXT_PUBLIC_SYSTEM_MANAGE_ITEM
ARG NEXT_PUBLIC_SYSTEM_HOST
ARG NEXT_PUBLIC_SYSTEM_DESCRIPTION
ARG NEXT_PUBLIC_MAX_HISTORY_COUNT
ARG MYSQL_HOST
ARG MYSQL_PORT
ARG MYSQL_DATABASE
ARG MYSQL_USER
ARG MYSQL_PASSWORD
ENV NEXT_PUBLIC_SYSTEM_NAME=${NEXT_PUBLIC_SYSTEM_NAME}
ENV NEXT_PUBLIC_SYSTEM_NAME_SHORT=${NEXT_PUBLIC_SYSTEM_NAME_SHORT}
ENV NEXT_PUBLIC_SYSTEM_MANAGE_ITEM=${NEXT_PUBLIC_SYSTEM_MANAGE_ITEM}
ENV NEXT_PUBLIC_SYSTEM_HOST=${NEXT_PUBLIC_SYSTEM_HOST}
ENV NEXT_PUBLIC_SYSTEM_DESCRIPTION=${NEXT_PUBLIC_SYSTEM_DESCRIPTION}
ENV NEXT_PUBLIC_MAX_HISTORY_COUNT=${NEXT_PUBLIC_MAX_HISTORY_COUNT}
ENV MYSQL_HOST=${MYSQL_HOST}
ENV MYSQL_PORT=${MYSQL_PORT}
ENV MYSQL_DATABASE=${MYSQL_DATABASE}
ENV MYSQL_USER=${MYSQL_USER}
ENV MYSQL_PASSWORD=${MYSQL_PASSWORD}

# Uncomment the following line to disable telemetry at run time
# ENV NEXT_TELEMETRY_DISABLED 1

# Note: Don't expose ports here, Compose will handle that for us

CMD ["node", "server.js"]
