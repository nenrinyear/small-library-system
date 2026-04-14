FROM node:22-alpine

WORKDIR /app

COPY package.json package-lock.json* ./
RUN npm ci

COPY src ./src
COPY public ./public
COPY next.config.mjs .
COPY tsconfig.json .
COPY next-env.d.ts .

EXPOSE 3000

CMD ["npm", "run", "dev"]
