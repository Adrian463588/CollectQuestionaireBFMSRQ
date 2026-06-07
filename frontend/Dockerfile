# Stage 1: Build the Next.js application
FROM node:20-alpine AS builder

WORKDIR /app

# Cache dependencies
COPY package*.json ./
RUN npm ci

# Copy source code
COPY . .

# Set dynamic build arguments for public environment variables
ARG NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL

# Build the Next.js application
RUN npm run build

# Stage 2: Runtime image
FROM node:20-alpine AS runner

WORKDIR /app

# Install security patches/updates and set environment
ENV NODE_ENV=production

# Create non-root user and group for security
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

# Copy essential files
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules

# Ensure proper file permissions
RUN chown -R nextjs:nodejs /app

# Switch to non-root user
USER nextjs

# Expose port
EXPOSE 3000

# Start Next.js server
CMD ["npm", "start"]
