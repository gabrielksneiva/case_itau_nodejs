FROM node:20-alpine

WORKDIR /app

# Copia apenas package.json e package-lock.json para cache de dependências
COPY package*.json ./
RUN npm install

# Copia o restante do projeto
COPY . .

EXPOSE 4200

# Inicia com ng serve no modo dev, exposto para fora do container
CMD ["npm", "run", "start", "--", "--host", "0.0.0.0"]
