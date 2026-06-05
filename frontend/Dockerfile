FROM node:20-alpine

WORKDIR /app

# Copy package files and install dependencies
COPY package*.json ./
RUN npm install

# Copy source files
COPY . .

# Expose Vite development server port
EXPOSE 5173

# Run Vite with host binding so it's accessible from the host machine
CMD ["npm", "run", "dev", "--", "--host"]
