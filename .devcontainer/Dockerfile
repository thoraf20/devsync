FROM mcr.microsoft.com/devcontainers/go:1.21

# Optional: Node + Yarn if you want to support NestJS as well
RUN apt-get update && \
  apt-get install -y curl gnupg && \
  curl -fsSL https://deb.nodesource.com/setup_18.x | bash - && \
  apt-get install -y nodejs && \
  npm install -g yarn
