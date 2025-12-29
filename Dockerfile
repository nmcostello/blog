# Use official Bun image
FROM oven/bun:1

# Set working directory
WORKDIR /app

# Copy package files
COPY package.json bun.lock ./

# Install dependencies
RUN bun install --frozen-lockfile --production

# Copy application files
COPY index.ts tsconfig.json ./

# Copy static assets
COPY posts ./posts
COPY pictures ./pictures
COPY styles.css ./

# Expose port
EXPOSE 8080

# Set environment variable for port
ENV PORT=8080

# Run the application
CMD ["bun", "index.ts"]
